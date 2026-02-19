package onebotv11

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Hafuunano/UniTransfer/protocol"
)

// CQ escape sequences (OneBot v11).
const (
	cqEscAmp   = "&amp;"
	cqEscLB    = "&#91;"
	cqEscRB    = "&#93;"
	cqEscComma = "&#44;"
)

// cqUnescape reverses CQ code escaping in value: &amp; -> &, &#91; -> [, &#93; -> ], &#44; -> ,
func cqUnescape(s string) string {
	s = strings.ReplaceAll(s, cqEscAmp, "&")
	s = strings.ReplaceAll(s, cqEscLB, "[")
	s = strings.ReplaceAll(s, cqEscRB, "]")
	s = strings.ReplaceAll(s, cqEscComma, ",")
	return s
}

// cqEscape escapes & [ ] , for use inside CQ param values.
func cqEscape(s string) string {
	s = strings.ReplaceAll(s, "&", cqEscAmp)
	s = strings.ReplaceAll(s, "[", cqEscLB)
	s = strings.ReplaceAll(s, "]", cqEscRB)
	s = strings.ReplaceAll(s, ",", cqEscComma)
	return s
}

// cqCodeRegex matches [CQ:type,key=value,...]. Captures type and the rest (params string).
var cqCodeRegex = regexp.MustCompile(`\[CQ:([^,\]]+)(?:,(.*?))?\]`)

// ParseCQ parses a CQ code string into a slice of segments. Plain text between CQ codes becomes text segments.
func ParseCQ(s string) protocol.Message {
	var out protocol.Message
	lastEnd := 0
	for _, loc := range cqCodeRegex.FindAllStringSubmatchIndex(s, -1) {
		fullStart, fullEnd := loc[0], loc[1]
		typeStart, typeEnd := loc[2], loc[3]
		paramsStart, paramsEnd := loc[4], loc[5]
		if fullStart > lastEnd {
			plain := s[lastEnd:fullStart]
			if plain != "" {
				out = append(out, Text(cqUnescape(plain)))
			}
		}
		segType := s[typeStart:typeEnd]
		data := map[string]any{}
		if paramsStart >= 0 && paramsEnd >= 0 && paramsStart < paramsEnd {
			paramsStr := cqUnescape(s[paramsStart:paramsEnd])
			for _, part := range strings.Split(paramsStr, ",") {
				eq := strings.Index(part, "=")
				if eq <= 0 {
					continue
				}
				key := strings.TrimSpace(part[:eq])
				val := strings.TrimSpace(part[eq+1:])
				data[key] = val
			}
		}
		out = append(out, protocol.Segment{Type: segType, Data: data})
		lastEnd = fullEnd
	}
	if lastEnd < len(s) {
		plain := s[lastEnd:]
		if plain != "" {
			out = append(out, Text(cqUnescape(plain)))
		}
	}
	return out
}

// ToCQ converts a segment to its CQ code string (no leading/trailing text).
func ToCQ(seg protocol.Segment) string {
	if seg.Type == protocol.SegmentTypeText {
		if t, ok := seg.Data["text"]; ok {
			return cqEscape(segmentDataToString(t))
		}
		return ""
	}
	var parts []string
	for k, v := range seg.Data {
		parts = append(parts, k+"="+cqEscape(segmentDataToString(v)))
	}
	if len(parts) == 0 {
		return "[CQ:" + seg.Type + "]"
	}
	return "[CQ:" + seg.Type + "," + strings.Join(parts, ",") + "]"
}

// segmentDataToString converts segment data value to string for CQ/JSON.
func segmentDataToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	default:
		return ""
	}
}

// MessageToCQ converts a full message to a single CQ code string (text segments and CQ codes concatenated).
func MessageToCQ(msg protocol.Message) string {
	var b strings.Builder
	for _, seg := range msg {
		if seg.Type == protocol.SegmentTypeText {
			if t, ok := seg.Data["text"]; ok {
				b.WriteString(cqEscape(segmentDataToString(t)))
			}
		} else {
			b.WriteString(ToCQ(seg))
		}
	}
	return b.String()
}
