package drain

import "io"

//go:generate go run github.com/Shadowsocks-NET/v2ray-go/v4/common/errors/errorgen

type Drainer interface {
	AcknowledgeReceive(size int)
	Drain(reader io.Reader) error
}
