package domain

type Liverole struct {
	ID             int `sql:"primary_key" alias:"liveroles.id"`
	DiscordRoleID  int `                  alias:"liveroles.discord_role_id"`
	DiscordGuildID int `                  alias:"liveroles.discord_guild_id"`
}
