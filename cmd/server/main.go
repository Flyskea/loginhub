package main

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/tx7do/kratos-transport/transport/gin"

	"loginhub/internal/base/validator"
	"loginhub/internal/conf"
	"loginhub/pkg/logx"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	rootCmd.Version = Version
	rootCmd.Flags().StringVarP(&flagconf, "conf", "c", "./configs/config-example.yaml", "config path")
}

func runApp() {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	config := &conf.Bootstrap{}
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Scan(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("config: %+v\n", config)
	conf.Init(config)
	validator.InitValidator()
	app, clean, err := initApp(
		config.GetLog(),
		config.GetData().GetDatabase(),
		config.GetData().GetRedis(),
		config.GetPassport(),
		config.GetCaptcha(),
		config.GetSmtp(),
		config.Server.GetHttp(),
		config.GetIp2Region(),
	)
	defer clean()
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newApp(
	logger *logx.KratosToSlog,
	httpServer *gin.Server,
) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			httpServer,
		),
	)
}

func main() {
	Execute()
}
