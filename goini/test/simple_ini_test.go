package test

import (
	"dailylib/goini/simple"
	"fmt"
	"testing"
)

func TestSimpleIni(t *testing.T)  {
	cfg, err := simple.Load("../my.ini")
	if err != nil {
		fmt.Printf("load err: %+v\n", err)
	}

	// fmt.Printf("cfg: %+v", cfg)
	fmt.Println(cfg.Section("mysql").Key("ip").String())

	// Todo 目前只实现获取 key value 这一个方法

}
