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
	"github.com/yaji1122/bookings-go/internal/model"
	"github.com/yaji1122/bookings-go/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":8081"

//宣告一個系統設定 AppConfig for same pkg use
var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	pool, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {

		}
	}(pool.SQL)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.Pool, error) {

	//what am I goin to put in the session
	gob.Register(model.Reservation{})

	//change to true when in production
	appConfig.InProduction = false

	//create Logger
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	//初始化session
	log.Println("初始化Session Manager")
	session = scs.New()
	session.Lifetime = 30 * time.Minute
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //localhost use http, in product will be true
	appConfig.Session = session

	//connect to db
	log.Println("Connecting to database")
	pool, err := driver.ConnectSQL("root:53434976@/test?charset=utf8")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	//產生 Template Cache
	log.Println("產生Template Cache")
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error Creating Template Cache")
		return nil, err
	}
	//將產生的Template Cache指定到 AppConfig中
	appConfig.TemplateCache = templateCache
	//傳入 AppConfig
	appConfig.UseCache = false //dev mode 設為False
	//初始化Validator
	model.InitialValidator()
	// render pkg 設定 appConfig
	render.NewRenderer(&appConfig)
	//set up configs
	repo := handler.NewRepo(&appConfig, pool)
	handler.NewHandlers(repo)
	helper.NewHelper(&appConfig)

	//here move the routes to the router.go
	//http.HandleFunc("/", handler.Repo.Home)
	//http.HandleFunc("/about", handler.Repo.About)

	log.Println(fmt.Sprintf("Starting application on port %s http://127.0.0.1%s", port, port))
	//_ = http.ListenAndServe(port, nil)
	return pool, nil
}

//func Divide(w http.ResponseWriter, r *http.Request) {
//	f, err := divideValues(100.0, 0.0)
//	if err != nil {
//		fmt.Fprintf(w, fmt.Sprintf("Error Message: %s", err))
//	} else {
//		fmt.Fprintf(w, fmt.Sprintf( "%f divided by %f is %f", 100.0, 0.0, f))
//	}
//}

//func divideValues(x, y float32) (float32, error) {
//	var result float32
//	if y <= 0.0 {
//		return result, errors.New("Can't not divide by 0")
//	}
//	result = x / y
//	return result, nil
//}
////add two values
//func addValues(x, y int) (int, error) {
//	var sum int
//	sum = x + y
//	return sum, nil
//}
