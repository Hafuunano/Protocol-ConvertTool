// Package onebotv11 provides OneBot v11 protocol compatibility: builders, CQ code, JSON, and Context implementation.
package onebotv11

import (
	"maps"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
)

// Text returns a text segment.
func Text(s string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeText,
		Data: map[string]any{"text": s},
	}
}

// Image returns an image segment. file is required; optional keys: type (e.g. "flash"), url, cache, proxy, timeout.
func Image(file string, opts map[string]any) protocol.Segment {
	data := map[string]any{"file": file}
	maps.Copy(data, opts)
	return protocol.Segment{
		Type: protocol.SegmentTypeImage,
		Data: data,
	}
}

// At returns an at segment for the given qq (user id).
func At(qq string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeAt,
		Data: map[string]any{"qq": qq},
	}
}

// AtAll returns an at-all segment.
func AtAll() protocol.Segment {
	return At("all")
}

// Face returns a QQ face segment.
func Face(id string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeFace,
		Data: map[string]any{"id": id},
	}
}

// Reply returns a reply segment referencing the given message id.
func Reply(id string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeReply,
		Data: map[string]any{"id": id},
	}
}

// Record returns a voice segment.
func Record(file string, opts map[string]any) protocol.Segment {
	data := map[string]any{"file": file}
	maps.Copy(data, opts)
	return protocol.Segment{
		Type: protocol.SegmentTypeRecord,
		Data: data,
	}
}

// Video returns a short-video segment.
func Video(file string, opts map[string]any) protocol.Segment {
	data := map[string]any{"file": file}
	maps.Copy(data, opts)
	return protocol.Segment{
		Type: protocol.SegmentTypeVideo,
		Data: data,
	}
}

// Poke returns a poke segment.
func Poke(typ, id string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypePoke,
		Data: map[string]any{"type": typ, "id": id},
	}
}

// Share returns a link share segment.
func Share(url, title string, opts map[string]any) protocol.Segment {
	data := map[string]any{"url": url, "title": title}
	maps.Copy(data, opts)
	return protocol.Segment{
		Type: protocol.SegmentTypeShare,
		Data: data,
	}
}

// Contact returns a contact (friend/group) recommendation segment. contactType is "qq" or "group".
func Contact(contactType, id string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeContact,
		Data: map[string]any{"type": contactType, "id": id},
	}
}

// Location returns a location segment.
func Location(lat, lon float64, opts map[string]any) protocol.Segment {
	data := map[string]any{"lat": lat, "lon": lon}
	maps.Copy(data, opts)
	return protocol.Segment{
		Type: protocol.SegmentTypeLocation,
		Data: data,
	}
}

// Music returns a music share segment. musicType: qq, 163, xm, custom. For custom, use opts for url, audio, title, etc.
func Music(musicType, id string, opts map[string]any) protocol.Segment {
	data := map[string]any{"type": musicType, "id": id}
	maps.Copy(data, opts)
	return protocol.Segment{
		Type: protocol.SegmentTypeMusic,
		Data: data,
	}
}

// Anonymous returns an anonymous send segment. ignore: "0" or "1".
func Anonymous(ignore string) protocol.Segment {
	data := make(map[string]any)
	if ignore != "" {
		data["ignore"] = ignore
	}
	return protocol.Segment{
		Type: protocol.SegmentTypeAnonymous,
		Data: data,
	}
}

// Forward returns a forward (merged forward) segment by id.
func Forward(id string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeForward,
		Data: map[string]any{"id": id},
	}
}

// Node returns a forward node (by message id, or custom with user_id, nickname, content).
func Node(opts map[string]any) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeNode,
		Data: opts,
	}
}

// Rps returns a rock-paper-scissors magic expression segment.
func Rps() protocol.Segment {
	return protocol.Segment{Type: protocol.SegmentTypeRps, Data: map[string]any{}}
}

// Dice returns a dice magic expression segment.
func Dice() protocol.Segment {
	return protocol.Segment{Type: protocol.SegmentTypeDice, Data: map[string]any{}}
}

// Shake returns a window shake segment.
func Shake() protocol.Segment {
	return protocol.Segment{Type: protocol.SegmentTypeShake, Data: map[string]any{}}
}

// XML returns an XML message segment.
func XML(data string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeXML,
		Data: map[string]any{"data": data},
	}
}

// JSON returns a JSON message segment.
func JSON(data string) protocol.Segment {
	return protocol.Segment{
		Type: protocol.SegmentTypeJSON,
		Data: map[string]any{"data": data},
	}
}
