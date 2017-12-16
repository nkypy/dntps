package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml"
)

type Conf struct {
	Domain    *string `toml:"domain"`
	HTTPS     *bool   `toml:"https"`
	DDNS      *bool   `toml:"ddns"`
	DNSServer *int    `toml:"dns_server"`
	CFMail    *string `toml:"cf_mail"`
	CFKey     *string `toml:"cf_key"`
}

var Root string

func init() {
	var err error
	Root, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("目录出错")
	}
}

func LoadConf() *Conf {
	body, err := ioutil.ReadFile(path.Join(Root, "conf", "conf.toml"))
	if err != nil {
		panic("配置文件出错")
	}
	conf := Conf{}
	err = toml.Unmarshal(body, &conf)
	if err != nil {
		panic("配置文件出错")
	}
	return &conf
}

func main() {
	conf := LoadConf()
	if len(strings.Split(*conf.Domain, ".")) != 3 {
		panic("域名长度错误")
	}
	if err := UpdateDNS(*conf.CFKey, *conf.CFMail, *conf.Domain); err != nil {
		panic("DNS 更新失败")
	}
}
