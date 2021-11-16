package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

const InProduction = false

type Configuration struct {
	UseCache      bool
	InProduction  bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
}
