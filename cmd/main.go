package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"time"
)

func main() {
	kurs := 100000.0
	label := ""

	for {
		bot, err := tgbotapi.NewBotAPI("6134240801:AAHgp7KV5erel6IkBRG2L2e-6ujf5MIi26k")
		if err != nil {
			log.Panic(err)
		}

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		c := colly.NewCollector()
		c.OnHTML("#USD_ds", func(e *colly.HTMLElement) {
			k1, _ := strconv.ParseFloat(e.Text, 32)
			if kurs > k1 {
				label = "kurs"
				kurs = k1
			}
		})

		err = c.Visit("https://kantorkurs.pl/")

		if err != nil {
			panic(err)
		}

		// -------------------------------------------
		c = colly.NewCollector()
		c.OnHTML("td.c_sell", func(e *colly.HTMLElement) {
			if e.Index == 2 {
				k1, _ := strconv.ParseFloat(e.Text, 32)
				if kurs > k1 {
					label = "cent"
					kurs = k1
				}
			}
		})

		err = c.Visit("https://m.centkantor.pl/")

		if err != nil {
			panic(err)
		}

		// -------------------------------------------
		c = colly.NewCollector()
		c.OnHTML("td.list-table__col--value", func(e *colly.HTMLElement) {
			if e.Index == 115 {
				k1, _ := strconv.ParseFloat(e.Text, 32)
				if kurs > k1 {
					label = "tavex"
					kurs = k1
				}
			}
		})

		err = c.Visit("https://tavex.pl/kantor-wroclaw-kursy-walut/")

		if err != nil {
			panic(err)
		}

		bot.Send(tgbotapi.NewMessage(-1001638150368, "В канторі "+label+" зараз найнижчий курс, долар можна купити по "+fmt.Sprintf("%f", kurs)+" злотих"))

		time.Sleep(1 * time.Hour)
	}
}
