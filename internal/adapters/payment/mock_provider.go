package payment

import (
	"context"
	"fmt"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
)

// MockProvider моковая реализация платежного провайдера
type MockProvider struct{}

// NewMockProvider создает новый mock провайдер
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// CreatePayment создает мок платеж
func (m *MockProvider) CreatePayment(ctx context.Context, amount float64, currency, description string) (string, string, error) {
	mockPaymentID := id.GenerateWithPrefix("mock_payment")
	mockURL := fmt.Sprintf("https://mock-payment.example.com/pay/%s", mockPaymentID)

	// В реальности здесь будет вызов API ЮKassa/Stripe
	return mockURL, mockPaymentID, nil
}

// CheckPaymentStatus проверяет статус платежа (мок)
func (m *MockProvider) CheckPaymentStatus(ctx context.Context, paymentID string) (string, error) {
	// В реальности здесь будет вызов API
	return string(core.PaymentStatusCompleted), nil
}
