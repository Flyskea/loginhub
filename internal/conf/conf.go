package conf

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

var GConf *Bootstrap

const (
	defaultPath = "./static/ip2region.xdb"
)

func Init(conf *Bootstrap) {
	if conf.GetPassport().GetDeviceCacheTtl().AsDuration() == 0 {
		conf.GetPassport().DeviceCacheTtl = durationpb.New(time.Hour * 24)
	}
	if conf.GetIp2Region() == nil || conf.GetIp2Region().GetPath() == "" {
		conf.Ip2Region = &IP2Region{
			Path: defaultPath,
		}
	}
	GConf = conf
}
