#include "_cgo_export.h"
#include <wlc/wlc.h>

/* output */
bool handle_output_created(wlc_handle output) {
	return _goHandleOutputCreated(output);
}

void handle_output_destroyed(wlc_handle output) {
	_goHandleOutputDestroyed(output);
}

void handle_output_focus(wlc_handle output, bool focus) {
	_goHandleOutputFocus(output, focus);
}

void handle_output_resolution(wlc_handle output, const struct wlc_size *from, const struct wlc_size *to) {
	struct wlc_size nc_from = {
		.w = from->w,
		.h = from->h
	};
	struct wlc_size nc_to = {
		.w = to->w,
		.h = to->h
	};
	_goHandleOutputResolution(output, &nc_from, &nc_to);
}

void handle_output_render_pre(wlc_handle output) {
	_goHandleOutputRenderPre(output);
}

void handle_output_render_post(wlc_handle output) {
	_goHandleOutputRenderPost(output);
}

void handle_output_context_created(wlc_handle output) {
	_goHandleOutputContextCreated(output);
}

void handle_output_context_destroyed(wlc_handle output) {
	_goHandleOutputContextDestroyed(output);
}


/* view */
bool handle_view_created(wlc_handle view) {
	return _goHandleViewCreated(view);
}

void handle_view_destroyed(wlc_handle view) {
	_goHandleViewDestroyed(view);
}

void handle_view_focus(wlc_handle view, bool focus) {
	_goHandleViewFocus(view, focus);
}

void handle_view_move_to_output(wlc_handle view, wlc_handle from_output, wlc_handle to_output) {
	_goHandleViewMoveToOutput(view, from_output, to_output);
}

void handle_view_request_geometry(wlc_handle view, const struct wlc_geometry *geometry) {
	struct wlc_geometry nc_geometry = {
		.origin = {
			.x = geometry->origin.x,
			.y = geometry->origin.y
		},
		.size = {
			.w = geometry->size.w,
			.h = geometry->size.h
		}
	};
	_goHandleViewRequestGeometry(view, &nc_geometry);
}

void handle_view_request_state(wlc_handle view, enum wlc_view_state_bit bit, bool toggle) {
	_goHandleViewRequestState(view, bit, toggle);
}

void handle_view_request_move(wlc_handle view, const struct wlc_point *point) {
	struct wlc_point nc_point = {
		.x = point->x,
		.y = point->y
	};
	_goHandleViewRequestMove(view, &nc_point);
}

void handle_view_request_resize(wlc_handle view, uint32_t edges, const struct wlc_point *point) {
	struct wlc_point nc_point = {
		.x = point->x,
		.y = point->y
	};
	_goHandleViewRequestResize(view, edges, &nc_point);
}

void handle_view_render_pre(wlc_handle view) {
	_goHandleViewRenderPre(view);
}

void handle_view_render_post(wlc_handle view) {
	_goHandleViewRenderPost(view);
}

void handle_view_properties_updated(wlc_handle view, uint32_t mask) {
	_goHandleViewPropertiesUpdated(view, mask);
}

/* keyboard */
bool handle_keyboard_key(wlc_handle view, uint32_t time, const struct wlc_modifiers *modifiers, uint32_t key, enum wlc_key_state state) {
	struct wlc_modifiers nc_modifiers = {
		.leds = modifiers->leds,
		.mods = modifiers->mods
	};
	return _goHandleKeyboardKey(view, time, &nc_modifiers, key, state);
}

/* pointer */
bool handle_pointer_button(wlc_handle view, uint32_t time, const struct wlc_modifiers *modifiers, uint32_t button, enum wlc_button_state state, const struct wlc_point *point) {
	struct wlc_modifiers nc_modifiers = {
		.leds = modifiers->leds,
		.mods = modifiers->mods
	};
	struct wlc_point nc_point = {
		.x = point->x,
		.y = point->y
	};
	return _goHandlePointerButton(view, time, &nc_modifiers, button, state, &nc_point);
}

bool handle_pointer_scroll(wlc_handle view, uint32_t time, const struct wlc_modifiers *modifiers, uint8_t axis_bits, double amount[2]) {
	struct wlc_modifiers nc_modifiers = {
		.leds = modifiers->leds,
		.mods = modifiers->mods
	};
	return _goHandlePointerScroll(view, time, &nc_modifiers, axis_bits, amount);
}

bool handle_pointer_motion(wlc_handle view, uint32_t time, const struct wlc_point *point) {
	struct wlc_point nc_point = {
		.x = point->x,
		.y = point->y
	};
	return _goHandlePointerMotion(view, time, &nc_point);
}

/* touch */
bool handle_touch_touch(wlc_handle view, uint32_t time, const struct wlc_modifiers *modifiers, enum wlc_touch_type touch, int32_t slot, const struct wlc_point *point) {
	struct wlc_modifiers nc_modifiers = {
		.leds = modifiers->leds,
		.mods = modifiers->mods
	};
	struct wlc_point nc_point = {
		.x = point->x,
		.y = point->y
	};
	return _goHandleTouchTouch(view, time, &nc_modifiers, touch, slot, &nc_point);
}

/* compositor */
void handle_compositor_ready(void) {
	_goHandleCompositorReady();
}

void handle_compositor_terminate(void) {
	_goHandleCompositorTerminate();
}

/* input */
bool handle_input_created(struct libinput_device *device) {
	return _goHandleInputCreated(device);
}

void handle_input_destroyed(struct libinput_device *device) {
	_goHandleInputDestroyed(device);
}



/* Callback wrappers */

void set_output_created_cb() {
	wlc_set_output_created_cb(handle_output_created);
}

void set_output_destroyed_cb() {
	wlc_set_output_destroyed_cb(handle_output_destroyed);
}

void set_output_focus_cb() {
	wlc_set_output_focus_cb(handle_output_focus);
}

void set_output_resolution_cb() {
	wlc_set_output_resolution_cb(handle_output_resolution);
}

void set_output_render_pre_cb() {
	wlc_set_output_render_pre_cb(handle_output_render_pre);
}

void set_output_render_post_cb() {
	wlc_set_output_render_post_cb(handle_output_render_post);
}

void set_output_context_created_cb() {
	wlc_set_output_context_created_cb(handle_output_context_created);
}

void set_output_context_destroyed_cb() {
	wlc_set_output_context_destroyed_cb(handle_output_context_destroyed);
}

void set_view_created_cb() {
	wlc_set_view_created_cb(handle_view_created);
}

void set_view_destroyed_cb() {
	wlc_set_view_destroyed_cb(handle_view_destroyed);
}

void set_view_focus_cb() {
	wlc_set_view_focus_cb(handle_view_focus);
}

void set_view_move_to_output_cb() {
	wlc_set_view_move_to_output_cb(handle_view_move_to_output);
}

void set_view_request_geometry_cb() {
	wlc_set_view_request_geometry_cb(handle_view_request_geometry);
}

void set_view_request_state_cb() {
	wlc_set_view_request_state_cb(handle_view_request_state);
}

void set_view_request_move_cb() {
	wlc_set_view_request_move_cb(handle_view_request_move);
}

void set_view_request_resize_cb() {
	wlc_set_view_request_resize_cb(handle_view_request_resize);
}

void set_view_render_pre_cb() {
	wlc_set_view_render_pre_cb(handle_view_render_pre);
}

void set_view_render_post_cb() {
	wlc_set_view_render_post_cb(handle_view_render_post);
}

void set_view_properties_updated_cb() {
	wlc_set_view_properties_updated_cb(handle_view_properties_updated);
}

void set_keyboard_key_cb() {
	wlc_set_keyboard_key_cb(handle_keyboard_key);
}

void set_pointer_button_cb() {
	wlc_set_pointer_button_cb(handle_pointer_button);
}

void set_pointer_scroll_cb() {
	wlc_set_pointer_scroll_cb(handle_pointer_scroll);
}

void set_pointer_motion_cb() {
	wlc_set_pointer_motion_cb(handle_pointer_motion);
}

void set_touch_cb() {
	wlc_set_touch_cb(handle_touch_touch);
}

void set_compositor_ready_cb() {
	wlc_set_compositor_ready_cb(handle_compositor_ready);
}

void set_compositor_terminate_cb() {
	wlc_set_compositor_terminate_cb(handle_compositor_terminate);
}

void set_input_created_cb() {
	wlc_set_input_created_cb(handle_input_created);
}

void set_input_destroyed_cb() {
	wlc_set_input_destroyed_cb(handle_input_destroyed);
}
