package decoder

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/util"
)

func (d *Decoder) mapNetFlowField(item *model.FlowRecord, fieldType uint16, raw any) {
	switch fieldType {
	case 1:
		item.Bytes = asUint64(raw)
	case 2:
		item.Packets = asUint64(raw)
	case 4:
		item.IPProtocol = uint8(asUint64(raw))
		item.L4ProtocolName = util.ProtocolName(item.IPProtocol)
	case 6:
		item.TCPFlags = uint16(asUint64(raw))
	case 7:
		item.SrcPort = uint16(asUint64(raw))
	case 8, 27:
		item.SrcIP = asIP(raw)
	case 10:
		item.InputInterface = uint32(asUint64(raw))
	case 11:
		item.DstPort = uint16(asUint64(raw))
	case 12, 28:
		item.DstIP = asIP(raw)
	case 14:
		item.OutputInterface = uint32(asUint64(raw))
	case 16:
		item.SrcASN = uint32(asUint64(raw))
	case 17:
		item.DstASN = uint32(asUint64(raw))
	case 34:
		item.SamplerRate = uint32(asUint64(raw))
	case 58, 59:
		item.VLANID = uint16(asUint64(raw))
	case 61:
		dir := asUint64(raw)
		if dir == 1 {
			item.FlowDirection = "ingress"
		} else if dir == 2 {
			item.FlowDirection = "egress"
		}
	case 152:
		if t := asTimeMilli(raw); !t.IsZero() {
			item.TimestampStart = t
		}
	case 153:
		if t := asTimeMilli(raw); !t.IsZero() {
			item.TimestampEnd = t
		}
	}
}

func asUint64(v any) uint64 {
	switch x := v.(type) {
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	case int8:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int16:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int32:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int64:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case int:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case float32:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case float64:
		if x < 0 {
			return 0
		}
		return uint64(x)
	case []byte:
		var out uint64
		for _, b := range x {
			out = (out << 8) | uint64(b)
		}
		return out
	default:
		return 0
	}
}

func asIP(v any) string {
	switch x := v.(type) {
	case net.IP:
		return x.String()
	case []byte:
		return net.IP(x).String()
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Array && rv.Type().Elem().Kind() == reflect.Uint8 {
			b := make([]byte, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				b[i] = byte(rv.Index(i).Uint())
			}
			return net.IP(b).String()
		}
		parsed := net.ParseIP(fmt.Sprint(v))
		if parsed == nil {
			return ""
		}
		return parsed.String()
	}
}

func asTimeMilli(v any) time.Time {
	ms := asUint64(v)
	if ms == 0 {
		return time.Time{}
	}
	return time.UnixMilli(int64(ms)).UTC()
}
