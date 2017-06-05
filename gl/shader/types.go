package shader

// EventType marks the event type received
type EventType uint32

type Event struct {
	Type EventType
	msg  string
}

const (
	_ EventType = iota
	EVENT_TYPE_RELOAD
)

// ShadelData represent the internal shader program
type ShadelData struct {
	Program  uint32
	reload   chan Event
	reloaded chan interface{}
	uniforms map[string]Uniform
	err      bool
}

type uniformData struct {
	name     string
	t        UniformType
	value    interface{}
	location int32
}

// Uniform represents a typed shader uniform found by introspection of the shader
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
	SetUniform3f(string, float32, float32, float32)
	GetProgram() uint32
	ShaderChanged() <-chan interface{}
}
