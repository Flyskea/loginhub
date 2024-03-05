package location

import (
	"context"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"

	"loginhub/internal/base/iface"
	"loginhub/internal/conf"
)

var _ iface.IP2RegionSearcher = (*XDBIPRegionSearcher)(nil)

type XDBIPRegionSearcher struct {
	searcher *xdb.Searcher
}

func (s *XDBIPRegionSearcher) RegionBasicByIPV4(ctx context.Context, ip string) (*iface.RegionBasic, error) {
	info, err := s.searcher.SearchByStr(ip)
	if err != nil {
		return nil, err
	}
	region := strings.Split(info, "|")
	var country, province, city string
	if len(region) > 0 {
		country = region[0]
	}
	if len(region) > 2 {
		province = region[2]
	}
	if len(region) > 3 {
		city = region[3]
	}
	if country == "0" {
		country = ""
	}
	if province == "0" {
		province = ""
	}
	if city == "0" {
		city = ""
	}
	return &iface.RegionBasic{
		Country:  country,
		Province: province,
		City:     city,
	}, nil
}

func NewIPRegionSearcher(conf *conf.IP2Region) (*XDBIPRegionSearcher, error) {
	buff, err := xdb.LoadContentFromFile(conf.GetPath())
	if err != nil {
		return nil, err
	}
	handle, err := xdb.NewWithBuffer(buff)
	if err != nil {
		return nil, err
	}
	return &XDBIPRegionSearcher{
		searcher: handle,
	}, nil
}
