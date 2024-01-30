package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

const (
	AccessToken    = "6887425159:AAERVk6Rs4W90PfczuHLN68EjG8OM-1Btbs" // Replace with your actual access token
	TelegramApiUrl = "https://api.telegram.org"
)

type Message struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date  int    `json:"date"`
		Text  string `json:"text"`
		Audio struct {
			FileID string `json:"file_id"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
		} `json:"document"`
		Video struct {
			FileID string `json:"file_id"`
		} `json:"video"`
		Photo []struct {
			FileID string `json:"file_id"`
		} `json:"photo"`
	} `json:"message"`
}

func webhook(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		log.Printf("failed to unmarshal body: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fmt.Println("message: ", message)

	if message.Message.Text != "" {
		err := sendTextMessage(message.Message.Chat.ID, message.Message.Text)
		if err != nil {
			handleSendError(c, err)
		}
	} else if len(message.Message.Photo) > 0 {
		err := sendPhotoMessage(message.Message.Chat.ID, message.Message.Photo[len(message.Message.Photo)-1].FileID)
		if err != nil {
			handleSendError(c, err)
		}
	} else if message.Message.Audio.FileID != "" {
		err := sendAudioMessage(message.Message.Chat.ID, message.Message.Audio.FileID)
		if err != nil {
			handleSendError(c, err)
		}
	} else if message.Message.Document.FileID != "" {
		err := sendDocumentMessage(message.Message.Chat.ID, message.Message.Document.FileID)
		if err != nil {
			handleSendError(c, err)
		}
	} else if message.Message.Video.FileID != "" {
		err := sendVideoMessage(message.Message.Chat.ID, message.Message.Video.FileID)
		if err != nil {
			handleSendError(c, err)
		}
	}

	c.Status(http.StatusOK)
}

func sendTextMessage(chatID int, message string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage?chat_id=%d&text=%s", TelegramApiUrl, AccessToken, chatID, message)
	return sendRequest(url)
}

func sendPhotoMessage(chatID int, imageURL string) error {
	url := fmt.Sprintf("%s/bot%s/sendPhoto?chat_id=%d&photo=%s", TelegramApiUrl, AccessToken, chatID, imageURL)
	return sendRequest(url)
}

func sendAudioMessage(chatID int, audioURL string) error {
	url := fmt.Sprintf("%s/bot%s/sendAudio?chat_id=%d&audio=%s", TelegramApiUrl, AccessToken, chatID, audioURL)
	return sendRequest(url)
}

func sendVideoMessage(chatID int, videoURL string) error {
	url := fmt.Sprintf("%s/bot%s/sendVideo?chat_id=%d&video=%s", TelegramApiUrl, AccessToken, chatID, videoURL)
	return sendRequest(url)
}

func sendDocumentMessage(chatID int, fileURL string) error {
	url := fmt.Sprintf("%s/bot%s/sendDocument?chat_id=%d&document=%s", TelegramApiUrl, AccessToken, chatID, fileURL)
	return sendRequest(url)
}

func sendRequest(reqUrl string) error {
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to wrap request: %w", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(res.Body)

	log.Printf("message sent successfully?\n%#v", res)
	return nil
}

func handleSendError(c *gin.Context, err error) {
	log.Printf("failed to send message: %v", err)
	c.AbortWithError(http.StatusInternalServerError, err)
}

func main() {
	router := gin.Default()
	router.POST("/", webhook)

	log.Printf("http server listening at localhost:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
