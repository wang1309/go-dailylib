package simple

import (
	"fmt"
	"github.com/imdario/mergo"
	"log"
	"testing"
)

type redisConfig struct {
	Address string
	Port    int
	DB      int
	User    []string
}

var defaultConfig = redisConfig{
	Address: "127.0.0.1",
	Port:    6381,
	DB:      1,
	User:    []string{"aaa", "bbb"},
}

func TestMG(t *testing.T) {
	var config redisConfig
	config.DB = 2
	if err := mergo.Merge(&config, defaultConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Println("redis address: ", config.Address)
	fmt.Println("redis port: ", config.Port)
	fmt.Println("redis db: ", config.DB)

	var m = make(map[string]interface{})
	if err := mergo.Map(&m, defaultConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m)
}
