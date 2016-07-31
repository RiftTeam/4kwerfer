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
	Program uint32
	reload  chan Event
}

type UniformType uint32

const (
	TypeInvalid = UniformType(iota)
	TypeFloat
	TypeInt
	TypeVec2
	TypeVec3
)

type Uniform struct {
	name string
	t    UniformType
}

type Shadel interface {
	ReplaceShadel(string, string) error
	Use()
	Uniforms() []*Uniform
	GetProgram() uint32
}
