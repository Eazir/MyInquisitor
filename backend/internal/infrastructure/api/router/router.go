package router

import (
	"github.com/gin-gonic/gin"
	"github.com/myinquisitor/backend/internal/infrastructure/api/handler"
	"github.com/myinquisitor/backend/internal/infrastructure/api/middleware"
)

func Setup(
	r *gin.Engine,
	authH *handler.AuthHandler,
	profileH *handler.ProfileHandler,
	debtH *handler.DebtHandler,
	expenseH *handler.ExpenseHandler,
	accH *handler.AccountingHandler,
	adminH *handler.AdminHandler,
	authMW *middleware.AuthMiddleware,
	adminMW *middleware.AdminMiddleware,
) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/refresh", authH.Refresh)
		}

		protected := v1.Group("")
		protected.Use(authMW.Authenticate())
		{
			debts := protected.Group("/debts")
			{
				debts.POST("", debtH.Create)
				debts.POST("/", debtH.Create)
				debts.GET("", debtH.List)
				debts.GET("/", debtH.List)
				debts.GET("/:id", debtH.GetByID)
				debts.PUT("/:id", debtH.Update)
				debts.DELETE("/:id", debtH.Delete)
				debts.GET("/:id/monthly", debtH.GetMonthlyStatus)
				debts.PUT("/:id/monthly/:year/:month/pay", debtH.MarkAsPaid)
			}

			expenses := protected.Group("/expenses")
			{
				expenses.POST("", expenseH.Create)
				expenses.POST("/", expenseH.Create)
				expenses.GET("", expenseH.List)
				expenses.GET("/", expenseH.List)
				expenses.GET("/:id", expenseH.GetByID)
				expenses.PUT("/:id", expenseH.Update)
				expenses.DELETE("/:id", expenseH.Delete)
				expenses.PUT("/:id/monthly/:year/:month/toggle", expenseH.TogglePaid)
			}

			accounting := protected.Group("/accounting")
			{
				accounting.POST("/transactions", accH.RecordTransaction)
				accounting.GET("/transactions", accH.ListTransactions)
				accounting.GET("/balance/:year/:month", accH.MonthlyBalance)
				accounting.GET("/cash-flow", accH.CashFlow)
				accounting.GET("/projections", accH.Projections)
			}

			categories := protected.Group("/categories")
			{
				categories.POST("", accH.CreateCategory)
				categories.POST("/", accH.CreateCategory)
				categories.GET("", accH.ListCategories)
				categories.GET("/", accH.ListCategories)
				categories.DELETE("/:id", accH.DeleteCategory)
			}

			profile := protected.Group("/profile")
			{
				profile.PUT("", profileH.UpdateProfile)
				profile.PUT("/", profileH.UpdateProfile)
				profile.PUT("/password", profileH.ChangePassword)
			}

			admin := protected.Group("/admin")
			admin.Use(adminMW.RequireSuperAdmin())
			{
				admin.GET("/users", adminH.ListUsers)
				admin.POST("/users", adminH.CreateUser)
				admin.PUT("/users/:id", adminH.UpdateUser)
				admin.PUT("/users/:id/activate/:active", adminH.SetActive)
				admin.POST("/invite", adminH.GenerateInvite)
			}
		}
	}
}
