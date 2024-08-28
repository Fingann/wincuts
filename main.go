package main

import (
	"fmt"
	"os"
	"os/signal"
	"wincuts/keyboard"
	"wincuts/keyboard/shortcut"
	//"wincuts/keyboard/code"
)

func main() {
	hook := keyboard.NewHook()
	//mapper := code.NewMapper()

	// Start the system state which sets the hook and processes key events
	if err := hook.Start(); err != nil {
		fmt.Println(err)
		return
	}
	defer hook.Stop()

	eventChan := hook.Subscribe()

	keybindingService := shortcut.NewKeybindingService(eventChan)
	bind,err := shortcut.NewBindingAction([]string{"VK_LSHIFT", "VK_A"}, true, func() error {
		fmt.Println("Ctrl+Shift+A pressed")
		return nil
	})
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to create binding action: %v", err))
		return
	}


	keybindingService.RegisterKeyBindingActions(bind).Start()

	// Process key events in a separate goroutine
//go func() {
//	for event := range hook.Subscribe() {
//		if !event.KeyDown {
//			continue
//		}
//		prettyString,err:=mapper.PrettyPrint(event.PressedKeys)
//		if err != nil {
//			fmt.Println(fmt.Errorf("Failed to get pretty string: %v", err))
//			continue
//		}
//		
//		fmt.Println("Pressed keys:", prettyString)
//	}
//}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
