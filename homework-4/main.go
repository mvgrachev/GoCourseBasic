package main

import (
	"database/sql"
	"flag"
	"os"
	"os/signal"

	"GoCourseBasic/homework-4/server"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
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
	db, err := sql.Open("mysql", "root:root@/blog")
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
