package models

type JournalRef struct {
	DocRef string `json:"docref"`
}

type JournalCommand struct {
	Command string      `json:"command"`
	Payload interface{} `json:"payload,omitempty"`
}

type DocRefPool struct {
	DocRef   string `json:"docref"`
	Username string `json:"username"`
}

type DocRefEvent struct {
	DocRef   string `json:"docref"`
	Username string `json:"username"`
	Status   string `json:"status"`
}
