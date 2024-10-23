package checkapi

// Info represents information about the service.
type Info struct {
	Status     string `json:"status,omitempty"`
	Host       string `json:"host,omitempty"`
	GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
}
