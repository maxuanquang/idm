package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/maxuanquang/idm/internal/configs"
)

var (
	//go:embed migrations/mysql
	migrationDirectoryMySQL embed.FS
)

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

func NewMigrator(dbConfig configs.Database) (Migrator, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	sourceInstance, err := iofs.New(migrationDirectoryMySQL, "migrations/mysql")
	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceInstance,
		dbConfig.Database,
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &migrator{
		instance: m,
	}, nil
}

type migrator struct {
	instance *migrate.Migrate
}

// Down implements Migrator.
func (m *migrator) Down(ctx context.Context) error {
	err := m.instance.Down()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

// Up implements Migrator.
func (m *migrator) Up(ctx context.Context) error {
	err := m.instance.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
