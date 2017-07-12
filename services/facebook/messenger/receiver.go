package messenger

import (
	"encoding/json"
	"io"
	"time"
)

type ReceivedMessage struct {
	PageID      string
	SenderID    string
	RecipientID string
	UpdateTime  time.Time
	Timestamp   time.Time
	Content     interface{}
}

func Parse(r io.Reader) ([]ReceivedMessage, error) {
	var messages []ReceivedMessage
	var err error
	var p callbackPayload
	if err = json.NewDecoder(r).Decode(&p); err != nil {
		return nil, err
	}
	for _, ent := range p.Entry {
		updatedTime := time.Unix(ent.Time/1000, ent.Time%1000*1000000)
		for _, raw := range ent.RawMessaging {
			var m ReceivedMessage
			var v map[string]json.RawMessage
			if err = json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			m.UpdateTime = updatedTime
			m.SenderID = parseUserID(v, "sender")
			m.RecipientID = parseUserID(v, "recipient")
			m.Timestamp = parseTimetamp(v, "timestamp")
			var content json.RawMessage
			var ok bool
			if content, ok = v["account_linking"]; ok {
				m.Content = &AccountLinking{}
			} else if content, ok = v["message"]; ok {
				m.Content = &ReceivedMessageContent{}
			}
			if m.Content != nil {
				if err = json.Unmarshal(content, m.Content); err != nil {
					m.Content = nil
				}
			}
			messages = append(messages, m)
		}
	}
	return messages, nil
}

func parseUserID(v map[string]json.RawMessage, attr string) string {
	if vv, ok := v[attr]; ok {
		var u userID
		json.Unmarshal(vv, &u)
		return u.ID
	}
	return ""
}

func parseTimetamp(v map[string]json.RawMessage, attr string) time.Time {
	if vv, ok := v[attr]; ok {
		var i int64
		json.Unmarshal(vv, &i)
		if i != 0 {
			return time.Unix(i/1000, i%1000*1000000)
		}
	}
	return time.Time{}
}

type callbackPayload struct {
	Object string   `json:"object"`
	Entry  []*entry `json:"entry"`
}

type entry struct {
	ID           string            `json:"id"`   // page id
	Time         int64             `json:"time"` // timestamp
	RawMessaging []json.RawMessage `json:"messaging"`
}

// AccountStatus is a enum type for account linking status.
type AccountStatus string

// AccountStatus constants
const (
	AccountStatusLinked   AccountStatus = "linked"
	AccountStatusUnlinked AccountStatus = "unlinked"
)

type AccountLinking struct {
	Status            AccountStatus `json:"status"`
	AuthorizationCode string        `json:"authorization_code"`
}

type userID struct {
	ID string `json:"id"`
}

type ReceivedMessageContent struct {
	MID         string               `json:"mid"`
	Seq         int                  `json:"seq"`
	Text        string               `json:"text"`
	QuickReply  *ReceivedQuickReply  `json:"quick_reply"`
	Attachments []ReceivedAttachment `json:"attachments"`
}

type ReceivedQuickReply struct {
	Payload string `json:"payload"`
}

type ReceivedAttachment struct {
	Type    string                     `json:"type"`
	Payload *ReceivedAttachmentPayload `json:"payload"`
}

type ReceivedAttachmentPayload struct {
	URL         string       `json:"url"`
	Coordinates *Coordinates `json:"coordinates"`
}
