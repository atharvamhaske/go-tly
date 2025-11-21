package models

type DomainHealth struct {
	DNSResolved   bool   `bson:"dns_resolved" json:"dns_resolved"`
	SSLValid      bool   `bson:"ssl_valid" json:"ssl_valid"`
	StatusCode    int    `bson:"status_code" json:"status_code"`
	LatencyMS     int64  `bson:"latency_ms" json:"latency_ms"`
	Message       string `bson:"message,omitempty" json:"message"`
	CloudflareSec string `bson:"cf_security" json:"cf_security"`
}
