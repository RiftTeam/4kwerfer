package shader

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// file names for shaders
const (
	VshFile     = "shader.vsh"
	FshFile     = "shader.fsh"
	ExeFshFile  = "frag.glsl"
	CombineFile = "ppfrag.glsl"
)

// ReplaceShadel replaces the current shader program with a new one loaded from file
func (s *ShadelData) ReplaceShadel(vshFile, fshFile string) error {
	p := s.Program
	vsh, _ := ioutil.ReadFile(VshFile)
	fsh, _ := ioutil.ReadFile(FshFile)
	if err := s.newProgram(string(vsh)+"\x00", string(fsh)+"\x00"); err != nil {
		log.Printf("Failed to replace shader, keeping last program: %s", err.Error())
		s.err = true
	} else {
		if !s.err {
			defer gl.DeleteProgram(p)
		}
		s.err = false
		s.Use()
	}
	return nil
}

// Use activates the shader
func (s *ShadelData) Use() {
	select {
	case ev := <-s.reload:
		log.Printf("reloading %#v", ev)
		s.ReplaceShadel(VshFile, FshFile)
		log.Print("done")
		gl.UseProgram(s.Program)
		// non blocking send, event is dropped if there is no one paying attention
		select {
		case s.reloaded <- nil:
		default:
		}

	default:
		gl.UseProgram(s.Program)
	}

}

// NewShadel creates a new shader from files
func NewShadel(vshFile string, fshFile string) Shadel {
	vsh, _ := ioutil.ReadFile(vshFile)
	fsh, _ := ioutil.ReadFile(fshFile)
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

// ShaderChanged returns a channel on which an event is sent after the shader has been reloaded
func (s *ShadelData) ShaderChanged() <-chan interface{} {
	return s.reloaded
}

// GetProgram gets the actual program handle
func (s *ShadelData) GetProgram() uint32 {
	return s.Program
}

func (s *ShadelData) newProgram(vertexShaderSource, fragmentShaderSource string) error {
	log.Printf("Compiling vert shader")
	vertexShader, err := s.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		println(err.Error())
		return err
	}

	log.Printf("Compiling frag shader")
	fragmentShader, err := s.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		println(err.Error())
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
		fmt.Printf("---\n%v\n", log)

		return fmt.Errorf("failed to link program: %v", log)
	}
	println(">>>")
	for glerr := gl.GetError(); glerr != gl.NO_ERROR; glerr = gl.GetError() {
		fmt.Printf("clearing error, %#v\n", glerr)
	}
	//s.uniforms =
	s.extractUniforms()

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	log.Printf("using new program: %d", s.Program)
	gl.UseProgram(s.Program)

	return err
}

func (s *ShadelData) extractUniforms() map[string]uniformData {
	var uniformCount int32 = -1
	gl.GetProgramiv(s.Program, gl.ACTIVE_UNIFORMS, &uniformCount)
	glPanicOnError()
	var maxLength int32 = -1
	gl.GetProgramiv(s.Program, gl.ACTIVE_UNIFORM_MAX_LENGTH, &maxLength)
	glPanicOnError()
	retVal := make(map[string]uniformData, uniformCount)
	for i := int32(0); i < uniformCount; i++ {
		nameBuffer := make([]uint8, maxLength)
		var nameLength int32
		var size int32
		var t uint32
		gl.GetActiveUniform(s.Program, uint32(i), 100, &nameLength, &size, &t, &nameBuffer[0])
		glPanicOnError()
		name := string(nameBuffer[0:nameLength])
		retVal[name] = uniformData{
			name:     name,
			t:        toType(t),
			location: gl.GetUniformLocation(s.Program, &nameBuffer[0]),
		}
		log.Printf("%s", retVal[name])
	}
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
		println(log)

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
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					if strings.HasPrefix(event.Name, "shader") {
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

// SetUniform3f updates the uniform with the fiven name to the given values
func (s *ShadelData) SetUniform3f(name string, x, y, z float32) {
	gl.UseProgram(s.Program)
	uniLoc := gl.GetUniformLocation(s.Program, gl.Str(name+"\x00"))
	gl.Uniform3f(uniLoc, x, y, z)
}
