package main

import (
	"fmt"
	"os/user"
	"path/filepath"
	lib "stunel/lib"
	"sync"

	"github.com/spf13/viper"
)

type Tunnel struct {
	PrivateKeyPath string
	User           string
	ServerAddr     string
	RemoteAddr     string
	LocalAddr      string
}

func getAbsolutePathPath(rawPath string) (absolutePath string, err error) {
	if rawPath[0] == []byte("~")[0] {
		// 获取当前用户信息
		var currentUser *user.User
		currentUser, err = user.Current()
		if err != nil {
			return
		}
		// 获取用户的主目录
		homeDir := currentUser.HomeDir
		absolutePath = filepath.Join(homeDir, rawPath[2:])
	} else {
		absolutePath = rawPath
	}
	return
}

func getTunnels() (tunnels []*Tunnel, err error) {
	// 读取配置
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("解析配置文件./config.yaml失败:", err.Error())
		return
	}
	viper.UnmarshalKey("tunnels", &tunnels)
	return
}

func logStr(t *Tunnel, msg string) {
	fmt.Printf("%s->%s %s\n", t.LocalAddr, t.RemoteAddr, msg)
}

func main() {
	ts, err := getTunnels()
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	for _, t := range ts {
		keyPath, err := getAbsolutePathPath(t.PrivateKeyPath)
		if err != nil {
			logStr(t, fmt.Sprintf("get key path fail %s", err.Error()))
			return
		}
		wg.Add(1)
		lib.Tunnel(t.User, keyPath, t.ServerAddr, t.RemoteAddr, t.LocalAddr, func() {
			logStr(t, "start")
		}, func() {
			logStr(t, "close")
			wg.Done()
		})
	}
	wg.Wait()

}
