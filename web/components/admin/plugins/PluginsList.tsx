import React from 'react';
import { Alert, Button, Popconfirm, Space, Switch, Table, Tag, Tooltip, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import dynamic from 'next/dynamic';
import { Plugin } from '../../../interfaces/plugin';
import { Localization } from '../../../types/localization';
import { permissionDescriptionKey, permissionNameKey } from './permissionDescriptions';
import { PluginIcon } from './PluginIcon';
import s from './PluginsList.module.scss';

const { Text } = Typography;

// Lazy-loaded icons match the pattern used elsewhere in admin/.
const ReloadOutlined = dynamic(() => import('@ant-design/icons/ReloadOutlined'), { ssr: false });
const WarningFilled = dynamic(() => import('@ant-design/icons/WarningFilled'), { ssr: false });
const AppstoreOutlined = dynamic(() => import('@ant-design/icons/AppstoreOutlined'), {
  ssr: false,
});
const DeleteOutlined = dynamic(() => import('@ant-design/icons/DeleteOutlined'), { ssr: false });

export type PluginsListProps = {
  plugins: Plugin[];
  loading: boolean;
  // Plugin names whose enable/disable request is in flight — drives the Switch
  // spinner. Per-plugin (not a single global flag) so two concurrent toggles
  // don't lock the whole table.
  togglingNames: Set<string>;
  // Same idea for the Reload button.
  reloadingNames: Set<string>;
  // And for the Uninstall button.
  uninstallingNames: Set<string>;
  onToggleEnabled: (plugin: Plugin, enabled: boolean) => Promise<void> | void;
  onReload: (plugin: Plugin) => Promise<void> | void;
  onUninstall: (plugin: Plugin) => Promise<void> | void;
  onSelect: (plugin: Plugin) => void;
};

// PluginsList renders one row per discovered plugin with metadata, an
// enable/disable toggle, a reload button, and a click-to-detail action.
// Plugins that declare manifest.admin.pages expose a "View" button that takes
// the admin to the per-plugin detail with embedded iframe(s).
export const PluginsList = ({
  plugins,
  loading,
  togglingNames,
  reloadingNames,
  uninstallingNames,
  onToggleEnabled,
  onReload,
  onUninstall,
  onSelect,
}: PluginsListProps) => {
  const { t } = useTranslation();
  const columns = [
    {
      title: t(Localization.Admin.Plugins.pluginColumn),
      key: 'name',
      render: (_: unknown, plugin: Plugin) => (
        <div className={s.pluginCell}>
          <PluginIcon plugin={plugin} size="list" />
          <Space direction="vertical" size={2}>
            <Text strong>{plugin.name}</Text>
            {plugin.version && (
              <Text type="secondary" className={s.secondaryText}>
                v{plugin.version}
              </Text>
            )}
            {plugin.description && (
              <Text type="secondary" className={s.secondaryText}>
                {plugin.description}
              </Text>
            )}
            {plugin.lastError && (
              <Alert
                type="error"
                showIcon
                icon={<WarningFilled />}
                message={t(Localization.Admin.Plugins.pluginFailedToLoad)}
                description={plugin.lastError}
                className={s.errorAlert}
              />
            )}
          </Space>
        </div>
      ),
    },
    {
      title: t(Localization.Admin.Plugins.permissionsColumn),
      key: 'permissions',
      render: (_: unknown, plugin: Plugin) => {
        if (!plugin.permissions || plugin.permissions.length === 0) {
          return <Text type="secondary">{t(Localization.Admin.Plugins.none)}</Text>;
        }
        return (
          <Space size={[4, 4]} wrap>
            {plugin.permissions.map(p => {
              const nameKey = permissionNameKey[p];
              const descKey = permissionDescriptionKey[p];
              const label = nameKey ? t(nameKey) : p;
              const description = descKey ? t(descKey) : p;
              return (
                <Tooltip key={p} title={description}>
                  <Tag>{label}</Tag>
                </Tooltip>
              );
            })}
          </Space>
        );
      },
    },
    {
      title: t(Localization.Admin.Plugins.statusColumn),
      key: 'status',
      render: (_: unknown, plugin: Plugin) => {
        if ((plugin.pendingPermissions?.length ?? 0) > 0) {
          return (
            <Tooltip title={t(Localization.Admin.Plugins.statusPendingApprovalTooltip)}>
              <Tag color="orange">{t(Localization.Admin.Plugins.statusPendingApproval)}</Tag>
            </Tooltip>
          );
        }
        if (plugin.lastError) {
          return <Tag color="red">{t(Localization.Admin.Plugins.statusError)}</Tag>;
        }
        if (plugin.autoDisabled) {
          return (
            <Tooltip title={t(Localization.Admin.Plugins.statusAutoDisabledTooltip)}>
              <Tag color="red">{t(Localization.Admin.Plugins.statusAutoDisabled)}</Tag>
            </Tooltip>
          );
        }
        if (plugin.loaded) {
          return <Tag color="green">{t(Localization.Admin.Plugins.statusRunning)}</Tag>;
        }
        if (plugin.enabled) {
          return <Tag color="orange">{t(Localization.Admin.Plugins.statusEnabledNotLoaded)}</Tag>;
        }
        return <Tag>{t(Localization.Admin.Plugins.statusDisabled)}</Tag>;
      },
    },
    {
      title: t(Localization.Admin.Plugins.enabledColumn),
      key: 'enabled',
      render: (_: unknown, plugin: Plugin) => {
        if ((plugin.pendingPermissions?.length ?? 0) > 0) {
          return (
            <Tooltip title={t(Localization.Admin.Plugins.approveTooltip)}>
              <Button
                type="primary"
                loading={togglingNames.has(plugin.name)}
                onClick={() => onToggleEnabled(plugin, true)}
              >
                {t(Localization.Admin.Plugins.approveButton)}
              </Button>
            </Tooltip>
          );
        }
        return (
          <Switch
            checked={plugin.enabled}
            loading={togglingNames.has(plugin.name)}
            onChange={checked => onToggleEnabled(plugin, checked)}
            aria-label={t(Localization.Admin.Plugins.toggleAria, { name: plugin.name })}
          />
        );
      },
    },
    {
      title: '',
      key: 'actions',
      render: (_: unknown, plugin: Plugin) => {
        const hasAdminPages = (plugin.adminPages?.length ?? 0) > 0;
        return (
          <Space>
            <Tooltip title={t(Localization.Admin.Plugins.reloadTooltip)}>
              <Button
                icon={<ReloadOutlined />}
                loading={reloadingNames.has(plugin.name)}
                onClick={() => onReload(plugin)}
                disabled={!plugin.enabled}
              />
            </Tooltip>
            {hasAdminPages && (
              <Tooltip title={t(Localization.Admin.Plugins.openPluginAdmin)}>
                <Button icon={<AppstoreOutlined />} onClick={() => onSelect(plugin)}>
                  {t(Localization.Admin.Plugins.configure)}
                </Button>
              </Tooltip>
            )}
            <Popconfirm
              title={
                <div className={s.uninstallPrompt}>
                  <div className={s.uninstallPromptTitle}>
                    {t(Localization.Admin.Plugins.uninstallConfirmTitle)}
                  </div>
                  <div className={s.uninstallPromptDescription}>
                    {t(Localization.Admin.Plugins.uninstallConfirmDescription, {
                      name: plugin.name,
                    })}
                  </div>
                </div>
              }
              okText={t(Localization.Admin.Plugins.uninstallConfirmOk)}
              cancelText={t(Localization.Admin.Plugins.uninstallConfirmCancel)}
              okButtonProps={{ danger: true }}
              onConfirm={() => onUninstall(plugin)}
            >
              <Tooltip title={t(Localization.Admin.Plugins.uninstallTooltip)}>
                <Button
                  danger
                  icon={<DeleteOutlined />}
                  loading={uninstallingNames.has(plugin.name)}
                  aria-label={t(Localization.Admin.Plugins.uninstallAria, { name: plugin.name })}
                />
              </Tooltip>
            </Popconfirm>
          </Space>
        );
      },
    },
  ];

  return (
    <Table
      rowKey={p => p.name}
      columns={columns}
      dataSource={plugins}
      loading={loading}
      pagination={false}
    />
  );
};
