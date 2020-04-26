package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"GoCourseBasic/homework-4/models"
	
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Server - объект сервера
type Server struct {
	lg            *logrus.Logger
	db            *sql.DB
	templatesDir  string
	indexTemplate string
	Page          models.Page
}

// New - создаёт новый экземпляр сервера
func New(lg *logrus.Logger, db *sql.DB) *Server {
	return &Server{
		lg:            lg,
		db:            db,
		templatesDir:  "./templates",
		indexTemplate: "index.html",
		Page: models.Page{
			Posts: models.PostItemSlice{
				{Id: "0", Title: "123", Date: "12345", Summary: "erferf", Body: "afasdf", File: "afasf"},
				{Id: "1", Title: "1sdfv23", Date: "1adfv2345", Summary: "eadfvrferf", Body: "aafvfasdf", File: "afvafasf"},
			},
		},
	}
}

// Start - запускает сервер
func (serv *Server) Start(addr string) error {
	r := chi.NewRouter()
	serv.bindRoutes(r)
	serv.lg.Debug("server is started ...")
	return http.ListenAndServe(addr, r)
}

// SendErr - отправляет ошибку пользователю и логирует её
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := models.ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - отправляет 500 ошибку
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}
