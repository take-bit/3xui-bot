package payment

import (
	"context"
	"fmt"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {

	return &MockProvider{}
}

func (m *MockProvider) CreatePayment(ctx context.Context, amount float64, currency, description string) (string, string, error) {
	mockPaymentID := id.GenerateWithPrefix("mock_payment")
	mockURL := fmt.Sprintf("https://mock-payment.example.com/pay/%s", mockPaymentID)

	return mockURL, mockPaymentID, nil
}

func (m *MockProvider) CheckPaymentStatus(ctx context.Context, paymentID string) (string, error) {

	return string(core.PaymentStatusCompleted), nil
}
