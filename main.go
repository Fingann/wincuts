package main

import (
	"fmt"
	"wincuts/keylogger"
)

func main() {
	hook := keylogger.NewKeyboardHook()

	// Start the system state which sets the hook and processes key events
	if err := hook.Start(); err != nil {
		fmt.Println(err)
		return
	}

	// Process key events in a separate goroutine
	go func() {
		for event := range hook.Subscribe() {
			if !event.KeyDown {
				continue
			}
			keynames := []string{}
			for vkCode := range event.PressedKeys {
				name, ok := keylogger.CodeToVKMap[vkCode]
				if !ok {
					continue
				}
				keynames = append(keynames, name)
			}
			fmt.Println("Pressed keys:", keynames)
		}
	}()

	// Example: Stop the hook after some time (10 seconds)
	// time.Sleep(10 * time.Second)
	// hook.Stop()

	// Prevent main from exiting
	select {}
}
