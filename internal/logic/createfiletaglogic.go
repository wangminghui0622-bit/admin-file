package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-file/file"
	"github.com/wangminghui0622-bit/admin-file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFileTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFileTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFileTagLogic {
	return &CreateFileTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateFileTagLogic) CreateFileTag(in *file.FileTagInfo) (*file.BaseIDResp, error) {
	if in.Name == nil || *in.Name == "" {
		return nil, errorx.New(errorx.InvalidArgument, "name is required")
	}
	status := int32(1)
	if in.Status != nil {
		status = int32(*in.Status)
	}
	row := &model.FmsFileTag{
		Status: &status,
		Name:   *in.Name,
		Remark: in.Remark,
	}
	if err := l.svcCtx.Query.FmsFileTag.WithContext(l.ctx).Create(row); err != nil {
		return nil, dbErr(err)
	}
	return &file.BaseIDResp{Id: uint64(row.ID), Msg: successMsg}, nil
}
