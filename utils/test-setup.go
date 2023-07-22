package utils
import(
	"github.com/negeek/short-access/db"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func Setup() {
	// env
	err := godotenv.Load("../../../internal/env/.env")
	
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	// DB connection
	dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL not set")
    }
	if err= db.Connect(dbURL); err != nil {
		log.Fatal(err)
	}
}