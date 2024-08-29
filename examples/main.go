package main

import (
	"fmt"

	"github.com/yannickkirschen/sip/pkg/sip"
)

type MyConfig struct {
	Profile string        `sip:"profile"`
	Bla     int8          `sip:"bla"`
	User    *MyUserConfig `sip:"user"`
}

type MyUserConfig struct {
	Username string `sip:"username"`
	Password string `sip:"password"`
}

func main() {
	config := &MyConfig{}

	sip.RegisterFile("examples/my.json")
	if err := sip.Fill(config, "my"); err != nil {
		panic(err)
	}
	fmt.Println(config)
	fmt.Println(config.User)
}
