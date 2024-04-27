package main

import (
	"database/sql"
	"fmt"
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
	host := database.Host // Should be 'db' in Docker Compose setup
	port := database.Port // Default should be '5432'
	user := database.User // Default should be 'postgres'
	password := database.Password
	dbname := database.Dbname
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
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
	err = c.AddFunc("@hourly", func() {
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
