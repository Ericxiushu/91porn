package main

import (
	"fmt"
	"math/rand"
	"os"
)

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

/// 区域随机整型数字
func random_int(min, max int) int {
	randNum := rand.Intn(max-min) + min
	return randNum
}

/// 生成随机ip
func random_ip() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		random_int(1, 255), random_int(1, 255), random_int(1, 255), random_int(1, 255))
}

// 判读文件夹是否存在
func isExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
