package cli

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/prodemmi/kodo/core/entities"
)

func ShowServerInfo(url string, config *entities.Config) {

	banner := figure.NewFigure("KODO", "o8", true)
	color.Cyan(banner.String())

	fmt.Println()

	fmt.Println(color.GreenString("======================================"))

	if config.Flags.Config != ".kodo" && config.Flags.Config != "./.kodo" {
		fmt.Println(color.YellowString("▶ Config Path: %s", config.Flags.Config))
	}

	if config.Flags.Investor {
		fmt.Println(color.BlueString("▶ Investor Mode: Enabled"))
	}

	fmt.Println(color.GreenString("▶ Running at: %s", url))

	fmt.Println(color.GreenString("======================================"))
}

func OpenBrowser(url string) error {
	var cmdName string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmdName = "open"
		args = []string{url}
	case "windows":
		cmdName = "cmd"
		args = []string{"/C", "start", url}
	default:
		cmdName = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmdName, args...).Start()
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func PrintHelp() {
	banner := figure.NewFigure("KODO ", "slant", true)
	color.Cyan(banner.String())
	fmt.Println()

	fmt.Println(color.YellowString("A source-code kanban and note app"))
	fmt.Println(color.GreenString("--------------------------------------------------"))
	fmt.Println(color.WhiteString("Usage:"))
	fmt.Println(color.WhiteString("  kodo [flags]"))
	fmt.Println()
	fmt.Println(color.WhiteString("Available Flags:"))
	fmt.Println(color.WhiteString("  -c, --config <path>     Path to config file (default .kodo)"))
	fmt.Println(color.WhiteString("  -i, --investor          Run in investor mode (default false)"))
	fmt.Println(color.WhiteString("  -h, --help              Show this help message"))
	fmt.Println(color.GreenString("--------------------------------------------------"))
	fmt.Println()
}
