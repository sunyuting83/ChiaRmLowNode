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
	"time"
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

func GetLinuxPath() string {
	cmd := exec.Command("uname", "-a")
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()
	cmd.Start()
	opBytes, _ := ioutil.ReadAll(stdout)
	if strings.Contains(string(opBytes), "Ubuntu") || strings.Contains(string(opBytes), "Debain") {
		return "/lib/chia-blockchain/resources/app.asar.unpacked/daemon/"
	}
	return ""
}

func Check(OS, ChiaPath, SpiltStr string) {
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
					fmt.Println(strings.Join([]string{"Remove nodeID", nodeID}, "-"))
					return
				}
			} else {
				if !strings.HasPrefix(v, "127.0.0.1") {
					c := strings.Index(v, "...")
					d := c - 8
					nodeID := v[d:c]
					co := strings.Join([]string{"chia show -r", nodeID}, " ")
					commanda := strings.Join([]string{ChiaPath, co}, "")
					RunCommand(OS, commanda)
					fmt.Println(strings.Join([]string{"Remove nodeID", nodeID}, "-"))
					return
				}
			}
		}
	}
}

func main() {
	var (
		OS       string = runtime.GOOS
		ChiaPath string
		SpiltStr string = "\n"
		Sleep    int    = 120
	)
	if len(os.Args) >= 2 {
		s, _ := strconv.Atoi(os.Args[1])
		if s < 0 {
			Sleep = s
		}
	}
	if OS == "linux" {
		p := GetLinuxPath()
		if len(p) > 0 {
			ChiaPath = p
		} else {
			if len(os.Args) >= 3 {
				ChiaPath = os.Args[2]
			} else {
				time.Sleep(time.Duration(10) * time.Second)
				fmt.Println("请将chia运行目录填入第三个参数")
				os.Exit(0)
			}
		}
	}
	if OS == "windows" {
		SpiltStr = "\r\n"
		homedir, err := GetUserInfo()
		if err != nil {
			time.Sleep(time.Duration(10) * time.Second)
			fmt.Println("获取用户目录失败")
			os.Exit(0)
		}
		rootPath := strings.Join([]string{homedir, `AppData\Local\chia-blockchain`}, `\`)
		files, _ := ioutil.ReadDir(rootPath)
		var versionNumber []string
		for _, f := range files {
			if strings.Contains(f.Name(), "app-") {
				versionNumber = append(versionNumber, string(f.Name()))
			}
		}
		ChiaPath = strings.Join([]string{rootPath, versionNumber[0], `resources\app.asar.unpacked\daemon`}, `\`)
		if len(versionNumber) > 1 {
			n := len(versionNumber) - 1
			ChiaPath = strings.Join([]string{rootPath, versionNumber[n], `resources\app.asar.unpacked\daemon`}, `\`)
		}
		ChiaPath = strings.Join([]string{ChiaPath, `\`}, "")
	}
	if !IsDir(ChiaPath) {
		fmt.Println("获取Chia运行目录失败")
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(Sleep))
	go func() {
		for range ticker.C {
			Check(OS, ChiaPath, SpiltStr)
		}
		ch <- 1
	}()
	<-ch

}
