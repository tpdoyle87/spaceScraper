package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"spaceScraper/internal/config"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	log.Info("Starting API")
	conf := config.New()

	d := conf.Database
	host := d.Host // Should be 'db' in Docker Compose setup
	port := d.Port // Default should be '5432'
	user := d.User // Default should be 'postgres'
	password := d.Password
	dbname := d.Dbname
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
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/articles", GetArticles(db))
	}
	err = router.Run("3000")
	if err != nil {
		return
	}
}

func GetArticles(db *sql.DB) gin.HandlerFunc {
	result, err := db.Exec("WITH RankedArticles AS (SELECT *, ROW_NUMBER() OVER (PARTITION BY Company ORDER BY Date DESC) AS rn FROM articles")
	if err != nil {
		log.Info("couldn't get articles")
	}
	return func(c *gin.Context) {
		c.JSON(200, result)
	}
}
