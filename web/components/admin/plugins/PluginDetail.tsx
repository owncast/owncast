import React, { useEffect, useMemo, useState } from 'react';
import { Alert, Space, Spin, Table, Tabs, Tag, Typography } from 'antd';
import { useTranslation } from 'next-export-i18n';
import dynamic from 'next/dynamic';
import ReactMarkdown from 'react-markdown';
import { Plugin, PluginPermission } from '../../../interfaces/plugin';
import { Localization } from '../../../types/localization';
import { fetchText, pluginInstructionsUrl } from '../../../utils/apis';
import { permissionDescriptionKey } from './permissionDescriptions';
import s from './PluginDetail.module.scss';

const { Title, Text, Paragraph } = Typography;

// Manifest `icon` fields are short semantic names. Map them to AntD icons
// so the gear/wrench/etc. a plugin author wrote shows up on the tab. Names
// without a mapping fall through to no icon (and we'll just show the
// title), which is what we'd otherwise have rendered anyway.
const SettingOutlined = dynamic(() => import('@ant-design/icons/SettingOutlined'), { ssr: false });
const ToolOutlined = dynamic(() => import('@ant-design/icons/ToolOutlined'), { ssr: false });
const UserOutlined = dynamic(() => import('@ant-design/icons/UserOutlined'), { ssr: false });
const LockOutlined = dynamic(() => import('@ant-design/icons/LockOutlined'), { ssr: false });
const InfoCircleOutlined = dynamic(() => import('@ant-design/icons/InfoCircleOutlined'), {
  ssr: false,
});
const AppstoreOutlined = dynamic(() => import('@ant-design/icons/AppstoreOutlined'), {
  ssr: false,
});
const FileTextOutlined = dynamic(() => import('@ant-design/icons/FileTextOutlined'), {
  ssr: false,
});
const BellOutlined = dynamic(() => import('@ant-design/icons/BellOutlined'), { ssr: false });

const pageIconForName = (name?: string): React.ReactNode => {
  if (!name) return null;
  switch (name.toLowerCase()) {
    case 'gear':
    case 'settings':
    case 'cog':
      return <SettingOutlined />;
    case 'wrench':
    case 'tool':
      return <ToolOutlined />;
    case 'user':
    case 'users':
      return <UserOutlined />;
    case 'lock':
    case 'security':
      return <LockOutlined />;
    case 'info':
    case 'about':
      return <InfoCircleOutlined />;
    case 'apps':
    case 'plugin':
      return <AppstoreOutlined />;
    case 'docs':
    case 'page':
      return <FileTextOutlined />;
    case 'bell':
    case 'notifications':
      return <BellOutlined />;
    default:
      return null;
  }
};

export type PluginDetailProps = {
  plugin: Plugin;
};

// pluginAdminUrl builds the iframe URL as a same-origin absolute path.
// Globs that end in /* are stripped to the literal prefix; the plugin's
// own router serves anything beneath. A trailing "/" is appended when
// the path looks like a directory (no file extension on the last
// segment) so the host doesn't 301-redirect to canonicalize, which
// Next.js dev's trailing-slash handling can otherwise compose into a
// redirect loop. Paths whose last segment has an extension (e.g.
// "/admin/page.html") are left alone so the host actually serves the
// file instead of looking up a same-named directory.
const pluginAdminUrl = (pluginSlug: string, path: string): string => {
  let normalized = path;
  if (normalized.endsWith('/*')) {
    normalized = normalized.slice(0, -2) || '/';
  } else if (normalized.endsWith('*')) {
    normalized = normalized.slice(0, -1);
  }
  if (!normalized.startsWith('/')) {
    normalized = `/${normalized}`;
  }
  const lastSegment = normalized.split('/').filter(Boolean).pop() ?? '';
  const looksLikeFile = lastSegment.includes('.');
  if (!looksLikeFile && !normalized.endsWith('/')) {
    normalized = `${normalized}/`;
  }
  return `/plugins/${encodeURIComponent(pluginSlug)}${normalized}`;
};

// PluginInstructions fetches the plugin's bundled INSTRUCTIONS.md and
// renders it as markdown. The endpoint serves raw markdown (admin-gated),
// so we read it as text and hand it to ReactMarkdown. Fetched on mount and
// whenever the slug changes (navigating between plugins).
const PluginInstructions = ({ slug }: { slug: string }) => {
  const { t } = useTranslation();
  const [markdown, setMarkdown] = useState<string | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    let cancelled = false;
    setMarkdown(null);
    setError(false);
    fetchText(pluginInstructionsUrl(slug))
      .then(text => {
        if (!cancelled) setMarkdown(text);
      })
      .catch(() => {
        if (!cancelled) setError(true);
      });
    return () => {
      cancelled = true;
    };
  }, [slug]);

  if (error) {
    return (
      <Alert type="error" showIcon message={t(Localization.Admin.Plugins.instructionsLoadError)} />
    );
  }
  if (markdown === null) {
    return <Spin />;
  }
  return (
    <div className={s.instructions}>
      <ReactMarkdown>{markdown}</ReactMarkdown>
    </div>
  );
};

// PluginDetail renders the per-plugin view: metadata header + a tab for
// each admin page declared in the plugin's manifest. Each tab body is an
// iframe pointing at the plugin's own HTTP namespace.
//
// Declared pages are de-duplicated by the iframe URL they resolve to —
// manifests commonly list both `/admin` and `/admin/*` so the host's
// auth gating covers the page and anything beneath it, but for the user
// those are one logical page, not two tabs.
export const PluginDetail = ({ plugin }: PluginDetailProps) => {
  const { t } = useTranslation();
  const uniquePages = useMemo(() => {
    const seen = new Set<string>();
    return (plugin.adminPages ?? [])
      .map(page => ({ ...page, url: pluginAdminUrl(plugin.slug, page.path) }))
      .filter(page => {
        if (seen.has(page.url)) return false;
        seen.add(page.url);
        return true;
      });
  }, [plugin.adminPages, plugin.slug]);

  const pendingSet = useMemo(
    () => new Set(plugin.pendingPermissions ?? []),
    [plugin.pendingPermissions],
  );

  const permissionRows = useMemo(
    () =>
      (plugin.permissions ?? []).map(perm => {
        const key = permissionDescriptionKey[perm];
        return {
          key: perm,
          permission: perm,
          description: key ? t(key) : '-',
          // network.fetch carries an extra dimension to the trust
          // decision: which hosts the plugin is allowed to reach.
          // Surface manifest.network.allowedHosts on this row so the
          // admin sees the host scope alongside the permission.
          allowedHosts:
            perm === PluginPermission.NetworkFetch ? (plugin.allowedHosts ?? []) : undefined,
          pending: pendingSet.has(perm),
        };
      }),
    [plugin.permissions, plugin.allowedHosts, pendingSet, t],
  );

  const permissionsTab = useMemo(
    () => ({
      key: '__permissions',
      label: t(Localization.Admin.Plugins.permissionsTab),
      children:
        permissionRows.length === 0 ? (
          <Alert
            type="info"
            showIcon
            message={t(Localization.Admin.Plugins.noPermissionsTitle)}
            description={t(Localization.Admin.Plugins.noPermissionsDescription)}
          />
        ) : (
          <Space direction="vertical" size={12} style={{ width: '100%' }}>
            {pendingSet.size > 0 && (
              <Alert
                type="warning"
                showIcon
                message={t(Localization.Admin.Plugins.statusPendingApproval)}
                description={t(Localization.Admin.Plugins.statusPendingApprovalTooltip)}
              />
            )}
            <Table
              dataSource={permissionRows}
              pagination={false}
              size="middle"
              columns={[
                {
                  title: t(Localization.Admin.Plugins.permissionColumnHeader),
                  dataIndex: 'permission',
                  key: 'permission',
                  width: 220,
                  render: (v: string, row: { pending: boolean; allowedHosts?: string[] }) =>
                    row.pending ? (
                      <Space size={6}>
                        <code>{v}</code>
                        <Tag color="orange">
                          {t(Localization.Admin.Plugins.statusPendingApproval)}
                        </Tag>
                      </Space>
                    ) : (
                      <code>{v}</code>
                    ),
                },
                {
                  title: t(Localization.Admin.Plugins.descriptionColumnHeader),
                  dataIndex: 'description',
                  key: 'description',
                  render: (v: string, row: { allowedHosts?: string[]; pending: boolean }) =>
                    row.allowedHosts && row.allowedHosts.length > 0 ? (
                      <Space direction="vertical" size={4}>
                        <span>{v}</span>
                        <span>
                          {t(Localization.Admin.Plugins.allowedHostsLabel)}{' '}
                          {row.allowedHosts.map((host, idx) => (
                            <React.Fragment key={host}>
                              {idx > 0 && ', '}
                              <code>{host}</code>
                            </React.Fragment>
                          ))}
                        </span>
                      </Space>
                    ) : (
                      v
                    ),
                },
              ]}
            />
          </Space>
        ),
    }),
    [permissionRows, pendingSet, t],
  );

  const pageTabs = useMemo(
    () =>
      uniquePages.map(page => {
        const icon = pageIconForName(page.icon);
        return {
          key: page.url,
          label: icon ? (
            <span>
              {icon} {page.title}
            </span>
          ) : (
            page.title
          ),
          children: (
            <iframe
              title={`${plugin.name} – ${page.title}`}
              src={page.url}
              sandbox="allow-same-origin allow-scripts allow-forms allow-popups"
              className={s.iframe}
            />
          ),
        };
      }),
    [uniquePages, plugin.name],
  );

  const instructionsTab = useMemo(
    () =>
      plugin.hasInstructions
        ? {
            key: '__instructions',
            label: (
              <span>
                <FileTextOutlined /> {t(Localization.Admin.Plugins.instructionsTab)}
              </span>
            ),
            children: <PluginInstructions slug={plugin.slug} />,
          }
        : null,
    [plugin.hasInstructions, plugin.slug, t],
  );

  const tabs = useMemo(() => {
    const base = [...pageTabs, permissionsTab];
    // Instructions sits as the second tab: after the first admin page (or
    // after Permissions when the plugin declares no admin pages), so the
    // primary page stays the landing tab.
    return instructionsTab ? [base[0], instructionsTab, ...base.slice(1)] : base;
  }, [instructionsTab, pageTabs, permissionsTab]);

  return (
    <div>
      <Space direction="vertical" size={16} className={s.container}>
        <Space direction="vertical" size={4}>
          <Title level={3} className={s.title}>
            {plugin.name}
            {plugin.version && (
              <Text type="secondary" className={s.version}>
                v{plugin.version}
              </Text>
            )}
          </Title>
          {plugin.description && <Paragraph>{plugin.description}</Paragraph>}
        </Space>

        {plugin.lastError && (
          <Alert
            type="error"
            showIcon
            message={t(Localization.Admin.Plugins.pluginErrorTitle)}
            description={plugin.lastError}
          />
        )}

        {/*
          Uncontrolled tabs: AntD picks defaultActiveKey on mount and
          manages selection state internally. The `key` prop ties the
          tab state to the plugin slug, so navigating between plugins
          (same component, new props) remounts the Tabs and resets the
          active tab to the new plugin's first one instead of carrying
          the stale selection across.
        */}
        <Tabs key={plugin.slug} defaultActiveKey={tabs[0]?.key} items={tabs} />
      </Space>
    </div>
  );
};
