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
						Cmd: h.start,
					},
				},
				{
					Path: SubcommandStop,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.stop,
						Autocomplete: h.autocompletePollList,
					},
				},
				{
					Path: SubcommandStatus,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.status,
						Autocomplete: h.autocompletePollList,
					},
				},
				{
					Path: SubcommandAddOption,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.addOption,
						Autocomplete: h.autocompletePollList,
					},
				},
				{
					Path: SubcommandRemoveOption,
					Type: disroute.TypeSubcommand,
					Handlers: disroute.Handlers{
						Cmd:          h.removeOption,
						Autocomplete: h.removeOptionAutocomplete,
					},
				},
			},
		},
	}
}
