package main

// Entry represents a host to route to either another host, or a
// file folder on disk.
type Entry struct {
	Source     string `json:"host"`
	DestHost   string `json:"dest,omitempty"`
	DestFolder string `json:"dest_folder,omitempty"`
}

func (e *Entry) dest() string {
	if e.DestHost != "" {
		return e.DestHost
	}
	return e.DestFolder
}
