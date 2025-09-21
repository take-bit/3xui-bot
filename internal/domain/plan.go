package domain

import "time"

// Plan представляет план подписки
type Plan struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Devices   int               `json:"devices"`
	Prices    map[string]Prices `json:"prices"` // валюта -> цены по периодам
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Prices представляет цены для разных периодов
type Prices struct {
	Days30  int `json:"30"`  // цена за 30 дней
	Days60  int `json:"60"`  // цена за 60 дней
	Days180 int `json:"180"` // цена за 180 дней
	Days365 int `json:"365"` // цена за 365 дней
}

// GetPrice возвращает цену для указанного количества дней
func (p *Plan) GetPrice(currency string, days int) (int, bool) {
	prices, exists := p.Prices[currency]
	if !exists {
		return 0, false
	}

	switch days {
	case 30:
		return prices.Days30, true
	case 60:
		return prices.Days60, true
	case 180:
		return prices.Days180, true
	case 365:
		return prices.Days365, true
	default:
		return 0, false
	}
}

// GetSupportedDurations возвращает поддерживаемые периоды подписки
func (p *Plan) GetSupportedDurations() []int {
	return []int{30, 60, 180, 365}
}
