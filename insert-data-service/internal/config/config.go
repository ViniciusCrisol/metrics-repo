package config

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	Env = ""

	SQS       = ""
	AWSLogin  = ""
	AWSSecret = ""
	AWSRegion = ""

	DBConnURL = ""
)

func init() {
	godotenv.Load()

	Env = os.Getenv("ENV")

	SQS = os.Getenv("SQS")
	AWSLogin = os.Getenv("AWS_LOGIN")
	AWSSecret = os.Getenv("AWS_SECRET")
	AWSRegion = os.Getenv("AWS_REGION")

	DBConnURL = os.Getenv("DB_CONN_URL")
}
