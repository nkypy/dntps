package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.kshih.com/dntps/dnsapi"

	toml "github.com/pelletier/go-toml"
)

type Conf struct {
	Domain  *string `toml:"domain"`
	HTTPS   *bool   `toml:"https"`
	DNS     *string `toml:"dns"`
	CFKey   *string `toml:"cf_key"`
	CFEmail *string `toml:"cf_email"`
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
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	lastRecordFn := path.Join(Root, "data", "record.toml")
	conf := LoadConf()
	switch {
	case conf.Domain == nil || conf.DNS == nil:
		panic("缺少必要的设置！")
	case len(strings.Split(*conf.Domain, ".")) != 3:
		panic("域名格式错误！")
	default:
		ipCurrent, err := ipGet()
		if err != nil {
			panic("当前 IP 获取失败！")
		}
		switch *conf.DNS {
		case "cf":
			if conf.CFKey == nil || conf.CFEmail == nil {
				panic("缺少 CloudFlare Key 和 Email!")
			}
			dnsapi.CFUpdate(
				*ipCurrent, lastRecordFn,
				*conf.CFKey, *conf.CFEmail, *conf.Domain)
			if conf.HTTPS != nil && *conf.HTTPS == true {
				// 申请 HTTPS
			}
		default:
			panic("暂不支持此 DNS 提供商，欢迎提交 PR !")
		}
	}
}
