package conf

import (
	"encoding/json"

	"github.com/Shadowsocks-NET/v2ray-go/v4/common/serial"
	"github.com/Shadowsocks-NET/v2ray-go/v4/proxy/blackhole"
	"github.com/golang/protobuf/proto"
)

type NoneResponse struct{}

func (*NoneResponse) Build() (proto.Message, error) {
	return new(blackhole.NoneResponse), nil
}

type HTTPResponse struct{}

func (*HTTPResponse) Build() (proto.Message, error) {
	return new(blackhole.HTTPResponse), nil
}

type BlackholeConfig struct {
	Response json.RawMessage `json:"response"`
}

func (v *BlackholeConfig) Build() (proto.Message, error) {
	config := new(blackhole.Config)
	if v.Response != nil {
		response, _, err := configLoader.Load(v.Response)
		if err != nil {
			return nil, newError("Config: Failed to parse Blackhole response config.").Base(err)
		}
		responseSettings, err := response.(Buildable).Build()
		if err != nil {
			return nil, err
		}
		config.Response = serial.ToTypedMessage(responseSettings)
	}

	return config, nil
}

var configLoader = NewJSONConfigLoader(
	ConfigCreatorCache{
		"none": func() interface{} { return new(NoneResponse) },
		"http": func() interface{} { return new(HTTPResponse) },
	},
	"type",
	"")
