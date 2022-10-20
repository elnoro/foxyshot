package logger

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultSize = 10 // megabytes
	defaultAge  = 15 // days
)

func FromArgs(args []string) {
	logFile, err := parseLog(args)
	if err != nil {
		log.Fatal("invalid log file specified", err)
	}

	setUp(logFile)
}

func parseLog(args []string) (string, error) {
	f := flag.NewFlagSet("logger", flag.ExitOnError)
	logFile := f.String("logfile", "", "path to file, empty means stdout")
	err := f.Parse(args)
	if err != nil {
		return "", fmt.Errorf("parsing loggin args")
	}

	return *logFile, nil
}

func setUp(file string) {
	log.SetOutput(&lumberjack.Logger{
		Filename: file,
		MaxSize:  defaultSize,
		MaxAge:   defaultAge,
	})
}
