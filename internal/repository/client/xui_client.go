package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"3xui-bot/internal/domain"
)

// XUIClient представляет клиент для работы с 3X-UI API
type XUIClient struct {
	baseURL    string
	username   string
	password   string
	token      string
	httpClient *http.Client
}

// XUIConfig представляет конфигурацию для XUI клиента
type XUIConfig struct {
	BaseURL  string
	Username string
	Password string
	Token    string
	Timeout  time.Duration
}

// NewXUIClient создает новый клиент для работы с 3X-UI API
func NewXUIClient(config XUIConfig) domain.XUIClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &XUIClient{
		baseURL:    config.BaseURL,
		username:   config.Username,
		password:   config.Password,
		token:      config.Token,
		httpClient: &http.Client{Timeout: config.Timeout},
	}
}

// XUIResponse представляет базовый ответ от 3X-UI API
type XUIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"msg"`
	Data    interface{} `json:"obj"`
}

// Login выполняет аутентификацию в 3X-UI панели
func (c *XUIClient) Login(ctx context.Context) error {
	if c.token != "" {
		// Если токен уже есть, проверяем его валидность
		return c.validateToken(ctx)
	}

	loginData := map[string]string{
		"username": c.username,
		"password": c.password,
	}

	response, err := c.makeRequest(ctx, "POST", "/login", loginData)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	var loginResp XUIResponse
	if err := json.Unmarshal(response, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	if !loginResp.Success {
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

	// Сохраняем токен из ответа
	if tokenData, ok := loginResp.Data.(map[string]interface{}); ok {
		if token, ok := tokenData["token"].(string); ok {
			c.token = token
		}
	}

	return nil
}

// validateToken проверяет валидность токена
func (c *XUIClient) validateToken(ctx context.Context) error {
	response, err := c.makeRequest(ctx, "GET", "/xui/inbound/list", nil)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	var resp XUIResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return fmt.Errorf("failed to parse validation response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("token validation failed: %s", resp.Message)
	}

	return nil
}

// GetInbounds получает список inbound'ов
func (c *XUIClient) GetInbounds(ctx context.Context) ([]domain.XUIServerInfo, error) {
	response, err := c.makeRequest(ctx, "GET", "/xui/inbound/list", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	var resp XUIResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse inbounds response: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get inbounds: %s", resp.Message)
	}

	// Парсим данные в список серверов
	var inbounds []domain.XUIServerInfo
	if data, ok := resp.Data.([]interface{}); ok {
		for _, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				inbound := domain.XUIServerInfo{
					ID:       int(itemMap["id"].(float64)),
					Remark:   itemMap["remark"].(string),
					Address:  itemMap["address"].(string),
					Port:     int(itemMap["port"].(float64)),
					Protocol: itemMap["protocol"].(string),
					Enable:   itemMap["enable"].(bool),
				}
				inbounds = append(inbounds, inbound)
			}
		}
	}

	return inbounds, nil
}

// GetClients получает список клиентов для указанного inbound
func (c *XUIClient) GetClients(ctx context.Context, inboundID int) ([]domain.XUIClientInfo, error) {
	response, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/xui/inbound/listClient/%d", inboundID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	var resp XUIResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse clients response: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("failed to get clients: %s", resp.Message)
	}

	// Парсим данные в список клиентов
	var clients []domain.XUIClientInfo
	if data, ok := resp.Data.([]interface{}); ok {
		for _, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				client := domain.XUIClientInfo{
					ID:         itemMap["id"].(string),
					Email:      itemMap["email"].(string),
					UUID:       itemMap["uuid"].(string),
					AlterID:    int(itemMap["alterId"].(float64)),
					Level:      int(itemMap["level"].(float64)),
					Enable:     itemMap["enable"].(bool),
					TotalGB:    int64(itemMap["totalGB"].(float64)),
					UsedGB:     int64(itemMap["usedGB"].(float64)),
					ExpiryTime: int64(itemMap["expiryTime"].(float64)),
				}
				clients = append(clients, client)
			}
		}
	}

	return clients, nil
}

// AddClient добавляет нового клиента к inbound
func (c *XUIClient) AddClient(ctx context.Context, inboundID int, userID int64, uuid string, totalGB int64, expiryTime int64) error {
	clientData := map[string]interface{}{
		"id":         "",
		"email":      fmt.Sprintf("user_%d@vpn.local", userID), // Генерируем email из User ID
		"uuid":       uuid,
		"alterId":    0,
		"level":      0,
		"enable":     true,
		"totalGB":    totalGB,
		"usedGB":     0,
		"expiryTime": expiryTime,
	}

	response, err := c.makeRequest(ctx, "POST", fmt.Sprintf("/xui/inbound/addClient/%d", inboundID), clientData)
	if err != nil {
		return fmt.Errorf("failed to add client: %w", err)
	}

	var resp XUIResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return fmt.Errorf("failed to parse add client response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to add client: %s", resp.Message)
	}

	return nil
}

// UpdateClient обновляет данные клиента
func (c *XUIClient) UpdateClient(ctx context.Context, inboundID int, clientID string, userID int64, totalGB int64, expiryTime int64) error {
	clientData := map[string]interface{}{
		"id":         clientID,
		"email":      fmt.Sprintf("user_%d@vpn.local", userID), // Генерируем email из User ID
		"totalGB":    totalGB,
		"expiryTime": expiryTime,
	}

	response, err := c.makeRequest(ctx, "POST", fmt.Sprintf("/xui/inbound/updateClient/%d", inboundID), clientData)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	var resp XUIResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return fmt.Errorf("failed to parse update client response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to update client: %s", resp.Message)
	}

	return nil
}

// DeleteClient удаляет клиента
func (c *XUIClient) DeleteClient(ctx context.Context, inboundID int, clientID string) error {
	response, err := c.makeRequest(ctx, "POST", fmt.Sprintf("/xui/inbound/delClient/%d", inboundID), map[string]string{"id": clientID})
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	var resp XUIResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return fmt.Errorf("failed to parse delete client response: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to delete client: %s", resp.Message)
	}

	return nil
}

// makeRequest выполняет HTTP запрос к 3X-UI API
func (c *XUIClient) makeRequest(ctx context.Context, method, endpoint string, data interface{}) ([]byte, error) {
	url := c.baseURL + endpoint

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
