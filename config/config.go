package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	local = "local"
	minio = "minio"
)

var config Config

// Get 获取初始化后的全局配置
func Get() *Config {
	return &config
}

type Config struct {
	Image   *Image   `yaml:"image"`
	HTTP    *HTTP    `yaml:"http"`
	Storage *Storage `yaml:"storage"`
	Dev     bool     `yaml:"dev"`
}

type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	TLS  bool   `yaml:"tls"`
}

type Storage struct {
	// local, minio
	Type string `yaml:"type"`

	// local config
	Path string `yaml:"path"`

	// minio config
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	TLS       bool   `yaml:"tls"`
}

type Image struct {
	Raw    bool     `yaml:"raw"`
	Preset []string `yaml:"preset"`
}

func Init(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return _init(&config)
}

func _init(config *Config) error {
	imageConfig := defaultImageConfig()
	HTTPConfig := defaultHTTPConfig()
	storageConfig := defaultStorageConfig()

	// 图片方面的配置
	if config.Image == nil {
		config.Image = imageConfig
	} else {
		if config.Image.Preset == nil {
			config.Image.Preset = imageConfig.Preset
		}
	}

	// HTTP服务方面的配置
	if config.HTTP == nil {
		config.HTTP = HTTPConfig
	} else {
		if config.HTTP.Host == "" {
			config.HTTP.Host = HTTPConfig.Host
		}
		if config.HTTP.Port == 0 {
			config.HTTP.Port = HTTPConfig.Port
		}
		if config.HTTP.Cert == "" {
			config.HTTP.Cert = HTTPConfig.Cert
		}
		if config.HTTP.Key == "" {
			config.HTTP.Key = HTTPConfig.Key
		}
	}

	// 存储方面的配置
	if config.Storage == nil {
		config.Storage = storageConfig
		return nil
	}
	switch config.Storage.Type {
	case local:
		if config.Storage.Path == "" {
			config.Storage.Path = storageConfig.Path
		}
		return nil
	case minio:
		if config.Storage.Endpoint == "" || config.Storage.AccessKey == "" || config.Storage.SecretKey == "" {
			return errors.New("please entry the storage config completely")
		}
		if config.Storage.Bucket == "" {
			config.Storage.Bucket = "gallery"
		}
		return nil
	default:
		return errors.New("the storage type is not supported")
	}
}

func defaultImageConfig() *Image {
	return &Image{
		Raw: true,
		Preset: []string{
			"480*480",
			"1080*1080",
		},
	}
}

func defaultHTTPConfig() *HTTP {
	return &HTTP{
		Host: "0.0.0.0",
		Port: 9779,
		Cert: "cert.pem",
		Key:  "key.pem",
	}
}

func defaultStorageConfig() *Storage {
	return &Storage{
		Type: local,
		Path: "./gallery",
	}
}
