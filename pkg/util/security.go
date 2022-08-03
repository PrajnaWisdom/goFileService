package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
    "math/rand"
    "time"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	"strings"
)

const (
	defaultPbkdf2Iterations = 150000
    letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
    letterIdxBits = 6
	// 先位移，然后减 1
	letterIdxMask = 1<<letterIdxBits - 1
	// 生成一个 int64,有 63 个 有效 bit
	letterIdxMax = 63 / letterIdxBits

)


var src = rand.NewSource(time.Now().UnixNano())


func GenerateRandomString(n int) string {
	if n <= 0 {
		panic("String length must be positive")
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; i-- {
		if remain == 0 {
			// 生成的随机数用完了，重新生成
			cache, remain = src.Int63(), letterIdxMax
		}
		b[i] = letterBytes[cache&letterIdxMask]
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}


// 生成加密后密码，源自 Python 的 werkzeug 库
func GeneratePasswordHash(password string) (string, error) {
	salt := GenerateRandomString(8)
	h, actualMethod, err := hashInternal("pbkdf2:sha256", salt, password)
	return fmt.Sprintf("%s$%s$%s", actualMethod, salt, h), err
}

func CheckPasswordHash(pwHash, password string) error {
	if strings.Count(pwHash, "$") < 2 {
		return errors.New("加密密码格式不合法")
	}
	// 把字符串分成三部分
	args := strings.SplitN(pwHash, "$", 3)
	tmp, _, err := hashInternal(args[0], args[1], password)
	if err != nil {
		return errors.New("加密错误")
	}
	if args[2] != tmp {
		return errors.New("密码错误")
	}
	return nil
}

func hashInternal(method, salt, password string) (string, string, error) {
	if method == "plain" {
		return password, method, nil
	}
	if strings.HasPrefix(method, "pbkdf2:") {
		args := strings.Split(method[7:], ":")
		method = args[0]
		var iterations int
		// 注意：error 要在这里定义，不然下面执行 strconv.Atoi 函数时，
		// 使用 := 赋值符会导致 iterations 作用域在当前大括号中
		var err error
		if len(args) == 1 {
			iterations = defaultPbkdf2Iterations
		} else if len(args) == 2 {
			iterations, err = strconv.Atoi(args[1])
			if err != nil || iterations <= 0 {
				iterations = defaultPbkdf2Iterations
			}
		} else {
			return "", "", errors.New("invalid number of arguments for PBKDF2")
		}
		actualMethod := fmt.Sprintf("pbkdf2:%s:%d", method, iterations)
		rv := pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, sha256.New)
		// 将 []byte 转换成十六进制小写字符串
		encodedStr := hex.EncodeToString(rv)
		return encodedStr, actualMethod, nil
	} else {
		return "", "", errors.New("don't support other method")
	}
}
