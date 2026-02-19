package onebotv11

import (
	"strings"

	"github.com/Hafuunano/UniTransfer/protocol"
)

// ExtractPlainText returns only the concatenated text of text segments (no CQ/media, no placeholders).
func ExtractPlainText(msg protocol.Message) string {
	var b strings.Builder
	for _, seg := range msg {
		if seg.Type == protocol.SegmentTypeText {
			if t, ok := seg.Data["text"]; ok {
				b.WriteString(segmentDataToString(t))
			}
		}
	}
	return b.String()
}

// AltMessage returns a plain-text alternative representation of the message (e.g. for logs or alt_message).
// Text segments are output as-is; other segments are replaced by human-readable placeholders like [image], [@123].
func AltMessage(msg protocol.Message) string {
	var b strings.Builder
	for _, seg := range msg {
		switch seg.Type {
		case protocol.SegmentTypeText:
			if t, ok := seg.Data["text"]; ok {
				b.WriteString(segmentDataToString(t))
			}
		case protocol.SegmentTypeImage:
			b.WriteString("[image]")
		case protocol.SegmentTypeFace:
			b.WriteString("[face]")
		case protocol.SegmentTypeRecord:
			b.WriteString("[voice]")
		case protocol.SegmentTypeVideo:
			b.WriteString("[video]")
		case protocol.SegmentTypeAt:
			qq := dataGetStr(seg.Data, "qq")
			if qq == "all" {
				b.WriteString("[@all]")
			} else {
				b.WriteString("[@")
				b.WriteString(qq)
				b.WriteString("]")
			}
		case protocol.SegmentTypeReply:
			b.WriteString("[reply]")
		case protocol.SegmentTypeShare:
			b.WriteString("[share]")
		case protocol.SegmentTypeContact:
			b.WriteString("[contact]")
		case protocol.SegmentTypeLocation:
			b.WriteString("[location]")
		case protocol.SegmentTypeMusic:
			b.WriteString("[music]")
		case protocol.SegmentTypeForward:
			b.WriteString("[forward]")
		case protocol.SegmentTypeNode:
			b.WriteString("[node]")
		case protocol.SegmentTypeXML:
			b.WriteString("[xml]")
		case protocol.SegmentTypeJSON:
			b.WriteString("[json]")
		case protocol.SegmentTypePoke:
			b.WriteString("[poke]")
		case protocol.SegmentTypeAnonymous:
			b.WriteString("[anonymous]")
		case protocol.SegmentTypeRps:
			b.WriteString("[rps]")
		case protocol.SegmentTypeDice:
			b.WriteString("[dice]")
		case protocol.SegmentTypeShake:
			b.WriteString("[shake]")
		default:
			b.WriteString("[")
			b.WriteString(seg.Type)
			b.WriteString("]")
		}
	}
	return b.String()
}

func dataGetStr(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	return segmentDataToString(v)
}
