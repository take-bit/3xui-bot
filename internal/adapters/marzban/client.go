package marzban

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"3xui-bot/internal/core"
)

// MarzbanRepository представляет репозиторий для работы с Marzban API
type MarzbanRepository struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	token      string
	tokenExp   time.Time
}

// NewMarzbanRepository создает новый экземпляр Marzban repository
func NewMarzbanRepository(baseURL, username, password string) *MarzbanRepository {
	return &MarzbanRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		username: username,
		password: password,
	}
}

// Authenticate выполняет аутентификацию в Marzban API (алиас для Login)
func (m *MarzbanRepository) Authenticate(ctx context.Context) error {
	return m.Login(ctx)
}

// Login выполняет аутентификацию в Marzban API
func (m *MarzbanRepository) Login(ctx context.Context) error {
	// Marzban использует OAuth2 форму (application/x-www-form-urlencoded), а не JSON
	formData := url.Values{}
	formData.Set("username", m.username)
	formData.Set("password", m.password)

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/api/admin/token", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	m.token = loginResp.AccessToken
	m.tokenExp = time.Now().Add(23 * time.Hour) // Токен действует 24 часа, обновляем через 23

	return nil
}

// ensureToken проверяет и обновляет токен при необходимости
func (m *MarzbanRepository) ensureToken(ctx context.Context) error {
	if m.token == "" || time.Now().After(m.tokenExp) {
		return m.Login(ctx)
	}
	return nil
}

// makeRequest выполняет HTTP запрос с автоматическим обновлением токена
func (m *MarzbanRepository) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	if err := m.ensureToken(ctx); err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, m.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Если токен истек, пробуем еще раз с новым токеном
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		if err := m.Login(ctx); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Повторяем запрос с новым токеном
		if body != nil {
			jsonData, _ := json.Marshal(body)
			reqBody = bytes.NewBuffer(jsonData)
		}
		req, err = http.NewRequestWithContext(ctx, method, m.baseURL+endpoint, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create retry request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+m.token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		return m.httpClient.Do(req)
	}

	return resp, nil
}

// CreateUser создает нового пользователя в Marzban
func (m *MarzbanRepository) CreateUser(ctx context.Context, userData *core.MarzbanUserData) (*core.MarzbanUserData, error) {
	resp, err := m.makeRequest(ctx, "POST", "/api/user", userData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа для логирования
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create user with status %d: %s", resp.StatusCode, string(body))
	}

	var createdUser core.MarzbanUserData
	if err := json.Unmarshal(body, &createdUser); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w, body: %s", err, string(body))
	}

	return &createdUser, nil
}

// GetUser получает данные пользователя из Marzban
func (m *MarzbanRepository) GetUser(ctx context.Context, username string) (*core.MarzbanUserData, error) {
	resp, err := m.makeRequest(ctx, "GET", "/api/user/"+username, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user with status %d: %s", resp.StatusCode, string(body))
	}

	var userData core.MarzbanUserData
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &userData, nil
}

// UpdateUser обновляет данные пользователя в Marzban
func (m *MarzbanRepository) UpdateUser(ctx context.Context, username string, userData *core.MarzbanUserData) (*core.MarzbanUserData, error) {
	resp, err := m.makeRequest(ctx, "PUT", "/api/user/"+username, userData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update user with status %d: %s", resp.StatusCode, string(body))
	}

	var updatedUser core.MarzbanUserData
	if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &updatedUser, nil
}

// DeleteUser удаляет пользователя из Marzban
func (m *MarzbanRepository) DeleteUser(ctx context.Context, username string) error {
	resp, err := m.makeRequest(ctx, "DELETE", "/api/user/"+username, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil // Пользователь уже удален или не существует
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetUsers получает список всех пользователей (с пагинацией)
func (m *MarzbanRepository) GetUsers(ctx context.Context, offset, limit int) ([]*core.MarzbanUserData, error) {
	endpoint := fmt.Sprintf("/api/users?offset=%d&limit=%d", offset, limit)
	resp, err := m.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get users with status %d: %s", resp.StatusCode, string(body))
	}

	var users []*core.MarzbanUserData
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode users response: %w", err)
	}

	return users, nil
}

// GetSystemStats получает статистику системы
func (m *MarzbanRepository) GetSystemStats(ctx context.Context) (map[string]interface{}, error) {
	resp, err := m.makeRequest(ctx, "GET", "/api/system", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get system stats with status %d: %s", resp.StatusCode, string(body))
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode system stats response: %w", err)
	}

	return stats, nil
}

// GetUserUsage получает статистику использования пользователя
func (m *MarzbanRepository) GetUserUsage(ctx context.Context, username string) (map[string]interface{}, error) {
	resp, err := m.makeRequest(ctx, "GET", "/api/user/"+username+"/usage", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user usage with status %d: %s", resp.StatusCode, string(body))
	}

	var usage map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("failed to decode user usage response: %w", err)
	}

	return usage, nil
}

// GetInbounds получает список доступных inbounds
func (m *MarzbanRepository) GetInbounds(ctx context.Context) ([]map[string]interface{}, error) {
	resp, err := m.makeRequest(ctx, "GET", "/api/inbounds", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get inbounds with status %d: %s", resp.StatusCode, string(body))
	}

	// API может вернуть либо массив, либо объект с массивом внутри
	// Сначала пробуем как массив
	var inboundsArray []map[string]interface{}
	if err := json.Unmarshal(body, &inboundsArray); err == nil {
		return inboundsArray, nil
	}

	// Если не массив, пробуем как объект с полем (например, {"inbounds": [...]})
	var inboundsObject map[string]interface{}
	if err := json.Unmarshal(body, &inboundsObject); err != nil {
		return nil, fmt.Errorf("failed to decode inbounds response: %w, body: %s", err, string(body))
	}

	// Ищем массив в объекте
	for _, value := range inboundsObject {
		if arr, ok := value.([]interface{}); ok {
			result := make([]map[string]interface{}, 0, len(arr))
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					result = append(result, m)
				}
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("unexpected inbounds response format: %s", string(body))
}

// GetStats получает статистику системы
func (m *MarzbanRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	if err := m.ensureTokenValid(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+"/api/system", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stats request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get stats with status %d: %s", resp.StatusCode, string(body))
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	return stats, nil
}

// ResetUserTraffic сбрасывает трафик пользователя
func (m *MarzbanRepository) ResetUserTraffic(ctx context.Context, username string) error {
	if err := m.ensureTokenValid(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.baseURL+"/api/user/"+username+"/reset", nil)
	if err != nil {
		return fmt.Errorf("failed to create reset request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to reset traffic with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ensureTokenValid проверяет валидность токена и переаутентифицируется при необходимости
func (m *MarzbanRepository) ensureTokenValid(ctx context.Context) error {
	if m.token == "" || time.Now().After(m.tokenExp) {
		return m.Login(ctx)
	}
	return nil
}
