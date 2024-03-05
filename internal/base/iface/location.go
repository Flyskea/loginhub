package iface

import (
	"context"
	"strings"
)

type RegionBasic struct {
	Country  string
	Province string
	City     string
}

func (r *RegionBasic) String() string {
	res := make([]string, 0)
	if r.Country != "" {
		res = append(res, r.Country)
	}
	if r.Province != "" {
		res = append(res, r.Province)
	}
	if r.City != "" {
		res = append(res, r.City)
	}
	return strings.Join(res, "-")
}

type IP2RegionSearcher interface {
	RegionBasicByIPV4(ctx context.Context, ip string) (*RegionBasic, error)
}
