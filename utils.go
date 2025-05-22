package utils

import (
	"fmt"
	"net"
	"os"
)

func isVirtualInterface(iface net.Interface) bool {
	// 方法1：检查标志位组合
	if iface.Flags&net.FlagLoopback != 0 {
		return true
	}

	// 方法2：检查MAC地址厂商标识
	if len(iface.HardwareAddr) >= 3 {
		oui := fmt.Sprintf("%02x:%02x:%02x",
			iface.HardwareAddr[0],
			iface.HardwareAddr[1],
			iface.HardwareAddr[2])

		// 已知虚拟化厂商OUI列表
		virtualOUI := []string{
			"00:05:69", // VMware
			"00:0c:29", // VMware
			"00:1c:14", // VMware
			"08:00:27", // VirtualBox
			"00:50:56", // VMware ESX
			"00:1c:42", // Parallels
		}
		for _, v := range virtualOUI {
			if oui == v {
				return true
			}
		}
	}

	return false
}

func GetMac() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && !isVirtualInterface(iface) {
			mac := iface.HardwareAddr.String()
			if mac != "" {
				return mac
			}
		}
	}
	return ""
}

// 获取本机IP，排除虚拟网卡和Docker容器
func GetIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && !isVirtualInterface(iface) {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

// 获取主机名
func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
