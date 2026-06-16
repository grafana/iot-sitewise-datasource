package models

type L4eAnomalyDiagnostics struct {
	Name  string  `json:"name,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type L4eAnomalyResult struct {
	AnomalyScore     float64                 `json:"anomaly_score,omitempty"`
	PredictionReason string                  `json:"prediction_reason,omitempty"`
	Diagnostics      []L4eAnomalyDiagnostics `json:"diagnostics,omitempty"`
}
