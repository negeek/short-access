package utils
import(
	"fmt"
	"github.com/negeek/short-access/db"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func Setup() {
	appEnv:=os.Getenv("APP_ENV")
	if appEnv=="dev"{
		err := godotenv.Load(".env")
		if err != nil {
			// try this directory
			err = godotenv.Load("../../internal/env/.env")
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}
	}
	// DB connection
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	
	if err:= db.Connect(dbURL); err != nil {
		log.Fatal(err)
	}
}