package dbtest

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	username = "postgres"
	password = "secret"
	dbName   = "test_db"
)

func CreateContainer() (*dockertest.Resource, *dockertest.Pool) {
	// Dockerに接続するプールを作成
	// プールは docker API への接続を表し、dockerイメージの作成と削除に使用される
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// スキーマファイルの絶対パスを取得
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Could not get current file path")
	}
	projectRoot := filepath.Join(filepath.Dir(filename), "../../..")
	schemaPath := filepath.Join(projectRoot, "internal/infrastructure/database/sqlc/schema.sql")
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		log.Fatalf("Could not get absolute path: %s", err)
	}

	// Dockerコンテナ起動時の細かいオプションを指定する
	runOptions := &dockertest.RunOptions{
		Repository: "postgis/postgis",
		Tag:        "18-3.6",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + username,
			"POSTGRES_DB=" + dbName,
		},
		Mounts: []string{
			absSchemaPath + ":/docker-entrypoint-initdb.d/schema.sql",
		},
	}

	// コンテナを起動
	resource, err := pool.RunWithOptions(runOptions,func(config *docker.HostConfig) {
		// 処理が終了したらコンテナを自動削除する設定
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	return resource, pool
}

func CloseContainer(resource *dockertest.Resource, pool *dockertest.Pool) {
	// コンテナの終了
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func ConnectDB(resource *dockertest.Resource, pool *dockertest.Pool) *pgxpool.Pool {
	// コンテナが起動して DB が受け付け可能になるまでリトライ
	var p *pgxpool.Pool

	hostAndPort := resource.GetHostPort("5432/tcp") // "localhost:12345" の形式で返される
	//データベース接続用のURLを作成
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, hostAndPort, dbName)

	log.Printf("Connecting to database at %s", databaseUrl)

	if err := pool.Retry(func() error {
		var err error

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		config, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return err
		}

		// 新しい接続プールを作成
		p, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			log.Printf("Failed to create connection pool: %v, retrying...", err)
			return err
		}
		// 疎通確認
		if err := p.Ping(ctx); err != nil {
			log.Printf("Failed to ping database: %v, retrying...", err)
			return err
		}
		log.Println("Successfully connected to database")
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	return p
}