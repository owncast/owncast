import React, { useState } from 'react';
import { Button, Card, Empty, Space, Spin, Tag, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../../types/localization';
import { PuzzlePiece } from './PluginIcon';
import s from './BrowseRegistry.module.scss';

const { Text, Paragraph } = Typography;

// One published plugin in the registry's browse payload. Mirrors the
// publishedPluginView shape from the directory's plugin_browse.go;
// keep these aligned.
export type RegistryPlugin = {
  name: string;
  summary?: string;
  homepage?: string;
  tags?: string[];
  iconURL?: string;
  authorName?: string;
  latest?: { version: string; sizeBytes?: number; sha256?: string };
  versions?: { version: string }[];
};

export type BrowseRegistryProps = {
  // Map of installed plugin name -> currently-installed version, so
  // each card can decide whether to render Install / Installed /
  // Update.
  installedVersions: Map<string, string>;
  // The registry's published plugins. Fetched by the parent so the
  // same data backs the Browse tab and the Installed tab's
  // "update available" tags.
  registry: RegistryPlugin[];
  registryLoading: boolean;
  // Called when the admin clicks Install or Update. The parent runs
  // the registry-install endpoint, refreshes both lists, and (for
  // non-enabled plugins) opens the InstallConfirmModal.
  onInstall: (name: string, version: string) => Promise<void>;
};

// BrowseRegistry renders the publicly-available plugins as a list of
// cards inside the Browse tab. The bytes never pass through the
// browser: the Owncast host downloads from the registry, verifies
// SHA256, and runs the same Install path as a manual .ocpkg upload.
export const BrowseRegistry = ({
  installedVersions,
  registry,
  registryLoading,
  onInstall,
}: BrowseRegistryProps) => {
  const { t } = useTranslation();
  // Names whose install POST is in flight; per-card spinner state so
  // multiple installs don't lock the whole list.
  const [installing, setInstalling] = useState<Set<string>>(new Set());

  const triggerInstall = async (plugin: RegistryPlugin) => {
    if (!plugin.latest) return;
    const { name } = plugin;
    setInstalling(prev => new Set(prev).add(name));
    try {
      await onInstall(name, plugin.latest.version);
    } finally {
      setInstalling(prev => {
        const next = new Set(prev);
        next.delete(name);
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

  if (registry.length === 0) {
    return <Empty description={t(Localization.Admin.Plugins.browseEmpty)} />;
  }

  return (
    <Space direction="vertical" size="middle" className={s.list}>
      {registry.map(plugin => {
        const installedVersion = installedVersions.get(plugin.name);
        const latestVersion = plugin.latest?.version;
        const isInstalled = installedVersion !== undefined;
        const hasUpdate =
          isInstalled && latestVersion !== undefined && installedVersion !== latestVersion;

        let actionButton: React.ReactNode;
        if (!plugin.latest) {
          actionButton = <Button disabled>{t(Localization.Admin.Plugins.browseInstall)}</Button>;
        } else if (hasUpdate) {
          actionButton = (
            <Button
              type="primary"
              loading={installing.has(plugin.name)}
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
              loading={installing.has(plugin.name)}
              onClick={() => triggerInstall(plugin)}
            >
              {t(Localization.Admin.Plugins.browseInstall)}
            </Button>
          );
        }

        return (
          <Card key={plugin.name} size="small">
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
                </div>
                {plugin.summary && <Paragraph className={s.summary}>{plugin.summary}</Paragraph>}
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
