package usecase

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized access")
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
)

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserTrialAlreadyUsed = errors.New("user trial already used")
)

var (
	ErrSubscriptionNotActive    = errors.New("subscription not active")
	ErrSubscriptionExpired      = errors.New("subscription expired")
	ErrPlanNotActive            = errors.New("plan not active")
	ErrSubscriptionLimitReached = errors.New("subscription limit reached")
)

var (
	ErrPaymentAlreadyPaid = errors.New("payment already paid")
	ErrPaymentCancelled   = errors.New("payment cancelled")
	ErrPaymentFailed      = errors.New("payment failed")
	ErrInvalidAmount      = errors.New("invalid amount")
)

var (
	ErrVPNConfigNotActive    = errors.New("VPN config not active")
	ErrVPNConfigLimitReached = errors.New("VPN config limit reached")
	ErrInvalidVPNType        = errors.New("invalid VPN type")
)

var (
	ErrSelfReferral          = errors.New("cannot refer yourself")
	ErrReferralAlreadyExists = errors.New("referral already exists")
)
