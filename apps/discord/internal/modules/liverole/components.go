package liverole

import "github.com/maxguuse/disroute"

func (h *Handler) GetComponents() []*disroute.Component {
	return make([]*disroute.Component, 0)
}
