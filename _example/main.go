package main

import (
	"github.com/Toyz/gogato"
	"github.com/Toyz/gogato/_example/actions"
	"log"
)

func main() {
	ga := gogato.NewGogato()
	if ga == nil {
		return
	}

	log.Println(ga.RegisterAction(&actions.CounterAction{Count: 0}))

	log.Println(ga.Run())
}