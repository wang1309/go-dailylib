package test

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"testing"
)

func TestGetIni(t *testing.T) {
	cfg, err := ini.Load("../my.ini")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}

	fmt.Println("App Name: ", cfg.Section("").Key("app_name").String())
	fmt.Println("Log level: ", cfg.Section("").Key("log_level").String())

	fmt.Println("MySQL IP:", cfg.Section("mysql").Key("ip").String())
	mysqlPort, err := cfg.Section("mysql").Key("port").Int()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MySQL Port:", mysqlPort)
	fmt.Println("MySQL User:", cfg.Section("mysql").Key("user").String())
	fmt.Println("MySQL Password:", cfg.Section("mysql").Key("password").String())
	fmt.Println("MySQL Database:", cfg.Section("mysql").Key("database").String())

	fmt.Println("Redis IP:", cfg.Section("redis").Key("ip").String())
	redisPort, err := cfg.Section("redis").Key("port").Int()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Redis Port:", redisPort)
}

func TestMustInt(t *testing.T) {
	cfg, err := ini.Load("../my.ini")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}

	redisPassword, err := cfg.Section("redis").Key("password").Int()
	if err != nil {
		fmt.Println("before must, get redis password error:", err)
	} else {
		fmt.Println("before must, get redis password:", redisPassword)
	}

	redisPassword = cfg.Section("redis").Key("password").MustInt(111)
	fmt.Println(redisPassword)

	redisPassword, err = cfg.Section("redis").Key("password").Int()
	if err != nil {
		fmt.Println("before must, get redis password error:", err)
	} else {
		fmt.Println("before must, get redis password:", redisPassword)
	}
}

func TestGetSectionName(t *testing.T) {
	cfg, err := ini.Load("../my.ini")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}

	sections := cfg.Sections()
	names := cfg.SectionStrings()

	fmt.Println("sections: ", sections)
	fmt.Println("names: ", names)
}

func TestParentChild(t *testing.T) {
	cfg, err := ini.Load("../parent_child.ini")
	if err != nil {
		fmt.Println("Fail to read file: ", err)
		return
	}

	fmt.Println("Clone url from package.sub:", cfg.Section("package.sub").Key("CLONE_URL").String())
}

func TestWrite(t *testing.T) {
	cfg := ini.Empty()

	defaultSection := cfg.Section("")
	defaultSection.NewKey("app_name", "awesome web")
	defaultSection.NewKey("log_level", "DEBUG")

	mysqlSection, err := cfg.NewSection("mysql")
	if err != nil {
		fmt.Println("new mysql section failed:", err)
		return
	}
	mysqlSection.NewKey("ip", "127.0.0.1")
	mysqlSection.NewKey("port", "3306")
	mysqlSection.NewKey("user", "root")
	mysqlSection.NewKey("password", "123456")
	mysqlSection.NewKey("database", "awesome")

	redisSection, err := cfg.NewSection("redis")
	if err != nil {
		fmt.Println("new redis section failed:", err)
		return
	}
	redisSection.NewKey("ip", "127.0.0.1")
	redisSection.NewKey("port", "6381")

	err = cfg.SaveTo("../my.ini")
	if err != nil {
		fmt.Println("SaveTo failed: ", err)
	}

	err = cfg.SaveToIndent("../my-pretty.ini", "\t")
	if err != nil {
		fmt.Println("SaveToIndent failed: ", err)
	}

	cfg.WriteTo(os.Stdout)
	fmt.Println()
	cfg.WriteToIndent(os.Stdout, "\t")
}

type Config struct {
	AppName  string `ini:"app_name"`
	LogLevel string `ini:"log_level"`

	MySQL MySQLConfig `ini:"mysql"`
	Redis RedisConfig `ini:"redis"`
}

type MySQLConfig struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

type RedisConfig struct {
	IP   string `ini:"ip"`
	Port int    `ini:"port"`
}

func TestMapTo(t *testing.T) {
	cfg, err := ini.Load("../my.ini")
	if err != nil {
		fmt.Println("load my.ini failed: ", err)
	}

	c := Config{}
	cfg.MapTo(&c)

	fmt.Printf("%+v\n", c)


	mysqlCfg := MySQLConfig{}
	err = cfg.Section("mysql").MapTo(&mysqlCfg)
	fmt.Printf("%+v\n", mysqlCfg)
}
