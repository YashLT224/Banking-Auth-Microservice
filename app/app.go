package app

import (
	"BankingAuth/domain"
	"BankingAuth/logger"
	"BankingAuth/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func createDbClient(dbUrl string) *sqlx.DB {
	DBID := os.Getenv("DBID")
	DBPSWD := os.Getenv("DBPSWD")
	DB := os.Getenv("DB")
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", DBID, DBPSWD, DB))
	if err != nil {

		panic(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 20)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	//defer db.Close()
	return db
}
func Start() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("godotenv error: %s", err)
		return
	}

	router := mux.NewRouter()
	dbClient := createDbClient(os.Getenv("CLEARDB_DATABASE_URL"))
	authRepository := domain.NewAuthRepository(dbClient)
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions())}

	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", ah.NotImplementedHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	logger.Info(fmt.Sprintf("Starting OAuth server on %s:%s ...", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}
