package test

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
	"testing"
)

func TestBp(t *testing.T) {
	b := bytebufferpool.Get()

	b.WriteString("hello")
	b.WriteByte(',')
	b.WriteString(" world")

	//fmt.Println(b.String())

	bytebufferpool.Put(b)


	b2 := bytebufferpool.Get()

	fmt.Printf("%s\n", b2.String())

}
