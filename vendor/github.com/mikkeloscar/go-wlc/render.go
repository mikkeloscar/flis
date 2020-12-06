package wlc

/*
#cgo LDFLAGS: -lwlc
#include <stdlib.h>
#include <wlc/wlc-render.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// PixelFormat describes the pixelformat used when writing/reading pixels.
type PixelFormat C.enum_wlc_pixel_format

const (
	// RGBA8888 defines a color format where each channel is 8 bits.
	RGBA8888 PixelFormat = iota
)

// PixelsWrite write pixel data with the specific format to outputs
// framebuffer. If geometry is out of bounds, it will be automatically clamped.
// TODO: make more go friendly
func PixelsWrite(format PixelFormat, geometry Geometry, data unsafe.Pointer) {
	cgeometry := geometry.c()
	defer C.free(unsafe.Pointer(cgeometry))
	C.wlc_pixels_write(C.enum_wlc_pixel_format(format), cgeometry, data)
}

// PixelsRead read pixel data from output's framebuffer.
// If the geometry is out of bounds, it will be automatically clamped.
// Potentially clamped geometry will be stored in out_geometry, to indicate
// width / height of the returned data.
// TODO: make more go friendly
func PixelsRead(format PixelFormat, geometry Geometry, outGeometry *Geometry, outData unsafe.Pointer) {
	cgeometry := geometry.c()
	defer C.free(unsafe.Pointer(cgeometry))
	var cgOut C.struct_wlc_geometry
	C.wlc_pixels_read(C.enum_wlc_pixel_format(format), cgeometry, &cgOut, outData)
	geometryCtoGo(outGeometry, &cgOut)
}

// Render renders surfaces inside post / pre render hooks.
func (s Resource) Render(geometry Geometry) {
	cgeometry := geometry.c()
	defer C.free(unsafe.Pointer(cgeometry))
	C.wlc_surface_render(C.wlc_resource(s), cgeometry)
}

// ScheduleRender schedules output for rendering next frame.
// If output was already scheduled this is no-op, if output is currently
// rendering, it will render immediately after.
func (o Output) ScheduleRender() {
	C.wlc_output_schedule_render(C.wlc_handle(o))
}

// FlushFrameCallbacks adds frame callbacks of the given surface for the next
// output frame.  It applies recursively to all subsurfaces.  Useful when the
// compositor creates custom animations which require disabling internal
// rendering, but still need to update the surface textures (for ex. video
// players).
func (s Resource) FlushFrameCallbacks() {
	C.wlc_surface_flush_frame_callbacks(C.wlc_resource(s))
}

// Renderer defines enabled renderers.
type Renderer C.enum_wlc_renderer

const (
	// RendererGLES2 defines a GLES2 renderer.
	RendererGLES2 Renderer = iota
	// NoRenderer defines no renderer.
	NoRenderer
)

// GetRenderer returns currently active renderer on the given output.
func (o Output) GetRenderer() Renderer {
	return Renderer(C.wlc_output_get_renderer(C.wlc_handle(o)))
}

// SurfaceFormat defines the format returned by GetTextures.
type SurfaceFormat C.enum_wlc_surface_format

const (
	// SurfaceRGB defines surface format RGB.
	SurfaceRGB SurfaceFormat = iota
	// SurfaceRGBA defines surface format RGBA.
	SurfaceRGBA
	// SurfaceEGL defines surface format EGL.
	SurfaceEGL
	// SurfaceYUv defines surface format yuv.
	SurfaceYUv
	// SurfaceYUV defines surface format YUV.
	SurfaceYUV
	// SurfaceYXUXV defines surface format YXUXV.
	SurfaceYXUXV
)

// GetTextures returns an array with the textures of a surface. Returns error
// if surface is invalid. Note that these are not only OpenGL textures but
// rather render-specific.
func (s Resource) GetTextures() ([3]uint32, SurfaceFormat, error) {
	var outTextures [3]C.uint32_t
	var format C.enum_wlc_surface_format
	val := bool(C.wlc_surface_get_textures(C.wlc_resource(s), &outTextures[0], &format))
	if val {
		var textures [3]uint32
		for i, t := range outTextures {
			textures[i] = uint32(t)
		}
		return textures, SurfaceFormat(format), nil
	}

	return [3]uint32{}, 0, fmt.Errorf("invalid surface")
}
