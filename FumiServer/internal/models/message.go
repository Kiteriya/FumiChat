package models

type Message struct {
	Type      string
	SenderID  string
	TargetID  string
	Payload   string
	Timestamp int64
}
