package cmd

type Channel struct {
	Name    string `json:"name"`
	CsvName string `json:"csvName"`
}

type Package struct {
	Name               string    `json:"name"`
	Channels           []Channel `json:"channels"`
	DefaultChannelName string    `json:"defaultChannelName"`
}
