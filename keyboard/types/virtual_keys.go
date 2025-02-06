package types

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// VKCode represents Microsoft defined virtual key codes.
//
// For more details, see the MSDN.
//
// https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
type VirtualKey uint32

func (vk VirtualKey) String() string {
	if name, exists := VKCodeToKeyName[vk]; exists {
		return fmt.Sprintf("VK_%s (%#x)", name,int(vk))
	}
	return fmt.Sprintf("VK_%x", int(vk))
}

func (vk VirtualKey) Name() string {
	if name, exists := VKCodeToKeyName[vk]; exists {
		return fmt.Sprintf("VK_%s", name)
	}
	return fmt.Sprintf("VK_%#x", vk)
}

func (vk VirtualKey) KeybindName() string { 
	if name, exists := VKCodeToKeyName[vk]; exists {
		return fmt.Sprintf("%s", name)
	}
	return fmt.Sprintf("%#x", vk)
}

func (vk VirtualKey) IsModifier() bool {
	_, exists := VKModifierMap[vk]
	return exists
}

type KeyBinding []VirtualKey

func NewKeybinding(vks ...VirtualKey) KeyBinding {
	return vks
}


// ExactMatch checks if the keybinding only contains the same keys and nothing else.
// The keys may be in any order.
func (kb KeyBinding) Match(match []VirtualKey) bool {
	if len(kb) != len(match) {
		return false
	}

	for _, key := range kb {
		if !slices.Contains[KeyBinding](match,key){
			return false
		}
	}

	return true
}

func (kb KeyBinding) SubsetOf(match KeyBinding) bool {
	if len(kb) > len(match) {
		return false
	}
	
	for _, key := range kb {
		if !match.Contains(key) {
			return false
		}
	}
	return true
}

func (kb KeyBinding) Contains(vk VirtualKey) bool {
	for _, key := range kb {
		if key == vk {
			return true
		}
	}
	return false
}

func (kb KeyBinding) String() string {
	var keys []string
	for _, key := range kb {
		keys = append(keys, key.String())
	}
	return fmt.Sprintf("%v", keys)
}
func (kb KeyBinding) PrettyString() string {
	KeybindSort(kb)
	var keys []string
	for _, key := range kb {
		keys = append(keys, key.KeybindName())
	}
	return strings.Join(keys, " + ")
}

func (kb KeyBinding) Name() string {
	var keys []string
	for _, key := range kb {
		keys = append(keys, key.Name())
	}
	return fmt.Sprintf("%v", keys)
}

func (kb KeyBinding) KeybindName() string {
	var keys []string
	for _, key := range kb {
		keys = append(keys, key.KeybindName())
	}
	return strings.Join(keys, " + ")
}

// VirtualKeySorter sorts VirtualKey in a nice printable order. 
// So that modifiers are first and by alphabetical order. 
// Then the rest of the keys are sorted alphabetically.
func KeybindSort(vks []VirtualKey) {
	sort.Slice(vks, func(i, j int) bool {
		if vks[i].IsModifier() && !vks[j].IsModifier() {
			return true
		}
		if !vks[i].IsModifier() && vks[j].IsModifier() {
			return false
		}
		return vks[i] < vks[j]
	})
}

func GetVirtualKey(name string) (VirtualKey, bool) {
	vk, exists := KeyNameToVKCode[name]
	return vk, exists
}
