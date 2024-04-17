package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()
}

func main() {
	initConfig()

	log.Println("Service started...")
	go scheduleDailyTask(viper.GetInt("schedule.hour"), viper.GetInt("schedule.minute"), viper.GetInt("schedule.second"), sendDutyReminder)
	select {}
}

func scheduleDailyTask(hour, min, sec int, task func()) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(next.Sub(now))
		task()
	}
}

func sendDutyReminder() {
	weekday := time.Now().Weekday()
	curweekday := "dutySchedule." + weekday.String()
	dutyPerson := viper.GetString(curweekday)
	message := dutyPerson + ", " + viper.GetString("messages.duty")

	hitokoto, err := getHitokoto()
	if err != nil {
		log.Println("Error getting hitokoto:", err)
	} else {
		message += "\n\n今日一言: " + hitokoto
	}

	isRoom := viper.GetBool("webhook.isRoom")

	payload := map[string]interface{}{
		"to":     viper.GetString("webhook.groupName"),
		"isRoom": isRoom,
		"data":   map[string]string{"content": message},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding payload:", err)
		return
	}

	resp, err := http.Post(viper.GetString("webhook.url"), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error sending duty reminder:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Duty reminder sent successfully: %s", string(body))
}

func getHitokoto() (string, error) {
	categoriesPath := filepath.Join("sentences", "categories.json")
	var categories []struct {
		Key  string `json:"key"`
		Path string `json:"path"`
	}

	hitokotobytes, err := os.ReadFile(categoriesPath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(hitokotobytes, &categories)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	selectedCategory := categories[rand.Intn(len(categories))]

	sentencesPath := filepath.Join("sentences", selectedCategory.Path)
	var sentences []struct {
		Hitokoto string `json:"hitokoto"`
		From     string `json:"from"`
	}

	hitokotobytes, err = os.ReadFile(sentencesPath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(hitokotobytes, &sentences)
	if err != nil {
		return "", err
	}

	selectedSentence := sentences[rand.Intn(len(sentences))]
	return fmt.Sprintf("%s ——《%s》", selectedSentence.Hitokoto, selectedSentence.From), nil
}
