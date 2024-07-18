package client

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (c *Client) registerLogger() {
	discordgo.Logger = func(msgL int, caller int, format string, a ...any) {
		var pcs [1]uintptr
		runtime.Callers(4, pcs[:])

		var lvl slog.Level
		switch msgL {
		case discordgo.LogError:
			lvl = slog.LevelError
		case discordgo.LogWarning:
			lvl = slog.LevelWarn
		case discordgo.LogDebug:
			lvl = slog.LevelDebug
		default:
			lvl = slog.LevelInfo
		}

		msg := fmt.Sprintf(format, a...)

		r := slog.NewRecord(time.Now(), lvl, msg, pcs[0])

		if !c.logger.Handler().Enabled(context.Background(), lvl) {
			return
		}

		_ = c.logger.Handler().Handle(context.Background(), r) //nolint: errcheck
	}
}
