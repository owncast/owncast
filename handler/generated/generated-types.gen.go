// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package generated

import (
	"encoding/json"
	"time"

	"github.com/oapi-codegen/runtime"
)

const (
	BasicAuthScopes  = "BasicAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for WebhookEventType.
const (
	CHAT               WebhookEventType = "CHAT"
	CHATACTION         WebhookEventType = "CHAT_ACTION"
	NAMECHANGE         WebhookEventType = "NAME_CHANGE"
	PING               WebhookEventType = "PING"
	PONG               WebhookEventType = "PONG"
	STREAMSTARTED      WebhookEventType = "STREAM_STARTED"
	STREAMSTOPPED      WebhookEventType = "STREAM_STOPPED"
	STREAMTITLEUPDATED WebhookEventType = "STREAM_TITLE_UPDATED"
	SYSTEM             WebhookEventType = "SYSTEM"
	USERJOINED         WebhookEventType = "USER_JOINED"
	USERPARTED         WebhookEventType = "USER_PARTED"
	VISIBILITYUPDATE   WebhookEventType = "VISIBILITY-UPDATE"
)

// ActionMessage defines model for ActionMessage.
type ActionMessage struct {
	Body      *string `json:"body,omitempty"`
	Id        *string `json:"id,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// AdminConfigValue defines model for AdminConfigValue.
type AdminConfigValue struct {
	Value *AdminConfigValue_Value `json:"value,omitempty"`
}

// AdminConfigValueValue0 defines model for .
type AdminConfigValueValue0 = string

// AdminConfigValueValue1 defines model for .
type AdminConfigValueValue1 = bool

// AdminConfigValue_Value defines model for AdminConfigValue.Value.
type AdminConfigValue_Value struct {
	union json.RawMessage
}

// AdminFederationConfig defines model for AdminFederationConfig.
type AdminFederationConfig struct {
	BlockedDomains *[]string `json:"blockedDomains,omitempty"`
	Enabled        *bool     `json:"enabled,omitempty"`
	GoLiveMessage  *string   `json:"goLiveMessage,omitempty"`
	IsPrivate      *bool     `json:"isPrivate,omitempty"`
	ShowEngagement *bool     `json:"showEngagement,omitempty"`
	Username       *string   `json:"username,omitempty"`
}

// AdminLog defines model for AdminLog.
type AdminLog struct {
	Level   *string    `json:"level,omitempty"`
	Message *string    `json:"message,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
}

// AdminMessageVisibilityUpdate defines model for AdminMessageVisibilityUpdate.
type AdminMessageVisibilityUpdate struct {
	IdArray *[]string `json:"idArray,omitempty"`
	Visible *bool     `json:"visible,omitempty"`
}

// AdminNotificationsConfig defines model for AdminNotificationsConfig.
type AdminNotificationsConfig struct {
	Browser *BrowserNotificationConfiguration `json:"browser,omitempty"`
	Discord *DiscordConfiguration             `json:"discord,omitempty"`
}

// AdminServerConfig defines model for AdminServerConfig.
type AdminServerConfig struct {
	AdminPassword           *string                   `json:"adminPassword,omitempty"`
	ChatDisabled            *bool                     `json:"chatDisabled,omitempty"`
	ChatEstablishedUserMode *bool                     `json:"chatEstablishedUserMode,omitempty"`
	ChatJoinMessagesEnabled *bool                     `json:"chatJoinMessagesEnabled,omitempty"`
	DisableSearchIndexing   *bool                     `json:"disableSearchIndexing,omitempty"`
	ExternalActions         *[]ExternalAction         `json:"externalActions,omitempty"`
	Federation              *AdminFederationConfig    `json:"federation,omitempty"`
	FfmpegPath              *string                   `json:"ffmpegPath,omitempty"`
	ForbiddenUsernames      *[]string                 `json:"forbiddenUsernames,omitempty"`
	HideViewerCount         *bool                     `json:"hideViewerCount,omitempty"`
	InstanceDetails         *AdminWebConfig           `json:"instanceDetails,omitempty"`
	Notifications           *AdminNotificationsConfig `json:"notifications,omitempty"`
	RtmpServerPort          *int                      `json:"rtmpServerPort,omitempty"`
	S3                      *S3Info                   `json:"s3,omitempty"`
	SocketHostOverride      *string                   `json:"socketHostOverride,omitempty"`
	StreamKeyOverridden     *bool                     `json:"streamKeyOverridden,omitempty"`
	StreamKeys              *[]StreamKey              `json:"streamKeys,omitempty"`
	SuggestedUsernames      *[]string                 `json:"suggestedUsernames,omitempty"`
	SupportedCodecs         *[]string                 `json:"supportedCodecs,omitempty"`
	VideoCodec              *string                   `json:"videoCodec,omitempty"`
	VideoServingEndpoint    *string                   `json:"videoServingEndpoint,omitempty"`
	VideoSettings           *AdminVideoSettings       `json:"videoSettings,omitempty"`
	WebServerIP             *string                   `json:"webServerIP,omitempty"`
	WebServerPort           *int                      `json:"webServerPort,omitempty"`
	Yp                      *AdminYPInfo              `json:"yp,omitempty"`
}

// AdminStatus defines model for AdminStatus.
type AdminStatus struct {
	Broadcaster            *Broadcaster          `json:"broadcaster,omitempty"`
	CurrentBroadcast       *CurrentBroadcast     `json:"currentBroadcast,omitempty"`
	Health                 *StreamHealthOverview `json:"health,omitempty"`
	Online                 *bool                 `json:"online,omitempty"`
	OverallPeakViewerCount *int                  `json:"overallPeakViewerCount,omitempty"`
	SessionPeakViewerCount *int                  `json:"sessionPeakViewerCount,omitempty"`
	StreamTitle            *string               `json:"streamTitle,omitempty"`
	VersionNumber          *string               `json:"versionNumber,omitempty"`
	ViewerCount            *int                  `json:"viewerCount,omitempty"`
}

// AdminVideoSettings defines model for AdminVideoSettings.
type AdminVideoSettings struct {
	LatencyLevel         *int                   `json:"latencyLevel,omitempty"`
	VideoQualityVariants *[]StreamOutputVariant `json:"videoQualityVariants,omitempty"`
}

// AdminWebConfig defines model for AdminWebConfig.
type AdminWebConfig struct {
	AppearanceVariables *map[string]string `json:"appearanceVariables,omitempty"`
	CustomJavascript    *string            `json:"customJavascript,omitempty"`
	CustomStyles        *string            `json:"customStyles,omitempty"`
	ExtraPageContent    *string            `json:"extraPageContent,omitempty"`
	Logo                *string            `json:"logo,omitempty"`
	Name                *string            `json:"name,omitempty"`
	Nsfw                *bool              `json:"nsfw,omitempty"`
	OfflineMessage      *string            `json:"offlineMessage,omitempty"`
	SocialHandles       *[]SocialHandle    `json:"socialHandles,omitempty"`
	StreamTitle         *string            `json:"streamTitle,omitempty"`
	Summary             *string            `json:"summary,omitempty"`
	Tags                *[]string          `json:"tags,omitempty"`
	Version             *string            `json:"version,omitempty"`
	WelcomeMessage      *string            `json:"welcomeMessage,omitempty"`
}

// AdminYPInfo defines model for AdminYPInfo.
type AdminYPInfo struct {
	Enabled     *bool   `json:"enabled,omitempty"`
	InstanceUrl *string `json:"instanceUrl,omitempty"`
}

// AnonymousUser defines model for AnonymousUser.
type AnonymousUser struct {
	AccessToken *string `json:"accessToken,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Id          *string `json:"id,omitempty"`
}

// AuthenticationConfig defines model for AuthenticationConfig.
type AuthenticationConfig struct {
	IndieAuthEnabled *bool `json:"indieAuthEnabled,omitempty"`
}

// BaseAPIResponse Simple API response
type BaseAPIResponse struct {
	Message *string `json:"message,omitempty"`
	Success *bool   `json:"success,omitempty"`
}

// Broadcaster defines model for Broadcaster.
type Broadcaster struct {
	RemoteAddr    *string               `json:"remoteAddr,omitempty"`
	StreamDetails *InboundStreamDetails `json:"streamDetails,omitempty"`
	Time          *time.Time            `json:"time,omitempty"`
}

// BrowserConfig defines model for BrowserConfig.
type BrowserConfig struct {
	Enabled   *bool   `json:"enabled,omitempty"`
	PublicKey *string `json:"publicKey,omitempty"`
}

// BrowserNotificationConfiguration defines model for BrowserNotificationConfiguration.
type BrowserNotificationConfiguration struct {
	Enabled       *bool   `json:"enabled,omitempty"`
	GoLiveMessage *string `json:"goLiveMessage,omitempty"`
}

// ChatClient defines model for ChatClient.
type ChatClient struct {
	ConnectedAt  *time.Time  `json:"connectedAt,omitempty"`
	Geo          *GeoDetails `json:"geo,omitempty"`
	MessageCount *int        `json:"messageCount,omitempty"`
	User         *User       `json:"user,omitempty"`
	UserAgent    *string     `json:"userAgent,omitempty"`
}

// ChatClients defines model for ChatClients.
type ChatClients = []ChatClient

// ChatMessages defines model for ChatMessages.
type ChatMessages = []ChatMessages_Item

// ChatMessages_Item defines model for ChatMessages.Item.
type ChatMessages_Item struct {
	union json.RawMessage
}

// CollectedMetrics defines model for CollectedMetrics.
type CollectedMetrics struct {
	Cpu    *[]TimestampedValue `json:"cpu,omitempty"`
	Disk   *[]TimestampedValue `json:"disk,omitempty"`
	Memory *[]TimestampedValue `json:"memory,omitempty"`
}

// CurrentBroadcast defines model for CurrentBroadcast.
type CurrentBroadcast struct {
	LatencyLevel   *LatencyLevel          `json:"latencyLevel,omitempty"`
	OutputSettings *[]StreamOutputVariant `json:"outputSettings,omitempty"`
}

// DiscordConfiguration defines model for DiscordConfiguration.
type DiscordConfiguration struct {
	Enabled       *bool   `json:"enabled,omitempty"`
	GoLiveMessage *string `json:"goLiveMessage,omitempty"`
	Webhook       *string `json:"webhook,omitempty"`
}

// Emoji Name and url for an emoji
type Emoji struct {
	// Name The name of the emoji
	Name *string `json:"name,omitempty"`

	// Url URL for the emoji image
	Url *string `json:"url,omitempty"`
}

// Emojis defines model for Emojis.
type Emojis = []Emoji

// Error Structure for an error response
type Error struct {
	Error *string `json:"error,omitempty"`
}

// Event defines model for Event.
type Event struct {
	Id        *string `json:"id,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// ExternalAPIUser defines model for ExternalAPIUser.
type ExternalAPIUser struct {
	AccessToken  *string    `json:"accessToken,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	DisplayColor *int       `json:"displayColor,omitempty"`
	DisplayName  *string    `json:"displayName,omitempty"`
	Id           *string    `json:"id,omitempty"`
	IsBot        *bool      `json:"isBot,omitempty"`
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty"`
	Scopes       *[]string  `json:"scopes,omitempty"`
	Type         *string    `json:"type,omitempty"`
}

// ExternalAction defines model for ExternalAction.
type ExternalAction struct {
	Color          *string `json:"color,omitempty"`
	Description    *string `json:"description,omitempty"`
	Html           *string `json:"html,omitempty"`
	Icon           *string `json:"icon,omitempty"`
	OpenExternally *bool   `json:"openExternally,omitempty"`
	Title          *string `json:"title,omitempty"`
	Url            *string `json:"url,omitempty"`
}

// FederatedAction defines model for FederatedAction.
type FederatedAction struct {
	Body      *string `json:"body,omitempty"`
	Id        *string `json:"id,omitempty"`
	Image     *string `json:"image,omitempty"`
	Link      *string `json:"link,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	Title     *string `json:"title,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// FederationConfig defines model for FederationConfig.
type FederationConfig struct {
	Account       *string `json:"account,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
	FollowerCount *int    `json:"followerCount,omitempty"`
}

// Follower defines model for Follower.
type Follower struct {
	// DisabledAt When this follower was rejected or disabled.
	DisabledAt *string `json:"disabledAt,omitempty"`

	// Image The avatar image of the follower.
	Image *string `json:"image,omitempty"`

	// Link The IRI of the remote actor.
	Link *string `json:"link,omitempty"`

	// Name The display name of the follower.
	Name *string `json:"name,omitempty"`

	// Timestamp When this follow request was created.
	Timestamp *string `json:"timestamp,omitempty"`

	// Username The account username of the remote actor.
	Username *string `json:"username,omitempty"`
}

// Followers defines model for Followers.
type Followers = []Follower

// GeoDetails defines model for GeoDetails.
type GeoDetails struct {
	CountryCode *string `json:"countryCode,omitempty"`
	RegionName  *string `json:"regionName,omitempty"`
	TimeZone    *string `json:"timeZone,omitempty"`
}

// IPAddress defines model for IPAddress.
type IPAddress struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	IpAddress *string    `json:"ipAddress,omitempty"`
	Notes     *string    `json:"notes,omitempty"`
}

// InboundStreamDetails defines model for InboundStreamDetails.
type InboundStreamDetails struct {
	AudioBitrate *int     `json:"audioBitrate,omitempty"`
	AudioCodec   *string  `json:"audioCodec,omitempty"`
	Encoder      *string  `json:"encoder,omitempty"`
	Framerate    *float32 `json:"framerate,omitempty"`
	Height       *int     `json:"height,omitempty"`
	VideoBitrate *int     `json:"videoBitrate,omitempty"`
	VideoCodec   *string  `json:"videoCodec,omitempty"`
	Width        *int     `json:"width,omitempty"`
}

// LatencyLevel defines model for LatencyLevel.
type LatencyLevel struct {
	Level *int `json:"level,omitempty"`
}

// MessageEvent defines model for MessageEvent.
type MessageEvent struct {
	Body *string `json:"body,omitempty"`
}

// NotificationConfig defines model for NotificationConfig.
type NotificationConfig struct {
	Browser *BrowserConfig `json:"browser,omitempty"`
}

// PaginatedFollowers defines model for PaginatedFollowers.
type PaginatedFollowers struct {
	Results *Followers `json:"results,omitempty"`
	Total   *int       `json:"total,omitempty"`
}

// PlaybackMetrics defines model for PlaybackMetrics.
type PlaybackMetrics struct {
	Bandwidth             *float64 `json:"bandwidth,omitempty"`
	DownloadDuration      *float64 `json:"downloadDuration,omitempty"`
	Errors                *float64 `json:"errors,omitempty"`
	Latency               *float64 `json:"latency,omitempty"`
	QualityVariantChanges *float64 `json:"qualityVariantChanges,omitempty"`
}

// S3Info defines model for S3Info.
type S3Info struct {
	AccessKey      *string `json:"accessKey,omitempty"`
	Acl            *string `json:"acl,omitempty"`
	Bucket         *string `json:"bucket,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
	Endpoint       *string `json:"endpoint,omitempty"`
	ForcePathStyle *bool   `json:"forcePathStyle,omitempty"`
	PathPrefix     *string `json:"pathPrefix,omitempty"`
	Region         *string `json:"region,omitempty"`
	Secret         *string `json:"secret,omitempty"`
}

// SocialHandle defines model for SocialHandle.
type SocialHandle struct {
	Icon     *string `json:"icon,omitempty"`
	Platform *string `json:"platform,omitempty"`
	Url      *string `json:"url,omitempty"`
}

// Status Response for status
type Status struct {
	LastConnectTime    *string `json:"lastConnectTime,omitempty"`
	LastDisconnectTime *string `json:"lastDisconnectTime,omitempty"`
	Online             *bool   `json:"online,omitempty"`
	ServerTime         *string `json:"serverTime,omitempty"`
	VersionNumber      *string `json:"versionNumber,omitempty"`
	ViewerCount        *int    `json:"viewerCount,omitempty"`
}

// StreamHealthOverview defines model for StreamHealthOverview.
type StreamHealthOverview struct {
	HealthPercentage *int    `json:"healthPercentage,omitempty"`
	Healthy          *bool   `json:"healthy,omitempty"`
	Message          *string `json:"message,omitempty"`
	Representation   *int    `json:"representation,omitempty"`
}

// StreamKey defines model for StreamKey.
type StreamKey struct {
	Comment *string `json:"comment,omitempty"`
	Key     *string `json:"key,omitempty"`
}

// StreamOutputVariant defines model for StreamOutputVariant.
type StreamOutputVariant struct {
	AudioBitrate     *int    `json:"audioBitrate,omitempty"`
	AudioPassthrough *bool   `json:"audioPassthrough,omitempty"`
	CpuUsageLevel    *int    `json:"cpuUsageLevel,omitempty"`
	Framerate        *int    `json:"framerate,omitempty"`
	Name             *string `json:"name,omitempty"`
	ScaledHeight     *int    `json:"scaledHeight,omitempty"`
	ScaledWidth      *int    `json:"scaledWidth,omitempty"`
	VideoBitrate     *int    `json:"videoBitrate,omitempty"`
	VideoPassthrough *bool   `json:"videoPassthrough,omitempty"`
}

// SystemActionEvent defines model for SystemActionEvent.
type SystemActionEvent struct {
	Body      *string `json:"body,omitempty"`
	Id        *string `json:"id,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// SystemMessage defines model for SystemMessage.
type SystemMessage struct {
	Body      *string `json:"body,omitempty"`
	Id        *string `json:"id,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// TimestampedValue defines model for TimestampedValue.
type TimestampedValue struct {
	Time  *time.Time `json:"time,omitempty"`
	Value *float64   `json:"value,omitempty"`
}

// User defines model for User.
type User struct {
	Authenticated *bool     `json:"authenticated,omitempty"`
	CreatedAt     *string   `json:"createdAt,omitempty"`
	DisabledAt    *string   `json:"disabledAt,omitempty"`
	DisplayColor  *int      `json:"displayColor,omitempty"`
	DisplayName   *string   `json:"displayName,omitempty"`
	Id            *string   `json:"id,omitempty"`
	IsBot         *bool     `json:"isBot,omitempty"`
	NameChangedAt *string   `json:"nameChangedAt,omitempty"`
	PreviousNames *[]string `json:"previousNames,omitempty"`
	Scopes        *[]string `json:"scopes,omitempty"`
}

// UserEvent defines model for UserEvent.
type UserEvent struct {
	ClientId *int    `json:"clientId,omitempty"`
	HiddenAt *string `json:"hiddenAt,omitempty"`
	User     *User   `json:"user,omitempty"`
}

// UserMessage defines model for UserMessage.
type UserMessage struct {
	Body      *string `json:"body,omitempty"`
	ClientId  *int    `json:"clientId,omitempty"`
	HiddenAt  *string `json:"hiddenAt,omitempty"`
	Id        *string `json:"id,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	Type      *string `json:"type,omitempty"`
	User      *User   `json:"user,omitempty"`
}

// Users defines model for Users.
type Users = []User

// VideoVariant defines model for VideoVariant.
type VideoVariant struct {
	Index *int    `json:"index,omitempty"`
	Name  *string `json:"name,omitempty"`
}

// Viewer defines model for Viewer.
type Viewer struct {
	ClientID  *string     `json:"clientID,omitempty"`
	FirstSeen *time.Time  `json:"firstSeen,omitempty"`
	Geo       *GeoDetails `json:"geo,omitempty"`
	IpAddress *string     `json:"ipAddress,omitempty"`
	UserAgent *string     `json:"userAgent,omitempty"`
}

// WebConfig defines model for WebConfig.
type WebConfig struct {
	AppearanceVariables  *map[string]string    `json:"appearanceVariables,omitempty"`
	Authentication       *AuthenticationConfig `json:"authentication,omitempty"`
	ChatDisabled         *bool                 `json:"chatDisabled,omitempty"`
	CustomStyles         *string               `json:"customStyles,omitempty"`
	ExternalActions      *[]ExternalAction     `json:"externalActions,omitempty"`
	ExtraPageContent     *string               `json:"extraPageContent,omitempty"`
	Federation           *FederationConfig     `json:"federation,omitempty"`
	HideViewerCount      *bool                 `json:"hideViewerCount,omitempty"`
	Logo                 *string               `json:"logo,omitempty"`
	MaxSocketPayloadSize *int                  `json:"maxSocketPayloadSize,omitempty"`
	Name                 *string               `json:"name,omitempty"`
	Notifications        *NotificationConfig   `json:"notifications,omitempty"`
	Nsfw                 *bool                 `json:"nsfw,omitempty"`
	OfflineMessage       *string               `json:"offlineMessage,omitempty"`
	SocialHandles        *[]SocialHandle       `json:"socialHandles,omitempty"`
	SocketHostOverride   *string               `json:"socketHostOverride,omitempty"`
	StreamTitle          *string               `json:"streamTitle,omitempty"`
	Summary              *string               `json:"summary,omitempty"`
	Tags                 *[]string             `json:"tags,omitempty"`
	Version              *string               `json:"version,omitempty"`
}

// Webhook defines model for Webhook.
type Webhook struct {
	Events    *[]WebhookEventType `json:"events,omitempty"`
	Id        *int                `json:"id,omitempty"`
	LastUsed  *time.Time          `json:"lastUsed,omitempty"`
	Timestamp *time.Time          `json:"timestamp,omitempty"`
	Url       *string             `json:"url,omitempty"`
}

// WebhookEventType defines model for WebhookEventType.
type WebhookEventType string

// YPDetails defines model for YPDetails.
type YPDetails struct {
	Description           *string         `json:"description,omitempty"`
	LastConnectTime       *string         `json:"lastConnectTime,omitempty"`
	Logo                  *string         `json:"logo,omitempty"`
	Name                  *string         `json:"name,omitempty"`
	Nsfw                  *bool           `json:"nsfw,omitempty"`
	Online                *bool           `json:"online,omitempty"`
	OverallMaxViewerCount *int            `json:"overallMaxViewerCount,omitempty"`
	SessionMaxViewerCount *int            `json:"sessionMaxViewerCount,omitempty"`
	Social                *[]SocialHandle `json:"social,omitempty"`
	StreamTitle           *string         `json:"streamTitle,omitempty"`
	Tags                  *[]string       `json:"tags,omitempty"`
	ViewerCount           *int            `json:"viewerCount,omitempty"`
}

// AccessToken defines model for AccessToken.
type AccessToken = string

// Limit The numbers of items to return
type Limit = int

// Offset The number of items to skip before starting to collect the result set
type Offset = int

// N400 Simple API response
type N400 = BaseAPIResponse

// N500 Structure for an error response
type N500 = Error

// N501 Structure for an error response
type N501 = Error

// Default Simple API response
type Default = BaseAPIResponse

// CreateExternalAPIUserJSONBody defines parameters for CreateExternalAPIUser.
type CreateExternalAPIUserJSONBody struct {
	Name   *string   `json:"name,omitempty"`
	Scopes *[]string `json:"scopes,omitempty"`
}

// DeleteExternalAPIUserJSONBody defines parameters for DeleteExternalAPIUser.
type DeleteExternalAPIUserJSONBody struct {
	Token *string `json:"token,omitempty"`
}

// UpdateUserEnabledJSONBody defines parameters for UpdateUserEnabled.
type UpdateUserEnabledJSONBody struct {
	Enabled *bool   `json:"enabled,omitempty"`
	UserId  *string `json:"userId,omitempty"`
}

// UpdateUserModeratorJSONBody defines parameters for UpdateUserModerator.
type UpdateUserModeratorJSONBody struct {
	IsModerator *bool   `json:"isModerator,omitempty"`
	UserId      *string `json:"userId,omitempty"`
}

// SetCustomColorVariableValuesJSONBody defines parameters for SetCustomColorVariableValues.
type SetCustomColorVariableValuesJSONBody struct {
	Value *map[string]string `json:"value,omitempty"`
}

// SetForbiddenUsernameListJSONBody defines parameters for SetForbiddenUsernameList.
type SetForbiddenUsernameListJSONBody struct {
	Value *[]string `json:"value,omitempty"`
}

// SetSuggestedUsernameListJSONBody defines parameters for SetSuggestedUsernameList.
type SetSuggestedUsernameListJSONBody struct {
	Value *[]string `json:"value,omitempty"`
}

// SetStreamKeysJSONBody defines parameters for SetStreamKeys.
type SetStreamKeysJSONBody struct {
	Value *[]StreamKey `json:"value,omitempty"`
}

// DeleteCustomEmojiJSONBody defines parameters for DeleteCustomEmoji.
type DeleteCustomEmojiJSONBody struct {
	Name *string `json:"name,omitempty"`
}

// UploadCustomEmojiJSONBody defines parameters for UploadCustomEmoji.
type UploadCustomEmojiJSONBody struct {
	// Data Base 64 encoded image data
	Data *string `json:"data,omitempty"`
	Name *string `json:"name,omitempty"`
}

// GetFollowersAdminParams defines parameters for GetFollowersAdmin.
type GetFollowersAdminParams struct {
	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
	Limit  *Limit  `form:"limit,omitempty" json:"limit,omitempty"`
}

// ApproveFollowerJSONBody defines parameters for ApproveFollower.
type ApproveFollowerJSONBody struct {
	ActorIRI *string `json:"actorIRI,omitempty"`
	Approved *bool   `json:"approved,omitempty"`
}

// GetViewersOverTimeParams defines parameters for GetViewersOverTime.
type GetViewersOverTimeParams struct {
	// WindowStart Start date in unix time
	WindowStart *string `form:"windowStart,omitempty" json:"windowStart,omitempty"`
}

// CreateWebhookJSONBody defines parameters for CreateWebhook.
type CreateWebhookJSONBody struct {
	Events *[]WebhookEventType `json:"events,omitempty"`
	Url    *string             `json:"url,omitempty"`
}

// DeleteWebhookJSONBody defines parameters for DeleteWebhook.
type DeleteWebhookJSONBody struct {
	Id *int `json:"id,omitempty"`
}

// GetChatMessagesParams defines parameters for GetChatMessages.
type GetChatMessagesParams struct {
	AccessToken AccessToken `form:"accessToken" json:"accessToken"`
}

// RegisterAnonymousChatUserJSONBody defines parameters for RegisterAnonymousChatUser.
type RegisterAnonymousChatUserJSONBody struct {
	DisplayName *string `json:"displayName,omitempty"`
}

// RegisterAnonymousChatUserParams defines parameters for RegisterAnonymousChatUser.
type RegisterAnonymousChatUserParams struct {
	XForwardedUser *string `json:"X-Forwarded-User,omitempty"`
}

// GetFollowersParams defines parameters for GetFollowers.
type GetFollowersParams struct {
	Offset *Offset `form:"offset,omitempty" json:"offset,omitempty"`
	Limit  *Limit  `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostNotificationsRegisterJSONBody defines parameters for PostNotificationsRegister.
type PostNotificationsRegisterJSONBody struct {
	// Channel Name of notification channel
	Channel *string `json:"channel,omitempty"`

	// Destination Target of the notification in the channel
	Destination *string `json:"destination,omitempty"`
}

// PostNotificationsRegisterParams defines parameters for PostNotificationsRegister.
type PostNotificationsRegisterParams struct {
	AccessToken AccessToken `form:"accessToken" json:"accessToken"`
}

// RemoteFollowJSONBody defines parameters for RemoteFollow.
type RemoteFollowJSONBody struct {
	Account *string `json:"account,omitempty"`
}

// CreateExternalAPIUserJSONRequestBody defines body for CreateExternalAPIUser for application/json ContentType.
type CreateExternalAPIUserJSONRequestBody CreateExternalAPIUserJSONBody

// DeleteExternalAPIUserJSONRequestBody defines body for DeleteExternalAPIUser for application/json ContentType.
type DeleteExternalAPIUserJSONRequestBody DeleteExternalAPIUserJSONBody

// UpdateMessageVisibilityJSONRequestBody defines body for UpdateMessageVisibility for application/json ContentType.
type UpdateMessageVisibilityJSONRequestBody = AdminMessageVisibilityUpdate

// BanIPAddressJSONRequestBody defines body for BanIPAddress for application/json ContentType.
type BanIPAddressJSONRequestBody = AdminConfigValue

// UnbanIPAddressJSONRequestBody defines body for UnbanIPAddress for application/json ContentType.
type UnbanIPAddressJSONRequestBody = AdminConfigValue

// UpdateUserEnabledJSONRequestBody defines body for UpdateUserEnabled for application/json ContentType.
type UpdateUserEnabledJSONRequestBody UpdateUserEnabledJSONBody

// UpdateUserModeratorJSONRequestBody defines body for UpdateUserModerator for application/json ContentType.
type UpdateUserModeratorJSONRequestBody UpdateUserModeratorJSONBody

// SetAdminPasswordJSONRequestBody defines body for SetAdminPassword for application/json ContentType.
type SetAdminPasswordJSONRequestBody = AdminConfigValue

// SetCustomColorVariableValuesJSONRequestBody defines body for SetCustomColorVariableValues for application/json ContentType.
type SetCustomColorVariableValuesJSONRequestBody SetCustomColorVariableValuesJSONBody

// SetChatDisabledJSONRequestBody defines body for SetChatDisabled for application/json ContentType.
type SetChatDisabledJSONRequestBody = AdminConfigValue

// SetEnableEstablishedChatUserModeJSONRequestBody defines body for SetEnableEstablishedChatUserMode for application/json ContentType.
type SetEnableEstablishedChatUserModeJSONRequestBody = AdminConfigValue

// SetForbiddenUsernameListJSONRequestBody defines body for SetForbiddenUsernameList for application/json ContentType.
type SetForbiddenUsernameListJSONRequestBody SetForbiddenUsernameListJSONBody

// SetChatJoinMessagesEnabledJSONRequestBody defines body for SetChatJoinMessagesEnabled for application/json ContentType.
type SetChatJoinMessagesEnabledJSONRequestBody = AdminConfigValue

// SetSuggestedUsernameListJSONRequestBody defines body for SetSuggestedUsernameList for application/json ContentType.
type SetSuggestedUsernameListJSONRequestBody SetSuggestedUsernameListJSONBody

// SetServerNameJSONRequestBody defines body for SetServerName for application/json ContentType.
type SetServerNameJSONRequestBody = AdminConfigValue

// SetCustomOfflineMessageJSONRequestBody defines body for SetCustomOfflineMessage for application/json ContentType.
type SetCustomOfflineMessageJSONRequestBody = AdminConfigValue

// SetExtraPageContentJSONRequestBody defines body for SetExtraPageContent for application/json ContentType.
type SetExtraPageContentJSONRequestBody = AdminConfigValue

// SetServerSummaryJSONRequestBody defines body for SetServerSummary for application/json ContentType.
type SetServerSummaryJSONRequestBody = AdminConfigValue

// SetStreamKeysJSONRequestBody defines body for SetStreamKeys for application/json ContentType.
type SetStreamKeysJSONRequestBody SetStreamKeysJSONBody

// SetStreamTitleJSONRequestBody defines body for SetStreamTitle for application/json ContentType.
type SetStreamTitleJSONRequestBody = AdminConfigValue

// SetVideoCodecJSONRequestBody defines body for SetVideoCodec for application/json ContentType.
type SetVideoCodecJSONRequestBody = AdminConfigValue

// SetServerWelcomeMessageJSONRequestBody defines body for SetServerWelcomeMessage for application/json ContentType.
type SetServerWelcomeMessageJSONRequestBody = AdminConfigValue

// DeleteCustomEmojiJSONRequestBody defines body for DeleteCustomEmoji for application/json ContentType.
type DeleteCustomEmojiJSONRequestBody DeleteCustomEmojiJSONBody

// UploadCustomEmojiJSONRequestBody defines body for UploadCustomEmoji for application/json ContentType.
type UploadCustomEmojiJSONRequestBody UploadCustomEmojiJSONBody

// ApproveFollowerJSONRequestBody defines body for ApproveFollower for application/json ContentType.
type ApproveFollowerJSONRequestBody ApproveFollowerJSONBody

// CreateWebhookJSONRequestBody defines body for CreateWebhook for application/json ContentType.
type CreateWebhookJSONRequestBody CreateWebhookJSONBody

// DeleteWebhookJSONRequestBody defines body for DeleteWebhook for application/json ContentType.
type DeleteWebhookJSONRequestBody DeleteWebhookJSONBody

// RegisterAnonymousChatUserJSONRequestBody defines body for RegisterAnonymousChatUser for application/json ContentType.
type RegisterAnonymousChatUserJSONRequestBody RegisterAnonymousChatUserJSONBody

// SendChatActionJSONRequestBody defines body for SendChatAction for application/json ContentType.
type SendChatActionJSONRequestBody = SystemActionEvent

// ExternalUpdateMessageVisibilityJSONRequestBody defines body for ExternalUpdateMessageVisibility for application/json ContentType.
type ExternalUpdateMessageVisibilityJSONRequestBody = AdminMessageVisibilityUpdate

// SendIntegrationChatMessageJSONRequestBody defines body for SendIntegrationChatMessage for application/json ContentType.
type SendIntegrationChatMessageJSONRequestBody = UserMessage

// SendSystemMessageJSONRequestBody defines body for SendSystemMessage for application/json ContentType.
type SendSystemMessageJSONRequestBody = SystemMessage

// SendSystemMessageToConnectedClientJSONRequestBody defines body for SendSystemMessageToConnectedClient for application/json ContentType.
type SendSystemMessageToConnectedClientJSONRequestBody = SystemMessage

// ExternalSetStreamTitleJSONRequestBody defines body for ExternalSetStreamTitle for application/json ContentType.
type ExternalSetStreamTitleJSONRequestBody = AdminConfigValue

// PostMetricsPlaybackJSONRequestBody defines body for PostMetricsPlayback for application/json ContentType.
type PostMetricsPlaybackJSONRequestBody = PlaybackMetrics

// PostNotificationsRegisterJSONRequestBody defines body for PostNotificationsRegister for application/json ContentType.
type PostNotificationsRegisterJSONRequestBody PostNotificationsRegisterJSONBody

// RemoteFollowJSONRequestBody defines body for RemoteFollow for application/json ContentType.
type RemoteFollowJSONRequestBody RemoteFollowJSONBody

// AsAdminConfigValueValue0 returns the union data inside the AdminConfigValue_Value as a AdminConfigValueValue0
func (t AdminConfigValue_Value) AsAdminConfigValueValue0() (AdminConfigValueValue0, error) {
	var body AdminConfigValueValue0
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAdminConfigValueValue0 overwrites any union data inside the AdminConfigValue_Value as the provided AdminConfigValueValue0
func (t *AdminConfigValue_Value) FromAdminConfigValueValue0(v AdminConfigValueValue0) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAdminConfigValueValue0 performs a merge with any union data inside the AdminConfigValue_Value, using the provided AdminConfigValueValue0
func (t *AdminConfigValue_Value) MergeAdminConfigValueValue0(v AdminConfigValueValue0) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsAdminConfigValueValue1 returns the union data inside the AdminConfigValue_Value as a AdminConfigValueValue1
func (t AdminConfigValue_Value) AsAdminConfigValueValue1() (AdminConfigValueValue1, error) {
	var body AdminConfigValueValue1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAdminConfigValueValue1 overwrites any union data inside the AdminConfigValue_Value as the provided AdminConfigValueValue1
func (t *AdminConfigValue_Value) FromAdminConfigValueValue1(v AdminConfigValueValue1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAdminConfigValueValue1 performs a merge with any union data inside the AdminConfigValue_Value, using the provided AdminConfigValueValue1
func (t *AdminConfigValue_Value) MergeAdminConfigValueValue1(v AdminConfigValueValue1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t AdminConfigValue_Value) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *AdminConfigValue_Value) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsUserMessage returns the union data inside the ChatMessages_Item as a UserMessage
func (t ChatMessages_Item) AsUserMessage() (UserMessage, error) {
	var body UserMessage
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromUserMessage overwrites any union data inside the ChatMessages_Item as the provided UserMessage
func (t *ChatMessages_Item) FromUserMessage(v UserMessage) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeUserMessage performs a merge with any union data inside the ChatMessages_Item, using the provided UserMessage
func (t *ChatMessages_Item) MergeUserMessage(v UserMessage) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsSystemMessage returns the union data inside the ChatMessages_Item as a SystemMessage
func (t ChatMessages_Item) AsSystemMessage() (SystemMessage, error) {
	var body SystemMessage
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromSystemMessage overwrites any union data inside the ChatMessages_Item as the provided SystemMessage
func (t *ChatMessages_Item) FromSystemMessage(v SystemMessage) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeSystemMessage performs a merge with any union data inside the ChatMessages_Item, using the provided SystemMessage
func (t *ChatMessages_Item) MergeSystemMessage(v SystemMessage) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsActionMessage returns the union data inside the ChatMessages_Item as a ActionMessage
func (t ChatMessages_Item) AsActionMessage() (ActionMessage, error) {
	var body ActionMessage
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromActionMessage overwrites any union data inside the ChatMessages_Item as the provided ActionMessage
func (t *ChatMessages_Item) FromActionMessage(v ActionMessage) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeActionMessage performs a merge with any union data inside the ChatMessages_Item, using the provided ActionMessage
func (t *ChatMessages_Item) MergeActionMessage(v ActionMessage) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsFederatedAction returns the union data inside the ChatMessages_Item as a FederatedAction
func (t ChatMessages_Item) AsFederatedAction() (FederatedAction, error) {
	var body FederatedAction
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromFederatedAction overwrites any union data inside the ChatMessages_Item as the provided FederatedAction
func (t *ChatMessages_Item) FromFederatedAction(v FederatedAction) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeFederatedAction performs a merge with any union data inside the ChatMessages_Item, using the provided FederatedAction
func (t *ChatMessages_Item) MergeFederatedAction(v FederatedAction) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t ChatMessages_Item) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ChatMessages_Item) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
