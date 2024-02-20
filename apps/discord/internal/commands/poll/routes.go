package poll

import "github.com/maxguuse/disroute"

func (h *Handler) GetRoutes() *disroute.Cmd {
	return &disroute.Cmd{
		Path: CommandPoll,
		Options: []*disroute.CmdOption{
			{
				Path:    SubcommandStart,
				Type:    disroute.TypeSubcommand,
				Handler: h.startPoll,
			},
			{
				Path:    SubcommandStop,
				Type:    disroute.TypeSubcommand,
				Handler: h.stopPoll,
			},
			{
				Path:    SubcommandStatus,
				Type:    disroute.TypeSubcommand,
				Handler: h.statusPoll,
			},
			{
				Path:    SubcommandAddOption,
				Type:    disroute.TypeSubcommand,
				Handler: h.addPollOption,
			},
			{
				Path:    SubcommandRemoveOption,
				Type:    disroute.TypeSubcommand,
				Handler: h.removePollOption,
			},
		},
	}
}
