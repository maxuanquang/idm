package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maxuanquang/idm/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database interface {
	AddError(err error) error
	Assign(attrs ...interface{}) (tx *gorm.DB)
	Association(column string) *gorm.Association
	Attrs(attrs ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Commit() *gorm.DB
	Connection(fc func(tx *gorm.DB) error) (err error)
	Count(count *int64) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	CreateInBatches(value interface{}, batchSize int) (tx *gorm.DB)
	DB() (*sql.DB, error)
	Debug() (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Distinct(args ...interface{}) (tx *gorm.DB)
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FirstOrInit(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Get(key string) (interface{}, bool)
	Group(name string) (tx *gorm.DB)
	Having(query interface{}, args ...interface{}) (tx *gorm.DB)
	InnerJoins(query string, args ...interface{}) (tx *gorm.DB)
	InstanceGet(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) *gorm.DB
	Joins(query string, args ...interface{}) (tx *gorm.DB)
	Last(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Limit(limit int) (tx *gorm.DB)
	Migrator() gorm.Migrator
	Model(value interface{}) (tx *gorm.DB)
	Not(query interface{}, args ...interface{}) (tx *gorm.DB)
	Offset(offset int) (tx *gorm.DB)
	Omit(columns ...string) (tx *gorm.DB)
	Or(query interface{}, args ...interface{}) (tx *gorm.DB)
	Order(value interface{}) (tx *gorm.DB)
	Pluck(column string, dest interface{}) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
	Raw(sql string, values ...interface{}) (tx *gorm.DB)
	Rollback() *gorm.DB
	RollbackTo(name string) *gorm.DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Save(value interface{}) (tx *gorm.DB)
	SavePoint(name string) *gorm.DB
	Scan(dest interface{}) (tx *gorm.DB)
	ScanRows(rows *sql.Rows, dest interface{}) error
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) (tx *gorm.DB)
	Select(query interface{}, args ...interface{}) (tx *gorm.DB)
	Session(config *gorm.Session) *gorm.DB
	Set(key string, value interface{}) *gorm.DB
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Table(name string, args ...interface{}) (tx *gorm.DB)
	Take(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	ToSQL(queryFn func(tx *gorm.DB) *gorm.DB) string
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Unscoped() (tx *gorm.DB)
	Update(column string, value interface{}) (tx *gorm.DB)
	UpdateColumn(column string, value interface{}) (tx *gorm.DB)
	UpdateColumns(values interface{}) (tx *gorm.DB)
	Updates(values interface{}) (tx *gorm.DB)
	Use(plugin gorm.Plugin) error
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	WithContext(ctx context.Context) *gorm.DB
}

func InitializeDB(dbConfig configs.Database) (Database, func(), error) {
	dbMigrator, err := NewMigrator(dbConfig)
	if err != nil {
		return nil, nil, err
	}
	err = dbMigrator.Up(context.Background())
	if err != nil {
		return nil, nil, err
	}

	db, cleanup, err := InitializeGORMDB(dbConfig)
	if err != nil {
		return nil, nil, err
	}

	return db, cleanup, nil
}

func InitializeGORMDB(dbConfig configs.Database) (*gorm.DB, func(), error) {
	// Create data source name (DSN) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	// Open GORM database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup, nil
}
