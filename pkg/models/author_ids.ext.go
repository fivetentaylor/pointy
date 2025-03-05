package models

import "strconv"

func (a *AuthorID) AuthorIDString() string {
	return strconv.FormatInt(int64(a.AuthorID), 16)
}
