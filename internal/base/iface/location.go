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
	return strings.Join([]string{r.Country, r.Province, r.City}, "-")
}

type IP2RegionSearcher interface {
	RegionBasicByIPV4(ctx context.Context, ip string) (*RegionBasic, error)
}
