package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	// "scanner/tcp_connect_scanner/scanner"
	"time"

	"github.com/malfunkt/iprange"
	// "test_scanner/scanner"
	// "sec-dev-in-action-src/scanner/tcp-connect-scanner-demo/util"
)

// Connect 函数建立一个到指定 IP 地址和端口的连接。
// 它接受两个参数，一个 IP 地址（字符串）和一个端口号（int）。
// 它返回一个 net.Conn 对象表示连接和一个错误（如果有）。
func Connect(ip string, port int) (net.Conn, error) {
	// DialTimeout 函数创建一个新的 TCP 连接到指定的 IP 和端口，超时时间为 1 秒。
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 1*time.Second)
	// defer 语句确保在函数退出前关闭连接。
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	// 返回连接和错误（如果有）
	return conn, err
}

// 封装一个GetIpList函数，可以根据输入的ipList返回一个[]net.IP的切片
func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	list := addressList.Expand()
	return list, err
}

// 多端口的处理需要支持",“与”-“分割的端口列表
// 可以使用strings包的Split函数先分割以”,“连接的ipList
// 然后再分割以”-"连接的ipList，最后返回一个[]int切片
func GetPorts(selection string) ([]int, error) {
	ports := []int{}
	if selection == "" {
		return ports, nil
	}
	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invaild port selection sequment: '%s'", r)
			}
			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invaild port number: '%s'", parts[0])
			}
			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invaild port number: '%s'", parts[1])
			}
			if p1 > p2 {
				return nil, fmt.Errorf("Invaild port range: %d-%d", p1, p2)
			}
			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}
		} else if port, err := strconv.Atoi(r); err != nil {
			return nil, fmt.Errorf("Invaild port number: '%s'", r)
		} else {
			ports = append(ports, port)
		}
	}
	return ports, nil
}

func main() {
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, _ := GetIpList(ipList)
		ports, _ := GetPorts(portList)
		for _, ip := range ips {
			for _, port := range ports {
				_, err := Connect(ip.String(), port)
				if err != nil {
					continue
				}
				fmt.Printf("ip: %v, port: %v is open \n", ip, port)
			}
		}
	} else {
		fmt.Printf("Need two Args, format is %v iplist port", os.Args[0])
	}

}
