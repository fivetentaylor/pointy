package rogue

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/storage/s3"
)

type Logger struct {
	ctx               context.Context
	prefix            string
	docID             string
	s3                *s3.S3
	buffer            *strings.Builder
	active            bool
	passthroughActive bool

	passthrough *log.Logger
	mu          sync.Mutex
}

func NewLogger(passthrough *log.Logger, s3 *s3.S3, docID string, prefix string) *Logger {
	l := &Logger{
		ctx:               context.Background(),
		prefix:            prefix,
		docID:             docID,
		s3:                s3,
		buffer:            &strings.Builder{},
		passthroughActive: true,
		active:            true,
		passthrough:       passthrough,
	}

	go l.poll()

	return l
}

func (l *Logger) Deactivate() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.active = false
}

func (l *Logger) Activate() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.active = true
}

func (l *Logger) PassthroughActive() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.passthroughActive
}

func (l *Logger) DeactivatePassthrough() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.passthroughActive = false
}

func (l *Logger) ActivatePassthrough() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.passthroughActive = true
}

func (l *Logger) poll() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			err := l.Flush()
			if err != nil {
				l.passthrough.Error(fmt.Sprintf("failed to flush rogue doc log: %s", err))
			}
		}
	}
}

type LogEntry struct {
	Timestamp time.Time
	UserID    string
	SessionID string
	Message   string
}

var idRegex = regexp.MustCompile(`(.*)\[([A-Za-z0-9]{5})\|([A-Za-z0-9]{5})\](.*)`)

func (l *Logger) Read() ([]LogEntry, error) {
	btss, err := l.s3.GetAllObjects(l.s3.Bucket, l.ContainingKey())
	if err != nil {
		return nil, err
	}

	lines := []string{}
	for _, bts := range btss {
		lines = append(lines, strings.Split(string(bts), "\n")...)
	}

	entries := make([]LogEntry, 0, len(lines))
	for _, l := range lines {
		if l == "" {
			continue
		}
		fmt.Println(l)
		parts := strings.SplitN(l, "@", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid log entry: %s", l)
		}
		timestamp, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid log entry timestamp: %s", l)
		}
		remainder := parts[1]
		entry := LogEntry{
			Timestamp: timestamp,
			Message:   remainder,
		}

		idMatch := idRegex.FindStringSubmatch(remainder)
		if len(idMatch) > 2 {
			entry.UserID = idMatch[2]
			entry.SessionID = idMatch[3]
			entry.Message = idMatch[1] + idMatch[4]
		}

		entries = append(entries, entry)
	}

	sortByTimestamp(entries)

	return entries, nil
}

func sortByTimestamp(entries []LogEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
}

func (l *Logger) Info(msg string) {
	l.writeString("INFO", msg)
	if l.passthroughActive {
		l.passthrough.Info(msg)
	}
}
func (l *Logger) Warn(msg string) {
	l.writeString("WARN", msg)
	if l.passthroughActive {
		l.passthrough.Warn(msg)
	}
}
func (l *Logger) Error(msg string) {
	l.writeString("ERROR", msg)
	if l.passthroughActive {
		l.passthrough.Error(msg)
	}
}
func (l *Logger) Debug(msg string) {
	l.writeString("DEBUG", msg)
	if l.passthroughActive {
		l.passthrough.Debug(msg)
	}
}
func (l *Logger) Fatal(msg string) {
	l.writeString("FATAL", msg)
	if l.passthroughActive {
		l.passthrough.Fatal(msg)
	}
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.writeString("INFO", fmt.Sprintf(format, args...))
	if l.passthroughActive {
		l.passthrough.Infof(format, args...)
	}
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.writeString("WARN", fmt.Sprintf(format, args...))
	if l.passthroughActive {
		l.passthrough.Warnf(format, args...)
	}
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.writeString("ERROR", fmt.Sprintf(format, args...))
	if l.passthroughActive {
		l.passthrough.Errorf(format, args...)
	}
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.writeString("DEBUG", fmt.Sprintf(format, args...))
	if l.passthroughActive {
		l.passthrough.Debugf(format, args...)
	}
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.writeString("FATAL", fmt.Sprintf(format, args...))
	if l.passthroughActive {
		l.passthrough.Fatalf(format, args...)
	}
}

func (l *Logger) writeString(level, msg string) {
	if !l.active {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.buffer.WriteString(
		fmt.Sprintf("%s@%s [%s]: %s\n", time.Now().Format(time.RFC3339Nano), level, l.prefix, msg),
	)
	if err != nil {
		l.passthrough.Error(fmt.Sprintf("failed to write rogue doc log: %s", err))
	}
}

func (l *Logger) ContainingKey() string {
	return fmt.Sprintf(constants.LogsPrefix, l.docID)
}

func (l *Logger) Key() string {
	return fmt.Sprintf("%s/%s.log", l.ContainingKey(), time.Now().Format(time.RFC3339Nano))
}

func (l *Logger) Flush() error {
	if l.buffer.Len() == 0 {
		return nil
	}

	if !l.active {
		l.passthrough.Infof("skipping flushing %d bytes to s3: %s", l.buffer.Len(), l.Key())
		return nil
	}

	l.passthrough.Debugf("flushing %d bytes to s3: %s", l.buffer.Len(), l.Key())

	l.mu.Lock()
	defer l.mu.Unlock()

	err := l.s3.PutObject(l.s3.Bucket, l.Key(), "text/plain", []byte(l.buffer.String()))
	if err != nil {
		return err
	}
	l.buffer.Reset()

	return nil
}

func (l *Logger) Close() {
	l.ctx.Done()
	err := l.Flush()
	if err != nil {
		l.passthrough.Error(fmt.Sprintf("failed to flush rogue doc log: %s", err))
	}
}
