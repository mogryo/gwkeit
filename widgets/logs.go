package widgets

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/slicelib"
	"github.com/rivo/tview"
)

type LogType int

const (
	ERROR_MESSAGE LogType = iota
	SUCCESS_MESSAGE
	TIMESTAMP_MESSAGE
)

type LogsWidget struct {
	View *tview.TextView
	list []string
}

func NewLogsWidget(app *tview.Application) *LogsWidget {
	lw := &LogsWidget{}
	lw.View = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	lw.View.SetBackgroundColor(tcell.ColorDefault)

	lw.list = make([]string, 0)
	return lw
}

func (lw *LogsWidget) addTimestampLog() {
	timestamp := "[grey]" + time.Now().Format("2006-01-02 15:04:05") + "[-]"
	lw.addLog(timestamp, TIMESTAMP_MESSAGE)
}

func (lw *LogsWidget) AddErrorLogs(logs []string) {
	if len(logs) == 0 {
		return
	}

	lw.addTimestampLog()
	for _, log := range logs {
		lw.addLog(log, ERROR_MESSAGE)
	}
	lw.View.SetText(strings.Join(slicelib.TakeLast(lw.list, 12), "\n"))
}

func (lw *LogsWidget) AddSuccessLogs(logs []string) {
	if len(logs) == 0 {
		return
	}

	lw.addTimestampLog()
	for _, log := range logs {
		lw.addLog(log, SUCCESS_MESSAGE)
	}
	lw.View.SetText(strings.Join(slicelib.TakeLast(lw.list, 12), "\n"))
}

func (lw *LogsWidget) addLog(message string, logType LogType) {
	messageWithType := message
	switch logType {
	case ERROR_MESSAGE:
		messageWithType = "[#ff4689]  " + messageWithType + "[-]"
	case SUCCESS_MESSAGE:
		messageWithType = "[#a6e22e]  " + messageWithType + "[-]"
	case TIMESTAMP_MESSAGE:
		messageWithType = "[grey]" + messageWithType + "[-]"
	}
	lw.list = append(lw.list, messageWithType)
}
