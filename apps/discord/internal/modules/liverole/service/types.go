package service

type AddLiveRoleRequest struct {
	GuildID string
	RoleID  string
}

type RemoveLiveRoleRequest struct {
	GuildID string
	RoleID  string
}
