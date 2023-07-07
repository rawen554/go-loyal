package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rawen554/go-loyal/internal/config"
	"github.com/rawen554/go-loyal/internal/middleware/auth"
	"github.com/rawen554/go-loyal/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	CreateUser(user *models.User) (int64, error)
	GetUser(login string) (*models.User, error)
	Ping() error
}

type App struct {
	Config *config.ServerConfig
	store  Store
}

func NewApp(config *config.ServerConfig, store Store) *App {
	return &App{
		Config: config,
		store:  store,
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

	user := models.User{
		Login:    userCreds.Login,
		Password: userCreds.Password,
	}
	if isRegister {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 7)
		if err != nil {
			log.Printf("cannot hash pass: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.Password = string(hash)

		rowsAffected, err := a.store.CreateUser(&user)
		if err != nil || rowsAffected == 0 {
			log.Printf("cannot operate user creds: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {
		u, err := a.store.GetUser(user.Login)
		if err != nil {
			log.Printf("cannot operate user creds: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		user.ID = u.ID
	}

	jwt, err := auth.BuildJWTString(strconv.FormatUint(user.ID, 10), a.Config.Seed)
	if err != nil {
		log.Printf("cannot build jwt string: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.SetCookie(auth.CookieName, jwt, 3600*24*30, "", "", false, true)
	res.WriteHeader(http.StatusOK)
}

func (a *App) PutOrder(c *gin.Context) {
	userID := c.GetString(auth.UserIDKey)

	if _, err := c.Writer.Write([]byte(userID)); err != nil {
		log.Printf("error writing bytes: %v", err)
		c.Writer.WriteHeader(http.StatusNotImplemented)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (a *App) GetOrders(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotImplemented)
}

func (a *App) GetWithdrawals(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotImplemented)
}

func (a *App) GetBalance(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotImplemented)
}

func (a *App) BalanceWithdraw(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotImplemented)
}

func (a *App) Ping(c *gin.Context) {
	if err := a.store.Ping(); err != nil {
		log.Printf("Error opening connection to DB: %v", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
