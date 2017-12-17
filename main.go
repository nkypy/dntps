package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml"
)

type Conf struct {
	Domain    *string `toml:"domain"`
	DDNS      *bool   `toml:"ddns"`
	HTTPS     *bool   `toml:"https"`
	DNSServer *string `toml:"dns_server"`
	DNSKey    *string `toml:"dns_key"`
	DNSEmail  *string `toml:"dns_email"`
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
		panic("配置文件不存在")
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
	switch {
	case conf.Domain == nil || conf.DNSServer == nil || conf.DNSKey == nil:
		panic("缺少必要设置")
	case len(strings.Split(*conf.Domain, ".")) != 3:
		panic("域名格式错误")
	case conf.DDNS != nil && conf.DNSServer != nil && *conf.DDNS == true:
		switch *conf.DNSServer {
		case "cf":
			if err := cfUpdateDNS(*conf.DNSKey, *conf.DNSEmail, *conf.Domain); err != nil {
				panic("DNS 更新失败")
			}
		default:
			panic("暂不支持此 DNS 提供商，欢迎提交 PR !")
		}
	case conf.HTTPS != nil && *conf.HTTPS == true:
		log.Println("HTTPS 待开发，欢迎 PR")
	default:
		panic("未进行任何设置")
	}
}
