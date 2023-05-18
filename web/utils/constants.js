// misc constants used throughout the app

export const URL_STATUS = `/api/status`;
export const URL_CHAT_HISTORY = `/api/chat`;
export const URL_CUSTOM_EMOJIS = `/api/emoji`;
export const URL_CONFIG = `/api/config`;
export const URL_VIEWER_PING = `/api/ping`;

// inline moderation actions
export const URL_HIDE_MESSAGE = `/api/chat/messagevisibility`;
export const URL_BAN_USER = `/api/chat/users/setenabled`;

// TODO: This directory is customizable in the config.  So we should expose this via the config API.
export const URL_STREAM = `/hls/stream.m3u8`;
// export const URL_WEBSOCKET = `${location.protocol === 'https:' ? 'wss' : 'ws'}://${
//   location.host
// }/ws`;
export const URL_CHAT_REGISTRATION = `/api/chat/register`;
export const URL_FOLLOWERS = `/api/followers`;
export const URL_PLAYBACK_METRICS = `/api/metrics/playback`;

export const URL_REGISTER_NOTIFICATION = `/api/notifications/register`;
export const URL_REGISTER_EMAIL_NOTIFICATION = `/api/notifications/register/email`;
export const URL_CHAT_INDIEAUTH_BEGIN = `/api/auth/indieauth`;

export const TIMER_STATUS_UPDATE = 5000; // ms
export const TIMER_DISABLE_CHAT_AFTER_OFFLINE = 5 * 60 * 1000; // 5 mins
export const TIMER_STREAM_DURATION_COUNTER = 1000;
export const TEMP_IMAGE =
  'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';

export const OWNCAST_LOGO_LOCAL = '/img/logo.svg';

export const MESSAGE_OFFLINE = 'Stream is offline.';
export const MESSAGE_ONLINE = 'Stream is online.';

export const DYNAMIC_PADDING_VALUE = '320px';
