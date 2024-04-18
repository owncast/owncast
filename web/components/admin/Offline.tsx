import { Card, Col, Row, Typography } from 'antd';
import Link from 'next/link';
import { FC, useContext } from 'react';
import dynamic from 'next/dynamic';
import { LogTable } from './LogTable';
import { OwncastLogo } from '../common/OwncastLogo/OwncastLogo';
import { NewsFeed } from './NewsFeed';
import { ConfigDetails } from '../../types/config-section';
import { ServerStatusContext } from '../../utils/server-status-context';

const { Paragraph, Text } = Typography;

const { Title } = Typography;
const { Meta } = Card;

// Lazy loaded components

const BookTwoTone = dynamic(() => import('@ant-design/icons/BookTwoTone'), {
  ssr: false,
});

const MessageTwoTone = dynamic(() => import('@ant-design/icons/MessageTwoTone'), {
  ssr: false,
});

const PlaySquareTwoTone = dynamic(() => import('@ant-design/icons/PlaySquareTwoTone'), {
  ssr: false,
});

const ProfileTwoTone = dynamic(() => import('@ant-design/icons/ProfileTwoTone'), {
  ssr: false,
});

function generateStreamURL(serverURL, rtmpServerPort) {
  return `rtmp://${serverURL.replace(/(^\w+:|^)\/\//, '')}:${rtmpServerPort}/live`;
}

export type OfflineProps = {
  logs: any[];
  config: ConfigDetails;
};

export const Offline: FC<OfflineProps> = ({ logs = [], config }) => {
  const serverStatusData = useContext(ServerStatusContext);

  const { serverConfig } = serverStatusData || {};
  const { rtmpServerPort, streamKeyOverridden } = serverConfig;
  const instanceUrl = global.window?.location.hostname || '';

  let rtmpURL;
  if (instanceUrl && rtmpServerPort) {
    rtmpURL = generateStreamURL(instanceUrl, rtmpServerPort);
  }

  const data = [
    {
      icon: <BookTwoTone twoToneColor="#6f42c1" />,
      title: 'Use your broadcasting software',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/broadcasting/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn how to point your existing software to your new server and start streaming your
            content.
          </a>
          <div className="stream-info-container">
            <Text strong className="stream-info-label">
              Streaming URL:
            </Text>
            {rtmpURL && (
              <Paragraph className="stream-info-box" copyable>
                {rtmpURL}
              </Paragraph>
            )}
            <Text strong className="stream-info-label">
              Streaming Keys:
            </Text>
            <Text strong className="stream-info-box">
              {!streamKeyOverridden ? (
                <Link href="/admin/config/server"> View </Link>
              ) : (
                <span style={{ paddingLeft: '10px', fontWeight: 'normal' }}>
                  Overridden via command line.
                </span>
              )}
            </Text>
          </div>
        </div>
      ),
    },
    {
      icon: <PlaySquareTwoTone twoToneColor="#f9826c" />,
      title: 'Embed your video onto other sites',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/embed?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn how you can add your Owncast stream to other sites you control.
          </a>
        </div>
      ),
    },
  ];

  if (!config?.chatDisabled) {
    data.push({
      icon: <MessageTwoTone twoToneColor="#0366d6" />,
      title: 'Chat is disabled',
      content: <span>Chat will continue to be disabled until you begin a live stream.</span>,
    });
  }

  if (!config?.yp?.enabled) {
    data.push({
      icon: <ProfileTwoTone twoToneColor="#D18BFE" />,
      title: 'Find an audience on the Owncast Directory',
      content: (
        <div>
          List yourself in the Owncast Directory and show off your stream. Enable it in{' '}
          <Link href="/admin/config/general/">settings.</Link>
        </div>
      ),
    });
  }

  if (!config?.federation?.enabled) {
    data.push({
      icon: <img alt="fediverse" width="20px" src="/img/fediverse-color.png" />,
      title: 'Add your Owncast instance to the Fediverse',
      content: (
        <div>
          <Link href="/admin/config-federation/">Enable Owncast social</Link> features to have your
          instance join the Fediverse, allowing people to follow, share and engage with your live
          stream.
        </div>
      ),
    });
  }

  return (
    <>
      <Row>
        <Col span={12} offset={6}>
          <div className="offline-intro">
            <span className="logo">
              <OwncastLogo variant="simple" />
            </span>
            <div>
              <Title level={2}>No stream is active</Title>
              <p>You should start one.</p>
            </div>
          </div>
        </Col>
      </Row>
      <Row gutter={[16, 16]} className="offline-content">
        <Col span={12} xs={24} sm={24} md={24} lg={12} className="list-section">
          {data.map(item => (
            <Card key={item.title} size="small" bordered={false}>
              <Meta avatar={item.icon} title={item.title} description={item.content} />
            </Card>
          ))}
        </Col>
        <Col span={12} xs={24} sm={24} md={24} lg={12}>
          <NewsFeed />
        </Col>
      </Row>
      <LogTable logs={logs} initialPageSize={5} />
    </>
  );
};
