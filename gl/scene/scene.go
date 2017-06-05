package scene

import (
	"bitbucket.org/rift_collabo/4kwerfer/gl/object"
	"bitbucket.org/rift_collabo/4kwerfer/gl/target"
	"time"
)

// Scene is just a collection of objects, very simple
// TODO objects would need to have positions and a scene needs a cam
type Scene struct {
	objects []object.Object
}

// NewScene creates a scene with the given objects
func NewScene(objects ...object.Object) *Scene {
	return &Scene{
		objects: objects,
	}
}

// Render the scene to the given target
func (s *Scene) Render(bindTarget target.RenderTarget) {
	bindTarget()
	for _, object := range s.objects {
		object.Render()
	}
}

// Update updates the scene with the given time
func (s *Scene) Update(t time.Duration) {
	for _, o := range s.objects {
		o.Update(t)
	}
}
