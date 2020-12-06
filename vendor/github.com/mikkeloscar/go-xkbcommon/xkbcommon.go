package xkb

/*
#cgo LDFLAGS: -lxkbcommon
#include <stdlib.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"
import "unsafe"

type KeySym C.xkb_keysym_t

type KeySymFlags C.enum_xkb_keysym_flags

const (
	KeySymNoFlags         KeySymFlags = 0
	KeySymCaseInsensitive             = (1 << 0)
)

func KeySymFromName(name string, flags KeySymFlags) KeySym {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return KeySym(C.xkb_keysym_from_name(cname, C.enum_xkb_keysym_flags(flags)))
}
