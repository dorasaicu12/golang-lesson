package models

import "github.com/dorasaicu12/booking/internal/form"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form *form.Form
	IsAuthenticated int
	UrlGlobal string
}
