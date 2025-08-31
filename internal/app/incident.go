package app

// Incident represents a detected file system anomaly with specific metadata.
type Incident struct {
	Severity  string
	EventType string
	FilePath  string
	Message   string
}
