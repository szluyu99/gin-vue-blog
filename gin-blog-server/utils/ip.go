package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"xojoc.pw/useragent"
)

var IP = new(ipUtil)

type ipUtil struct{}

// 获取通用信息: ipAddress, ipSource, browser, os
func (i *ipUtil) GetInfo(c *gin.Context) (ipAddress, browser, os string) {
	// fmt.Println("IP GetInfo")
	ipAddress = i.GetIpAddress(c)
	userAgent := i.GetUserAgent(c)
	// fmt.Println("ipAddress: ", ipAddress)
	browser = userAgent.Name + " " + userAgent.Version.String()
	os = userAgent.OS + " " + userAgent.OSVersion.String()
	return
}

// FIXME: 获取用户发送请求的 IP 地址
func (*ipUtil) GetIpAddress(c *gin.Context) (ipAddress string) {
	ipAddress = c.Request.Header.Get("X-Real-IP")
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("Proxy-Client-IP")
	}
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("WL-Proxy-Client-IP")
	}
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.RemoteAddr
	}
	if ipAddress == "127.0.0.1" {
		ip, err := externalIP()
		log.Println(err)
		ipAddress = ip.String()
	}
	if ipAddress != "" && len(ipAddress) > 15 {
		if strings.Index(ipAddress, ",") > 0 {
			ipAddress = ipAddress[:strings.Index(ipAddress, ",")]
		}
	}
	return
}

// 获取 IP 来源
func (*ipUtil) GetIpSource(ipAddress string) string {
	var dbPath = "./config/ip2region.xdb"
	searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		fmt.Printf("failed to create searcher: %s\n", err.Error())
		return ""
	}
	defer searcher.Close()
	region, err := searcher.SearchByStr(ipAddress)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ipAddress, err)
		return "未知"
	}
	return region
}

func (i *ipUtil) GetIpSourceSimpleIdle(ipAddress string) string {
	return i.IpSourceSimpleSplit(i.GetIpSource(ipAddress))
}

func (*ipUtil) IpSourceSimpleSplit(region string) string {
	ipSource := strings.Split(region, "|")
	if ipSource[0] != "中国" {
		return ipSource[0]
	}
	if ipSource[2] == "0" {
		ipSource[2] = ""
	}
	if ipSource[3] == "0" {
		ipSource[3] = ""
	}
	if ipSource[4] == "0" {
		ipSource[4] = ""
	}
	if ipSource[2] == "" && ipSource[3] == "" && ipSource[4] == "" {
		return ipSource[0]
	}
	return ipSource[2] + ipSource[3] + " " + ipSource[4]
}

func (*ipUtil) GetUserAgent(c *gin.Context) *useragent.UserAgent {
	ua := useragent.Parse(c.Request.UserAgent())
	return ua
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}
