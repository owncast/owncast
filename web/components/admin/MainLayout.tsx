import { FC, ReactNode, useContext, useEffect, useState } from 'react';
import Link from 'next/link';
import Head from 'next/head';
import { differenceInSeconds } from 'date-fns';
import { useRouter } from 'next/router';
import { Layout, Menu, Alert, Button, Space, Tooltip, Badge } from 'antd';

import { useTranslation } from 'next-export-i18n';
import classNames from 'classnames';
import dynamic from 'next/dynamic';
import { upgradeVersionAvailable } from '../../utils/apis';
import { PluginsContext } from '../../utils/plugins-context';
import { PluginIcon } from './plugins/PluginIcon';
import { Localization } from '../../types/localization';
import { parseSecondsToDurationString } from '../../utils/format';

import { OwncastLogo } from '../common/OwncastLogo/OwncastLogo';
import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';

import { TextFieldWithSubmit } from './TextFieldWithSubmit';
import { TEXTFIELD_PROPS_STREAM_TITLE } from '../../utils/config-constants';
import { ComposeFederatedPost } from './ComposeFederatedPost';
import { usePendingFeatureRequestCount } from '../../hooks/useFeatureRequests';
import { UpdateArgs } from '../../types/config-section';
import { FatalErrorStateModal } from '../modals/FatalErrorStateModal/FatalErrorStateModal';

// Lazy loaded components

const SettingOutlined = dynamic(() => import('@ant-design/icons/SettingOutlined'), {
  ssr: false,
}); // Lazy loaded components

const HomeOutlined = dynamic(() => import('@ant-design/icons/HomeOutlined'), {
  ssr: false,
});

const TeamOutlined = dynamic(() => import('@ant-design/icons/TeamOutlined'), {
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

const StarOutlined = dynamic(() => import('@ant-design/icons/StarOutlined'), {
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

  // Drives the badge on the Featured Streams sidebar item.
  const pendingFeatureRequestCount = usePendingFeatureRequestCount(federationEnabled);

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
  // submenu. Read from the shared PluginsContext so an install or
  // uninstall on the Plugins page updates the submenu live, without the
  // admin needing to refresh the whole page.
  const { plugins } = useContext(PluginsContext);

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
  const upgradeMessage = t(Localization.Admin.MainLayout.upgrade, {
    version: upgradeVersionString,
  });
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
  const statusMessage = online
    ? t(Localization.Admin.MainLayout.statusOnline, { duration: streamDurationString })
    : t(Localization.Admin.MainLayout.statusOffline);

  const statusIndicator = (
    <div className="online-status-indicator">
      <span className="status-label">{statusMessage}</span>
      <span className="status-icon">{statusIcon}</span>
    </div>
  );

  const integrationsMenu = [
    {
      label: <Link href="/admin/webhooks">{t(Localization.Admin.MainLayout.webhooks)}</Link>,
      key: '/admin/webhooks',
    },
    {
      label: (
        <Link href="/admin/access-tokens">{t(Localization.Admin.MainLayout.accessTokens)}</Link>
      ),
      key: '/admin/access-tokens',
    },
    {
      label: <Link href="/admin/actions">{t(Localization.Admin.MainLayout.externalActions)}</Link>,
      key: '/admin/actions',
    },
  ];

  const chatMenu = [
    {
      label: (
        <Link href="/admin/chat/messages">{t(Localization.Admin.MainLayout.chatMessages)}</Link>
      ),
      key: '/admin/chat/messages',
    },
    {
      label: <Link href="/admin/chat/emojis">{t(Localization.Admin.MainLayout.chatEmojis)}</Link>,
      key: '/admin/chat/emojis',
    },
  ];

  const utilitiesMenu = [
    {
      label: <Link href="/admin/hardware-info">{t(Localization.Admin.MainLayout.hardware)}</Link>,
      key: '/admin/hardware-info',
    },
    {
      label: (
        <Link href="/admin/stream-health">{t(Localization.Admin.MainLayout.streamHealth)}</Link>
      ),
      key: '/admin/stream-health',
    },
    {
      label: <Link href="/admin/logs">{t(Localization.Admin.MainLayout.logs)}</Link>,
      key: '/admin/logs',
    },
    ...(federationEnabled
      ? [
          {
            label: (
              <Link href="/admin/federation/actions">
                {t(Localization.Admin.MainLayout.socialActions)}
              </Link>
            ),
            key: '/admin/federation/actions',
          },
        ]
      : []),
  ];

  const configurationMenu = [
    {
      label: (
        <Link href="/admin/config/general">{t(Localization.Admin.MainLayout.configGeneral)}</Link>
      ),
      key: '/admin/config/general',
    },
    {
      label: (
        <Link href="/admin/config/server">{t(Localization.Admin.MainLayout.configServer)}</Link>
      ),
      key: '/admin/config/server',
    },
    {
      label: <Link href="/admin/config-video">{t(Localization.Admin.MainLayout.configVideo)}</Link>,
      key: '/admin/config-video',
    },
    {
      label: <Link href="/admin/config-chat">{t(Localization.Admin.MainLayout.configChat)}</Link>,
      key: '/admin/config-chat',
    },
    {
      label: (
        <Link href="/admin/config-federation">{t(Localization.Admin.MainLayout.configSocial)}</Link>
      ),
      key: '/admin/config-federation',
    },
    {
      label: (
        <Link href="/admin/config-notify">
          {t(Localization.Admin.MainLayout.configNotifications)}
        </Link>
      ),
      key: '/admin/config-notify',
    },
  ];

  const menuItems = [
    {
      label: <Link href="/admin">{t(Localization.Admin.MainLayout.home)}</Link>,
      icon: <HomeOutlined />,
      key: '/admin',
    },
    {
      label: <Link href="/admin/viewer-info">{t(Localization.Admin.MainLayout.viewers)}</Link>,
      icon: <LineChartOutlined />,
      key: '/admin/viewer-info',
    },
    {
      label: <Link href="/admin/users">{t(Localization.Admin.MainLayout.users)}</Link>,
      icon: <TeamOutlined />,
      key: '/admin/users',
    },
    ...(!chatDisabled
      ? [
          {
            label: <span>{t(Localization.Admin.MainLayout.chat)}</span>,
            icon: <MessageOutlined />,
            children: chatMenu,
            key: 'chat',
          },
        ]
      : []),
    ...(federationEnabled
      ? [
          {
            key: '/admin/federation/followers',
            label: (
              <Link href="/admin/federation/followers">
                {t(Localization.Admin.MainLayout.followers)}
              </Link>
            ),
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
        ]
      : []),
    ...(federationEnabled
      ? [
          {
            key: '/admin/config-featured',
            label: (
              <Link href="/admin/config-featured">
                {t(Localization.Admin.MainLayout.featuredStreams)}
                {pendingFeatureRequestCount > 0 && (
                  <Badge
                    count={pendingFeatureRequestCount}
                    size="small"
                    style={{ marginInlineStart: 8 }}
                  />
                )}
              </Link>
            ),
            icon: <StarOutlined />,
          },
        ]
      : []),
    {
      key: 'configuration',
      label: t(Localization.Admin.MainLayout.configuration),
      icon: <SettingOutlined />,
      children: configurationMenu,
    },
    {
      key: 'utilities',
      label: t(Localization.Admin.MainLayout.utilities),
      icon: <ToolOutlined />,
      children: utilitiesMenu,
    },
    {
      key: 'integrations',
      label: t(Localization.Admin.MainLayout.integrations),
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
        // One entry per enabled plugin that declares at least one admin
        // page, so the admin can jump straight to a plugin's config
        // without going through the overview + Configure button. Disabled
        // plugins are omitted — their admin pages don't serve. URL is
        // a static route plus an id query param (the plugin's slug)
        // because the admin UI is statically exported and can't
        // enumerate plugin identifiers at build time. Sidebar labels use
        // the human-readable display name.
        ...plugins
          .filter(p => p.enabled && Object.keys(p.adminPages ?? {}).length > 0)
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
    ...(upgradeVersion
      ? [
          {
            type: 'divider' as const,
            key: 'upgrade-divider',
          },
        ]
      : []),
    ...(upgradeVersion
      ? [
          {
            key: '/admin/upgrade',
            label: (
              <Link href="/admin/upgrade">
                <strong>{upgradeMessage}</strong>
              </Link>
            ),
            icon: <DownloadOutlined />,
          },
        ]
      : []),
    {
      key: '/admin/help',
      label: <Link href="/admin/help">{t(Localization.Admin.MainLayout.help)}</Link>,
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
        <title>{t(Localization.Admin.MainLayout.title)}</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {serverError?.type === 'OWNCAST_SERVICE_UNREACHABLE' && (
        <FatalErrorStateModal
          title={t(Localization.Admin.MainLayout.serverUnreachable)}
          message={serverError.msg}
        />
      )}

      <Sider width={240} className="side-nav">
        <h1 className="owncast-title">
          <span className="logo-container">
            <OwncastLogo variant="simple" />
          </span>
          <span className="title-label">{t(Localization.Admin.MainLayout.title)}</span>
        </h1>
        <Menu
          mode="inline"
          theme="dark"
          className="menu-container"
          items={menuItems}
          selectedKeys={[route || '/admin']}
          openKeys={openKeys}
          onOpenChange={onOpenChange}
        />
      </Sider>

      <Layout className="layout-main">
        <Header className="layout-header">
          <Space orientation="horizontal">
            <Tooltip title={t(Localization.Admin.MainLayout.composePostTooltip)}>
              <Button
                type="link"
                icon={<EditOutlined />}
                size="small"
                onClick={handleCreatePostButtonPressed}
                style={{ display: federationEnabled ? 'block' : 'none', margin: '10px' }}
              >
                {t(Localization.Admin.MainLayout.composePost)}
              </Button>
            </Tooltip>
          </Space>
          <div className="global-stream-title-container">
            <TextFieldWithSubmit
              fieldName="streamTitle"
              {...TEXTFIELD_PROPS_STREAM_TITLE}
              placeholder={t(Localization.Admin.MainLayout.streamTitlePlaceholder)}
              value={currentStreamTitle}
              initialValue={instanceDetails.streamTitle}
              onChange={handleStreamTitleChanged}
            />
          </div>
          <Space orientation="horizontal">{statusIndicator}</Space>
        </Header>

        {headerAlertMessage}

        <Content className="main-content-container">{children}</Content>

        <Footer className="footer-container">
          <a href="https://owncast.online/?source=admin" target="_blank" rel="noopener noreferrer">
            {t(Localization.Admin.MainLayout.aboutOwncast, { version: versionNumber })}
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
