package database

import (
	"context" // context: コンテキスト、処理の文脈情報
	"os"      // os: operating system（オペレーティングシステム）
	"testing" // testing: テスト機能
	"time"    // time: 時間操作機能
)

// TestPostgreSQLDriverIntegration tests PostgreSQL driver with actual database
// TestPostgreSQLDriverIntegration: 実際のデータベースでPostgreSQLドライバーをテストする統合テスト関数
// integration: 統合、actual: 実際の
func TestPostgreSQLDriverIntegration(t *testing.T) {
	// Skip integration test if not in CI/CD or integration test environment
	// skip: スキップ、integration: 統合、ci/cd: 継続的インテグレーション/継続的デプロイメント、environment: 環境
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=1 to run") // skipping: スキップしている
	}

	// Set up test environment variables for Docker Compose setup
	// set: 設定する、up: 上に、variables: 変数、docker: Docker、compose: 構成
	testEnvVars := map[string]string{
		"DB_HOST":     "localhost", // localhost: ローカルホスト
		"DB_PORT":     "5432",
		"DB_USER":     "sift_user",
		"DB_PASSWORD": "sift_password_2024",
		"DB_NAME":     "sift_app_db",
		"DB_SSL_MODE": "disable",
	}

	// Set environment variables for test
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

	// Create PostgreSQL driver instance
	// create: 作成する、instance: インスタンス
	driver, err := NewPostgreSQLDriver()
	if err != nil {
		t.Fatalf("Failed to create PostgreSQL driver: %v", err) // failed: 失敗した
	}

	// Test database connection
	// connection: 接続
	if err := driver.Connect(); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer driver.Close() // defer: 遅延実行、close: 閉じる

	// Verify connection is active
	// verify: 検証する、active: アクティブ
	if !driver.IsConnected() {
		t.Error("Expected database connection to be active") // expected: 期待した
	}

	// Test basic database operations
	// basic: 基本的な、operations: 操作（複数形）
	t.Run("TestBasicOperations", func(t *testing.T) {
		testBasicDatabaseOperations(t, driver)
	})

	// Test connection recovery
	// recovery: 回復
	t.Run("TestConnectionRecovery", func(t *testing.T) {
		testConnectionRecovery(t, driver)
	})

	// Test connection statistics
	// statistics: 統計
	t.Run("TestConnectionStats", func(t *testing.T) {
		testConnectionStatistics(t, driver)
	})
}

// testBasicDatabaseOperations tests basic CRUD operations
// testBasicDatabaseOperations: 基本的なCRUD操作をテストする関数
// crud: Create, Read, Update, Delete（作成、読み取り、更新、削除）
func testBasicDatabaseOperations(t *testing.T, driver *PostgreSQLDriver) {
	db := driver.GetDB()
	if db == nil {
		t.Fatal("Database connection is nil") // nil: ヌル値
		return
	}

	// Test simple query execution
	// simple: 簡単な、query: クエリ、execution: 実行
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // context: コンテキスト、timeout: タイムアウト、background: バックグラウンド、second: 秒
	defer cancel()

	// Test SELECT query
	var version string
	err := db.QueryRowContext(ctx, "SELECT version()").Scan(&version) // queryrow: クエリ行、scan: スキャン
	if err != nil {
		t.Errorf("Failed to execute SELECT query: %v", err) // execute: 実行する、select: 選択
		return
	}

	if version == "" {
		t.Error("Expected non-empty version string") // non-empty: 空でない、string: 文字列
	}

	t.Logf("PostgreSQL version: %s", version) // logf: ログフォーマット

	// Test table existence in app schema
	// table: テーブル、existence: 存在、schema: スキーマ
	var tableExists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = 'app' AND table_name = 'users'
		)` // exists: 存在する、information: 情報、schema: スキーマ、where: どこで

	err = db.QueryRowContext(ctx, query).Scan(&tableExists)
	if err != nil {
		t.Errorf("Failed to check table existence: %v", err) // check: 確認する
		return
	}

	if !tableExists {
		t.Error("Expected users table to exist in app schema")
	}

	// Test insert operation (if table exists)
	// insert: 挿入、operation: 操作
	if tableExists {
		testUserID := "test-user-" + time.Now().Format("20060102150405") // format: フォーマット、日時フォーマット
		testEmail := testUserID + "@test.com"

		insertQuery := `
			INSERT INTO app.users (email, password_hash, first_name, last_name, is_active, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id` // returning: 返す、戻り値

		var insertedID string
		err = db.QueryRowContext(ctx, insertQuery, testEmail, "test_hash", "Test", "User", true, false).Scan(&insertedID) // hash: ハッシュ値
		if err != nil {
			t.Errorf("Failed to insert test user: %v", err)
			return
		}

		if insertedID == "" {
			t.Error("Expected non-empty user ID after insert")
		}

		t.Logf("Inserted test user with ID: %s", insertedID)

		// Clean up test data
		// clean: 清掃する、up: 上に、data: データ
		deleteQuery := "DELETE FROM app.users WHERE id = $1"  // delete: 削除
		_, err = db.ExecContext(ctx, deleteQuery, insertedID) // exec: 実行、execute: 実行する
		if err != nil {
			t.Errorf("Failed to clean up test user: %v", err) // clean: 清掃する、up: 上に
		}
	}
}

// testConnectionRecovery tests connection recovery functionality
// testConnectionRecovery: 接続回復機能をテストする関数
// recovery: 回復、functionality: 機能
func testConnectionRecovery(t *testing.T, driver *PostgreSQLDriver) {
	// Test initial connection state
	// initial: 初期の、state: 状態
	if !driver.IsConnected() {
		t.Error("Expected initial connection to be active")
		return
	}

	// Close connection to simulate connection loss
	// simulate: シミュレートする、loss: 損失
	originalDB := driver.GetDB()
	if originalDB != nil {
		originalDB.Close() // これは内部的な接続を閉じる
	}

	// Wait a moment for connection to be recognized as closed
	// wait: 待つ、moment: 瞬間、recognized: 認識される、closed: 閉じられた
	time.Sleep(100 * time.Millisecond) // sleep: 休止、millisecond: ミリ秒

	// Test reconnection
	// reconnection: 再接続
	err := driver.Reconnect()
	if err != nil {
		t.Errorf("Failed to reconnect to database: %v", err) // reconnect: 再接続
		return
	}

	// Verify connection is active again
	// again: 再び
	if !driver.IsConnected() {
		t.Error("Expected connection to be active after reconnect")
	}

	// Test that new connection works
	// works: 動作する
	db := driver.GetDB()
	if db == nil {
		t.Error("Expected valid database connection after reconnect") // valid: 有効な
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&result) // 簡単な接続テスト
	if err != nil {
		t.Errorf("Failed to execute query after reconnect: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected result 1, got: %d", result) // result: 結果
	}
}

// testConnectionStatistics tests connection statistics functionality
// testConnectionStatistics: 接続統計機能をテストする関数
func testConnectionStatistics(t *testing.T, driver *PostgreSQLDriver) {
	// Get connection statistics
	// statistics: 統計
	stats := driver.GetConnectionStats()

	// Verify basic statistics structure
	// structure: 構造
	if stats.MaxOpenConnections <= 0 {
		t.Error("Expected MaxOpenConnections to be greater than 0") // greater: より大きい、than: より
	}

	if stats.MaxOpenConnections > 100 {
		t.Error("Expected reasonable MaxOpenConnections limit") // reasonable: 合理的な、limit: 制限
	}

	t.Logf("Connection Statistics:") // statistics: 統計
	t.Logf("  Max Open Connections: %d", stats.MaxOpenConnections)
	t.Logf("  Open Connections: %d", stats.OpenConnections)
	t.Logf("  In Use: %d", stats.InUse) // in: 中に、use: 使用
	t.Logf("  Idle: %d", stats.Idle)    // idle: アイドル、待機中

	// Test that we have at least one open connection
	// least: 最少、one: 1つ、open: 開いている
	if stats.OpenConnections < 1 {
		t.Error("Expected at least one open connection")
	}
}

// TestDriverWithDockerCompose tests driver integration with Docker Compose setup
// TestDriverWithDockerCompose: Docker Compose設定でのドライバー統合をテストする関数
func TestDriverWithDockerCompose(t *testing.T) {
	// This test is designed to run with docker-compose up
	// designed: 設計された、run: 実行する、up: 起動
	if os.Getenv("DOCKER_COMPOSE_TEST") == "" {
		t.Skip("Skipping Docker Compose integration test. Set DOCKER_COMPOSE_TEST=1 to run")
	}

	// Use Docker Compose environment variables
	// use: 使用する
	testEnvVars := map[string]string{
		"DB_HOST":     "postgres", // Docker Compose service name
		"DB_PORT":     "5432",
		"DB_USER":     "sift_user",
		"DB_PASSWORD": "sift_password_2024",
		"DB_NAME":     "sift_app_db",
		"DB_SSL_MODE": "disable",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	// Test driver creation and connection
	// creation: 作成
	driver, err := NewPostgreSQLDriver()
	if err != nil {
		t.Fatalf("Failed to create driver for Docker Compose test: %v", err)
	}

	// Test connection with retry logic for Docker Compose startup
	// retry: 再試行、logic: ロジック、startup: 起動
	var connected bool
	maxRetries := 10 // maximum: 最大、retries: 再試行（複数形）
	for i := 0; i < maxRetries; i++ {
		if err := driver.Connect(); err == nil {
			connected = true
			break
		}
		t.Logf("Connection attempt %d failed, retrying...", i+1) // attempt: 試行、retrying: 再試行
		time.Sleep(2 * time.Second)
	}

	if !connected {
		t.Fatalf("Failed to connect to PostgreSQL in Docker Compose after %d attempts", maxRetries) // attempts: 試行（複数形）
	}

	defer driver.Close()

	// Verify database initialization
	// initialization: 初期化
	db := driver.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if admin user was created
	// admin: 管理者、created: 作成された
	var adminExists bool
	checkAdminQuery := "SELECT EXISTS(SELECT 1 FROM app.users WHERE email = 'admin@siftapp.com')"
	err = db.QueryRowContext(ctx, checkAdminQuery).Scan(&adminExists)
	if err != nil {
		t.Errorf("Failed to check admin user existence: %v", err)
		return
	}

	if !adminExists {
		t.Error("Expected admin user to be created during database initialization")
	}

	t.Log("Docker Compose integration test completed successfully") // completed: 完了した、successfully: 成功して
}
