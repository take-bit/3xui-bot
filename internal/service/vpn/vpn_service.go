package vpn

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"
)

// VPNService представляет сервис для работы с VPN подключениями
type VPNService struct {
	userRepo         domain.UserRepository
	subscriptionRepo domain.SubscriptionRepository
	serverRepo       domain.ServerRepository
	xuiClient        domain.XUIClient
}

// NewVPNService создает новый VPN сервис
func NewVPNService(
	userRepo domain.UserRepository,
	subscriptionRepo domain.SubscriptionRepository,
	serverRepo domain.ServerRepository,
	xuiClient domain.XUIClient,
) *VPNService {
	return &VPNService{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		serverRepo:       serverRepo,
		xuiClient:        xuiClient,
	}
}

// CreateConnection создает VPN подключение для пользователя
func (s *VPNService) CreateConnection(ctx context.Context, userID int64) (*domain.VPNConnection, error) {
	return s.CreateConnectionWithRegion(ctx, userID, "")
}

// CreateConnectionWithRegion создает VPN подключение для пользователя с указанием региона
func (s *VPNService) CreateConnectionWithRegion(ctx context.Context, userID int64, region string) (*domain.VPNConnection, error) {
	// Проверяем, что пользователь существует
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем активную подписку
	subscription, err := s.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Выбираем доступный сервер (простая логика)
	server, err := s.selectAvailableServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select server: %w", err)
	}

	// Аутентифицируемся в 3X-UI
	if err := s.xuiClient.Login(ctx); err != nil {
		return nil, fmt.Errorf("failed to login to 3X-UI: %w", err)
	}

	// Получаем список inbound'ов
	inbounds, err := s.xuiClient.GetInbounds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	if len(inbounds) == 0 {
		return nil, fmt.Errorf("no inbounds available")
	}

	// Используем первый доступный inbound
	inbound := inbounds[0]

	// Генерируем UUID для клиента
	uuid, err := s.generateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Вычисляем лимит трафика (в GB)
	totalGB := s.calculateTrafficLimit(subscription)

	// Вычисляем время истечения
	expiryTime := subscription.EndDate.Unix()

	// Добавляем клиента в 3X-UI
	if err := s.xuiClient.AddClient(ctx, inbound.ID, userID, uuid, totalGB, expiryTime); err != nil {
		return nil, fmt.Errorf("failed to add client to 3X-UI: %w", err)
	}

	// Обновляем счетчик клиентов на сервере
	if err := s.serverRepo.IncrementClients(ctx, server.ID); err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("failed to increment server clients: %v\n", err)
	}

	// Создаем URL для конфигурации
	configURL := s.generateConfigURL(server, inbound, uuid)

	// Создаем объект подключения
	connection := &domain.VPNConnection{
		UserID:       userID,
		ServerID:     server.ID,
		XUIInboundID: inbound.ID,
		XUIClientID:  "", // Будет заполнено при получении списка клиентов
		UUID:         uuid,
		Email:        fmt.Sprintf("user_%d@vpn.local", userID),
		ConfigURL:    configURL,
		CreatedAt:    time.Now(),
		ExpiresAt:    subscription.EndDate,
	}

	return connection, nil
}

// DeleteConnection удаляет VPN подключение пользователя
func (s *VPNService) DeleteConnection(ctx context.Context, userID int64) error {
	// Получаем информацию о пользователе
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Аутентифицируемся в 3X-UI
	if err := s.xuiClient.Login(ctx); err != nil {
		return fmt.Errorf("failed to login to 3X-UI: %w", err)
	}

	// Получаем список inbound'ов
	inbounds, err := s.xuiClient.GetInbounds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get inbounds: %w", err)
	}

	// Ищем клиента по сгенерированному email
	expectedEmail := fmt.Sprintf("user_%d@vpn.local", userID)

	for _, inbound := range inbounds {
		clients, err := s.xuiClient.GetClients(ctx, inbound.ID)
		if err != nil {
			continue // Пропускаем ошибки и продолжаем поиск
		}

		for _, client := range clients {
			if client.Email == expectedEmail {
				// Удаляем клиента
				if err := s.xuiClient.DeleteClient(ctx, inbound.ID, client.ID); err != nil {
					return fmt.Errorf("failed to delete client: %w", err)
				}

				// Уменьшаем счетчик клиентов на сервере
				// TODO: Нужно получить serverID из базы данных
				// if err := s.serverRepo.DecrementClients(ctx, serverID); err != nil {
				//     fmt.Printf("failed to decrement server clients: %v\n", err)
				// }

				return nil
			}
		}
	}

	return fmt.Errorf("client not found")
}

// GetConnectionInfo получает информацию о VPN подключении пользователя
func (s *VPNService) GetConnectionInfo(ctx context.Context, userID int64) (*domain.VPNConnection, error) {
	// Получаем информацию о пользователе
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Аутентифицируемся в 3X-UI
	if err := s.xuiClient.Login(ctx); err != nil {
		return nil, fmt.Errorf("failed to login to 3X-UI: %w", err)
	}

	// Получаем список inbound'ов
	inbounds, err := s.xuiClient.GetInbounds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	// Ищем клиента по сгенерированному email
	expectedEmail := fmt.Sprintf("user_%d@vpn.local", userID)

	for _, inbound := range inbounds {
		clients, err := s.xuiClient.GetClients(ctx, inbound.ID)
		if err != nil {
			continue
		}

		for _, client := range clients {
			if client.Email == expectedEmail {
				// Получаем информацию о сервере
				server, err := s.serverRepo.GetByID(ctx, 1) // TODO: Нужно сохранять serverID
				if err != nil {
					server = &domain.Server{ID: 1, Name: "Unknown"} // Fallback
				}

				// Создаем URL для конфигурации
				configURL := s.generateConfigURL(server, inbound, client.UUID)

				connection := &domain.VPNConnection{
					UserID:       userID,
					ServerID:     server.ID,
					XUIInboundID: inbound.ID,
					XUIClientID:  client.ID,
					UUID:         client.UUID,
					Email:        client.Email,
					ConfigURL:    configURL,
					CreatedAt:    time.Now(), // TODO: Нужно сохранять в БД
					ExpiresAt:    time.Unix(client.ExpiryTime, 0),
				}

				return connection, nil
			}
		}
	}

	return nil, fmt.Errorf("connection not found")
}

// UpdateConnectionExpiry обновляет время истечения подключения
func (s *VPNService) UpdateConnectionExpiry(ctx context.Context, userID int64, newExpiryTime time.Time) error {
	// Получаем активную подписку
	subscription, err := s.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Аутентифицируемся в 3X-UI
	if err := s.xuiClient.Login(ctx); err != nil {
		return fmt.Errorf("failed to login to 3X-UI: %w", err)
	}

	// Получаем список inbound'ов
	inbounds, err := s.xuiClient.GetInbounds(ctx)
	if err != nil {
		return fmt.Errorf("failed to get inbounds: %w", err)
	}

	// Ищем клиента по сгенерированному email
	expectedEmail := fmt.Sprintf("user_%d@vpn.local", userID)

	for _, inbound := range inbounds {
		clients, err := s.xuiClient.GetClients(ctx, inbound.ID)
		if err != nil {
			continue
		}

		for _, client := range clients {
			if client.Email == expectedEmail {
				// Обновляем время истечения
				if err := s.xuiClient.UpdateClient(ctx, inbound.ID, client.ID, userID, client.TotalGB, newExpiryTime.Unix()); err != nil {
					return fmt.Errorf("failed to update client expiry: %w", err)
				}

				// Обновляем подписку в базе данных
				subscription.EndDate = newExpiryTime
				subscription.UpdatedAt = time.Now()

				if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
					return fmt.Errorf("failed to update subscription: %w", err)
				}

				return nil
			}
		}
	}

	return fmt.Errorf("client not found")
}

// Вспомогательные методы

// selectAvailableServer выбирает доступный сервер
func (s *VPNService) selectAvailableServer(ctx context.Context) (*domain.Server, error) {
	// Получаем доступные серверы
	servers, err := s.serverRepo.GetAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available servers: %w", err)
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no available servers")
	}

	// Выбираем сервер с наименьшей нагрузкой
	var selectedServer *domain.Server
	minLoad := float64(100)

	for _, server := range servers {
		load := server.GetLoadPercentage()
		if load < minLoad {
			minLoad = load
			selectedServer = server
		}
	}

	if selectedServer == nil {
		// Если не удалось выбрать по нагрузке, берем первый доступный
		selectedServer = servers[0]
	}

	return selectedServer, nil
}

// generateUUID генерирует UUID для клиента
func (s *VPNService) generateUUID() (string, error) {
	// Простая генерация UUID v4
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// Устанавливаем версию (4) и variant bits
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// calculateTrafficLimit вычисляет лимит трафика на основе подписки
func (s *VPNService) calculateTrafficLimit(subscription *domain.Subscription) int64 {
	// Получаем план подписки
	// TODO: Нужно получать план из базы данных
	// Пока используем базовое значение
	days := int(subscription.EndDate.Sub(subscription.StartDate).Hours() / 24)

	// Примерная логика: 10GB на день
	return int64(days * 10)
}

// generateConfigURL генерирует URL для скачивания конфигурации
func (s *VPNService) generateConfigURL(server *domain.Server, inbound domain.XUIServerInfo, uuid string) string {
	// Генерируем URL для подписки
	// Формат: https://server:port/path/uuid
	return fmt.Sprintf("https://%s:%d/user/%s", server.Host, inbound.Port, uuid)
}
