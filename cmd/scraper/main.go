package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"net/http"
	"spaceScraper/internal"
	"spaceScraper/internal/config"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	log.Info("Starting scraper")
	conf := config.New()

	database := conf.Database
	connStr := "user=" + database.User + " dbname=" + database.Dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Errorln(err)
		}
	}(db)

	c := cron.New()
	err = c.AddFunc("0 * * * *", func() {
		log.Info("Fetching articles")
		internal.FetchArticles(db)
	})
	if err != nil {
		return
	}

	log.Info("Starting cron")
	c.Start()
	printCronEntries(c.Entries())
	defer c.Stop()

	http.ListenAndServe(":8080", nil)

}

func printCronEntries(cronEntries []*cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}
