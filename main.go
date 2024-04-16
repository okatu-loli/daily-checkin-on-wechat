package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

const (
	webhookURL  = "http://192.168.0.71:3001/webhook/msg/v2?token=L9N1P6fOCyDJ"
	groupName   = "千石"
	dutyMessage = "是时候进行日常值班签到了！\n\n值班内容：\n1. 给猫铲屎\n2. 与猫互动"
)

var dutySchedule = map[time.Weekday]string{
	time.Monday:    "七海",
	time.Tuesday:   "夏夏",
	time.Wednesday: "千石",
	time.Thursday:  "七海",
	time.Friday:    "小杰",
	time.Saturday:  "小杰",
	time.Sunday:    "千石",
}

func main() {
	log.Println("Service started...")
	go scheduleDailyTask(18, 11, 0, sendDutyReminder) // 每天9:00触发
	select {}
}

func scheduleDailyTask(hour, min, sec int, task func()) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(next.Sub(now)) // 等待直到下一次触发时间
		task()                    // 执行任务
	}
}

func sendDutyReminder() {
	weekday := time.Now().Weekday()
	dutyPerson := dutySchedule[weekday] // 仍然确定当天负责的人，但是消息将发到群组
	message := dutyPerson + ", " + dutyMessage

	// 添加一言到消息末尾
	hitokoto, err := getHitokoto()
	if err != nil {
		log.Println("Error getting hitokoto:", err)
	} else {
		message += "\n\n今日一言: " + hitokoto
	}

	payload := map[string]interface{}{
		"to":     groupName,
		"isRoom": false, // 指示消息应发送到群组
		"data":   map[string]string{"content": message},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding payload:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error sending duty reminder to the group:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Duty reminder sent successfully to group %s: %s", groupName, string(body))
}

func getHitokoto() (string, error) {
	// 读取categories.json以获取所有句子文件的路径
	categoriesPath := filepath.Join("sentences", "categories.json")
	var categories []struct {
		Key  string `json:"key"`
		Path string `json:"path"`
	}

	hitokotobytes, err := ioutil.ReadFile(categoriesPath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(hitokotobytes, &categories)
	if err != nil {
		return "", err
	}

	// 随机选择一个类别
	rand.Seed(time.Now().UnixNano())
	selectedCategory := categories[rand.Intn(len(categories))]

	// 读取选中类别的句子文件
	sentencesPath := filepath.Join("sentences", selectedCategory.Path)
	var sentences []struct {
		Hitokoto string `json:"hitokoto"`
		From     string `json:"from"`
	}

	hitokotobytes, err = ioutil.ReadFile(sentencesPath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(hitokotobytes, &sentences)
	if err != nil {
		return "", err
	}

	// 随机选择一句话
	selectedSentence := sentences[rand.Intn(len(sentences))]

	// 返回选中的一言和来源
	return fmt.Sprintf("%s ——《%s》", selectedSentence.Hitokoto, selectedSentence.From), nil
}
