package handler

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/model"
	"github.com/yaji1122/bookings-go/internal/pageRenderer"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var appConfig config.Configuration
var session *scs.SessionManager
var functions = template.FuncMap{}
var rootPath = "./../../templates"

func getRoutes() http.Handler {

	//what am I goin to put in the session
	gob.Register(model.Reservation{})

	//change to true when in production
	appConfig.InProduction = false

	//初始化session
	log.Println("初始化Session Manager")
	session = scs.New()
	session.Lifetime = 30 * time.Minute
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //localhost use http, in product will be true
	appConfig.Session = session
	//產生 Template Cache
	log.Println("產生Template Cache")
	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Error Creating Template Cache")
	}
	//將產生的Template Cache指定到 AppConfig中
	appConfig.TemplateCache = templateCache
	//傳入 Configuration
	appConfig.UseCache = true //dev mode 設為False
	//初始化Validator
	model.InitialValidator()
	// pageRenderer pkg 設定 appConfig
	//set up configs
	repo := InitiateRepository(&appConfig)
	CreateHandler(repo)
	pageRenderer.CreatePageRenderer(&appConfig)
	//here move the routes to the router.go
	//http.HandleFunc("/", Repo.Home)
	//http.HandleFunc("/about", Repo.About)
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(WriteToConsole)
	mux.Get("/booking", Repo.Booking)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/room", Repo.Room)
	mux.Get("/reservation", Repo.Reservation)
	mux.Get("/", Repo.Index)
	mux.Get("/summary", Repo.Summary)

	mux.Post("/search-availability", Repo.Availability)
	mux.Post("/reservation", Repo.PostReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	//_ = http.ListenAndServe(port, nil)
	return mux
}

// NoSurf for csrf check
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	//it uses cookies to make sure that the token it generates is available on per page basis
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// WriteToConsole simplee example
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		url := r.RequestURI
		if !strings.Contains(url, "static") {
			log.Printf("ip[%s] access url - %s", ip, url)
		}
		next.ServeHTTP(w, r)
	})
}

//CreateTestTemplateCache 產生網頁資料，並存成map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//templateMapping := make(map[string]*template.Template)

	//Create a map with index<->template
	templateMapping := map[string]*template.Template{}

	//找出所有的page
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", rootPath))
	if err != nil {
		return templateMapping, err
	}

	//找出layout
	matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", rootPath))
	if err != nil {
		return templateMapping, err
	}

	for _, page := range pages {
		//取得頁面檔名
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateMapping, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", rootPath))
			if err != nil {
				return templateMapping, err
			}
		}
		templateMapping[name] = ts
	}
	return templateMapping, err
}
