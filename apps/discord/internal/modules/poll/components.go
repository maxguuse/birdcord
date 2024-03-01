package poll

import "github.com/maxguuse/disroute"

func (h *Handler) GetComponents() []*disroute.Component {
	return []*disroute.Component{
		{
			Key:     "poll-vote-btn",
			Handler: h.VoteBtnHandler,
		},
	}
}
