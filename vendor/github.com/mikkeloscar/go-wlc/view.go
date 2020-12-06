package wlc

/*
#cgo LDFLAGS: -lwlc
#include <stdlib.h>
#include <wlc/wlc.h>
*/
import "C"

import "unsafe"

// View is a wlc_handle describing a view object in wlc.
type View C.wlc_handle

// Focus focuses view.
func (v View) Focus() {
	C.wlc_view_focus(C.wlc_handle(v))
}

// ViewUnfocus unfocuses all views.
func ViewUnfocus() {
	C.wlc_view_focus(0)
}

// Close closes view.
func (v View) Close() {
	C.wlc_view_close(C.wlc_handle(v))
}

// GetOutput gets output of view.
func (v View) GetOutput() Output {
	return Output(C.wlc_view_get_output(C.wlc_handle(v)))
}

// SetOutput sets output for view. Alternatively output.SetViews() can be used.
func (v View) SetOutput(output Output) {
	C.wlc_view_set_output(C.wlc_handle(v), C.wlc_handle(output))
}

// SendToBack sends view behind everything.
func (v View) SendToBack() {
	C.wlc_view_send_to_back(C.wlc_handle(v))
}

// SendBelow sends view below another view.
func (v View) SendBelow(other View) {
	C.wlc_view_send_below(C.wlc_handle(v), C.wlc_handle(other))
}

// BringAbove brings view above another view.
func (v View) BringAbove(other View) {
	C.wlc_view_bring_above(C.wlc_handle(v), C.wlc_handle(other))
}

// BringToFront brings view to front of everything.
func (v View) BringToFront() {
	C.wlc_view_bring_to_front(C.wlc_handle(v))
}

// GetMask gets current visibility bitmask.
func (v View) GetMask() uint32 {
	return uint32(C.wlc_view_get_mask(C.wlc_handle(v)))
}

// SetMask sets visibility bitmask.
func (v View) SetMask(mask uint32) {
	C.wlc_view_set_mask(C.wlc_handle(v), C.uint32_t(mask))
}

// GetGeometry gets current geometry (what the client sees).
func (v View) GetGeometry() *Geometry {
	cgeometry := C.wlc_view_get_geometry(C.wlc_handle(v))
	return geometryCtoGo(&Geometry{}, cgeometry)
}

// PositionerGetSize gets size requested by positioner, as defined in xdg-shell
// v6.
func (v View) PositionerGetSize() *Size {
	csize := C.wlc_view_positioner_get_size(C.wlc_handle(v))
	return sizeCtoGo(csize)
}

// PositionerGetAnchorRect gets anchor rectangle requested by positioner, as
// defined in xdg-shell v6.
// Returns nil if view has no valid positioner.
func (v View) PositionerGetAnchorRect() *Geometry {
	cgeometry := C.wlc_view_positioner_get_anchor_rect(C.wlc_handle(v))
	return geometryCtoGo(&Geometry{}, cgeometry)
}

// PositionerGetOffset gets offset requested by positioner, as defined in
// xdg-shell v6.
// Returns NULL if view has no valid positioner, or default value (0, 0) if
// positioner has no offset set.
func (v View) PositionerGetOffset() *Point {
	cpoint := C.wlc_view_positioner_get_offset(C.wlc_handle(v))
	return pointCtoGo(cpoint)
}

// PositionerGetAnchor gets anchor requested by positioner, as defined in
// xdg-shell v6.
// Returns default value WLC_BIT_ANCHOR_NONE if view has no valid positioner or
// if positioner has no anchor set.
func (v View) PositionerGetAnchor() PositionerAnchorBit {
	return PositionerAnchorBit(C.wlc_view_positioner_get_anchor(C.wlc_handle(v)))
}

// PositionerGetGravity gets anchor requested by positioner, as defined in
// xdg-shell v6.
// Returns default value WLC_BIT_GRAVITY_NONE if view has no valid positioner
// or if positioner has no gravity set.
func (v View) PositionerGetGravity() PositionerGravityBit {
	return PositionerGravityBit(C.wlc_view_positioner_get_gravity(C.wlc_handle(v)))
}

// PositionerGetConstraintAdjustment gets constraint adjustment requested by
// positioner, as defined in xdg-shell v6.
// Returns default value WLC_BIT_CONSTRAINT_ADJUSTMENT_NONE if view has no
// valid positioner or if positioner has no constraint adjustment set.
func (v View) PositionerGetConstraintAdjustment() PositionerConstraintAdjustmentBit {
	return PositionerConstraintAdjustmentBit(
		C.wlc_view_positioner_get_constraint_adjustment(C.wlc_handle(v)),
	)
}

// GetVisibleGeometry gets current visible geometry (what wlc displays).
func (v View) GetVisibleGeometry() Geometry {
	cgeometry := C.struct_wlc_geometry{}
	C.wlc_view_get_visible_geometry(C.wlc_handle(v), &cgeometry)
	return *geometryCtoGo(&Geometry{}, &cgeometry)
}

// SetGeometry sets geometry. Set edges if the geometry change is caused by
// interactive resize.
func (v View) SetGeometry(edges uint32, geometry Geometry) {
	cgeometry := geometry.c()
	defer C.free(unsafe.Pointer(cgeometry))
	C.wlc_view_set_geometry(C.wlc_handle(v), C.uint32_t(edges), cgeometry)
}

// GetType gets type bitfield for view.
func (v View) GetType() uint32 {
	return uint32(C.wlc_view_get_type(C.wlc_handle(v)))
}

// SetType sets type bit. TOggle indicates whether it is set or not.
func (v View) SetType(typ ViewTypeBit, toggle bool) {
	C.wlc_view_set_type(C.wlc_handle(v), uint32(typ), C._Bool(toggle))
}

// GetState gets current state bitfield.
func (v View) GetState() uint32 {
	return uint32(C.wlc_view_get_state(C.wlc_handle(v)))
}

// SetState sets state bit. Toggle indicates whether it is set or not.
func (v View) SetState(state ViewStateBit, toggle bool) {
	C.wlc_view_set_state(C.wlc_handle(v), uint32(state), C._Bool(toggle))
}

// GetParent gets parent view.
func (v View) GetParent() View {
	return View(C.wlc_view_get_parent(C.wlc_handle(v)))
}

// SetParent sets parent view.
func (v View) SetParent(parent View) {
	C.wlc_view_set_parent(C.wlc_handle(v), C.wlc_handle(parent))
}

// Title gets title.
func (v View) Title() string {
	ctitle := C.wlc_view_get_title(C.wlc_handle(v))
	return C.GoString(ctitle)
}

// Instance gets instance (shell-surface only).
func (v View) Instance() string {
	cinstance := C.wlc_view_get_instance(C.wlc_handle(v))
	return C.GoString(cinstance)
}

// GetClass gets class. (shell-surface only).
func (v View) GetClass() string {
	cclass := C.wlc_view_get_class(C.wlc_handle(v))
	return C.GoString(cclass)
}

// GetAppID gets app id. (xdg-surface only).
func (v View) GetAppID() string {
	capp := C.wlc_view_get_app_id(C.wlc_handle(v))
	return C.GoString(capp)
}

// GetPID gets pid of program owning the view.
func (v View) GetPID() int {
	return int(C.wlc_view_get_pid(C.wlc_handle(v)))
}
