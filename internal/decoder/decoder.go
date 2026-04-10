package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/util"
	"github.com/netsampler/goflow2/v2/decoders/netflow"
	"github.com/netsampler/goflow2/v2/decoders/netflowlegacy"
	sflowdec "github.com/netsampler/goflow2/v2/decoders/sflow"
)

type Decoder struct {
	templates netflow.NetFlowTemplateSystem
}

func New() *Decoder {
	return &Decoder{templates: netflow.CreateTemplateSystem()}
}

func (d *Decoder) Decode(sourceType string, payload []byte, exporterIP string) ([]model.FlowRecord, error) {
	switch sourceType {
	case "netflow_v5":
		return d.decodeNetFlowV5(payload, exporterIP)
	case "netflow_v9", "ipfix":
		return d.decodeNetFlowIPFIX(payload, exporterIP, sourceType)
	case "sflow":
		return d.decodeSFlow(payload, exporterIP)
	default:
		return nil, fmt.Errorf("unsupported decoder type: %s", sourceType)
	}
}

func (d *Decoder) decodeNetFlowV5(payload []byte, exporterIP string) ([]model.FlowRecord, error) {
	buf := bytes.NewBuffer(payload)
	packet := netflowlegacy.PacketNetFlowV5{}
	if err := netflowlegacy.DecodeMessageVersion(buf, &packet); err != nil {
		return nil, err
	}
	exportTime := time.Unix(int64(packet.UnixSecs), int64(packet.UnixNSecs)).UTC()
	out := make([]model.FlowRecord, 0, len(packet.Records))
	for _, rec := range packet.Records {
		start := exportTime
		end := exportTime
		if packet.SysUptime > rec.First {
			start = exportTime.Add(-time.Duration(packet.SysUptime-rec.First) * time.Millisecond)
		}
		if packet.SysUptime > rec.Last {
			end = exportTime.Add(-time.Duration(packet.SysUptime-rec.Last) * time.Millisecond)
		}
		item := model.FlowRecord{
			TimestampStart:  start,
			TimestampEnd:    end,
			ExporterID:      fmt.Sprintf("%s:%d", exporterIP, packet.EngineId),
			ExporterIP:      exporterIP,
			ObservationID:   uint32(packet.EngineId),
			SrcIP:           asIP(rec.SrcAddr),
			DstIP:           asIP(rec.DstAddr),
			SrcPort:         rec.SrcPort,
			DstPort:         rec.DstPort,
			IPProtocol:      rec.Proto,
			L4ProtocolName:  util.ProtocolName(rec.Proto),
			Bytes:           uint64(rec.DOctets),
			Packets:         uint64(rec.DPkts),
			InputInterface:  uint32(rec.Input),
			OutputInterface: uint32(rec.Output),
			SrcASN:          uint32(rec.SrcAS),
			DstASN:          uint32(rec.DstAS),
			TCPFlags:        uint16(rec.TCPFlags),
			SamplerRate:     uint32(packet.SamplingInterval),
			SourceType:      "netflow_v5",
		}
		out = append(out, item)
	}
	return out, nil
}

func (d *Decoder) decodeNetFlowIPFIX(payload []byte, exporterIP, sourceType string) ([]model.FlowRecord, error) {
	buf := bytes.NewBuffer(payload)
	nf := netflow.NFv9Packet{}
	ipfix := netflow.IPFIXPacket{}
	if err := netflow.DecodeMessageVersion(buf, d.templates, &nf, &ipfix); err != nil {
		return nil, err
	}
	var flowSets []interface{}
	obsID := uint32(0)
	exportTime := time.Now().UTC()
	if sourceType == "netflow_v9" {
		flowSets = nf.FlowSets
		obsID = nf.SourceId
		exportTime = time.Unix(int64(nf.UnixSeconds), 0).UTC()
	} else {
		flowSets = ipfix.FlowSets
		obsID = ipfix.ObservationDomainId
		exportTime = time.Unix(int64(ipfix.ExportTime), 0).UTC()
	}
	out := make([]model.FlowRecord, 0, 128)
	for _, fs := range flowSets {
		flowSet, ok := fs.(netflow.DataFlowSet)
		if !ok {
			continue
		}
		for _, rec := range flowSet.Records {
			item := model.FlowRecord{
				TimestampStart: exportTime,
				TimestampEnd:   exportTime,
				ExporterID:     fmt.Sprintf("%s:%d", exporterIP, obsID),
				ExporterIP:     exporterIP,
				ObservationID:  obsID,
				SourceType:     sourceType,
			}
			for _, field := range rec.Values {
				d.mapNetFlowField(&item, field.Type, field.Value)
			}
			out = append(out, item)
		}
	}
	return out, nil
}

func (d *Decoder) decodeSFlow(payload []byte, exporterIP string) ([]model.FlowRecord, error) {
	buf := bytes.NewBuffer(payload)
	packet := sflowdec.Packet{}
	if err := sflowdec.DecodeMessageVersion(buf, &packet); err != nil {
		return nil, err
	}
	obsID := packet.SubAgentId
	now := time.Now().UTC()
	out := make([]model.FlowRecord, 0, 128)
	for _, sample := range packet.Samples {
		sv := reflect.ValueOf(sample)
		if sv.Kind() == reflect.Pointer {
			sv = sv.Elem()
		}
		if sv.Kind() != reflect.Struct {
			continue
		}
		samplingRate := getUintField(sv, "SamplingRate")
		inputIf := getUintField(sv, "Input")
		outputIf := getUintField(sv, "Output")
		recordsField := sv.FieldByName("Records")
		if !recordsField.IsValid() || recordsField.Kind() != reflect.Slice {
			continue
		}
		for i := 0; i < recordsField.Len(); i++ {
			rec := recordsField.Index(i)
			if rec.Kind() == reflect.Pointer {
				rec = rec.Elem()
			}
			if rec.Kind() != reflect.Struct {
				continue
			}
			recData := rec.FieldByName("Data")
			if !recData.IsValid() {
				continue
			}
			data := recData.Interface()
			item := model.FlowRecord{
				TimestampStart:  now,
				TimestampEnd:    now,
				ExporterID:      fmt.Sprintf("%s:%d", exporterIP, obsID),
				ExporterIP:      exporterIP,
				ObservationID:   obsID,
				InputInterface:  uint32(inputIf),
				OutputInterface: uint32(outputIf),
				SamplerRate:     uint32(samplingRate),
				SourceType:      "sflow",
			}
			switch typed := data.(type) {
			case sflowdec.SampledIPv4:
				item.SrcIP = fmt.Sprint(typed.SrcIP)
				item.DstIP = fmt.Sprint(typed.DstIP)
				item.SrcPort = uint16(typed.SrcPort)
				item.DstPort = uint16(typed.DstPort)
				item.IPProtocol = uint8(typed.Protocol)
				item.L4ProtocolName = util.ProtocolName(uint8(typed.Protocol))
				item.TCPFlags = uint16(typed.TcpFlags)
				item.Bytes = uint64(typed.Length)
			case sflowdec.SampledIPv6:
				item.SrcIP = fmt.Sprint(typed.SrcIP)
				item.DstIP = fmt.Sprint(typed.DstIP)
				item.SrcPort = uint16(typed.SrcPort)
				item.DstPort = uint16(typed.DstPort)
				item.IPProtocol = uint8(typed.Protocol)
				item.L4ProtocolName = util.ProtocolName(uint8(typed.Protocol))
				item.TCPFlags = uint16(typed.TcpFlags)
				item.Bytes = uint64(typed.Length)
			default:
				continue
			}
			if item.Bytes == 0 {
				item.Bytes = 1500
			}
			if samplingRate > 0 {
				item.Packets = uint64(samplingRate)
			} else {
				item.Packets = 1
			}
			out = append(out, item)
		}
	}
	return out, nil
}

func getUintField(v reflect.Value, name string) uint64 {
	f := v.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	switch f.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Uint()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if f.Int() < 0 {
			return 0
		}
		return uint64(f.Int())
	default:
		return 0
	}
}

func decodeVersion(payload []byte) uint16 {
	if len(payload) < 2 {
		return 0
	}
	return binary.BigEndian.Uint16(payload[:2])
}
