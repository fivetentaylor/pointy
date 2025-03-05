package v3

import "fmt"

type ErrorEmptySpan struct{}

func (e ErrorEmptySpan) Error() string {
	return "empty span"
}

type ErrorEmptyFormat struct{}

func (e ErrorEmptyFormat) Error() string {
	return "empty format"
}

type ErrorStopIteration struct{}

func (e ErrorStopIteration) Error() string {
	return "stop iteration"
}

type ErrorParentNotFound struct {
	message string
}

func (e ErrorParentNotFound) Error() string {
	return e.message
}

type ErrorInvalidOffset struct {
	ID     *ID
	Offset int
}

func (e ErrorInvalidOffset) Error() string {
	return fmt.Sprintf("invalid offset: %d at id: %+v", e.Offset, e.ID)
}

type ErrorNoLeftVisSibling struct {
	ID ID
}

func (e ErrorNoLeftVisSibling) Error() string {
	return fmt.Sprintf("no left visible sibling for id: %+v", e.ID)
}

type ErrorNoRightVisSibling struct {
	ID ID
}

func (e ErrorNoRightVisSibling) Error() string {
	return fmt.Sprintf("no right visible sibling for id: %+v", e.ID)
}

type ErrorNoLeftTotSibling struct {
	ID ID
}

func (e ErrorNoLeftTotSibling) Error() string {
	return fmt.Sprintf("no left total sibling for id: %+v", e.ID)
}

type ErrorNoRightTotSibling struct {
	ID ID
}

func (e ErrorNoRightTotSibling) Error() string {
	return fmt.Sprintf("no right total sibling for id: %+v", e.ID)
}

type ErrorNoRightSiblingAt struct {
	ID      ID
	Address string
}

func (e ErrorNoRightSiblingAt) Error() string {
	return fmt.Sprintf("no right sibling at address %q for id: %+v", e.Address, e.ID)
}

type ErrorNoLeftSiblingAt struct {
	ID      ID
	Address string
}

func (e ErrorNoLeftSiblingAt) Error() string {
	return fmt.Sprintf("no left sibling at address %q for id: %+v", e.Address, e.ID)
}

type ErrorNotCodeblock struct {
	ID      ID
	Address string
}

func (e ErrorNotCodeblock) Error() string {
	return fmt.Sprintf("not a codeblock at address %q for id: %v", e.Address, e.ID)
}
