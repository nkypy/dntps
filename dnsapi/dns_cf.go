package dnsapi

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

func cfLastRecordGet(fn string) (*cloudflare.DNSRecord, error) {
	if _, err := os.Stat(fn); err == nil {
		body, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, err
		}
		record := cloudflare.DNSRecord{}
		if err := toml.Unmarshal(body, &record); err != nil {
			return nil, err
		}
		return &record, nil
	}
	return nil, errors.New("配置文件不存在")
}

func cfNewRecordGet(fn, key, mail, domain string) (*cloudflare.DNSRecord, error) {
	api, err := cloudflare.New(key, mail)
	if err != nil {
		return nil, err
	}
	zoneName := strings.Join(strings.Split(domain, ".")[1:], ".")
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return nil, err
	}
	dns, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: domain})
	if err != nil {
		return nil, err
	}
	if len(dns) != 1 {
		return nil, errors.New("没有此域名")
	}
	record := dns[0]
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(record); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(fn, buf.Bytes(), 0644); err != nil {
		return nil, err
	}
	return &record, nil
}

func CFUpdate(ip, fn, key, mail, domain string) error {
	record, err := cfLastRecordGet(fn)
	if err != nil || ip != record.Content {
		api, errNew := cloudflare.New(key, mail)
		if errNew != nil {
			return errNew
		}
		if err != nil {
			record, err = cfNewRecordGet(fn, key, mail, domain)
			if err != nil {
				return err
			}
		}
		if err := api.UpdateDNSRecord(record.ZoneID, record.ID, *record); err != nil {
			return err
		}
		log.Println("更新 IP：", ip)
	}
	log.Println("不需要更新 IP！")
	return nil
}
