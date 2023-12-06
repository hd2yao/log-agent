package main

import (
    "fmt"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/load"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "time"

    "github.com/shirou/gopsutil/cpu"
)

// cpu info
func getCpuInfo() {
    cpuInfos, err := cpu.Info()
    if err != nil {
        fmt.Printf("get cpu info failed, err:%v\n", err)
    }
    for _, ci := range cpuInfos {
        fmt.Println(ci)
    }
    // cpu 使用率
    for {
        percent, _ := cpu.Percent(time.Second, false)
        fmt.Printf("cpu percent:%v\n", percent)
    }
}

// cpu 负载
func getLoad() {
    info, err := load.Avg()
    if err != nil {
        fmt.Printf("load.Avg() failed, err:%v", err)
        return
    }
    fmt.Println(info)
}

// 内存信息
func getMemInfo() {
    info, err := mem.VirtualMemory()
    if err != nil {
        fmt.Printf("get mem info failed, err:%v", err)
        return
    }
    fmt.Println(info)
}

// host info
func getHostInfo() {
    hInfo, _ := host.Info()
    fmt.Printf("host info:%v uptime:%v boottime:%v\n", hInfo, hInfo.Uptime, hInfo.BootTime)
}

// 磁盘 信息
func getDiskInfo() {
    // 获取所有分区信息
    parts, err := disk.Partitions(true)
    if err != nil {
        fmt.Printf("get disk partitions failed, err:%v", err)
        return
    }
    fmt.Println(parts)
    for _, part := range parts {
        partInfo, err := disk.Usage(part.Mountpoint)
        if err != nil {
            fmt.Printf("get part stat failed, err:%v", err)
            return
        }
        fmt.Println(partInfo)
    }
    // 磁盘IO
    ioStat, _ := disk.IOCounters()
    for k, v := range ioStat {
        fmt.Printf("%v:%v\n", k, v)
    }
}

// 网络信息
func getNetInfo() {
    netIOs, err := net.IOCounters(true)
    if err != nil {
        fmt.Printf("get net io counters failed, err:%v", err)
        return
    }
    for _, netIO := range netIOs {
        fmt.Println(netIO)
    }
}

func main() {
    //getCpuInfo()
    //getLoad()
    //getMemInfo()
    //getHostInfo()
    //getDiskInfo()
    getNetInfo()
}
