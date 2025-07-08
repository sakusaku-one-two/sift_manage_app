-- PostgreSQL Database Initialization Script
-- PostgreSQL データベース初期化スクリプト
-- initialization: 初期化、script: スクリプト

-- Create database if not exists (handled by POSTGRES_DB environment variable)
-- database: データベース、exists: 存在する、handled: 処理される、environment: 環境、variable: 変数
-- Note: Database creation is automatically handled by PostgreSQL Docker image
-- note: 注意、creation: 作成、automatically: 自動的に、handled: 処理される、image: イメージ

-- Set timezone to UTC
-- timezone: タイムゾーン、utc: 協定世界時
SET timezone = 'UTC';

-- Enable necessary PostgreSQL extensions
-- enable: 有効にする、necessary: 必要な、extensions: 拡張機能（複数形）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- uuid: 汎用一意識別子、ossp: UUID生成機能
CREATE EXTENSION IF NOT EXISTS "pgcrypto";   -- pgcrypto: PostgreSQL暗号化機能

-- Create application schema
-- create: 作成する、application: アプリケーション、schema: スキーマ
CREATE SCHEMA IF NOT EXISTS app;

-- Grant privileges to application user
-- grant: 付与する、privileges: 権限（複数形）、user: ユーザー
GRANT ALL PRIVILEGES ON SCHEMA app TO sift_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA app TO sift_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app TO sift_user;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA app TO sift_user;

-- Set default privileges for future objects
-- default: デフォルト、future: 将来の、objects: オブジェクト（複数形）
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL ON TABLES TO sift_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL ON SEQUENCES TO sift_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL ON FUNCTIONS TO sift_user;

-- Create migration history table
-- migration: マイグレーション、history: 履歴、table: テーブル
CREATE TABLE IF NOT EXISTS app.schema_migrations (
    version BIGINT PRIMARY KEY,           -- version: バージョン、bigint: 大整数型、primary: 主、key: キー
    dirty BOOLEAN NOT NULL DEFAULT FALSE, -- dirty: 汚れた、boolean: 真偽値型、not: しない、null: ヌル、default: デフォルト、false: 偽
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- applied: 適用された、at: に、timestamp: タイムスタンプ、zone: ゾーン、current: 現在の
);

-- Create users table example
-- users: ユーザー（複数形）、example: 例
CREATE TABLE IF NOT EXISTS app.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),                    -- id: 識別子、uuid: 汎用一意識別子、generate: 生成する
    email VARCHAR(255) UNIQUE NOT NULL,                                -- email: メールアドレス、varchar: 可変長文字列、unique: 一意
    password_hash VARCHAR(255) NOT NULL,                               -- password: パスワード、hash: ハッシュ値
    first_name VARCHAR(100),                                           -- first: 最初の、name: 名前
    last_name VARCHAR(100),                                            -- last: 最後の
    is_active BOOLEAN DEFAULT TRUE,                                    -- active: アクティブ、活発な、true: 真
    is_verified BOOLEAN DEFAULT FALSE,                                 -- verified: 検証済み
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,    -- created: 作成された
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP     -- updated: 更新された
);

-- Create index on email for faster lookups
-- index: インデックス、faster: より速い、lookups: 検索（複数形）
CREATE INDEX IF NOT EXISTS idx_users_email ON app.users(email);

-- Create index on active status
-- status: 状態
CREATE INDEX IF NOT EXISTS idx_users_active ON app.users(is_active);

-- Create sessions table example
-- sessions: セッション（複数形）
CREATE TABLE IF NOT EXISTS app.sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES app.users(id) ON DELETE CASCADE, -- references: 参照する、delete: 削除、cascade: カスケード
    token_hash VARCHAR(255) UNIQUE NOT NULL,                           -- token: トークン
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,                      -- expires: 期限切れ、at: に
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for sessions
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON app.sessions(user_id);

-- Create index on token_hash for sessions
CREATE INDEX IF NOT EXISTS idx_sessions_token ON app.sessions(token_hash);

-- Create index on expires_at for cleanup
-- cleanup: 清掃、クリーンアップ
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON app.sessions(expires_at);

-- Create logs table for application logging
-- logs: ログ（複数形）、logging: ログ記録
CREATE TABLE IF NOT EXISTS app.application_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level VARCHAR(20) NOT NULL,                                        -- level: レベル、ログレベル
    message TEXT NOT NULL,                                              -- message: メッセージ、text: テキスト型
    context JSONB,                                                      -- context: コンテキスト、jsonb: JSON Binary型
    user_id UUID REFERENCES app.users(id) ON DELETE SET NULL,          -- set: 設定する、null: ヌル
    ip_address INET,                                                    -- ip: Internet Protocol、address: アドレス、inet: IP アドレス型
    user_agent TEXT,                                                    -- agent: エージェント、ユーザーエージェント
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on logs for efficient querying
-- efficient: 効率的な、querying: クエリ実行
CREATE INDEX IF NOT EXISTS idx_logs_level ON app.application_logs(level);
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON app.application_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_logs_user_id ON app.application_logs(user_id);

-- Create function to update updated_at timestamp
-- function: 関数、update: 更新する
CREATE OR REPLACE FUNCTION app.update_updated_at_column()          -- replace: 置換する
RETURNS TRIGGER AS $$                                               -- returns: 返す、trigger: トリガー
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;                             -- new: 新しい
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;                                                -- language: 言語、plpgsql: PostgreSQL手続き言語

-- Create trigger to automatically update updated_at
-- trigger: トリガー、automatically: 自動的に
CREATE TRIGGER update_users_updated_at                              -- trigger: トリガー
    BEFORE UPDATE ON app.users                                      -- before: 前に、update: 更新、on: に対して
    FOR EACH ROW                                                     -- each: 各、row: 行
    EXECUTE FUNCTION app.update_updated_at_column();                 -- execute: 実行する、function: 関数

-- Create database user with limited privileges for read-only access
-- limited: 制限された、privileges: 権限、read-only: 読み取り専用、access: アクセス
-- CREATE USER readonly_user WITH PASSWORD 'readonly_password_2024';
-- GRANT CONNECT ON DATABASE sift_app_db TO readonly_user;
-- GRANT USAGE ON SCHEMA app TO readonly_user;
-- GRANT SELECT ON ALL TABLES IN SCHEMA app TO readonly_user;

-- Insert initial admin user (for development/testing only)
-- insert: 挿入する、initial: 初期の、admin: 管理者、development: 開発、testing: テスト、only: のみ
INSERT INTO app.users (email, password_hash, first_name, last_name, is_active, is_verified)
VALUES (
    'admin@siftapp.com',                                            -- 管理者メールアドレス
    crypt('admin_password_2024', gen_salt('bf')),                   -- crypt: 暗号化、gen_salt: ソルト生成、bf: Blowfish暗号化
    'Admin',                                                        -- 管理者名
    'User',                                                         -- 管理者姓
    TRUE,                                                           -- アクティブ状態
    TRUE                                                            -- 検証済み状態
) ON CONFLICT (email) DO NOTHING;                                   -- conflict: 競合、do: する、nothing: 何もしない

-- Log initialization completion
-- log: ログ、completion: 完了
INSERT INTO app.application_logs (level, message, context)
VALUES (
    'INFO',                                                         -- 情報レベル
    'Database initialization completed successfully',                -- 初期化完了メッセージ
    '{"event": "database_init", "timestamp": "' || CURRENT_TIMESTAMP || '"}'::JSONB  -- イベント情報
);

-- Display initialization summary
-- display: 表示する、summary: 要約
DO $$                                                               -- do: 実行する、無名コードブロック
BEGIN
    RAISE NOTICE 'PostgreSQL database initialization completed';     -- raise: 発生させる、notice: 通知
    RAISE NOTICE 'Database: sift_app_db';
    RAISE NOTICE 'Schema: app';
    RAISE NOTICE 'User: sift_user';
    RAISE NOTICE 'Extensions: uuid-ossp, pgcrypto';
    RAISE NOTICE 'Tables created: users, sessions, application_logs, schema_migrations';
END $$; 