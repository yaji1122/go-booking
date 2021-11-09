package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/model"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

// NewRenderer sets the config for the template package
var appConfig *config.AppConfig

var rootPath = "./templates"

func NewRenderer(config *config.AppConfig) {
	appConfig = config
}

func AddDefaultData(td *model.TemplateData, r *http.Request) *model.TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")
	td.Error = appConfig.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template 回應請求，回傳對應的Template頁面
func Template(w http.ResponseWriter, r *http.Request, name string, data *model.TemplateData) error {
	name = name + ".page.gohtml"
	var templateCache map[string]*template.Template
	//get the template cache from the appConfig config
	if appConfig.UseCache {
		templateCache = appConfig.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	//map 如果key沒有對應的value, 回傳 nil, false
	t, ok := templateCache[name]
	if !ok {
		return errors.New("cant get Template from cache")
	}
	byteBuffer := new(bytes.Buffer)

	data = AddDefaultData(data, r)

	_ = t.Execute(byteBuffer, data)

	_, err := byteBuffer.WriteTo(w)
	if err != nil {
		log.Fatal("Error writing template to browser", err)
		return err
	}
	//
	//parseTemplate, _ := template.ParseFiles("./templates/" + tmpl + ".page.gohtml")
	////Execute sending template to web browser
	//err = parseTemplate.Execute(w, nil)
	//if err != nil {
	//	fmt.Println("error parsing template:", err)
	//}
	return nil
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
