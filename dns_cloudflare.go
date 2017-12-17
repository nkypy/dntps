package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

func lastRecordGet() (*cloudflare.DNSRecord, error) {
	if _, err := os.Stat(path.Join(Root, "data", "record.toml")); err == nil {
		body, err := ioutil.ReadFile(path.Join(Root, "data", "record.toml"))
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

func newRecordGet(key, mail, domain string) (*cloudflare.API, *cloudflare.DNSRecord, error) {
	api, err := cloudflare.New(key, mail)
	if err != nil {
		return nil, nil, err
	}
	zoneName := strings.Join(strings.Split(domain, ".")[1:], ".")
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return nil, nil, err
	}
	dns, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: domain})
	if err != nil {
		return nil, nil, err
	}
	if len(dns) != 1 {
		return nil, nil, errors.New("没有此域名")
	}
	record := dns[0]
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(record); err != nil {
		return nil, nil, err
	}
	if err := ioutil.WriteFile(path.Join(Root, "data", "record.toml"), buf.Bytes(), 0644); err != nil {
		return nil, nil, err
	}
	return api, &record, nil
}

func cfUpdateDNS(key, mail, domain string) error {
	newIP, err := ipGet()
	if err != nil {
		return err
	}
	var api *cloudflare.API
	record, err := lastRecordGet()
	if err != nil {
		api, record, err = newRecordGet(key, mail, domain)
		if err != nil {
			return err
		}
	}
	if *newIP != record.Content {
		record.Content = *newIP
		if err := api.UpdateDNSRecord(record.ZoneID, record.ID, *record); err != nil {
			return err
		}
	}
	return nil
}
