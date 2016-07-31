package gl

type EventType uint32

type Event struct {
	Type EventType
	msg  string
}

const (
	_ EventType = iota
	EVENT_TYPE_RELOAD
)

type ShadelData struct {
	Program  uint32
	reload   chan Event
	uniforms map[string]Uniform
}

type UniformType uint32

//go:generate stringer -type=UniformType
const (
	TypeInvalid UniformType = iota
	TypeFloat
	TypeInt
	TypeVec2
	TypeVec3
)

type uniformData struct {
	name     string
	t        UniformType
	value    interface{}
	location int32
}

type Uniform interface {
	Name() string
	Type() UniformType
	SetValue(interface{}) error
	GetValue() interface{}
	ValueString() string
}

type Shadel interface {
	ReplaceShadel(string, string) error
	Use()
	//	Uniforms() map[string]Uniform
	GetProgram() uint32
}

type Vec2 struct {
	x, y float32
}
type Vec3 struct {
	Vec2
	z float32
}
