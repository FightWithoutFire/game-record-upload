package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type GameInfo struct {
	GameName   string
	LaunchFile string
	RecordPath string
}

var paths map[string]GameInfo
var userPath string
var pids []int32
var ossClient OssClient

func init() {

	ossClient = MClient
	userPath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("get user path error \n")
		return
	}
	paths = make(map[string]GameInfo)
	absPath, err := filepath.Abs("./internal/app/gamepath.txt")
	if err != nil {
		panic("get path failure")
	}
	file, err := os.Open(absPath)
	defer file.Close()
	if err != nil {
		panic("gamepath opening failure")
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				break
			}
			panic("gamepath parse failure")
		}
		if len(line) == 0 {
			break
		}
		c := strings.Split(string(line), ",")
		recordPath := strings.Replace(c[2], "{USERPATH}", userPath, -1)
		fmt.Printf("%v \n", recordPath)

		info := GameInfo{
			GameName:   c[0],
			LaunchFile: c[1],
			RecordPath: recordPath,
		}
		paths[c[1]] = info
		fmt.Printf("%v: %v \n", c[1], info)
	}
}

func main() {

	//processes, err := process.Processes()
	//if err != nil {
	//	panic("read process failure")
	//}
	//for _, p := range processes {
	//	name, err := p.Name()
	//	if err != nil {
	//		panic("read process failure")
	//	}
	//	info, exist := paths[name]
	//	if exist {
	//		pid := p.Pid
	//		notInPids := true
	//		for _, num := range pids {
	//			if num == pid {
	//				notInPids = false
	//			}
	//
	//		}
	//		if notInPids {
	//			pids = append(pids, pid)
	//			fmt.Printf("name: %s \n", name)
	//			go func(pid int32, info GameInfo) {
	//				for {
	//					exists, err := process.PidExists(pid)
	//					if err != nil {
	//						panic("check pid exist failure")
	//					}
	//					if !exists {
	//						index := -1
	//						ossClient.upload(info)
	//						for i, num := range pids {
	//							if num == pid {
	//								index = i
	//								notInPids = false
	//							}
	//
	//						}
	//
	//						pids = append(pids[:index], pids[index+1:]...)
	//						break
	//					}
	//					time.Sleep(5 * time.Second)
	//				}
	//			}(pid, info)
	//		}
	//	}
	//}

	//filepath.Walk("d:/", func(path string, info fs.FileInfo, err error) error {
	filepath.Walk("D:\\gameSoftware\\steam", func(path string, info fs.FileInfo, err error) error {
		fmt.Printf("check %s \n", info.Name())

		gameInfo, exist := paths[info.Name()]

		if exist {
			fmt.Printf("find game %s \n", info.Name())

			fmt.Printf("dest %s \n", gameInfo.RecordPath)

			fmt.Printf(" %s is %s`s record dir \n", gameInfo.RecordPath, gameInfo.GameName)
			files, err := ioutil.ReadDir(gameInfo.RecordPath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				ossClient.download(gameInfo)
			} else {
				ossClient.upload(gameInfo)
			}

		}

		return nil
	})

}
