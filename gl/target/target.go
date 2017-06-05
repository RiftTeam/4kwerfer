package target

import (
	gl "github.com/go-gl/gl/v4.1-core/gl"
)

// RenderTarget is just an abstraction over targets like framebuffers and screens
type RenderTarget func()

type PingPong struct {
	fbo   uint32
	rtt   []uint32
	index uint8
}

// NewPingPong creates a new ping pong render target
func NewPingPong(windowWidth, windowHeight int32) *PingPong {
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

	target := &PingPong{
		rtt:   rtt,
		fbo:   fbo,
		index: 0,
	}
	return target
}

// BindLastTexture bind the last texture rendered to, so you can apply it to another object (e.g. post processing)
func (t *PingPong) BindLastTexture() {
	gl.BindTexture(gl.TEXTURE_2D, t.rtt[t.index])
}

// Bind is an implementation of the RenderTarget function type
func (t *PingPong) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, t.rtt[t.index], 0)
	gl.BindTexture(gl.TEXTURE_2D, t.rtt[1-t.index])
	t.index = 1 - t.index
}

// Screen is the target for rendering to screen
var Screen = func() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
