package main

import (
	"fmt"
	//	"go/build"
	"flag"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/RiftTeam/4kwerfer/config"
	"github.com/RiftTeam/4kwerfer/gl/object"
	"github.com/RiftTeam/4kwerfer/gl/scene"
	"github.com/RiftTeam/4kwerfer/gl/shader"
	"github.com/RiftTeam/4kwerfer/gl/target"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	//"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 1280
const windowHeight = 720

//var program uint32

var pLock = &sync.Mutex{}

var cfg = config.Config{}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	flag.IntVar(&cfg.MaxIters, "i", 256, "The maximum number of iterations")
	flag.IntVar(&cfg.MaxDuration, "d", 30, "The maximum number of seconds to do progressive rendering")
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Duck! This is 4kwerfer", nil, nil)
	if err != nil {
		panic(err)
	}
	abortRender := make(chan interface{})
	window.SetKeyCallback(
		func(w *glfw.Window,
			key glfw.Key,
			scancode int,
			action glfw.Action,
			mods glfw.ModifierKey,
		) {

			if key == glfw.KeyEscape {
				os.Exit(0)
			}
			if key == glfw.KeySpace {
				go func() { abortRender <- nil }()
			}
		},
	)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	progressiveShader := shader.NewShadel(shader.VshFile, shader.ExeFshFile)
	//presentationProg := shader.NewShadel(shader.VshFile, shader.CombineFile)

	//redrawChan := progressiveShader.ShaderChanged()

	progressiveQuad := object.NewFullScreenQuad(progressiveShader, windowWidth, windowHeight)
	//finalQuad := object.NewFullScreenQuad(progressiveShader, windowWidth, windowHeight)

	//	pingPong := target.NewPingPong(windowWidth, windowHeight)
	scene := scene.NewScene(progressiveQuad)
	startTime := time.Now()
	//for {
	scene.Update(time.Since(startTime))
	scene.Render(target.Screen)
	window.SwapBuffers()
	glfw.PollEvents()
	//}
	<-time.After(10 * time.Second)

	//	for {

	/*timeOut := time.After(time.Duration(cfg.MaxDuration) * time.Second)
	var i uint64
	for i = 0; i < uint64(cfg.MaxIters); i++ {
		glfw.PollEvents()
		log.Printf("rendering iteration %d", i)
		select {
		case <-timeOut:
			break
		case <-abortRender:
			break
		default:
			scene.Update(time.Duration(i) * time.Second)
			scene.Render(pingPong.Bind)
		}
	}

	glfw.PollEvents()
	*/
	//scene.Update(time.Duration(i) * time.Second)
	//		scene.Render(target.Screen)
	//		window.SwapBuffers()
	//	<-progressiveShader.ShaderChanged()
	//	}
}

func shouldClose(window *glfw.Window) bool {
	return window.ShouldClose()
}

var fsQuad = []float32{
	1.0, 1.0, 0.0, // 1.0, 1.0, // vertex 0
	-1.0, 1.0, 0.0, // 0.0, 1.0, // vertex 1
	1.0, -1.0, 0.0, // 1.0, 0.0, // vertex 2
	-1.0, -1.0, 0.0, // 0.0, 0.0, // vertex 3
}

func breakRender() {

}
