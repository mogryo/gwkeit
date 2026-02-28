package widgets

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/slicelib"
	"github.com/rivo/tview"
)

type LogType int

const (
	ErrorMessage LogType = iota
	SuccessMessage
	InfoMessage
	TimestampMessage
)

type LogsWidget struct {
	View *tview.TextView
	list []string
}

func NewLogsWidget(tools *apptools.Tools) *LogsWidget {
	lw := &LogsWidget{}
	lw.View = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			tools.RefreshScreen()
		})
	lw.View.SetBackgroundColor(tcell.ColorDefault)

	lw.list = make([]string, 0)
	return lw
}

func (lw *LogsWidget) addTimestampLog() {
	timestamp := "[grey]" + time.Now().Format("2006-01-02 15:04:05") + "[-]"
	lw.addLog(timestamp, TimestampMessage)
}

func (lw *LogsWidget) AddErrorLogs(logs []string) {
	if len(logs) == 0 {
		return
	}

	lw.addTimestampLog()
	for _, log := range logs {
		lw.addLog(log, ErrorMessage)
	}
	lw.View.SetText(strings.Join(slicelib.TakeLast(lw.list, 12), "\n"))
}

func (lw *LogsWidget) AddSuccessLogs(logs []string) {
	if len(logs) == 0 {
		return
	}

	lw.addTimestampLog()
	for _, log := range logs {
		lw.addLog(log, SuccessMessage)
	}
	lw.View.SetText(strings.Join(slicelib.TakeLast(lw.list, 12), "\n"))
}

func (lw *LogsWidget) AddInfoLogs(logs []string) {
	if len(logs) == 0 {
		return
	}
	lw.addTimestampLog()
	for _, log := range logs {
		lw.addLog(log, InfoMessage)
	}
	lw.View.SetText(strings.Join(slicelib.TakeLast(lw.list, 12), "\n"))
}

func (lw *LogsWidget) addLog(message string, logType LogType) {
	messageWithType := message
	switch logType {
	case ErrorMessage:
		messageWithType = "[#ff4689]  " + messageWithType + "[-]"
	case SuccessMessage:
		messageWithType = "[#a6e22e]  " + messageWithType + "[-]"
	case InfoMessage:
		messageWithType = "[#0087ff]  " + messageWithType + "[-]"
	case TimestampMessage:
		messageWithType = "[grey]" + messageWithType + "[-]"
	}
	lw.list = append(lw.list, messageWithType)
}
