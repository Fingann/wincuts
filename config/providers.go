// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"wincuts/keyboard/types"
)

// DefaultActionProvider provides the default set of actions.
// This follows the Open/Closed Principle by allowing new actions to be added without modifying existing code.
type DefaultActionProvider struct{}

// GetActions implements ActionProvider.
func (d *DefaultActionProvider) GetActions() map[string]Action {
	return map[string]Action{
		"SwitchDesktop": {
			Name:        "SwitchDesktop",
			Description: "Switch to the specified virtual desktop",
			ParamTypes:  []string{"desktop"},
			Validator:   validateSwitchDesktop,
		},
		"MoveWindowToDesktop": {
			Name:        "MoveWindowToDesktop",
			Description: "Move the active window to specified desktop and switch to it",
			ParamTypes:  []string{"desktop"},
			Validator:   validateMoveWindowToDesktop,
		},
		"CreateDesktop": {
			Name:        "CreateDesktop",
			Description: "Create a new virtual desktop",
			ParamTypes:  []string{},
			Validator:   validateCreateDesktop,
		},
	}
}

// DefaultKeyProvider provides the default set of valid keys.
// This follows the Open/Closed Principle by allowing new keys to be added without modifying existing code.
type DefaultKeyProvider struct{}

// GetValidKeys implements KeyProvider.
func (d *DefaultKeyProvider) GetValidKeys() map[string]types.VirtualKey {
	return map[string]types.VirtualKey{
		// Modifiers
		"LAlt":   types.VK_LMENU,
		"RAlt":   types.VK_RMENU,
		"LShift": types.VK_LSHIFT,
		"RShift": types.VK_RSHIFT,
		"LCtrl":  types.VK_LCONTROL,
		"RCtrl":  types.VK_RCONTROL,
		// Numbers
		"1": types.VK_1,
		"2": types.VK_2,
		"3": types.VK_3,
		"4": types.VK_4,
		"5": types.VK_5,
		"6": types.VK_6,
		"7": types.VK_7,
		"8": types.VK_8,
		"9": types.VK_9,
		// Add more keys as needed
	}
}

// Action validators
func validateSwitchDesktop(params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("SwitchDesktop requires exactly one parameter")
	}
	return nil
}

func validateMoveWindowToDesktop(params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("MoveWindowToDesktop requires exactly one parameter")
	}
	return nil
}

func validateCreateDesktop(params []string) error {
	if len(params) != 0 {
		return fmt.Errorf("CreateDesktop takes no parameters")
	}
	return nil
}
