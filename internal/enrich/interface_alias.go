package enrich

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/flowscope/flowscope/internal/model"
	"gopkg.in/yaml.v3"
)

type InterfaceAliasEnricher struct {
	aliases map[uint32]string
}

func NewInterfaceAliasEnricher(path string) (*InterfaceAliasEnricher, error) {
	if strings.TrimSpace(path) == "" {
		return &InterfaceAliasEnricher{aliases: map[uint32]string{}}, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	parsed := map[string]string{}
	if err := yaml.Unmarshal(b, &parsed); err != nil {
		return nil, err
	}
	out := &InterfaceAliasEnricher{aliases: map[uint32]string{}}
	for k, v := range parsed {
		i, err := strconv.Atoi(strings.TrimSpace(k))
		if err != nil {
			continue
		}
		out.aliases[uint32(i)] = v
	}
	return out, nil
}

func (e *InterfaceAliasEnricher) Enrich(_ context.Context, record *model.FlowRecord) {
	if e == nil {
		return
	}
	if alias := e.aliases[record.InputInterface]; alias != "" {
		record.InputIfAlias = alias
	}
	if alias := e.aliases[record.OutputInterface]; alias != "" {
		record.OutputIfAlias = alias
	}
}
