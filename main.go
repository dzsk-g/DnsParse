package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

var (
	api         = "http://api.ip33.com/dns/resolver"
	domains     = [3]string{"api.themoviedb.org", "image.tmdb.org", "www.themoviedb.org"}
	dnsProvider = [2]string{"156.154.70.1", "208.67.222.222"}
)

// Data /**服务器返回数据
type Data struct {
	Dns    string `json:"dns"`
	Type   string `json:"type"`
	Domain string `json:"domain"`
	State  bool   `json:"state"`
	Record []struct {
		Ip  string `json:"ip"`
		Ttl int    `json:"ttl"`
	} `json:"record"`
}

// ResolveIP 解析IP
func ResolveIP(domain string, dns string) ([]string, error) {
	var ips []string
	v := url.Values{}
	v.Add("domain", domain)
	v.Add("type", "A")
	v.Add("dns", dns)
	resp, err := http.PostForm(api, v)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	reb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := &Data{}
	err = json.Unmarshal(reb, d)
	if err != nil {
		return nil, err
	}
	for _, s := range d.Record {
		ips = append(ips, s.Ip)
	}
	return ips, nil
}

// IsIPReachable 简单检测IP是否可达，这里通过尝试连接443端口来判断，可按需调整逻辑
func IsIPReachable(ip string) bool {
	conn, err := net.Dial("tcp", ip+":443")
	if err == nil {
		_ = conn.Close()
		return true
	}
	return false
}

func main() {
	var entries []HostEntry
	for _, domain := range domains {
		for _, dns := range dnsProvider {
			ips, err := ResolveIP(domain, dns)
			if err != nil {
				fmt.Printf("解析IP错误 %s: %s\n", domain, err)
				continue
			}
			for _, ip := range ips {
				isUse := IsIPReachable(ip)
				fmt.Printf("%s\t%s\t%t\n", ip, domain, isUse)
				if isUse {
					var hosts []string
					hosts = append(hosts, domain)
					entry := HostEntry{
						ip,
						hosts,
					}
					entries = append(entries, entry)
				}
			}
		}
	}
	err := UpdateHosts(entries)
	if err != nil {
		fmt.Printf("hosts 更新失败: %s\n", err)
	}
}
