package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"GoCourseBasic/homework-5/server"

	"github.com/jinzhu/gorm"
  	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	db, err := gorm.Open("mysql", "root:root@/blog")

	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	defer db.Close()
	serv := server.New(lg, db)

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
