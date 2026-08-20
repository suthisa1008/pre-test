package testdb

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	pg "github.com/demo/product-api/internal/repository/postgres"
)

func Start(t *testing.T) *gorm.DB {
	t.Helper()
	if os.Getenv("SKIP_INTEGRATION") == "1" {
		t.Skip("SKIP_INTEGRATION=1")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	schema := "t_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	for _, dsn := range candidateURLs() {
		db, err := openIsolated(dsn, schema)
		if err == nil {
			t.Cleanup(func() {
				_ = db.Exec("DROP SCHEMA IF EXISTS " + schema + " CASCADE").Error
				_ = pg.Close(db)
			})
			return db
		}
	}

	dsn, err := startContainer(ctx)
	if err != nil {
		t.Skipf("postgres unavailable (start docker compose or Testcontainers): %v", err)
	}
	db, err := openIsolated(dsn, schema)
	if err != nil {
		t.Skipf("postgres migrate failed: %v", err)
	}
	t.Cleanup(func() { _ = pg.Close(db) })
	return db
}

func candidateURLs() []string {
	var urls []string
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		urls = append(urls, v)
	}
	urls = append(urls, "postgres://product:product@localhost:5433/product?sslmode=disable")
	return urls
}

func openIsolated(dsn, schema string) (*gorm.DB, error) {
	admin, err := pg.Open(dsn)
	if err != nil {
		return nil, err
	}
	if err := admin.Exec("CREATE SCHEMA IF NOT EXISTS " + schema).Error; err != nil {
		_ = pg.Close(admin)
		return nil, err
	}
	_ = pg.Close(admin)

	isolated, err := pg.Open(withSearchPath(dsn, schema))
	if err != nil {
		return nil, err
	}
	return isolated, nil
}

func withSearchPath(dsn, schema string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	return u.String()
}

func startContainer(ctx context.Context) (dsn string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("testcontainers panic: %v", rec)
		}
	}()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("product"),
		postgres.WithUsername("product"),
		postgres.WithPassword("product"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return "", err
	}
	return container.ConnectionString(ctx, "sslmode=disable")
}
