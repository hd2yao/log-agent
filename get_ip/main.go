package main

import (
    "fmt"
    "log"
    "net"
)

func GetLocalIP() (ip string, err error) {
    // 获取本机的所有网络接口的地址
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return
    }

    // 遍历所有网络接口的地址
    for _, addr := range addrs {
        ipAddr, ok := addr.(*net.IPNet) // 类型断言
        if !ok {
            continue
        }
        fmt.Println(ipAddr)
        // 排除回环接口地址
        if ipAddr.IP.IsLoopback() {
            continue
        }

        // 排除非全局单播地址（例如，排除私有 IP 地址）
        if !ipAddr.IP.IsGlobalUnicast() {
            continue
        }
        fmt.Println(ipAddr)
        return ipAddr.IP.String(), nil
    }
    return
}

func GetOutboundIP() string {
    conn, err := net.Dial("tcp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.TCPAddr)
    fmt.Println(localAddr.String())
    return localAddr.IP.String()
}

func main() {
    //GetLocalIP()
    GetOutboundIP()
}
