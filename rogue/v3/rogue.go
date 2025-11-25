package v3

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
)

var (
	ErrUnknownRequiredID = errors.New("unknown ID for required ID field")
	ErrUnknownOpType     = errors.New("unknown op type")
)

var (
	RootID = ID{Author: "root", Seq: 0}
	LastID = ID{Author: "q", Seq: 1}
	NoID   = ID{Author: "", Seq: -1}
)

func ParseID(idStr string) (ID, error) {
	parts := strings.Split(idStr, "_")
	if len(parts) != 2 {
		return ID{}, fmt.Errorf("invalid ID format")
	}
	author := parts[0]
	seq, err := strconv.Atoi(parts[1])
	if err != nil {
		return ID{}, err
	}
	return ID{Author: author, Seq: seq}, nil
}

type InsertAction struct {
	Index int
	Text  string
}

func (action InsertAction) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"type":  "insert",
		"index": action.Index,
		"text":  action.Text,
	}
}

type DeleteAction struct {
	Index int
	Count int
}

func (action DeleteAction) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"type":  "delete",
		"index": action.Index,
		"count": action.Count,
	}
}

type FormatAction struct {
	Index  int
	Length int
	Format FormatV3
}

func (action FormatAction) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"type":   "format",
		"index":  action.Index,
		"length": action.Length,
		"format": action.Format.AsMap(),
	}
}

type NoOpAction struct{}

func (action NoOpAction) AsMap() map[string]interface{} {
	return map[string]interface{}{}
}

type Action interface {
	ActionType() int
	AsMap() map[string]interface{}
}

type Actions []Action

func (actions Actions) AsJS() []interface{} {
	var out []interface{}
	for _, action := range actions {
		if action == nil {
			continue
		}
		out = append(out, action.AsMap())
	}
	return out
}

func (action InsertAction) ActionType() int { return 1 }
func (action DeleteAction) ActionType() int { return 2 }
func (action FormatAction) ActionType() int { return 3 }
func (action NoOpAction) ActionType() int   { return 4 }

func (r *Rogue) GetFirstID() (ID, error) {
	return r.Rope.GetVisID(0)
}

func (r *Rogue) GetFirstTotID() (ID, error) {
	return r.Rope.GetTotID(0)
}

func (r *Rogue) GetLastID() (ID, error) {
	return r.Rope.GetVisID(r.VisSize - 1)
}

func (r *Rogue) GetLastTotID() (ID, error) {
	return r.Rope.GetTotID(r.TotSize - 1)
}

func (r *Rogue) GetWrappingIDs() (ID, ID, error) {
	first, err := r.GetFirstID()
	if err != nil {
		return NoID, NoID, err
	}
	last, err := r.GetLastID()
	if err != nil {
		return NoID, NoID, err
	}
	return first, last, nil
}

func (r *Rogue) GetWrappingTotIDs() (ID, ID, error) {
	first, err := r.GetFirstTotID()
	if err != nil {
		return NoID, NoID, err
	}
	last, err := r.GetLastTotID()
	if err != nil {
		return NoID, NoID, err
	}
	return first, last, nil
}

func (r *Rogue) GetUint16() []uint16 {
	out := make([]uint16, 0, r.VisSize)

	// use rope to just traverse the visible text
	err := r.Rope.Root.DftVis(func(rn *RopeNode) error {
		fv := rn.Val.Explode().Visible()
		out = append(out, fv.Text...)
		return nil
	})
	if err != nil {
		log.Errorf("failed to get text: %v", err)
	}

	return out
}

func (r *Rogue) GetText() string {
	return Uint16ToStr(r.GetUint16())
}

func (r *Rogue) GetTotText() string {
	out := make([]uint16, 0, r.TotSize)

	// use rope to just traverse the visible text
	err := r.Rope.Root.Dft(func(rn *RopeNode) error {
		tot := rn.Val.Explode()
		out = append(out, tot.Text...)
		return nil
	})
	if err != nil {
		log.Errorf("failed to get text: %v", err)
	}

	return Uint16ToStr(out)
}

func (r *Rogue) GetPlaintext(startID, endID ID, addr *ContentAddress) (string, error) {
	if addr == nil {
		vis, err := r.Rope.GetBetween(startID, endID)
		if err != nil {
			return "", err
		}

		return Uint16ToStr(vis.Text), nil
	} else {
		vis, err := r.Filter(startID, endID, addr)
		if err != nil {
			return "", err
		}

		return Uint16ToStr(vis.Text), nil
	}
}

func (r *Rogue) GetAllNodes() []FugueNode {
	out := []FugueNode{}

	for _, root := range r.Roots {
		root.dft(func(node *FugueNode) error {
			out = append(out, *node)
			return nil
		})
	}

	return out
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

// Helper function to create a new ParentNotFoundError
func newParentNotFoundError(message string) ErrorParentNotFound {
	return ErrorParentNotFound{message: message}
}

type UndoState struct {
	address   *ContentAddress
	redoStack []Op
}

type Rogue struct {
	Author       string
	VisSize      int
	TotSize      int
	LamportClock int
	Roots        []*FugueNode
	RootSeqs     []int
	CharHistory  CharHistory
	ScrubState   *ScrubState
	UndoState    *UndoState
	RopeIndex    RopeIndex
	OpIndex      OpIndex
	FailedOps    *FailedOps
	Rope         *Rope
	Formats      *Formats
	NOS          *NOSV2
}

func NewRogue(author string) *Rogue {
	index := RopeIndex{}
	rope := NewRope(index)
	formats := NewFormats(rope)
	nos := NewNOSV2(rope)
	rogue := &Rogue{
		Author:       author,
		VisSize:      0,
		TotSize:      0,
		LamportClock: 0,
		Roots:        []*FugueNode{},
		CharHistory:  NewCharHistory(),
		UndoState:    nil,
		RopeIndex:    index,
		OpIndex:      NewOpIndex(),
		FailedOps:    NewFailedOps(),
		Rope:         rope,
		Formats:      formats,
		NOS:          nos,
	}

	return rogue
}

func NewRogueForQuill(author string) *Rogue {
	rogue := NewRogue(author)

	root := InsertOp{RootID, "x", ID{"", 0}, Root}
	newline := InsertOp{LastID, "\n", RootID, Right}
	delroot := DeleteOp{ID{"q", 2}, RootID, 1}

	_, err := rogue.MergeOp(root)
	if err != nil {
		panic(fmt.Sprintf("failed to merge root: %v", err))
	}

	_, err = rogue.MergeOp(newline)
	if err != nil {
		panic(fmt.Sprintf("failed to merge newline: %v", err))
	}

	_, err = rogue.MergeOp(delroot)
	if err != nil {
		panic(fmt.Sprintf("failed to merge delroot: %v", err))
	}

	return rogue
}

func (r *Rogue) Reset() {
	index := RopeIndex{}
	rope := NewRope(index)

	r.VisSize = 0
	r.TotSize = 0
	r.LamportClock = 0
	r.Roots = []*FugueNode{}
	r.RootSeqs = []int{}
	r.CharHistory = NewCharHistory()
	r.UndoState = nil
	r.RopeIndex = index
	r.OpIndex = NewOpIndex()
	r.FailedOps = NewFailedOps()
	r.Rope = rope
	r.Formats = NewFormats(rope)
	r.NOS = NewNOSV2(rope)
}

func (r *Rogue) Copy(r2 *Rogue) {
	r.Author = r2.Author
	r.VisSize = r2.VisSize
	r.TotSize = r2.TotSize
	r.LamportClock = r2.LamportClock
	r.Roots = r2.Roots
	r.RootSeqs = r2.RootSeqs
	r.CharHistory = r2.CharHistory
	r.RopeIndex = r2.RopeIndex
	r.OpIndex = r2.OpIndex
	r.FailedOps = r2.FailedOps
	r.Rope = r2.Rope
	r.Formats = r2.Formats
	r.NOS = r2.NOS
}

func (r *Rogue) Reload() {
	index := RopeIndex{}
	rope := NewRope(index)

	r.VisSize = 0
	r.TotSize = 0
	r.LamportClock = 0
	r.Roots = []*FugueNode{}
	r.RootSeqs = []int{}
	r.CharHistory = NewCharHistory()
	r.UndoState = nil
	r.RopeIndex = index
	r.Rope = rope
	r.Formats = NewFormats(rope)
	r.NOS = NewNOSV2(rope)

	ops, err := r.ToOps()
	if err != nil {
		panic(fmt.Sprintf("failed to get ops: %v", err))
	}

	r.FailedOps = NewFailedOps()
	r.OpIndex = NewOpIndex()

	for _, op := range ops {
		_, err := r.MergeOp(op)
		if err != nil {
			panic(fmt.Sprintf("failed to merge op: %v", err))
		}
	}
}

func (r *Rogue) Insert(visIx int, sText string) (InsertOp, error) {
	r.UndoState = nil

	text := StrToUint16(sText)
	if !IsValidUTF16(text) {
		return InsertOp{}, stackerr.New(fmt.Errorf("invalid utf16: %v", text))
	}

	if visIx < 0 || visIx > r.VisSize {
		return InsertOp{}, stackerr.New(fmt.Errorf("index out of bounds: Unable to insert node at index %d", visIx))
	}

	if len(text) == 0 {
		return InsertOp{}, stackerr.New(fmt.Errorf("Insert(%d, %q): unable to insert empty string", visIx, sText))
	}

	var err error
	var node *FugueNode
	var ropeNode *RopeNode
	if len(r.Roots) == 0 {
		// totally empty rogue
		id := r.NextID(len(text))
		newNode := NewFugueNode(id, text, Root, nil)
		r.VisSize += len(text)
		r.TotSize += len(text)
		r.Roots = append(r.Roots, newNode)
		r.RootSeqs = append(r.RootSeqs, 0)

		_, err = r.Rope.InsertWithIx(0, newNode)
		if err != nil {
			return InsertOp{}, err
		}

		op := InsertOp{ID: id, Text: Uint16ToStr(text), ParentID: ID{"", 0}, Side: Root}
		r.OpIndex.Put(op)
		return op, nil
	}

	if r.VisSize == 0 {
		// empty visibly but there are tombstoned nodes
		// so insert to the right of the first tombstoned node
		_, ropeNode, err = r.Rope.GetTotNode(0)
		if err != nil {
			return InsertOp{}, err
		}
		node := ropeNode.Val

		id := r.NextID(len(text))
		newNode := NewFugueNode(id, text, Root, nil)
		// node.insertLeft(newNode)
		node.insertRight(newNode)

		_, err = r.Rope.Insert(newNode)
		if err != nil {
			return InsertOp{}, err
		}

		parentId := adjustedParentId(newNode)
		r.VisSize += len(text)
		r.TotSize += len(text)

		op := InsertOp{ID: id, Text: Uint16ToStr(text), ParentID: parentId, Side: newNode.Side}
		r.OpIndex.Put(op)
		return op, nil
	}

	if visIx == 0 {
		// insert to the left of the first visible node
		_, ropeNode, err = r.Rope.GetNode(0)
		if err != nil {
			return InsertOp{}, fmt.Errorf("GetNode(%v): %w", visIx-1, err)
		}
		node = ropeNode.Val

		id := r.NextID(len(text))
		newNode := NewFugueNode(id, text, Root, nil)
		node.insertLeft(newNode)

		_, err := r.Rope.Insert(newNode)
		if err != nil {
			return InsertOp{}, fmt.Errorf("Rope.Insert(%v): %w", newNode, err)
		}

		parentId := adjustedParentId(newNode)
		r.VisSize += len(text)
		r.TotSize += len(text)

		op := InsertOp{ID: id, Text: Uint16ToStr(text), ParentID: parentId, Side: newNode.Side}
		r.OpIndex.Put(op)
		return op, nil
	}

	// else insert to the right of the visibly indexed node
	visOffset := -1
	visOffset, ropeNode, err = r.Rope.GetNode(visIx - 1)
	if err != nil {
		return InsertOp{}, fmt.Errorf("GetNode(%v): %w", visIx-1, err)
	}
	node = ropeNode.Val

	totOffset, err := node.getTotOffset(visOffset)
	if err != nil {
		return InsertOp{}, fmt.Errorf("getTotOffset(%v): %w", visOffset, err)
	}

	_, rNode, err := node.splitNode(totOffset)
	if err != nil {
		return InsertOp{}, fmt.Errorf("splitNode(%v): %w", totOffset, err)
	}

	if rNode != nil {
		ropeNode.updateWeight()

		_, err = r.Rope.Insert(rNode)
		if err != nil {
			return InsertOp{}, fmt.Errorf("Rope.Insert(%v): %w", rNode, err)
		}
	}

	id := r.NextID(len(text))
	if node.ID.Author == id.Author && node.ID.Seq+len(node.Text) == id.Seq {
		// Special case: append text to existing node
		node.Text = append(node.Text, text...)
		node.IsDeleted = append(node.IsDeleted, make([]bool, len(text))...)
		if len(node.Text) != len(node.IsDeleted) {
			panic(fmt.Sprintf("text and isDeleted length mismatch: %d != %d", len(node.Text), len(node.IsDeleted)))
		}
		ropeNode.updateWeight()

		r.VisSize += len(text)
		r.TotSize += len(text)
		parentId := ID{Author: id.Author, Seq: id.Seq - 1}

		op := InsertOp{ID: id, Text: Uint16ToStr(text), ParentID: parentId, Side: Right}
		r.OpIndex.Put(op)
		return op, nil
	} else {
		newNode := NewFugueNode(id, text, Root, nil)
		node.insertRight(newNode)
		_, err := r.Rope.Insert(newNode)
		if err != nil {
			return InsertOp{}, fmt.Errorf("Rope.Insert(%v): %w", newNode, err)
		}
		parentId := adjustedParentId(newNode)
		r.VisSize += len(text)
		r.TotSize += len(text)

		op := InsertOp{ID: id, Text: Uint16ToStr(text), ParentID: parentId, Side: newNode.Side}
		r.OpIndex.Put(op)
		return op, nil
	}
}

func (r *Rogue) InsertRightOf(id ID, text string) (InsertOp, error) {
	ropeNode := r.RopeIndex.Get(id)
	if ropeNode == nil {
		return InsertOp{}, fmt.Errorf("node with ID %+v doesn't exist", id)
	}

	sib := ropeNode.Val
	totOffset := id.Seq - sib.ID.Seq

	_, rNode, err := sib.splitNode(totOffset)
	if err != nil {
		return InsertOp{}, fmt.Errorf("splitNode(%v): %w", totOffset, err)
	}

	if rNode != nil {
		ropeNode.updateWeight()

		_, err = r.Rope.Insert(rNode)
		if err != nil {
			return InsertOp{}, fmt.Errorf("Rope.Insert(%v): %w", rNode, err)
		}
	}

	u16Text := StrToUint16(text)
	newId := r.NextID(len(u16Text))
	newNode := NewFugueNode(newId, u16Text, Root, nil)
	r.VisSize += len(u16Text)
	r.TotSize += len(u16Text)
	sib.insertRight(newNode)
	parentID := adjustedParentId(newNode)

	_, err = r.Rope.Insert(newNode)
	if err != nil {
		return InsertOp{}, fmt.Errorf("Rope.Insert(%v): %w", newNode, err)
	}

	return InsertOp{ID: newId, Text: text, ParentID: parentID, Side: newNode.Side}, nil
}

func (r *Rogue) Delete(visIx, length int) (Op, error) {
	r.UndoState = nil

	if visIx < 0 || visIx+length > r.VisSize {
		return nil, stackerr.New(fmt.Errorf("index: %d length: %d out of bounds for rogue size: %d", visIx, length, r.VisSize))
	}

	if length < 1 {
		return nil, stackerr.New(fmt.Errorf("length must be greater than 0"))
	}

	startChar, err := r.GetChar(visIx)
	if err != nil {
		return nil, err
	}

	if IsLowSurrogate(startChar) {
		visIx--
	}

	endChar, err := r.GetChar(visIx + length - 1)
	if err != nil {
		return nil, err
	}

	if IsHighSurrogate(endChar) {
		length++
	}

	visOffset, ropeNode, err := r.Rope.GetNode(visIx)
	if err != nil {
		return nil, err
	}

	totOffset, err := ropeNode.Val.getTotOffset(visOffset)
	if err != nil {
		return nil, err
	}

	mop := MultiOp{}
	for {
		node := ropeNode.Val
		spanLength := 0
		id := r.NextID(1)

		for i := totOffset; i < len(node.Text); i++ {
			spanLength++
			if node.IsDeleted[i] == false {
				node.IsDeleted[i] = true
				targetID := ID{Author: node.ID.Author, Seq: node.ID.Seq + i}
				r.CharHistory.Add(targetID, &Marker{ID: id, IsDel: true})
				r.VisSize--
				length--

				if length == 0 {
					break
				}
			}
		}

		targetID := ID{Author: node.ID.Author, Seq: node.ID.Seq + totOffset}
		dop := DeleteOp{
			ID:         id,
			TargetID:   targetID,
			SpanLength: spanLength,
		}

		mop = mop.Append(dop)
		ropeNode.updateWeight()

		if length == 0 {
			op := FlattenMop(mop)
			r.OpIndex.Put(op)
			return op, nil
		}

		ropeNode, err = ropeNode.RightVisSibling()
		if err != nil {
			return mop, err
		}
		totOffset = 0
	}
}

func (r *Rogue) _trimOpIDs(startIx, endIx int, iStartID, iEndID ID) (trimStartID, trimEndID *ID, err error) {
	_, iStartIx, err := r.Rope.GetIndex(iStartID)
	if err != nil {
		return nil, nil, fmt.Errorf("GetIndex(%v): %w", iStartID, err)
	}

	_, iEndIx, err := r.Rope.GetIndex(iEndID)
	if err != nil {
		return nil, nil, fmt.Errorf("GetIndex(%v): %w", iEndID, err)
	}

	// no overlap, so continue
	if iEndIx < startIx || endIx < iStartIx {
		return nil, nil, nil
	}

	for iStartIx < startIx {
		iStartID = ID{Author: iStartID.Author, Seq: iStartID.Seq + 1}
		_, iStartIx, err = r.Rope.GetIndex(iStartID)
		if err != nil {
			return nil, nil, fmt.Errorf("GetIndex(%v): %w", iStartID, err)
		}
	}

	for endIx < iEndIx {
		iEndID = ID{Author: iEndID.Author, Seq: iEndID.Seq - 1}
		_, iEndIx, err = r.Rope.GetIndex(iEndID)
		if err != nil {
			return nil, nil, fmt.Errorf("GetIndex(%v): %w", iEndID, err)
		}
	}

	return &iStartID, &iEndID, nil
}

func (r *Rogue) _invertOp(targetID ID, op Op, startIx, endIx int) (undoOp Op, redoOp Op, err error) {
	switch op := op.(type) {
	case InsertOp:
		iStartID := op.ID
		uText := StrToUint16(op.Text)
		iEndID := ID{Author: op.ID.Author, Seq: op.ID.Seq + len(uText) - 1}

		trimStartID, trimEndID, err := r._trimOpIDs(startIx, endIx, iStartID, iEndID)
		if err != nil {
			return nil, nil, fmt.Errorf("_trimOpIDs(%v, %v, %v, %v): %w", startIx, endIx, iStartID, iEndID, err)
		}

		if trimStartID == nil || trimEndID == nil {
			return nil, nil, nil
		}

		spanLength := trimEndID.Seq - trimStartID.Seq + 1
		dop := DeleteOp{ID: r.NextID(1), TargetID: *trimStartID, SpanLength: spanLength}
		rop := ShowOp{ID: NoID, TargetID: *trimStartID, SpanLength: spanLength}
		return dop, rop, nil
	case DeleteOp:
		iStartID := op.TargetID
		iEndID := ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + op.SpanLength - 1}

		trimStartID, trimEndID, err := r._trimOpIDs(startIx, endIx, iStartID, iEndID)
		if err != nil {
			return nil, nil, fmt.Errorf("_trimOpIDs(%v, %v, %v, %v): %w", startIx, endIx, iStartID, iEndID, err)
		}

		if trimStartID == nil || trimEndID == nil {
			return nil, nil, nil
		}

		spanLength := trimEndID.Seq - trimStartID.Seq + 1
		sop := ShowOp{ID: r.NextID(1), TargetID: *trimStartID, SpanLength: spanLength}
		rop := DeleteOp{ID: NoID, TargetID: *trimStartID, SpanLength: spanLength}
		return sop, rop, nil
	case ShowOp:
		iStartID := op.TargetID
		iEndID := ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + op.SpanLength - 1}

		trimStartID, trimEndID, err := r._trimOpIDs(startIx, endIx, iStartID, iEndID)
		if err != nil {
			return nil, nil, fmt.Errorf("_trimOpIDs(%v, %v, %v, %v): %w", startIx, endIx, iStartID, iEndID, err)
		}

		if trimStartID == nil || trimEndID == nil {
			return nil, nil, nil
		}

		spanLength := trimEndID.Seq - trimStartID.Seq + 1
		dop := DeleteOp{ID: r.NextID(1), TargetID: *trimStartID, SpanLength: spanLength}
		rop := ShowOp{ID: NoID, TargetID: *trimStartID, SpanLength: spanLength}
		return dop, rop, nil
	case FormatOp:
		uop, err := r.invertFormatOp(targetID, op)
		if err != nil {
			return nil, nil, fmt.Errorf("invertFormatOp(%v): %w", op, err)
		}

		return uop, op, nil
	case RewindOp:
		iStartID := op.Address.StartID
		iEndID := op.Address.EndID

		trimStartID, trimEndID, err := r._trimOpIDs(startIx, endIx, iStartID, iEndID)
		if err != nil {
			return nil, nil, fmt.Errorf("_trimOpIDs(%v, %v, %v, %v): %w", startIx, endIx, iStartID, iEndID, err)
		}

		if trimStartID == nil || trimEndID == nil {
			return nil, nil, nil
		}

		address := ContentAddress{StartID: *trimStartID, EndID: *trimEndID, MaxIDs: op.UndoAddress.MaxIDs}
		undoAddress := ContentAddress{StartID: *trimStartID, EndID: *trimEndID, MaxIDs: op.Address.MaxIDs}
		rop := RewindOp{ID: r.NextID(1), Address: address, UndoAddress: undoAddress}
		return rop, rop, nil
	case MultiOp:
		umop, rmop := MultiOp{}, MultiOp{}

		for i := len(op.Mops) - 1; i >= 0; i-- {
			uop, rop, err := r._invertOp(targetID, op.Mops[i], startIx, endIx)
			if err != nil {
				return nil, nil, fmt.Errorf("_invertOp(%v, %v, %v): %w", op, startIx, endIx, err)
			}

			umop = umop.Append(uop)
			rmop = rmop.Append(rop)
		}

		return FlattenMop(umop), FlattenMop(rmop), nil
	default:
		return nil, nil, fmt.Errorf("unknown op type: %w", ErrUnknownOpType)
	}
}

func (r *Rogue) _undoNext(startID, endID ID, ca *ContentAddress) (uop, rop Op, err error) {
	if ca == nil {
		return nil, nil, nil
	}

	if len(ca.MaxIDs) == 0 {
		return nil, nil, nil
	}

	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	curID := ca.MaxID()
	if curID == nil {
		return nil, nil, nil
	}

	op := r.OpIndex.Get(*curID)
	if op == nil {
		return nil, nil, fmt.Errorf("OpIndex.Get(%v): %w", curID, ErrUnknownRequiredID)
	}

	opID := op.GetID()
	nextID := r.OpIndex.nextSmallest(opID)
	if nextID == nil {
		delete(ca.MaxIDs, opID.Author)
	} else {
		ca.MaxIDs[opID.Author] = nextID.Seq
	}

	uop, rop, err = r._invertOp(op.GetID(), op, startIx, endIx)
	if err != nil {
		return nil, nil, fmt.Errorf("_invertOp(%v, %v, %v): %w", op, startIx, endIx, err)
	}

	return uop, rop, nil
}

func (r *Rogue) TrimIDSpan(startID, endID, toTrimStart, toTrimEnd ID) (tStartID, tEndID ID, err error) {
	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return NoID, NoID, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return NoID, NoID, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	_, tStartIx, err := r.Rope.GetIndex(toTrimStart)
	if err != nil {
		return NoID, NoID, fmt.Errorf("GetIndex(%v): %w", toTrimStart, err)
	}

	_, tEndIx, err := r.Rope.GetIndex(toTrimEnd)
	if err != nil {
		return NoID, NoID, fmt.Errorf("GetIndex(%v): %w", toTrimEnd, err)
	}

	startIx = max(startIx, tStartIx)
	endIx = min(endIx, tEndIx)

	if endIx < startIx {
		return NoID, NoID, fmt.Errorf("totEndIx < totStartIx: %d < %d", endIx, startIx)
	}

	tStartID, err = r.Rope.GetTotID(startIx)
	if err != nil {
		return NoID, NoID, fmt.Errorf("GetTotID(%d): %w", startIx, err)
	}

	tEndID, err = r.Rope.GetTotID(endIx)
	if err != nil {
		return NoID, NoID, fmt.Errorf("GetTotID(%d): %w", endIx, err)
	}

	return tStartID, tEndID, nil
}

func (r *Rogue) TotRightOf(id ID) (ID, error) {
	return r.Rope.TotRightOf(id)
}

func (r *Rogue) TotLeftOf(id ID) (ID, error) {
	return r.Rope.TotLeftOf(id)
}

func (r *Rogue) VisRightOf(id ID) (ID, error) {
	return r.Rope.VisRightOf(id)
}

func (r *Rogue) VisLeftOf(id ID) (ID, error) {
	return r.Rope.VisLeftOf(id)
}

func (r *Rogue) VisScanLeftOf(id ID, value uint16) (*ID, error) {
	return r.VisScanLeftOfFunc(id, func(v uint16) bool {
		return v == value
	})
}

func (r *Rogue) VisScanLeftOfFunc(id ID, fn func(v uint16) bool) (*ID, error) {
	var (
		leftValue uint16
		err       error
	)

	isDel, err := r.IsDeleted(id)
	if err != nil {
		return nil, err
	}

	if isDel {
		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return nil, nil
			} else {
				return nil, err
			}
		}
	}

	for {
		leftValue, err = r.GetCharByID(id)
		if err != nil {
			return nil, err
		}

		if fn(leftValue) {
			return &id, nil
		}

		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return nil, nil
			}

			return nil, err
		}
	}
}

func (r *Rogue) VisScanRightOf(id ID, value uint16) (*ID, error) {
	return r.VisScanRightOfFunc(id, func(v uint16) bool {
		return v == value
	})
}

func (r *Rogue) VisScanRightOfFunc(id ID, fn func(uint16) bool) (*ID, error) {
	var (
		rightValue uint16
		err        error
	)

	isDel, err := r.IsDeleted(id)
	if err != nil {
		return nil, err
	}

	if isDel {
		id, err = r.VisRightOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return nil, nil
			}

			return nil, err
		}
	}

	for {
		rightValue, err = r.GetCharByID(id)
		if err != nil {
			return nil, err
		}

		if fn(rightValue) {
			return &id, nil
		}

		id, err = r.VisRightOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return nil, nil
			}

			return nil, err
		}
	}
}

func (r *Rogue) ValidateFormat(op FormatOp) error {
	if !op.Format.IsSpan() {
		if op.StartID != op.EndID {
			return fmt.Errorf("line format op has different start and end id: %v", op)
		}
		return nil
	}

	format := op.Format.(FormatV3Span)

	sticky, noSticky := format.SplitSticky()

	if len(sticky) > 0 && len(noSticky) > 0 {
		return fmt.Errorf("format has both sticky and non-sticky attributes: %v", op)
	}

	_, startIx, err := r.Rope.GetIndex(op.StartID)
	if err != nil {
		return fmt.Errorf("GetIndex(%v): %w", op.StartID, err)
	}

	_, endIx, err := r.Rope.GetIndex(op.EndID)
	if err != nil {
		return fmt.Errorf("GetIndex(%v): %w", op.EndID, err)
	}

	if len(sticky) > 0 {
		endIx--
	}

	if endIx < startIx {
		return fmt.Errorf("endIx < startIx: %d < %d for op: %v", endIx, startIx, op)
	}

	return nil
}

func (r *Rogue) _insertListFormat(ix, length int, format FormatV3) (Op, error) {
	lm := getListMeta(format)
	if !lm.isList {
		return nil, nil
	}

	startID, err := r.Rope.GetVisID(ix)
	if err != nil {
		return nil, err
	}

	id := startID

	// look back to change any equal indent list items
	mop := MultiOp{}
	for {
		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				break
			}
			return nil, err
		}

		lid, err := r.VisScanLeftOf(id, '\n')
		if err != nil {
			return nil, err
		}

		if lid == nil {
			break
		}

		prevFormat, err := r.GetCurLineFormat(*lid, *lid)
		if err != nil {
			return nil, err
		}

		prevLm := getListMeta(prevFormat)
		if prevLm.isList {
			isDiffList := (prevLm.isBullet != lm.isBullet || prevLm.isOrdered != lm.isOrdered) && !prevLm.isIndent
			if prevLm.indent < lm.indent {
				break
			} else if prevLm.indent == lm.indent && isDiffList {
				fop := FormatOp{
					ID:      r.NextID(1),
					StartID: *lid,
					EndID:   *lid,
					Format:  format,
				}

				mop = mop.Append(fop)
			}
		} else {
			break
		}

		id = *lid
	}

	// look forward until end of format span
	endIx := min(ix+length-1, r.VisSize-1)
	id = startID
	for {
		rid, err := r.VisScanRightOf(id, '\n')
		if err != nil {
			return nil, err
		}

		if rid == nil {
			op, err := r.Insert(r.VisSize, "\n")
			if err != nil {
				return nil, err
			}

			mop = mop.Append(op)
			rid = &op.ID
		}

		targetIx, _, err := r.Rope.GetIndex(*rid)
		if err != nil {
			return nil, err
		}

		mop = mop.Append(FormatOp{
			ID:      r.NextID(1),
			StartID: *rid,
			EndID:   *rid,
			Format:  format,
		})

		if targetIx >= endIx {
			break
		}

		id, err = r.TotRightOf(*rid)
		if err != nil {
			if errors.As(err, &ErrorNoRightTotSibling{}) {
				break
			}

			return nil, err
		}
	}

	// keep going forward to change any equal indent list items
	for {
		id, err = r.VisRightOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				break
			}
			return nil, err
		}

		rid, err := r.VisScanRightOf(id, '\n')
		if err != nil {
			return nil, err
		}

		if rid == nil {
			break
		}

		nextFormat, err := r.GetCurLineFormat(*rid, *rid)
		if err != nil {
			return nil, err
		}

		nextLm := getListMeta(nextFormat)
		if nextLm.isList {
			isDiffList := (nextLm.isBullet != lm.isBullet || nextLm.isOrdered != lm.isOrdered) && !nextLm.isIndent
			if nextLm.indent < lm.indent {
				break
			} else if nextLm.indent == lm.indent && isDiffList {
				fop := FormatOp{
					ID:      r.NextID(1),
					StartID: *rid,
					EndID:   *rid,
					Format:  format,
				}

				mop = mop.Append(fop)
			}
		} else {
			break
		}

		id = *rid
	}

	return FlattenMop(mop), nil
}

func (r *Rogue) _insertLineFormat(ix, length int, format FormatV3) (Op, error) {
	mop := MultiOp{}

	if _, ok := format.(FormatV3Image); ok {
		endIx := min(ix+length-1, r.VisSize-1)
		endID, err := r.Rope.GetVisID(endIx)
		if err != nil {
			return nil, err
		}

		rid, err := r.VisScanRightOf(endID, '\n')
		if err != nil {
			return nil, err
		}

		if rid == nil {
			op, err := r.Insert(r.VisSize, "\n")
			if err != nil {
				return nil, err
			}

			mop = mop.Append(op)
			rid = &op.ID
		}

		fop := FormatOp{
			ID:      r.NextID(1),
			StartID: *rid,
			EndID:   *rid,
			Format:  format,
		}

		_, err = r.NOS.Insert(fop)
		if err != nil {
			return nil, err
		}

		err = r.Formats.Insert(fop)
		if err != nil {
			return nil, err
		}

		mop = mop.Append(fop)

		return FlattenMop(mop), nil
	}

	id, err := r.Rope.GetVisID(ix)
	if err != nil {
		return nil, fmt.Errorf("GetVisId(%d): %w", ix, err)
	}

	endIx := min(ix+length-1, r.VisSize-1)
	for {
		rid, err := r.VisScanRightOf(id, '\n')
		if err != nil {
			return nil, err
		}

		if rid == nil {
			op, err := r.Insert(r.VisSize, "\n")
			if err != nil {
				return nil, err
			}

			mop = mop.Append(op)
			rid = &op.ID
		}

		targetIx, _, err := r.Rope.GetIndex(*rid)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", *rid, err)
		}

		fop := FormatOp{
			ID:      r.NextID(1),
			StartID: *rid,
			EndID:   *rid,
			Format:  format,
		}

		_, err = r.NOS.Insert(fop)
		if err != nil {
			return nil, err
		}

		err = r.Formats.Insert(fop)
		if err != nil {
			return nil, err
		}

		mop = mop.Append(fop)

		if targetIx >= endIx {
			break
		}

		id, err = r.TotRightOf(*rid)
		if err != nil {
			if errors.As(err, &ErrorNoRightTotSibling{}) {
				break
			}

			return nil, fmt.Errorf("TotRightOf(%v): %w", *rid, err)
		}
	}

	return FlattenMop(mop), nil
}

var (
	stickyTag      = map[string]string{"e": "true"}
	noStickyTag    = map[string]string{"en": "true"}
	stickyAndNoTag = map[string]string{"e": "true", "en": "true"}
)

func fopWithTag(op FormatOp, tags map[string]string) Op {
	if format, ok := op.Format.(FormatV3Span); ok {
		format = format.Copy().(FormatV3Span)
		for k, v := range tags {
			format[k] = v
		}

		return FormatOp{
			ID:      op.ID,
			StartID: op.StartID,
			EndID:   op.EndID,
			Format:  format,
		}
	}

	return op
}

func (r *Rogue) Format(ix, length int, format FormatV3) (Op, error) {
	r.UndoState = nil
	mop := MultiOp{}

	if ix < 0 || ix+length > r.VisSize {
		return nil, stackerr.New(fmt.Errorf("index: %d length: %d out of bounds for rogue size: %d", ix, length, r.VisSize))
	}

	startID, err := r.Rope.GetVisID(ix)
	if err != nil {
		return nil, err
	}

	if format.IsSpan() {
		f := format.(FormatV3Span)
		sticky, noSticky := f.SplitSticky()

		if len(sticky) > 0 {
			if ix+length >= r.VisSize {
				return nil, stackerr.New(fmt.Errorf("sticky span index: %d length: %d out of bounds for rogue size: %d", ix, length, r.VisSize))
			}

			endID, err := r.Rope.GetVisID(ix + length - 1)
			if err != nil {
				return nil, err
			}

			afterID, err := r.TotRightOf(endID)
			if err != nil {
				return nil, err
			}

			inserted, err := r.NOS.Insert(FormatOp{
				StartID: startID,
				EndID:   afterID,
				Format:  sticky,
			})

			for _, op := range inserted {
				op.ID = r.NextID(1)
				mop = mop.Append(op)

				err = r.Formats.Sticky.Insert(op)
				if err != nil {
					return nil, err
				}
			}
		}

		if len(noSticky) > 0 {
			endID, err := r.Rope.GetVisID(ix + length - 1)
			if err != nil {
				return nil, err
			}

			inserted, err := r.NOS.Insert(FormatOp{
				StartID: startID,
				EndID:   endID,
				Format:  noSticky,
			})

			for _, op := range inserted {
				op.ID = r.NextID(1)
				mop = mop.Append(op)

				err = r.Formats.NoSticky.Insert(op)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		// clamp indent if format is a list
		op, err := r._insertLineFormat(ix, length, format)
		if err != nil {
			return nil, fmt.Errorf("_insertLineFormat(%d, %d, %v): %w", ix, length, format, err)
		}
		mop = mop.Append(op)
	}

	op := FlattenMop(mop)

	/*_, err = r.MergeOp(op)
	if err != nil {
		return nil, fmt.Errorf("MergeOp(%v): %w", op, err)
	}*/

	r.OpIndex.Put(op)

	return op, nil
}

func (r *Rogue) FormatLineByID(id ID, format FormatV3) (Op, error) {
	op := FormatOp{
		ID:      r.NextID(1),
		StartID: id,
		EndID:   id,
		Format:  format,
	}

	_, err := r.MergeOp(op)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (r *Rogue) NextID(length int) ID {
	lc := r.LamportClock
	r.LamportClock += length

	return ID{Author: r.Author, Seq: lc}
}

func (node *FugueNode) splitNode(totOffset int) (*FugueNode, *FugueNode, error) {
	if totOffset < 0 || totOffset >= len(node.Text)-1 {
		return node, nil, nil
	}

	totOffset++
	if IsLowSurrogate(node.Text[totOffset]) {
		// surPair := Uint16ToStr(node.Text[totOffset-1 : totOffset+1])
		vis := node.Explode().Visible()
		return nil, nil, stackerr.New(fmt.Errorf("can not split node: %s totText: %q visText: %q at ix: %d", node.ID, Uint16ToStr(node.Text), Uint16ToStr(vis.Text), totOffset))
	}

	newNode := NewFugueNode(ID{node.ID.Author, node.ID.Seq + totOffset},
		node.Text[totOffset:],
		Right,
		node)

	lDel, rDel := node.IsDeleted[:totOffset], node.IsDeleted[totOffset:]
	node.IsDeleted = lDel
	newNode.IsDeleted = rDel

	newNode.Parent = node
	node.Text = node.Text[:totOffset]
	newNode.RightChildren = node.RightChildren
	for _, child := range newNode.RightChildren {
		child.Parent = newNode
	}
	node.RightChildren = []*FugueNode{newNode}

	return node, newNode, nil
}

func (r *Rogue) insertNewRoot(op InsertOp) (Actions, error) {
	text := StrToUint16(op.Text)
	newRoot := NewFugueNode(op.ID, text, Root, nil)

	ropeIx := 0
	seqIx := bisectLeft(r.RootSeqs, op.ParentID.Seq, func(rootSeq int) int {
		return rootSeq
	})

	var err error
	if len(r.RootSeqs) == 0 {
		// do nothing
	} else if seqIx < len(r.RootSeqs) {
		sibRoot := r.Roots[seqIx]
		sib := sibRoot.leftmost()
		_, ropeIx, err = r.Rope.GetIndex(sib.ID)
		if err != nil {
			return nil, fmt.Errorf("Rope.GetIndex(%v): %w", sib.ID, err)
		}
	} else {
		sibRoot := r.Roots[seqIx-1]
		sib := sibRoot.rightmost()
		_, ropeIx, err = r.Rope.GetIndex(sib.ID)
		if err != nil {
			return nil, fmt.Errorf("Rope.GetIndex(%v): %w", sib.ID, err)
		}
		ropeIx++
	}

	r.RootSeqs = InsertAt(r.RootSeqs, seqIx, op.ParentID.Seq)
	r.Roots = InsertAt(r.Roots, seqIx, newRoot)
	_, err = r.Rope.InsertWithIx(ropeIx, newRoot)
	if err != nil {
		return nil, fmt.Errorf("Rope.InsertWithIx(%v): %w", ropeIx, err)
	}

	r.LamportClock = max(r.LamportClock, op.ID.Seq+len(op.Text))
	r.VisSize += len(text)
	r.TotSize += len(text)

	visIx, _, err := r.Rope.GetIndex(op.ID)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetIndex(%v): %v", op.ID, err)
	}

	return Actions{InsertAction{Index: visIx, Text: op.Text}}, nil
}

func (r *Rogue) insertOp(op InsertOp) (Actions, error) {
	existingNode := r.RopeIndex.Get(op.ID)
	if existingNode != nil {
		return nil, nil
	}

	if op.Side == Root {
		return r.insertNewRoot(op)
	}

	text := StrToUint16(op.Text)
	parentRope := r.RopeIndex.Get(op.ParentID)
	if parentRope == nil {
		return nil, stackerr.New(newParentNotFoundError(fmt.Sprintf("Cannot insert: Parent node with ID %+v doesn't exist.", op.ParentID)))
	}
	parentNode := parentRope.Val

	if op.Side == Right && op.ParentID.Author == op.ID.Author && op.ParentID.Seq == op.ID.Seq-1 {
		parentNode.Text = append(parentNode.Text, text...)
		parentNode.IsDeleted = append(parentNode.IsDeleted, make([]bool, len(text))...)
		if len(parentNode.Text) != len(parentNode.IsDeleted) {
			panic(fmt.Sprintf("text and isDeleted length mismatch: %d != %d", len(parentNode.Text), len(parentNode.IsDeleted)))
		}
		r.LamportClock = max(r.LamportClock, op.ID.Seq+len(op.Text))
		r.VisSize += len(text)
		r.TotSize += len(text)

		parentRope.updateWeight()

		visIx, _, err := r.Rope.GetIndex(op.ID)
		if err != nil {
			return nil, err
		}

		return Actions{InsertAction{Index: visIx, Text: op.Text}}, nil
	}

	newNode := NewFugueNode(op.ID, text, Root, parentNode)
	totOffset := op.ParentID.Seq - parentNode.ID.Seq

	if op.Side == Right {
		_, rNode, err := parentNode.splitNode(totOffset)
		if err != nil {
			return nil, err
		}

		if rNode != nil {
			parentRope.updateWeight()

			_, err = r.Rope.Insert(rNode)
			if err != nil {
				return nil, err
			}
		}
		newNode.Side = Right
		parentNode.insertChild(Right, newNode)
		_, err = r.Rope.Insert(newNode)
		if err != nil {
			return nil, err
		}
	} else {
		lNode, rNode, err := parentNode.splitNode(totOffset - 1)
		if err != nil {
			return nil, err
		}

		if rNode != nil {
			parentRope.updateWeight()

			parentNode = rNode
			_, err = r.Rope.Insert(rNode)
			if err != nil {
				return nil, err
			}
		} else {
			parentNode = lNode
		}

		newNode.Side = Left
		newNode.Parent = parentNode
		parentNode.insertChild(Left, newNode)
		_, err = r.Rope.Insert(newNode)
		if err != nil {
			return nil, err
		}
	}

	r.LamportClock = max(r.LamportClock, op.ID.Seq+len(op.Text))
	r.VisSize += len(text)
	r.TotSize += len(text)

	visIx, _, err := r.Rope.GetIndex(op.ID)
	if err != nil {
		return nil, err
	}

	return Actions{InsertAction{Index: visIx, Text: op.Text}}, nil
}

func minIdSeq(a, b ID) ID {
	if a.Seq < b.Seq {
		return a
	}
	return b
}

func (r *Rogue) deleteOp(op DeleteOp) (Actions, error) {
	id, targetID, spanLength := op.ID, op.TargetID, op.SpanLength
	ropeNode := r.RopeIndex.Get(targetID)
	if ropeNode == nil {
		return nil, stackerr.New(newParentNotFoundError(fmt.Sprintf("Cannot delete: Node with ID %+v doesn't exist.", targetID)))
	}
	if r.TotSize == 0 {
		return nil, stackerr.New(fmt.Errorf("cannot delete: Rope is empty"))
	}

	var err error
	node := ropeNode.Val
	deleted, visIx := 0, -1
	totOffset := targetID.Seq - node.ID.Seq
	r.LamportClock = max(r.LamportClock, id.Seq+1)

	actions := make(Actions, 0)
	for i := 0; i < spanLength; i++ {
		if totOffset == len(node.Text) {
			if visIx >= 0 {
				actions = append(actions, DeleteAction{Index: visIx, Count: deleted})
			}

			ropeNode.updateWeight()

			ropeNode = r.RopeIndex.Get(targetID)
			if ropeNode == nil {
				return nil, stackerr.New(newParentNotFoundError(fmt.Sprintf("Cannot delete: Node with ID %+v doesn't exist.", targetID)))
			}

			node = ropeNode.Val
			totOffset, deleted, visIx = 0, 0, -1
		}

		r.CharHistory.Add(targetID, &Marker{ID: id, IsDel: true})
		maxMarker := r.CharHistory.Max(targetID)

		if node.IsDeleted[totOffset] == false && maxMarker.IsDel == true {
			if visIx < 0 { // We haven't found a visible character yet
				visIx, _, err = r.Rope.GetIndex(ID{targetID.Author, node.ID.Seq + totOffset})
				if err != nil {
					return nil, err
				}
			}

			node.IsDeleted[totOffset] = true
			r.VisSize--
			deleted++
		} else if visIx >= 0 {
			actions = append(actions, DeleteAction{Index: visIx, Count: deleted})
			deleted, visIx = 0, -1
		}

		totOffset++
		targetID = ID{node.ID.Author, node.ID.Seq + totOffset}
	}

	ropeNode.updateWeight()

	if visIx >= 0 {
		actions = append(actions, DeleteAction{Index: visIx, Count: deleted})
	}

	return actions, nil
}

func (r *Rogue) showOp(op ShowOp) (Actions, error) {
	id, targetID, spanLength := op.ID, op.TargetID, op.SpanLength
	ropeNode := r.RopeIndex.Get(targetID)
	if ropeNode == nil {
		return nil, stackerr.New(newParentNotFoundError(fmt.Sprintf("Cannot show: Node with ID %+v doesn't exist.", targetID)))
	}
	if r.TotSize == 0 {
		return nil, stackerr.New(fmt.Errorf("cannot show: Rope is empty"))
	}

	var err error
	node := ropeNode.Val
	visIx := -1
	totOffset := targetID.Seq - node.ID.Seq
	r.LamportClock = max(r.LamportClock, id.Seq+1)

	actions := make(Actions, 0)
	for i := 0; i < spanLength; i++ {
		if totOffset == len(node.Text) {
			if visIx >= 0 {
				// actions = append(actions, DeleteAction{Index: visIx, Count: 1})
			}

			ropeNode.updateWeight()

			ropeNode = r.RopeIndex.Get(targetID)
			if ropeNode == nil {
				return nil, stackerr.New(newParentNotFoundError(fmt.Sprintf("Cannot delete: Node with ID %+v doesn't exist.", targetID)))
			}

			node = ropeNode.Val
			totOffset, visIx = 0, -1
		}

		r.CharHistory.Add(targetID, &Marker{ID: id, IsDel: false})
		maxMarker := r.CharHistory.Max(targetID)

		if node.IsDeleted[totOffset] == true && maxMarker.IsDel == false {
			if visIx < 0 { // We haven't found a visible character yet
				visIx, _, err = r.Rope.GetIndex(ID{targetID.Author, node.ID.Seq + totOffset})
				if err != nil {
					return nil, err
				}
			}

			node.IsDeleted[totOffset] = false
			r.VisSize++
		} else if visIx >= 0 {
			actions = append(actions, DeleteAction{Index: visIx, Count: 1})
			visIx = -1
		}

		totOffset++
		targetID = ID{node.ID.Author, node.ID.Seq + totOffset}
	}

	ropeNode.updateWeight()

	if visIx >= 0 {
		actions = append(actions, DeleteAction{Index: visIx, Count: 1})
	}

	return actions, nil
}

func (r *Rogue) FormatOp(formatOp FormatOp) ([]Action, error) {
	err := r.insertFormatOp(formatOp)
	if err != nil {
		return nil, fmt.Errorf("indexFormatOp(%v): %w", formatOp, err)
	}

	r.LamportClock = max(r.LamportClock, formatOp.ID.Seq+1)

	lineActions, err := r.LineFormatAction(formatOp)
	if err != nil {
		return nil, fmt.Errorf("LineFormatAction(%v): %w", formatOp, err)
	}

	spanActions, err := r.SpanFormatAction(formatOp)
	if err != nil {
		return nil, fmt.Errorf("SpanFormatAction(%v): %w", formatOp, err)
	}

	actions := make([]Action, 0, len(lineActions)+len(spanActions))
	actions = append(actions, lineActions...)
	actions = append(actions, spanActions...)

	return actions, nil
}

func (r *Rogue) IsDeleted(id ID) (bool, error) {
	return r.Rope.IsDeleted(id)
}

func (r *Rogue) LineFormatAction(formatOp FormatOp) ([]Action, error) {
	if formatOp.Format.IsSpan() {
		return nil, nil
	}

	if formatOp.StartID != formatOp.EndID {
		log.Warnf("LineFormatAction: startID != endID: %v", formatOp)
		return nil, nil
	}

	isDel, err := r.IsDeleted(formatOp.StartID)
	if err != nil {
		return nil, fmt.Errorf("IsDeleted(%v): %w", formatOp.StartID, err)
	}
	if isDel {
		return []Action{}, nil
	}

	c, err := r.GetCharByID(formatOp.StartID)
	if err != nil {
		return nil, fmt.Errorf("GetCharByID(%v): %w", formatOp.StartID, err)
	}

	if c != '\n' {
		log.Warnf("LineFormatAction: char != \\n: %v", formatOp)
		return nil, nil
	}

	startIx, _, err := r.Rope.GetIndex(formatOp.StartID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", formatOp.StartID, err)
	}

	return []Action{FormatAction{
		Index:  startIx,
		Length: 1,
		Format: formatOp.Format,
	}}, nil
}

func (r *Rogue) SpanFormatAction(formatOp FormatOp) ([]Action, error) {
	if !formatOp.Format.IsSpan() {
		return nil, nil
	}

	startID := formatOp.StartID
	isDeleted, err := r.IsDeleted(startID)
	if err != nil {
		return nil, fmt.Errorf("IsDeleted(%v): %w", startID, err)
	}

	if isDeleted {
		startID, err = r.VisRightOf(startID)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return nil, nil
			}
			return nil, fmt.Errorf("VisRightOf(%v): %w", startID, err)
		}
	}

	startIx, _, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	actions := make([]Action, 0)
	if formatOp.Format.IsSpan() {
		f := formatOp.Format.(FormatV3Span)
		sticky, noSticky := f.SplitSticky()

		if len(sticky) > 0 {
			endID, err := r.VisLeftOf(formatOp.EndID)
			if err != nil {
				if errors.As(err, &ErrorNoLeftVisSibling{}) {
					return nil, nil
				}
				return nil, fmt.Errorf("VisLeftOf(%v): %w", formatOp.EndID, err)
			}

			endIx, _, err := r.Rope.GetIndex(endID)
			if err != nil {
				return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
			}

			if startIx <= endIx {
				actions = append(actions, FormatAction{
					Index:  startIx,
					Length: endIx - startIx + 1,
					Format: sticky,
				})
			}
		}

		if len(noSticky) > 0 {
			endID := formatOp.EndID
			isDeleted, err := r.IsDeleted(endID)
			if err != nil {
				return nil, fmt.Errorf("IsDeleted(%v): %w", startID, err)
			}

			if isDeleted {
				endID, err = r.VisLeftOf(formatOp.EndID)
				if err != nil {
					if errors.As(err, &ErrorNoLeftVisSibling{}) {
						return []Action{}, nil
					}
					return nil, fmt.Errorf("VisLeftOf(%v): %w", formatOp.EndID, err)
				}
			}

			endIx, _, err := r.Rope.GetIndex(endID)
			if err != nil {
				return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
			}

			if startIx <= endIx {
				actions = append(actions, FormatAction{
					Index:  startIx,
					Length: endIx - startIx + 1,
					Format: noSticky,
				})
			}
		}
	}

	return actions, nil
}

func (r *Rogue) snapshotOp(op SnapshotOp) ([]Action, error) {
	err := r.DeserRogue(op.Snapshot)
	if err != nil {
		return nil, err
	}

	actions := make([]Action, 0)

	text := r.GetText()
	if len(text) > 0 {
		insertAction := InsertAction{
			Index: 0,
			Text:  text[:len(text)-1], // drop the \n that's already in quill
		}
		actions = append(actions, insertAction)
	}

	// TODO: add format actions if we want to support that

	return actions, nil
}

// ValidateFugue validates the Fugue tree
//   - Checks that all children have the correct parent and side
func (r *Rogue) ValidateFugues() error {
	for _, node := range r.Roots {
		err := node.ValidateParentSide()
		if err != nil {
			return fmt.Errorf("ValidateParentSide(): %w", err)
		}
	}

	return nil
}

func (r *Rogue) _mergeOp(op Op) (Actions, error) {
	if r.OpIndex.GetExact(op.GetID()) != nil {
		return nil, nil
	}

	switch op := op.(type) {
	case InsertOp:
		actions, err := r.insertOp(op)
		if err != nil {
			r.FailedOps.Put(op)
			return nil, err
		}

		return actions, nil
	case DeleteOp:
		actions, err := r.deleteOp(op)
		if err != nil {
			r.FailedOps.Put(op)
			return nil, err
		}

		return actions, nil
	case ShowOp:
		actions, err := r.showOp(op)
		if err != nil {
			r.FailedOps.Put(op)
			return nil, err
		}

		return actions, nil
	case FormatOp:
		actions, err := r.FormatOp(op)
		if err != nil {
			r.FailedOps.Put(op)
			return nil, err
		}

		return actions, nil
	case RewindOp:
		err := r.rewindOp(op)
		if err != nil {
			r.FailedOps.Put(op)
			return nil, err
		}

		// TODO: return actions from rewind
		return nil, nil
	case MultiOp:
		return nil, stackerr.Errorf("Nested MultiOp not supported: %v", op)
	case SnapshotOp:
		return nil, stackerr.Errorf("Nested SnapshotOp not supported: %v", op)
	default:
		return nil, stackerr.Errorf("unknown operation type: %T (Value: %+v)", op, op)
	}
}

func (r *Rogue) MergeOp(op Op) (Actions, error) {
	if op == nil {
		return nil, nil
	}

	ops := []Op{op}

	switch op := op.(type) {
	case SnapshotOp:
		return r.snapshotOp(op)
	case MultiOp:
		if len(op.Mops) == 0 {
			return nil, nil
		}

		opSeen := r.OpIndex.GetExact(op.GetID())
		if mop, ok := opSeen.(MultiOp); ok {
			if len(mop.Mops) != len(op.Mops) {
				return nil, stackerr.Errorf("MultiOp with same ID but different length: %v != %v", len(mop.Mops), len(op.Mops))
			}

			return nil, nil // already merged
		}

		ops = op.Mops
	}

	outActions := make([]Action, 0)

	for _, op2 := range ops {
		actions, err := r._mergeOp(op2)
		if err != nil {
			r.OpIndex.Remove(op)
			r.FailedOps.Put(op)

			return nil, err
		}
		outActions = append(outActions, actions...)
	}

	r.OpIndex.Put(op)
	return outActions, nil
}

func (r *Rogue) GetChar(visIx int) (uint16, error) {
	visOffset, node, err := r.Rope.GetNode(visIx)
	if err != nil {
		return 0, err
	}

	totOffset, err := node.Val.getTotOffset(visOffset)
	if err != nil {
		return 0, err
	}
	return node.Val.Text[totOffset], nil
}

func (r *Rogue) GetCharByID(id ID) (uint16, error) {
	rn := r.RopeIndex.Get(id)
	if rn == nil {
		return 0, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	totIx := id.Seq - rn.Val.ID.Seq
	if rn.Val.IsDeleted[totIx] == true {
		return 0, stackerr.New(fmt.Errorf("node with ID %+v is deleted", id))
	}

	return rn.Val.Text[totIx], nil
}

func (r *Rogue) GetTotCharByID(id ID) (uint16, error) {
	rn := r.RopeIndex.Get(id)
	if rn == nil {
		return 0, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	totIx := id.Seq - rn.Val.ID.Seq
	return rn.Val.Text[totIx], nil
}

func (r *Rogue) ContainsID(id ID) bool {
	op := r.OpIndex.Get(id)
	return op != nil
}

func PrintRogueTree(rogue *Rogue) {
	fmt.Printf("Rogue (Tot: %d Vis: %d)\n", rogue.TotSize, rogue.VisSize)
	for _, root := range rogue.Roots {
		PrintTree(root, "", true)
	}
}

// Function to pretty print the tree
func PrintTree(root *FugueNode, prefix string, isTail bool) {
	sideLabel := ""
	switch root.Side {
	case Left:
		sideLabel = "Left"
	case Root:
		sideLabel = "Root"
	case Right:
		sideLabel = "Right"
	}

	if isTail {
		fmt.Printf("%s └──  %s(%d) [%s] %s %s\n", prefix, root.ID.Author, root.ID.Seq, sideLabel, fugueNodeText(root), fugueNodeDel(root))
	} else {
		fmt.Printf("%s ├──  %s(%d) [%s] %s %s\n", prefix, root.ID.Author, root.ID.Seq, sideLabel, fugueNodeText(root), fugueNodeDel(root))
	}

	children := append(root.LeftChildren, root.RightChildren...)
	for i := 0; i < len(children); i++ {
		if isTail {
			PrintTree(children[i], prefix+"    ", i == len(children)-1)
		} else {
			PrintTree(children[i], prefix+"│   ", i == len(children)-1)
		}
	}
}

func fugueNodeText(node *FugueNode) string {
	var stringbuilder strings.Builder
	for i, char := range node.Text {
		str := strings.ReplaceAll(fmt.Sprintf("%#v", Uint16ToStr([]uint16{char})), `"`, "")
		if node.IsDeleted[i] == false {
			stringbuilder.WriteString(str)
			continue
		}

		stringbuilder.WriteString(fmt.Sprintf("\033[41m%v\033[0m", str))
	}
	return stringbuilder.String()
}

func fugueNodeDel(node *FugueNode) string {
	var ids []string
	for _, isDel := range node.IsDeleted {
		if isDel == false {
			ids = append(ids, " ")
		} else {
			ids = append(ids, "x")
		}
	}
	return "[" + strings.Join(ids, ",") + "]"
}

func (r *Rogue) ApplyAction(action Action) (Op, error) {
	mop := MultiOp{}

	switch action := action.(type) {
	case InsertAction:
		op, err := r.Insert(action.Index, action.Text)
		if err != nil {
			return nil, fmt.Errorf("Insert(%d, %q): %w", action.Index, action.Text, err)
		}
		mop = mop.Append(op)
	case DeleteAction:
		dop, err := r.Delete(action.Index, action.Count)
		if err != nil {
			return nil, fmt.Errorf("Delete(%d, %d): %w", action.Index, action.Count, err)
		}
		mop = mop.Append(dop)
	case FormatAction:
		fop, err := r.Format(action.Index, action.Length, action.Format)
		if err != nil {
			return nil, fmt.Errorf("Format(%d, %d, %v): %w", action.Index, action.Length, action.Format, err)
		}
		mop = mop.Append(fop)
	default:
		return nil, fmt.Errorf("unknown action type: %T (Value: %+v)", action, action)
	}

	return FlattenMop(mop), nil
}

func findNewlineIndices(s []uint16) []int {
	var indices []int
	for i, char := range s {
		if char == '\n' {
			indices = append(indices, i)
		}
	}
	return indices
}

func (r *Rogue) isEmptyLine(visIx int) (bool, error) {
	c, err := r.GetChar(visIx)
	if err != nil {
		return false, fmt.Errorf("GetChar(%d): %w", visIx, err)
	}

	if visIx == 0 {
		if c == '\n' {
			return true, nil
		} else {
			return false, nil
		}
	}

	leftChar, err := r.GetChar(visIx - 1)
	if err != nil {
		return false, fmt.Errorf("GetChar(%d): %w", visIx-1, err)
	}

	return leftChar == '\n' && c == '\n', nil
}

func (r *Rogue) CanRedo() bool {
	if r.UndoState == nil {
		return false
	}
	return len(r.UndoState.redoStack) > 0
}

func (r *Rogue) CanUndo() bool {
	isUndoStateGood := true

	if r.UndoState != nil && r.UndoState.address != nil {
		mids := r.UndoState.address.MaxIDs
		_, rootOk := mids["root"]
		_, qOk := mids["q"]

		isUndoStateGood = !(rootOk && qOk && len(mids) == 2)
	}

	return isUndoStateGood && r.LamportClock > 3
}
