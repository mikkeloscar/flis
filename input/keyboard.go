package input

import xkb "github.com/mikkeloscar/go-xkbcommon"

// KeyState holds the current state of pressed keys.
type KeyState struct {
	keySyms  map[xkb.KeySym]struct{}
	keyCodes map[uint32]struct{}
}

// NewKeyState initializes a new empty KeyState.
func NewKeyState() *KeyState {
	return &KeyState{
		keySyms:  make(map[xkb.KeySym]struct{}),
		keyCodes: make(map[uint32]struct{}),
	}
}

// type keys struct {
// 	keySym xkb.KeySym
// 	altSym xkb.KeySym
// }

// PressKey adds a key to the keystate.
func (k *KeyState) PressKey(keySym xkb.KeySym, keyCode uint32) {
	if keyCode == 0 {
		return
	}

	if !k.IsPressed(keySym, keyCode) {
		k.keyCodes[keyCode] = struct{}{}
		k.keySyms[keySym] = struct{}{}
	}
}

// ReleaseKey removes a key from the keystate.
func (k *KeyState) ReleaseKey(keySym xkb.KeySym, keyCode uint32) {
	if k.IsPressed(keySym, keyCode) {
		delete(k.keySyms, keySym)
		delete(k.keyCodes, keyCode)
	}
}

// IsPressed returns true if key is pressed.
func (k *KeyState) IsPressed(keySym xkb.KeySym, keyCode uint32) bool {
	if _, ok := k.keySyms[keySym]; ok {
		return true
	}

	if _, ok := k.keyCodes[keyCode]; ok {
		return true
	}

	return false
}
