package app

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rawen554/go-loyal/internal/config"
	"github.com/rawen554/go-loyal/internal/middleware/auth"
	"github.com/rawen554/go-loyal/internal/models"
	"github.com/rawen554/go-loyal/internal/store"
	"github.com/rawen554/go-loyal/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Config  *config.ServerConfig
	store   store.Store
	accrual Accrual
}

func NewApp(config *config.ServerConfig, store store.Store, accrual Accrual) *App {
	return &App{
		Config:  config,
		store:   store,
		accrual: accrual,
	}
}

func (a *App) Authz(c *gin.Context) {
	req := c.Request
	res := c.Writer
	isRegister := strings.Contains(c.Request.RequestURI, "register")

	userCreds := models.User{}
	if err := json.NewDecoder(req.Body).Decode(&userCreds); err != nil {
		log.Printf("body cannot be decoded: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userReq := models.User{
		Login:    userCreds.Login,
		Password: userCreds.Password,
	}
	if isRegister {
		hash, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), 7)
		if err != nil {
			log.Printf("cannot hash pass: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		userReq.Password = string(hash)

		if _, err = a.store.CreateUser(&userReq); err != nil {
			if errors.Is(err, store.ErrDuplicateLogin) {
				log.Printf("login already taken: %v", err)
				res.WriteHeader(http.StatusConflict)
				return
			} else {
				log.Printf("cannot operate user creds: %v", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	} else {
		u, err := a.store.GetUser(&models.User{Login: userReq.Login})
		if err != nil {
			if errors.Is(err, store.ErrLoginNotFound) {
				log.Printf("login not found: %v", err)
				res.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				log.Printf("cannot operate user creds: %v", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userReq.Password)); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		userReq.ID = u.ID
	}

	jwt, err := auth.BuildJWTString(userReq.ID, a.Config.Seed)
	if err != nil {
		log.Printf("cannot build jwt string: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.SetCookie(auth.CookieName, jwt, 3600*24*30, "", "", false, true)
	res.WriteHeader(http.StatusOK)
}

func (a *App) PutOrder(c *gin.Context) {
	userID := c.GetUint64(auth.UserIDKey)
	req := c.Request
	res := c.Writer

	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	number := string(body)

	if isValidLuhn := utils.IsValidLuhn(number); !isValidLuhn {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := a.store.PutOrder(number, userID); err != nil {
		if errors.Is(err, models.ErrOrderHasBeenProcessedByAnotherUser) {
			res.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, models.ErrOrderHasBeenProcessedByUser) {
			res.WriteHeader(http.StatusOK)
			return
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	go func() {
		info, err := a.accrual.GetOrderInfo(number)
		if err != nil {
			log.Printf("error interacting with accrual: %v", err)
			return
		}
		log.Printf("accrual OK: %v", info)

	}()

	res.WriteHeader(http.StatusAccepted)
}

func (a *App) GetOrders(c *gin.Context) {
	userID := c.GetUint64(auth.UserIDKey)
	res := c.Writer

	orders, err := a.store.GetUserOrders(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserHasNoItems) {
			res.WriteHeader(http.StatusNoContent)
			return
		} else {
			log.Printf("error getting user orders: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(orders); err != nil {
		log.Printf("Error writing response in JSON: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) GetWithdrawals(c *gin.Context) {
	userID := c.GetUint64(auth.UserIDKey)
	res := c.Writer

	withdrawals, err := a.store.GetWithdrawals(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserHasNoItems) {
			res.WriteHeader(http.StatusNoContent)
			return
		} else {
			log.Printf("unknown error: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(withdrawals); err != nil {
		log.Printf("Error writing response in JSON: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) GetBalance(c *gin.Context) {
	userID := c.GetUint64(auth.UserIDKey)
	res := c.Writer

	balance, err := a.store.GetUserBalance(userID)
	if err != nil {
		log.Printf("Error getting user balance: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(balance); err != nil {
		log.Printf("Error writing response in JSON: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) BalanceWithdraw(c *gin.Context) {
	userID := c.GetUint64(auth.UserIDKey)
	res := c.Writer
	req := c.Request

	var withdrawRequest models.BalanceWithdrawShema
	if err := json.NewDecoder(req.Body).Decode(&withdrawRequest); err != nil {
		log.Printf("Body cannot be decoded: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	u, err := a.store.GetUser(&models.User{ID: userID})
	if err != nil {
		log.Printf("error getting user by ID: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if u.Balance < float64(withdrawRequest.Sum) {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}

	if isValidLuhn := utils.IsValidLuhn(withdrawRequest.Order); !isValidLuhn {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := a.store.CreateWithdraw(userID, withdrawRequest); err != nil {
		log.Printf("cant save withdraw: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (a *App) Ping(c *gin.Context) {
	if err := a.store.Ping(); err != nil {
		log.Printf("Error opening connection to DB: %v", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
