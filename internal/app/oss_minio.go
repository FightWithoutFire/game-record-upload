package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var GoClient *minio.Client
var bucketName string
var MClient *MinioClient

type MinioClient struct {
	OssClient
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	fmt.Printf("%s \n", userHomeDir)
	abs, err := filepath.Abs("./internal/app/oss_minio.config")
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

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
		// Set this value so that the underlying transport round-tripper
		// doesn't try to auto decode the body of objects with
		// content-encoding set to `gzip`.
		//
		// Refer:
		//    https://golang.org/src/net/http/transport.go?h=roundTrip#L1843
		DisableCompression: true,
	}
	fmt.Printf("config:  \n %v \n", config)

	GoClient, err = minio.New(config["endpoint"], &minio.Options{
		Creds:     credentials.NewStaticV4(config["id"], config["secret"], ""),
		Secure:    false,
		Transport: transport,
	})

	bucketName = config["bucket"]

	fmt.Printf("client: \n %v \n", GoClient.EndpointURL())

	fmt.Printf("minio is online %v \n", GoClient.IsOnline())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

}

func (client *MinioClient) upload(info GameInfo) {
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

	format := time.Now().Format("2006-01-02 15:04:05")
	for _, file := range files {
		fmt.Printf("path %s \n", recordPath+"/"+file.Name())
		openFile, err := os.Open(recordPath + "/" + file.Name())
		defer openFile.Close()
		if err != nil {
			fmt.Printf("open file err %v \n", err)
			continue
		}

		fileKey := objectKey + "/" + file.Name()
		fmt.Printf(" %s \n", fileKey)
		fmt.Printf("upload %s \n", file.Name())
		_, err = openFile.Stat()

		fmt.Printf("bucket name : %v \n", bucketName)
		_, err = GoClient.FPutObject(context.Background(), bucketName, info.GameName+"/"+format+"/"+file.Name(), recordPath+"/"+file.Name(), minio.PutObjectOptions{})
		if err != nil {
			fmt.Printf("upload error %v \n", err)
			continue
		}
	}
	fmt.Printf("upload finished \n")
}

func (client *MinioClient) download(info GameInfo) {

	prefix := info.GameName + "/"

	objects := GoClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{Prefix: prefix})
	dest := strings.Replace(info.RecordPath, "{USERPATH}", userHomeDir, -1)
	for o := range objects {
		split := strings.Split(o.Key, "/")
		fileName := split[len(split)-1]
		err := GoClient.FGetObject(context.Background(), "", o.Key, dest+"\\"+fileName, minio.GetObjectOptions{})
		if err != nil {
			return
		}
	}
}
