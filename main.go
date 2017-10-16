package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// считываем в один стринг весь файл и сплитим на слайсы
	dat, err := ioutil.ReadFile("C:\\Users\\Chalyi\\Desktop\\BEST\\IT_dept\\infobot\\numbers")
	check(err)
	premapinfo := strings.Fields(string(dat))

	// заполняем мап
	var mapmap map[string]string
	mapmap=make(map[string]string)
	for i := 0; i < len(premapinfo); i+=2 {
		mapmap[premapinfo[i]] = premapinfo[i+1]
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
				reply := mapmap[update.Message.CommandArguments()]
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
			case "all":
				for i := 0; i < len(premapinfo); i+=2{
					bot.Send(tgbotapi.NewMessage(ChatID,"Name: " + premapinfo[i]+"; Phone: "+premapinfo[i+1]))
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
