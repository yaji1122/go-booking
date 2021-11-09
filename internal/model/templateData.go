package model

import "github.com/yaji1122/bookings-go/internal/forms"

//TemplateData holds data sent from handler for template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	//如果不確定type，直接用interface
	Data      map[string]interface{}
	CSRFToken string
	//A Flash message to the user
	Flash   string
	Warning string
	//Error message
	Error string
	Form  *forms.Form
}
