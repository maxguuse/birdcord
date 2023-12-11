package client

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"runtime"
	"time"
)

func (c *Client) registerLogger() {
	discordgo.Logger = func(msgL int, caller int, format string, a ...interface{}) {
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
		fullMsg := fmt.Sprintf("%s", msg)

		r := slog.NewRecord(time.Now(), lvl, fullMsg, pcs[0])

		_ = c.Log.Handler().Handle(context.Background(), r)
	}
}
