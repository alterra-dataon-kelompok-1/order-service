package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config interface {
	Get(key string) string
}

type config struct{}

func (c *config) Get(key string) string {
	return os.Getenv(key)
}

func New(filenames ...string) Config {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Println("error loading environment variables:", err)
		panic(err)
	}
	return &config{}
}
