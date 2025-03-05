package v3

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/charmbracelet/log"
)

func (r *Rogue) ToOps() ([]Op, error) {
	out := make([]Op, 0, r.OpIndex.Size()+r.FailedOps.Size())

	for _, tree := range r.OpIndex.AuthorOps {
		tree.Dft(func(op Op) error {
			out = append(out, op)
			return nil
		})
	}

	r.FailedOps.Tree.Dft(func(op Op) error {
		out = append(out, op)
		return nil
	})

	slices.SortFunc(out, func(a, b Op) int {
		aID := a.GetID()
		bID := b.GetID()

		if aID.Seq < bID.Seq {
			return -1
		} else if aID.Seq > bID.Seq {
			return 1
		}

		if aID.Author < bID.Author {
			return -1
		} else if aID.Author > bID.Author {
			return 1
		}

		return 0
	})

	return out, nil
}

func (r *Rogue) Serializable() (*SerializedRogue, error) {
	ops, err := r.ToOps()
	if err != nil {
		return nil, fmt.Errorf("Serializable(): %w", err)
	}

	return &SerializedRogue{
		Version: StringPtr("v0"),
		Ops:     ops,
	}, nil
}

func (r *Rogue) MarshalJSON() ([]byte, error) {
	sr, err := r.Serializable()
	if err != nil {
		return nil, fmt.Errorf("MarshalJSON(): %w", err)
	}
	return json.Marshal(sr)
}

func (r *Rogue) _migrateV0(op Op, nos *NOSV2) error {
	mop := MultiOp{}

	ops := []Op{op}
	if m, ok := op.(MultiOp); ok {
		ops = m.Mops
		mop.Mops = make([]Op, 0, len(ops)*2)
	}

	for _, op := range ops {
		if fop, ok := op.(FormatOp); ok {
			if _, ok := fop.Format.(FormatV3Span); ok {
				inserted, err := nos.Insert(fop)
				if err != nil {
					log.Errorf("DeserRogue() failed to merge op: %v with err: %v", op, err)
					continue
				}

				for _, n := range inserted {
					mop.Mops = append(mop.Mops, n)
				}
			} else {
				mop.Mops = append(mop.Mops, op)
			}
		} else {
			mop.Mops = append(mop.Mops, op)
		}
	}

	_, err := r.MergeOp(FlattenMop(mop))
	if err != nil {
		return err
	}

	return nil
}

func (r *Rogue) DeserRogue(sr *SerializedRogue) error {
	if r == nil {
		return fmt.Errorf("DeserRogue() nil receiver")
	}

	r.Reset()

	if sr.Version == nil {
		// Handle migration
		nos := NewNOSV2(r.Rope)

		for _, op := range sr.Ops {
			r._migrateV0(op, nos)
		}

		return nil
	}

	for _, op := range sr.Ops {
		_, err := r.MergeOp(op)
		if err != nil {
			log.Errorf("DeserRogue() failed to merge op: %v with err: %v", op, err)
		}
	}

	// TODO: Speed up deser by not setting the NOS
	// until all of the ops are loaded. Requires a different
	// MergeOp function.
	/*err := r.ResetNOS()
	if err != nil {
		return err
	}*/

	return nil
}

func (r *Rogue) UnmarshalJSON(data []byte) error {
	// printMemUsage()
	var sr SerializedRogue
	if err := json.Unmarshal(data, &sr); err != nil {
		return fmt.Errorf("rogue.UnmarshalJSON() failed to unmarshal SerializedRogue: %w", err)
	}

	err := r.DeserRogue(&sr)
	if err != nil {
		return fmt.Errorf("rogue.UnmarshalJSON() failed to deserRogue: %w", err)
	}
	// printMemUsage()

	return nil
}

// TODO: don't use MarshalJSON and UnmarshalJSON
func (r *Rogue) DeepCopy() (*Rogue, error) {
	c := &Rogue{}

	bs, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bs, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (r *Rogue) NewSnapshotOp() (SnapshotOp, error) {
	sr, err := r.Serializable()
	if err != nil {
		return SnapshotOp{}, err
	}

	return SnapshotOp{
		ID:       ID{"root", 0},
		Snapshot: sr,
	}, nil
}
