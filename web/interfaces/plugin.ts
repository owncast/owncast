// Plugin describes a discovered plugin as returned by GET /api/admin/plugins.
// Mirrors services/plugins.DiscoveredEntry on the backend.
export interface Plugin {
  name: string;
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
  // /api/plugins/<name>/icon and renders them in the list and sidebar.
  hasIcon?: boolean;
  // pendingPermissions lists permissions the manifest now declares that
  // the admin has not yet approved. Non-empty means the plugin was
  // updated on disk to ask for more access than was originally granted;
  // the plugin is held in a not-loaded state until the admin re-enables
  // it (which captures a fresh approval snapshot covering the new set).
  pendingPermissions?: string[];
  lastError?: string;
  discoveredAt: string;
  adminPages?: PluginAdminPage[];
}

// PluginAdminPage is a single admin-only page declared in a plugin's
// manifest.admin.pages entry. The body is rendered as an iframe to
// /plugins/<plugin-name><path>.
export interface PluginAdminPage {
  title: string;
  path: string;
  icon?: string;
}
