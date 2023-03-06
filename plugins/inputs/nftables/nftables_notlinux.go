//go:build !linux

package nftables

import (
	_ "embed"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type Nftables struct {
	Log telegraf.Logger `toml:"-"`
}

func (n *Nftables) Init() error {
	n.Log.Warn("current platform is not supported")
	return nil
}
func (*Nftables) SampleConfig() string                { return sampleConfig }
func (*Nftables) Gather(_ telegraf.Accumulator) error { return nil }

func init() {
	inputs.Add("nftables", func() telegraf.Input {
		return &Nftables{}
	})
}
