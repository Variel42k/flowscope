package inventory

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/flowscope/flowscope/internal/model"
	"gopkg.in/yaml.v3"
)

type Index struct {
	assets []assetEntry
}

type assetEntry struct {
	asset model.InventoryAsset
	net   *net.IPNet
}

func Load(path string, format string) (*Index, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var items []model.InventoryAsset
	switch strings.ToLower(format) {
	case "json":
		if err := json.Unmarshal(b, &items); err != nil {
			return nil, err
		}
	default:
		if err := yaml.Unmarshal(b, &items); err != nil {
			return nil, err
		}
	}
	idx := &Index{}
	for _, item := range items {
		_, ipnet, err := net.ParseCIDR(item.CIDR)
		if err != nil {
			return nil, fmt.Errorf("invalid inventory cidr %s: %w", item.CIDR, err)
		}
		idx.assets = append(idx.assets, assetEntry{asset: item, net: ipnet})
	}
	return idx, nil
}

func (i *Index) Match(ip string) *model.InventoryAsset {
	if i == nil {
		return nil
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil
	}
	for _, item := range i.assets {
		if item.net.Contains(parsed) {
			copy := item.asset
			return &copy
		}
	}
	return nil
}

func (i *Index) Assets() []model.InventoryAsset {
	if i == nil {
		return nil
	}
	out := make([]model.InventoryAsset, 0, len(i.assets))
	for _, a := range i.assets {
		out = append(out, a.asset)
	}
	return out
}
