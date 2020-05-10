package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	"GoCourseBasic/homework-6/server"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
)

// NewLogger - Создаёт новый логгер
func NewLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)
	return lg
}

func main() {
	flagServAddr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	lg := NewLogger()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	db := client.Database("blog")

	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	serv := server.New(lg, ctx, db)

	go func() {
		err := serv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt, os.Kill)
	<-stopSig
}
