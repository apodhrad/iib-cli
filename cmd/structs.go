package cmd

type Service struct {
	Name    string
	Methods []string
}

type Channel struct {
	Name    string `json:"name"`
	CsvName string `json:"csvName"`
}

type Package struct {
	Name               string    `json:"name"`
	Channels           []Channel `json:"channels"`
	DefaultChannelName string    `json:"defaultChannelName"`
}

type Bundle struct {
	CsvName     string `json:"csvName"`
	PackageName string `json:"packageName"`
	ChannelName string `json:"channelname"`
	BundlePath  string `json:"bundlePath"`
	Version     string `json:"version"`
	Replaces    string `json:"replaces,omitempty"`
}
