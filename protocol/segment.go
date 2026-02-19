// Package protocol defines protocol-agnostic unified message types (Segment, Message)
// and the plugin context interface. Protocol-specific implementations live in subpackages
// (e.g. onebotv11, zerobot).
package protocol

// Segment type constants (aligned with OneBot v11; used as canonical segment kinds).
const (
	SegmentTypeText      = "text"
	SegmentTypeFace      = "face"
	SegmentTypeImage     = "image"
	SegmentTypeRecord    = "record"
	SegmentTypeVideo     = "video"
	SegmentTypeAt        = "at"
	SegmentTypeRps       = "rps"
	SegmentTypeDice      = "dice"
	SegmentTypeShake     = "shake"
	SegmentTypePoke      = "poke"
	SegmentTypeAnonymous = "anonymous"
	SegmentTypeShare     = "share"
	SegmentTypeContact   = "contact"
	SegmentTypeLocation  = "location"
	SegmentTypeMusic     = "music"
	SegmentTypeReply     = "reply"
	SegmentTypeForward   = "forward"
	SegmentTypeNode      = "node"
	SegmentTypeXML       = "xml"
	SegmentTypeJSON      = "json"
)

// Segment is a single message segment (type + data), the unified "integrated" unit.
type Segment struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

// Message is an ordered list of segments; mixed content = multiple segments in one message.
type Message []Segment
