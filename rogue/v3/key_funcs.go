package v3

import (
	"golang.org/x/exp/constraints"
)

func lessThan[T constraints.Ordered](a, b T) bool {
	return a < b
}

func opGetID(op Op) (ID, error) {
	return op.GetID(), nil
}

func lessThanRevID(a, b ID) bool {
	if a.Seq < b.Seq {
		return true
	} else if a.Seq > b.Seq {
		return false
	} else {
		return a.Author < b.Author
	}
}
