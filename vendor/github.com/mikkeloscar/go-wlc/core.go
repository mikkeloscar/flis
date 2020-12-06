package wlc

/*
#cgo LDFLAGS: -lwlc
#include <stdlib.h>
#include <wlc/wlc.h>

// handle wlc_log_set_handler callback.
extern void log_handler_cb(enum wlc_log_type type, const char *str);
extern void wrap_wlc_log_set_handler();

// handle wlc_event_loop_add_fd callback.
extern int event_loop_fd_cb(int fd, uint32_t mask, void *arg);
extern struct wlc_event_source *wrap_wlc_event_loop_add_fd(int fd, uint32_t mask);

// handle wlc_event_loop_add_timer callback.
extern int event_loop_timer_cb(void *arg);
extern struct wlc_event_source *wrap_wlc_event_loop_add_timer(uint32_t id);

// internal wlc_interface reference.
extern struct wlc_interface interface_wlc;
extern void init_interface(uint32_t mask);
*/
import "C"

import (
	"math/rand"
	"time"
	"unsafe"
)

var logHandler func(LogType, string)

//export _goLogHandlerCb
func _goLogHandlerCb(typ C.enum_wlc_log_type, msg *C.char) {
	logHandler(LogType(typ), C.GoString(msg))
}

// LogSetHandler sets log handler. Can be set before Init.
func LogSetHandler(handler func(LogType, string)) {
	logHandler = handler
	C.wrap_wlc_log_set_handler()
}

// Init initializeses wlc. Returns false on failure.
//
// Avoid running unverified code before Init as wlc compositor may be run
// with higher privileges on non logind systems where compositor binary needs
// to be suid.
//
// Init's purpose is to initialize and drop privileges as soon as possible.
func Init() bool {
	return bool(C.wlc_init())
}

// Terminate wlc.
func Terminate() {
	C.wlc_terminate()
}

// GetBackendType queries for the backend wlc is using.
func GetBackendType() BackendType {
	return BackendType(C.wlc_get_backend_type())
}

// Exec program.
func Exec(bin string, arg ...string) {
	// prepend bin to start of args slice as expected by wlc.
	args := append([]string{bin}, arg...)
	cbin := C.CString(bin)
	defer C.free(unsafe.Pointer(cbin))
	cargs := strSlicetoCArray(args)
	defer freeCStrArray(cargs)
	C.wlc_exec(cbin, cargs)
}

// Run event loop.
func Run() {
	C.wlc_run()
}

// TODO make more go friendly

// HandleSetUserData can be used to link custom data to handle.
// Client must allocate and handle the data as some C type.
func HandleSetUserData(handle View, userdata unsafe.Pointer) {
	C.wlc_handle_set_user_data(C.wlc_handle(handle), userdata)
}

// HandleGetUserData gets custom linked user data from handle.
func HandleGetUserData(handle View) unsafe.Pointer {
	return C.wlc_handle_get_user_data(C.wlc_handle(handle))
}

type fdEvent struct {
	cb     func(int, uint32, interface{})
	arg    interface{}
	source EventSource
}

var eventLoopFd = make(map[int]fdEvent)

//export _goEventLoopFdCb
func _goEventLoopFdCb(fd C.int, mask C.uint32_t) {
	if event, ok := eventLoopFd[int(fd)]; ok {
		event.cb(int(fd), uint32(mask), event.arg)
	}
}

// EventLoopAddFd adds fd to event loop.
func EventLoopAddFd(fd int, mask uint32, cb func(int, uint32, interface{}), arg interface{}) EventSource {
	source := EventSource(C.wrap_wlc_event_loop_add_fd(
		C.int(fd),
		C.uint32_t(mask),
	))

	if source != nil {
		eventLoopFd[fd] = fdEvent{
			cb:     cb,
			arg:    arg,
			source: source,
		}
	}

	return source
}

type timerEvent struct {
	cb     func(interface{})
	arg    interface{}
	source EventSource
}

var eventLoopTimer = make(map[uint32]timerEvent)
var timerEventRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func timerEventID() uint32 {
	newID := uint32(0)
	for {
		id := timerEventRand.Uint32()
		if _, ok := eventLoopTimer[id]; !ok {
			newID = id
			break
		}
	}

	return newID
}

//export _goEventLoopTimerCb
func _goEventLoopTimerCb(id C.int32_t) {
	if event, ok := eventLoopTimer[uint32(id)]; ok {
		event.cb(event.arg)
	}
}

// EventLoopAddTimer adds timer to event loop.
func EventLoopAddTimer(cb func(interface{}), arg interface{}) EventSource {
	id := timerEventID()

	source := EventSource(C.wrap_wlc_event_loop_add_timer(C.uint32_t(id)))

	if source != nil {
		eventLoopTimer[id] = timerEvent{
			cb:     cb,
			arg:    arg,
			source: source,
		}
	}

	return source
}

// EventSourceTimerUpdate updates timer to trigger after delay.
// Returns true on success.
func EventSourceTimerUpdate(source EventSource, msDelay int32) bool {
	return bool(C.wlc_event_source_timer_update(
		source,
		C.int32_t(msDelay),
	))
}

// EventSourceRemove removes event source from event loop.
func EventSourceRemove(source EventSource) {
	found := false

	for fd, event := range eventLoopFd {
		if source == event.source {
			found = true
			delete(eventLoopFd, fd)
			break
		}
	}

	if !found {
		for id, event := range eventLoopFd {
			if source == event.source {
				delete(eventLoopFd, id)
			}
		}
	}
	C.wlc_event_source_remove(source)
}
