package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"os"
)

type Article struct {
	Title   string
	Link    string
	Company string
	Date    string
}

// article.Title = e.ChildText(".article-name")
func main() {
	file, err := os.Create("articles.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Company", "Title", "Link", "Date"}
	writer.Write(headers)

	urls := []string{"https://www.space.com/", "https://www.seti.org/news"}
	c := colly.NewCollector()
	article := Article{}
	for _, url := range urls {
		switch url {
		case "https://www.space.com/":
			c.OnHTML(".article-link", func(e *colly.HTMLElement) {
				WriteSpace(e, article, writer)
			})
		case "https://www.seti.org/news":
			c.OnHTML("div[class=row]", func(e *colly.HTMLElement) {
				e.ForEach(".col-md-3", func(_ int, e *colly.HTMLElement) {
					WriteSeti(article, e, writer)
				})
			})
		}
		c.Visit(url)
	}
}

func WriteSpace(e *colly.HTMLElement, article Article, writer *csv.Writer) {
	article.Company = "space.com"
	article.Title = e.ChildText(".article-name")
	article.Link = e.Attr("href")
	err := writer.Write([]string{article.Company, article.Title, article.Link, article.Date})
	if err != nil {
		panic(err)
	}
}

func WriteSeti(article Article, e *colly.HTMLElement, writer *csv.Writer) {
	article.Company = "seti.org"
	article.Title = e.ChildText("a")
	article.Link = "https://www.seti.org" + e.ChildAttr("a", "href")
	article.Date = e.ChildText(".views-field-field-date")
	err := writer.Write([]string{article.Company, article.Title, article.Link, article.Date})
	if err != nil {
		panic(err)
	}
}
