package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/PuerkitoBio/goquery"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/jinzhu/configor"
)

type Config struct {
    Token     string
    SourceUrl string
}

func main() {
    var config Config
    err := configor.Load(&config, "config.json")
    if err != nil {
        log.Printf("[ERROR] Failed to read config.json: %v", err)
    }

    bot, err := tgbotapi.NewBotAPI(config.Token)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = false

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore any non-Message Updates
            continue
        }

        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        var answer string
        if update.Message.Text == "/start" {
            answer = "Welcome to Covid19 news bot!\nType any message to get stat\nCreated by deface (t.me/deface90)"
        }

        res, err := http.Get(config.SourceUrl)
        if err != nil {
            log.Fatal(err)
        }

        if res.StatusCode != 200 {
            log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
        }

        // Load the HTML document
        doc, err := goquery.NewDocumentFromReader(res.Body)
        if err != nil {
            log.Fatal(err)
        }

        doc.Find("#maincounter-wrap").Each(func(i int, s *goquery.Selection) {
            h := s.Find("h1").Text()
            if h == "Coronavirus Cases:" {
                answer += "Coronavirus Cases: " + s.Find("span").Text()
            }
            if h == "Recovered:" {
                answer += "\nRecovered: " + s.Find("span").Text()
            }
            if h == "Deaths:" {
                answer += "\nDeaths: " + s.Find("span").Text()
            }
        })

        doc.Find("#main_table_countries_today a.mt_a").Each(func(i int, s *goquery.Selection) {
            if s.Text() != "Russia" {
                return
            }

            answer += "\n\n Russia:"
            cels := s.ParentsFiltered("tr").ChildrenFiltered("td")
            answer += fmt.Sprintf("\nCases: %v (%v)",cels.Eq(1).Text(), cels.Eq(2).Text())
            answer += fmt.Sprintf("\nRecovered: %v", cels.Eq(5).Text())
            answer += fmt.Sprintf("\nDeaths: %v (%v)",cels.Eq(3).Text(), cels.Eq(4).Text())
        })

        _ = res.Body.Close()

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)

        _, err = bot.Send(msg)
        if err != nil {
            log.Printf("[ERROR] Failed to send message: %v", err)
        }
    }
}
