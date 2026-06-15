import React, { FC, ReactNode, useContext, useEffect, useState } from 'react';
import Link from 'next/link';
import Head from 'next/head';
import { differenceInSeconds } from 'date-fns';
import { useRouter } from 'next/router';
import { Layout, Menu, Alert, Button, Space, Tooltip } from 'antd';

import { useTranslation } from 'next-export-i18n';
import classNames from 'classnames';
import dynamic from 'next/dynamic';
import { fetchData, PLUGINS_LIST, upgradeVersionAvailable } from '../../utils/apis';
import { Plugin } from '../../interfaces/plugin';
import { PluginIcon } from './plugins/PluginIcon';
import { Localization } from '../../types/localization';
import { parseSecondsToDurationString } from '../../utils/format';

import { OwncastLogo } from '../common/OwncastLogo/OwncastLogo';
import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';

import { TextFieldWithSubmit } from './TextFieldWithSubmit';
import { TEXTFIELD_PROPS_STREAM_TITLE } from '../../utils/config-constants';
import { ComposeFederatedPost } from './ComposeFederatedPost';
import { UpdateArgs } from '../../types/config-section';
import { FatalErrorStateModal } from '../modals/FatalErrorStateModal/FatalErrorStateModal';

// Lazy loaded components

const SettingOutlined = dynamic(() => import('@ant-design/icons/SettingOutlined'), {
  ssr: false,
}); // Lazy loaded components

const HomeOutlined = dynamic(() => import('@ant-design/icons/HomeOutlined'), {
  ssr: false,
});

const LineChartOutlined = dynamic(() => import('@ant-design/icons/LineChartOutlined'), {
  ssr: false,
});

const ToolOutlined = dynamic(() => import('@ant-design/icons/ToolOutlined'), {
  ssr: false,
});

const PlayCircleFilled = dynamic(() => import('@ant-design/icons/PlayCircleFilled'), {
  ssr: false,
});

const MinusSquareFilled = dynamic(() => import('@ant-design/icons/MinusSquareFilled'), {
  ssr: false,
});

const QuestionCircleOutlined = dynamic(() => import('@ant-design/icons/QuestionCircleOutlined'), {
  ssr: false,
});

const MessageOutlined = dynamic(() => import('@ant-design/icons/MessageOutlined'), {
  ssr: false,
});

const ExperimentOutlined = dynamic(() => import('@ant-design/icons/ExperimentOutlined'), {
  ssr: false,
});

const AppstoreOutlined = dynamic(() => import('@ant-design/icons/AppstoreOutlined'), {
  ssr: false,
});

const EditOutlined = dynamic(() => import('@ant-design/icons/EditOutlined'), {
  ssr: false,
});

const DownloadOutlined = dynamic(() => import('@ant-design/icons/DownloadOutlined'), {
  ssr: false,
});

const FediverseOutlined = dynamic(() => import('../../assets/images/icons/fediverse.svg'), {
  ssr: false,
});

export type MainLayoutProps = {
  children: ReactNode;
};

export const MainLayout: FC<MainLayoutProps> = ({ children }) => {
  const { t } = useTranslation();
  const context = useContext(ServerStatusContext);
  const { serverConfig, online, broadcaster, versionNumber, error: serverError } = context || {};
  const { instanceDetails, chatDisabled, federation } = serverConfig;
  const { enabled: federationEnabled } = federation;

  const [currentStreamTitle, setCurrentStreamTitle] = useState('');
  const [postModalDisplayed, setPostModalDisplayed] = useState(false);

  const alertMessage = useContext(AlertMessageContext);

  const router = useRouter();
  const { route } = router || {};

  const { Header, Footer, Content, Sider } = Layout;

  const [upgradeVersion, setUpgradeVersion] = useState('');
  const checkForUpgrade = async () => {
    if (versionNumber === '0.0.0') {
      return;
    }
    try {
      const result = await upgradeVersionAvailable(versionNumber);
      setUpgradeVersion(result);
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    checkForUpgrade();
  }, [versionNumber]);

  // Plugins with declared admin pages drive the sidebar's Plugins
  // submenu. Fetched once on mount; the admin would refresh the page to
  // see the submenu reflect newly enabled/disabled plugins.
  const [plugins, setPlugins] = useState<Plugin[]>([]);
  useEffect(() => {
    let cancelled = false;
    fetchData(PLUGINS_LIST)
      .then(result => {
        if (cancelled) return;
        setPlugins(Array.isArray(result) ? result : []);
      })
      .catch(() => {
        // Sidebar still works without the submenu — the top-level Plugins
        // link to the overview page remains. Silent on error.
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    setCurrentStreamTitle(instanceDetails.streamTitle);
  }, [instanceDetails]);

  const handleStreamTitleChanged = ({ value }: UpdateArgs) => {
    setCurrentStreamTitle(value);
  };

  const handleCreatePostButtonPressed = () => {
    setPostModalDisplayed(true);
  };

  const appClass = classNames({
    'app-container': true,
    online,
  });

  const upgradeVersionString = `${upgradeVersion}` || '';
  const upgradeMessage = `Upgrade to v${upgradeVersionString}`;
  const openMenuItems = upgradeVersion ? ['utilities-menu'] : [];

  const clearAlertMessage = () => {
    alertMessage.setMessage(null);
  };

  const headerAlertMessage = alertMessage.message ? (
    <Alert message={alertMessage.message} afterClose={clearAlertMessage} banner closable />
  ) : null;

  // status indicator items
  const streamDurationString = broadcaster
    ? parseSecondsToDurationString(differenceInSeconds(new Date(), new Date(broadcaster.time)))
    : '';

  const statusIcon = online ? <PlayCircleFilled /> : <MinusSquareFilled />;
  const statusMessage = online ? `Online ${streamDurationString}` : 'Offline';

  const statusIndicator = (
    <div className="online-status-indicator">
      <span className="status-label">{statusMessage}</span>
      <span className="status-icon">{statusIcon}</span>
    </div>
  );

  const integrationsMenu = [
    {
      label: <Link href="/admin/webhooks">Webhooks</Link>,
      key: '/admin/webhooks',
    },
    {
      label: <Link href="/admin/access-tokens">Access Tokens</Link>,
      key: '/admin/access-tokens',
    },
    {
      label: <Link href="/admin/actions">External Actions</Link>,
      key: '/admin/actions',
    },
  ];

  const chatMenu = [
    {
      label: <Link href="/admin/chat/messages">Messages</Link>,
      key: '/admin/chat/messages',
    },
    {
      label: <Link href="/admin/chat/users">Users</Link>,
      key: '/admin/chat/users',
    },
    {
      label: <Link href="/admin/chat/emojis">Emojis</Link>,
      key: '/admin/chat/emojis',
    },
  ];

  const utilitiesMenu = [
    {
      label: <Link href="/admin/hardware-info">Hardware</Link>,
      key: '/admin/hardware-info',
    },
    {
      label: <Link href="/admin/stream-health">Stream Health</Link>,
      key: '/admin/stream-health',
    },
    {
      label: <Link href="/admin/logs">Logs</Link>,
      key: '/admin/logs',
    },
    federationEnabled && {
      label: <Link href="/admin/federation/actions">Social Actions</Link>,
      key: '/admin/federation/actions',
    },
  ];

  const configurationMenu = [
    {
      label: <Link href="/admin/config/general">General</Link>,
      key: '/admin/config/general',
    },
    {
      label: <Link href="/admin/config/server">Server Setup</Link>,
      key: '/admin/config/server',
    },
    {
      label: <Link href="/admin/config-video">Video</Link>,
      key: '/admin/config-video',
    },
    {
      label: <Link href="/admin/config-chat">Chat</Link>,
      key: '/admin/config-chat',
    },
    {
      label: <Link href="/admin/config-federation">Social</Link>,
      key: '/admin/config-federation',
    },
    {
      label: <Link href="/admin/config-featured">Featured</Link>,
      key: '/admin/config-featured',
    },
    {
      label: <Link href="/admin/config-notify">Notifications</Link>,
      key: '/admin/config-notify',
    },
  ];

  const menuItems = [
    { label: <Link href="/admin">Home</Link>, icon: <HomeOutlined />, key: '/admin' },
    {
      label: <Link href="/admin/viewer-info">Viewers</Link>,
      icon: <LineChartOutlined />,
      key: '/admin/viewer-info',
    },
    !chatDisabled && {
      label: <span>Chat &amp; Users</span>,
      icon: <MessageOutlined />,
      children: chatMenu,
      key: 'chat-and-users',
    },
    federationEnabled && {
      key: '/admin/federation/followers',
      label: <Link href="/admin/federation/followers">Followers</Link>,
      icon: (
        <span
          role="img"
          aria-label="message"
          className="anticon anticon-message ant-menu-item-icon"
        >
          {/* Wrapping the icon in span for consistency with other icons used
            directly from antd */}
          <FediverseOutlined />
        </span>
      ),
    },
    {
      key: 'configuration',
      label: 'Configuration',
      icon: <SettingOutlined />,
      children: configurationMenu,
    },
    {
      key: 'utilities',
      label: 'Utilities',
      icon: <ToolOutlined />,
      children: utilitiesMenu,
    },
    {
      key: 'integrations',
      label: 'Integrations',
      icon: <ExperimentOutlined />,
      children: integrationsMenu,
    },
    {
      key: 'plugins-menu',
      label: t(Localization.Admin.Plugins.sidebarTitle),
      icon: <AppstoreOutlined />,
      children: [
        {
          key: '/admin/plugins',
          label: <Link href="/admin/plugins">{t(Localization.Admin.Plugins.overview)}</Link>,
        },
        // One entry per loaded plugin that declares at least one admin
        // page, so the admin can jump straight to a plugin's config
        // without going through the overview + Configure button. URL is
        // a static route plus an id query param (the plugin's slug)
        // because the admin UI is statically exported and can't
        // enumerate plugin identifiers at build time. Sidebar labels use
        // the human-readable display name.
        ...plugins
          .filter(p => (p.adminPages?.length ?? 0) > 0)
          .map(p => ({
            key: `/admin/plugins/configure?id=${p.slug}`,
            label: (
              <Link href={{ pathname: '/admin/plugins/configure', query: { id: p.slug } }}>
                {p.name}
              </Link>
            ),
            icon: <PluginIcon plugin={p} size="sidebar" />,
          })),
      ],
    },
    upgradeVersion && {
      type: 'divider',
      key: 'upgrade-divider',
    },
    upgradeVersion && {
      key: '/admin/upgrade',
      label: (
        <Link href="/admin/upgrade">
          <strong>{upgradeMessage}</strong>
        </Link>
      ),
      icon: <DownloadOutlined />,
    },
    {
      key: '/admin/help',
      label: <Link href="/admin/help">Help</Link>,
      icon: <QuestionCircleOutlined />,
    },
  ];

  const [openKeys, setOpenKeys] = useState(openMenuItems);

  const onOpenChange = (keys: string[]) => {
    setOpenKeys(keys);
  };

  useEffect(() => {
    menuItems.forEach(item =>
      item?.children?.forEach(child => {
        if (child?.key === route) setOpenKeys([...openMenuItems, item.key]);
      }),
    );
  }, []);

  // The per-plugin configure page is a query-string route
  // (/admin/plugins/configure?id=...), so the literal-key match above
  // doesn't fire when navigating to a specific plugin. Open the plugins
  // submenu whenever the URL is anywhere in /admin/plugins.
  useEffect(() => {
    if (route && route.startsWith('/admin/plugins')) {
      setOpenKeys(prev => (prev.includes('plugins-menu') ? prev : [...prev, 'plugins-menu']));
    }
  }, [route]);

  return (
    <Layout id="admin-page" className={appClass}>
      <Head>
        <title>Owncast Admin</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {serverError?.type === 'OWNCAST_SERVICE_UNREACHABLE' && (
        <FatalErrorStateModal title="Server Unreachable" message={serverError.msg} />
      )}

      <Sider width={240} className="side-nav">
        <h1 className="owncast-title">
          <span className="logo-container">
            <OwncastLogo variant="simple" />
          </span>
          <span className="title-label">Owncast Admin</span>
        </h1>
        <Menu
          mode="inline"
          className="menu-container"
          items={menuItems}
          selectedKeys={[route || '/admin']}
          openKeys={openKeys}
          onOpenChange={onOpenChange}
        />
      </Sider>

      <Layout className="layout-main">
        <Header className="layout-header">
          <Space direction="horizontal">
            <Tooltip title="Compose post to your social followers">
              <Button
                type="link"
                icon={<EditOutlined />}
                size="small"
                onClick={handleCreatePostButtonPressed}
                style={{ display: federationEnabled ? 'block' : 'none', margin: '10px' }}
              >
                Compose Post
              </Button>
            </Tooltip>
          </Space>
          <div className="global-stream-title-container">
            <TextFieldWithSubmit
              fieldName="streamTitle"
              {...TEXTFIELD_PROPS_STREAM_TITLE}
              placeholder="What are you streaming now? (Stream title)"
              value={currentStreamTitle}
              initialValue={instanceDetails.streamTitle}
              onChange={handleStreamTitleChanged}
            />
          </div>
          <Space direction="horizontal">{statusIndicator}</Space>
        </Header>

        {headerAlertMessage}

        <Content className="main-content-container">{children}</Content>

        <Footer className="footer-container">
          <a href="https://owncast.online/?source=admin" target="_blank" rel="noopener noreferrer">
            About Owncast v{versionNumber}
          </a>
        </Footer>
      </Layout>

      <ComposeFederatedPost
        open={postModalDisplayed}
        handleClose={() => setPostModalDisplayed(false)}
      />
    </Layout>
  );
};
