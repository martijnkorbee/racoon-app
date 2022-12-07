package handlers

import (
	"context"
	"net/http"

	_ "github.com/martijnkorbee/goracoon"
)

// render helper function to render pages
func (h *Handlers) render(
	w http.ResponseWriter,
	r *http.Request,
	tmpl string,
	variables, data interface{},
) error {
	return h.App.Render.Page(w, r, tmpl, variables, data)
}

// logError helper function to log errors
func (h *Handlers) logError(v ...any) {
	h.App.ErrorLog.Println(v...)
}

// crypto helpers

func (h *Handlers) encrypt(text string) (string, error) {
	crypto := goracoon.Encryption{
		Key: []byte(h.App.EncryptionKey),
	}

	encrypted, err := crypto.Encrypt(text)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func (h *Handlers) decrypt(encrypted string) (string, error) {
	crypto := goracoon.Encryption{
		Key: []byte(h.App.EncryptionKey),
	}

	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		return "", err
	}

	return decrypted, nil
}

// session helpers

func (h *Handlers) sessionPut(ctx context.Context, key string, value interface{}) {
	h.App.SessionManager.Put(ctx, key, value)
}

func (h *Handlers) sessionHas(ctx context.Context, key string) bool {
	return h.App.SessionManager.Exists(ctx, key)
}

func (h *Handlers) sessionGet(ctx context.Context, key string) interface{} {
	return h.App.SessionManager.Get(ctx, key)
}

func (h *Handlers) sessionRemove(ctx context.Context, key string) {
	h.App.SessionManager.Remove(ctx, key)
}

func (h *Handlers) sessionRenew(ctx context.Context) error {
	return h.App.SessionManager.RenewToken(ctx)
}

func (h *Handlers) sessionDestroy(ctx context.Context) {
	h.App.SessionManager.Destroy(ctx)
}

func (h *Handlers) randomString(n int) string {
	return h.App.RandomStringGenerator(n)
}
