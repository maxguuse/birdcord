package liverole

import (
	"github.com/maxguuse/disroute"
)

func (h *Handler) GetRoutes() *disroute.Cmd {
	return &disroute.Cmd{
		Path: "liverole",
		Options: []*disroute.CmdOption{
			{
				Path:    SubcommandAdd,
				Type:    disroute.TypeSubcommand,
				Handler: h.addLiveRole,
			},
			{
				Path:    SubcommandRemove,
				Type:    disroute.TypeSubcommand,
				Handler: h.removeLiveRole,
			},
			{
				Path:    SubcommandList,
				Type:    disroute.TypeSubcommand,
				Handler: h.listLiveRoles,
			},
			{
				Path:    SubcommandClear,
				Type:    disroute.TypeSubcommand,
				Handler: h.clearLiveRoles,
			},
		},
	}
}
