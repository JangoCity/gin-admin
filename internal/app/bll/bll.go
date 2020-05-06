package bll

import (
	"context"

	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/pkg/auth"
)

// GetUserID 获取用户ID
func GetUserID(ctx context.Context) string {
	userID, ok := icontext.FromUserID(ctx)
	if ok {
		return userID
	}
	return ""
}

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, transModel *model.Trans, fn TransFunc) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}
	trans, err := transModel.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(icontext.NewTrans(ctx, trans))
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}

// Common 提供统一的业务逻辑处理
type Common struct {
	Demo *Demo
}

// NewCommon 创建统一的业务逻辑处理
func NewCommon(m *model.Common, a auth.Auther) *Common {
	return &Common{
		Demo: NewDemo(m),
	}
}
