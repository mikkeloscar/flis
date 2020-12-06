package wlc

/*
#cgo LDFLAGS: -lwlc
#include <wlc/wlc.h>

// output
extern bool handle_output_created(wlc_handle output);
extern void handle_output_destroyed(wlc_handle output);
extern void handle_output_focus(wlc_handle output, bool focus);
extern void handle_output_resolution(wlc_handle output, const struct wlc_size *from, const struct wlc_size *to);
extern void handle_output_pre_render(wlc_handle output);
extern void handle_output_post_render(wlc_handle output);
extern void handle_output_context_created(wlc_handle output);
extern void handle_output_context_destroyed(wlc_handle output);
// view
extern bool handle_view_created(wlc_handle view);
extern void handle_view_destroyed(wlc_handle view);
extern void handle_view_focus(wlc_handle view, bool focus);
extern void handle_view_move_to_output(wlc_handle view, wlc_handle from_output, wlc_handle to_output);
extern void handle_view_geometry_request(wlc_handle view, const struct wlc_geometry*);
extern void handle_view_state_request(wlc_handle view, enum wlc_view_state_bit, bool toggle);
extern void handle_view_move_request(wlc_handle view, const struct wlc_point*);
extern void handle_view_resize_request(wlc_handle view, uint32_t edges, const struct wlc_point*);
extern void handle_view_pre_render(wlc_handle view);
extern void handle_view_post_render(wlc_handle view);
extern void handle_view_properties_updated(wlc_handle view, uint32_t mask);
// keyboard
extern bool handle_keyboard_key(wlc_handle view, uint32_t time, const struct wlc_modifiers*, uint32_t key, enum wlc_key_state);
// pointer
extern bool handle_pointer_button(wlc_handle view, uint32_t time, const struct wlc_modifiers*, uint32_t button, enum wlc_button_state, const struct wlc_point*);
extern bool handle_pointer_scroll(wlc_handle view, uint32_t time, const struct wlc_modifiers*, uint8_t axis_bits, double amount[2]);
extern bool handle_pointer_motion(wlc_handle view, uint32_t time, const struct wlc_point*);
// touch
extern bool handle_touch_touch(wlc_handle view, uint32_t time, const struct wlc_modifiers*, enum wlc_touch_type, int32_t slot, const struct wlc_point*);
// compositor
extern void handle_compositor_ready(void);
extern void handle_compositor_terminate(void);
// input
extern bool handle_input_created(struct libinput_device *device);
extern void handle_input_destroyed(struct libinput_device *device);

// callback wrappers
void set_output_created_cb();
void set_output_destroyed_cb();
void set_output_focus_cb();
void set_output_resolution_cb();
void set_output_render_pre_cb();
void set_output_render_post_cb();
void set_output_context_created_cb();
void set_output_context_destroyed_cb();
void set_view_created_cb();
void set_view_destroyed_cb();
void set_view_focus_cb();
void set_view_move_to_output_cb();
void set_view_request_geometry_cb();
void set_view_request_state_cb();
void set_view_request_move_cb();
void set_view_request_resize_cb();
void set_view_render_pre_cb();
void set_view_render_post_cb();
void set_view_properties_updated_cb();
void set_keyboard_key_cb();
void set_pointer_button_cb();
void set_pointer_scroll_cb();
void set_pointer_motion_cb();
void set_touch_cb();
void set_compositor_ready_cb();
void set_compositor_terminate_cb();
void set_input_created_cb();
void set_input_destroyed_cb();
*/
import "C"
import "unsafe"

var wlcInterface internalInterface

// Interface is used for commication with wlc.
type internalInterface struct {
	Output struct {
		Created    func(Output) bool
		Destroyed  func(Output)
		Focus      func(Output, bool)
		Resolution func(Output, *Size, *Size)
		Render     struct {
			Pre  func(Output)
			Post func(Output)
		}
		Context struct {
			Created   func(Output)
			Destroyed func(Output)
		}
	}
	View struct {
		Created      func(View) bool
		Destroyed    func(View)
		Focus        func(View, bool)
		MoveToOutput func(View, Output, Output)
		Request      struct {
			Geometry func(View, *Geometry)
			State    func(View, ViewStateBit, bool)
			Move     func(View, *Point)
			Resize   func(View, uint32, *Point)
		}
		Render struct {
			Pre  func(View)
			Post func(View)
		}
		PropertiesUpdated func(View, ViewPropertyUpdateBit)
	}
	Keyboard struct {
		Key func(View, uint32, Modifiers, uint32, KeyState) bool
	}
	Pointer struct {
		Button func(View, uint32, Modifiers, uint32, ButtonState, *Point) bool
		Scroll func(View, uint32, Modifiers, uint8, [2]float64) bool
		Motion func(View, uint32, *Point) bool
	}
	Touch struct {
		Touch func(View, uint32, Modifiers, TouchType, int32, *Point) bool
	}
	Compositor struct {
		Ready     func()
		Terminate func()
	}
	Input struct {
		Created   func(*C.struct_libinput_device) bool
		Destroyed func(*C.struct_libinput_device)
	}
}

// SetOutputCreatedCb sets callback to trigger when output is created. Callback
// should return false if you want to destroy the output. (e.g. failed to
// allocate data related to view)
func SetOutputCreatedCb(cb func(Output) bool) {
	wlcInterface.Output.Created = cb
	C.set_output_created_cb()
}

// SetOutputDestroyedCb sets callback to trigger when output is destroyed.
func SetOutputDestroyedCb(cb func(Output)) {
	wlcInterface.Output.Destroyed = cb
	C.set_output_destroyed_cb()
}

// SetOutputFocusCb sets callback to trigger when output got or lost focus.
func SetOutputFocusCb(fn func(Output, bool)) {
	wlcInterface.Output.Focus = fn
	C.set_output_focus_cb()
}

// SetOutputResolutionCb sets callback to trigger when output resolution
// changed.
func SetOutputResolutionCb(cb func(Output, *Size, *Size)) {
	wlcInterface.Output.Resolution = cb
	C.set_output_resolution_cb()
}

// SetOutputRenderPreCb sets the pre render hook for output.
func SetOutputRenderPreCb(cb func(Output)) {
	wlcInterface.Output.Render.Pre = cb
	C.set_output_render_pre_cb()
}

// SetOutputRenderPostCb sets the post render hook for output.
func SetOutputRenderPostCb(cb func(Output)) {
	wlcInterface.Output.Render.Post = cb
	C.set_output_render_post_cb()
}

// SetOutputContextCreated sets callback to trigger when output context is
// created. This generally happens on startup and when current tty changes.
func SetOutputContextCreated(cb func(Output)) {
	wlcInterface.Output.Context.Created = cb
	C.set_output_context_created_cb()
}

// SetOutputContextDestroyed sets callback to trigger when output context is
// destroyed.
func SetOutputContextDestroyed(cb func(Output)) {
	wlcInterface.Output.Context.Destroyed = cb
	C.set_output_context_destroyed_cb()
}

// SetViewCreatedCb sets callback to trigger when view is created. Callback
// should return false if you want to destroy the view. (e.g. failed to
// allocate data related to view).
func SetViewCreatedCb(cb func(View) bool) {
	wlcInterface.View.Created = cb
	C.set_view_created_cb()
}

// SetViewDestroyedCb sets callback to trigger when view is destroyed.
func SetViewDestroyedCb(cb func(View)) {
	wlcInterface.View.Destroyed = cb
	C.set_view_destroyed_cb()
}

// SetViewFocusCb sets callback to trigger when view got or lost focus.
func SetViewFocusCb(cb func(View, bool)) {
	wlcInterface.View.Focus = cb
	C.set_view_focus_cb()
}

// SetViewMoveToOutputCb sets callback to trigger when view is moved to an
// output.
func SetViewMoveToOutputCb(cb func(View, Output, Output)) {
	wlcInterface.View.MoveToOutput = cb
	C.set_view_move_to_output_cb()
}

// SetViewRequestGeometryCb sets callback to trigger when a view requests to
// set geometry. Apply using View.SetGeometry to agree.
func SetViewRequestGeometryCb(cb func(View, *Geometry)) {
	wlcInterface.View.Request.Geometry = cb
	C.set_view_request_geometry_cb()
}

// SetViewRequestStateCb sets callback to trigger when a view requests to
// disable or enable the given state. Apply using View.SetState to agree.
func SetViewRequestStateCb(cb func(View, ViewStateBit, bool)) {
	wlcInterface.View.Request.State = cb
	C.set_view_request_state_cb()
}

// SetViewRequestMoveCb sets callback to trigger when view requests to move
// itself. Start an interactive move to agree.
func SetViewRequestMoveCb(cb func(View, *Point)) {
	wlcInterface.View.Request.Move = cb
	C.set_view_request_move_cb()
}

// SetViewRequestResizeCb sets callback to trigger when view requests to resize
// iteself with the given edge. Start an interactive resize to agree.
func SetViewRequestResizeCb(cb func(View, uint32, *Point)) {
	wlcInterface.View.Request.Resize = cb
	C.set_view_request_resize_cb()
}

// SetViewRenderPreCb sets the pre render hook for view.
func SetViewRenderPreCb(cb func(View)) {
	wlcInterface.View.Render.Pre = cb
	C.set_view_render_pre_cb()
}

// SetViewRenderPostCb sets the post render hook for view.
func SetViewRenderPostCb(cb func(View)) {
	wlcInterface.View.Render.Post = cb
	C.set_view_render_post_cb()
}

// SetViewPropertiesUpdatedCb sets callback to trigger when view properties are
// updated.
func SetViewPropertiesUpdatedCb(cb func(View, ViewPropertyUpdateBit)) {
	wlcInterface.View.PropertiesUpdated = cb
	C.set_view_properties_updated_cb()
}

// SetKeyboardKeyCb sets callback to trigger when key event was triggered, view
// handle will be zero if there was no focus. Callback can return true to
// prevent sending the event to clients.
func SetKeyboardKeyCb(cb func(View, uint32, Modifiers, uint32, KeyState) bool) {
	wlcInterface.Keyboard.Key = cb
	C.set_keyboard_key_cb()
}

// SetPointerButtonCb sets callback to trigger when button event was triggered,
// view handle will be zero if there was no focus. Callback can return true
// to prevent sending the event to clients.
func SetPointerButtonCb(cb func(View, uint32, Modifiers, uint32, ButtonState, *Point) bool) {
	wlcInterface.Pointer.Button = cb
	C.set_pointer_button_cb()
}

// SetPointerScrollCb sets callback to trigger when scroll event was triggered,
// view handle will be zero if there was no focus. Callback can return true
// to prevent sending the event to clients.
func SetPointerScrollCb(cb func(View, uint32, Modifiers, uint8, [2]float64) bool) {
	wlcInterface.Pointer.Scroll = cb
	C.set_pointer_scroll_cb()
}

// SetPointerMotionCb sets callback to trigger when motion event was triggered,
// view handle will be zero if there was no focus. Apply with
// wlc_pointer_set_position to agree. Callback can return true to prevent
// sending the event to clients.
func SetPointerMotionCb(cb func(View, uint32, *Point) bool) {
	wlcInterface.Pointer.Motion = cb
	C.set_pointer_motion_cb()
}

// SetTouchCb sets callback to trigger when touch event was triggered, view
// handle will be zero if there was no focus. Callback can return true to
// prevent sending the event to clients.
func SetTouchCb(cb func(View, uint32, Modifiers, TouchType, int32, *Point) bool) {
	wlcInterface.Touch.Touch = cb
	C.set_touch_cb()
}

// SetCompositorReadyCb sets callback to trigger when compositor is ready to
// accept clients.
func SetCompositorReadyCb(cb func()) {
	wlcInterface.Compositor.Ready = cb
	C.set_compositor_ready_cb()
}

// SetCompositorTerminateCb sets callback to trigger when compositor is about
// to terminate.
func SetCompositorTerminateCb(cb func()) {
	wlcInterface.Compositor.Terminate = cb
	C.set_compositor_terminate_cb()
}

// SetInputCreatedCb sets callback to trigger when input device is created.
// Return value of callback does nothing. (Experimental).
func SetInputCreatedCb(cb func(*C.struct_libinput_device) bool) {
	wlcInterface.Input.Created = cb
	C.set_input_created_cb()
}

// SetInputDestroyedCb sets callback to trigger when input device was
// destroyed. (Experimental).
func SetInputDestroyedCb(cb func(*C.struct_libinput_device)) {
	wlcInterface.Input.Destroyed = cb
	C.set_input_destroyed_cb()
}

// output wrappers

//export _goHandleOutputCreated
func _goHandleOutputCreated(output C.wlc_handle) C._Bool {
	return C._Bool(wlcInterface.Output.Created(Output(output)))
}

//export _goHandleOutputDestroyed
func _goHandleOutputDestroyed(output C.wlc_handle) {
	wlcInterface.Output.Destroyed(Output(output))
}

//export _goHandleOutputFocus
func _goHandleOutputFocus(output C.wlc_handle, focus bool) {
	wlcInterface.Output.Focus(Output(output), focus)
}

//export _goHandleOutputResolution
func _goHandleOutputResolution(output C.wlc_handle, from *C.struct_wlc_size, to *C.struct_wlc_size) {
	wlcInterface.Output.Resolution(Output(output), sizeCtoGo(from), sizeCtoGo(to))
}

//export _goHandleOutputRenderPre
func _goHandleOutputRenderPre(output C.wlc_handle) {
	wlcInterface.Output.Render.Pre(Output(output))
}

//export _goHandleOutputRenderPost
func _goHandleOutputRenderPost(output C.wlc_handle) {
	wlcInterface.Output.Render.Post(Output(output))
}

//export _goHandleOutputContextCreated
func _goHandleOutputContextCreated(output C.wlc_handle) {
	wlcInterface.Output.Context.Created(Output(output))
}

//export _goHandleOutputContextDestroyed
func _goHandleOutputContextDestroyed(output C.wlc_handle) {
	wlcInterface.Output.Context.Destroyed(Output(output))
}

// view wrappers

//export _goHandleViewCreated
func _goHandleViewCreated(view C.wlc_handle) C._Bool {
	return C._Bool(wlcInterface.View.Created(View(view)))
}

//export _goHandleViewDestroyed
func _goHandleViewDestroyed(view C.wlc_handle) {
	wlcInterface.View.Destroyed(View(view))
}

//export _goHandleViewFocus
func _goHandleViewFocus(view C.wlc_handle, focus bool) {
	wlcInterface.View.Focus(View(view), focus)
}

//export _goHandleViewMoveToOutput
func _goHandleViewMoveToOutput(view C.wlc_handle, fromOutput C.wlc_handle, toOutput C.wlc_handle) {
	wlcInterface.View.MoveToOutput(View(view), Output(fromOutput), Output(toOutput))
}

//export _goHandleViewRequestGeometry
func _goHandleViewRequestGeometry(view C.wlc_handle, geometry *C.struct_wlc_geometry) {
	wlcInterface.View.Request.Geometry(View(view), geometryCtoGo(&Geometry{}, geometry))
}

//export _goHandleViewRequestState
func _goHandleViewRequestState(view C.wlc_handle, state C.enum_wlc_view_state_bit, toggle bool) {
	wlcInterface.View.Request.State(View(view), ViewStateBit(state), toggle)
}

//export _goHandleViewRequestMove
func _goHandleViewRequestMove(view C.wlc_handle, point *C.struct_wlc_point) {
	wlcInterface.View.Request.Move(View(view), pointCtoGo(point))
}

//export _goHandleViewRequestResize
func _goHandleViewRequestResize(view C.wlc_handle, edges C.uint32_t, point *C.struct_wlc_point) {
	wlcInterface.View.Request.Resize(View(view), uint32(edges), pointCtoGo(point))
}

//export _goHandleViewRenderPre
func _goHandleViewRenderPre(view C.wlc_handle) {
	wlcInterface.View.Render.Pre(View(view))
}

//export _goHandleViewRenderPost
func _goHandleViewRenderPost(view C.wlc_handle) {
	wlcInterface.View.Render.Post(View(view))
}

//export _goHandleViewPropertiesUpdated
func _goHandleViewPropertiesUpdated(view C.wlc_handle, mask C.uint32_t) {
	wlcInterface.View.PropertiesUpdated(View(view), ViewPropertyUpdateBit(mask))
}

// keyboard wrapper

//export _goHandleKeyboardKey
func _goHandleKeyboardKey(view C.wlc_handle, time C.uint32_t, modifiers *C.struct_wlc_modifiers, key C.uint32_t, state C.enum_wlc_key_state) C._Bool {
	return C._Bool(wlcInterface.Keyboard.Key(
		View(view),
		uint32(time),
		modsCtoGo(modifiers),
		uint32(key),
		KeyState(state),
	))
}

// pointer wrapper

//export _goHandlePointerButton
func _goHandlePointerButton(view C.wlc_handle, time C.uint32_t, modifiers *C.struct_wlc_modifiers, button C.uint32_t, state C.enum_wlc_button_state, point *C.struct_wlc_point) C._Bool {
	return C._Bool(wlcInterface.Pointer.Button(
		View(view),
		uint32(time),
		modsCtoGo(modifiers),
		uint32(button),
		ButtonState(state),
		pointCtoGo(point),
	))
}

//export _goHandlePointerScroll
func _goHandlePointerScroll(view C.wlc_handle, time C.uint32_t, modifiers *C.struct_wlc_modifiers, axisBits C.uint8_t, amount *C.double) C._Bool {
	// convert double[2] to [2]float64
	goAmount := [2]float64{
		*(*float64)(amount),
		*(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(amount)) + unsafe.Sizeof(*amount))),
	}
	return C._Bool(wlcInterface.Pointer.Scroll(
		View(view),
		uint32(time),
		modsCtoGo(modifiers),
		uint8(axisBits),
		goAmount,
	))
}

//export _goHandlePointerMotion
func _goHandlePointerMotion(view C.wlc_handle, time C.uint32_t, point *C.struct_wlc_point) C._Bool {
	return C._Bool(wlcInterface.Pointer.Motion(
		View(view),
		uint32(time),
		pointCtoGo(point),
	))
}

// touch wrapper

//export _goHandleTouchTouch
func _goHandleTouchTouch(view C.wlc_handle, time C.uint32_t, modifiers *C.struct_wlc_modifiers, touch C.enum_wlc_touch_type, slot C.int32_t, point *C.struct_wlc_point) C._Bool {
	return C._Bool(wlcInterface.Touch.Touch(
		View(view),
		uint32(time),
		modsCtoGo(modifiers),
		TouchType(touch),
		int32(slot),
		pointCtoGo(point),
	))
}

// compositor wrapper

//export _goHandleCompositorReady
func _goHandleCompositorReady() {
	wlcInterface.Compositor.Ready()
}

//export _goHandleCompositorTerminate
func _goHandleCompositorTerminate() {
	wlcInterface.Compositor.Terminate()
}

// input wrapper

//export _goHandleInputCreated
func _goHandleInputCreated(device *C.struct_libinput_device) C._Bool {
	return C._Bool(wlcInterface.Input.Created(device))
}

//export _goHandleInputDestroyed
func _goHandleInputDestroyed(device *C.struct_libinput_device) {
	wlcInterface.Input.Destroyed(device)
}
