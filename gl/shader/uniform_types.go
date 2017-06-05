package shader

import (
	"fmt"
)

type UniformType uint32

//go:generate stringer -type=UniformType
const (
	TypeInvalid UniformType = iota
	TypeFloat
	TypeInt
	TypeVec2
	TypeVec3
)

var UniformTypeNames = map[UniformType]string{
	TypeInvalid: "invalid",
	TypeFloat:   "float",
	TypeInt:     "int",
	TypeVec2:    "vec2",
	TypeVec3:    "vec3",
}

type UniformValue interface {
	String() string
}

type Vec2 struct {
	value [2]float32
}

func (v *Vec2) String() string {
	return fmt.Sprintf("vec2(%.2f,%.2f)", v.value[0], v.value[1])
}

type Vec3 struct {
	value [3]float32
}

func (v *Vec3) String() string {
	return fmt.Sprintf("vec3(%.2f,%.2f,%.2f)", v.value[0], v.value[1], v.value[0])
}

type Float struct {
	value float32
}

func (v *Float) String() string {
	return fmt.Sprintf("%.2f", v.value)
}

type Int struct {
	value int32
}

func (v *Int) String() string {
	return fmt.Sprintf("%d", v.value)
}
