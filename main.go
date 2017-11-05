package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"encoding/csv"
	"bufio"
	"io"
)

type Person struct {
	lastname string
	info     string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// считываем в один стринг весь файл и сплитим на слайсы
	dataCSV, err := os.Open("data.csv")
	check(err)
	reader := csv.NewReader(bufio.NewReader(dataCSV))
	var activemembers []Person
	//зчитуємо дані в слайс структур
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		activemembers = append(activemembers, Person{
			lastname: line[0],
			info:     line[1],
		})
	}
	// заполняем мап
	//var mapmap map[string]*Person
	mapmap := make(map[string]*Person)
	for i := 0; i < len(activemembers); i++ {
		mapmap[activemembers[i].lastname] = &activemembers[i]
	}
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("442632858:AAGT6aDU-axkUJIyQ1M6dmAwNustMGfcPEA")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, _ := bot.GetUpdatesChan(ucfg)
	// читаем обновления из канала
	for {
		for update := range updates {
			// Пользователь, который написал боту
			if update.Message == nil {
				continue
			}
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			// в зависимости от команды - выводим имя из мапки или всю инфу
			switch update.Message.Command() {
			case "get":
				reply := mapmap[update.Message.CommandArguments()].info
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			case "all":
				for i := 0; i < len(activemembers); i++ {
					bot.Send(tgbotapi.NewMessage(ChatID, "Name: "+activemembers[i].lastname+"; Phone: "+activemembers[i].info))
				}
			}

			// Ответим пользователю его же сообщением
			//reply := Text
			// Созадаем сообщение
			//msg := tgbotapi.NewMessage(ChatID, reply)
			// и отправляем его

		}

	}
}
