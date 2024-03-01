package liverole

import (
	"github.com/maxguuse/disroute"
)

func (h *Handler) GetRoutes() *disroute.Cmd {
	return &disroute.Cmd{
		Path: "liverole",
		Options: []*disroute.CmdOption{
			{
				Path: SubcommandAdd,
				Type: disroute.TypeSubcommand,
				Handlers: disroute.Handlers{
					Cmd: h.addLiveRole,
				},
			},
			{
				Path: SubcommandRemove,
				Type: disroute.TypeSubcommand,
				Handlers: disroute.Handlers{
					Cmd: h.removeLiveRole,
				},
			},
			{
				Path: SubcommandList,
				Type: disroute.TypeSubcommand,
				Handlers: disroute.Handlers{
					Cmd: h.listLiveRoles,
				},
			},
			{
				Path: SubcommandClear,
				Type: disroute.TypeSubcommand,
				Handlers: disroute.Handlers{
					Cmd: h.clearLiveRoles,
				},
			},
		},
	}
}
