package main

import (
	"fmt"
	//	"go/build"
	"log"
	"runtime"
	"sync"

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

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.EXTENSIONS))
	fmt.Println("OpenGL version", version)

	shadel := NewShadel()

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
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(shadel.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

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

func shouldClose(window *glfw.Window) bool {
	return window.ShouldClose()
}

var vertexShader = `
#version 410
in vec4 vert;
void main() {
    gl_Position = vert;// projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 410
out vec4 outputColor;
void main() {
    outputColor = vec4(0);//texture(tex, fragTexCoord);
}
` + "\x00"

var cubeVertices = []float32{
	1.0, 1.0, 0.0, // 1.0, 1.0, // vertex 0
	-1.0, 1.0, 0.0, // 0.0, 1.0, // vertex 1
	1.0, -1.0, 0.0, // 1.0, 0.0, // vertex 2
	-1.0, -1.0, 0.0, // 0.0, 0.0, // vertex 3
}
