package telegram

import (
	"context"
	"fmt"

	"golang-pocket/pkg/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) initAuthProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthLink(message.Chat.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", b.messages.Responses)

	msgText := fmt.Sprintf(b.messages.Responses.Start, authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	fmt.Println("AS:JOFJOQWFHASFHIOP")
	fmt.Println(b.tokenRepository.Get(chatID, repository.AccessTokens))
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectLink(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)

	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectLink(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
