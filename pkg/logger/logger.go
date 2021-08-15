package logger

import (
	"log"
	"os"
)

type SnippetLogger struct {
	Info  *log.Logger
	Error *log.Logger
}

func New() *SnippetLogger {
	return &SnippetLogger{
		Info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		Error: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
