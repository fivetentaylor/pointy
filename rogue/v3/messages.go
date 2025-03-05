package v3

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

const InsertOpType = 0
const DeleteOpType = 1
const FormatOpType = 2
const SnapshotOpType = 3
const UndoOpType = 4       // DEPRECATED
const UndoFormatOpType = 5 // DEPRECATED
const MultiOpType = 6
const RewindOpType = 7
const ShowOpType = 8

type InsertOp struct {
	ID       ID
	Text     string
	ParentID ID
	Side     Side
}

type ID struct {
	Author string
	Seq    int
}

func (id ID) lessThan(other ID) bool {
	return id.Seq < other.Seq || (id.Seq == other.Seq && id.Author < other.Author)
}

func (id ID) String() string {
	return fmt.Sprintf("%s_%d", id.Author, id.Seq)
}

func (id ID) AsJS() []interface{} {
	return []interface{}{id.Author, id.Seq}
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{id.Author, id.Seq})
}

func (id *ID) UnmarshalJSON(b []byte) error {
	var temp []interface{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	if len(temp) != 2 {
		return fmt.Errorf("expected 2 elements in the array, got %d", len(temp))
	}

	author, ok := temp[0].(string)
	if !ok {
		return fmt.Errorf("expected a string for Author, got %T", temp[0])
	}

	seqFloat, ok := temp[1].(float64) // JSON numbers are unmarshaled as float64
	if !ok {
		return fmt.Errorf("expected a number for Seq, got %T", temp[1])
	}

	id.Author = author
	id.Seq = int(seqFloat)

	return nil
}

type Side int

func (s Side) String() string {
	switch s {
	case Left:
		return "left"
	case Root:
		return "root"
	case Right:
		return "right"
	default:
		return "unknown"
	}
}

const (
	Left  Side = -1
	Root  Side = 0
	Right Side = 1
)

type ContentAddress struct {
	StartID ID             `json:"startID"`
	EndID   ID             `json:"endID"`
	MaxIDs  map[string]int `json:"maxIDs"`
}

func (ca ContentAddress) AsArr() []interface{} {
	maxIDs := make([]interface{}, 0, len(ca.MaxIDs))
	for author, seq := range ca.MaxIDs {
		maxIDs = append(maxIDs, []interface{}{author, seq})
	}

	return []interface{}{ca.StartID.AsJS(), ca.EndID.AsJS(), maxIDs}
}

func (a ContentAddress) MarshalJSON() ([]byte, error) {
	maxIDs := make([]ID, 0, len(a.MaxIDs))
	for k, v := range a.MaxIDs {
		maxIDs = append(maxIDs, ID{k, v})
	}
	sort.Slice(maxIDs, func(i, j int) bool {
		return maxIDs[i].Seq < maxIDs[j].Seq
	})

	return json.Marshal([]interface{}{a.StartID, a.EndID, maxIDs})
}

type _oldAddress struct {
	StartID []interface{}            `json:"startID"`
	EndID   []interface{}            `json:"endID"`
	MaxIDs  []map[string]interface{} `json:"maxIDs"`
}

func (a *ContentAddress) _parseOldAddress(data []byte) error {
	var oldAddr _oldAddress
	err := json.Unmarshal(data, &oldAddr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal old address: %v", err)
	}

	// Parse StartID
	if len(oldAddr.StartID) != 2 {
		return fmt.Errorf("invalid startID format")
	}
	a.StartID.Author = oldAddr.StartID[0].(string)
	a.StartID.Seq = int(oldAddr.StartID[1].(float64))

	// Parse EndID
	if len(oldAddr.EndID) != 2 {
		return fmt.Errorf("invalid endID format")
	}
	a.EndID.Author = oldAddr.EndID[0].(string)
	a.EndID.Seq = int(oldAddr.EndID[1].(float64))

	// Parse MaxIDs
	a.MaxIDs = make(map[string]int)
	for _, maxID := range oldAddr.MaxIDs {
		key := maxID["key"].(string)
		value := int(maxID["value"].(float64))
		a.MaxIDs[key] = value
	}

	return nil
}

func (a *ContentAddress) UnmarshalJSON(data []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(data, &serializedOp); err != nil {
		err2 := a._parseOldAddress(data)
		if err2 == nil {
			return nil
		}

		return err
	}

	if len(serializedOp) != 3 {
		return fmt.Errorf("invalid number of fields in serialized ContentAddress")
	}

	if err := json.Unmarshal(serializedOp[0], &a.StartID); err != nil {
		return err
	}

	if err := json.Unmarshal(serializedOp[1], &a.EndID); err != nil {
		return err
	}

	maxIDs := make([]ID, 0)
	if err := json.Unmarshal(serializedOp[2], &maxIDs); err != nil {
		return err
	}

	a.MaxIDs = make(map[string]int)
	for _, id := range maxIDs {
		a.MaxIDs[id.Author] = id.Seq
	}

	return nil
}

func (op InsertOp) String() string {
	return fmt.Sprintf("InsertOp{ID: %s, Text: %q, ParentID: %s, Side: %s, UTF-16: %s}", op.ID, op.Text, op.ParentID, op.Side, Uint16SliceToHexString(StrToUint16(op.Text)))
}
func (op InsertOp) GetID() ID { return op.ID }
func (op InsertOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       op.ID.String(),
		"text":     op.Text,
		"parentID": op.ParentID.String(),
		"side":     op.Side.String(),
	}
}
func (op InsertOp) MarshalJSON() ([]byte, error) {
	serializedOp := []interface{}{
		InsertOpType,
		op.ID,
		op.Text,
		op.ParentID,
		op.Side,
	}

	return json.Marshal(serializedOp)
}
func (op *InsertOp) UnmarshalJSON(b []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(b, &serializedOp); err != nil {
		return err
	}
	// [0,["8cd8-dk",3],"H",["q",1],1]

	if len(serializedOp) != 5 {
		return fmt.Errorf("invalid number of fields in serialized InsertOp")
	}

	if err := json.Unmarshal(serializedOp[1], &op.ID); err != nil {
		return err
	}
	if err := json.Unmarshal(serializedOp[2], &op.Text); err != nil {
		return err
	}
	if err := json.Unmarshal(serializedOp[3], &op.ParentID); err != nil {
		return err
	}
	if err := json.Unmarshal(serializedOp[4], &op.Side); err != nil {
		return err
	}

	return nil
}

// [0,["8cd8-dk",3],"H",["q",1],1]
func (op InsertOp) AsArr() []interface{} {
	return []interface{}{
		InsertOpType,
		op.ID.AsJS(),
		op.Text,
		op.ParentID.AsJS(),
		int(op.Side),
	}
}

type DeleteOp struct {
	ID         ID
	TargetID   ID
	SpanLength int
}

func (op DeleteOp) String() string {
	return fmt.Sprintf("DeleteOp{ID: %s, TargetID: %s, SpanLength: %d}", op.ID, op.TargetID, op.SpanLength)
}
func (op DeleteOp) GetID() ID { return op.ID }
func (op DeleteOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         op.ID.String(),
		"targetID":   op.TargetID.String(),
		"spanLength": op.SpanLength,
	}
}
func (op DeleteOp) AsArr() []interface{} {
	return []interface{}{
		DeleteOpType,
		op.ID.AsJS(),
		op.TargetID.AsJS(),
		op.SpanLength,
	}
}
func (op DeleteOp) MarshalJSON() ([]byte, error) {
	serializedOp := []interface{}{
		DeleteOpType,
		op.ID,
		op.TargetID,
		op.SpanLength,
	}

	return json.Marshal(serializedOp)
}
func (op *DeleteOp) UnmarshalJSON(b []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(b, &serializedOp); err != nil {
		return err
	}

	if len(serializedOp) < 3 || len(serializedOp) > 4 {
		return fmt.Errorf("invalid number of fields in serialized DeleteOp")
	}

	op.SpanLength = 1

	if err := json.Unmarshal(serializedOp[1], &op.ID); err != nil {
		return err
	}
	if err := json.Unmarshal(serializedOp[2], &op.TargetID); err != nil {
		return err
	}

	if len(serializedOp) == 4 {
		if err := json.Unmarshal(serializedOp[3], &op.SpanLength); err != nil {
			return err
		}
	}

	return nil
}

type ShowOp struct {
	ID         ID
	TargetID   ID
	SpanLength int
}

func (op ShowOp) String() string {
	return fmt.Sprintf("ShowOp{ID: %s, TargetID: %s, SpanLength: %d}", op.ID, op.TargetID, op.SpanLength)
}
func (op ShowOp) GetID() ID { return op.ID }
func (op ShowOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         op.ID.String(),
		"targetID":   op.TargetID.String(),
		"spanLength": op.SpanLength,
	}
}
func (op ShowOp) AsArr() []interface{} {
	return []interface{}{
		ShowOpType,
		op.ID.AsJS(),
		op.TargetID.AsJS(),
		op.SpanLength,
	}
}
func (op ShowOp) MarshalJSON() ([]byte, error) {
	serializedOp := []interface{}{
		ShowOpType,
		op.ID,
		op.TargetID,
		op.SpanLength,
	}

	return json.Marshal(serializedOp)
}
func (op *ShowOp) UnmarshalJSON(b []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(b, &serializedOp); err != nil {
		return err
	}

	if len(serializedOp) != 4 {
		return fmt.Errorf("invalid number of fields in serialized ShowOp")
	}

	if err := json.Unmarshal(serializedOp[1], &op.ID); err != nil {
		return err
	}
	if err := json.Unmarshal(serializedOp[2], &op.TargetID); err != nil {
		return err
	}

	if err := json.Unmarshal(serializedOp[3], &op.SpanLength); err != nil {
		return err
	}

	return nil
}

type RewindOp struct {
	ID          ID
	Address     ContentAddress
	UndoAddress ContentAddress
}

func (op RewindOp) String() string {
	return fmt.Sprintf("RewindOp{ID: %s, Addresses: %v}", op.ID, op.Address)
}
func (op RewindOp) GetID() ID { return op.ID }
func (op RewindOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          op.ID.String(),
		"address":     op.Address,
		"undoAddress": op.UndoAddress,
	}
}
func (op RewindOp) AsArr() []interface{} {
	return []interface{}{
		RewindOpType,
		op.ID.AsJS(),
		op.Address.AsArr(),
		op.UndoAddress.AsArr(),
	}
}
func (op RewindOp) MarshalJSON() ([]byte, error) {
	serializedOp := []interface{}{
		RewindOpType,
		op.ID,
		op.Address.AsArr(),
		op.UndoAddress.AsArr(),
	}

	return json.Marshal(serializedOp)
}
func (op *RewindOp) UnmarshalJSON(b []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(b, &serializedOp); err != nil {
		return err
	}

	if len(serializedOp) != 4 {
		return fmt.Errorf("invalid number of fields in serialized RewindOp")
	}

	if err := json.Unmarshal(serializedOp[1], &op.ID); err != nil {
		return err
	}

	if err := json.Unmarshal(serializedOp[2], &op.Address); err != nil {
		return err
	}

	if err := json.Unmarshal(serializedOp[3], &op.UndoAddress); err != nil {
		return err
	}

	return nil
}

type MultiOp struct {
	Mops []Op
}

func (op MultiOp) String() string {
	return fmt.Sprintf("MultiOp{Ops: %v}", op.Mops)
}

func (op MultiOp) GetID() ID {
	if len(op.Mops) == 0 {
		return NoID
	}

	id := op.Mops[0].GetID()
	for _, o := range op.Mops[1:] {
		if o.GetID().lessThan(id) {
			id = o.GetID()
		}
	}

	return id
}
func (op MultiOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"ops": op.Mops,
	}
}

func (op MultiOp) AsArr() []interface{} {
	ops := make([]interface{}, 0, len(op.Mops))
	for _, o := range op.Mops {
		ops = append(ops, o.AsArr())
	}

	return []interface{}{
		MultiOpType,
		op.GetID().AsJS(),
		ops,
	}
}

func (op MultiOp) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		MultiOpType,
		op.GetID(),
		op.Mops,
	})
}

// Hello World!
// {startIx: 0, endIx: 7, format: {"b": "true"}}
// {startIx: 5, endIx: 12, format: {"i": "true"}}
// <b>Hello</b><b><i> W</i></b><i>orld</i>
// **Hello_ W_**_orld_

func (op *MultiOp) UnmarshalJSON(b []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(b, &serializedOp); err != nil {
		return err
	}

	if len(serializedOp) != 3 {
		return fmt.Errorf("invalid number of fields in serialized MultiOp")
	}

	var msgs []Message
	if err := json.Unmarshal(serializedOp[2], &msgs); err != nil {
		return err
	}

	op.Mops = make([]Op, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Op != nil {
			op.Mops = append(op.Mops, msg.Op)
		}
	}

	return nil
}

type FormatOp struct {
	ID      ID
	StartID ID
	EndID   ID
	Format  FormatV3
}

func (op FormatOp) String() string {
	return fmt.Sprintf("FormatOp{ID: %s, StartID: %s, EndID: %s, Format: %s}", op.ID, op.StartID, op.EndID, op.Format)
}
func (op FormatOp) GetID() ID { return op.ID }
func (op FormatOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      op.ID.String(),
		"startID": op.StartID.String(),
		"endID":   op.EndID.String(),
		"format":  op.Format.AsMap(),
	}
}
func (op FormatOp) AsArr() []interface{} {
	return []interface{}{
		FormatOpType,
		op.ID.AsJS(),
		op.StartID.AsJS(),
		op.EndID.AsJS(),
		op.Format.AsMap(),
	}
}

func (op FormatOp) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		FormatOpType,
		op.ID,
		op.StartID,
		op.EndID,
		op.Format.AsMap(),
	})
}

func (op *FormatOp) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 5 {
		return fmt.Errorf("FormatOpV3 should have five elements")
	}

	var opType int
	if err := json.Unmarshal(raw[0], &opType); err != nil {
		return err
	}

	if opType != FormatOpType {
		return fmt.Errorf("invalid op type: %d", opType)
	}

	if err := json.Unmarshal(raw[1], &op.ID); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[2], &op.StartID); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[3], &op.EndID); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(raw[4], &m); err != nil {
		return err
	}

	f, err := MapToFormatV3(m)
	if err != nil {
		return err
	}

	op.Format = f

	return nil
}

type SerializedRogue struct {
	Version *string    `json:"version"`
	Ops     MessageOps `json:"ops"`
}

type SnapshotOp struct {
	ID       ID
	Snapshot *SerializedRogue
}

func (op SnapshotOp) String() string {
	return fmt.Sprintf("SnapshotOp{ID: %s, Snapshot: %v}", op.ID, op.Snapshot)
}
func (op SnapshotOp) GetID() ID { return op.ID }
func (op SnapshotOp) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       op.ID.String(),
		"snapshot": op.Snapshot,
	}
}
func (op SnapshotOp) AsArr() []interface{} {
	return []interface{}{
		SnapshotOpType,
		op.ID.AsJS(),
		// 	op.Snapshot.AsJS(), // TODO: fix this interface
	}
}
func (op SnapshotOp) MarshalJSON() ([]byte, error) {
	serializedOp := []interface{}{
		SnapshotOpType,
		op.Snapshot,
	}

	return json.Marshal(serializedOp)
}
func (op *SnapshotOp) UnmarshalJSON(b []byte) error {
	var serializedOp []json.RawMessage
	if err := json.Unmarshal(b, &serializedOp); err != nil {
		return err
	}

	if len(serializedOp) != 2 {
		return fmt.Errorf("invalid number of fields in serialized SnapshotOp")
	}

	if err := json.Unmarshal(serializedOp[1], &op.Snapshot); err != nil {
		return err
	}

	return nil
}

type Op interface {
	GetID() ID
	MarshalJSON() ([]byte, error)
	AsMap() map[string]interface{}
	AsArr() []interface{}
	fmt.Stringer
}

type Message struct {
	Op Op
}

type MessageOps []Op

func (ops MessageOps) AsJS() []interface{} {
	var arr []interface{}
	for _, op := range ops {
		arr = append(arr, op.AsArr())
	}
	return arr
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Op)
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var op []json.RawMessage
	if err := json.Unmarshal(data, &op); err != nil {
		return err
	}

	var msgType int
	if err := json.Unmarshal(op[0], &msgType); err != nil {
		return err
	}

	switch msgType {
	case InsertOpType:
		var op InsertOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	case DeleteOpType:
		var op DeleteOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	case ShowOpType:
		var op ShowOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	case FormatOpType:
		var op FormatOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	case RewindOpType:
		var op RewindOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	case MultiOpType:
		var op MultiOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	case SnapshotOpType:
		var op SnapshotOp
		if err := json.Unmarshal(data, &op); err != nil {
			return err
		}
		m.Op = op
	default:
		m.Op = nil
		log.Warnf("invalid message type: %d", msgType)
	}

	return nil
}

func (ops *MessageOps) String() string {
	s := strings.Builder{}
	for _, op := range *ops {
		s.WriteString(fmt.Sprintf("%#v\n", op))
	}
	return s.String()
}

func (ops *MessageOps) UnmarshalJSON(data []byte) error {
	var y []json.RawMessage
	if err := json.Unmarshal(data, &y); err != nil {
		return err
	}

	out := make(MessageOps, 0, len(y))
	for _, op := range y {
		var msg Message
		err := json.Unmarshal(op, &msg)
		if err != nil {
			return err
		}

		if msg.Op != nil {
			out = append(out, msg.Op)
		}
	}

	*ops = out

	return nil
}

func handleListValue(i interface{}) (int, error) {
	switch v := i.(type) {
	case string:
		if v == "" || v == "null" {
			return -1, nil
		}

		li, err := strconv.Atoi(v)
		if err != nil {
			return -1, err
		}

		if 0 <= li && li <= 6 {
			return li, nil
		} else {
			return 0, nil
			// return li, fmt.Errorf("bullet list value must be between 0 and 6")
		}
	case float64:
		li := int(v)

		if 0 <= li && li <= 6 {
			return li, nil
		} else {
			return 0, nil
			// return li, fmt.Errorf("bullet list value must be between 0 and 6")
		}
	case int:
		if 0 <= v && v <= 6 {
			return v, nil
		} else {
			return 0, nil
			// return v, fmt.Errorf("bullet list value must be between 0 and 6")
		}
	default:
		return -1, fmt.Errorf("bullet list value must be a string or a number")
	}
}

func handleHeaderValue(i interface{}) (int, error) {
	switch v := i.(type) {
	case string:
		if v == "" || v == "null" {
			return -1, nil
		}

		li, err := strconv.Atoi(v)
		if err != nil {
			return -1, err
		}

		if 1 <= li && li <= 6 {
			return li, nil
		} else {
			return li, fmt.Errorf("header value must be between 1 and 6")
		}
	case float64:
		i := int(v)

		if 1 <= i && i <= 6 {
			return i, nil
		} else {
			return i, fmt.Errorf("header value must be between 1 and 6")
		}
	case int:
		if 1 <= v && v <= 6 {
			return v, nil
		} else {
			return v, fmt.Errorf("header value must be between 1 and 6")
		}
	default:
		return -1, fmt.Errorf("header value must be a string or a number not %T", v)
	}
}

func interfaceToString(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func handleImgValue(m map[string]interface{}) FormatV3Image {
	return FormatV3Image{
		Src:    interfaceToString(m["img"]),
		Alt:    interfaceToString(m["alt"]),
		Width:  interfaceToString(m["width"]),
		Height: interfaceToString(m["height"]),
	}
}

var spanKeys = map[string]string{
	"bold":      "b",
	"italic":    "i",
	"underline": "u",
	"strike":    "s",
	"link":      "a",
}

func handleSpanFormat(m map[string]interface{}) (FormatV3Span, error) {
	span := FormatV3Span{}
	for k, v := range m {
		if key, ok := spanKeys[k]; ok {
			k = key
		}

		if vb, ok := v.(bool); ok {
			span[k] = strconv.FormatBool(vb)
		} else if vs, ok := v.(string); ok {
			span[k] = strings.Trim(vs, `"`)
		} else {
			return nil, fmt.Errorf("value must be a boolean or a string")
		}
	}

	return span, nil
}

func MapToFormatV3(m map[string]interface{}) (FormatV3, error) {
	if _, ok := m["img"]; ok {
		format := handleImgValue(m)

		if format.Src != "" {
			return format, nil
		}

		delete(m, "img")
		delete(m, "alt")
		delete(m, "height")
		delete(m, "width")
	}

	// old style header
	if h, ok := m["header"]; ok {
		hi, err := handleHeaderValue(h)
		if err != nil {
			return nil, err
		}

		if hi != -1 {
			return FormatV3Header(hi), nil
		}

		delete(m, "header")
	}

	// new style header
	if h, ok := m["h"]; ok {
		hi, err := handleHeaderValue(h)
		if err != nil {
			return nil, err
		}

		if hi != -1 {
			return FormatV3Header(hi), nil
		}

		delete(m, "h")
	}

	// old style quill list
	if l, ok := m["list"]; ok {
		indent := 0
		var err error

		if i, ok := m["indent"]; ok {
			indent, err = handleListValue(i)
			if err != nil {
				return nil, err
			}

			delete(m, "indent")
		}

		if indent != -1 {
			if l == "ordered" || l == "\"ordered\"" {
				return FormatV3OrderedList(indent), nil
			} else if l == "bullet" || l == "\"bullet\"" {
				return FormatV3BulletList(indent), nil
			}
		}

		delete(m, "list")
	}

	// new style ordered list
	if l, ok := m["ol"]; ok {
		i, err := handleListValue(l)
		if err != nil {
			return nil, err
		}

		if i != -1 {
			return FormatV3OrderedList(i), nil
		}

		delete(m, "ol")
	}

	// new style bullet list
	if l, ok := m["ul"]; ok {
		i, err := handleListValue(l)
		if err != nil {
			return nil, err
		}

		if i != -1 {
			return FormatV3BulletList(i), nil
		}

		delete(m, "ul")
	}

	// new style indented line
	if l, ok := m["il"]; ok {
		i, err := handleListValue(l)
		if err != nil {
			return nil, err
		}

		if i != -1 {
			return FormatV3IndentedLine(i), nil
		}

		delete(m, "il")
	}

	// old style code block
	if c, ok := m["code-block"]; ok {
		if cs, ok := c.(string); ok {
			if cs != "" && cs != "null" {
				if lang, ok := m["language"]; ok {
					l := strings.Trim(lang.(string), `"`)
					return FormatV3CodeBlock(l), nil
				} else {
					return FormatV3CodeBlock(""), nil
				}
			}
		} else {
			return nil, fmt.Errorf("code-block value must be a string")
		}

		delete(m, "code-block")
	}

	// new style code block
	if c, ok := m["cb"]; ok {
		if cs, ok := c.(string); ok {
			if cs != "" && cs != "null" {
				return FormatV3CodeBlock(cs), nil
			}
		} else {
			return nil, fmt.Errorf("code-block value must be a string")
		}

		delete(m, "cb")
	}

	// old style blockquote
	if b, ok := m["blockquote"]; ok {
		if bs, ok := b.(string); ok {
			if bs != "" && bs != "null" {
				return FormatV3BlockQuote{}, nil
			}
		} else {
			return nil, fmt.Errorf("blockquote value must be a string")
		}

		delete(m, "blockquote")
	}

	// new style blockquote
	if b, ok := m["bq"]; ok {
		if bs, ok := b.(string); ok {
			if bs != "" && bs != "null" {
				return FormatV3BlockQuote{}, nil
			}
		} else {
			return nil, fmt.Errorf("blockquote value must be a string")
		}

		delete(m, "bq")
	}

	if r, ok := m["r"]; ok {
		if rs, ok := r.(string); ok {
			if rs != "" && rs != "null" {
				return FormatV3Rule{}, nil
			}
		} else {
			return nil, fmt.Errorf("rule value must be a string")
		}

		delete(m, "r")
	}

	if len(m) == 0 {
		return FormatV3Line{}, nil
	}

	return handleSpanFormat(m)
}

func FlattenMop(mop MultiOp) Op {
	if len(mop.Mops) == 1 {
		return mop.Mops[0]
	}

	return mop
}

func (mop MultiOp) Append(op Op) MultiOp {
	if op == nil {
		return mop
	}

	switch op := op.(type) {
	case MultiOp:
		mop.Mops = append(mop.Mops, op.Mops...)
	default:
		mop.Mops = append(mop.Mops, op)
	}

	return mop
}

func MaxID(op Op) ID {
	maxOp := op

	if mop, ok := op.(MultiOp); ok {
		maxOp = mop.Mops[0]
		for _, op := range mop.Mops[1:] {
			if maxOp.GetID().lessThan(op.GetID()) {
				maxOp = op
			}
		}
	}

	if op, ok := maxOp.(InsertOp); ok {
		return ID{Author: op.ID.Author, Seq: op.ID.Seq + UTF16Length(op.Text) - 1}
	} else {
		return maxOp.GetID()
	}
}
