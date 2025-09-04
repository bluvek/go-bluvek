package bvutils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"runtime/debug"
	"strings"
)

var (
	ServerAddr  string
	ServerIsTLS bool
)

// GetServerAddr 获取服务地址
func GetServerAddr() string {
	prefix := "http"
	if ServerIsTLS {
		prefix = "https"
	}

	addr := ServerAddr
	addrArr := strings.Split(ServerAddr, ":")
	if addrArr[0] == "" {
		addrArr[0] = GetLocalIP()
		addr = strings.Join(addrArr, ":")
	}

	return fmt.Sprintf("%s://%s", prefix, addr)
}

// SafeGo 运行一个函数，并捕获panic
func SafeGo(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[run func] panic: %v, funcName: %s，stack=%s \n",
				err, GetCallerName(fn), string(debug.Stack()))
		}
	}()

	fn()
}

// GetCallerName 获取调用者的名称
func GetCallerName(caller interface{}) string {
	typ := reflect.TypeOf(caller)

	switch typ.Kind() {
	case reflect.Ptr: // 指针类型
		if typ.Elem().Kind() == reflect.Struct {
			return typ.Elem().Name()
		}
		return fmt.Sprintf("%v", reflect.ValueOf(caller).Elem())
	case reflect.Struct: // 结构体类型
		return typ.Name()
	default: // 其他类型
		return fmt.Sprintf("%v", caller)
	}
}

// GetModuleName 获取当前模块名称
func GetModuleName() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get module name: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}

// 获取操作系统
func GetPlatform(userAgent string) string {
	ua := strings.ToLower(userAgent)

	// 移动端
	if strings.Contains(ua, "android") {
		return "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod") {
		return "iOS"
	}

	// 桌面端
	if strings.Contains(ua, "windows") {
		return "Windows"
	} else if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		return "MacOS"
	} else if strings.Contains(ua, "linux") {
		return "Linux"
	}

	return "Unknown"
}

// 获取浏览器类型
func GetBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg") {
		return "Google Chrome"
	} else if strings.Contains(ua, "edg") {
		return "Microsoft Edge"
	} else if strings.Contains(ua, "firefox") {
		return "Mozilla Firefox"
	} else if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		return "Apple Safari"
	} else if strings.Contains(ua, "opr") || strings.Contains(ua, "opera") {
		return "Opera"
	} else if strings.Contains(ua, "msie") || strings.Contains(ua, "trident") {
		return "Internet Explorer"
	}

	return "Unknown"
}
