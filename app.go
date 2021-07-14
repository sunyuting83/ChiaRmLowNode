package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

// RunCommand run command
func RunCommand(OS, command string) (k string, err error) {
	var cmd *exec.Cmd
	if OS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("/bin/sh", "-c", command)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return "", err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if len(bytesErr) != 0 {
		return "", errors.New("0")

	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(bytes), nil
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
func GetUserInfo() (homedir string, err error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

func main() {
	var (
		OS       string = runtime.GOOS
		ChiaPath string
		SpiltStr string = "\n"
	)
	if runtime.GOOS == "windows" {
		SpiltStr = "\r\n"
	}
	if OS == "windows" {
		homedir, err := GetUserInfo()
		if err != nil {
			fmt.Println("获取用户目录失败")
			os.Exit(0)
		}
		ChiaPath = strings.Join([]string{homedir, `AppData\Local\chia-blockchain`}, `\`)
	}
	if !IsDir(ChiaPath) {
		fmt.Println("获取Chia运行目录失败")
		os.Exit(0)
	}
	ChiaPath = "/lib/chia-blockchain/resources/app.asar.unpacked/daemon/"
	command := strings.Join([]string{ChiaPath, "chia show -c"}, "")
	f, _ := RunCommand(OS, command)
	list := strings.Split(f, SpiltStr)

	commanda := strings.Join([]string{ChiaPath, "chia show -s"}, "")
	h, _ := RunCommand(OS, commanda)
	hlist := strings.Split(h, SpiltStr)
	var (
		intN int
		intH int
	)
	for _, v := range hlist {
		if strings.Contains(v, "Time:") {
			height := strings.Replace(v, " ", "", -1)

			a := strings.Index(height, "Height:") + 7
			height = height[a:]
			intHeight, _ := strconv.Atoi(height)
			ilan := len(height)
			intH = (ilan / 2) - 1
			var x string
			for i := 0; i < ilan; i++ {
				if i == intH {
					x += "1"
				} else {
					x += "0"
				}
			}
			intA, _ := strconv.Atoi(x)
			intN = intHeight - intA
		}
	}

	for i, v := range list {
		if i >= 2 {
			if len(v) > 0 {
				if strings.HasPrefix(v, "FULL_NODE") {
					height := list[i+1]
					height = strings.Replace(height, " ", "", -1)
					a := strings.Index(height, "Height:") + 7
					b := strings.Index(height, "-Has")
					height = height[a:b]
					intHeight, _ := strconv.Atoi(height)
					if intHeight < intN {
						c := strings.Index(v, "...")
						d := c - 8
						nodeID := v[d:c]
						co := strings.Join([]string{"chia show -r", nodeID}, " ")
						commanda := strings.Join([]string{ChiaPath, co}, "")
						RunCommand(OS, commanda)
						return
					}
				}
			}
		}
	}
}
