package g

import (
	"log"
	"os"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}
