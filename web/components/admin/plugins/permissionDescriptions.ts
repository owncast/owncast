import { Localization } from '../../../types/localization';

// permissionDescriptionKey maps a plugin permission identifier (e.g.
// "chat.send") to the i18n key for its plain-language description. Used
// in the plugins list (as tooltips on each permission tag) and in the
// per-plugin detail view (as the description column in the Permissions
// tab). Mirrors the permission constants in services/plugins/hostfns.go;
// keep in lock-step with Localization.Admin.Plugins.Permissions.
export const permissionDescriptionKey: Record<string, string> = {
  'storage.kv': Localization.Admin.Plugins.Permissions.storageKv,
  'storage.upload': Localization.Admin.Plugins.Permissions.storageUpload,
  'chat.send': Localization.Admin.Plugins.Permissions.chatSend,
  'chat.history': Localization.Admin.Plugins.Permissions.chatHistory,
  'chat.moderate': Localization.Admin.Plugins.Permissions.chatModerate,
  'network.fetch': Localization.Admin.Plugins.Permissions.networkFetch,
  'events.emit': Localization.Admin.Plugins.Permissions.eventsEmit,
  'http.serve': Localization.Admin.Plugins.Permissions.httpServe,
  'http.sse': Localization.Admin.Plugins.Permissions.httpSse,
  'server.read': Localization.Admin.Plugins.Permissions.serverRead,
  'notifications.send': Localization.Admin.Plugins.Permissions.notificationsSend,
  'users.read': Localization.Admin.Plugins.Permissions.usersRead,
  'users.moderate': Localization.Admin.Plugins.Permissions.usersModerate,
  'fediverse.post': Localization.Admin.Plugins.Permissions.fediversePost,
  'videoconfig.read': Localization.Admin.Plugins.Permissions.videoconfigRead,
  'videoconfig.write': Localization.Admin.Plugins.Permissions.videoconfigWrite,
  'ui.modify': Localization.Admin.Plugins.Permissions.uiModify,
  'chat.filter': Localization.Admin.Plugins.Permissions.chatFilter,
};

// permissionNameKey maps a permission identifier to the i18n key for
// its short plain-language label (e.g. "Send chat messages"). Used on
// the permission Tags in the plugins list, alongside the full
// description in a hover tooltip.
export const permissionNameKey: Record<string, string> = {
  'storage.kv': Localization.Admin.Plugins.PermissionNames.storageKv,
  'storage.upload': Localization.Admin.Plugins.PermissionNames.storageUpload,
  'chat.send': Localization.Admin.Plugins.PermissionNames.chatSend,
  'chat.history': Localization.Admin.Plugins.PermissionNames.chatHistory,
  'chat.moderate': Localization.Admin.Plugins.PermissionNames.chatModerate,
  'network.fetch': Localization.Admin.Plugins.PermissionNames.networkFetch,
  'events.emit': Localization.Admin.Plugins.PermissionNames.eventsEmit,
  'http.serve': Localization.Admin.Plugins.PermissionNames.httpServe,
  'http.sse': Localization.Admin.Plugins.PermissionNames.httpSse,
  'server.read': Localization.Admin.Plugins.PermissionNames.serverRead,
  'notifications.send': Localization.Admin.Plugins.PermissionNames.notificationsSend,
  'users.read': Localization.Admin.Plugins.PermissionNames.usersRead,
  'users.moderate': Localization.Admin.Plugins.PermissionNames.usersModerate,
  'fediverse.post': Localization.Admin.Plugins.PermissionNames.fediversePost,
  'videoconfig.read': Localization.Admin.Plugins.PermissionNames.videoconfigRead,
  'videoconfig.write': Localization.Admin.Plugins.PermissionNames.videoconfigWrite,
  'ui.modify': Localization.Admin.Plugins.PermissionNames.uiModify,
  'chat.filter': Localization.Admin.Plugins.PermissionNames.chatFilter,
};
