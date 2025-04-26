package websocket

type Message struct {
	Type     string      `json:"type"`
	RoomID   string      `json:"roomId"`
	UserID   string      `json:"userId"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata,omitempty"`
}
