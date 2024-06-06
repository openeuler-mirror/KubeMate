package util

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"ops-entry/constValue"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

/**
* @Description: 多版本数据合并，新版本数据中未配置的数据则用低版本数据进行填充
* @param mData旧版本数据
* @param mData2新版本数据
* return
*   @resp 待修改的数据
*
 */

func MergeMap(mData, mData2 map[string]string) map[string]string {
	if len(mData2) < 1 {
		return mData
	}
	if len(mData) < 1 {
		return mData2
	}
	for k, val := range mData2 {
		mData[k] = val
	}
	return mData
}

func IsElementInSlice(slice []string, target string) bool {
	if len(slice) == 0 || len(target) == 0 {
		return false
	}

	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

// isValidPath 检查路径是否包含非法字符。
func isValidPath(path string) bool {
	// 检查路径是否为空
	if path == "" {
		return false
	}
	// 检查路径是否包含空字符
	if strings.ContainsRune(path, '\x00') {
		return false
	}

	// 检查路径是否包含非法字符
	for _, r := range path {
		if r < ' ' { // 控制字符
			return false
		}
		if runtime.GOOS == "windows" {
			break
		}

		switch r {
		case '\\', ':', '*', '?', '"', '<', '>', '|':
			return false
		}
	}

	// 检查路径是否是有效的UTF-8字符串
	if !utf8.ValidString(path) {
		return false
	}

	return true
}

/**
* @Description: 文件拷贝
* @param src源文件
* @param dst目标文件，如果目标文件不存在则创建相应的目录
* return
*   @resp
*
 */

func CopyFile(src, dst string) error {
	if len(src) == 0 || len(dst) == 0 {
		return errors.New("empty src or dst")
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return &os.PathError{Op: "copy", Path: src, Err: os.ErrInvalid}
	}

	if sourceFileStat.Size() > constValue.MaxCopyFileSize || sourceFileStat.Size() == 0 {
		return fmt.Errorf("%s is too large or small", src)
	}

	// 创建目标路径（如果它不存在）
	if err := CreatePath(dst); err != nil {
		return fmt.Errorf("dst create failed [err:%s]", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

type Context struct {
	RequestId      string
	ginContextKeys map[string]interface{}
}

func CreateContext(requestId string) Context {
	var c Context

	if len(requestId) > 0 {
		c.RequestId = requestId
	} else {
		timestamp := time.Now().UnixNano()
		rand.Seed(timestamp)
		cur_rand := rand.Intn(31415926)
		c.RequestId = fmt.Sprintf("%d-%d", cur_rand, timestamp)
	}

	return c
}

func (c *Context) P() string {
	return fmt.Sprintf("[%s] ", c.RequestId)
}
