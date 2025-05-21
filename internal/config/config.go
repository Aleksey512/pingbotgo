package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Settings struct {
	SQLitePath string            `validate:"required,omitempty"`
	TgBotToken string            `validate:"required"`
	Servers    map[string]string `validate:"required,min=1,dive,keys,required,endkeys,required"`
}

func (s *Settings) RedisURL() string {
	return ""
}

func NewSettings() (*Settings, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}

	sqlitePath := os.Getenv("SQLLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "sqlite3.db"
	}
	tgBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	serversEnv := os.Getenv("SERVERS")

	var servers map[string]string
	if err := json.Unmarshal([]byte(serversEnv), &servers); err != nil {
		return nil, fmt.Errorf("invalid SERVERS JSON: %v", err)
	}

	cfg := &Settings{
		SQLitePath: sqlitePath,
		TgBotToken: tgBotToken,
		Servers:    servers,
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validation error: %v", err)
	}

	return cfg, nil
}
