import Link from 'next/link';
import { Result, Card, Row, Col } from 'antd';
import {
  MessageTwoTone,
  QuestionCircleTwoTone,
  BookTwoTone,
  PlaySquareTwoTone,
  ProfileTwoTone,
} from '@ant-design/icons';
import OwncastLogo from '../components/logo';
import LogTable from '../components/log-table';

const { Meta } = Card;

export default function Offline({ logs = [], config }) {
  const data = [
    {
      icon: <BookTwoTone twoToneColor="#6f42c1" />,
      title: 'Use your broadcasting software',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/broadcasting/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn how to point your existing software to your new server and start streaming your
            content.
          </a>
        </div>
      ),
    },
    {
      icon: <MessageTwoTone twoToneColor="#0366d6" />,
      title: 'Chat is disabled',
      content: 'Chat will continue to be disabled until you begin a live stream.',
    },
    {
      icon: <PlaySquareTwoTone twoToneColor="#f9826c" />,
      title: 'Embed your video onto other sites',
      content: (
        <div>
          <a href="https://owncast.online/docs/embed" target="_blank" rel="noopener noreferrer">
            Learn how you can add your Owncast stream to other sites you control.
          </a>
        </div>
      ),
    },
    {
      icon: <QuestionCircleTwoTone twoToneColor="#ffd33d" />,
      title: 'Not sure what to do next?',
      content: (
        <div>
          If you&apos;re having issues or would like to know how to customize and configure your
          Owncast server visit <Link href="/help">the help page.</Link>
        </div>
      ),
    },
  ];

  if (!config?.yp?.enabled) {
    data.push({
      icon: <ProfileTwoTone twoToneColor="#D18BFE" />,
      title: 'Find an audience on the Owncast Directory',
      content: (
        <div>
          List yourself in the Owncast Directory and show off your stream. Enable it in{' '}
          <Link href="/config-public-details">settings.</Link>
        </div>
      ),
    });
  }

  return (
    <>
      <Row gutter={[16, 16]} className="offline-content">
        <Col span={12} xs={24} sm={24} md={24} lg={12} className="logo-section">
          <Result
            icon={<OwncastLogo />}
            title="No stream is active."
            subTitle="You should start one."
          />
        </Col>

        <Col span={12} xs={24} sm={24} md={24} lg={12} className="list-section">
          {data.map(item => (
            <Card key={item.title} size="small" bordered={false}>
              <Meta avatar={item.icon} title={item.title} description={item.content} />
            </Card>
          ))}
        </Col>
      </Row>
      <LogTable logs={logs} pageSize={5} />
    </>
  );
}
