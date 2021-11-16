package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/driver"
	"github.com/yaji1122/bookings-go/internal/handler"
	"github.com/yaji1122/bookings-go/internal/helper"
	"github.com/yaji1122/bookings-go/internal/logger"
	"github.com/yaji1122/bookings-go/internal/model"
	"github.com/yaji1122/bookings-go/internal/pageRenderer"
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":8081"
const inProduction = false

var configuration config.Configuration
var session *scs.SessionManager

func main() {
	//初始化伺服器，產生所需要的設定
	pool, err := initiate()
	checkErr(err)
	//main方法結束時，關閉資料庫連線
	defer func(pool *sql.DB) {
		err := pool.Close()
		checkErr(err)
	}(pool)

	//開啟Http Server 並監聽port
	server := http.Server{
		Addr:    port,
		Handler: routes(&configuration),
	}

	err = server.ListenAndServe()
	checkErr(err)
}

func initiate() (*sql.DB, error) {

	//what am I going to put in the session
	gob.Register(model.Reservation{})

	//產生 Template Cache
	templateCache, err := pageRenderer.CreateTemplateCache()
	checkErr(err)

	//產生 http Session
	session = scs.New()
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //localhost use http, in product will be true
	session.Lifetime = 30 * time.Minute

	//設定載入Configuration
	log.Println("初始化 Configuration")
	configuration.InProduction = config.InProduction
	if configuration.InProduction {
		configuration.UseCache = false
	} else {
		log.Println("開發模式：不使用Cache")
		configuration.UseCache = true
	}
	//初始化Logger
	logger.CreateLogger()

	configuration.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	configuration.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	configuration.Session = session
	configuration.TemplateCache = templateCache

	pool := driver.CreateDatabaseConnectionPool("root:53434976@/test?charset=utf8")
	logInstance := logger.CreateLogger()

	//初始化Validator
	model.InitialValidator()
	// pageRenderer pkg 設定 configuration
	pageRenderer.CreatePageRenderer(&configuration)
	//set up configs
	handler.CreateHandler(logInstance, session, pool)
	helper.NewHelper(&configuration)

	log.Println(fmt.Sprintf("Starting application on port %s http://127.0.0.1%s", port, port))

	return pool, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
