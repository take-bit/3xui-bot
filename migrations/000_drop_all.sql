-- ================================================================
-- DROP ALL TABLES (Down Migration)
-- ================================================================
-- Этот файл удаляет все таблицы для чистой миграции

-- Удаляем таблицы в обратном порядке (из-за foreign key constraints)
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS referral_links CASCADE;
DROP TABLE IF EXISTS referrals CASCADE;
DROP TABLE IF EXISTS vpn_connections CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS subscriptions CASCADE;
DROP TABLE IF EXISTS plans CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Удаляем типы (если есть)
DROP TYPE IF EXISTS payment_status CASCADE;
DROP TYPE IF EXISTS subscription_status CASCADE;
DROP TYPE IF EXISTS notification_type CASCADE;

-- Сообщение об успешном удалении
SELECT 'All tables dropped successfully' AS status;

