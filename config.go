package blink

import (
	"fmt"
	"github.com/mzky/weblink/internal/log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// 临时文件夹，用于释放 DLL 以及其他临时文件
	tempPath string
	// dll文件路径，非绝对路径将在临时文件夹内创建
	dllFile string
	// 设置storage本地文件目录
	storagePath string
	// 设置cookie文件名
	cookieFile string
}

func NewConfig(setups ...func(*Config)) (*Config, error) {

	tempPath := filepath.Join(os.TempDir(), "mini-blink")

	conf := &Config{
		tempPath:    tempPath,
		dllFile:     "blink.dll",
		storagePath: "LocalStorage",
		cookieFile:  "cookie.dat",
	}

	for _, setup := range setups {
		setup(conf)
	}

	log.Debug("临时文件夹：%s", conf.tempPath)
	if err := os.MkdirAll(conf.tempPath, 0644); err != nil {
		return nil, fmt.Errorf("临时文件夹(%s)不存在，且创建不成功，请确认文件夹权限。", conf.tempPath)
	}

	return conf, nil
}

func WithTempPath(path string) func(*Config) {
	return func(conf *Config) {
		conf.tempPath = path
	}
}

func WithDllFile(dllFile string) func(*Config) {
	return func(conf *Config) {
		conf.dllFile = dllFile
	}
}

func WithStoragePath(path string) func(*Config) {
	return func(conf *Config) {
		conf.storagePath = path
	}
}

func WithCookieFile(path string) func(*Config) {
	return func(conf *Config) {
		conf.cookieFile = path
	}
}

func (conf *Config) GetDllFile() string {
	return conf.dllFile
}

func (conf *Config) GetTempPath() string {
	return conf.tempPath
}

func (conf *Config) GetDllFileABS() string {

	if filepath.IsAbs(conf.dllFile) {
		return conf.dllFile
	}

	return filepath.Join(conf.tempPath, conf.dllFile)
}

func (conf *Config) GetStoragePath() string {

	if filepath.IsAbs(conf.storagePath) {
		return conf.storagePath
	}

	return filepath.Join(conf.tempPath, conf.storagePath)
}

func (conf *Config) GetCookieFileABS() string {

	if filepath.IsAbs(conf.cookieFile) {
		return conf.cookieFile
	}

	return filepath.Join(conf.tempPath, conf.cookieFile)
}

func (conf *Config) ParseCookie(cookieFile string) []*http.Cookie {
	file, err := os.ReadFile(cookieFile)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(file), "\r\n")
	cookies := make([]*http.Cookie, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 7 {
			continue
		}

		expiration, _ := strconv.ParseInt(parts[4], 10, 64)
		expires := time.Unix(expiration, 0)

		// 创建http.Cookie实例
		cookie := &http.Cookie{
			Name:     parts[5],
			Value:    parts[6],
			Path:     parts[2],
			Domain:   parts[0],
			HttpOnly: strings.Contains(parts[1], "TRUE"),
			Secure:   strings.Contains(parts[3], "TRUE"),
			Expires:  expires,
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}
