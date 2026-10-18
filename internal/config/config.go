package config

import "github.com/zeromicro/go-zero/zrpc"

type UploadConf struct {
	PrivateStorePath    string
	PublicStorePath     string
	ServerURL           string
	MaxImageSize        int64
	MaxVideoSize        int64
	MaxAudioSize        int64
	MaxOtherSize        int64
	DeleteFileWithCloud bool
}

type Config struct {
	zrpc.RpcServerConf
	DatabaseDSN string
	RedisConf   struct {
		Host string
		Type string
		Pass string
	}
	UploadConf UploadConf
}
