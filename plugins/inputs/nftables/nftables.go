//go:build linux

package nftables

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type NftOutput struct {
	Nftables []struct {
		Rule struct {
			Family string `json:"family"`
			Table  string `json:"table"`
			Chain  string `json:"chain"`
			Expr   []struct {
				Counter struct {
					Packets uint64 `json:"packets"`
					Bytes   uint64 `json:"bytes"`
				} `json:"counter,omitempty"`
			} `json:"expr"`
		} `json:"rule,omitempty"`
	} `json:"nftables"`
}

type Nftables struct {
	UseSudo   bool
	NftBinary string
}

func (*Nftables) SampleConfig() string {
	return sampleConfig
}

func (n *Nftables) Gather(acc telegraf.Accumulator) error {
	command := []string{n.NftBinary, "--json", "--numeric", "list", "ruleset"}
	if n.UseSudo {
		command = append([]string{"sudo", "--non-interactive"}, command...)
	}

	execCmd := exec.Command(command[0], command[1:]...)
	out, err := internal.StdOutputTimeout(execCmd, 5*time.Second)
	if err != nil {
		return fmt.Errorf(
			"failed to run command %s: %w - %s", strings.Join(command, " "), err, string(out),
		)
	}

	data := NftOutput{}
	err = json.Unmarshal(out, &data)
	if err != nil {
		return fmt.Errorf(
			"failed to unmarshall JSON: %w", err,
		)
	}

	for _, item := range data.Nftables {
		tags := map[string]string{
			"chain":  item.Rule.Chain,
			"family": item.Rule.Family,
			"table":  item.Rule.Table,
		}
		for _, expr := range item.Rule.Expr {
			fields := map[string]any{
				"bytes":   expr.Counter.Bytes,
				"packets": expr.Counter.Packets,
			}
			acc.AddFields("nftables", fields, tags)
		}

	}

	return nil
}

func init() {
	inputs.Add("nftables", func() telegraf.Input {
		return &Nftables{
			NftBinary: "/usr/bin/nft",
		}
	})
}
