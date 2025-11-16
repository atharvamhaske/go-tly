package models

//for realtime dashboard and clickhouse

type AnalyticsSummary struct {
    TotalClicks int64            `json:"total_clicks"`
    Countries   map[string]int64 `json:"countries"`
    Browsers    map[string]int64 `json:"browsers"`
    Devices     map[string]int64 `json:"devices"`
}