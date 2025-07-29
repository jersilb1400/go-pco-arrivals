package database

import (
	"database/sql"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/models"
)

// DatabaseType represents the type of database being used
type DatabaseType string

const (
	SQLiteDB  DatabaseType = "sqlite"
	MongoDBDB DatabaseType = "mongodb"
)

// Database interface for both SQLite and MongoDB
type Database interface {
	Connect() error
	Close() error
	Migrate() error
	GetType() DatabaseType
}

// SQLiteDatabase wraps GORM SQLite connection
type SQLiteDatabase struct {
	config config.DatabaseConfig
	db     *gorm.DB
}

// MongoDBDatabase wraps MongoDB connection
type MongoDBDatabase struct {
	mongodb *MongoDB
}

func Connect(cfg config.DatabaseConfig) (Database, error) {
	dbType := DatabaseType(os.Getenv("DATABASE_TYPE"))
	if dbType == "" {
		dbType = SQLiteDB // Default to SQLite
	}

	switch dbType {
	case MongoDBDB:
		return connectMongoDB()
	case SQLiteDB:
		fallthrough
	default:
		return connectSQLite(cfg)
	}
}

func connectSQLite(cfg config.DatabaseConfig) (*SQLiteDatabase, error) {
	db, err := gorm.Open(sqlite.Open(cfg.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &SQLiteDatabase{
		config: cfg,
		db:     db,
	}, nil
}

func connectMongoDB() (*MongoDBDatabase, error) {
	mongodb, err := NewMongoDB()
	if err != nil {
		return nil, err
	}

	return &MongoDBDatabase{
		mongodb: mongodb,
	}, nil
}

// SQLiteDatabase methods
func (s *SQLiteDatabase) Connect() error {
	// Already connected in constructor
	return nil
}

func (s *SQLiteDatabase) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *SQLiteDatabase) Migrate() error {
	return s.db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Event{},
		&models.Notification{},
		&models.CheckIn{},
		&models.Location{},
		&models.BillboardState{},
		&models.SecurityCode{},
	)
}

func (s *SQLiteDatabase) GetType() DatabaseType {
	return SQLiteDB
}

// GetGormDB returns the underlying GORM DB for SQLite
func (s *SQLiteDatabase) GetGormDB() *gorm.DB {
	return s.db
}

// MongoDBDatabase methods
func (m *MongoDBDatabase) Connect() error {
	return m.mongodb.Migrate()
}

func (m *MongoDBDatabase) Close() error {
	return m.mongodb.Close()
}

func (m *MongoDBDatabase) Migrate() error {
	return m.mongodb.Migrate()
}

func (m *MongoDBDatabase) GetType() DatabaseType {
	return MongoDBDB
}

// GetMongoDB returns the underlying MongoDB connection
func (m *MongoDBDatabase) GetMongoDB() *MongoDB {
	return m.mongodb
}

// Legacy functions for backward compatibility
func ConnectLegacy(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConfigureConnectionPool(db *gorm.DB, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func MigrateLegacy(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Event{},
		&models.Notification{},
		&models.CheckIn{},
		&models.Location{},
		&models.BillboardState{},
		&models.SecurityCode{},
	)
}

func GetDBStats(db *gorm.DB) (*sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return &stats, nil
}
