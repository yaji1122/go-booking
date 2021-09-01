package render

import (
	"bytes"
	"github.com/yaji1122/bookings-go/pkg/config"
	"github.com/yaji1122/bookings-go/pkg/model"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

// NewTemplates sets the config for the template package
var appConfig *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultdata(td *model.TemplateData) *model.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *model.TemplateData) {
	tmpl = tmpl + ".page.gohtml"
	var tc map[string]*template.Template
	//get the template cache from the appConfig config
	if appConfig.UseCache {
		tc = appConfig.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//map 如果key沒有對應的value, 回傳 nil, false
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("can't get matching template")
	}
	byteBuffer := new(bytes.Buffer)

	data = AddDefaultdata(data)

	_ = t.Execute(byteBuffer, data)

	_, err := byteBuffer.WriteTo(w)
	if err != nil {
		log.Fatal("Error writing template to browser", err)
	}
	//
	//parseTemplate, _ := template.ParseFiles("./templates/" + tmpl + ".page.gohtml")
	////Execute sending template to web browser
	//err = parseTemplate.Execute(w, nil)
	//if err != nil {
	//	fmt.Println("error parsing template:", err)
	//}
}

//CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)

	//Create a map with index<->template
	myCache := map[string]*template.Template{}

	//get the gohtml that match the pattern
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		//get the name of the template page
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		//find any layout match the template
		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, err
}
