package ui

import (
	"testing"
	"time"

	"github.com/anchore/bubbly/bubbles/taskprogress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
	"github.com/wagoodman/go-partybus"
	"github.com/wagoodman/go-progress"

	"github.com/xeol-io/xeol/xeol/event"
)

func TestHandler_handleUpdateEolDatabase(t *testing.T) {

	tests := []struct {
		name       string
		eventFn    func(*testing.T) partybus.Event
		iterations int
	}{
		{
			name: "downloading DB",
			eventFn: func(t *testing.T) partybus.Event {
				prog := &progress.Manual{}
				prog.SetTotal(100)
				prog.Set(50)

				mon := struct {
					progress.Progressable
					progress.Stager
				}{
					Progressable: prog,
					Stager: &progress.Stage{
						Current: "current",
					},
				}

				return partybus.Event{
					Type:  event.UpdateEolDatabase,
					Value: mon,
				}
			},
		},
		{
			name: "DB download complete",
			eventFn: func(t *testing.T) partybus.Event {
				prog := &progress.Manual{}
				prog.SetTotal(100)
				prog.Set(100)
				prog.SetCompleted()

				mon := struct {
					progress.Progressable
					progress.Stager
				}{
					Progressable: prog,
					Stager: &progress.Stage{
						Current: "current",
					},
				}

				return partybus.Event{
					Type:  event.UpdateEolDatabase,
					Value: mon,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.eventFn(t)
			handler := New(DefaultHandlerConfig())
			handler.WindowSize = tea.WindowSizeMsg{
				Width:  100,
				Height: 80,
			}

			models := handler.Handle(e)
			require.Len(t, models, 1)
			model := models[0]

			tsk, ok := model.(taskprogress.Model)
			require.True(t, ok)

			got := runModel(t, tsk, tt.iterations, taskprogress.TickMsg{
				Time:     time.Now(),
				Sequence: tsk.Sequence(),
				ID:       tsk.ID(),
			})
			t.Log(got)
			snaps.MatchSnapshot(t, got)
		})
	}
}