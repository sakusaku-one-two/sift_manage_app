package database

import (
	"log" // log: ログ出力機能
)

// ExampleUsage demonstrates how to use the PostgreSQL driver
// ExampleUsage: PostgreSQLドライバーの使用方法を示すサンプル関数
// demonstrates: 実演する、使用方法を示す
func ExampleUsage() {
	// Method 1: Create driver with environment variables
	// method: 方法、create: 作成する、environment: 環境、variables: 変数
	driver, err := NewPostgreSQLDriver()
	if err != nil {
		log.Printf("Failed to create driver: %v", err) // failed: 失敗した
		return
	}

	// Connect to database
	// connect: 接続する
	if err := driver.Connect(); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer driver.Close() // defer: 遅延実行、close: 閉じる

	// Check if connected
	// check: 確認する、connected: 接続された
	if driver.IsConnected() {
		log.Println("Successfully connected to database") // successfully: 成功して
	}

	// Get database connection for queries
	// queries: クエリ（複数形）、問い合わせ
	db := driver.GetDB()
	if db != nil {
		log.Println("Database connection is available") // available: 利用可能な
	}

	// Get connection statistics
	// statistics: 統計
	stats := driver.GetConnectionStats()
	log.Printf("Open connections: %d", stats.OpenConnections) // open: 開いている、connections: 接続
}

// ExampleUsageWithCustomConfig demonstrates how to use the driver with custom configuration
// ExampleUsageWithCustomConfig: カスタム設定でのドライバー使用方法を示すサンプル関数
// custom: カスタム、独自の
func ExampleUsageWithCustomConfig() {
	// Method 2: Create driver with custom configuration
	// custom: カスタム、独自の
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		Database: "mydatabase",
		SSLMode:  "require",
	}

	driver, err := NewPostgreSQLDriverWithConfig(config)
	if err != nil {
		log.Printf("Failed to create driver with custom config: %v", err)
		return
	}

	// Connect to database
	if err := driver.Connect(); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer driver.Close()

	// Use the database connection
	// use: 使用する
	db := driver.GetDB()
	if db != nil {
		// Example query execution
		// example: 例、query: クエリ、execution: 実行
		rows, err := db.Query("SELECT version()") // version: バージョン
		if err != nil {
			log.Printf("Failed to execute query: %v", err) // execute: 実行する
			return
		}
		defer rows.Close() // rows: 行（複数形）

		// Process query results
		// process: 処理する、results: 結果（複数形）
		for rows.Next() {
			var version string
			if err := rows.Scan(&version); err != nil {
				log.Printf("Failed to scan row: %v", err) // scan: スキャン、row: 行
				continue
			}
			log.Printf("PostgreSQL version: %s", version)
		}
	}
}

// ExampleReconnection demonstrates connection recovery
// ExampleReconnection: 接続回復を実演するサンプル関数
// demonstrates: 実演する、recovery: 回復
func ExampleReconnection() {
	driver, err := NewPostgreSQLDriver()
	if err != nil {
		log.Printf("Failed to create driver: %v", err)
		return
	}

	// Initial connection
	// initial: 初期の
	if err := driver.Connect(); err != nil {
		log.Printf("Failed to connect initially: %v", err) // initially: 初期に
		return
	}

	// Simulate connection loss and recovery
	// simulate: シミュレートする、loss: 損失、recovery: 回復
	log.Println("Simulating connection loss...") // simulating: シミュレートしている
	driver.Close()

	// Attempt to reconnect
	// attempt: 試行する、reconnect: 再接続
	log.Println("Attempting to reconnect...") // attempting: 試行している
	if err := driver.Reconnect(); err != nil {
		log.Printf("Failed to reconnect: %v", err)
		return
	}

	log.Println("Successfully reconnected to database")
	driver.Close()
}

// ExampleEnvironmentVariables shows required environment variables
// ExampleEnvironmentVariables: 必要な環境変数を示すサンプル関数
// shows: 示す、required: 必要な
func ExampleEnvironmentVariables() {
	log.Println("Required environment variables:")                     // required: 必要な
	log.Println("DB_HOST=localhost (optional, defaults to localhost)") // optional: オプション、defaults: デフォルト
	log.Println("DB_PORT=5432 (optional, defaults to 5432)")
	log.Println("DB_USER=your_username (required)")
	log.Println("DB_PASSWORD=your_password (required)")
	log.Println("DB_NAME=your_database (required)")
	log.Println("DB_SSL_MODE=require (optional, defaults to require)")
	log.Println("")
	log.Println("Valid SSL modes: disable, require, verify-ca, verify-full") // valid: 有効な, modes: モード
}
