package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// HostEntry 结构体表示hosts文件中的一行记录，包含IP地址和主机名列表
type HostEntry struct {
	IP    string
	Hosts []string
}

var hostsFilePath = ""

const startTag = "#TMDB_HOST_START"
const endTag = "#TMDB_HOST_END"

// 根据系统类型获取Host文件路径
func getHostsFilePath() string {
	if runtime.GOOS == "linux" {
		return "/etc/hosts"
	} else if runtime.GOOS == "windows" {
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	}
	return ""
}

// UpdateHosts 更新hosts文件
func UpdateHosts(entries []HostEntry) error {
	hostsFilePath = getHostsFilePath()
	if hostsFilePath == "" {
		return fmt.Errorf("未知操作系统")
	}
	hostsData, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return fmt.Errorf("读取HOSTS文件失败: %v", err)
	}
	lines := strings.Split(string(hostsData), "\n")
	inUpdateRange := false
	var updateLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, startTag) {
			inUpdateRange = true
			continue
		}
		if strings.HasPrefix(line, endTag) {
			inUpdateRange = false
			continue
		}
		//不在指定范围内，则直接添加原有行
		if !inUpdateRange {
			updateLines = append(updateLines, line)
		}
	}
	updateLine(&updateLines, entries)
	// 将更新后的内容重新组合成字符串
	newHostsData := strings.Join(updateLines, "\n")
	// 写入更新后的HOSTS文件
	err = os.WriteFile(hostsFilePath, []byte(newHostsData), 0644)
	if err != nil {
		return fmt.Errorf("写入HOSTS文件失败: %v", err)
	}
	return nil
}

// 更新行
func updateLine(updateLines *[]string, entries []HostEntry) {
	//添加标志位
	st := fmt.Sprintf("%s\t%s\n", startTag, "Update on"+time.Now().Format("2006-01-02 15:04:05"))
	*updateLines = append(*updateLines, st)
	for _, entry := range entries {
		line := fmt.Sprintf("%s\t%s\n", entry.IP, strings.Join(entry.Hosts, "\t"))
		*updateLines = append(*updateLines, line)
	}
	et := fmt.Sprintf("%s\n", endTag)
	*updateLines = append(*updateLines, et)
}
