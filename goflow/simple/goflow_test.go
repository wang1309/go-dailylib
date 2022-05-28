package simple_gf

import (
	"fmt"
	goflow "github.com/kamildrazkiewicz/go-flow"
	"testing"
	"time"
)

func TestGoFlow(t *testing.T)  {
	f1 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function1 started")
		time.Sleep(time.Millisecond * 1000)
		return 1, nil
	}

	f2 := func(r map[string]interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("function2 started", r["f1"])
		return "some results f2", nil
	}

	f3 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function3 started", r["f1"])
		return "some results f3", nil
	}

	f4 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function4 started", r)
		return nil, nil
	}

	res, err := goflow.New().
		Add("f1", nil, f1).
		Add("f2", []string{"f1"}, f2).
		Add("f3", []string{"f1"}, f3).
		Add("f4", []string{"f2", "f3"}, f4).
		Do()

	fmt.Println(res, err)
}
