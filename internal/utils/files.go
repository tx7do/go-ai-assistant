package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GetExePath 获取应用程序路径
func GetExePath() string {
	strExePath, _ := os.Executable()
	strCurrentPath := filepath.Dir(strExePath)
	return strCurrentPath
}

// GetCurrentPath 获取当前工作路径
func GetCurrentPath() string {
	strCurrentPath, _ := os.Getwd()
	return strCurrentPath
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDir 创建文件夹
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// RemoveDir 删除文件夹
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// WriteBinaryFile 写入二进制文件
func WriteBinaryFile(filePath string, bytesData []byte) {
	_ = CreateDir(filepath.Dir(filePath))

	fp, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(fp *os.File) {
		err = fp.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(fp)

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, bytesData)

	_, _ = fp.Write(buf.Bytes())
}

// ReadBinaryFile 读取二进制文件
func ReadBinaryFile(filePath string) ([]byte, error) {
	bytesData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytesData, nil
}

// LoadYaml 加载YAML配置文件
func LoadYaml[T interface{}](filePath string, cfg *T) error {
	var err error

	var dataBytes []byte
	if dataBytes, err = os.ReadFile(filePath); err != nil {
		log.Println("读取 yaml 配置文件失败：", err)
		return err
	}

	if err = yaml.Unmarshal(dataBytes, &cfg); err != nil {
		log.Println("解析 yaml 配置文件失败：", err)
		return err
	}

	return nil
}

// SaveYaml 保存YAML配置文件
func SaveYaml[T interface{}](filePath string, cfg *T) error {
	var err error

	var dataBytes []byte
	if dataBytes, err = yaml.Marshal(cfg); err != nil {
		log.Println("序列化 yaml 配置文件 失败：", err)
		return err
	}

	if err = os.WriteFile(filePath, dataBytes, 0644); err != nil {
		log.Println("写入 yaml 配置文件失败：", err)
		return err
	}

	return nil
}
