import { PluginPermission } from '../../../interfaces/plugin';
import { Localization } from '../../../types/localization';

// permissionDescriptionKey maps a plugin permission identifier to the
// i18n key for its plain-language description. Used in the plugins
// list (as tooltips on each permission tag) and in the per-plugin
// detail view (as the description column in the Permissions tab).
// Keep in lock-step with Localization.Admin.Plugins.Permissions and
// the host-side constants in services/plugins/hostfns.go.
export const permissionDescriptionKey: Record<string, string> = {
  [PluginPermission.StorageKV]: Localization.Admin.Plugins.Permissions.storageKv,
  [PluginPermission.StorageUpload]: Localization.Admin.Plugins.Permissions.storageUpload,
  [PluginPermission.StorageFS]: Localization.Admin.Plugins.Permissions.storageFs,
  [PluginPermission.ChatSend]: Localization.Admin.Plugins.Permissions.chatSend,
  [PluginPermission.ChatHistory]: Localization.Admin.Plugins.Permissions.chatHistory,
  [PluginPermission.ChatModerate]: Localization.Admin.Plugins.Permissions.chatModerate,
  [PluginPermission.NetworkFetch]: Localization.Admin.Plugins.Permissions.networkFetch,
  [PluginPermission.EventsEmit]: Localization.Admin.Plugins.Permissions.eventsEmit,
  [PluginPermission.HttpServe]: Localization.Admin.Plugins.Permissions.httpServe,
  [PluginPermission.HttpSse]: Localization.Admin.Plugins.Permissions.httpSse,
  [PluginPermission.ServerRead]: Localization.Admin.Plugins.Permissions.serverRead,
  [PluginPermission.NotificationsSend]: Localization.Admin.Plugins.Permissions.notificationsSend,
  [PluginPermission.UsersRead]: Localization.Admin.Plugins.Permissions.usersRead,
  [PluginPermission.UsersModerate]: Localization.Admin.Plugins.Permissions.usersModerate,
  [PluginPermission.FediversePost]: Localization.Admin.Plugins.Permissions.fediversePost,
  [PluginPermission.VideoConfigRead]: Localization.Admin.Plugins.Permissions.videoconfigRead,
  [PluginPermission.VideoConfigWrite]: Localization.Admin.Plugins.Permissions.videoconfigWrite,
  [PluginPermission.UIModify]: Localization.Admin.Plugins.Permissions.uiModify,
  [PluginPermission.ChatFilter]: Localization.Admin.Plugins.Permissions.chatFilter,
};

// permissionNameKey maps a permission identifier to the i18n key for
// its short plain-language label (e.g. "Send chat messages"). Used on
// the permission Tags in the plugins list, alongside the full
// description in a hover tooltip.
export const permissionNameKey: Record<string, string> = {
  [PluginPermission.StorageKV]: Localization.Admin.Plugins.PermissionNames.storageKv,
  [PluginPermission.StorageUpload]: Localization.Admin.Plugins.PermissionNames.storageUpload,
  [PluginPermission.StorageFS]: Localization.Admin.Plugins.PermissionNames.storageFs,
  [PluginPermission.ChatSend]: Localization.Admin.Plugins.PermissionNames.chatSend,
  [PluginPermission.ChatHistory]: Localization.Admin.Plugins.PermissionNames.chatHistory,
  [PluginPermission.ChatModerate]: Localization.Admin.Plugins.PermissionNames.chatModerate,
  [PluginPermission.NetworkFetch]: Localization.Admin.Plugins.PermissionNames.networkFetch,
  [PluginPermission.EventsEmit]: Localization.Admin.Plugins.PermissionNames.eventsEmit,
  [PluginPermission.HttpServe]: Localization.Admin.Plugins.PermissionNames.httpServe,
  [PluginPermission.HttpSse]: Localization.Admin.Plugins.PermissionNames.httpSse,
  [PluginPermission.ServerRead]: Localization.Admin.Plugins.PermissionNames.serverRead,
  [PluginPermission.NotificationsSend]:
    Localization.Admin.Plugins.PermissionNames.notificationsSend,
  [PluginPermission.UsersRead]: Localization.Admin.Plugins.PermissionNames.usersRead,
  [PluginPermission.UsersModerate]: Localization.Admin.Plugins.PermissionNames.usersModerate,
  [PluginPermission.FediversePost]: Localization.Admin.Plugins.PermissionNames.fediversePost,
  [PluginPermission.VideoConfigRead]: Localization.Admin.Plugins.PermissionNames.videoconfigRead,
  [PluginPermission.VideoConfigWrite]: Localization.Admin.Plugins.PermissionNames.videoconfigWrite,
  [PluginPermission.UIModify]: Localization.Admin.Plugins.PermissionNames.uiModify,
  [PluginPermission.ChatFilter]: Localization.Admin.Plugins.PermissionNames.chatFilter,
};
