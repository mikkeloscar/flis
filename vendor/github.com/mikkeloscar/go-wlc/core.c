#include "_cgo_export.h"
#include <stdlib.h>
#include <string.h>
#include <wlc/wlc.h>

void log_handler_cb(enum wlc_log_type type, const char *str) {
	char *msg = strdup(str);
	_goLogHandlerCb(type, msg);
	free(msg);
}

void wrap_wlc_log_set_handler() {
	wlc_log_set_handler(log_handler_cb);
}

int event_loop_fd_cb(int fd, uint32_t mask, void *arg) {
	_goEventLoopFdCb(fd, mask);
}

struct wlc_event_source *wrap_wlc_event_loop_add_fd(int fd, uint32_t mask) {
	wlc_event_loop_add_fd(fd, mask, event_loop_fd_cb, NULL);
}

int event_loop_timer_cb(void *arg) {
	_goEventLoopTimerCb(*((uint32_t*)arg));
}

struct wlc_event_source *wrap_wlc_event_loop_add_timer(uint32_t id) {
	wlc_event_loop_add_timer(event_loop_timer_cb, &id);
}
