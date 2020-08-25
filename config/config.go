package config

//easyjson:json
type Config struct {
	LinesProvider Provider `json:"linesProvider"`
	HTTPPort      string   `json:"httpPort"`
	GrpcPort      string   `json:"grpcPort"`
	DBDataSource  string   `json:"database"`
	LogLevel      string   `json:"logLevel"`
}

//easyjson:json
type Provider struct {
	URL       string   `json:"url"`
	Sports    []string `json:"sports"`
	Intervals []int    `json:"intervals"`
}
