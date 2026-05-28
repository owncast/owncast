import React, { ReactElement, useCallback, useEffect, useState } from 'react';
import { Alert, Button, message, Space, Typography, Upload } from 'antd';
import type { UploadProps } from 'antd';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { useTranslation } from 'next-export-i18n';
import { AdminLayout } from '../../components/layouts/AdminLayout';
import { fetchData, PLUGIN_UPLOAD, PLUGINS_LIST, pluginActionUrl } from '../../utils/apis';
import { Plugin } from '../../interfaces/plugin';
import { PluginsList } from '../../components/admin/plugins/PluginsList';
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

  useEffect(() => {
    loadPlugins();
  }, [loadPlugins]);

  const handleToggleEnabled = async (plugin: Plugin, enabled: boolean) => {
    setTogglingNames(prev => {
      const next = new Set(prev);
      next.add(plugin.name);
      return next;
    });
    try {
      await fetchData(pluginActionUrl(plugin.name, enabled ? 'enable' : 'disable'), {
        method: 'POST',
      });
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setTogglingNames(prev => {
        const next = new Set(prev);
        next.delete(plugin.name);
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
      next.add(plugin.name);
      return next;
    });
    try {
      await fetchData(pluginActionUrl(plugin.name, 'uninstall'), { method: 'POST' });
      message.success(t(Localization.Admin.Plugins.uninstallSuccess, { name: plugin.name }));
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setUninstallingNames(prev => {
        const next = new Set(prev);
        next.delete(plugin.name);
        return next;
      });
    }
  };

  const handleReload = async (plugin: Plugin) => {
    setReloadingNames(prev => {
      const next = new Set(prev);
      next.add(plugin.name);
      return next;
    });
    try {
      await fetchData(pluginActionUrl(plugin.name, 'reload'), { method: 'POST' });
      await loadPlugins();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setReloadingNames(prev => {
        const next = new Set(prev);
        next.delete(plugin.name);
        return next;
      });
    }
  };

  return (
    <div>
      <Space direction="horizontal" className={s.titleRow}>
        <Title>{t(Localization.Admin.Plugins.pageTitle)}</Title>
        <Upload {...uploadProps}>
          <Button type="primary" icon={<UploadOutlined />}>
            {t(Localization.Admin.Plugins.uploadButton)}
          </Button>
        </Upload>
        <Button icon={<ReloadOutlined />} onClick={loadPlugins}>
          {t(Localization.Admin.Plugins.refresh)}
        </Button>
      </Space>

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

      <PluginsList
        plugins={plugins}
        loading={loading}
        togglingNames={togglingNames}
        reloadingNames={reloadingNames}
        uninstallingNames={uninstallingNames}
        onToggleEnabled={handleToggleEnabled}
        onReload={handleReload}
        onUninstall={handleUninstall}
        onSelect={p =>
          router.push({ pathname: '/admin/plugins/configure', query: { name: p.name } })
        }
      />
    </div>
  );
};

Plugins.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default Plugins;
