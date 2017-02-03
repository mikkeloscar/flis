package wlc

/*
#cgo LDFLAGS: -lwlc
#include <wlc/wlc.h>

struct wlc_point *init_point(int32_t x, int32_t y) {
	struct wlc_point *point = malloc(sizeof(struct wlc_point));
	point->x = x;
	point->y = y;
	return point;
}

struct wlc_size *init_size(uint32_t w, uint32_t h) {
	struct wlc_size *size = malloc(sizeof(struct wlc_size));
	size->w = w;
	size->h = h;
	return size;
}

struct wlc_geometry *init_geometry(int32_t x, int32_t y, uint32_t w, uint32_t h) {
	struct wlc_geometry *geometry = malloc(sizeof(struct wlc_geometry));
	geometry->origin.x = x;
	geometry->origin.y = y;
	geometry->size.w = w;
	geometry->size.h = h;
	return geometry;
}
*/
import "C"
import "math"

// Point is a fixed 2D point.
type Point struct {
	X, Y int32
}

func (p *Point) c() *C.struct_wlc_point {
	return C.init_point(C.int32_t(p.X), C.int32_t(p.Y))
}

func pointCtoGo(c *C.struct_wlc_point) *Point {
	if c != nil {
		return &Point{
			X: int32((*c).x),
			Y: int32((*c).y),
		}
	}

	return nil
}

// Size is a fixed 2D size.
type Size struct {
	W, H uint32
}

func (s *Size) c() *C.struct_wlc_size {
	return C.init_size(C.uint32_t(s.W), C.uint32_t(s.H))
}

func sizeCtoGo(c *C.struct_wlc_size) *Size {
	if c != nil {
		return &Size{
			W: uint32((*c).w),
			H: uint32((*c).h),
		}
	}

	return nil
}

// Geometry is a fixed 2D point, size pair.
type Geometry struct {
	Origin Point
	Size   Size
}

func (g *Geometry) c() *C.struct_wlc_geometry {
	return C.init_geometry(
		C.int32_t(g.Origin.X),
		C.int32_t(g.Origin.Y),
		C.uint32_t(g.Size.W),
		C.uint32_t(g.Size.H),
	)
}

func geometryCtoGo(g *Geometry, c *C.struct_wlc_geometry) *Geometry {
	if c != nil {
		g.Origin = *pointCtoGo((*C.struct_wlc_point)(&(*c).origin))
		g.Size = *sizeCtoGo((*C.struct_wlc_size)(&(*c).size))
		return g
	}

	return nil
}

var (
	// PointZero defines a point at (0,0).
	PointZero = Point{0, 0}
	// SizeZero defines a size 0x0.
	SizeZero = Size{0, 0}
	// GeometryZero defines a geometry with size 0x0 at point (0,0).
	GeometryZero = Geometry{Point{0, 0}, Size{0, 0}}
)

// PointMin returns the smallest values of two points.
func PointMin(a, b Point) Point {
	return Point{
		X: int32(math.Min(float64(a.X), float64(b.X))),
		Y: int32(math.Min(float64(a.Y), float64(b.Y))),
	}
}

// PointMax returns the biggest values of two points.
func PointMax(a, b Point) Point {
	return Point{
		X: int32(math.Max(float64(a.X), float64(b.X))),
		Y: int32(math.Max(float64(a.Y), float64(b.Y))),
	}
}

// SizeMin returns the smallest values of two sizes.
func SizeMin(a, b Size) Size {
	return Size{
		W: uint32(math.Min(float64(a.W), float64(b.W))),
		H: uint32(math.Min(float64(a.H), float64(b.H))),
	}
}

// SizeMax returns the biggest values of two sizes.
func SizeMax(a, b Size) Size {
	return Size{
		W: uint32(math.Max(float64(a.W), float64(b.W))),
		H: uint32(math.Max(float64(a.H), float64(b.H))),
	}
}

// PointEquals compares two points.
func PointEquals(a, b Point) bool {
	return a.X == b.X && a.Y == b.Y
}

// SizeEquals compares two sizes.
func SizeEquals(a, b Size) bool {
	return a.W == b.W && a.H == b.H
}

// GeometryEquals compares two geometries.
func GeometryEquals(a, b Geometry) bool {
	return PointEquals(a.Origin, b.Origin) && SizeEquals(a.Size, b.Size)
}

// GeometryContains check if b is contained in a.
func GeometryContains(a, b Geometry) bool {
	return a.Origin.X <= b.Origin.X &&
		a.Origin.Y <= b.Origin.Y &&
		a.Origin.X+int32(a.Size.W) >= b.Origin.X+int32(b.Size.W) &&
		a.Origin.Y+int32(a.Size.H) >= b.Origin.Y+int32(b.Size.H)
}
