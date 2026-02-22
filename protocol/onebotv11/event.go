// Package onebotv11: event types aligned with OneBot v11 message event spec.
// See https://github.com/botuniverse/onebot-11/blob/master/event/message.md

package onebotv11

import (
	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
)

// MessageEvent is the common message event payload (private and group).
// Fields match OneBot v11 event/message.md.
type MessageEvent struct {
	Time        int64            `json:"time"`
	SelfID      int64            `json:"self_id"`
	PostType    string           `json:"post_type"`    // "message"
	MessageType string           `json:"message_type"` // "private" | "group"
	SubType     string           `json:"sub_type"`
	MessageID   int32            `json:"message_id"`
	GroupID     int64            `json:"group_id"` // 0 for private
	UserID      int64            `json:"user_id"`
	Message     protocol.Message `json:"-"` // parsed; fill from raw or array
	RawMessage  string           `json:"raw_message"`
	Font        int32            `json:"font,omitempty"`
	Sender      *SenderInfo      `json:"sender,omitempty"`
	Anonymous   *AnonymousInfo   `json:"anonymous,omitempty"`
}

// SenderInfo matches OneBot v11 sender object (private: user_id, nickname, sex, age; group: + card, role, title, etc.).
type SenderInfo struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname,omitempty"`
	Card     string `json:"card,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Age      int32  `json:"age,omitempty"`
	Area     string `json:"area,omitempty"`
	Level    string `json:"level,omitempty"`
	Role     string `json:"role,omitempty"` // "owner" | "admin" | "member"
	Title    string `json:"title,omitempty"`
}

// AnonymousInfo is present for anonymous group messages.
type AnonymousInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// DisplayName returns sender display name (card preferred in group, else nickname).
func (s *SenderInfo) DisplayName() string {
	if s == nil {
		return ""
	}
	if s.Card != "" {
		return s.Card
	}
	return s.Nickname
}
