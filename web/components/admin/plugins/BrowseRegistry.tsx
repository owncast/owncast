import React, { useState } from 'react';
import { Alert, Button, Card, Empty, Space, Spin, Tag, Tooltip, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../../types/localization';
import { isPluginUpdateAvailable } from '../../../utils/apis';
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

  return (
    <Space direction="vertical" size="middle" className={s.list}>
      {registry.map(plugin => {
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
          actionButton = <Button disabled>{t(Localization.Admin.Plugins.browseInstalled)}</Button>;
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
                    <Text type="secondary"> &middot; {readableBytes(plugin.latest.sizeBytes)}</Text>
                  ) : null}
                </div>
                {plugin.authorName && (
                  <Text type="secondary" className={s.author}>
                    {t(Localization.Admin.Plugins.browseAuthor, { name: plugin.authorName })}
                  </Text>
                )}
                {plugin.summary && <Paragraph className={s.summary}>{plugin.summary}</Paragraph>}
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
  );
};
