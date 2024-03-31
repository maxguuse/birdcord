package liverole

import (
	"strings"

	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
	"github.com/maxguuse/disroute"
)

func (h *Handler) addLiveRole(ctx *disroute.Ctx) disroute.Response {
	role := ctx.Options["role"].RoleValue(ctx.Session(), ctx.Interaction().GuildID)

	err := h.Service.Add(ctx.Context(), &service.AddLiveRoleRequest{
		GuildID: ctx.Interaction().GuildID,
		RoleID:  role.ID,
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	return disroute.Response{
		Message: "Live-роль успешно добавлена.",
	}
}

func (h *Handler) clearLiveRoles(ctx *disroute.Ctx) disroute.Response {
	err := h.Service.Clear(ctx.Context(), ctx.Interaction().GuildID)
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	return disroute.Response{
		Message: "Live-роли успешно удалены.",
	}
}

func (h *Handler) listLiveRoles(ctx *disroute.Ctx) disroute.Response {
	rolesList, err := h.Service.List(ctx.Context(), ctx.Interaction().GuildID)
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	return disroute.Response{
		Message: "Список live-ролей: \n" + strings.Join(rolesList, "\n"),
	}
}

func (h *Handler) removeLiveRole(ctx *disroute.Ctx) disroute.Response {
	role := ctx.Options["role"].RoleValue(ctx.Session(), ctx.Interaction().GuildID)

	err := h.Service.Remove(ctx.Context(), &service.RemoveLiveRoleRequest{
		GuildID: ctx.Interaction().GuildID,
		RoleID:  role.ID,
	})
	if err != nil {
		return disroute.Response{
			Err: err,
		}
	}

	return disroute.Response{
		Message: "Live-роль успешно удалена.",
	}
}
