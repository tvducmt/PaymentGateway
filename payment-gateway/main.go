package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"gitlab.com/rockship/payment-gateway/api"

	_ "github.com/lib/pq"
)

// STARTSERVER ...
const STARTSERVER = `Server's running!
        .---.  
      .'_:___".
      |__ --==| Beep Boop ...
      [  ]  :[| 
      |__| I=[| 
      / / ____| 
     |-/.____.' 
    /___\ /___\
>`

var (
	postgresURI string
	db          *sqlx.DB
	log         *logrus.Entry
	workDir     string
	staticPath  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	raven.SetDSN(os.Getenv("RAVEN_DSN"))
	postgresURI = os.Getenv("POSTGRES_URI")
	db, err = sqlx.Connect("postgres", postgresURI)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Database is connected.")
}

func createRouter(env api.Env) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RequestID,
		middleware.RedirectSlashes,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
		SentryLoggingMiddleware,
	)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	router.Use(cors.Handler)

	router.Get("/", api.HandleAPI(env.Index).Reg)
	router.Post("/login", api.HandleAPI(env.LoginProcess).Reg)
	router.Post("/register", api.HandleAPI(env.RegisterProcess).Reg)
	router.Post("/user/forgot-password", api.HandleAPI(env.ForgotPassword).Reg)
	router.Post("/user/reset-password", api.HandleAPI(env.ResetPassword).Reg)

	router.Put("/user/otp", api.HandleAPI(env.VerifyOtp).Reg)
	router.Group(func(r chi.Router) {
		r.Use(OauthAuthorizationMiddleware)

		r.Put("/user/info", api.HandleAPI(env.UpdateUserInfo).Reg)
		r.Put("/user/password", api.HandleAPI(env.UpdatePassword).Reg)

		r.Get("/user/balance", api.HandleAPI(env.GetBalance).Reg)
		r.Get("/user/info", api.HandleAPI(env.GetUserInfo).Reg)
		r.Get("/user/coupons", api.HandleAPI(env.GetCouponInfo).Reg)
		r.Get("/customer/cards", api.HandleAPI(env.GetAllCustomerStripe).Reg)
		r.Get("/transaction", api.HandleAPI(env.GetTransactionDetail).Reg)
		r.Post("/credit-card", api.HandleAPI(env.CreditCardProcess).Reg)
		r.Post("/bank-transfer", api.HandleAPI(env.BankTransferProcess).Reg)
		r.Post("/transactions", api.HandleAPI(env.HistoryTransaction).Reg)

		r.Post("/ethereum-transfer", api.HandleAPI(env.EthereumTransferProcess).Reg)
	})

	router.Route("/secret", func(r chi.Router) {
		r.Use(TrustDeviceTokenMiddleware)

		r.Put("/transactions/timeout", api.HandleAPI(env.UpdateTransactionTimeOut).Reg)
		r.Post("/receipt", api.HandleAPI(env.ReceiveReceipt).Reg)
	})
	return router
}

func main() {
	env := api.Env{
		DB:  db,
		Log: log,
	}
	router := createRouter(env)
	//models.GenarateCodes(db)

	fmt.Println(STARTSERVER)
	server := &http.Server{
		Addr:         ":3000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go EthereumMonitoring(env.DB)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
