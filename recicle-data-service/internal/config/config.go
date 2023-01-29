package config

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	Env = ""

	AWSLogin  = ""
	AWSSecret = ""
	AWSRegion = ""

	DBConnURL = ""
)

func init() {
	godotenv.Load()

	Env = os.Getenv("ENV")

	AWSLogin = os.Getenv("AWS_LOGIN")
	AWSSecret = os.Getenv("AWS_SECRET")
	AWSRegion = os.Getenv("AWS_REGION")

	DBConnURL = os.Getenv("DB_CONN_URL")
}
