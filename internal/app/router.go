package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rawen554/go-loyal/internal/middleware/auth"
	"github.com/rawen554/go-loyal/internal/middleware/compress"
	ginLogger "github.com/rawen554/go-loyal/internal/middleware/logger"
)

const (
	emptyRoute   = ""
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

	r.POST("/api/user/register", a.Authz)
	r.POST("/api/user/login", a.Authz)

	protectedUserAPI := r.Group(userAPIRoute)
	protectedUserAPI.Use(auth.AuthMiddleware(a.Config.Seed))
	{
		protectedUserAPI.GET("withdrawals", a.GetWithdrawals)
		ordersAPI := protectedUserAPI.Group("orders")
		{
			ordersAPI.POST(emptyRoute, a.PutOrder)
			ordersAPI.GET(emptyRoute, a.GetOrders)
		}

		balanceAPI := protectedUserAPI.Group("balance")
		{
			balanceAPI.GET(emptyRoute, a.GetBalance)
			balanceAPI.POST("withdraw", a.BalanceWithdraw)
		}
	}

	return r
}
