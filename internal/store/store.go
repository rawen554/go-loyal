package store

import (
	"context"
	"errors"
	"fmt"
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
var ErrLoginNotFound = errors.New("login not found")

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
	if err := conn.AutoMigrate(&models.Order{}); err != nil {
		return nil, err
	}
	if err := conn.AutoMigrate(&models.Withdraw{}); err != nil {
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

func (db *DBStore) GetUser(u *models.User) (*models.User, error) {
	var user models.User
	result := db.conn.Where(u).First(&user)

	if result.RowsAffected == 0 {
		return nil, ErrLoginNotFound
	}

	return &user, result.Error
}

func (db *DBStore) GetUserBalance(userID uint64) (*models.UserBalanceShema, error) {
	var user models.User
	var userBalance models.UserBalanceShema
	result := db.conn.Model(&user).Where(&models.User{ID: userID}).Take(&userBalance)

	return &userBalance, result.Error
}

func (db *DBStore) PutOrder(number string, userID uint64) error {
	var order models.Order
	result := db.conn.Where(models.Order{Number: number}).Attrs(models.Order{UserID: userID, Status: models.NEW}).FirstOrCreate(&order)
	if err := result.Error; err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}

	if order.UserID == userID && order.Number == number && result.RowsAffected == 0 {
		return models.ErrOrderHasBeenProcessedByUser
	}

	return nil
}

func (db *DBStore) GetUserOrders(userID uint64) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	result := db.conn.Order("uploaded_at desc").Where(&models.Order{UserID: userID}).Find(&orders)

	if err := result.Error; err != nil {
		return nil, fmt.Errorf("error getting all user orders: %w", err)
	}

	if len(orders) == 0 {
		return nil, models.ErrUserHasNoOrders
	}

	return orders, nil
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
