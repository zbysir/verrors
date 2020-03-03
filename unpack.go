package verrors

import (
	"encoding/json"
)

type Store interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	GetAll() map[string]interface{}
}

type Setter interface {
	Set(Store)
}

type PackError struct {
	values map[string]interface{}
	err    error
}

func (p PackError) Error() string {
	if p.err == nil {
		return "nil"
	}
	return p.err.Error()
}

func (p *PackError) GetAll() map[string]interface{} {
	return p.values
}

func (p *PackError) Get(k string) (interface{}, bool) {
	if p.values == nil {
		return nil, false
	}
	i, ok := p.values[k]
	return i, ok
}

func (p *PackError) Set(k string, v interface{}) {
	if p.values == nil {
		p.values = map[string]interface{}{}
	}
	p.values[k] = v
}

func (p PackError) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"err": p.err.Error(),
	}
	if p.values != nil {
		m["values"] = p.values
	}
	return json.Marshal(m)
}

// Cause取得内部错误
func (p PackError) Cause() error {
	return p.err
}

type PackErrors []PackError

func (u PackErrors) Error() string {
	e := u.Merge()
	return e.Error()
}

func (u PackErrors) Last() (err PackError) {
	if len(u) == 0 {
		return
	}

	return u[len(u)-1]
}

func (u PackErrors) First() (err PackError) {
	if len(u) == 0 {
		return
	}

	return u[0]
}

// 将多层错误合并成一层, 最外层的数据会覆盖内层.
func (u PackErrors) Merge() (err PackError) {
	if len(u) == 0 {
		return
	}
	var p PackError

	for i := len(u) - 1; i >= 0; i-- {
		v := u[i]
		p.err = v.err
		d := v.GetAll()
		for k, v := range d {
			p.Set(k, v)
		}
	}

	return p
}

// Unpack会连续Unwrap, 返回错误数组
func Unpack(err error) (es PackErrors) {
	next := err
	for next != nil {
		var p PackError
		next, p = UnpackOnce(next)
		es = append(es, p)
	}

	return
}

// 只解包一层
// unpackOnce会将解析err下的所有internalError, 并将internalError的数据放入packError
// next可能是真正错误, 也可能是internalError
// packError.err 一定是真正的错误
func UnpackOnce(err error) (next error, packError PackError) {
	var internal []error
	// 跳过internalError
	if _, isInternal := err.(InternalError); isInternal {
		// 如果自己是InternalError, 则返回的internal包含自己
		for i, ok := err.(InternalError); ok; i, ok = err.(InternalError) {
			internal = append(internal, err)
			err = i.InternalError()
		}
	}

	// 跳过internalError之后的错误一定是真正错误
	packError.err = err

	// Set
	// 后写的代码(外层的InternalError)覆盖之前
	for i := len(internal) - 1; i >= 0; i-- {
		v := internal[i]
		if s, ok := v.(Setter); ok {
			s.Set(&packError)
		}
	}

	// 将真正错误Unwrap
	if w, ok := err.(Wrapper); ok {
		next = w.Unwrap()
	} else {
		next = nil
	}
	return
}
