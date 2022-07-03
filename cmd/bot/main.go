package main

import (
	"golang-pocket/pkg/config"
	"golang-pocket/pkg/server"
	"golang-pocket/pkg/storage"
	"golang-pocket/pkg/storage/boltdb"
	"golang-pocket/pkg/telegram"
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	botApi, err := tgbotapi.NewBotAPI("5495765071:AAEN_SE6C48PzqsOAIY8BBAQ0GUw9UTPoZg")
	if err != nil {
		log.Fatal(err)
	}
	botApi.Debug = true

	pocketClient, err := pocket.NewClient("102646-93045e4d8e7b84f9a733a75")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initBolt()
	if err != nil {
		log.Fatal(err)
	}
	storage := boltdb.NewTokenStorage(db)

	bot := telegram.NewBot(botApi, pocketClient, "http://localhost/", storage, cfg.Messages)

	redirectServer := server.NewAuthServer(cfg.BotURL, storage, pocketClient)

	go func() {
		if err := redirectServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initBolt() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(storage.RequestTokens))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
