package models

type Notifications []Notification

type Notification struct {
	ID           int         `json:"pk"`
	Message      string      `json:"message"`
	Details      string      `json:"details"`
	Category     string      `json:"category"`
	Topic        string      `json:"topic"`
	Level        string      `json:"level"`
	Announcement interface{} `json:"announcement"`
	BuildLog     BuildLog    `json:"build_log"`
	ExtraData    ExtraData   `json:"extra_data"`
	Read         bool        `json:"read"`
	Seen         bool        `json:"seen"`
	Issued       string      `json:"issued"`
}

type BuildLog struct {
	ID            int    `json:"pk"`
	Cluster       int    `json:"cluster"`
	EventCategory string `json:"event_category"`
	EventType     string `json:"event_type"`
	EventState    string `json:"event_state"`
	Message       string `json:"message"`
	Reference     string `json:"reference"`
	Created       string `json:"created"`
}

type ExtraData struct {
	Org struct {
		ID   int    `json:"pk"`
		Name string `json:"name"`
	} `json:"org"`
	Cluster struct {
		ID   int    `json:"pk"`
		Name string `json:"name"`
	} `json:"cluster"`
}
