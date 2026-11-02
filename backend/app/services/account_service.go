package services

import (
	"context"

	"leochat/app/apperr"
)

// AccountService closes accounts.
type AccountService struct {
	users  UserRepository
	hasher PasswordHasher
	files  FileRemover
}

func NewAccountService(users UserRepository, hasher PasswordHasher, files FileRemover) *AccountService {
	return &AccountService{users: users, hasher: hasher, files: files}
}

// Delete permanently removes the account after checking the password again, so
// an unlocked phone with an open session is not enough. The database removes
// everything tied to the user; photo files are cleaned up afterwards.
func (s *AccountService) Delete(ctx context.Context, userID uint, password string) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !s.hasher.Check(password, user.Password) {
		return apperr.ErrWrongPassword
	}

	if err := s.users.Delete(ctx, userID); err != nil {
		return err
	}

	if len(user.Photos) > 0 {
		paths := make([]string, len(user.Photos))
		for i, p := range user.Photos {
			paths[i] = p.Path
		}
		// The account is already gone; leftover files must not fail the request.
		_ = s.files.Delete(paths...)
	}
	return nil
}
