package code

import (
	"fmt"
	"slices"
	"strings"
)

// KeyMapper interface defines methods for querying key codes and names.
type Mapper interface {
    GetKeyCode(name string) (uint32, error)
    GetKeyName(code uint32) (string, error)
    GetKeyCodes(names []string) ([]uint32, error)
    GetKeyNames(codes []uint32) ([]string, error)
    KeysToCodeMap(names []string) (map[uint32]string, error)
	PrettyPrint(codes []uint32) (string, error)

}

func NewMapper() Mapper {
	return &KeyMapperImpl{}
}

var _ Mapper = &KeyMapperImpl{}
// KeyMapperImpl is an implementation of the KeyMapper interface.
type KeyMapperImpl struct{}

// PrettyPrint returns a string representation of the key codes.
func (km *KeyMapperImpl) PrettyPrint(codes []uint32) (string, error) {
    slices.Sort[[]uint32](codes)
    slices.Reverse(codes)
    result,err := km.GetKeyNames(codes)
    if err != nil {
        return "", fmt.Errorf("failed to get key names: %v", err)
    }

	return strings.Join(result," + "), nil
}

//KeysToCodeMap returns a map of key names to their corresponding key codes.
func (km *KeyMapperImpl) KeysToCodeMap(names []string) (map[uint32]string, error) {
    result := make(map[uint32]string)
    for _, name := range names {
        code, err := km.GetKeyCode(name)
        if err != nil {
            return nil, err
        }
        result[code] = name
    }
    return result, nil
}


// GetKeyCode returns the key code for a given key name.
func (km *KeyMapperImpl) GetKeyCode(name string) (uint32, error) {
    if keyCode, exists := VKMap[name]; exists {
        return keyCode, nil
    }
    return 0, fmt.Errorf("key name %s not found in VKToCodeMap", name)
}

// GetKeyName returns the key name for a given key code.
func (km *KeyMapperImpl) GetKeyName(code uint32) (string, error) {
    for name, keyCode := range VKMap {
        if keyCode == code {
            return name, nil
        }
    }
    return "", fmt.Errorf("key code %d not found in VKToCodeMap", code)
}

// GetKeyCodes returns a map of key names to their corresponding key codes.
func (km *KeyMapperImpl) GetKeyCodes(names []string) ([]uint32, error) {
    result := make([]uint32, len(names))
    for i, name := range names {
        code, err := km.GetKeyCode(name)
        if err != nil {
            return nil, err
        }
        result[i] = code
    }
    return result, nil
}

// GetKeyNames returns a map of key codes to their corresponding key names.
func (km *KeyMapperImpl) GetKeyNames(codes []uint32) ([]string, error) {
    result := make([]string, len(codes))
    for i, code := range codes {
        name, err := km.GetKeyName(code)
        if err != nil {
            return nil, err
        }
        result[i] = name
    }
    return result, nil
}