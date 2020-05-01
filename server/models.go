package server

// Clip represents a message in a register
type Clip struct {
	Reg     string `json:"reg,omitempty"`
	Content string `json:"content"`
}
