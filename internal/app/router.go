package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rawen554/go-loyal/internal/middleware/auth"
	"github.com/rawen554/go-loyal/internal/middleware/compress"
	ginLogger "github.com/rawen554/go-loyal/internal/middleware/logger"
)

const (
	rootRoute    = "/"
	userAPIRoute = "/api/user"
)

func (a *App) SetupRouter() *gin.Engine {
	r := gin.New()
	ginLoggerMiddleware, err := ginLogger.Logger()
	if err != nil {
		log.Fatal(err)
	}
	r.Use(ginLoggerMiddleware)
	r.Use(compress.Compress())

	publicUserAPI := r.Group(userAPIRoute)
	{
		publicUserAPI.POST("register", a.Authz)
		publicUserAPI.POST("login", a.Authz)
	}

	protectedUserAPI := r.Group(userAPIRoute)
	protectedUserAPI.Use(auth.AuthMiddleware(a.Config.Seed))
	{
		protectedUserAPI.GET("withdrawals", a.GetWithdrawals)
		ordersAPI := protectedUserAPI.Group("orders")
		{
			ordersAPI.POST(rootRoute, a.PutOrder)
			ordersAPI.GET(rootRoute, a.GetOrders)
		}

		balanceAPI := protectedUserAPI.Group("balance")
		{
			balanceAPI.GET(rootRoute, a.GetBalance)
			balanceAPI.POST("withdraw", a.BalanceWithdraw)
		}
	}

	return r
}