package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rawen554/go-loyal/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBStore struct {
	conn *gorm.DB
}

type Store interface {
	CreateUser(user *models.User) (int64, error)
	GetUser(u *models.User) (*models.User, error)
	PutOrder(number string, userID uint64) error
	GetUserOrders(userID uint64) ([]models.Order, error)
	GetUserBalance(userID uint64) (*models.UserBalanceShema, error)
	CreateWithdraw(userID uint64, w models.BalanceWithdrawShema) error
	GetWithdrawals(userID uint64) ([]models.Withdraw, error)
	Ping() error
	Close()
}

var ErrDBInsertConflict = errors.New("conflict insert into table, returned stored value")
var ErrURLDeleted = errors.New("url is deleted")
var ErrLoginNotFound = errors.New("login not found")
var ErrDuplicateLogin = errors.New("login already registered")

func NewPostgresStore(ctx context.Context, dsn string, logLevel logger.LogLevel) (Store, error) {
	conn, err := ConnectLoop(dsn, 5*time.Second, time.Minute)
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

func ConnectLoop(dsn string, tick time.Duration, timeout time.Duration) (*gorm.DB, error) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", time.Minute)
		case <-ticker.C:
			conn, err := gorm.Open(postgres.New(postgres.Config{
				DSN: dsn,
			}), &gorm.Config{})
			if err != nil {
				log.Printf("error connecting to db: %v", err)
			} else {
				return conn, err
			}
		}
	}
}

func (db *DBStore) CreateUser(user *models.User) (int64, error) {
	result := db.conn.Create(user)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, ErrDuplicateLogin
			}
		}

		log.Printf("error saving user to db: %v", result.Error)
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

	if order.UserID != userID {
		return models.ErrOrderHasBeenProcessedByAnotherUser
	}

	if order.UserID == userID && order.Number == number && result.RowsAffected == 0 {
		return models.ErrOrderHasBeenProcessedByUser
	}

	return nil
}

func (db *DBStore) GetUserOrders(userID uint64) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	result := db.conn.Order("uploaded_at asc").Where(&models.Order{UserID: userID}).Find(&orders)

	if err := result.Error; err != nil {
		return nil, fmt.Errorf("error getting all user orders: %w", err)
	}

	if len(orders) == 0 {
		return nil, models.ErrUserHasNoItems
	}

	return orders, nil
}

func (db *DBStore) CreateWithdraw(userID uint64, w models.BalanceWithdrawShema) error {
	result := db.conn.Create(&models.Withdraw{OrderNum: w.Order, Sum: w.Sum, UserID: userID})
	return result.Error
}

func (db *DBStore) GetWithdrawals(userID uint64) ([]models.Withdraw, error) {
	withdrawals := make([]models.Withdraw, 0)
	result := db.conn.Order("processed_at asc").Where(&models.Withdraw{UserID: userID}).Find(&withdrawals)

	if err := result.Error; err != nil {
		return nil, fmt.Errorf("error getting all user withdrawals: %w", err)
	}

	if len(withdrawals) == 0 {
		return nil, models.ErrUserHasNoItems
	}

	return withdrawals, nil
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
