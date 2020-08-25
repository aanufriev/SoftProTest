package config

// Config is a struct into which
// the config.json will be parsed
//easyjson:json
type Config struct {
	LinesProvider provider `json:"linesProvider"`
	HTTPPort      string   `json:"httpPort"`
	GrpcPort      string   `json:"grpcPort"`
	DBDataSource  string   `json:"database"`
	LogLevel      string   `json:"logLevel"`
}

//easyjson:json
type provider struct {
	URL       string   `json:"url"`
	Sports    []string `json:"sports"`
	Intervals []int    `json:"intervals"`
}
