package store

import (
	"context"
	"errors"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/rawen554/go-loyal/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBStore struct {
	conn *gorm.DB
}

var ErrDBInsertConflict = errors.New("conflict insert into table, returned stored value")
var ErrURLDeleted = errors.New("url is deleted")

func NewPostgresStore(ctx context.Context, dsn string, logLevel logger.LogLevel) (*DBStore, error) {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn.Logger = logger.Default.LogMode(logLevel)
	if err := conn.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	log.Println("successfully connected to the database")

	return &DBStore{conn: conn}, nil
}

func (db *DBStore) CreateUser(user *models.User) (int64, error) {
	result := db.conn.Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (db *DBStore) GetUser(login string) (*models.User, error) {
	var user models.User
	result := db.conn.Where("login = ?", login).First(&user)

	return &user, result.Error
}

func (db *DBStore) Ping() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func (db *DBStore) Close() {
	sqlDB, err := db.conn.DB()
	if err != nil {
		log.Printf("gorm cant get sql.DB interface: %v", err)
	}

	sqlDB.Close()
}
