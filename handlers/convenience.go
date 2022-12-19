package handlers

import (
	"github.com/martijnkorbee/goracoon"
)

// crypto helpers

func (h *Handlers) encrypt(text string) (string, error) {
	crypto := goracoon.Encryption{
		Key: []byte(h.Racoon.EncryptionKey),
	}

	encrypted, err := crypto.Encrypt(text)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func (h *Handlers) decrypt(encrypted string) (string, error) {
	crypto := goracoon.Encryption{
		Key: []byte(h.Racoon.EncryptionKey),
	}

	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		return "", err
	}

	return decrypted, nil
}
