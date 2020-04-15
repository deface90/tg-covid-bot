package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "log"
    "net/http"
    "strings"
    "time"
)

func main() {
    res, err := http.Get("https://www.worldometers.info/coronavirus/")
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

    doc.Find("#newsdate" + time.Now().Format("2006-01-02") + " li.news_li").Each(func(i int, s *goquery.Selection) {
        if i > 5 {
            return
        }
        fmt.Println(strings.Replace(s.Text(), "[source]", "", -1))
    })
}
