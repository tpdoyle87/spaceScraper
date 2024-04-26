package internal

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gocolly/colly"
)

type Article struct {
	Id      string
	Title   string
	URL     string
	Company string
	Date    string
	Hash    string
}

func FetchArticles(db *sql.DB) {
	urls := []string{"https://www.space.com/", "https://www.seti.org/news"}
	c := colly.NewCollector(
		colly.Async(true),
	)
	article := Article{}
	for _, url := range urls {
		switch url {
		case "https://www.space.com/":
			c.OnHTML(".article-link", func(e *colly.HTMLElement) {
				WriteSpace(e, article, db)
			})
		case "https://www.seti.org/news":
			c.OnHTML("div[class=row]", func(e *colly.HTMLElement) {
				e.ForEach(".col-md-3", func(_ int, e *colly.HTMLElement) {
					WriteSeti(article, e, db)
				})
			})
		}
		c.Visit(url)
		c.Wait()
	}
}

func WriteSpace(e *colly.HTMLElement, article Article, db *sql.DB) {
	article.Company = "space.com"
	article.Title = e.ChildText(".article-name")
	article.URL = e.Attr("href")

	toHash := article.Company + article.Title + article.URL
	article.Hash = byteArrayToString(sha256.Sum256([]byte(toHash)))
	UpsertArticles(db, article)

}

func WriteSeti(article Article, e *colly.HTMLElement, db *sql.DB) {
	article.Company = "seti.org"
	article.Title = e.ChildText("a")
	article.URL = "https://www.seti.org" + e.ChildAttr("a", "href")
	article.Date = e.ChildText(".views-field-field-date")

	toHash := article.Company + article.Title + article.URL
	article.Hash = byteArrayToString(sha256.Sum256([]byte(toHash)))
	UpsertArticles(db, article)
}

func byteArrayToString(b [32]byte) string {
	return fmt.Sprintf("%x", b)
}

func UpsertArticles(db *sql.DB, article Article) {
	_, err := db.Exec("INSERT INTO articles(company, title, url, date, hash) VALUES($1, $2, $3, $4, $5) ON CONFLICT (hash) DO NOTHING", article.Company, article.Title, article.URL, article.Date, article.Hash)
	if err != nil {
		panic(err)
	}
}
