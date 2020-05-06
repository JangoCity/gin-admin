package app

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/google/gops/agent"
)

// Object 对象集合
type Object struct {
	Auth  auth.Auther
	Model *model.Common
	Bll   *bll.Common
}

// Init 应用初始化
func Init(ctx context.Context) func() {
	loggerCall, err := InitLogger()
	if err != nil {
		panic(err)
	}

	if c := config.GetGlobalConfig().Monitor; c.Enable {
		err = agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir, ShutdownCleanup: true})
		if err != nil {
			logger.StartSpan(ctx, "开启[agent]服务监听", "ginadmin.Init").Errorf(err.Error())
		}
	}

	InitCaptcha()

	obj, objCall, err := InitObject(ctx)
	if err != nil {
		panic(err)
	}

	app := InitWeb(ctx, obj)
	httpCall := InitHTTPServer(ctx, app)

	return func() {
		if httpCall != nil {
			httpCall()
		}

		if objCall != nil {
			objCall()
		}

		if loggerCall != nil {
			loggerCall()
		}
	}
}

// InitObject 初始化对象数据
func InitObject(ctx context.Context) (*Object, func(), error) {
	obj := &Object{}

	auth, err := InitJWTAuth()
	if err != nil {
		return nil, nil, err
	}
	obj.Auth = auth

	m, storeCall, err := model.Init()
	if err != nil {
		return nil, nil, err
	}
	obj.Model = m
	obj.Bll = bll.NewCommon(obj.Model, obj.Auth)

	return obj, func() {
		if storeCall != nil {
			storeCall()
		}

		if auth != nil {
			err := auth.Release()
			if err != nil {
				logger.StartSpan(ctx, "释放认证资源", "ginadmin.InitObject").Errorf(err.Error())
			}
		}
	}, nil
}
