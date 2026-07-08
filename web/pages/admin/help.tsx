import { Alert, Button, Card, Col, Divider, Input, Row, Space, Typography, message } from 'antd';
import Link from 'next/link';
import dynamic from 'next/dynamic';
import { useContext, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
import type { ReactElement } from 'react';
import type { VideoVariant } from '../../types/config-section';
import { Localization } from '../../types/localization';
import { ServerStatusContext } from '../../utils/server-status-context';
import { LOGS_WARN, fetchData, upgradeVersionAvailable } from '../../utils/apis';
import { AdminLayout } from '../../components/layouts/AdminLayout';

const { Title, Text, Paragraph } = Typography;

// Lazy loaded components

const BugTwoTone = dynamic(() => import('@ant-design/icons/BugTwoTone'), {
  ssr: false,
});

const MessageTwoTone = dynamic(() => import('@ant-design/icons/MessageTwoTone'), {
  ssr: false,
});

const ToolTwoTone = dynamic(() => import('@ant-design/icons/ToolTwoTone'), {
  ssr: false,
});

const CopyOutlined = dynamic(() => import('@ant-design/icons/CopyOutlined'), {
  ssr: false,
});

const LinkOutlined = dynamic(() => import('@ant-design/icons/LinkOutlined'), {
  ssr: false,
});

const SettingOutlined = dynamic(() => import('@ant-design/icons/SettingOutlined'), {
  ssr: false,
});

// The help page ships with the next release, alongside the new owncast.online
// docs site, so all documentation links use the new site's URL structure.
const DOCS_ROOT = 'https://owncast.online';

const docsUrl = (path: string): string => `${DOCS_ROOT}${path}?source=admin`;

type LogEntry = {
  time: string;
  level: string;
  message: string;
};

type HelpTask = {
  title: string;
  adminPage?: string;
  docsPath?: string;
};

const describeVariant = (variant: VideoVariant): string => {
  if (variant.videoPassthrough) {
    return 'video passthrough';
  }
  const size =
    variant.scaledWidth || variant.scaledHeight
      ? `${variant.scaledWidth || 'auto'}x${variant.scaledHeight || 'auto'} `
      : '';
  return `${size}${variant.videoBitrate}kbps ${variant.framerate}fps cpu:${variant.cpuUsageLevel}`;
};

export default function Help() {
  const { t } = useTranslation();
  const serverStatus = useContext(ServerStatusContext);
  const { online, versionNumber, overallPeakViewerCount, health, serverConfig } = serverStatus;
  const { videoSettings, videoCodec } = serverConfig;

  const [upgradeVersion, setUpgradeVersion] = useState('');
  const [releaseString, setReleaseString] = useState('');
  const [recentErrors, setRecentErrors] = useState<LogEntry[]>([]);

  useEffect(() => {
    if (!versionNumber || versionNumber === '0.0.0') {
      return;
    }
    upgradeVersionAvailable(versionNumber)
      .then(result => setUpgradeVersion(result || ''))
      .catch(() => {
        // Not being able to check for updates is not a problem worth surfacing
        // on the help page.
      });
  }, [versionNumber]);

  useEffect(() => {
    // The public config carries the full release string, including the build
    // platform, which the admin status API does not.
    fetch('/api/config')
      .then(response => response.json())
      .then(config => setReleaseString(config?.version || ''))
      .catch(() => {});

    fetchData(LOGS_WARN)
      .then(logs => {
        if (Array.isArray(logs)) {
          setRecentErrors(logs.slice(0, 20));
        }
      })
      .catch(() => {});
  }, []);

  // The support info is assembled from an explicit allowlist of safe fields.
  // The server config held client-side also contains secrets (stream keys,
  // admin password, S3 credentials) that must never end up in text people
  // paste into public issues and chat rooms.
  const supportInfo = useMemo(() => {
    const lines = [
      releaseString || `Owncast v${versionNumber}`,
      `Stream status: ${online ? 'online' : 'offline'}`,
    ];
    if (online && health?.message) {
      lines.push(`Stream health: ${health.healthPercentage}% ${health.message}`);
    }
    lines.push(`Latency level: ${videoSettings.latencyLevel}`);
    if (videoCodec) {
      lines.push(`Video codec: ${videoCodec}`);
    }
    videoSettings.videoQualityVariants.forEach((variant, index) => {
      lines.push(`Output ${index + 1}: ${describeVariant(variant)}`);
    });
    if (recentErrors.length > 0) {
      lines.push('', 'Recent warnings and errors:');
      recentErrors.forEach(entry => {
        lines.push(`[${entry.time}] ${entry.level}: ${entry.message}`);
      });
    }
    return lines.join('\n');
  }, [releaseString, versionNumber, online, health, videoSettings, videoCodec, recentErrors]);

  const copySupportInfo = async () => {
    try {
      await navigator.clipboard.writeText(supportInfo);
      message.success(t(Localization.Admin.Help.copied));
    } catch {
      // Clipboard access can be denied; the text is visible for manual copy.
    }
  };

  const searchDocs = (query: string) => {
    if (!query) {
      return;
    }
    window.open(`${DOCS_ROOT}/search?q=${encodeURIComponent(query)}`, '_blank', 'noopener');
  };

  const helpActions = [
    {
      icon: <ToolTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.fixTitle),
      description: t(Localization.Admin.Help.fixDescription),
      button: t(Localization.Admin.Help.fixButton),
      href: docsUrl('/troubleshoot'),
    },
    {
      icon: <MessageTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.communityTitle),
      description: t(Localization.Admin.Help.communityDescription),
      button: t(Localization.Admin.Help.communityButton),
      href: docsUrl('/chat'),
    },
    {
      icon: <BugTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.bugTitle),
      description: t(Localization.Admin.Help.bugDescription),
      button: t(Localization.Admin.Help.bugButton),
      href: 'https://github.com/owncast/owncast/issues/new/choose',
    },
  ];

  const taskGroups: { title: string; tasks: HelpTask[] }[] = [
    {
      title: t(Localization.Admin.Help.groupSetup),
      tasks: [
        {
          title: t(Localization.Admin.Help.taskBroadcasting),
          adminPage: '/admin/config/server',
          docsPath: '/docs/broadcasting',
        },
        {
          title: t(Localization.Admin.Help.taskVideoQuality),
          adminPage: '/admin/config-video',
          docsPath: '/docs/video',
        },
        {
          title: t(Localization.Admin.Help.taskStorage),
          adminPage: '/admin/config/server',
          docsPath: '/docs/storage',
        },
        {
          title: t(Localization.Admin.Help.taskSsl),
          docsPath: '/docs/sslproxies',
        },
      ],
    },
    {
      title: t(Localization.Admin.Help.groupCustomize),
      tasks: [
        {
          title: t(Localization.Admin.Help.taskWebsite),
          adminPage: '/admin/config/general',
          docsPath: '/docs/configuration/website',
        },
        {
          title: t(Localization.Admin.Help.taskChat),
          adminPage: '/admin/config-chat',
          docsPath: '/docs/chat/moderation',
        },
        {
          title: t(Localization.Admin.Help.taskNotifications),
          adminPage: '/admin/config-notify',
          docsPath: '/docs/configuration/notifications',
        },
      ],
    },
    {
      title: t(Localization.Admin.Help.groupGrow),
      tasks: [
        {
          title: t(Localization.Admin.Help.taskFediverse),
          adminPage: '/admin/config-federation',
          docsPath: '/docs/social',
        },
        {
          title: t(Localization.Admin.Help.taskDirectory),
          adminPage: '/admin/config/general',
          docsPath: '/docs/directory',
        },
        {
          title: t(Localization.Admin.Help.taskEmbed),
          docsPath: '/docs/embed',
        },
      ],
    },
    {
      title: t(Localization.Admin.Help.groupExtend),
      tasks: [
        {
          title: t(Localization.Admin.Help.taskPlugins),
          adminPage: '/admin/plugins',
          docsPath: '/docs/plugins',
        },
        {
          title: t(Localization.Admin.Help.taskApis),
          docsPath: '/docs/extend',
        },
      ],
    },
  ];

  const moreResources = [
    { title: t(Localization.Admin.Help.documentation), href: docsUrl('/docs') },
    { title: t(Localization.Admin.Help.releaseNotes), href: docsUrl('/releases') },
    {
      title: t(Localization.Admin.Help.discussions),
      href: 'https://github.com/owncast/owncast/discussions',
    },
    { title: t(Localization.Admin.Help.fediverse), href: 'https://social.owncast.online/@owncast' },
  ];

  return (
    <div className="help-page">
      <Title style={{ textAlign: 'center' }}>{t(Localization.Admin.Help.title)}</Title>

      <Row justify="center" style={{ marginBottom: '24px' }}>
        <Col xs={24} md={16} lg={12}>
          <Input.Search
            size="large"
            placeholder={t(Localization.Admin.Help.searchPlaceholder)}
            onSearch={searchDocs}
            enterButton
            allowClear
          />
        </Col>
      </Row>

      {upgradeVersion && (
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: '16px' }}
          message={`${t(Localization.Admin.Help.upgradeAvailable)} (v${upgradeVersion})`}
          action={
            <Button size="small" href="/admin/upgrade">
              {t(Localization.Admin.Help.upgradeLink)}
            </Button>
          }
        />
      )}

      {!online && overallPeakViewerCount === 0 && (
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: '16px' }}
          message={t(Localization.Admin.Help.gettingStarted)}
          action={
            <Button
              size="small"
              href={docsUrl('/quickstart')}
              target="_blank"
              rel="noopener noreferrer"
            >
              {t(Localization.Admin.Help.gettingStartedLink)}
            </Button>
          }
        />
      )}

      <Row gutter={[16, 16]}>
        {helpActions.map(action => (
          <Col xs={24} lg={8} key={action.title}>
            <Card style={{ height: '100%' }}>
              <Card.Meta
                avatar={action.icon}
                title={action.title}
                description={
                  <Space direction="vertical">
                    {action.description}
                    <Button
                      type="primary"
                      href={action.href}
                      target="_blank"
                      rel="noopener noreferrer"
                      icon={<LinkOutlined />}
                    >
                      {action.button}
                    </Button>
                  </Space>
                }
              />
            </Card>
          </Col>
        ))}
      </Row>

      <Divider />
      <Title level={2}>{t(Localization.Admin.Help.commonTasks)}</Title>
      {taskGroups.map(group => (
        <div key={group.title}>
          <Title level={4}>{group.title}</Title>
          <Row gutter={[16, 16]} style={{ marginBottom: '16px' }}>
            {group.tasks.map(task => (
              <Col xs={24} lg={12} xl={8} key={task.title}>
                <Card size="small" style={{ height: '100%' }}>
                  <Text strong>{task.title}</Text>
                  <div style={{ marginTop: '8px' }}>
                    <Space size="large">
                      {task.adminPage && (
                        <Link href={task.adminPage}>
                          <SettingOutlined /> {t(Localization.Admin.Help.openSettings)}
                        </Link>
                      )}
                      {task.docsPath && (
                        <a href={docsUrl(task.docsPath)} target="_blank" rel="noopener noreferrer">
                          <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
                        </a>
                      )}
                    </Space>
                  </div>
                </Card>
              </Col>
            ))}
          </Row>
        </div>
      ))}

      <Divider />
      <Title level={2}>{t(Localization.Admin.Help.supportInfo)}</Title>
      <Paragraph>{t(Localization.Admin.Help.supportInfoDescription)}</Paragraph>
      <Button type="primary" icon={<CopyOutlined />} onClick={copySupportInfo}>
        {t(Localization.Admin.Help.copySupportInfo)}
      </Button>

      <Divider />
      <Title level={2}>{t(Localization.Admin.Help.moreResources)}</Title>
      <Space direction="vertical">
        {moreResources.map(resource => (
          <a key={resource.title} href={resource.href} target="_blank" rel="noopener noreferrer">
            <LinkOutlined /> {resource.title}
          </a>
        ))}
      </Space>
    </div>
  );
}

Help.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
