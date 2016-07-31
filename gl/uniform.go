package gl

import (
	"errors"
)

func NewUniform(name string, t UniformType) Uniform {
	return &uniformData{
		name: name,
		t:    t,
	}
}

func (u *uniformData) GetValue() interface{} {
	return 0
}

func (u *uniformData) Name() string {
	return u.name
}

func (u *uniformData) SetValue(v interface{}) error {
	if err := u.checkValueType(v); err != nil {
		return err
	}
	return nil
}

func (u *uniformData) Type() UniformType {
	return u.t
}

func (u *uniformData) ValueString() string {
	return ""
}

func (u *uniformData) checkValueType(v interface{}) error {
	ok := false
	switch u.t {
	case TypeFloat:
		_, ok = v.(float32)
	case TypeVec2:
		_, ok = v.(Vec2)
	case TypeVec3:
		_, ok = v.(Vec3)
	case TypeInt:
		_, ok = v.(int)
	}
	if !ok {
		return errors.New("expected and actual type didn't match")
	}
	return nil
}
