package up

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"                    // migrate: マイグレーション機能
    "github.com/golang-migrate/migrate/v4/database/postgres"  // postgres: PostgreSQLデータベース対応
    "github.com/golang-migrate/migrate/v4/source/file"        // file: ファイルソース機能
    "github.com/joho/godotenv"                               // godotenv: 環境変数読み込み
    _ "github.com/lib/pq"                                    // pq: PostgreSQLドライバー（blank import）
)
)

func migrate_up() {
	fmt.Println("up")

}
