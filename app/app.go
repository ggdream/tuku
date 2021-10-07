package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/tuku/config"
	"github.com/ggdream/tuku/controller"
	"github.com/ggdream/tuku/lib/fs"
)

type TuKu struct {
	http *gin.Engine
	fs   fs.FS
}

func New(configPath string) (*TuKu, error) {
	if err := config.Init(configPath); err != nil {
		return nil, err
	}
	globalConfig := config.Get()

	if !globalConfig.Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	var fsIns fs.FS
	var err error
	if globalConfig.Storage.Type == "local" {
		fsIns, err = fs.NewLocal(globalConfig.Storage.Path)
	} else {
		fsIns, err = fs.NewMinIO(globalConfig.Storage.Endpoint, globalConfig.Storage.Bucket, globalConfig.Storage.AccessKey, globalConfig.Storage.SecretKey, globalConfig.Storage.TLS)
	}
	if err != nil {
		return nil, err
	}

	tuKu := &TuKu{
		http: gin.New(),
		fs:   fsIns,
	}
	tuKu.init()

	return tuKu, nil
}

func (k *TuKu) Run() error {
	globalConfig := config.Get()
	addr := fmt.Sprintf("%s:%d", globalConfig.HTTP.Host, globalConfig.HTTP.Port)

	if !globalConfig.HTTP.TLS {
		return k.http.Run(addr)
	} else {
		return k.http.RunTLS(addr, globalConfig.HTTP.Cert, globalConfig.HTTP.Key)
	}
}

func (k *TuKu) init() {
	k.setMiddle()
	k.setRoutes()
}

// setMiddle 注册中间件
func (k *TuKu) setMiddle() {
	k.http.Use(gin.Logger())
	k.http.Use(gin.Recovery())
}

// setRoutes 注册HTTP服务路由
func (k *TuKu) setRoutes() {
	api := k.http.Group("/api")
	api.GET("/image/:object", controller.GetImage(k.fs))
	api.POST("/upload", controller.SimpleUpload(k.fs))
	api.POST("/uploads", controller.MultipleUpload(k.fs))
}
