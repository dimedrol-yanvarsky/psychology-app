package errors

import "errors"

// User domain errors
var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrUserNotFound   = errors.New("user not found")
	ErrWrongPassword  = errors.New("wrong password")
	ErrUserDeleted    = errors.New("user deleted")
	ErrUserBlocked    = errors.New("user blocked")
	ErrUserExists     = errors.New("user exists")
	ErrInvalidEmail   = errors.New("invalid email")
	ErrPasswordsMatch = errors.New("passwords mismatch")
)

// Common errors
var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrDatabase  = errors.New("database error")
	ErrInvalidID = errors.New("invalid id")
)

// Test errors
var (
	ErrNoQuestions = errors.New("no questions")
)

// Review errors
var (
	ErrReviewExists = errors.New("review already exists")
)

// Recommendation errors
var (
	ErrResequenceFailed = errors.New("resequence failed")
	ErrListFetchFailed  = errors.New("list fetch failed")
)
