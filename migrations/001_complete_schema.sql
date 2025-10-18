CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    language_code VARCHAR(10),
    is_blocked BOOLEAN DEFAULT FALSE,
    has_trial BOOLEAN DEFAULT FALSE, -- Использовал ли пользователь пробный период
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plans (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    days INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR(50) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    plan_id VARCHAR(50) NOT NULL REFERENCES plans(id),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(50) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    payment_method VARCHAR(255),
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- VPN ПОДКЛЮЧЕНИЯ (MARZBAN)
-- =============================================================================

-- Таблица VPN подключений (связь с Marzban)
CREATE TABLE IF NOT EXISTS vpn_connections (
    id VARCHAR(50) PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    marzban_username VARCHAR(100) NOT NULL UNIQUE, -- Username в Marzban
    name VARCHAR(255), -- Локальное имя подключения
    is_active BOOLEAN DEFAULT TRUE, -- Флаг активности в нашей системе
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- РЕФЕРАЛЬНАЯ СИСТЕМА
-- =============================================================================

-- Таблица реферальных связей
CREATE TABLE IF NOT EXISTS referrals (
    id BIGSERIAL PRIMARY KEY,
    referrer_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    referee_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(referrer_id, referee_id)
);

-- Таблица реферальных ссылок
CREATE TABLE IF NOT EXISTS referral_links (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    link VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- УВЕДОМЛЕНИЯ
-- =============================================================================

-- Таблица уведомлений
CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(50) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- ИНДЕКСЫ
-- =============================================================================

-- Индексы для пользователей
-- Индекс на telegram_id не нужен - это уже PRIMARY KEY

-- Индексы для подписок
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_active ON subscriptions(user_id, is_active, end_date);

-- Индексы для платежей
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

-- Индексы для VPN подключений
CREATE INDEX IF NOT EXISTS idx_vpn_connections_telegram_user_id ON vpn_connections(telegram_user_id);
CREATE INDEX IF NOT EXISTS idx_vpn_connections_marzban_username ON vpn_connections(marzban_username);
CREATE INDEX IF NOT EXISTS idx_vpn_connections_is_active ON vpn_connections(is_active);

-- Индексы для рефералов
CREATE INDEX IF NOT EXISTS idx_referrals_referrer_id ON referrals(referrer_id);
CREATE INDEX IF NOT EXISTS idx_referrals_referee_id ON referrals(referee_id);
CREATE INDEX IF NOT EXISTS idx_referral_links_user_id ON referral_links(user_id);
CREATE INDEX IF NOT EXISTS idx_referral_links_link ON referral_links(link);

-- Индексы для уведомлений
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);

-- =============================================================================
-- КОММЕНТАРИИ
-- =============================================================================

-- Комментарии к таблицам
COMMENT ON TABLE users IS 'Пользователи Telegram';
COMMENT ON TABLE plans IS 'Тарифные планы подписок';
COMMENT ON TABLE subscriptions IS 'Подписки пользователей';
COMMENT ON TABLE payments IS 'Платежи пользователей';
COMMENT ON TABLE vpn_connections IS 'VPN подключения пользователей в Marzban (только связи и локальные данные)';
COMMENT ON TABLE referrals IS 'Реферальные связи между пользователями';
COMMENT ON TABLE referral_links IS 'Реферальные ссылки пользователей';
COMMENT ON TABLE notifications IS 'Уведомления для пользователей';

-- Комментарии к ключевым полям
COMMENT ON COLUMN users.telegram_id IS 'Уникальный ID пользователя в Telegram';
COMMENT ON COLUMN users.has_trial IS 'Использовал ли пользователь пробный период (ограничение: один раз)';
COMMENT ON COLUMN users.created_at IS 'Дата создания аккаунта пользователя';

COMMENT ON COLUMN plans.price IS 'Цена плана в рублях';
COMMENT ON COLUMN plans.days IS 'Количество дней действия плана';

COMMENT ON COLUMN subscriptions.name IS 'Название подписки (задается пользователем)';
COMMENT ON COLUMN subscriptions.start_date IS 'Дата начала подписки';
COMMENT ON COLUMN subscriptions.end_date IS 'Дата окончания подписки';

COMMENT ON COLUMN payments.amount IS 'Сумма платежа в рублях';
COMMENT ON COLUMN payments.currency IS 'Валюта платежа';
COMMENT ON COLUMN payments.status IS 'Статус платежа: pending, completed, failed, cancelled';

COMMENT ON COLUMN vpn_connections.telegram_user_id IS 'ID пользователя Telegram';
COMMENT ON COLUMN vpn_connections.marzban_username IS 'Уникальный username в Marzban API';
COMMENT ON COLUMN vpn_connections.name IS 'Локальное имя подключения для пользователя';
COMMENT ON COLUMN vpn_connections.is_active IS 'Флаг активности подключения в нашей системе';

COMMENT ON COLUMN referrals.referrer_id IS 'ID пользователя, который пригласил';
COMMENT ON COLUMN referrals.referee_id IS 'ID пользователя, которого пригласили';

COMMENT ON COLUMN notifications.type IS 'Тип уведомления: info, warning, error';
COMMENT ON COLUMN notifications.is_read IS 'Прочитано ли уведомление';

-- =============================================================================
-- БАЗОВЫЕ ДАННЫЕ
-- =============================================================================

-- Заполнение таблицы планов базовыми данными
INSERT INTO plans (id, name, description, price, days, is_active) VALUES
('plan_1w', '1 неделя', 'Пробная подписка на 1 неделю', 25.00, 7, true),
('plan_1m', '1 месяц', 'Подписка на 1 месяц', 100.00, 30, true),
('plan_3m', '3 месяца', 'Подписка на 3 месяца', 250.00, 90, true),
('plan_6m', '6 месяцев', 'Подписка на 6 месяцев', 450.00, 180, true),
('plan_12m', '12 месяцев', 'Подписка на 12 месяцев', 800.00, 365, true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    days = EXCLUDED.days,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;
