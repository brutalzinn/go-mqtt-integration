package command

import "encoding/json"

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Audio struct {
	Src  string `json:"src"`
	Name string `json:"name"`
}

type Notify struct {
	Title   string `json:"title"`
	Message string `json:"Message"`
}

type Command struct {
	OS struct {
		Windows struct {
			Name string `json:"name"`
			Args string `json:"args"`
		} `json:"windows"`
		Mac struct {
			Name string `json:"name"`
			Args string `json:"args"`
		} `json:"mac"`
		Linux struct {
			Name string `json:"name"`
			Args string `json:"args"`
		} `json:"linux"`
	}
}

type Youtube struct {
	Src       string `json:"src"`
	Download  bool   `json:"download"`
	OnlyAudio bool   `json:"onlyAudio"`
}
