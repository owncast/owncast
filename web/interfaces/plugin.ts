// PluginPermission enumerates the permission identifiers the host
// accepts in plugin.manifest.json. Mirrors the constants in the host's
// services/plugins/hostfns.go; keep in lock-step.
//
// Use this enum anywhere the frontend needs to reason about a specific
// permission (e.g. branching on `network.fetch` to surface allowed
// hosts). Avoid bare string literals so refactors and rename checks
// are mechanical, and the typechecker catches typos.
export const PluginPermission = {
  StorageKV: 'storage.kv',
  StorageUpload: 'storage.upload',
  ChatSend: 'chat.send',
  ChatHistory: 'chat.history',
  ChatModerate: 'chat.moderate',
  ChatFilter: 'chat.filter',
  UsersRead: 'users.read',
  UsersModerate: 'users.moderate',
  NetworkFetch: 'network.fetch',
  EventsEmit: 'events.emit',
  HttpServe: 'http.serve',
  HttpSse: 'http.sse',
  ServerRead: 'server.read',
  VideoConfigRead: 'videoconfig.read',
  VideoConfigWrite: 'videoconfig.write',
  NotificationsSend: 'notifications.send',
  FediversePost: 'fediverse.post',
  UIModify: 'ui.modify',
} as const;

// The const above and this type intentionally share a name: TypeScript
// merges value-space and type-space declarations, so `PluginPermission`
// works as both an enum-like namespace and a string-literal union. The
// eslint no-redeclare rule doesn't model this; suppress it here.
// eslint-disable-next-line no-redeclare
export type PluginPermission = (typeof PluginPermission)[keyof typeof PluginPermission];

// Plugin describes a discovered plugin as returned by GET /api/admin/plugins.
// Mirrors services/plugins.DiscoveredEntry on the backend. Two name-like
// fields: `slug` is the canonical identifier used in URLs and action
// endpoints; `name` is the human-readable display name shown in lists.
export interface Plugin {
  slug: string;
  name: string;
  // botDisplayName, when present, is what the plugin uses as its chat
  // identity. Empty/undefined means "use name". Shown in the admin's
  // plugin list so the operator knows what viewers will see in chat.
  botDisplayName?: string;
  version?: string;
  description?: string;
  permissions?: string[];
  path: string;
  enabled: boolean;
  loaded: boolean;
  // autoDisabled is set when the host's strike system stopped invoking
  // the plugin after consecutive filter failures. The admin's enabled
  // choice is preserved, but the plugin isn't doing any work until it's
  // reloaded or rebuilt.
  autoDisabled?: boolean;
  // hasIcon is true when the plugin ships an icon.png alongside its
  // manifest. The admin UI fetches the bytes from
  // /api/plugins/<slug>/icon and renders them in the list and sidebar.
  hasIcon?: boolean;
  // hasInstructions is true when the plugin ships an INSTRUCTIONS.md
  // alongside its manifest. The admin UI fetches the markdown from
  // /api/admin/plugins/<slug>/instructions and renders it in an
  // Instructions tab on the plugin's details page.
  hasInstructions?: boolean;
  // pendingPermissions lists permissions the manifest now declares that
  // the admin has not yet approved. Non-empty means the plugin was
  // updated on disk to ask for more access than was originally granted;
  // the plugin is held in a not-loaded state until the admin re-enables
  // it (which captures a fresh approval snapshot covering the new set).
  pendingPermissions?: string[];
  // allowedHosts mirrors manifest.network.allowedHosts. Populated when
  // the plugin requests network.fetch; the admin UI surfaces this list
  // alongside the network.fetch permission entry so an admin reviewing
  // the plugin sees the host scope the plugin is allowed to reach.
  allowedHosts?: string[];
  lastError?: string;
  discoveredAt: string;
  adminPages?: PluginAdminPage[];
}

// PluginAdminPage is a single admin-only page declared in a plugin's
// manifest.admin.pages entry. The body is rendered as an iframe to
// /plugins/<slug><path>.
export interface PluginAdminPage {
  title: string;
  path: string;
  icon?: string;
}
