package gl

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"log"
	"strings"
)

const (
	VSH_FILE = "shader.vsh"
	FSH_FILE = "shader.fsh"
)

func (s *ShadelData) ReplaceShadel(vshFile, fshFile string) error {
	vsh, _ := ioutil.ReadFile("shader.vsh")
	fsh, _ := ioutil.ReadFile("shader.fsh")

	defer gl.DeleteProgram(s.Program)
	return s.newProgram(string(vsh)+"\x00", string(fsh)+"\x00")
}

func (s *ShadelData) Use() {
	select {
	case ev := <-s.reload:
		log.Printf("reloading %s", ev)
		s.ReplaceShadel(VSH_FILE, FSH_FILE)

	default:
		gl.UseProgram(s.Program)
	}
}

func NewShadel() Shadel {
	vsh, _ := ioutil.ReadFile("shader.vsh")
	fsh, _ := ioutil.ReadFile("shader.fsh")
	reload := make(chan Event)
	s := &ShadelData{
		reload: reload,
	}
	s.newProgram(string(vsh)+"\x00", string(fsh)+"\x00")
	vertAttrib := uint32(gl.GetAttribLocation(s.Program, gl.Str("v\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	go shaderWatcher(reload)
	//s.Uniforms()
	return s
}

func (s *ShadelData) Uniforms() []*Uniform {
	var uniformCount int32 = -1
	gl.GetProgramiv(s.Program, gl.ACTIVE_UNIFORMS, &uniformCount)
	glPanicOnError()
	retVal := make([]*Uniform, uniformCount)
	for i := int32(0); i < uniformCount; i++ {
		name := [100]uint8{}
		var nameLength int32
		var size int32
		var t uint32
		gl.GetActiveUniform(s.Program, uint32(i), 100, &nameLength, &size, &t, &name[0])
		retVal[i] = &Uniform{
			name: string(name[0:nameLength]),
			t:    toType(t),
		}
	}
	log.Printf(`Uniforms are:
	------
	%#v
	------`, retVal)
	return retVal
}

func toType(t uint32) UniformType {
	switch t {
	case gl.FLOAT:
		return TypeFloat
	case gl.INT:
		return TypeInt
	case gl.FLOAT_VEC2:
		return TypeVec2
	case gl.FLOAT_VEC3:
		return TypeVec3
	}
	return TypeInvalid
}

func glPanicOnError() {
	err := gl.GetError()
	if err != gl.NO_ERROR {
		var r string
		switch err {
		case gl.INVALID_VALUE:
			r = "invalid value"
		case gl.INVALID_OPERATION:
			r = "invalid operation"
		case gl.INVALID_ENUM:
			r = "invalid enum"
		}
		panic(r)
	}
}

func (s *ShadelData) GetProgram() uint32 {
	return s.Program
}

func (s *ShadelData) newProgram(vertexShaderSource, fragmentShaderSource string) error {
	vertexShader, err := s.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragmentShader, err := s.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	s.Program = gl.CreateProgram()

	gl.AttachShader(s.Program, vertexShader)
	gl.AttachShader(s.Program, fragmentShader)
	gl.LinkProgram(s.Program)

	var status int32
	gl.GetProgramiv(s.Program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.Program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.Program, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}
	s.Uniforms()

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return err
}

func (s *ShadelData) compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func shaderWatcher(reload chan<- Event) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					log.Println("modified file:", event.Name)
					if strings.HasPrefix(event.Name, "shader") {
						log.Printf("sending reload event for %s %s", event.Name, event.Op)
						reload <- Event{EVENT_TYPE_RELOAD, ""}
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
