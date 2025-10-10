-- ================================================================
-- SEED DATA: Планы подписки
-- ================================================================
-- Этот файл добавляет начальные планы подписки в систему

-- Очистка существующих планов (опционально, раскомментировать если нужно)
-- DELETE FROM plans;

-- Базовые планы подписки
INSERT INTO plans (id, name, description, price, days, is_active, created_at, updated_at) VALUES
    (
        'trial_7',
        '🎁 Пробный период',
        'Бесплатный пробный период на 7 дней для новых пользователей',
        0.00,
        7,
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'monthly',
        '📅 Месячная подписка',
        '30 дней безлимитного VPN со скоростью до 100 Мбит/с',
        299.00,
        30,
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'quarterly',
        '📆 Квартальная подписка',
        '90 дней VPN со скидкой 15% • Экономия 135₽',
        749.00,
        90,
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'yearly',
        '🎯 Годовая подписка',
        '365 дней VPN со скидкой 30% • Экономия 1089₽ • Самое выгодное!',
        2499.00,
        365,
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    days = EXCLUDED.days,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- Вывод добавленных планов
SELECT 
    id,
    name,
    price || ' ₽' as price,
    days || ' дней' as duration,
    CASE WHEN is_active THEN '✓' ELSE '✗' END as active
FROM plans
ORDER BY days ASC;

-- Статистика
SELECT COUNT(*) as total_plans FROM plans WHERE is_active = true;

