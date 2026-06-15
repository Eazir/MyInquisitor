package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	appAcc "github.com/myinquisitor/backend/internal/application/accounting"
	appAdmin "github.com/myinquisitor/backend/internal/application/admin"
	appAuth "github.com/myinquisitor/backend/internal/application/auth"
	appDebt "github.com/myinquisitor/backend/internal/application/debt"
	appExpense "github.com/myinquisitor/backend/internal/application/expense"
	appProfile "github.com/myinquisitor/backend/internal/application/profile"
	"github.com/myinquisitor/backend/internal/infrastructure/api/handler"
	"github.com/myinquisitor/backend/internal/infrastructure/api/middleware"
	"github.com/myinquisitor/backend/internal/infrastructure/api/router"
	"github.com/myinquisitor/backend/internal/infrastructure/auth"
	"github.com/myinquisitor/backend/internal/infrastructure/config"
	"github.com/myinquisitor/backend/internal/infrastructure/persistence"
	"github.com/myinquisitor/backend/internal/infrastructure/security"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := persistence.NewPostgresDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connected to database")

	migrationsDir := "internal/infrastructure/persistence/migrations"
	if env := os.Getenv("MIGRATIONS_DIR"); env != "" {
		migrationsDir = env
	}
	if err := persistence.RunMigrations(ctx, db, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied successfully")

	encryptSvc, err := security.NewEncryptionService(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("failed to initialize encryption: %v", err)
	}

	pwdSvc := security.NewPasswordService()
	jwtSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration, cfg.RefreshExpiration)

	userRepo := persistence.NewUserRepository(db, encryptSvc)
	debtRepo := persistence.NewDebtRepository(db)
	debtMonthlyRepo := persistence.NewDebtMonthlyStatusRepository(db)
	expenseRepo := persistence.NewRecurringExpenseRepository(db)
	expenseMonthlyRepo := persistence.NewExpenseMonthlyStatusRepository(db)
	txRepo := persistence.NewTransactionRepository(db)
	summaryRepo := persistence.NewMonthlySummaryRepository(db)
	catRepo := persistence.NewCategoryRepository(db)
	inviteRepo := persistence.NewInviteTokenRepository(db)

	registerUC := appAuth.NewRegisterUseCase(userRepo, inviteRepo, pwdSvc, jwtSvc)
	loginUC := appAuth.NewLoginUseCase(userRepo, pwdSvc, jwtSvc)
	refreshUC := appAuth.NewRefreshUseCase(userRepo, jwtSvc)

	debtCreateUC := appDebt.NewCreateUseCase(debtRepo, debtMonthlyRepo)
	debtListUC := appDebt.NewListUseCase(debtRepo)
	debtGetByIDUC := appDebt.NewGetByIDUseCase(debtRepo)
	debtUpdateUC := appDebt.NewUpdateUseCase(debtRepo, debtMonthlyRepo)
	debtDeleteUC := appDebt.NewDeleteUseCase(debtRepo)
	debtMarkPaidUC := appDebt.NewMarkPaidUseCase(debtMonthlyRepo, debtRepo, txRepo)
	debtMonthlyStatusUC := appDebt.NewGetMonthlyStatusUseCase(debtMonthlyRepo)

	expCreateUC := appExpense.NewCreateUseCase(expenseRepo)
	expListUC := appExpense.NewListUseCase(expenseRepo)
	expGetByIDUC := appExpense.NewGetByIDUseCase(expenseRepo)
	expUpdateUC := appExpense.NewUpdateUseCase(expenseRepo)
	expDeleteUC := appExpense.NewDeleteUseCase(expenseRepo)
	expTogglePaidUC := appExpense.NewTogglePaidUseCase(expenseMonthlyRepo, expenseRepo, txRepo)
	expMonthlyStatusUC := appExpense.NewGetMonthlyStatusUseCase(expenseMonthlyRepo)

	accRecordTxUC := appAcc.NewRecordTransactionUseCase(txRepo)
	accListTxUC := appAcc.NewListTransactionsUseCase(txRepo)
	accBalanceUC := appAcc.NewGetMonthlyBalanceUseCase(summaryRepo, txRepo)
	accCashFlowUC := appAcc.NewGetCashFlowUseCase(txRepo)
	accProjectionsUC := appAcc.NewGetProjectionsUseCase(txRepo, expenseRepo, debtMonthlyRepo)
	accCreateCatUC := appAcc.NewCreateCategoryUseCase(catRepo)
	accListCatUC := appAcc.NewListCategoriesUseCase(catRepo)
	accDeleteCatUC := appAcc.NewDeleteCategoryUseCase(catRepo)

	profileUpdateUC := appProfile.NewUpdateProfileUseCase(userRepo)
	profileChangePasswordUC := appProfile.NewChangePasswordUseCase(userRepo, pwdSvc)

	adminListUsersUC := appAdmin.NewListUsersUseCase(userRepo)
	adminCreateUserUC := appAdmin.NewCreateUserUseCase(userRepo, pwdSvc)
	adminUpdateUserUC := appAdmin.NewUpdateUserUseCase(userRepo, pwdSvc)
	adminDeactivateUserUC := appAdmin.NewDeactivateUserUseCase(userRepo)
	adminGenerateInviteUC := appAdmin.NewGenerateInviteUseCase(inviteRepo)
	adminListInvitesUC := appAdmin.NewListInvitesUseCase(inviteRepo)
	adminDeleteInviteUC := appAdmin.NewDeleteInviteUseCase(inviteRepo)

	authH := handler.NewAuthHandler(registerUC, loginUC, refreshUC)
	profileH := handler.NewProfileHandler(profileUpdateUC, profileChangePasswordUC)
	debtH := handler.NewDebtHandler(debtCreateUC, debtListUC, debtGetByIDUC, debtUpdateUC, debtDeleteUC, debtMarkPaidUC, debtMonthlyStatusUC)
	expenseH := handler.NewExpenseHandler(expCreateUC, expListUC, expGetByIDUC, expUpdateUC, expDeleteUC, expTogglePaidUC, expMonthlyStatusUC)
	accH := handler.NewAccountingHandler(accRecordTxUC, accListTxUC, accBalanceUC, accCashFlowUC, accProjectionsUC, accCreateCatUC, accListCatUC, accDeleteCatUC)
	adminH := handler.NewAdminHandler(adminListUsersUC, adminCreateUserUC, adminUpdateUserUC, adminDeactivateUserUC, adminGenerateInviteUC, adminListInvitesUC, adminDeleteInviteUC)

	authMW := middleware.NewAuthMiddleware(jwtSvc)
	adminMW := middleware.NewAdminMiddleware()
	corsMW := middleware.NewCORSMiddleware(cfg.AllowedOrigins)

	r := gin.Default()
	r.Use(corsMW.Handler())
	router.Setup(r, authH, profileH, debtH, expenseH, accH, adminH, authMW, adminMW)

	srv := make(chan os.Signal, 1)
	signal.Notify(srv, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("starting server on port %s", cfg.ServerPort)
		if err := r.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-srv
	log.Println("shutting down server...")
}
