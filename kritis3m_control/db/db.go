package db

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/philslol/kritis3m_scale/kritis3m_control/db/utils"
	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"github.com/philslol/kritis3m_scale/kritis3m_control/util"
)

// var errDatabaseNotSupported = errors.New("database type not supported")

type KV struct {
	Key   string
	Value string
}

type KSDatabase struct {
	DB *gorm.DB
}

// TODO(kradalby): assemble this struct from toptions or something typed
// rather than arguments.
func NewKritis3mScaleDatabase(
	cfg types.DatabaseConfig,
	db_logger zerolog.Logger,
) (*KSDatabase, error) {

	dbConn, err := openDB(cfg, db_logger)
	if err != nil {
		return nil, err
	}

	dbConn.Logger.LogMode(5) // Log detailed GORM output

	dbConn.AutoMigrate(&types.DBIdentity{})
	dbConn.AutoMigrate(&types.DBNode{},
		&types.DBNodes{},
		&types.DBApplication{},
		&types.DBIdentity{},
		&types.DBNode{},
		&types.HardwareConfig{},
		&types.DBWhitelist{},
		&types.DBNodeConfig{},
		&types.DBTrustedClients{},
		&types.SelectedConfiguration{},
		&types.DBAslEndpointConfig{},
	)

	db := KSDatabase{
		DB: dbConn,
	}

	return &db, err
}

func openDB(cfg types.DatabaseConfig, db_logger zerolog.Logger) (*gorm.DB, error) {
	// TODO(kradalby): Integrate this with zerolog
	var log_level int
	switch cfg.LogLevel {
	case "info":
		log_level = int(logger.Info)
	case "warn":
		log_level = int(logger.Warn)
	case "error":
		log_level = int(logger.Error)
	default:
		log.Fatal().Msg("wrong debug level for database logger")
	}
	log.Info().Msgf("log level %d", log_level)
	gormLogger := utils.NewGormZerologger(db_logger, logger.LogLevel(logger.Error), 200*time.Millisecond)

	dir := filepath.Dir(cfg.Sqlite.Path)
	err := util.EnsureDir(dir)
	if err != nil {
		return nil, fmt.Errorf("creating directory for sqlite: %w", err)
	}

	log.Info().
		Str("database", types.DatabaseSqlite).
		Str("path", cfg.Sqlite.Path).
		Msg("Opening database")
	//NOTE:
	// - journal mode: changes are writen to a write-ahead-log and are regularily committed
	//	- even during one open write transaction
	//- synchronous 1(=Normal) engine will sync at the most critical moments, but less often than in FULL(2), default is full
	//	- WAL + NORMAL is always consistent
	//- REFERENCE:
	//	- [Blog Post]: https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
	//	- [SQL Post]: https://www.sqlite.org/pragma.html#pragma_synchronous

	db, err := gorm.Open(
		sqlite.Open(cfg.Sqlite.Path+"?synchronous=1&_journal_mode=WAL"),
		&gorm.Config{
			Logger: gormLogger,
		},
	)
	if err != nil {
		log.Err(err).Msg("cant open database")
	}
	db.Exec("PRAGMA mmap_size=268435456")

	// // The pure Go SQLite library does not handle locking in
	// // the same way as the C based one and we cant use the gorm
	// // connection pool as of 2022/02/23.
	sqlDB, _ := db.DB() //returns the db
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1) // single instance can write to db
	sqlDB.SetConnMaxIdleTime(time.Hour)

	return db, err

}

func (ksdb *KSDatabase) PingDB(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	sqlDB, err := ksdb.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func (ksdb *KSDatabase) Close() error {
	db, err := ksdb.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (ksdb *KSDatabase) Read(fn func(rx *gorm.DB) error) error {
	rx := ksdb.DB.Begin()
	defer rx.Rollback()
	return fn(rx)
}

func Read[T any](db *gorm.DB, fn func(rx *gorm.DB) (T, error)) (T, error) {
	rx := db.Begin()
	defer rx.Rollback()
	ret, err := fn(rx)
	if err != nil {
		var no T
		return no, err
	}
	return ret, nil
}

func (ksdb *KSDatabase) Write(fn func(tx *gorm.DB) error) error {
	tx := ksdb.DB.Begin()
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit().Error
}

func Write[T any](db *gorm.DB, fn func(tx *gorm.DB) (T, error)) (T, error) {
	tx := db.Begin()
	defer tx.Rollback()
	ret, err := fn(tx)
	if err != nil {
		var no T
		return no, err
	}
	return ret, tx.Commit().Error
}
