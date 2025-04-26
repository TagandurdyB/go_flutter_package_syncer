package helpers

import (
	// "fmt"
	"log"
)

func ErrH(err ...interface{}) bool {
	e := err[len(err)-1]
	if e != nil {
		logSave(err...)
		return true
	}
	return false
}

func logSave(args ...interface{}) {
	log.Print(args...)
	// text := now() + "   " + fmt.Sprint(args...) + "\n"
	// date, _ := newDay()
	// appendFile(date, text)
}

func logSavef(format string, v ...any) {
	log.Printf(format, v...)
	// text := now() + "   " + fmt.Sprintf(format, v...)
	// date, _ := newDay()
	// appendFile(date, text)
}
