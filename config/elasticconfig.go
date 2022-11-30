package config

type TimesResponse struct {
	ServiceName   string `json:"ServiceName"`
	ReleaseNumber string `json:"ReleaseNumber"`
	BuildNumber   string `json:"BuildNumber"`
	Percentile99  int64  `json:"percentile99"`
	Label         string `json:"Label"`
	Average       int64  `json:"Average"`
	Scenario      string `json:"Scenario"`
	Iteration     int64  `json:"Iteration"`
}
