package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"

	"github.com/charmbracelet/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/teamreviso/code/pkg/stackerr"
)

type DeltaType int

const (
	DeltaTypeInsert DeltaType = iota
	DeltaTypeDelete
	DeltaTypeEqual
)

func (d DeltaType) String() string {
	return [...]string{"Insert", "Delete", "Equal"}[d]
}

type FugueDiff struct {
	Text       []uint16
	TotIxs     []int
	IDs        []ID
	DeltaTypes []DeltaType
}

func (f *FugueDiff) String() string {
	return fmt.Sprintf(
		"FugueVisIx{Text: %v, TotIxs: %v, IDs: %v, DeltaTypes: %v}",
		Uint16ToStr(f.Text),
		f.TotIxs,
		f.IDs,
		f.DeltaTypes,
	)
}

type FilteredNOS struct {
	*FugueDiff
	spanNos *NOS
	lineNos *NOS
}

func (f *FilteredNOS) String() string {
	return fmt.Sprintf(
		"FilteredNOS{Text: %v, spanNos: %v, lineNos: %v}",
		Uint16ToStr(f.Text),
		f.spanNos,
		f.lineNos,
	)
}

func (r *Rogue) Filter(startID, endID ID, addr *ContentAddress) (*FugueDiff, error) {
	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	if startIx > endIx {
		return nil, nil
	}

	addrStartIx, addrEndIx := 0, r.TotSize

	if addr != nil {
		_, addrStartIx, err = r.Rope.GetIndex(addr.StartID)
		if err != nil {
			return nil, err
		}

		_, addrEndIx, err = r.Rope.GetIndex(addr.EndID)
		if err != nil {
			return nil, err
		}
	}

	out := &FugueDiff{
		Text:       make([]uint16, 0, endIx-startIx),
		TotIxs:     make([]int, 0, endIx-startIx),
		IDs:        make([]ID, 0, endIx-startIx),
		DeltaTypes: make([]DeltaType, 0, endIx-startIx),
	}

	tot, err := r.Rope.GetTotBetween(startID, endID)
	if err != nil {
		return nil, err
	}

	for ix, id := range tot.IDs {
		totIx := startIx + ix

		if addr == nil || totIx < addrStartIx || totIx > addrEndIx {
			if !tot.IsDeleted[ix] {
				out.Text = append(out.Text, tot.Text[ix])
				out.TotIxs = append(out.TotIxs, totIx)
				out.IDs = append(out.IDs, id)
				out.DeltaTypes = append(out.DeltaTypes, DeltaTypeEqual)
			}
			continue
		}

		if addr.Contains(id) {
			ch := r.CharHistory[id]
			if ch == nil || !addr.Contains(ch.Min().ID) {
				out.Text = append(out.Text, tot.Text[ix])
				out.TotIxs = append(out.TotIxs, totIx)
				out.IDs = append(out.IDs, id)
				if tot.IsDeleted[ix] {
					out.DeltaTypes = append(out.DeltaTypes, DeltaTypeInsert)
				} else {
					out.DeltaTypes = append(out.DeltaTypes, DeltaTypeEqual)
				}
			} else {
				ch.ReverseDft(func(m *Marker) error {
					if addr.Contains(m.ID) {
						if !m.IsDel {
							out.Text = append(out.Text, tot.Text[ix])
							out.TotIxs = append(out.TotIxs, totIx)
							out.IDs = append(out.IDs, id)
							if tot.IsDeleted[ix] {
								out.DeltaTypes = append(out.DeltaTypes, DeltaTypeInsert)
							} else {
								out.DeltaTypes = append(out.DeltaTypes, DeltaTypeEqual)
							}
						}
						return ErrorStopIteration{}
					}

					return nil
				})
			}
		}
	}

	return out, nil
}

func (r *Rogue) GetFilteredNOSBetween(startID, endID ID, fromAddress, toAddress *ContentAddress, smartQuote bool) (*FilteredNOS, error) {
	filtered, err := r.Filter(startID, endID, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("Filter(%v, %v, %v): %w", startID, endID, fromAddress, err)
	}

	curVis, spanNOS, lineNOS, err := r.ToIndexNos(startID, endID, toAddress, smartQuote)
	if err != nil {
		return nil, fmt.Errorf("r.ToIndexNos(%v, %v, %v): %w", startID, endID, toAddress, err)
	}
	if curVis == nil {
		return nil, fmt.Errorf("curVis is nil")
	}

	diffs := DiffWords(Uint16ToStr(filtered.Text), Uint16ToStr(curVis.Text))

	length := 0
	for _, diff := range diffs {
		length += len(diff.Text)
	}
	diffText := FugueDiff{
		Text:       make([]uint16, 0, length),
		TotIxs:     make([]int, 0, length),
		IDs:        make([]ID, 0, length),
		DeltaTypes: make([]DeltaType, 0, length),
	}
	diffSpan, diffLine := NewNOS(), NewNOS()

	visIx, delIx, diffIx := 0, 0, 0
	visToDiffIx := map[int]int{}
	var diffNOSNodes []NOSNode
	for _, diff := range diffs {
		uText := StrToUint16(diff.Text)
		var updateFunc func()
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			updateFunc = func() {
				visToDiffIx[visIx] = diffIx
				diffText.IDs = append(diffText.IDs, curVis.IDs[visIx])
				diffText.TotIxs = append(diffText.TotIxs, curVis.TotIxs[visIx])
				diffText.DeltaTypes = append(diffText.DeltaTypes, DeltaTypeEqual)
				visIx++
				delIx++
			}
		case diffmatchpatch.DiffInsert:
			diffNOSNodes = append(diffNOSNodes, NOSNode{
				StartIx: diffIx,
				EndIx:   diffIx + len(uText) - 1,
				Format:  FormatV3Span{"ins": "true"},
			})

			updateFunc = func() {
				visToDiffIx[visIx] = diffIx
				diffText.IDs = append(diffText.IDs, curVis.IDs[visIx])
				diffText.TotIxs = append(diffText.TotIxs, curVis.TotIxs[visIx])
				diffText.DeltaTypes = append(diffText.DeltaTypes, DeltaTypeInsert)
				visIx++
			}
		case diffmatchpatch.DiffDelete:
			diffNOSNodes = append(diffNOSNodes, NOSNode{
				StartIx: diffIx,
				EndIx:   diffIx + len(uText) - 1,
				Format:  FormatV3Span{"del": "true", "noid": "true"},
			})
			updateFunc = func() {
				diffText.IDs = append(diffText.IDs, filtered.IDs[delIx])
				diffText.TotIxs = append(diffText.TotIxs, filtered.TotIxs[delIx])
				diffText.DeltaTypes = append(diffText.DeltaTypes, DeltaTypeDelete)
				delIx++
			}
		}

		for _, c := range uText {
			diffText.Text = append(diffText.Text, c)
			updateFunc()
			diffIx++
		}
	}

	spanNOS.tree.Dft(func(node *NOSNode) error {
		startIx, ok := visToDiffIx[node.StartIx]
		if !ok {
			log.Warnf("node.StartIx %v not found in visToDiffIx", node.StartIx)
			return nil
		}
		endIx, ok := visToDiffIx[node.EndIx]
		if !ok {
			log.Warnf("node.StartIx %v not found in visToDiffIx", node.StartIx)
			return nil
		}

		diffSpan.Insert(NOSNode{
			StartIx: startIx,
			EndIx:   endIx,
			Format:  node.Format,
		})

		return nil
	})

	for _, node := range diffNOSNodes {
		diffSpan.Insert(node)
	}

	/*lineNOS.tree.Dft(func(node *NOSNode) error {
		startIx := visToDiffIx[node.StartIx]
		endIx := visToDiffIx[node.EndIx]
		fmt.Printf("startIx: %d, endIx: %d, lineNOS: %v\n", startIx, endIx, node)
		return nil
	})*/

	ix := 0
	lineNOS.tree.Dft(func(node *NOSNode) error {
		startIx := visToDiffIx[node.StartIx]
		endIx := visToDiffIx[node.EndIx]

		if ix < startIx {
			startIx = ix
		}

		diffLine.Insert(NOSNode{
			StartIx: startIx,
			EndIx:   endIx,
			Format:  node.Format,
		})
		ix = endIx + 1

		return nil
	})

	if ix < len(diffText.Text)-1 {
		diffLine.Insert(NOSNode{
			StartIx: ix,
			EndIx:   len(diffText.Text) - 1,
			Format:  FormatV3Span{"del": "true", "noid": "true"},
		})
	}

	// spanNOS = spanNOS.mergeNeighbors()

	return &FilteredNOS{
		FugueDiff: &diffText,
		spanNos:   diffSpan,
		lineNos:   diffLine,
	}, nil
}

func (r *Rogue) GetFilteredNOS(startID, endID ID, address *ContentAddress, smartQuote bool) (*FilteredNOS, error) {
	return r.GetFilteredNOSBetween(startID, endID, address, nil, smartQuote)
}

type IDItem struct {
	ID   ID
	Char uint16
}

func (r *Rogue) WalkRightFromAt(startID ID, addr *ContentAddress) iter.Seq2[*IDItem, error] {
	return func(yield func(idItem *IDItem, err error) bool) {
		_, ix, err := r.Rope.GetIndex(startID)
		if err != nil {
			yield(nil, err)
			return
		}

		addrStartIx, addrEndIx := 0, r.TotSize
		if addr != nil {
			_, addrStartIx, err = r.Rope.GetIndex(addr.StartID)
			if err != nil {
				yield(nil, err)
				return
			}

			_, addrEndIx, err = r.Rope.GetIndex(addr.EndID)
			if err != nil {
				yield(nil, err)
				return
			}
		}

		node := r.Rope.Index.Get(startID)
		if node == nil {
			yield(nil, stackerr.New(fmt.Errorf("node %v doesn't exist", startID)))
			return
		}

		for {
			isStart := node.Val.ContainsID(startID)

			fn := node.Val.Explode()
			for i, id := range fn.IDs {
				if isStart && id.Seq < startID.Seq {
					continue
				}

				var idItem *IDItem
				if addr == nil || ix < addrStartIx || addrEndIx < ix {
					if !fn.IsDeleted[i] {
						idItem = &IDItem{id, fn.Text[i]}
					}
				} else if addr.Contains(id) {
					ch := r.CharHistory[id]
					if ch == nil || !addr.Contains(ch.Min().ID) {
						idItem = &IDItem{id, fn.Text[i]}
					} else {
						err := ch.ReverseDft(func(m *Marker) error {
							if addr.Contains(m.ID) {
								if !m.IsDel {
									idItem = &IDItem{id, fn.Text[i]}
								}
								return ErrorStopIteration{}
							}

							return nil
						})

						if err != nil {
							if !errors.As(err, &ErrorStopIteration{}) {
								yield(nil, stackerr.New(err))
								return
							}
						}
					}
				}

				if idItem != nil && !yield(idItem, nil) {
					return
				}

				ix++
			}

			if addr == nil {
				node, err = node.RightVisSibling()
				if err != nil {
					if errors.As(err, &ErrorNoRightVisSibling{}) {
						return
					}
					yield(nil, err)
					return
				}
			} else {
				node, err = node.RightTotSibling()
				if err != nil {
					if errors.As(err, &ErrorNoRightTotSibling{}) {
						return
					}
					yield(nil, err)
					return
				}
			}
		}
	}
}

func (r *Rogue) WalkLeftFromAt(startID ID, addr *ContentAddress) iter.Seq2[*IDItem, error] {
	return func(yield func(idItem *IDItem, err error) bool) {
		_, ix, err := r.Rope.GetIndex(startID)
		if err != nil {
			yield(nil, err)
			return
		}

		addrStartIx, addrEndIx := 0, r.TotSize
		if addr != nil {
			_, addrStartIx, err = r.Rope.GetIndex(addr.StartID)
			if err != nil {
				yield(nil, err)
				return
			}

			_, addrEndIx, err = r.Rope.GetIndex(addr.EndID)
			if err != nil {
				yield(nil, err)
				return
			}
		}

		node := r.Rope.Index.Get(startID)
		if node == nil {
			yield(nil, stackerr.New(fmt.Errorf("node %v doesn't exist", startID)))
			return
		}

		for {
			isStart := node.Val.ContainsID(startID)

			fn := node.Val.Explode()
			for i := len(fn.IDs) - 1; i >= 0; i-- {
				id := fn.IDs[i]
				if isStart && id.Seq > startID.Seq {
					continue
				}

				var idItem *IDItem
				if addr == nil || ix < addrStartIx || addrEndIx < ix {
					if !fn.IsDeleted[i] {
						idItem = &IDItem{id, fn.Text[i]}
					}
				} else if addr.Contains(id) {
					ch := r.CharHistory[id]
					if ch == nil || !addr.Contains(ch.Min().ID) {
						idItem = &IDItem{id, fn.Text[i]}
					} else {
						err := ch.ReverseDft(func(m *Marker) error {
							if addr.Contains(m.ID) {
								if !m.IsDel {
									idItem = &IDItem{id, fn.Text[i]}
								}
								return ErrorStopIteration{}
							}

							return nil
						})

						if err != nil {
							if !errors.As(err, &ErrorStopIteration{}) {
								yield(nil, stackerr.New(err))
								return
							}
						}
					}
				}

				if idItem != nil && !yield(idItem, nil) {
					return
				}

				ix--
			}

			if addr == nil {
				node, err = node.LeftVisSibling()
				if err != nil {
					if errors.As(err, &ErrorNoLeftVisSibling{}) {
						return
					}
					yield(nil, err)
					return
				}
			} else {
				node, err = node.LeftTotSibling()
				if err != nil {
					if errors.As(err, &ErrorNoLeftTotSibling{}) {
						return
					}
					yield(nil, err)
					return
				}
			}
		}
	}
}

type TotIDItem struct {
	ID    ID
	TotIx int
	Char  uint16
}

func (r *Rogue) WalkRightFromTot(startID ID) iter.Seq2[*TotIDItem, error] {
	return func(yield func(idItem *TotIDItem, err error) bool) {
		_, ix, err := r.Rope.GetIndex(startID)
		if err != nil {
			yield(nil, err)
			return
		}

		node := r.Rope.Index.Get(startID)
		if node == nil {
			yield(nil, stackerr.New(fmt.Errorf("node %v doesn't exist", startID)))
			return
		}

		for {
			isStart := node.Val.ContainsID(startID)

			fn := node.Val.Explode()
			for i, id := range fn.IDs {
				if isStart && id.Seq < startID.Seq {
					continue
				}

				totIDItem := &TotIDItem{
					ID:    id,
					TotIx: ix,
					Char:  fn.Text[i],
				}

				if !yield(totIDItem, nil) {
					return
				}

				ix++
			}

			node, err = node.RightTotSibling()
			if err != nil {
				if errors.As(err, &ErrorNoRightTotSibling{}) {
					return
				}
				yield(nil, err)
				return
			}
		}
	}
}

func (r *Rogue) WalkLeftFromTot(startID ID) iter.Seq2[*TotIDItem, error] {
	return func(yield func(idItem *TotIDItem, err error) bool) {
		_, ix, err := r.Rope.GetIndex(startID)
		if err != nil {
			yield(nil, err)
			return
		}

		node := r.Rope.Index.Get(startID)
		if node == nil {
			yield(nil, stackerr.New(fmt.Errorf("node %v doesn't exist", startID)))
			return
		}

		for {
			isStart := node.Val.ContainsID(startID)

			fn := node.Val.Explode()
			for i := len(fn.IDs) - 1; i >= 0; i-- {
				id := fn.IDs[i]
				if isStart && id.Seq > startID.Seq {
					continue
				}

				totIDItem := &TotIDItem{
					ID:    id,
					TotIx: ix,
					Char:  fn.Text[i],
				}

				if !yield(totIDItem, nil) {
					return
				}

				ix--
			}

			node, err = node.LeftTotSibling()
			if err != nil {
				if errors.As(err, &ErrorNoLeftTotSibling{}) {
					return
				}
				yield(nil, err)
				return
			}
		}
	}
}

func (r *Rogue) IDFromIDAndOffset(startID ID, offset int, addr *ContentAddress) (ID, error) {
	if offset == 0 {
		return startID, nil
	}

	if addr == nil {
		visIx, _, err := r.Rope.GetIndex(startID)
		if err != nil {
			return ID{}, fmt.Errorf("GetVisIndex(%v): %w", startID, err)
		}

		id, err := r.Rope.GetVisID(visIx + offset)
		if err != nil {
			return ID{}, fmt.Errorf("GetVisID(%v): %w", visIx+offset, err)
		}

		return id, nil
	}

	var id ID
	for v, err := range r.WalkRightFromAt(startID, addr) {
		if err != nil {
			return ID{}, err
		}

		if offset == 0 {
			id = v.ID
		}

		offset--
	}

	return id, nil
}

func (r *Rogue) IDsToEnclosingSpan(ids []ID, address *ContentAddress) (startID, endID ID, err error) {
	totIxs := make([]int, 0, len(ids))
	for _, id := range ids {
		_, totIx, err := r.Rope.GetIndex(id)
		if err != nil {
			return startID, endID, err
		}

		totIxs = append(totIxs, totIx)
	}

	minIx := SliceMinIx(totIxs)
	maxIx := SliceMaxIx(totIxs)

	startID = ids[minIx]
	endID = ids[maxIx]
	var lastVis *IDItem

	sid, err := r.VisLeftOf(startID)
	if err != nil {
		if !errors.As(err, &ErrorNoLeftVisSibling{}) {
			return startID, endID, err
		}
		// leave startID alone if it's the first id in the doc
		sid = startID
	}

	for v, err := range r.WalkLeftFromAt(sid, address) {
		if err != nil {
			return startID, endID, err
		}

		lastVis = v

		if v.Char == '\n' {
			break
		}
	}

	if lastVis != nil {
		startID = lastVis.ID
		if lastVis.Char == '\n' {
			startID, err = r.VisRightOf(lastVis.ID)
			if err != nil {
				return startID, endID, err
			}
		}
	}

	lastVis = nil
	for v, err := range r.WalkRightFromAt(endID, address) {
		if err != nil {
			return startID, endID, err
		}

		lastVis = v

		if v.Char == '\n' {
			break
		}
	}

	if lastVis == nil {
		ca, err := json.Marshal(address)
		if err != nil {
			return startID, endID, err
		}

		html, err := r.GetHtmlAt(startID, endID, address, true, false)
		if err != nil {
			return startID, endID, err
		}

		return startID, endID, stackerr.New(fmt.Errorf("ID: %v and all right siblings are deleted at content address: %s with html: %q", endID, ca, html))
	}

	endID = lastVis.ID

	return startID, endID, nil
}

func (r *Rogue) NearestAt(id ID, address *ContentAddress) (ID, error) {
	isDel, err := r.IsDeletedAt(id, address)
	if err != nil {
		return NoID, err
	}

	if !isDel {
		return id, nil
	}

	rid, err := r.RightOfAt(id, address)
	if err != nil {
		if !errors.As(err, &ErrorNoRightSiblingAt{}) {
			return NoID, err
		}

		lid, err := r.LeftOfAt(id, address)
		if err != nil {
			if errors.As(err, &ErrorNoLeftSiblingAt{}) {
				return NoID, stackerr.Errorf("doc is empty at address %v", address)
			}

			return NoID, err
		}

		return lid, nil
	}

	return rid, nil
}

func (r *Rogue) GetLineAt(id ID, address *ContentAddress) (startID, endID ID, offset int, err error) {
	if r.VisSize == 0 {
		return NoID, NoID, -1, stackerr.Errorf("doc is empty")
	}

	id, err = r.NearestAt(id, address)
	if err != nil {
		return NoID, NoID, -1, err
	}

	var lastVis *IDItem
	isFirst := true
	for v, err := range r.WalkLeftFromAt(id, address) {
		if err != nil {
			return startID, endID, -1, err
		}

		if isFirst {
			isFirst = false
		} else {
			if v.Char == '\n' {
				break
			}

			lastVis = v
			offset++
		}
	}

	if lastVis == nil {
		startID = id
	} else {
		startID = lastVis.ID
	}

	lastVis, endID = nil, startID
	for v, err := range r.WalkRightFromAt(id, address) {
		if err != nil {
			return startID, endID, -1, err
		}

		lastVis = v

		if v.Char == '\n' {
			break
		}
	}

	if lastVis != nil {
		endID = lastVis.ID
	}

	return startID, endID, offset, nil
}

func (r *Rogue) ScanLeftOfAt(id ID, value uint16, address *ContentAddress) (*ID, error) {
	for v, err := range r.WalkLeftFromAt(id, address) {
		if err != nil {
			return nil, err
		}

		if v.Char == value {
			return &v.ID, nil
		}
	}

	return nil, nil
}

func (r *Rogue) ScanRightOfAt(id ID, value uint16, address *ContentAddress) (*ID, error) {
	for v, err := range r.WalkRightFromAt(id, address) {
		if err != nil {
			return nil, err
		}

		if v.Char == value {
			return &v.ID, nil
		}
	}

	return nil, nil
}

func (r *Rogue) RightOfAt(id ID, address *ContentAddress) (ID, error) {
	isFirst := true
	for v, err := range r.WalkRightFromAt(id, address) {
		if err != nil {
			return NoID, err
		}

		if isFirst && v.ID == id {
			isFirst = false
			continue
		}

		return v.ID, nil
	}

	ca, err := json.Marshal(address)
	if err != nil {
		return NoID, stackerr.New(err)
	}

	return NoID, stackerr.New(ErrorNoRightSiblingAt{ID: id, Address: string(ca)})
}

func (r *Rogue) LeftOfAt(id ID, address *ContentAddress) (ID, error) {
	isFirst := true
	for v, err := range r.WalkLeftFromAt(id, address) {
		if err != nil {
			return NoID, err
		}

		if isFirst && v.ID == id {
			isFirst = false
			continue
		}

		return v.ID, nil
	}

	ca, err := json.Marshal(address)
	if err != nil {
		return NoID, stackerr.New(err)
	}

	return NoID, stackerr.New(ErrorNoLeftSiblingAt{ID: id, Address: string(ca)})
}

func (r *Rogue) IsDeletedAt(id ID, address *ContentAddress) (bool, error) {
	var err error
	addrStartIx, addrEndIx := 0, r.TotSize

	if address != nil {
		_, addrStartIx, err = r.Rope.GetIndex(address.StartID)
		if err != nil {
			return false, err
		}

		_, addrEndIx, err = r.Rope.GetIndex(address.EndID)
		if err != nil {
			return false, err
		}
	}

	_, ix, err := r.Rope.GetIndex(id)
	if err != nil {
		return false, err
	}

	if address == nil || ix < addrStartIx || addrEndIx < ix {
		return r.Rope.IsDeleted(id)
	}

	if !address.Contains(id) {
		return true, nil
	}

	ch := r.CharHistory[id]

	if ch == nil {
		return false, nil
	}

	n, err := ch.FindLeftSibNode(*address.MaxID())
	if err != nil {
		return false, err
	}

	for {
		if n == nil {
			return false, nil
		}

		if address.Contains(n.Value.ID) {
			return n.Value.IsDel, nil
		}

		n = n.StepLeft()
	}
}
