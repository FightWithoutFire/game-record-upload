package main

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var userHomeDir string
var machineId string

type OssClient interface {
	upload(info GameInfo)
	download(info GameInfo)
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	fmt.Printf("%s \n", userHomeDir)
	abs, err := filepath.Abs("./internal/app/oss_ali.config")
	if err != nil {
		panic("open ali oss config failure")
	}
	file, err := os.Open(abs)
	reader := bufio.NewReader(file)
	config := make(map[string]string)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic("read line failure")
		}
		c := strings.Split(string(line), ":")
		config[c[0]] = strings.TrimSpace(c[1])
	}
	client, err := oss.New(config["endpoint"], config["id"], config["secret"])

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取存储空间。
	DefaultBucket, err = client.Bucket(config["bucket"])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

}
