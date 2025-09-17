package main

import (
	"fmt"
	//	"go/build"
	"log"
	"runtime"
	"sync"

	object "github.com/RiftTeam/4kwerfer/gl/object"

	riftgl "github.com/RiftTeam/4kwerfer/gl/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	//"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 1280
const windowHeight = 720

//var program uint32

var pLock = &sync.Mutex{}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	shadel := riftgl.NewShadel(riftgl.VshFile, riftgl.FshFile)

	// Configure the vertex and fragment shaders

	shadel.Use()
	//	gl.UseProgram(program)

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(fsQuad)*4, gl.Ptr(fsQuad), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(shadel.GetProgram(), gl.Str("v\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	rtt := make([]uint32, 2)

	gl.GenTextures(2, &rtt[0])

	gl.BindTexture(gl.TEXTURE_2D, rtt[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)

	gl.BindTexture(gl.TEXTURE_2D, rtt[1])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)

	var fbo uint32
	gl.GenFramebuffers(1, &fbo)

	//previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		// Update
		//time := glfw.GetTime()
		//elapsed := time - previousTime
		//previousTime = time

		// Render
		//gl.UseProgram(program)

		//	shadel.Use()
		shadel.Use()

		gl.BindVertexArray(vao)

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func Draw() {
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
