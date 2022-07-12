package test

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"testing"
)

type Person struct {
	Name string `mapstructure:"name"`
	Age  int    `mapstructure:"age"`
	Job  string `mapstructure:",omitempty"`
}

func TestS2M(t *testing.T) {


	p1 := Person{}

	m := S2M()

	mapstructure.Decode(m, &p1)

	fmt.Printf("%+v\n",p1)

	/*data, _ := json.Marshal(m)
	fmt.Println(string(data))*/
}

func S2M() interface{} {
	p := &Person{
		Name: "dj",
		Age:  18,
	}

	var m map[string]interface{}
	mapstructure.Decode(p, &m)

	fmt.Println(m)

	return m
}