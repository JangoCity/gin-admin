package model

import (
	"context"
	"os"
	"path/filepath"

	"github.com/LyricTian/gin-admin/internal/app/config"
	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/internal/app/model/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"github.com/jinzhu/gorm"
)

// Init 初始化存储
func Init() (*Common, func(), error) {
	db, err := initGorm()
	if err != nil {
		return nil, nil, err
	}

	storeCall := func() {
		db.Close()
	}

	SetTablePrefix(config.GetGlobalConfig().Gorm.TablePrefix)
	err = AutoMigrate(db)
	if err != nil {
		return nil, nil, err
	}

	m := NewModel(db)
	return m, storeCall, nil
}

// initGorm 实例化gorm存储
func initGorm() (*gormplus.DB, error) {
	cfg := config.GetGlobalConfig()

	var dsn string
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		os.MkdirAll(filepath.Dir(dsn), 0777)
	case "postgres":
		dsn = cfg.Postgres.DSN()
	default:
		return nil, errors.New("unknown db")
	}

	return gormplus.New(&gormplus.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
	})
}

// SetTablePrefix 设定表名前缀
func SetTablePrefix(prefix string) {
	entity.SetTablePrefix(prefix)
}

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gormplus.DB, fn func(context.Context) error) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	transModel := NewTrans(db)
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

// WrapPageQuery 包装带有分页的查询
func WrapPageQuery(db *gorm.DB, pp *schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp != nil {
		total, err := gormplus.Wrap(db).FindPage(db, pp.PageIndex, pp.PageSize, out)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{
			Total: total,
		}, nil
	}

	result := db.Find(out)
	return nil, result.Error
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gormplus.DB) error {
	return db.AutoMigrate(
		new(entity.Demo),
	).Error
}

// Common 提供统一的存储接口
type Common struct {
	Trans *Trans
	Demo  *Demo
}

// NewModel 创建gorm存储，实现统一的存储接口
func NewModel(db *gormplus.DB) *Common {
	return &Common{
		Trans: NewTrans(db),
		Demo:  NewDemo(db),
	}
}
