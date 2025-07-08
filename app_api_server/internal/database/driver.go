package database

import (
	"database/sql" // sql: データベース操作用パッケージ、Structured Query Language（構造化照会言語）
	"fmt"          // fmt: format（フォーマット）、文字列フォーマット機能
	"log"          // log: ログ出力機能
	"os"           // os: operating system（オペレーティングシステム）、OS操作機能
	"strconv"      // strconv: string conversion（文字列変換）、文字列と数値の変換
	"time"         // time: 時間操作機能

	"github.com/joho/godotenv" // godotenv: 環境変数読み込み
	_ "github.com/lib/pq"      // pq: PostgreSQLドライバー（blank import）
)

// DatabaseConfig represents database configuration settings
// DatabaseConfig: データベース設定を表す構造体
// represents: 表現する、configuration: 設定、settings: 設定（複数形）
type DatabaseConfig struct {
	Host     string // host: ホスト、データベースサーバーのアドレス
	Port     int    // port: ポート、接続用のポート番号
	User     string // user: ユーザー、データベースユーザー名
	Password string // password: パスワード、認証用パスワード
	Database string // database: データベース、データベース名
	SSLMode  string // sslmode: SSL mode（セキュリティ層）、SSL接続モード
}

// PostgreSQLDriver represents PostgreSQL database driver
// PostgreSQLDriver: PostgreSQLデータベースドライバーを表す構造体
// represents: 表現する、driver: ドライバー
type PostgreSQLDriver struct {
	config *DatabaseConfig // config: 設定、configuration: 構成
	db     *sql.DB         // db: database（データベース）、データベース接続
}

// LoadDatabaseConfig loads database configuration from environment variables
// LoadDatabaseConfig: 環境変数からデータベース設定を読み込む関数
// loads: 読み込む、environment: 環境、variables: 変数（複数形）
func LoadDatabaseConfig() (*DatabaseConfig, error) {
	// Load environment variables from .env file
	// environment: 環境、variables: 変数、from: から
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err) // warning: 警告、found: 見つかった
	}

	// Get database configuration from environment variables
	// configuration: 設定
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // default: デフォルト、既定値
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "5432" // default PostgreSQL port
	}
	port, err := strconv.Atoi(portStr) // Atoi: ASCII to integer（ASCII文字列から整数へ）
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %v", err) // invalid: 無効な、number: 数
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return nil, fmt.Errorf("DB_USER environment variable is required") // required: 必要な
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	database := os.Getenv("DB_NAME")
	if database == "" {
		return nil, fmt.Errorf("DB_NAME environment variable is required")
	}

	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "require" // default: secure SSL mode
	}

	return &DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
	}, nil
}

// BuildConnectionString builds PostgreSQL connection string from configuration
// BuildConnectionString: 設定からPostgreSQL接続文字列を構築する関数
// builds: 構築する、connection: 接続、string: 文字列
func (c *DatabaseConfig) BuildConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
		c.SSLMode,
	)
}

// NewPostgreSQLDriver creates a new PostgreSQL driver instance
// NewPostgreSQLDriver: 新しいPostgreSQLドライバーインスタンスを作成するファクトリー関数
// creates: 作成する、instance: インスタンス
func NewPostgreSQLDriver() (*PostgreSQLDriver, error) {
	// Load database configuration
	// load: 読み込む
	config, err := LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err) // failed: 失敗した
	}

	// Create PostgreSQL driver instance
	// create: 作成する
	driver := &PostgreSQLDriver{
		config: config,
	}

	return driver, nil
}

// NewPostgreSQLDriverWithConfig creates a new PostgreSQL driver with custom configuration
// NewPostgreSQLDriverWithConfig: カスタム設定で新しいPostgreSQLドライバーを作成するファクトリー関数
// custom: カスタム、独自の
func NewPostgreSQLDriverWithConfig(config *DatabaseConfig) (*PostgreSQLDriver, error) {
	if config == nil {
		return nil, fmt.Errorf("database configuration cannot be nil") // cannot: できない、nil: ヌル値
	}

	// Validate configuration
	// validate: 検証する
	if err := validateDatabaseConfig(config); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err) // invalid: 無効な
	}

	driver := &PostgreSQLDriver{
		config: config,
	}

	return driver, nil
}

// validateDatabaseConfig validates database configuration
// validateDatabaseConfig: データベース設定を検証する関数
// validates: 検証する
func validateDatabaseConfig(config *DatabaseConfig) error {
	if config.Host == "" {
		return fmt.Errorf("host cannot be empty") // empty: 空の
	}

	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535") // between: 間に、must: しなければならない
	}

	if config.User == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if config.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if config.Database == "" {
		return fmt.Errorf("database name cannot be empty") // name: 名前
	}

	// Valid SSL modes for PostgreSQL
	// valid: 有効な、modes: モード（複数形）
	validSSLModes := []string{"disable", "require", "verify-ca", "verify-full"}
	validMode := false
	for _, mode := range validSSLModes {
		if config.SSLMode == mode {
			validMode = true
			break
		}
	}

	if !validMode {
		return fmt.Errorf("invalid SSL mode: %s", config.SSLMode)
	}

	return nil
}

// Connect establishes a connection to the PostgreSQL database
// Connect: PostgreSQLデータベースへの接続を確立する関数
// establishes: 確立する、connection: 接続
func (d *PostgreSQLDriver) Connect() error {
	// Build connection string
	// build: 構築する
	connectionString := d.config.BuildConnectionString()

	// Open database connection
	// open: 開く
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	// configure: 設定する、pool: プール、接続プール
	db.SetMaxOpenConns(25)                 // maximum: 最大の、open: 開いている、connections: 接続（複数形）
	db.SetMaxIdleConns(5)                  // idle: アイドル、待機中の
	db.SetConnMaxLifetime(5 * time.Minute) // lifetime: 寿命、minute: 分

	// Test database connection
	// test: テスト、試験
	if err := db.Ping(); err != nil {
		db.Close()                                            // Close database if ping fails
		return fmt.Errorf("failed to ping database: %w", err) // ping: 接続確認
	}

	d.db = db
	log.Printf("Successfully connected to PostgreSQL database: %s", d.config.Database) // successfully: 成功して
	return nil
}

// GetDB returns the database connection
// GetDB: データベース接続を返す関数
// returns: 返す
func (d *PostgreSQLDriver) GetDB() *sql.DB {
	return d.db
}

// GetConfig returns the database configuration
// GetConfig: データベース設定を返す関数
func (d *PostgreSQLDriver) GetConfig() *DatabaseConfig {
	return d.config
}

// Close closes the database connection
// Close: データベース接続を閉じる関数
// closes: 閉じる
func (d *PostgreSQLDriver) Close() error {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err) // close: 閉じる
		}
		log.Println("Database connection closed successfully")
	}
	return nil
}

// IsConnected checks if the database connection is active
// IsConnected: データベース接続がアクティブかどうかを確認する関数
// checks: 確認する、active: アクティブ、活発な
func (d *PostgreSQLDriver) IsConnected() bool {
	if d.db == nil {
		return false
	}

	// Test connection with ping
	if err := d.db.Ping(); err != nil {
		return false
	}

	return true
}

// Reconnect attempts to reconnect to the database
// Reconnect: データベースへの再接続を試行する関数
// attempts: 試行する、reconnect: 再接続
func (d *PostgreSQLDriver) Reconnect() error {
	// Close existing connection if any
	// existing: 既存の、if: もし、any: 何らかの
	if d.db != nil {
		d.db.Close()
	}

	// Attempt to reconnect
	// attempt: 試行する
	return d.Connect()
}

// GetConnectionStats returns database connection statistics
// GetConnectionStats: データベース接続統計を返す関数
// statistics: 統計
func (d *PostgreSQLDriver) GetConnectionStats() sql.DBStats {
	if d.db == nil {
		return sql.DBStats{}
	}
	return d.db.Stats()
}
