package repository

import (
	"testing"

	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	dbTest "github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/db_test"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	fixtures *testfixtures.Loader
	testQueries *dbgen.Queries
)

func GetTestQueries() *dbgen.Queries {
	return testQueries
}

// TestMainはパッケージ内の全てのテストを実行する前に一度だけ実行されます
func TestMain(m *testing.M) {
	var err error

	// DBの立ち上げ
	resource, pool := dbTest.CreateContainer()
	defer dbTest.CloseContainer(resource, pool)

	// DBへ接続する（pgxpool.Poolを使用）
	dbPool := dbTest.ConnectDB(resource, pool)
	defer dbPool.Close()

	// testfixturesのために*sql.DBを取得
	sqlDB := stdlib.OpenDBFromPool(dbPool)
	defer sqlDB.Close()

	// テストデータの準備
	fixturePath := "../fixtures"
	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB), // testfixturesは*sql.DBが必要
		testfixtures.Dialect("postgres"), // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory(fixturePath), // The directory containing the YAML files
	)
	if err != nil {
		panic(err)
	}
	
	testQueries = dbgen.New(dbPool)

	// 全てのテストを実行します
	m.Run()
}

// resetTestDataは各テスト実行前に呼び出され、テストデータをリセットします
func resetTestData(t *testing.T) {
	t.Helper()
	//fixtures.Load()でテーブルの既存データを削除し、fixturesから再度データをロードする
	if err := fixtures.Load(); err != nil {
		t.Fatal(err)
	}
}