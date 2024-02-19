package sever

import (
	"bytes"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetGOOS() string {
	return runtime.GOOS
}

// GetMacHardwareUUID 获取 mac 的系统 uuid
func GetMacHardwareUUID() (string, error) {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := strings.TrimSpace(out.String())
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Hardware UUID") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				//log.Print("parts", strings.TrimSpace(parts[1]))
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("hardware UUID not found")
}

func GetWinHardwareUUID() (string, error) {
	// 初始化 COM
	err := ole.CoInitialize(0)
	if err != nil {
		return "", err
	}
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return "", err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return "", err
	}
	defer wmi.Release()

	// 连接到 WMI
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return "", err
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// 执行查询获取 UUID
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT UUID FROM Win32_ComputerSystemProduct")
	if err != nil {
		return "", err
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	// 获取查询结果
	itemRaw, err := oleutil.CallMethod(result, "ItemIndex", 0)
	if err != nil {
		return "", err
	}
	item := itemRaw.ToIDispatch()
	defer item.Release()

	// 获取 UUID 属性
	uuid, err := oleutil.GetProperty(item, "UUID")
	if err != nil {
		return "", err
	}

	uuidStr := uuid.ToString()

	// 检查 UUID 是否为特定值
	if uuidStr == "03000200-0400-0500-0006-000700080009" {
		// 获取主MAC地址
		interfaces, err := net.Interfaces()
		if err != nil {
			return "", err
		}

		for _, interf := range interfaces {
			if interf.HardwareAddr.String() != "" && (interf.Flags&net.FlagLoopback) == 0 {
				return interf.HardwareAddr.String(), nil
			}

		}

		return "", fmt.Errorf("no MAC address found")
	}

	return uuidStr, nil
}
