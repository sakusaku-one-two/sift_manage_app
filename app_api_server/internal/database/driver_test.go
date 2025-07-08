package database

import (
	"os"      // os: operating system（オペレーティングシステム）
	"testing" // testing: テスト機能
)

// TestLoadDatabaseConfig tests database configuration loading
// TestLoadDatabaseConfig: データベース設定読み込みをテストする関数
// tests: テストする
func TestLoadDatabaseConfig(t *testing.T) {
	// Set up test environment variables
	// set: 設定する、up: 上に、test: テスト、environment: 環境、variables: 変数
	testEnvVars := map[string]string{
		"DB_HOST":     "test-host",     // host: ホスト
		"DB_PORT":     "5433",          // port: ポート
		"DB_USER":     "test-user",     // user: ユーザー
		"DB_PASSWORD": "test-password", // password: パスワード
		"DB_NAME":     "test-db",       // name: 名前
		"DB_SSL_MODE": "disable",       // ssl: セキュリティ層、mode: モード、disable: 無効
	}

	// Set environment variables
	// environment: 環境
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up environment variables after test
	// clean: 清掃する、up: 上に、after: 後に
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	// Test configuration loading
	// configuration: 設定
	config, err := LoadDatabaseConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err) // expected: 期待した、error: エラー
	}

	// Validate loaded configuration
	// validate: 検証する、loaded: 読み込まれた
	if config.Host != "test-host" {
		t.Errorf("Expected host 'test-host', got: %s", config.Host) // expected: 期待した
	}

	if config.Port != 5433 {
		t.Errorf("Expected port 5433, got: %d", config.Port)
	}

	if config.User != "test-user" {
		t.Errorf("Expected user 'test-user', got: %s", config.User)
	}

	if config.Password != "test-password" {
		t.Errorf("Expected password 'test-password', got: %s", config.Password)
	}

	if config.Database != "test-db" {
		t.Errorf("Expected database 'test-db', got: %s", config.Database)
	}

	if config.SSLMode != "disable" {
		t.Errorf("Expected SSL mode 'disable', got: %s", config.SSLMode)
	}
}

// TestLoadDatabaseConfigDefaults tests default configuration values
// TestLoadDatabaseConfigDefaults: デフォルト設定値をテストする関数
// defaults: デフォルト値（複数形）、values: 値（複数形）
func TestLoadDatabaseConfigDefaults(t *testing.T) {
	// Set only required environment variables
	// only: のみ、required: 必要な
	requiredEnvVars := map[string]string{
		"DB_USER":     "test-user",
		"DB_PASSWORD": "test-password",
		"DB_NAME":     "test-db",
	}

	// Set environment variables
	for key, value := range requiredEnvVars {
		os.Setenv(key, value)
	}

	// Clean up environment variables after test
	defer func() {
		for key := range requiredEnvVars {
			os.Unsetenv(key)
		}
		// Also clean up any default values that might have been set
		// also: また、default: デフォルト、might: かもしれない、have: 持つ、been: された
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_SSL_MODE")
	}()

	config, err := LoadDatabaseConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check default values
	// check: 確認する、default: デフォルト
	if config.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got: %s", config.Host)
	}

	if config.Port != 5432 {
		t.Errorf("Expected default port 5432, got: %d", config.Port)
	}

	if config.SSLMode != "require" {
		t.Errorf("Expected default SSL mode 'require', got: %s", config.SSLMode)
	}
}

// TestLoadDatabaseConfigMissingRequired tests error handling for missing required variables
// TestLoadDatabaseConfigMissingRequired: 必要な変数が不足している場合のエラーハンドリングをテスト
// missing: 不足している、handling: ハンドリング、処理
func TestLoadDatabaseConfigMissingRequired(t *testing.T) {
	// Test cases for missing required environment variables
	// cases: ケース（複数形）
	testCases := []struct {
		name        string            // name: 名前
		envVars     map[string]string // environment: 環境、variables: 変数
		expectError bool              // expect: 期待する、error: エラー
	}{
		{
			name:        "Missing DB_USER",
			envVars:     map[string]string{"DB_PASSWORD": "pass", "DB_NAME": "db"},
			expectError: true,
		},
		{
			name:        "Missing DB_PASSWORD",
			envVars:     map[string]string{"DB_USER": "user", "DB_NAME": "db"},
			expectError: true,
		},
		{
			name:        "Missing DB_NAME",
			envVars:     map[string]string{"DB_USER": "user", "DB_PASSWORD": "pass"},
			expectError: true,
		},
		{
			name:        "All required present",
			envVars:     map[string]string{"DB_USER": "user", "DB_PASSWORD": "pass", "DB_NAME": "db"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up all environment variables first
			// clean: 清掃する、up: 上に、all: 全て、first: 最初に
			envVarsToClean := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE"}
			for _, env := range envVarsToClean {
				os.Unsetenv(env)
			}

			// Set test environment variables
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				for key := range tc.envVars {
					os.Unsetenv(key)
				}
			}()

			// Test configuration loading
			_, err := LoadDatabaseConfig()

			if tc.expectError && err == nil {
				t.Errorf("Expected error for test case '%s', but got none", tc.name)
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for test case '%s', but got: %v", tc.name, err)
			}
		})
	}
}

// TestBuildConnectionString tests connection string building
// TestBuildConnectionString: 接続文字列構築をテストする関数
// string: 文字列、building: 構築
func TestBuildConnectionString(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "require",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=require"
	actual := config.BuildConnectionString()

	if actual != expected {
		t.Errorf("Expected connection string '%s', got: '%s'", expected, actual)
	}
}

// TestValidateDatabaseConfig tests database configuration validation
// TestValidateDatabaseConfig: データベース設定検証をテストする関数
// validation: 検証
func TestValidateDatabaseConfig(t *testing.T) {
	// Test cases for configuration validation
	// validation: 検証
	testCases := []struct {
		name        string
		config      *DatabaseConfig
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "require",
			},
			expectError: false,
		},
		{
			name: "Empty host",
			config: &DatabaseConfig{
				Host:     "",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "require",
			},
			expectError: true,
		},
		{
			name: "Invalid port - too low",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     0,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "require",
			},
			expectError: true,
		},
		{
			name: "Invalid port - too high",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     65536,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "require",
			},
			expectError: true,
		},
		{
			name: "Empty user",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "",
				Password: "pass",
				Database: "db",
				SSLMode:  "require",
			},
			expectError: true,
		},
		{
			name: "Empty password",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "",
				Database: "db",
				SSLMode:  "require",
			},
			expectError: true,
		},
		{
			name: "Empty database",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "",
				SSLMode:  "require",
			},
			expectError: true,
		},
		{
			name: "Invalid SSL mode",
			config: &DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "pass",
				Database: "db",
				SSLMode:  "invalid",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDatabaseConfig(tc.config)

			if tc.expectError && err == nil {
				t.Errorf("Expected error for test case '%s', but got none", tc.name)
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for test case '%s', but got: %v", tc.name, err)
			}
		})
	}
}

// TestNewPostgreSQLDriver tests PostgreSQL driver factory function
// TestNewPostgreSQLDriver: PostgreSQLドライバーファクトリー関数をテストする関数
// factory: ファクトリー、工場
func TestNewPostgreSQLDriver(t *testing.T) {
	// Set up test environment variables
	testEnvVars := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_USER":     "testuser",
		"DB_PASSWORD": "testpass",
		"DB_NAME":     "testdb",
		"DB_SSL_MODE": "disable",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up environment variables after test
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	// Test driver creation
	// creation: 作成
	driver, err := NewPostgreSQLDriver()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if driver == nil {
		t.Fatal("Expected driver instance, got nil") // instance: インスタンス
	}

	// Test driver configuration
	config := driver.GetConfig()
	if config.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got: %s", config.Host)
	}
}

// TestNewPostgreSQLDriverWithConfig tests PostgreSQL driver factory with custom config
// TestNewPostgreSQLDriverWithConfig: カスタム設定でのPostgreSQLドライバーファクトリーをテスト
// custom: カスタム、独自の
func TestNewPostgreSQLDriverWithConfig(t *testing.T) {
	// Test with valid configuration
	// valid: 有効な
	validConfig := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "require",
	}

	driver, err := NewPostgreSQLDriverWithConfig(validConfig)
	if err != nil {
		t.Fatalf("Expected no error with valid config, got: %v", err)
	}

	if driver == nil {
		t.Fatal("Expected driver instance, got nil")
	}

	// Test with nil configuration
	_, err = NewPostgreSQLDriverWithConfig(nil)
	if err == nil {
		t.Error("Expected error with nil config, got none")
	}

	// Test with invalid configuration
	invalidConfig := &DatabaseConfig{
		Host:     "",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "require",
	}

	_, err = NewPostgreSQLDriverWithConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error with invalid config, got none")
	}
}

// TestDriverMethods tests driver methods without actual database connection
// TestDriverMethods: 実際のデータベース接続なしでドライバーメソッドをテストする関数
// methods: メソッド（複数形）、without: なしで、actual: 実際の
func TestDriverMethods(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "require",
	}

	driver, err := NewPostgreSQLDriverWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to create driver: %v", err)
	}

	// Test GetConfig method
	// method: メソッド
	retrievedConfig := driver.GetConfig()
	if retrievedConfig.Host != config.Host {
		t.Errorf("Expected host '%s', got: '%s'", config.Host, retrievedConfig.Host)
	}

	// Test GetDB method (should return nil since not connected)
	// since: なぜなら、connected: 接続された
	db := driver.GetDB()
	if db != nil {
		t.Error("Expected nil database connection, got non-nil")
	}

	// Test IsConnected method (should return false since not connected)
	if driver.IsConnected() {
		t.Error("Expected IsConnected to return false, got true")
	}

	// Test Close method (should not error even if not connected)
	// even: たとえ、if: もし
	if err := driver.Close(); err != nil {
		t.Errorf("Expected no error on Close, got: %v", err)
	}

	// Test GetConnectionStats method (should return empty stats)
	// empty: 空の、stats: 統計
	stats := driver.GetConnectionStats()
	if stats.OpenConnections != 0 {
		t.Errorf("Expected 0 open connections, got: %d", stats.OpenConnections)
	}
}
