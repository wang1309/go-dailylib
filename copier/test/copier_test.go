package test

import (
	"fmt"
	"github.com/jinzhu/copier"
	"testing"
)

type User struct {
	Name string
	Age  int
}

type Employee struct {
	Name string
	Age  int
	Role string
}

func TestCp(t *testing.T) {
	user := User{Name: "dj", Age: 18}
	employee := Employee{}

	if err := copier.Copy(&employee, &user); err != nil {
		fmt.Printf("%+v", err)
	}

	fmt.Printf("%+v", employee)
}

func TestSlice2Slice(t *testing.T) {
	users := []User{
		{Name: "dj", Age: 18},
		{Name: "dj2", Age: 18},
	}
	employees := []Employee{}

	if err := copier.Copy(&employees, &users); err != nil {
		fmt.Printf("%+v", err)
	}

	fmt.Printf("%#v\n", employees)
}

func TestStruct2Slice(t *testing.T) {
	user := User{Name: "dj", Age: 18}
	employees := []Employee{}

	if err := copier.Copy(&employees, &user); err != nil {
		fmt.Printf("%+v", err)
	}

	fmt.Printf("%#v\n", employees)
}
