package handlers

import (
	"log/slog"
	"context"
	"fmt"
	"strings"

	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// VPNHandler обработчик VPN подключений
type VPNHandler struct {
	bot   *tgbotapi.BotAPI
	vpnUC *usecase.VPNUseCase
}

// NewVPNHandler создает новый обработчик VPN
func NewVPNHandler(
	bot *tgbotapi.BotAPI,
	vpnUC *usecase.VPNUseCase,
) *VPNHandler {
	return &VPNHandler{
		bot:   bot,
		vpnUC: vpnUC,
	}
}

// HandleShowVPNs показывает список VPN подключений пользователя
func (h *VPNHandler) HandleShowVPNs(ctx context.Context, userID int64, chatID int64) error {
	slog.Info("Showing VPNs for user %d", userID)

	// Получаем VPN подключения с данными из Marzban через UseCase
	vpns, err := h.vpnUC.GetUserVPNWithStats(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get VPNs: %w", err)
	}

	if len(vpns) == 0 {
		msg := tgbotapi.NewMessage(chatID,
			"📭 У вас пока нет активных VPN подключений.\n\n"+
				"Приобретите подписку, чтобы получить доступ к VPN.",
		)
		h.bot.Send(msg)
		return nil
	}

	// Формируем сообщение со списком VPN
	var message strings.Builder
	message.WriteString("🔐 *Ваши VPN подключения:*\n\n")

	for i, vpn := range vpns {
		statusEmoji := "✅"
		if !vpn.IsActive || vpn.Status != "active" {
			statusEmoji = "❌"
		}

		message.WriteString(fmt.Sprintf(
			"%d. %s *%s*\n"+
				"   Статус: %s %s\n"+
				"   Username: `%s`\n",
			i+1,
			statusEmoji,
			vpn.Name,
			statusEmoji,
			vpn.Status,
			vpn.MarzbanUsername,
		))

		// Добавляем статистику если есть
		if vpn.DataLimitBytes != nil && *vpn.DataLimitBytes > 0 {
			usedGB := float64(*vpn.DataUsedBytes) / (1024 * 1024 * 1024)
			limitGB := float64(*vpn.DataLimitBytes) / (1024 * 1024 * 1024)
			message.WriteString(fmt.Sprintf("   Трафик: %.2f / %.2f GB\n", usedGB, limitGB))
		}

		if vpn.ExpireAt != nil {
			message.WriteString(fmt.Sprintf("   Истекает: %s\n", vpn.ExpireAt.Format("02.01.2006 15:04")))
		}

		message.WriteString("\n")
	}

	message.WriteString("Выберите VPN для получения конфигурации:")

	// Создаем клавиатуру с VPN подключениями
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, vpn := range vpns {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("📥 %s", vpn.Name),
				fmt.Sprintf("vpn_config_%s", vpn.ID),
			),
		)
		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, message.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// HandleGetVPNConfig отправляет конфигурацию VPN
func (h *VPNHandler) HandleGetVPNConfig(ctx context.Context, userID int64, chatID int64, vpnID string) error {
	slog.Info("Getting VPN config %s for user %d", vpnID, userID)

	// Получаем VPN подключение через UseCase
	vpn, err := h.vpnUC.GetVPNConnectionWithStats(ctx, vpnID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ VPN подключение не найдено.")
		h.bot.Send(msg)
		return fmt.Errorf("failed to get VPN: %w", err)
	}

	// Проверяем что VPN принадлежит пользователю
	if vpn.TelegramUserID != userID {
		msg := tgbotapi.NewMessage(chatID, "❌ Доступ запрещен.")
		h.bot.Send(msg)
		return fmt.Errorf("unauthorized access to VPN")
	}

	// TODO: Генерировать реальные конфигурационные файлы из данных Marzban
	// Пока отправляем информацию о подключении

	configText := fmt.Sprintf(
		"🔐 *Конфигурация VPN: %s*\n\n"+
			"Username: `%s`\n"+
			"Статус: %s\n\n"+
			"📝 *Инструкция по подключению:*\n"+
			"1. Скачайте приложение VPN клиента\n"+
			"2. Импортируйте конфигурацию\n"+
			"3. Подключитесь к серверу\n\n"+
			"⚠️ Не делитесь конфигурацией с другими!",
		vpn.Name,
		vpn.MarzbanUsername,
		vpn.Status,
	)

	// Создаем клавиатуру с действиями
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", fmt.Sprintf("vpn_stats_%s", vpn.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", fmt.Sprintf("vpn_refresh_%s", vpn.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "vpn_list"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, configText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// HandleVPNStats показывает статистику VPN
func (h *VPNHandler) HandleVPNStats(ctx context.Context, userID int64, chatID int64, messageID int, vpnID string) error {
	slog.Info("Showing stats for VPN %s", vpnID)

	// Получаем VPN с актуальными данными через UseCase
	vpn, err := h.vpnUC.GetVPNConnectionWithStats(ctx, vpnID)
	if err != nil {
		return fmt.Errorf("failed to get VPN: %w", err)
	}

	usedGB := 0.0
	limitGB := 0.0
	if vpn.DataUsedBytes != nil {
		usedGB = float64(*vpn.DataUsedBytes) / (1024 * 1024 * 1024)
	}
	if vpn.DataLimitBytes != nil {
		limitGB = float64(*vpn.DataLimitBytes) / (1024 * 1024 * 1024)
	}
	usagePercent := 0.0
	if limitGB > 0 {
		usagePercent = (usedGB / limitGB) * 100
	}

	statsText := fmt.Sprintf(
		"📊 *Статистика VPN: %s*\n\n"+
			"📈 Использовано: %.2f GB / %.2f GB (%.1f%%)\n"+
			"📅 Истекает: %s\n"+
			"✅ Статус: %s\n\n"+
			"Обновлено: %s",
		vpn.Name,
		usedGB,
		limitGB,
		usagePercent,
		vpn.ExpireAt.Format("02.01.2006 15:04"),
		vpn.Status,
		vpn.UpdatedAt.Format("02.01.2006 15:04"),
	)

	// Обновляем сообщение
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, statsText)
	editMsg.ParseMode = "Markdown"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", fmt.Sprintf("vpn_stats_%s", vpnID)),
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", fmt.Sprintf("vpn_config_%s", vpnID)),
		),
	)
	editMsg.ReplyMarkup = &keyboard

	if _, err := h.bot.Send(editMsg); err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

// HandleVPNRefresh обновляет данные VPN из Marzban
func (h *VPNHandler) HandleVPNRefresh(ctx context.Context, userID int64, chatID int64, messageID int, vpnID string) error {
	slog.Info("Refreshing VPN %s", vpnID)

	// Синхронизируем с Marzban через UseCase
	if err := h.vpnUC.SyncVPNStatus(ctx, vpnID); err != nil {
		return fmt.Errorf("failed to sync VPN: %w", err)
	}

	// Показываем обновленную конфигурацию
	return h.HandleGetVPNConfig(ctx, userID, chatID, vpnID)
}
