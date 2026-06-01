import React, { ReactElement, useCallback, useEffect, useMemo, useState } from 'react';
import { Alert, Button, message, Space, Tabs, Typography, Upload } from 'antd';
import type { UploadProps } from 'antd';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { useTranslation } from 'next-export-i18n';
import { AdminLayout } from '../../components/layouts/AdminLayout';
import {
  fetchData,
  isPluginUpdateAvailable,
  PLUGIN_REGISTRY_INSTALL,
  PLUGIN_REGISTRY_LIST,
  PLUGIN_UPLOAD,
  PLUGINS_LIST,
  pluginActionUrl,
} from '../../utils/apis';
import { Plugin } from '../../interfaces/plugin';
import { PluginsList } from '../../components/admin/plugins/PluginsList';
import { BrowseRegistry, RegistryPlugin } from '../../components/admin/plugins/BrowseRegistry';
import { InstallConfirmModal } from '../../components/admin/plugins/InstallConfirmModal';
import { Localization } from '../../types/localization';
import s from './plugins.module.scss';

const { Title, Paragraph } = Typography;

const ReloadOutlined = dynamic(() => import('@ant-design/icons/ReloadOutlined'), { ssr: false });
const UploadOutlined = dynamic(() => import('@ant-design/icons/UploadOutlined'), { ssr: false });

const Plugins = () => {
  const { t } = useTranslation();
  const router = useRouter();
  const [plugins, setPlugins] = useState<Plugin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  // Names of plugins whose enable/disable POST is in flight. Drives the per-row
  // spinner on the Switch so the admin sees the action is happening (the
  // enable/disable round-trip includes a wasm load and can take a moment).
  const [togglingNames, setTogglingNames] = useState<Set<string>>(new Set());
  // Same idea for the Reload button — independent from togglingNames because
  // the two actions are independent (and visually distinct).
  const [reloadingNames, setReloadingNames] = useState<Set<string>>(new Set());
  // And for the Uninstall button.
  const [uninstallingNames, setUninstallingNames] = useState<Set<string>>(new Set());
  // Tracks which tab is active so the tab-bar extras (Upload +
  // Refresh) only show when they're contextually relevant to the
  // Installed view.
  const [activeTab, setActiveTab] = useState('installed');
  // The just-installed plugin awaiting an enable/cancel decision in
  // the InstallConfirmModal. Null when no modal is open.
  const [pendingEnable, setPendingEnable] = useState<Plugin | null>(null);
  // Registry list, used both by the Browse tab and by the Installed
  // tab's "update available" tags. Fetched at the page level so a
  // single source of truth feeds both views.
  const [registryPlugins, setRegistryPlugins] = useState<RegistryPlugin[]>([]);
  const [registryLoading, setRegistryLoading] = useState(true);
  // registryError is set when the registry fetch fails (network error, host
  // returned 5xx, etc.) so the Browse tab can distinguish "the catalog is
  // unreachable" from "the catalog is reachable but empty". `null` covers
  // both the loading and the empty-but-reachable cases.
  const [registryError, setRegistryError] = useState<string | null>(null);

  const loadPlugins = useCallback(async () => {
    try {
      const result = await fetchData(PLUGINS_LIST);
      setPlugins(Array.isArray(result) ? result : []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  // Registry fetch is independent of the installed-list fetch: a
  // misconfigured or unreachable registry shouldn't block the
  // Installed tab. On failure we set registryError so the Browse tab
  // renders a "catalog unavailable" message with a retry button,
  // distinct from the legitimately-empty case.
  const loadRegistry = useCallback(async () => {
    setRegistryLoading(true);
    try {
      const result = await fetchData(PLUGIN_REGISTRY_LIST);
      setRegistryPlugins(Array.isArray(result) ? result : []);
      setRegistryError(null);
    } catch (e) {
      setRegistryPlugins([]);
      setRegistryError(e instanceof Error ? e.message : String(e));
    } finally {
      setRegistryLoading(false);
    }
  }, []);

  useEffect(() => {
    loadPlugins();
    loadRegistry();
  }, [loadPlugins, loadRegistry]);

  // Map of installed plugin slug -> currently-installed version.
  // Browse cards use this to decide whether to render "Install" /
  // "Installed" / "Update". Keys on slug (the canonical identifier),
  // not display name, so two plugins with the same display string
  // don't collide.
  const installedVersions = useMemo(
    () =>
      new Map<string, string>(
        plugins.filter(p => Boolean(p.version)).map(p => [p.slug, p.version as string]),
      ),
    [plugins],
  );

  // Map of installed plugin slug -> newer version available in the
  // registry. Empty when the installed and registry versions match
  // (or the plugin isn't in the registry at all). PluginsList renders
  // the entries as "update available" tags in the version line.
  const availableUpdates = useMemo(
    () =>
      new Map<string, string>(
        plugins.flatMap(p => {
          const reg = registryPlugins.find(r => r.slug === p.slug);
          const latest = reg?.latest?.version;
          if (latest && isPluginUpdateAvailable(p.version, latest)) {
            return [[p.slug, latest] as [string, string]];
          }
          return [];
        }),
      ),
    [plugins, registryPlugins],
  );

  const handleToggleEnabled = async (plugin: Plugin, enabled: boolean) => {
    setTogglingNames(prev => {
      const next = new Set(prev);
      next.add(plugin.slug);
      return next;
    });
    try {
      await fetchData(pluginActionUrl(plugin.slug, enabled ? 'enable' : 'disable'), {
        method: 'POST',
      });
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setTogglingNames(prev => {
        const next = new Set(prev);
        next.delete(plugin.slug);
        return next;
      });
    }
  };

  // AntD's Upload calls customRequest with one file at a time. The request
  // is multipart with a single 'plugin' field, matching the server-side
  // form field name. We avoid AntD's automatic file-list display because
  // the result we care about is the refreshed plugin list, not the
  // intermediate upload state.
  const uploadProps: UploadProps = {
    name: 'plugin',
    accept: '.ocpkg',
    showUploadList: false,
    multiple: false,
    customRequest: async ({ file, onSuccess, onError }) => {
      const blob = file as Blob;
      const form = new FormData();
      form.append('plugin', blob, (file as File).name);
      try {
        const res = await fetch(PLUGIN_UPLOAD, { method: 'POST', body: form });
        if (!res.ok) {
          const body = await res.text();
          let detail = body;
          try {
            detail = JSON.parse(body).error ?? body;
          } catch {
            /* not JSON, use raw body */
          }
          throw new Error(detail || `upload failed: ${res.status}`);
        }
        const entry = (await res.json()) as Plugin;
        message.success(t(Localization.Admin.Plugins.uploadSuccess, { name: entry.name }));
        await loadPlugins();
        onSuccess?.(entry);
      } catch (e) {
        const msg = e instanceof Error ? e.message : String(e);
        setError(msg);
        onError?.(e as Error);
      }
    },
  };

  const handleUninstall = async (plugin: Plugin) => {
    setUninstallingNames(prev => {
      const next = new Set(prev);
      next.add(plugin.slug);
      return next;
    });
    try {
      await fetchData(pluginActionUrl(plugin.slug, 'uninstall'), { method: 'POST' });
      message.success(t(Localization.Admin.Plugins.uninstallSuccess, { name: plugin.name }));
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setUninstallingNames(prev => {
        const next = new Set(prev);
        next.delete(plugin.slug);
        return next;
      });
    }
  };

  const handleReload = async (plugin: Plugin) => {
    setReloadingNames(prev => {
      const next = new Set(prev);
      next.add(plugin.slug);
      return next;
    });
    try {
      await fetchData(pluginActionUrl(plugin.slug, 'reload'), { method: 'POST' });
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setReloadingNames(prev => {
        const next = new Set(prev);
        next.delete(plugin.slug);
        return next;
      });
    }
  };

  // Called when the admin clicks "Enable plugin" in the post-install
  // confirmation modal. Reuses the same enable endpoint the Switch
  // hits, so the re-approval and load-failure paths are identical.
  const handleEnableAfterInstall = async () => {
    const plugin = pendingEnable;
    setPendingEnable(null);
    if (!plugin) return;
    try {
      await fetchData(pluginActionUrl(plugin.slug, 'enable'), { method: 'POST' });
      message.success(t(Localization.Admin.Plugins.installEnabledSuccess, { name: plugin.name }));
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  // Common path for "fetch the registry's bytes for (name, version)
  // and run Manager.Install". Shared between the Browse tab's Update
  // button and the Installed tab's "update available" tag so both
  // surface the same post-install confirmation flow.
  const installFromRegistry = useCallback(
    async (slug: string, version: string) => {
      try {
        const entry = (await fetchData(PLUGIN_REGISTRY_INSTALL, {
          method: 'POST',
          data: { slug, version },
        })) as Plugin;
        setError(null);
        message.success(t(Localization.Admin.Plugins.uploadSuccess, { name: entry.name }));
        await loadPlugins();
        await loadRegistry();
        // Only prompt to enable for plugins that aren't already
        // running. Updates of enabled plugins keep their enabled state
        // unless the new manifest expanded permissions (in which case
        // the host marks them auto-disabled and the modal shows the
        // new perms).
        if (!entry.enabled) {
          setPendingEnable(entry);
        }
      } catch (e) {
        setError(e instanceof Error ? e.message : String(e));
      }
    },
    [loadPlugins, loadRegistry, t],
  );

  const installedTab = (
    <PluginsList
      plugins={plugins}
      loading={loading}
      togglingNames={togglingNames}
      reloadingNames={reloadingNames}
      uninstallingNames={uninstallingNames}
      availableUpdates={availableUpdates}
      onToggleEnabled={handleToggleEnabled}
      onReload={handleReload}
      onUninstall={handleUninstall}
      onUpdate={(plugin, version) => installFromRegistry(plugin.slug, version)}
      onSelect={p => router.push({ pathname: '/admin/plugins/configure', query: { id: p.slug } })}
    />
  );

  const browseTab = (
    <BrowseRegistry
      installedVersions={installedVersions}
      registry={registryPlugins}
      registryLoading={registryLoading}
      registryError={registryError}
      onInstall={installFromRegistry}
      onRetry={loadRegistry}
    />
  );

  // Tab-bar actions for the Installed view: Upload + Refresh. Lives
  // in the AntD Tabs `tabBarExtraContent` slot so the buttons sit
  // inline with the tab labels on the right edge. We hide them on
  // the Browse tab since neither action applies there.
  const installedActions = (
    <Space>
      <Upload {...uploadProps}>
        <Button icon={<UploadOutlined />}>{t(Localization.Admin.Plugins.uploadButton)}</Button>
      </Upload>
      <Button icon={<ReloadOutlined />} onClick={loadPlugins}>
        {t(Localization.Admin.Plugins.refresh)}
      </Button>
    </Space>
  );

  return (
    <div>
      <Title>{t(Localization.Admin.Plugins.pageTitle)}</Title>
      <Paragraph>{t(Localization.Admin.Plugins.pageDescription)}</Paragraph>

      {error && (
        <Alert
          type="error"
          showIcon
          closable
          message={t(Localization.Admin.Plugins.errorTitle)}
          description={error}
          onClose={() => setError(null)}
          className={s.errorAlert}
        />
      )}

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        tabBarExtraContent={activeTab === 'installed' ? installedActions : undefined}
        items={[
          {
            key: 'installed',
            label: t(Localization.Admin.Plugins.tabInstalled),
            children: installedTab,
          },
          { key: 'browse', label: t(Localization.Admin.Plugins.tabBrowse), children: browseTab },
        ]}
      />

      <InstallConfirmModal
        plugin={pendingEnable}
        onCancel={() => setPendingEnable(null)}
        onEnable={handleEnableAfterInstall}
      />
    </div>
  );
};

Plugins.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default Plugins;
