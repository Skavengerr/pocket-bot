package main

import (
	"fmt"
	"log"

	"golang-pocket/pkg/config"
	"golang-pocket/pkg/repository"
	"golang-pocket/pkg/repository/boltdb"
	"golang-pocket/pkg/server"
	"golang-pocket/pkg/telegram"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		return
	}

	fmt.Printf("441114444\n")
	fmt.Println(cfg)
	bot, err := tgbotapi.NewBotAPI("5495765071:AAEN_SE6C48PzqsOAIY8BBAQ0GUw9UTPoZg")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("102646-93045e4d8e7b84f9a733a75")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/", &cfg.Messages)

	authServer := server.NewAuthServer(pocketClient, tokenRepository, "https://t.me/golang_r_pocket_bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open("db_file", 0600, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return db, nil
}
