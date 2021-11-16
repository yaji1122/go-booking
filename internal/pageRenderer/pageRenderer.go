package pageRenderer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/helper"
	"github.com/yaji1122/bookings-go/internal/model"
	"html/template"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var templateCache map[string]*template.Template
var configuration *config.Configuration

var rootPath = "./templates"

func CreatePageRenderer(config *config.Configuration) {
	//建立cache
	cache, err := CreateTemplateCache()
	if err != nil {
		panic(err)
	}
	templateCache = cache
	configuration = config
}

func AddDefaultData(td *model.TemplateData, r *http.Request) *model.TemplateData {
	td.Flash = configuration.Session.PopString(r.Context(), "flash")
	td.Warning = configuration.Session.PopString(r.Context(), "warning")
	td.Error = configuration.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template 回應請求，回傳對應的Template頁面
func Template(w http.ResponseWriter, r *http.Request, name string, data *model.TemplateData) {
	name = name + ".page.gohtml"
	//get the template cache from the configuration config
	if !configuration.UseCache {
		templateCache, _ = CreateTemplateCache()
	}

	//map 如果key沒有對應的value, 回傳 nil, false
	t, ok := templateCache[name]
	if !ok {
		checkErr(w, errors.New("cant get Template from cache"))
	}
	byteBuffer := new(bytes.Buffer)

	data = AddDefaultData(data, r)

	_ = t.Execute(byteBuffer, data)

	_, err := byteBuffer.WriteTo(w)
	checkErr(w, err)
}

//CreateTemplateCache 產生網頁資料，並存成map
func CreateTemplateCache() (map[string]*template.Template, error) {
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

func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		helper.ServerError(w, err)
	}
}
