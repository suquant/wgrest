package wgrest

import (
	"log"
	"os"
)

var (
	// Logger default logger
	Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
)
