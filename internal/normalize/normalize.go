package normalize

import (
	"time"

	"github.com/google/uuid"

	"github.com/flowscope/flowscope/internal/model"
	"github.com/flowscope/flowscope/internal/util"
)

func Apply(record *model.FlowRecord) {
	if record.FlowID == "" {
		record.FlowID = uuid.NewString()
	}
	if record.TimestampStart.IsZero() {
		record.TimestampStart = time.Now().UTC()
	}
	if record.TimestampEnd.IsZero() || record.TimestampEnd.Before(record.TimestampStart) {
		record.TimestampEnd = record.TimestampStart
	}
	if record.L4ProtocolName == "" {
		record.L4ProtocolName = util.ProtocolName(record.IPProtocol)
	}
	if record.SrcSubnet == "" {
		record.SrcSubnet = util.SubnetForIP(record.SrcIP)
	}
	if record.DstSubnet == "" {
		record.DstSubnet = util.SubnetForIP(record.DstIP)
	}
	record.SrcIsPrivate = util.IsPrivateIP(record.SrcIP)
	record.DstIsPrivate = util.IsPrivateIP(record.DstIP)
	if record.MinuteBucket.IsZero() {
		record.MinuteBucket = util.MinuteBucket(record.TimestampStart)
	}
	if record.HourBucket.IsZero() {
		record.HourBucket = util.HourBucket(record.TimestampStart)
	}
	if record.FlowKeyHash == "" {
		record.FlowKeyHash = util.FlowHash(
			record.ExporterID,
			record.SrcIP,
			record.DstIP,
			record.SrcPort,
			record.DstPort,
			record.IPProtocol,
			record.MinuteBucket.Unix(),
		)
	}
	if record.FlowDirection == "" {
		switch {
		case record.SrcIsPrivate && !record.DstIsPrivate:
			record.FlowDirection = "egress"
		case !record.SrcIsPrivate && record.DstIsPrivate:
			record.FlowDirection = "ingress"
		default:
			record.FlowDirection = "lateral"
		}
	}
	if record.SourceType == "" {
		record.SourceType = "unknown"
	}
}
