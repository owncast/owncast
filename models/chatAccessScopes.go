package models

const (
	// ScopeCanSendChatMessages will allow sending chat messages as itself.
	ScopeCanSendChatMessages = "CAN_SEND_MESSAGES"
	// ScopeCanSendSystemMessages will allow sending chat messages as the system.
	ScopeCanSendSystemMessages = "CAN_SEND_SYSTEM_MESSAGES"
	// ScopeHasAdminAccess will allow performing administrative actions on the server.
	ScopeHasAdminAccess = "HAS_ADMIN_ACCESS"

	ModeratorScopeKey = "MODERATOR"
)
