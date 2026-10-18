package logic

import (
	"context"
	_ "embed"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// initSQL is the idempotent DDL + seed script for the file service.
//
//go:embed init.sql
var initSQL string

const initKey = "initDatabase:file"

type InitDatabaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// InitDatabase creates the fms tables, seeds sys_apis/casbin_rules/sys_menus
// and marks the database as initialized. Repeated calls are no-ops.
func (l *InitDatabaseLogic) InitDatabase(in *file.Empty) (*file.BaseResp, error) {
	if ok, _ := l.svcCtx.Redis.Exists(initKey); ok {
		return &file.BaseResp{Msg: "already initialized"}, nil
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Exec(initSQL).Error; err != nil {
		logx.Errorw("init sql failed", logx.Field("detail", err.Error()))
		return nil, errorx.New(errorx.ServerError, "初始化数据库失败: "+err.Error())
	}

	if err := l.svcCtx.Redis.Set(initKey, "1"); err != nil {
		return nil, err
	}

	// Providers may have been seeded/imported by another instance; reload the
	// in-memory S3 client group so cloud features work immediately.
	l.svcCtx.CloudStorage.Reload()

	return &file.BaseResp{Msg: successMsg}, nil
}
