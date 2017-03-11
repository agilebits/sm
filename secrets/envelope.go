package secrets

// Envelope defines JSON structure that wraps the encrypted content
type Envelope struct {
	Env         string `json:"env"`
	Region      string `json:"region,omitempty"`
	MasterKeyID string `json:"master,omitempty"`

	Key EncryptionKey `json:"key"`

	Data string `json:"data"`
}
