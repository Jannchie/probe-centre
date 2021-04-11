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
	StartMsg    = CentreMsg{Msg: "Start sending tasks to the probe", Code: 1}
	PausedMsg   = CentreMsg{Msg: "The tasks have stopped sending", Code: 2}
	ResumedMsg  = CentreMsg{Msg: "Recovery tasks transfer", Code: 3}
	FinishedMsg = CentreMsg{Msg: "Disconnected from probe.", Code: 4}
	EmptyMsg    = CentreMsg{Msg: "There are currently no tasks to perform.", Code: 5}
)
