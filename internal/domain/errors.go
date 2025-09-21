package domain

import "errors"

// Предопределенные ошибки домена
var (
	// Ошибки пользователей
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
	ErrUserBlocked       = errors.New("пользователь заблокирован")

	// Ошибки подписок
	ErrSubscriptionNotFound  = errors.New("подписка не найдена")
	ErrSubscriptionExpired   = errors.New("подписка истекла")
	ErrSubscriptionNotActive = errors.New("подписка не активна")
	ErrTrialAlreadyUsed      = errors.New("пробный период уже использован")

	// Ошибки планов
	ErrPlanNotFound    = errors.New("план не найден")
	ErrPlanInactive    = errors.New("план неактивен")
	ErrInvalidDuration = errors.New("неверная продолжительность подписки")
	ErrInvalidCurrency = errors.New("неверная валюта")

	// Ошибки серверов
	ErrServerNotFound    = errors.New("сервер не найден")
	ErrServerUnavailable = errors.New("сервер недоступен")
	ErrServerOverloaded  = errors.New("сервер перегружен")

	// Ошибки платежей
	ErrPaymentNotFound      = errors.New("платеж не найден")
	ErrPaymentAlreadyPaid   = errors.New("платеж уже оплачен")
	ErrPaymentExpired       = errors.New("платеж истек")
	ErrInvalidPaymentMethod = errors.New("неверный способ оплаты")

	// Ошибки промокодов
	ErrPromocodeNotFound    = errors.New("промокод не найден")
	ErrPromocodeExpired     = errors.New("промокод истек")
	ErrPromocodeUsedUp      = errors.New("промокод исчерпан")
	ErrPromocodeAlreadyUsed = errors.New("промокод уже использован")

	// Ошибки рефералов
	ErrReferralNotFound  = errors.New("реферальная связь не найдена")
	ErrSelfReferral      = errors.New("нельзя пригласить самого себя")
	ErrAlreadyReferred   = errors.New("пользователь уже приглашен")
	ErrInvalidReferrer   = errors.New("неверный реферер")

	// Ошибки уведомлений
	ErrNotificationNotFound    = errors.New("уведомление не найдено")
	ErrNotificationSent        = errors.New("уведомление уже отправлено")
	ErrInvalidNotificationType = errors.New("неверный тип уведомления")
)
