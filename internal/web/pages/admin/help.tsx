import { Button, Card, Col, Divider, Result, Row } from 'antd';
import Meta from 'antd/lib/card/Meta';
import Title from 'antd/lib/typography/Title';

import React, { ReactElement } from 'react';
import dynamic from 'next/dynamic';

import { AdminLayout } from '../../components/layouts/AdminLayout';

// Lazy loaded components

const ApiTwoTone = dynamic(() => import('@ant-design/icons/ApiTwoTone'), {
  ssr: false,
});

const BugTwoTone = dynamic(() => import('@ant-design/icons/BugTwoTone'), {
  ssr: false,
});

const CameraTwoTone = dynamic(() => import('@ant-design/icons/CameraTwoTone'), {
  ssr: false,
});

const DatabaseTwoTone = dynamic(() => import('@ant-design/icons/DatabaseTwoTone'), {
  ssr: false,
});

const EditTwoTone = dynamic(() => import('@ant-design/icons/EditTwoTone'), {
  ssr: false,
});

const Html5TwoTone = dynamic(() => import('@ant-design/icons/Html5TwoTone'), {
  ssr: false,
});

const LinkOutlined = dynamic(() => import('@ant-design/icons/LinkOutlined'), {
  ssr: false,
});

const QuestionCircleTwoTone = dynamic(() => import('@ant-design/icons/QuestionCircleTwoTone'), {
  ssr: false,
});

const SettingTwoTone = dynamic(() => import('@ant-design/icons/SettingTwoTone'), {
  ssr: false,
});

const SlidersTwoTone = dynamic(() => import('@ant-design/icons/SlidersTwoTone'), {
  ssr: false,
});

export default function Help() {
  const questions = [
    {
      icon: <SettingTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to configure my owncast instance',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/configuration/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <CameraTwoTone style={{ fontSize: '24px' }} />,
      title: 'Help configuring my broadcasting software',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/broadcasting/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <Html5TwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to embed my stream into another site',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/embed/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <EditTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to customize my website',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/website/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <SlidersTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to tweak my video output',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/encoding/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <DatabaseTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to use an external storage provider',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/storage/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
  ];

  const otherResources = [
    {
      icon: <BugTwoTone style={{ fontSize: '24px' }} />,
      title: 'I found a bug',
      content: (
        <div>
          If you found a bug, then please
          <a
            href="https://github.com/owncast/owncast/issues/new/choose"
            target="_blank"
            rel="noopener noreferrer"
          >
            {' '}
            let us know
          </a>
        </div>
      ),
    },
    {
      icon: <QuestionCircleTwoTone style={{ fontSize: '24px' }} />,
      title: 'I have a general question',
      content: (
        <div>
          Most general questions are answered in our
          <a
            href="https://owncast.online/faq/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            {' '}
            FAQ
          </a>{' '}
          or exist in our{' '}
          <a
            href="https://github.com/owncast/owncast/discussions"
            target="_blank"
            rel="noopener noreferrer"
          >
            discussions
          </a>
        </div>
      ),
    },
    {
      icon: <ApiTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to build add-ons for Owncast',
      content: (
        <div>
          You can build your own bots, overlays, tools and add-ons with our
          <a
            href="https://owncast.online/thirdparty?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            &nbsp;developer APIs.&nbsp;
          </a>
        </div>
      ),
    },
  ];

  return (
    <div className="help-page">
      <Title style={{ textAlign: 'center' }}>How can we help you?</Title>
      <Row gutter={[16, 16]} justify="space-around" align="middle">
        <Col xs={24} lg={12} style={{ textAlign: 'center' }}>
          <Result status="500" />
          <Title level={2}>Troubleshooting</Title>
          <Button
            target="_blank"
            rel="noopener noreferrer"
            href="https://owncast.online/docs/troubleshooting/?source=admin"
            icon={<LinkOutlined />}
            type="primary"
          >
            Fix your problems
          </Button>
        </Col>
        <Col xs={24} lg={12} style={{ textAlign: 'center' }}>
          <Result status="404" />
          <Title level={2}>Documentation</Title>
          <Button
            target="_blank"
            rel="noopener noreferrer"
            href="https://owncast.online/docs?source=admin"
            icon={<LinkOutlined />}
            type="primary"
          >
            Read the Docs
          </Button>
        </Col>
      </Row>
      <Divider />
      <Title level={2}>Common tasks</Title>
      <Row gutter={[16, 16]}>
        {questions.map(question => (
          <Col xs={24} lg={12} key={question.title}>
            <Card>
              <Meta avatar={question.icon} title={question.title} description={question.content} />
            </Card>
          </Col>
        ))}
      </Row>
      <Divider />
      <Title level={2}>Other</Title>
      <Row gutter={[16, 16]}>
        {otherResources.map(question => (
          <Col xs={24} lg={12} key={question.title}>
            <Card>
              <Meta avatar={question.icon} title={question.title} description={question.content} />
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}

Help.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
