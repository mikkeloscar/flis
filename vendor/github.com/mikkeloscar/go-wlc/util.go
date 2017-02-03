package wlc

/*
#cgo LDFLAGS: -lwlc
#include <stdlib.h>
#include <wlc/wlc.h>

char **char_array_init(size_t len) {
	char **arr = malloc(len * sizeof(char*));
	return arr;
}

void char_array_insert(char **arr, char *item, int index) {
	arr[index] = item;
}

void char_array_free(char **arr) {
	int i = 0;
	for (char **ptr = arr; *ptr; ++ptr) {
		free(*ptr);
	}
	free(arr);
}

wlc_handle *handle_array_init(size_t len) {
	wlc_handle *arr = malloc(len * sizeof(wlc_handle));
	return arr;
}

void handle_array_insert(wlc_handle *arr, wlc_handle item, int index) {
	arr[index] = item;
}
*/
import "C"

import "unsafe"

type handleType uint8

const (
	viewHandle handleType = iota
	outputHandle
)

// Initialize a C NULL terminated *char[] from a []string
func strSlicetoCArray(arr []string) **C.char {
	carr := C.char_array_init(C.size_t(len(arr) + 1))
	for i, s := range arr {
		C.char_array_insert(carr, C.CString(s), C.int(i))
	}
	C.char_array_insert(carr, nil, C.int(len(arr)))
	return carr
}

// Free a *char[]
func freeCStrArray(arr **C.char) {
	C.char_array_free(arr)
}

func outputHandlesCArraytoGoSlice(handles *C.wlc_handle, len int) []Output {
	goHandles := make([]Output, len)
	size := int(unsafe.Sizeof(*handles))
	for i := 0; i < len; i++ {
		ptr := unsafe.Pointer(uintptr(unsafe.Pointer(handles)) + uintptr(size*i))
		goHandles[i] = *(*Output)(ptr)
	}

	return goHandles
}

func viewHandlesCArraytoGoSlice(handles *C.wlc_handle, len int) []View {
	goHandles := make([]View, len)
	size := int(unsafe.Sizeof(*handles))
	for i := 0; i < len; i++ {
		ptr := unsafe.Pointer(uintptr(unsafe.Pointer(handles)) + uintptr(size*i))
		goHandles[i] = *(*View)(ptr)
	}

	return goHandles
}

func viewHandlesSliceToCArray(arr []View) (*C.wlc_handle, C.size_t) {
	carr := C.handle_array_init(C.size_t(len(arr)))
	for i, h := range arr {
		C.handle_array_insert(carr, C.wlc_handle(h), C.int(i))
	}

	return carr, C.size_t(len(arr))
}
