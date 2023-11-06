package utils

import (
	"log"
	"os"
	"path/filepath"
)

var InfoLogger *log.Logger
var ErrorLogger *log.Logger

func init() {
	path, err := filepath.Abs("./logs")
	if err != nil {
		log.Println("Error riding absolute path: ", err)
		return
	}

	myLog, err := os.OpenFile(path+"/momentum-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}

	InfoLogger = log.New(myLog, "[Info]:\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	ErrorLogger = log.New(myLog, "[Error]:\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
}
