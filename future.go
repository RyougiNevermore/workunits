package workunits

import (
	"fmt"
	"reflect"
)

type Future interface {
	Put(v interface{})
	Scan(v interface{}) error
}

type future struct {
	closed bool
	ch     chan interface{}
}

func (f *future) Put(v interface{}) {
	f.ch <- v
	close(f.ch)
}

func (f *future) Scan(target interface{}) error {
	promise := <-f.ch
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Kind() != reflect.Slice || targetValue.Kind() != reflect.Map || targetValue.Kind() != reflect.Chan {
		return fmt.Errorf("target should be a poniter, target kine: %v", targetValue.Kind())
	}
	pValue := reflect.ValueOf(promise)
	if pValue.Kind() == reflect.Ptr {
		pValue = pValue.Elem()
	}
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	if targetValue.Kind() != pValue.Kind() {
		return fmt.Errorf("kind not match, promis is %v, want is %v", pValue.Kind(), targetValue.Kind())
	}
	targetValue.Set(pValue)
	return nil
}
