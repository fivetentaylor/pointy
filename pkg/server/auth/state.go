package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/teamreviso/code/pkg/server/auth/types"
)

const initialStateLength = 16

func (a *Manager) generateStateString(value types.State) (string, error) {
	bytes := make([]byte, initialStateLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(bytes)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	stateWithValues := state + ":" + timestamp + ":" + value.Next + ":" + value.Sidebar

	return a.encrypt(stateWithValues)
}

func (a *Manager) verifyStateString(encryptedState string) (*types.State, bool, error) {
	// Decrypt the state parameter
	stateWithValues, err := a.decrypt(encryptedState)
	if err != nil {
		return nil, false, fmt.Errorf("invalid state parameter: %w", err)
	}

	// Split the state and timestamp
	parts := strings.Split(stateWithValues, ":")
	if len(parts) != 4 {
		return nil, false, fmt.Errorf("invalid state parameter format: %w", err)
	}

	timestampStr, next, sidebar := parts[1], parts[2], parts[3]

	// Verify the timestamp
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, false, fmt.Errorf("invalid timestamp: %w", err)
	}
	if time.Now().Unix()-timestamp > 300 { // 300 seconds (5 minutes)
		return nil, false, fmt.Errorf("expired state parameter")
	}

	return &types.State{
		Next:    next,
		Sidebar: sidebar,
	}, true, nil
}

func (a *Manager) encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(a.secret)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (a *Manager) decrypt(encryptedText string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(a.secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func stateToNextURL(state *types.State) string {
	base := ""

	if state == nil {
		return base
	}

	if state.Next != "" && state.Next[0] == '/' {
		base += state.Next
		if state.Sidebar != "" {
			base += "?sb=" + state.Sidebar
		}
	}

	return base
}
