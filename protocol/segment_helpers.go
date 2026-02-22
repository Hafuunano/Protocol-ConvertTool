// Package protocol (segment_helpers.go): protocol-agnostic segment constructors.
// Types and constants stay in segment.go; helpers for building segments (e.g. PokeUser) live here.
package protocol

// PokeUser returns a poke segment targeting the given user ID (protocol-agnostic).
// Data includes "id" and "qq" for compatibility with OneBot/ZeroBot (qq is used by ZeroBot for poke target).
func PokeUser(targetUserID string) Segment {
	return Segment{
		Type: SegmentTypePoke,
		Data: map[string]any{"id": targetUserID, "qq": targetUserID},
	}
}
