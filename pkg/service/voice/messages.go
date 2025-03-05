package voice

// Messages are sent to the client

type MessageInterface interface {
	GetType() string
}

type MessageBase struct {
	Type string `json:"type"`
}

func (m MessageBase) GetType() string {
	return m.Type
}

type SpeakingStarted struct {
	IsSpeaking bool `json:"is_speaking"`
	MessageBase
}

type Connected struct {
	MessageBase
}

type NewMessage struct {
	MessageBase
}

type Failure struct {
	Reason string `json:"reason"`
	MessageBase
}

type ResponseAudio struct {
	Delta  string `json:"delta"`
	ItemID string `json:"item_id"`
	MessageBase
}
