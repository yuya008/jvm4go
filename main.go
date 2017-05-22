package main

import (
	"github.com/yuya008/jvm4go/cmd"
	"log"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Panicln(err)
	}
}