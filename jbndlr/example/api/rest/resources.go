package rest

// InfoResource : Publicly accessible service information struct.
type InfoResource struct {
	Self    string `json:"_self"`
	Service string `json:"service"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}
