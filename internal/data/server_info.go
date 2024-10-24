package data

type ServerInfo struct {
	Status     string               `json:"status,omitempty"`
	SystemInfo ServerInfoSystemInfo `json:"system_info,omitempty"`
}
