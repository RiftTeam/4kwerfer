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
	Apply() error
}

type Shadel interface {
	ReplaceShadel(string, string) error
	Use()
	//	Uniforms() map[string]Uniform
	GetProgram() uint32
}
