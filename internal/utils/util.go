package utils

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

// RandomInt 区域随机整型数字
func RandomInt(min, max int) int {
	randNum := rand.Intn(max-min) + min
	return randNum
}

// RandomIP 生成随机ip
func RandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		RandomInt(1, 255), RandomInt(1, 255), RandomInt(1, 255), RandomInt(1, 255))
}

// GetNameFromUrl 通过url得到名字
func GetNameFromUrl(url string) string {
	arr := strings.Split(url, "/")
	return arr[len(arr)-1]
}

// IsExist 判读文件夹是否存在
func IsExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// ConvertToString 转换字符串编码
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// DownloadFile 下载文件
func DownloadFile(uri string, fileName string, c chan int) error {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		//panic(err)
		return errors.New(fmt.Sprintf("[%s] failed download Error:%X", uri, err))
	}

	req.Header.Set("Accept-Language", `zh-CN,zh;q=0.9`)
	req.Header.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36`)
	req.Header.Set("X-Forwarded-For", RandomIP())
	req.Header.Set("Referer", uri)
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`)
	req.Header.Set("Content-Type", `text/html; charset=UTF-8`)

	http.DefaultClient.Timeout = 10 * time.Second
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		//fmt.Println("failed download ")
		//panic(err)
		return errors.New(fmt.Sprintf("[%s] failed download Error:%X", uri, err))
	}
	if resp.StatusCode != http.StatusOK {
		//fmt.Println("failed download " + uri)
		return errors.New(fmt.Sprintf("[%s] failed download StatusCode:%d", uri, resp.StatusCode))
	}

	defer func() {
		_ = resp.Body.Close()
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		c <- 0
	}()

	localFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("[%s] failed openFile Error:%X", fileName, err))
	}

	if _, err := io.Copy(localFile, resp.Body); err != nil {
		//panic("failed save " + fileName)
		return errors.New(fmt.Sprintf("[%s] failed save Error:%X", fileName, err))
	}

	_ = localFile.Close()

	fmt.Println("success download " + fileName)
	return nil
}

// CheckValidPort 检查给定的端口是否有效
func CheckValidPort(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(l)
	return err == nil
}

// RandomPort 获取本地随机的可用端口
func RandomPort() int {
	var port int
	for {
		port = RandomInt(10000, 20000)
		if CheckValidPort(port) {
			return port
		}
	}
}

// GetOutBoundIP 获取本地IP地址
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
