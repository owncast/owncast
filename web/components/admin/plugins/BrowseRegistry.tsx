import React, { useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Empty,
  Input,
  Select,
  Space,
  Spin,
  Tag,
  Tooltip,
  Typography,
} from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../../types/localization';
import { isPluginUpdateAvailable } from '../../../utils/apis';
import { formatLinkHostname } from '../../../utils/format';
import { readableBytes } from '../../../utils/images';
import { PuzzlePiece } from './PluginIcon';
import { permissionDescriptionKey, permissionNameKey } from './permissionDescriptions';
import s from './BrowseRegistry.module.scss';

const { Text, Paragraph } = Typography;

// RegistryManifest is the embedded plugin.manifest.json from inside the
// .ocpkg, surfaced by the registry so the Browse tab can preview what
// a plugin asks for before the admin clicks Install. Only the fields
// the UI actually renders are typed here; the registry can send extra
// fields without breaking the page.
export type RegistryManifest = {
  permissions?: string[];
  // Canonical category slug (e.g. "chat-bots"); see categoryNameKey.
  category?: string;
};

// One published plugin in the registry's browse payload. Mirrors the
// publishedPluginView shape from the directory's plugin_browse.go;
// keep these aligned. `slug` is the canonical identifier (URL segment,
// install request body); `name` is the human-readable display name.
export type RegistryPlugin = {
  slug: string;
  name: string;
  summary?: string;
  homepage?: string;
  tags?: string[];
  iconURL?: string;
  // Author-supplied screenshot/preview image for the listing. Wider
  // than the square icon; rendered below the summary when present.
  previewURL?: string;
  authorName?: string;
  // Optional public link for the author (personal site, GitHub
  // profile). When present the "by ..." line links to it.
  authorURL?: string;
  // Canonical category slug from the taxonomy below. The registry
  // mirrors it top-level from the manifest; older registries only
  // embed it inside latest.manifest, so read via pluginCategory().
  category?: string;
  latest?: {
    version: string;
    sizeBytes?: number;
    sha256?: string;
    // Embedded manifest for the latest version. The registry inlines
    // it so Browse cards can show permissions (and other manifest
    // metadata) without a second round-trip per card.
    manifest?: RegistryManifest;
  };
  versions?: { version: string }[];
};

// Canonical plugin category taxonomy (slug -> label i18n key). Keep in
// lock-step with the registry's category list and
// Localization.Admin.Plugins.Categories.
const categoryNameKey: Record<string, string> = {
  'chat-bots': Localization.Admin.Plugins.Categories.chatBots,
  'chat-filters': Localization.Admin.Plugins.Categories.chatFilters,
  moderation: Localization.Admin.Plugins.Categories.moderation,
  authentication: Localization.Admin.Plugins.Categories.authentication,
  themes: Localization.Admin.Plugins.Categories.themes,
  overlays: Localization.Admin.Plugins.Categories.overlays,
  notifications: Localization.Admin.Plugins.Categories.notifications,
  integrations: Localization.Admin.Plugins.Categories.integrations,
  video: Localization.Admin.Plugins.Categories.video,
  analytics: Localization.Admin.Plugins.Categories.analytics,
  games: Localization.Admin.Plugins.Categories.games,
  'admin-utilities': Localization.Admin.Plugins.Categories.adminUtilities,
  examples: Localization.Admin.Plugins.Categories.examples,
  other: Localization.Admin.Plugins.Categories.other,
};

// A plugin's category, preferring the registry's top-level mirror and
// falling back to the embedded manifest so catalogs from a registry
// that predates the top-level field still filter.
const pluginCategory = (p: RegistryPlugin) => p.category ?? p.latest?.manifest?.category;

export type BrowseRegistryProps = {
  // Map of installed plugin slug -> currently-installed version, so
  // each card can decide whether to render Install / Installed /
  // Update.
  installedVersions: Map<string, string>;
  // The registry's published plugins. Fetched by the parent so the
  // same data backs the Browse tab and the Installed tab's
  // "update available" tags.
  registry: RegistryPlugin[];
  registryLoading: boolean;
  // Non-null when the registry fetch failed (network error, host
  // returned a non-2xx, etc.). Renders the "catalog unavailable"
  // state with a retry button. `null` means the fetch succeeded
  // (the array may still be empty: that's the "no plugins published
  // yet" state, distinct from "catalog is offline").
  registryError?: string | null;
  // Called when the admin clicks Install or Update. The parent runs
  // the registry-install endpoint, refreshes both lists, and (for
  // non-enabled plugins) opens the InstallConfirmModal.
  onInstall: (slug: string, version: string) => Promise<void>;
  // Called when the admin clicks Retry after a catalog fetch
  // failure. Should re-attempt the registry list fetch.
  onRetry?: () => void;
};

// BrowseRegistry renders the publicly-available plugins as a list of
// cards inside the Browse tab. The bytes never pass through the
// browser: the Owncast host downloads from the registry, verifies
// SHA256, and runs the same Install path as a manual .ocpkg upload.
export const BrowseRegistry = ({
  installedVersions,
  registry,
  registryLoading,
  registryError = null,
  onInstall,
  onRetry,
}: BrowseRegistryProps) => {
  const { t } = useTranslation();
  // Names whose install POST is in flight; per-card spinner state so
  // multiple installs don't lock the whole list.
  const [installing, setInstalling] = useState<Set<string>>(new Set());
  // Client-side filters over the loaded registry payload. Options are
  // derived from what's actually present so the controls never offer a
  // choice that matches nothing.
  const [permissionFilter, setPermissionFilter] = useState<string[]>([]);
  const [authorFilter, setAuthorFilter] = useState<string>();
  const [categoryFilter, setCategoryFilter] = useState<string>();
  // Free-text search over plugin names and summaries.
  const [searchTerm, setSearchTerm] = useState('');

  // Unknown slug falls back to the raw value so a newer registry
  // doesn't render blank labels.
  const categoryLabel = (slug: string) => (categoryNameKey[slug] ? t(categoryNameKey[slug]) : slug);

  const permissionOptions = useMemo(() => {
    const perms = new Set<string>();
    registry.forEach(p => p.latest?.manifest?.permissions?.forEach(perm => perms.add(perm)));
    return [...perms]
      .map(perm => ({
        value: perm,
        label: permissionNameKey[perm] ? t(permissionNameKey[perm]) : perm,
      }))
      .sort((a, b) => a.label.localeCompare(b.label));
  }, [registry, t]);

  const authorOptions = useMemo(() => {
    const authors = new Set<string>();
    registry.forEach(p => p.authorName && authors.add(p.authorName));
    return [...authors].sort().map(name => ({ value: name, label: name }));
  }, [registry]);

  // Empty until the registry starts publishing categories; the whole
  // category Select is hidden in that case.
  const categoryOptions = useMemo(() => {
    const categories = new Set<string>();
    registry.forEach(p => {
      const category = pluginCategory(p);
      if (category) categories.add(category);
    });
    return [...categories]
      .map(slug => ({
        value: slug,
        label: categoryNameKey[slug] ? t(categoryNameKey[slug]) : slug,
      }))
      .sort((a, b) => a.label.localeCompare(b.label));
  }, [registry, t]);

  const filteredRegistry = useMemo(
    () =>
      registry.filter(p => {
        const search = searchTerm.trim().toLowerCase();
        if (
          search &&
          !p.name.toLowerCase().includes(search) &&
          !(p.summary ?? '').toLowerCase().includes(search)
        ) {
          return false;
        }
        if (authorFilter && p.authorName !== authorFilter) return false;
        if (categoryFilter && pluginCategory(p) !== categoryFilter) return false;
        const perms = p.latest?.manifest?.permissions ?? [];
        // A plugin matches when it declares ALL selected permissions.
        return permissionFilter.every(perm => perms.includes(perm));
      }),
    [registry, permissionFilter, authorFilter, categoryFilter, searchTerm],
  );

  const clearFilters = () => {
    setPermissionFilter([]);
    setAuthorFilter(undefined);
    setCategoryFilter(undefined);
    setSearchTerm('');
  };

  const triggerInstall = async (plugin: RegistryPlugin) => {
    if (!plugin.latest) return;
    const { slug } = plugin;
    setInstalling(prev => new Set(prev).add(slug));
    try {
      await onInstall(slug, plugin.latest.version);
    } finally {
      setInstalling(prev => {
        const next = new Set(prev);
        next.delete(slug);
        return next;
      });
    }
  };

  if (registryLoading) {
    return (
      <div className={s.loader}>
        <Spin />
      </div>
    );
  }

  // Catalog unreachable: distinct from the empty-but-reachable state.
  // Shows the actual error (network, 502, etc.) so an admin
  // troubleshooting a misconfigured OWNCAST_PLUGIN_REGISTRY env or a
  // down upstream has something concrete to act on.
  if (registryError !== null) {
    return (
      <Alert
        type="warning"
        showIcon
        message={t(Localization.Admin.Plugins.browseUnavailableTitle)}
        description={
          <Space direction="vertical" size={8}>
            <Text>{t(Localization.Admin.Plugins.browseUnavailableDescription)}</Text>
            <Text type="secondary">
              <code>{registryError}</code>
            </Text>
            {onRetry && (
              <Button size="small" onClick={onRetry}>
                {t(Localization.Admin.Plugins.browseUnavailableRetry)}
              </Button>
            )}
          </Space>
        }
      />
    );
  }

  if (registry.length === 0) {
    return <Empty description={t(Localization.Admin.Plugins.browseEmpty)} />;
  }

  const filterBar = (
    <Space wrap className={s.filters}>
      <Input
        allowClear
        className={s.filterSelect}
        placeholder={t(Localization.Admin.Plugins.browseSearch)}
        value={searchTerm}
        onChange={e => setSearchTerm(e.target.value)}
      />
      {categoryOptions.length > 0 && (
        <Select
          allowClear
          className={s.filterSelect}
          placeholder={t(Localization.Admin.Plugins.browseFilterCategory)}
          value={categoryFilter}
          onChange={setCategoryFilter}
          options={categoryOptions}
        />
      )}
      <Select
        mode="multiple"
        allowClear
        className={s.filterSelect}
        placeholder={t(Localization.Admin.Plugins.browseFilterPermissions)}
        value={permissionFilter}
        onChange={setPermissionFilter}
        options={permissionOptions}
      />
      <Select
        allowClear
        className={s.filterSelect}
        placeholder={t(Localization.Admin.Plugins.browseFilterAuthor)}
        value={authorFilter}
        onChange={setAuthorFilter}
        options={authorOptions}
      />
    </Space>
  );

  // Filters excluded everything: keep the bar visible so the admin can
  // adjust, and offer a one-click reset.
  if (filteredRegistry.length === 0) {
    return (
      <>
        {filterBar}
        <Empty description={t(Localization.Admin.Plugins.browseFilteredEmpty)}>
          <Button onClick={clearFilters}>{t(Localization.Admin.Plugins.browseClearFilters)}</Button>
        </Empty>
      </>
    );
  }

  return (
    <>
      {filterBar}
      <Space direction="vertical" size="middle" className={s.list}>
        {filteredRegistry.map(plugin => {
          const installedVersion = installedVersions.get(plugin.slug);
          const latestVersion = plugin.latest?.version;
          const isInstalled = installedVersion !== undefined;
          const hasUpdate = isInstalled && isPluginUpdateAvailable(installedVersion, latestVersion);

          let actionButton: React.ReactNode;
          if (!plugin.latest) {
            actionButton = <Button disabled>{t(Localization.Admin.Plugins.browseInstall)}</Button>;
          } else if (hasUpdate) {
            actionButton = (
              <Button
                type="primary"
                loading={installing.has(plugin.slug)}
                onClick={() => triggerInstall(plugin)}
              >
                {t(Localization.Admin.Plugins.browseUpdate, { version: latestVersion })}
              </Button>
            );
          } else if (isInstalled) {
            actionButton = (
              <Button disabled>{t(Localization.Admin.Plugins.browseInstalled)}</Button>
            );
          } else {
            actionButton = (
              <Button
                type="primary"
                loading={installing.has(plugin.slug)}
                onClick={() => triggerInstall(plugin)}
              >
                {t(Localization.Admin.Plugins.browseInstall)}
              </Button>
            );
          }

          return (
            <Card key={plugin.slug} size="small">
              <div className={s.row}>
                <div className={s.icon}>
                  {plugin.iconURL ? (
                    <img src={plugin.iconURL} alt="" />
                  ) : (
                    <PuzzlePiece className={s.iconFallback} />
                  )}
                </div>
                <div className={s.body}>
                  <div className={s.title}>
                    <strong>{plugin.name}</strong>
                    {plugin.latest && <Text type="secondary"> v{plugin.latest.version}</Text>}
                    {plugin.latest?.sizeBytes ? (
                      <Text type="secondary">
                        {' '}
                        &middot; {readableBytes(plugin.latest.sizeBytes)}
                      </Text>
                    ) : null}
                    {pluginCategory(plugin) && (
                      <Tag className={s.categoryTag}>{categoryLabel(pluginCategory(plugin))}</Tag>
                    )}
                  </div>
                  {plugin.authorName &&
                    (plugin.authorURL ? (
                      <Typography.Link
                        href={plugin.authorURL}
                        target="_blank"
                        rel="noopener noreferrer"
                        className={s.author}
                      >
                        {t(Localization.Admin.Plugins.browseAuthor, { name: plugin.authorName })}
                      </Typography.Link>
                    ) : (
                      <Text type="secondary" className={s.author}>
                        {t(Localization.Admin.Plugins.browseAuthor, { name: plugin.authorName })}
                      </Text>
                    ))}
                  {plugin.summary && <Paragraph className={s.summary}>{plugin.summary}</Paragraph>}
                  {plugin.homepage && (
                    <Typography.Link
                      href={plugin.homepage}
                      target="_blank"
                      rel="noopener noreferrer"
                      className={s.homepage}
                    >
                      {formatLinkHostname(plugin.homepage)}
                    </Typography.Link>
                  )}
                  {plugin.previewURL && (
                    <img
                      className={s.preview}
                      src={plugin.previewURL}
                      alt={t(Localization.Admin.Plugins.browsePreviewAlt, { name: plugin.name })}
                      loading="lazy"
                    />
                  )}
                  {plugin.latest?.manifest?.permissions &&
                    plugin.latest.manifest.permissions.length > 0 && (
                      // Permissions the plugin's manifest declares, so the
                      // admin sees the scope they'd be granting before
                      // clicking Install. Each tag shows the short label;
                      // the tooltip carries the full plain-language
                      // description. Same maps as the Installed tab so the
                      // copy stays in lock-step between views.
                      <Space size={[4, 4]} wrap className={s.permissions}>
                        {plugin.latest.manifest.permissions.map(perm => {
                          const nameKey = permissionNameKey[perm];
                          const descKey = permissionDescriptionKey[perm];
                          const label = nameKey ? t(nameKey) : perm;
                          const description = descKey ? t(descKey) : perm;
                          return (
                            <Tooltip key={perm} title={description}>
                              <Tag>{label}</Tag>
                            </Tooltip>
                          );
                        })}
                      </Space>
                    )}
                  {plugin.tags && plugin.tags.length > 0 && (
                    <Space size={[4, 4]} wrap>
                      {plugin.tags.map(tag => (
                        <Tag key={tag}>{tag}</Tag>
                      ))}
                    </Space>
                  )}
                </div>
                <div className={s.actions}>{actionButton}</div>
              </div>
            </Card>
          );
        })}
      </Space>
    </>
  );
};
