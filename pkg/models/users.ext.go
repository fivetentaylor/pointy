package models

import (
	"fmt"
	"strings"

	"github.com/teamreviso/code/pkg/constants"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func (user *User) ShortName() string {
	return strings.Split(user.Name, " ")[0]
}

func (user *User) Initials() string {
	parts := strings.Split(user.Name, " ")
	initials := ""
	for _, part := range parts {
		initials += part[0:1]
	}
	return initials
}

func (user *User) HighlightColor() string {
	uid, err := uuid.Parse(user.ID)
	if err != nil {
		log.Error("error parsing uuid", "err", err)
		return constants.DefaultHighlightColor
	}
	var intVal uint64
	for i := 0; i < 8; i++ {
		intVal = (intVal << 8) | uint64(uid[i])
	}

	return constants.HighlightColors[intVal%uint64(len(constants.HighlightColors))]
}

func (user *User) IsReviso() bool {
	return user.ID == constants.RevisoUserID
}

func (user *User) AvatarS3Key() string {
	return fmt.Sprintf(constants.UserAvatarKeyFormat, user.ID)
}
