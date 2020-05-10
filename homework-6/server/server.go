package server

import (
	//"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"GoCourseBasic/homework-6/models"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
)

// Server - объект сервера
type Server struct {
	lg            *logrus.Logger
	ctx           context.Context
	db            *mongo.Database
	templatesDir  string
	indexTemplate string
	Page          models.Page
}

// New - создаёт новый экземпляр сервера
func New(lg *logrus.Logger, ctx context.Context, db *mongo.Database) *Server {
	return &Server{
		lg:            lg,
		ctx:           ctx,
		db:            db,
		templatesDir:  "./templates",
		indexTemplate: "index.html",
		Page: models.Page{
			Posts: models.PostItemSlice{
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
