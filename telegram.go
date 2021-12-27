package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendTelegramMessage(apiKey, chatId, text string) error {
	data := url.Values{}
	data.Set("chat_id", chatId)
	data.Set("text", text)

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiKey), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("telegram api return error status code: %d", resp.StatusCode)
	}

	return nil
}
