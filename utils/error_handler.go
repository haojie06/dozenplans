package utils

import (
	"log"
)

func PanicErr(err error, info string) {
	if err != nil {
		log.Panic(info, "\n", err.Error())
	}
}

func LogErr(err error, info string) {
	if err != nil {
		log.Println(info, "\n", err.Error())
	}
}
