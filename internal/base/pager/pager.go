package pager

import "loginhub/internal/base/constant"

type PageCond struct {
	Num  int64 `json:"pn" form:"pn"`
	Size int64 `json:"ps" form:"ps"`
}

func (p *PageCond) Validate() {
	if p.Num < constant.MinNum {
		p.Num = constant.DefaultNum
	}
	if p.Size < constant.MinSize {
		p.Size = constant.DefaultSize
	}

	if p.Size > constant.MaxSize {
		p.Size = constant.MaxSize
	}
}

func (p PageCond) GetOffset() int64 {
	return (p.Num - 1) * p.Size
}

type PageResp[T any] struct {
	Count int64 `json:"count"`
	List  []T   `json:"list"`
}
