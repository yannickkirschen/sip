# sip

[![Lint commit message](https://github.com/yannickkirschen/sip/actions/workflows/commit-lint.yml/badge.svg)](https://github.com/yannickkirschen/sip/actions/workflows/commit-lint.yml)
[![Push](https://github.com/yannickkirschen/sip/actions/workflows/push.yml/badge.svg)](https://github.com/yannickkirschen/sip/actions/workflows/push.yml)
[![Release](https://github.com/yannickkirschen/sip/actions/workflows/release.yml/badge.svg)](https://github.com/yannickkirschen/sip/actions/workflows/release.yml)
[![GitHub release](https://img.shields.io/github/release/yannickkirschen/sip.svg)](https://github.com/yannickkirschen/sip/releases/)

sip fills a Go struct with data from any sources.

> [!WARNING]
> This project is currently under development and not stable.

## Usage

```go
type MyConfig struct {
 Profile string        `sip:"profile"`
 User    *MyUserConfig `sip:"user"`
}

type MyUserConfig struct {
 Username string `sip:"username"`
 Password string `sip:"password"`
}

func main() {
 config := &MyConfig{}
 prefix := "my"

 sip.RegisterFile("examples/my.json")
 if err := sip.Fill(config, prefix); err != nil {
  panic(err)
 }
}
```
