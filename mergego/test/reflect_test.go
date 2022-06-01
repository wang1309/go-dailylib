package test

import (
	"fmt"
	"reflect"
	"testing"
)

type integer int

func TestReflect(t *testing.T)  {
	var a integer
	typeOfA := reflect.TypeOf(a)

	fmt.Println(typeOfA.Name(), typeOfA.Kind())
}
