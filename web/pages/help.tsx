import { Button, Card, Col, Divider, Result, Row } from 'antd';
import Meta from 'antd/lib/card/Meta';
import Title from 'antd/lib/typography/Title';
import {
  AlertOutlined,
  ApiTwoTone,
  BookOutlined,
  BugTwoTone,
  CameraTwoTone,
  DatabaseTwoTone,
  EditTwoTone,
  Html5TwoTone,
  LinkOutlined,
  QuestionCircleTwoTone,
  SettingTwoTone,
  SlidersTwoTone,
} from '@ant-design/icons';
import React from 'react';

interface Props {}

export default function Help(props: Props) {
  const questions = [
    {
      icon: <SettingTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to configure my owncast instance',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/configuration/"
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
      title: 'I need help configuring my broadcasting software',
      content: (
        <div>
          <a
            href="https://owncast.online/docs/broadcasting/"
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
          <a href="https://owncast.online/docs/embed/" target="_blank" rel="noopener noreferrer">
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
          <a href="https://owncast.online/docs/website/" target="_blank" rel="noopener noreferrer">
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <SlidersTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to tweak my encoding quality or performance',
      content: (
        <div>
          <a href="https://owncast.online/docs/encoding/" target="_blank" rel="noopener noreferrer">
            <LinkOutlined /> Learn more
          </a>
        </div>
      ),
    },
    {
      icon: <DatabaseTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to offload my video to an external storage provider',
      content: (
        <div>
          <a href="https://owncast.online/docs/storage/" target="_blank" rel="noopener noreferrer">
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
          <a href="https://owncast.online/docs/faq/" target="_blank" rel="noopener noreferrer">
            {' '}
            FAQ
          </a>{' '}
          or exist in our <a href="https://github.com/owncast/owncast/discussions">discussions</a>
        </div>
      ),
    },
    {
      icon: <ApiTwoTone style={{ fontSize: '24px' }} />,
      title: 'I want to build add-ons for my Owncast server',
      content: (
        <div>
          You can build your own bots, overlays, tools and add-ons with our
          <a href="https://owncast.online/thirdparty" target="_blank" rel="noopener noreferrer">
            &nbsp;developer APIs.&nbsp;
          </a>
        </div>
      ),
    },
  ];

  return (
    <div>
      <Title style={{ textAlign: 'center' }}>How can we help you?</Title>
      <Row gutter={[16, 16]} justify="space-around" align="middle">
        <Col xs={24} lg={12} style={{ textAlign: 'center' }}>
          <Result status="500" />
          <Title level={2}>Troubleshooting</Title>
          <Button
            target="_blank"
            rel="noopener noreferrer"
            href="https://owncast.online/docs/troubleshooting/"
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
            href="https://owncast.online/docs"
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
          <Col xs={24} lg={12}>
            <Card key={question.title}>
              <Meta avatar={question.icon} title={question.title} description={question.content} />
            </Card>
          </Col>
        ))}
      </Row>
      <Divider />
      <Title level={2}>Other</Title>
      <Row gutter={[16, 16]}>
        {otherResources.map(question => (
          <Col xs={24} lg={12}>
            <Card key={question.title}>
              <Meta avatar={question.icon} title={question.title} description={question.content} />
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}
