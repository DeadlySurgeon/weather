package weather

// Report is what we return after consuming from the weather report.
type Report struct {
	Condition      string  `json:"condition"`
	Temperature    string  `json:"temperature"`
	TemperatureRaw float32 `json:"temperature_raw"`
}
