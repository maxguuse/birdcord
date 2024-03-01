package poll

import (
	"github.com/maxguuse/disroute"
)

func (h *Handler) GetRoutes() []*disroute.Cmd {
	return []*disroute.Cmd{
		{
			Path: CommandPoll,
			Options: []*disroute.CmdOption{
				{
					Path: SubcommandStart,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd: h.startPoll,
					},
				},
				{
					Path: SubcommandStop,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.stopPoll,
						Autocomplete: h.autocompletePollList,
					},
				},
				{
					Path: SubcommandStatus,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.statusPoll,
						Autocomplete: h.autocompletePollList,
					},
				},
				{
					Path: SubcommandAddOption,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.addPollOption,
						Autocomplete: h.autocompletePollList,
					},
				},
				{
					Path: SubcommandRemoveOption,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.removePollOption,
						Autocomplete: h.removeOptionAutocomplete,
					},
				},
			},
		},
	}
}
