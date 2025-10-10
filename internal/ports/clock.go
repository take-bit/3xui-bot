package ports

import "time"

// Clock интерфейс для работы со временем (для тестируемости)
type Clock interface {
	Now() time.Time
}

// SystemClock реальная реализация часов
type SystemClock struct{}

// Now возвращает текущее время
func (c *SystemClock) Now() time.Time {
	return time.Now()
}
