package wlc

/*
#cgo LDFLAGS: -lwlc
#include <stdlib.h>
#include <wlc/wlc.h>
*/
import "C"

import "unsafe"

// KeyboardGetXKBState exposes xkb_state. Can be used for more advanced key
// handling. This is currently only exposed as a C struct.
func KeyboardGetXKBState() *C.struct_xkb_state {
	return C.wlc_keyboard_get_xkb_state()
}

// KeyboardGetXKBKeymap exposes xkb_keymap. Can be used for more advanced key
// handling. This is currently only exposed as a C struct.
func KeyboardGetXKBKeymap() *C.struct_xkb_keymap {
	return C.wlc_keyboard_get_xkb_keymap()
}

// KeyboardGetCurrentKeys gets currently held keys.
func KeyboardGetCurrentKeys() []uint32 {
	var len C.size_t
	keys := C.wlc_keyboard_get_current_keys(&len)
	goKeys := make([]uint32, 0, int(len))
	size := int(unsafe.Sizeof(*keys))
	for i := 0; i < int(len); i++ {
		ptr := unsafe.Pointer(uintptr(unsafe.Pointer(keys)) + uintptr(size*i))
		goKeys[i] = *(*uint32)(ptr)
	}

	return goKeys
}

// KeyboardGetKeysymForKey is an utility function to convert raw keycode to
// keysym. Passed modifiers may transform the key.
func KeyboardGetKeysymForKey(key uint32, mods *Modifiers) uint32 {
	if mods != nil {
		cmods := mods.c()
		defer C.free(unsafe.Pointer(cmods))
		return uint32(C.wlc_keyboard_get_keysym_for_key(C.uint32_t(key), cmods))
	}

	return uint32(C.wlc_keyboard_get_keysym_for_key(C.uint32_t(key), nil))
}

// KeyboardGetUtf32ForKey is an utility function to convert raw keycode to
// Unicdoe/UTF-32 codepoint. Passed modifiers may transform the key.
func KeyboardGetUtf32ForKey(key uint32, mods *Modifiers) uint32 {
	if mods != nil {
		cmods := mods.c()
		defer C.free(unsafe.Pointer(cmods))
		return uint32(C.wlc_keyboard_get_utf32_for_key(C.uint32_t(key), cmods))
	}

	return uint32(C.wlc_keyboard_get_utf32_for_key(C.uint32_t(key), nil))
}

// PointerGetPosition gets current pointer position.
func PointerGetPosition() *Point {
	var pos C.struct_wlc_point
	C.wlc_pointer_get_position(&pos)
	return pointCtoGo(&pos)
}

// PointerSetPosition sets pointer position.
func PointerSetPosition(pos Point) {
	cpos := pos.c()
	defer C.free(unsafe.Pointer(cpos))
	C.wlc_pointer_set_position(cpos)
}
