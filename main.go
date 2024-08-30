package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"wincuts/keyboard"
	"wincuts/keyboard/types"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}



func run() error {
	hook,err:= keyboard.NewHook()
	if err != nil {
		return fmt.Errorf("failed to create new hook: %v", err)
	}
	hook.Start()


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")

	for {
		select {
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case k := <- hook.Subscribe():
			if k.KeyDown {
				continue
			}
			keyB := types.NewKeybinding(k.PressedKeys...)
			fmt.Println("Keys: ", keyB.PrettyString())
			continue
		}
	}
}