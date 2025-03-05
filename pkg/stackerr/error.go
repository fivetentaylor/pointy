package stackerr

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const Header = "Stacktrace:\n"

func New(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), Header) {
		return err
	}

	return fmt.Errorf("%w\n%s%s", err, Header, string(debug.Stack()))
}

func Wrap(err error) error {
	return New(err)
}

func Errorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	if strings.Contains(err.Error(), Header) {
		return err
	}

	return fmt.Errorf("%w\n%s%s", err, Header, string(debug.Stack()))
}
