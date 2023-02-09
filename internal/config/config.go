package config

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/dorasaicu12/booking/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	ORM *sql.DB
	MailChan chan models.MailData
}
