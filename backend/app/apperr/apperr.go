// Package apperr defines the error vocabulary shared by services and the HTTP layer.
// Services return these errors; the HTTP layer maps each Kind to a status code.
package apperr

import "errors"

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }

func newError(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

// As unwraps err into an *Error when it is one.
func As(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

var (
	ErrInvalidCredentials = newError(KindUnauthorized, "invalid_credentials", "Email atau password salah.")
	ErrEmailTaken         = newError(KindConflict, "email_taken", "Email sudah terdaftar.")
	ErrUnderage           = newError(KindInvalid, "underage", "Kamu harus berusia minimal 18 tahun.")

	ErrUserNotFound     = newError(KindNotFound, "user_not_found", "Pengguna tidak ditemukan.")
	ErrLocationRequired = newError(KindInvalid, "location_required", "Aktifkan lokasi untuk mulai mencari.")

	ErrPhotoLimitReached = newError(KindInvalid, "photo_limit_reached", "Maksimal 6 foto per profil.")
	ErrPhotoNotFound     = newError(KindNotFound, "photo_not_found", "Foto tidak ditemukan.")

	ErrCannotSwipeSelf = newError(KindInvalid, "cannot_swipe_self", "Tidak bisa swipe profil sendiri.")
	ErrAlreadySwiped   = newError(KindConflict, "already_swiped", "Profil ini sudah kamu swipe.")

	ErrMatchNotFound = newError(KindNotFound, "match_not_found", "Match tidak ditemukan.")
	ErrMessageEmpty  = newError(KindInvalid, "message_empty", "Pesan tidak boleh kosong.")
)
