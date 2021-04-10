package model

type LoginForm struct {
	Mail     string `form:"Mail" binding:"required"`
	Password string `form:"Password" binding:"required"`
}

type RawDataForm struct {
	Data   string `form:"Data"`
	TaskID uint64 `form:"TaskID"`
	Number uint64 `form:"Number"`
}

type ProbeCmd struct {
	CMD  string      `json:"CMD"`
	Data RawDataForm `json:"Data,omitempty"`
}

type CentreMsg struct {
	Msg  string `json:"Msg"`
	Code int    `json:"Code"`
	URL  string `json:"URL,omitempty"`
}

var (
	StartMsg    = CentreMsg{Msg: "Probe Started", Code: 1}
	PausedMsg   = CentreMsg{Msg: "Probe Paused", Code: 2}
	ResumedMsg  = CentreMsg{Msg: "Probe Resumed", Code: 3}
	FinishedMsg = CentreMsg{Msg: "Probe Finished", Code: 4}
)
