package v3

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/stackerr"
	"github.com/teamreviso/code/rogue/v3/set"
)

type Formats struct {
	Sticky   *IntervalTree
	NoSticky *IntervalTree
	Lines    *Lines
}

func NewFormats(rope *Rope) *Formats {
	return &Formats{
		Sticky:   NewIntervalTree(rope, true),
		NoSticky: NewIntervalTree(rope, false),
		Lines:    NewLines(rope),
	}
}

func (f FormatV3Span) SplitSticky() (sticky, noSticky FormatV3Span) {
	sticky, noSticky = FormatV3Span{}, FormatV3Span{}
	noStick := set.NewSet("a", "ql", "qr", "en", "del", "ins")

	for k, v := range f {
		if noStick.Has(k) {
			noSticky[k] = v
		} else {
			sticky[k] = v
		}
	}

	return sticky, noSticky
}

type FormatV3 interface {
	IsSpan() bool
	AsMap() map[string]interface{}
	DropNull() FormatV3
	Empty() bool
	Equals(FormatV3) bool
	Copy() FormatV3
	String() string
}

type FormatV3Span map[string]string

func (f FormatV3Span) IsSpan() bool {
	return true
}

func (f FormatV3Span) AsMap() map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range f {
		out[k] = v
	}

	return out
}

func (f FormatV3Span) DropNull() FormatV3 {
	out := FormatV3Span{}

	for k, v := range f {
		if v != "null" && v != "" {
			out[k] = v
		}
	}

	return out
}

func (f FormatV3Span) Empty() bool {
	return len(f) == 0
}

func (f FormatV3Span) Equals(other FormatV3) bool {
	if _, ok := other.(FormatV3Span); !ok {
		return false
	}

	n := other.(FormatV3Span)

	if len(f) != len(n) {
		return false
	}

	for k, v := range f {
		if n[k] != v {
			return false
		}
	}

	return true
}

func (f FormatV3Span) Copy() FormatV3 {
	out := FormatV3Span{}

	for k, v := range f {
		out[k] = v
	}

	return out
}

func (f FormatV3Span) String() string {
	out := "FormatV3Span{"
	for k, v := range f {
		out += fmt.Sprintf("%s: %s, ", k, v)
	}
	out += "}"
	return out
}

type FormatV3NoStick map[string]string

func (f FormatV3NoStick) IsSpan() bool {
	return true
}

func (f FormatV3NoStick) AsMap() map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range f {
		out[k] = v
	}

	return out
}

func (f FormatV3NoStick) DropNull() FormatV3 {
	out := FormatV3NoStick{}

	for k, v := range f {
		if v != "null" && v != "" {
			out[k] = v
		}
	}

	return out
}

func (f FormatV3NoStick) Empty() bool {
	return len(f) == 0
}

func (f FormatV3NoStick) Equals(other FormatV3) bool {
	if _, ok := other.(FormatV3NoStick); !ok {
		return false
	}

	n := other.(FormatV3NoStick)

	if len(f) != len(n) {
		return false
	}

	for k, v := range f {
		if n[k] != v {
			return false
		}
	}

	return true
}

func (f FormatV3NoStick) Copy() FormatV3 {
	out := FormatV3NoStick{}

	for k, v := range f {
		out[k] = v
	}

	return out
}

func (f FormatV3NoStick) String() string {
	out := "FormatV3NoStick{"
	for k, v := range f {
		out += fmt.Sprintf("%s: %s, ", k, v)
	}
	out += "}"
	return out
}

type FormatV3Header int

func (f FormatV3Header) IsSpan() bool {
	return false
}

func (f FormatV3Header) AsMap() map[string]interface{} {
	return map[string]interface{}{"h": strconv.Itoa(int(f))}
}

func (f FormatV3Header) DropNull() FormatV3 {
	return f
}

func (f FormatV3Header) Empty() bool {
	return false
}

func (f FormatV3Header) Equals(new FormatV3) bool {
	n, ok := new.(FormatV3Header)
	if !ok {
		return false
	}

	return f == n
}

func (f FormatV3Header) String() string {
	return fmt.Sprintf("FormatV3Header(%d)", int(f))
}

func (f FormatV3Header) Copy() FormatV3 {
	return f
}

type FormatV3OrderedList int

func (f FormatV3OrderedList) IsSpan() bool {
	return false
}

func (f FormatV3OrderedList) AsMap() map[string]interface{} {
	return map[string]interface{}{"ol": strconv.Itoa(int(f))}
}

func (f FormatV3OrderedList) DropNull() FormatV3 {
	return f
}

func (f FormatV3OrderedList) Empty() bool {
	return false
}

func (f FormatV3OrderedList) Equals(new FormatV3) bool {
	n, ok := new.(FormatV3OrderedList)
	if !ok {
		return false
	}

	return f == n
}

func (f FormatV3OrderedList) Copy() FormatV3 {
	return f
}

func (f FormatV3OrderedList) String() string {
	return fmt.Sprintf("FormatV3OrderedList(%d)", int(f))
}

type FormatV3BulletList int

func (f FormatV3BulletList) IsSpan() bool {
	return false
}

func (f FormatV3BulletList) AsMap() map[string]interface{} {
	return map[string]interface{}{"ul": strconv.Itoa(int(f))}
}

func (f FormatV3BulletList) DropNull() FormatV3 {
	return f
}

func (f FormatV3BulletList) Empty() bool {
	return false
}

func (f FormatV3BulletList) Equals(other FormatV3) bool {
	n, ok := other.(FormatV3BulletList)
	if !ok {
		return false
	}

	return f == n
}

func (f FormatV3BulletList) Copy() FormatV3 {
	return f
}

func (f FormatV3BulletList) String() string {
	return fmt.Sprintf("FormatV3BulletList(%d)", int(f))
}

type FormatV3IndentedLine int

func (f FormatV3IndentedLine) IsSpan() bool {
	return false
}

func (f FormatV3IndentedLine) AsMap() map[string]interface{} {
	return map[string]interface{}{"il": strconv.Itoa(int(f))}
}

func (f FormatV3IndentedLine) DropNull() FormatV3 {
	return f
}

func (f FormatV3IndentedLine) Empty() bool {
	return false
}

func (f FormatV3IndentedLine) Equals(new FormatV3) bool {
	n, ok := new.(FormatV3IndentedLine)
	if !ok {
		return false
	}

	return f == n
}

func (f FormatV3IndentedLine) Copy() FormatV3 {
	return f
}

func (f FormatV3IndentedLine) String() string {
	return fmt.Sprintf("FormatV3IndentedLine(%d)", int(f))
}

type FormatV3CodeBlock string

func (f FormatV3CodeBlock) IsSpan() bool {
	return false
}

func (f FormatV3CodeBlock) AsMap() map[string]interface{} {
	return map[string]interface{}{"cb": string(f)}
}

func (f FormatV3CodeBlock) DropNull() FormatV3 {
	return f
}

func (f FormatV3CodeBlock) Empty() bool {
	return false
}

func (f FormatV3CodeBlock) Equals(new FormatV3) bool {
	n, ok := new.(FormatV3CodeBlock)
	if !ok {
		return false
	}

	return f == n
}

func (f FormatV3CodeBlock) Copy() FormatV3 {
	return f
}

func (f FormatV3CodeBlock) String() string {
	return fmt.Sprintf("FormatV3CodeBlock(%s)", string(f))
}

type FormatV3BlockQuote struct{}

func (f FormatV3BlockQuote) IsSpan() bool {
	return false
}

func (f FormatV3BlockQuote) AsMap() map[string]interface{} {
	return map[string]interface{}{"bq": "true"}
}

func (f FormatV3BlockQuote) DropNull() FormatV3 {
	return f
}

func (f FormatV3BlockQuote) Empty() bool {
	return false
}

func (f FormatV3BlockQuote) Equals(new FormatV3) bool {
	_, ok := new.(FormatV3BlockQuote)
	return ok
}

func (f FormatV3BlockQuote) Copy() FormatV3 {
	return f
}

func (f FormatV3BlockQuote) String() string {
	return fmt.Sprintf("FormatV3BlockQuote")
}

type FormatV3Line struct{}

func (f FormatV3Line) IsSpan() bool {
	return false
}

func (f FormatV3Line) AsMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (f FormatV3Line) DropNull() FormatV3 {
	return f
}

func (f FormatV3Line) Empty() bool {
	return false
}

func (f FormatV3Line) Equals(new FormatV3) bool {
	_, ok := new.(FormatV3Line)
	return ok
}

func (f FormatV3Line) Copy() FormatV3 {
	return f
}

func (f FormatV3Line) String() string {
	return fmt.Sprintf("FormatV3Line")
}

type FormatV3Rule struct{}

func (f FormatV3Rule) IsSpan() bool {
	return false
}

func (f FormatV3Rule) AsMap() map[string]interface{} {
	return map[string]interface{}{"r": "true"}
}

func (f FormatV3Rule) DropNull() FormatV3 {
	return f
}

func (f FormatV3Rule) Empty() bool {
	return false
}

func (f FormatV3Rule) Equals(new FormatV3) bool {
	_, ok := new.(FormatV3Rule)
	return ok
}

func (f FormatV3Rule) Copy() FormatV3 {
	return f
}

func (f FormatV3Rule) String() string {
	return fmt.Sprintf("FormatV3Rule")
}

type FormatV3Image struct {
	Src    string
	Alt    string
	Height string
	Width  string
}

func (f FormatV3Image) IsSpan() bool {
	return false
}

func (f FormatV3Image) AsMap() map[string]interface{} {
	return map[string]interface{}{"img": f.Src, "alt": f.Alt, "height": f.Height, "width": f.Width}
}

func (f FormatV3Image) DropNull() FormatV3 {
	return f
}

func (f FormatV3Image) Empty() bool {
	return false
}

func (f FormatV3Image) Equals(new FormatV3) bool {
	_, ok := new.(FormatV3Image)
	return ok
}

func (f FormatV3Image) Copy() FormatV3 {
	return f
}

func (f FormatV3Image) String() string {
	return fmt.Sprintf("FormatV3Image{Src: %s, Alt: %s, Height: %s, Width: %s}", f.Src, f.Alt, f.Height, f.Width)
}

type Span struct {
	StartIx int
	EndIx   int
}

func (s Span) Difference(b Span) []Span {
	// Case 1: b is completely outside s
	if b.EndIx < s.StartIx || b.StartIx > s.EndIx {
		return []Span{s}
	}

	// Case 2: b completely covers s
	if b.StartIx <= s.StartIx && b.EndIx >= s.EndIx {
		return nil
	}

	var result []Span

	// Left part of s (if any)
	if s.StartIx < b.StartIx {
		result = append(result, Span{s.StartIx, b.StartIx - 1})
	}

	// Right part of s (if any)
	if b.EndIx < s.EndIx {
		result = append(result, Span{b.EndIx + 1, s.EndIx})
	}

	return result
}

func (s Span) Intersection(b Span) *Span {
	// No overlap
	if b.EndIx < s.StartIx || s.EndIx < b.StartIx {
		return nil
	}

	// Calculate the intersection
	start := max(s.StartIx, b.StartIx)
	end := min(s.EndIx, b.EndIx)

	// Return the intersection span
	return &Span{start, end}
}

func (formats *Formats) Insert(formatOp FormatOp) (err error) {
	switch f := formatOp.Format.(type) {
	case FormatV3Span:
		sticky, noSticky := f.SplitSticky()

		if len(sticky) > 0 {
			fop := FormatOp{
				ID:      formatOp.ID,
				StartID: formatOp.StartID,
				EndID:   formatOp.EndID,
				Format:  sticky,
			}

			err = formats.Sticky.Insert(fop)
			if err != nil {
				return err
			}
		}

		if len(noSticky) > 0 {
			fop := FormatOp{
				ID:      formatOp.ID,
				StartID: formatOp.StartID,
				EndID:   formatOp.EndID,
				Format:  noSticky,
			}

			err = formats.NoSticky.Insert(fop)
			if err != nil {
				return err
			}
		}
	default:
		err = formats.Lines.Insert(formatOp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Rogue) _invertFormatOpSpans(targetID ID, op FormatOp, isSticky bool) (iop Op, err error) {
	_, startIx, err := r.Rope.GetIndex(op.StartID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := r.Rope.GetIndex(op.EndID)
	if err != nil {
		return nil, err
	}

	if f, ok := op.Format.(FormatV3Span); ok {
		_, nostick := f.SplitSticky()
		if len(nostick) == 0 {
			endIx--
		}
	}

	mop := MultiOp{}
	spans := []Span{{StartIx: startIx, EndIx: endIx}}

	for len(spans) > 0 {
		s := spans[0]
		spans = spans[1:]

		var fop *FormatOp

		if isSticky {
			fop, err = r.Formats.Sticky.MaxFormatBefore(s.StartIx, s.EndIx, targetID)
			if err != nil {
				return nil, err
			}
		} else {
			fop, err = r.Formats.NoSticky.MaxFormatBefore(s.StartIx, s.EndIx, targetID)
			if err != nil {
				return nil, err
			}
		}

		if fop == nil {
			// there's no format covering this span
			// add empty op to the mop

			startID, err := r.Rope.GetTotID(s.StartIx)
			if err != nil {
				return nil, err
			}

			if isSticky {
				s.EndIx = min(r.TotSize-1, s.EndIx+1)
			}

			endID, err := r.Rope.GetTotID(s.EndIx)
			if err != nil {
				return nil, err
			}

			var spanFormat FormatV3Span
			if isSticky {
				spanFormat = FormatV3Span{"e": "true"}
			} else {
				spanFormat = FormatV3Span{"en": "true"}
			}

			mop = mop.Append(FormatOp{
				ID:      r.NextID(1),
				StartID: startID,
				EndID:   endID,
				Format:  spanFormat,
			})

			continue
		}

		_, fStartIx, err := r.Rope.GetIndex(fop.StartID)
		if err != nil {
			return nil, err
		}

		_, fEndIx, err := r.Rope.GetIndex(fop.EndID)
		if err != nil {
			return nil, err
		}

		fs := Span{StartIx: fStartIx, EndIx: fEndIx}
		spanFormat, ok := fop.Format.Copy().(FormatV3Span)
		if !ok {
			log.Warnf("line formats shouldn't be in the overlapping list: %v", fop.Format)
			continue
		}

		if isSticky {
			fs.EndIx--
			spanFormat["e"] = "true"
		} else {
			spanFormat["en"] = "true"
		}

		intersection := fs.Intersection(s)
		if intersection == nil {
			return nil, stackerr.Errorf("no intersection between %v and %v", fs, s)
		}

		diffs := s.Difference(*intersection)
		spans = append(spans, diffs...)

		startID, err := r.Rope.GetTotID(intersection.StartIx)
		if err != nil {
			return nil, fmt.Errorf("GetTotID(%v): %w", intersection.StartIx, err)
		}

		if isSticky {
			intersection.EndIx++
		}

		endID, err := r.Rope.GetTotID(intersection.EndIx)
		if err != nil {
			return nil, fmt.Errorf("GetTotID(%v): %w", intersection.EndIx, err)
		}

		mop = mop.Append(FormatOp{
			ID:      r.NextID(1),
			StartID: startID,
			EndID:   endID,
			Format:  spanFormat,
		})
	}

	return FlattenMop(mop), nil
}

func (r *Rogue) invertFormatOp(targetID ID, op FormatOp) (iop Op, err error) {
	// FORMAT IS A LINE
	if !op.Format.IsSpan() {
		_, ix, err := r.Rope.GetIndex(op.StartID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.StartID, err)
		}

		line, err := r.Formats.Lines.Tree.Get(ix)
		if err != nil {
			return nil, fmt.Errorf("Formats.Lines.Tree.Get(%v): %w", ix, err)
		}

		// toID := ID{op.ID.Author, op.ID.Seq - 1}
		toID := ID{targetID.Author, targetID.Seq - 1}
		lop, _ := line.Formats.FindLeftSibNode(toID)

		out := FormatOp{
			ID:      r.NextID(1),
			StartID: op.StartID,
			EndID:   op.StartID,
		}

		if lop == nil {
			out.Format = FormatV3Line{}
		} else {
			out.Format = lop.Value.Format
		}

		return out, nil
	}

	// FORMAT IS A SPAN
	mop := MultiOp{}

	stickyOps, err := r._invertFormatOpSpans(targetID, op, true)
	if err != nil {
		return nil, err
	}
	mop = mop.Append(stickyOps)

	noStickyOps, err := r._invertFormatOpSpans(targetID, op, false)
	if err != nil {
		return nil, err
	}
	mop = mop.Append(noStickyOps)

	return FlattenMop(mop), nil
}

func (r *Rogue) _rewindFormatOpSpans(startIx, endIx int, address *ContentAddress, isSticky bool, opID ID) (ops []FormatOp, err error) {
	spans := []Span{{StartIx: startIx, EndIx: endIx}}

	for len(spans) > 0 {
		s := spans[0]
		spans = spans[1:]

		var fop *FormatOp

		if isSticky {
			fop, err = r.Formats.Sticky.MaxFormatAt(s.StartIx, s.EndIx, address)
			if err != nil {
				return nil, err
			}
		} else {
			fop, err = r.Formats.NoSticky.MaxFormatAt(s.StartIx, s.EndIx, address)
			if err != nil {
				return nil, err
			}
		}

		if fop == nil {
			// there's no format covering this span
			// add empty op to the mop

			startID, err := r.Rope.GetTotID(s.StartIx)
			if err != nil {
				return nil, err
			}

			var spanFormat FormatV3Span
			if isSticky {
				s.EndIx = min(r.TotSize-1, s.EndIx+1)
				spanFormat = FormatV3Span{"e": "true"}
			} else {
				spanFormat = FormatV3Span{"en": "true"}
			}

			endID, err := r.Rope.GetTotID(s.EndIx)
			if err != nil {
				return nil, err
			}

			ops = append(ops, FormatOp{
				ID:      opID,
				StartID: startID,
				EndID:   endID,
				Format:  spanFormat,
			})

			continue
		}

		_, fStartIx, err := r.Rope.GetIndex(fop.StartID)
		if err != nil {
			return nil, err
		}

		_, fEndIx, err := r.Rope.GetIndex(fop.EndID)
		if err != nil {
			return nil, err
		}

		fs := Span{StartIx: fStartIx, EndIx: fEndIx}
		spanFormat, ok := fop.Format.Copy().(FormatV3Span)
		if !ok {
			log.Warnf("line formats shouldn't be in the overlapping list: %v", fop.Format)
			continue
		}

		if isSticky {
			fs.EndIx--
			spanFormat["e"] = "true"
		} else {
			spanFormat["en"] = "true"
		}

		intersection := fs.Intersection(s)
		if intersection == nil {
			return nil, stackerr.Errorf("no intersection between %v and %v", fs, s)
		}

		diffs := s.Difference(*intersection)
		spans = append(spans, diffs...)

		startID, err := r.Rope.GetTotID(intersection.StartIx)
		if err != nil {
			return nil, err
		}

		if isSticky {
			intersection.EndIx = min(r.TotSize-1, intersection.EndIx+1)
		}

		endID, err := r.Rope.GetTotID(intersection.EndIx)
		if err != nil {
			return nil, err
		}

		ops = append(ops, FormatOp{
			ID:      opID,
			StartID: startID,
			EndID:   endID,
			Format:  spanFormat,
		})
	}

	return ops, nil
}

func (r *Rogue) rewindFormatTo(startID, endID ID, address *ContentAddress, opID ID) (ops []FormatOp, err error) {
	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	addrStartIx, addrEndIx := 0, r.TotSize
	if address != nil {
		_, addrStartIx, err = r.Rope.GetIndex(address.StartID)
		if err != nil {
			return nil, err
		}

		_, addrEndIx, err = r.Rope.GetIndex(address.EndID)
		if err != nil {
			return nil, err
		}

		if addrEndIx < addrStartIx {
			addrStartIx, addrEndIx = 0, r.TotSize
		}
	}

	// REWIND SPAN FORMATS
	ix0, ix1 := startIx, endIx

	if startIx < addrStartIx && addrStartIx < endIx {
		stickyOps, err := r._rewindFormatOpSpans(startIx, addrStartIx-1, nil, true, opID)
		if err != nil {
			return nil, err
		}
		ops = append(ops, stickyOps...)

		noStickyOps, err := r._rewindFormatOpSpans(startIx, addrStartIx-1, nil, false, opID)
		if err != nil {
			return nil, err
		}
		ops = append(ops, noStickyOps...)

		ix0 = addrStartIx
	}

	if addrEndIx < endIx && startIx < addrEndIx {
		stickyOps, err := r._rewindFormatOpSpans(addrEndIx+1, endIx, nil, true, opID)
		if err != nil {
			return nil, err
		}
		ops = append(ops, stickyOps...)

		noStickyOps, err := r._rewindFormatOpSpans(addrEndIx+1, endIx, nil, false, opID)
		if err != nil {
			return nil, err
		}
		ops = append(ops, noStickyOps...)

		ix1 = addrEndIx
	}

	stickyOps, err := r._rewindFormatOpSpans(ix0, ix1, address, true, opID)
	if err != nil {
		return nil, err
	}
	ops = append(ops, stickyOps...)

	noStickyOps, err := r._rewindFormatOpSpans(ix0, ix1, address, false, opID)
	if err != nil {
		return nil, err
	}
	ops = append(ops, noStickyOps...)

	// REWIND LINE FORMATS
	err = r.Formats.Lines.Tree.Slice(startIx, endIx, func(lh *LineHistory) error {
		n := lh.Formats.MaxNode()

		_, ix, err := r.Rope.GetIndex(lh.TargetID)
		if err != nil {
			return err
		}

		if address == nil || ix < addrStartIx || addrEndIx < ix {
			ops = append(ops, FormatOp{
				ID:      opID,
				StartID: lh.TargetID,
				EndID:   lh.TargetID,
				Format:  n.Value.Format,
			})
		} else {
			for n != nil {
				if address.Contains(n.Value.ID) {
					ops = append(ops, FormatOp{
						ID:      opID,
						StartID: lh.TargetID,
						EndID:   lh.TargetID,
						Format:  n.Value.Format,
					})
					break
				}

				n = n.StepLeft()
			}

			if n == nil {
				ops = append(ops, FormatOp{
					ID:      opID,
					StartID: lh.TargetID,
					EndID:   lh.TargetID,
					Format:  FormatV3Line{},
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ops, nil
}

func (r *Rogue) SetNOS(startID, endID ID, address *ContentAddress) error {
	ops, err := r.rewindFormatTo(startID, endID, address, NoID)
	if err != nil {
		return err
	}

	for _, op := range ops {
		_, err = r.NOS.Insert(op)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Rogue) ResetNOS() error {
	r.NOS = NewNOSV2(r.Rope)

	firstID, err := r.GetFirstID()
	if err != nil {
		return nil // Just means the doc is empty
	}

	lastID, err := r.GetLastID()
	if err != nil {
		return nil // Just means the doc is empty
	}

	return r.SetNOS(firstID, lastID, nil)
}

func (formats *Formats) SearchOverlapping(startID, endID ID) (ops []FormatOp, err error) {
	_, startIx, err := formats.Lines.Rope.GetIndex(startID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := formats.Lines.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	sticky, err := formats.Sticky.SearchOverlapping(startIx, endIx)
	if err != nil {
		return nil, err
	}
	ops = append(ops, sticky...)

	noSticky, err := formats.NoSticky.SearchOverlapping(startIx, endIx)
	if err != nil {
		return nil, err
	}
	ops = append(ops, noSticky...)

	formats.Lines.Tree.Slice(startIx, endIx, func(lh *LineHistory) error {
		lh.Formats.Dft(func(op FormatOp) error {
			ops = append(ops, op)
			return nil
		})
		return nil
	})

	slices.SortFunc(ops, func(a, b FormatOp) int {
		aSeq, aAuthor := a.ID.Seq, a.ID.Author
		bSeq, bAuthor := b.ID.Seq, b.ID.Author

		if aSeq < bSeq {
			return -1
		} else if aSeq > bSeq {
			return 1
		} else if aAuthor < bAuthor {
			return -1
		} else if aAuthor > bAuthor {
			return 1
		} else {
			return 0
		}
	})

	return ops, nil
}

func (formats *Formats) Validate() error {
	err := formats.Sticky.Validate()
	if err != nil {
		return fmt.Errorf("Sticky.ValidateSpans(): %w", err)
	}

	err = formats.NoSticky.Validate()
	if err != nil {
		return fmt.Errorf("NoSticky.ValidateSpans(): %w", err)
	}

	err = formats.Lines.Validate()
	if err != nil {
		return fmt.Errorf("Lines.Validate(): %w", err)
	}

	return nil
}

func (r *Rogue) insertFormatOp(fop FormatOp) error {
	err := r.ValidateFormat(fop)
	if err != nil {
		log.Errorf("ValidateFormat: %v", err)
		return nil
	}

	err = r.Formats.Insert(fop)
	if err != nil {
		return err
	}

	err = r.SetNOS(fop.StartID, fop.EndID, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rope) printTotIxFormatOp(op FormatOp) {
	_, startIx, err := r.GetIndex(op.StartID)
	if err != nil {
		fmt.Printf("GetIndex(%v): %v\n", op.StartID, err)
		return
	}

	_, endIx, err := r.GetIndex(op.EndID)
	if err != nil {
		fmt.Printf("GetIndex(%v): %v\n", op.EndID, err)
		return
	}

	fmt.Printf("%s %d %d %v\n", op.ID, startIx, endIx, op.Format)
}

func (r *Rope) printVisIxFormatOp(op Op) {
	switch op := op.(type) {
	case FormatOp:
		startID, err := r.NearestVisRightOf(op.StartID)
		if err != nil {
			fmt.Printf("GetNearestVisRight(%v): %v\n", op.StartID, err)
			return
		}

		startIx, _, err := r.GetIndex(startID)
		if err != nil {
			fmt.Printf("GetIndex(%v): %v\n", startID, err)
			return
		}

		endID, err := r.NearestVisLeftOf(op.EndID)
		if err != nil {
			fmt.Printf("GetNearestVisLeft(%v): %v\n", op.EndID, err)
			return
		}

		endIx, _, err := r.GetIndex(endID)
		if err != nil {
			fmt.Printf("GetIndex(%v): %v\n", endID, err)
			return
		}

		fmt.Printf("%s %d %d %v\n", op.ID, startIx, endIx, op.Format)
	case MultiOp:
		for _, op := range op.Mops {
			r.printVisIxFormatOp(op)
		}
	}
}
