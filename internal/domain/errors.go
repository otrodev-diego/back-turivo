package domain

import "errors"

var (
	// Common errors
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternalError = errors.New("internal server error")

	// User specific errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBlocked        = errors.New("user is blocked")

	// Company specific errors
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyAlreadyExists = errors.New("company already exists")
	ErrCompanySuspended     = errors.New("company is suspended")

	// Hotel specific errors
	ErrHotelNotFound = errors.New("hotel not found")

	// Driver specific errors
	ErrDriverNotFound      = errors.New("driver not found")
	ErrDriverAlreadyExists = errors.New("driver already exists")
	ErrDriverNotAvailable  = errors.New("driver not available")
	ErrDriverInactive      = errors.New("driver is inactive")

	// Request specific errors
	ErrRequestNotFound         = errors.New("request not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrRequestAlreadyAssigned  = errors.New("request already assigned")

	// Reservation specific errors
	ErrReservationNotFound = errors.New("reservation not found")
	ErrReservationPastDate = errors.New("reservation date is in the past")

	// Payment specific errors
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrPaymentFailed      = errors.New("payment failed")
	ErrPaymentAlreadyPaid = errors.New("payment already processed")

	// Feedback specific errors
	ErrFeedbackNotFound      = errors.New("feedback not found")
	ErrFeedbackAlreadyExists = errors.New("feedback already exists for this trip")

	// Auth specific errors
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)
