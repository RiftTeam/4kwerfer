package object

import (
	"time"

	shader "bitbucket.org/rift_collabo/4kwerfer/gl/shader"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Object is some random drawable object
type Object interface {
	Render()
	Update(time.Duration)
}

type shadedObject struct {
	obj    Object
	shader shader.Shadel
}

// ApplyShader applies a shader to a given object
func ApplyShader(shader shader.Shadel, obj Object) Object {
	return &shadedObject{
		obj:    obj,
		shader: shader,
	}
}

func (o *shadedObject) Render() {
	o.shader.Use()
}

func (o *shadedObject) Update(t time.Duration) {
	o.obj.Update(t)
}

type vaoScene struct {
	vaoID  uint32
	vboID  uint32
	shader shader.Shadel
}

type fsQuad struct {
	*vaoScene
	w int32
	h int32
}

// NewFullScreenQuad creates a simple fullscreen quad
func NewFullScreenQuad(shader shader.Shadel, windowWidth, windowHeight int32) Object {
	// Configure the vertex data
	fsQuadVertices := []float32{
		1.0, 1.0, 0.0, // 1.0, 1.0, // vertex 0
		-1.0, 1.0, 0.0, // 0.0, 1.0, // vertex 1
		1.0, -1.0, 0.0, // 1.0, 0.0, // vertex 2
		-1.0, -1.0, 0.0, // 0.0, 0.0, // vertex 3
	}
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(fsQuadVertices)*4, gl.Ptr(fsQuadVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(shader.GetProgram(), gl.Str("v\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	return &fsQuad{
		vaoScene: &vaoScene{
			vaoID:  vao,
			vboID:  vbo,
			shader: shader,
		},
		w: windowWidth,
		h: windowHeight,
	}
}

func (o *fsQuad) Render() {
	o.vaoScene.Render()
}

func (o *fsQuad) Update(t time.Duration) {
	//log.Printf("%f", float32(t/time.Millisecond)/1000.0)
	o.shader.SetUniform3f("u", float32(o.w)*2, float32(o.h)*2, float32(t/time.Millisecond)/1000.0)
}

func (o *vaoScene) Render() {
	o.shader.Use()
	gl.BindBuffer(gl.ARRAY_BUFFER, o.vboID)

	gl.BindVertexArray(o.vaoID)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

}

func (o *vaoScene) Update(t time.Duration) {

}
