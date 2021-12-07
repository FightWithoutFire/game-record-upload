package main

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var AliClient *oss.Client
var DefaultBucket *oss.Bucket

type AliOssClient struct {
	OssClient
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

func (client *AliOssClient) upload(info GameInfo) {
	fmt.Printf("upload started \n")
	dir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	recordPath := strings.Replace(info.RecordPath, "{USERPATH}", dir, -1)

	if err != nil {
		return
	}
	objectKey := info.GameName

	files, err := ioutil.ReadDir(recordPath)
	if err != nil {
		fmt.Printf("list record file error \n")
		return
	}
	for _, file := range files {
		fmt.Printf("upload %s \n", file.Name())
		err = DefaultBucket.PutObjectFromFile(objectKey+"/"+file.Name(), recordPath+"/"+file.Name())
		if err != nil {
			fmt.Printf("upload error %v \n", err)
			return
		}
	}
	fmt.Printf("upload finished \n")
}

func (client *AliOssClient) download(info GameInfo) {
	marker := ""
	prefix := oss.Prefix(info.GameName + "/")
	objects, err := DefaultBucket.ListObjects(oss.Marker(marker), prefix)
	if err != nil {
		return
	}
	dest := strings.Replace(info.RecordPath, "{USERPATH}", userHomeDir, -1)

	for _, o := range objects.Objects {
		split := strings.Split(o.Key, "/")
		fileName := split[len(split)-1]

		err := DefaultBucket.GetObjectToFile(o.Key, dest+"\\"+fileName)
		if err != nil {
			return
		}
	}
}
