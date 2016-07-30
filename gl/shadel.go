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

type Shadel interface {
	ReplaceShadel(string, string) error
	Use()
}

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

func NewShadel() *ShadelData {
	vsh, _ := ioutil.ReadFile("shader.vsh")
	fsh, _ := ioutil.ReadFile("shader.fsh")
	reload := make(chan Event)
	s := &ShadelData{
		reload: reload,
	}
	s.newProgram(string(vsh)+"\x00", string(fsh)+"\x00")
	go ExampleNewWatcher(reload)
	return s
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

func ExampleNewWatcher(reload chan<- Event) {
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
