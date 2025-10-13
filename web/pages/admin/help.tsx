import { Button, Card, Col, Divider, Result, Row } from 'antd';
import Meta from 'antd/lib/card/Meta';
import Title from 'antd/lib/typography/Title';

import React, { ReactElement } from 'react';
import dynamic from 'next/dynamic';

import { useTranslation } from 'next-export-i18n';
import { Localization } from '../../types/localization';
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
  const { t } = useTranslation();

  const questions = [
    {
      icon: <SettingTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.configureInstance),
      content: (
        <div>
          <a
            href="https://owncast.online/docs/configuration/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
          </a>
        </div>
      ),
    },
    {
      icon: <CameraTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.configureBroadcasting),
      content: (
        <div>
          <a
            href="https://owncast.online/docs/broadcasting/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
          </a>
        </div>
      ),
    },
    {
      icon: <Html5TwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.embedStream),
      content: (
        <div>
          <a
            href="https://owncast.online/docs/embed/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
          </a>
        </div>
      ),
    },
    {
      icon: <EditTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.customizeWebsite),
      content: (
        <div>
          <a
            href="https://owncast.online/docs/website/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
          </a>
        </div>
      ),
    },
    {
      icon: <SlidersTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.tweakVideo),
      content: (
        <div>
          <a
            href="https://owncast.online/docs/encoding/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
          </a>
        </div>
      ),
    },
    {
      icon: <DatabaseTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.useStorage),
      content: (
        <div>
          <a
            href="https://owncast.online/docs/storage/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            <LinkOutlined /> {t(Localization.Admin.Help.learnMore)}
          </a>
        </div>
      ),
    },
  ];

  const otherResources = [
    {
      icon: <BugTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.foundBug),
      content: (
        <div>
          {t(Localization.Admin.Help.bugPlease)}
          <a
            href="https://github.com/owncast/owncast/issues/new/choose"
            target="_blank"
            rel="noopener noreferrer"
          >
            {' '}
            {t(Localization.Admin.Help.letUsKnow)}
          </a>
        </div>
      ),
    },
    {
      icon: <QuestionCircleTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.generalQuestion),
      content: (
        <div>
          {t(Localization.Admin.Help.generalAnswered)}
          <a
            href="https://owncast.online/faq/?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            {' '}
            {t(Localization.Admin.Help.faq)}
          </a>{' '}
          {t(Localization.Admin.Help.orExist)}{' '}
          <a
            href="https://github.com/owncast/owncast/discussions"
            target="_blank"
            rel="noopener noreferrer"
          >
            {t(Localization.Admin.Help.discussions)}
          </a>
        </div>
      ),
    },
    {
      icon: <ApiTwoTone style={{ fontSize: '24px' }} />,
      title: t(Localization.Admin.Help.buildAddons),
      content: (
        <div>
          {t(Localization.Admin.Help.buildTools)}
          <a
            href="https://owncast.online/thirdparty?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            &nbsp;{t(Localization.Admin.Help.developerApis)}&nbsp;
          </a>
        </div>
      ),
    },
  ];

  return (
    <div className="help-page">
      <Title style={{ textAlign: 'center' }}>{t(Localization.Admin.Help.title)}</Title>
      <Row gutter={[16, 16]} justify="space-around" align="middle">
        <Col xs={24} lg={12} style={{ textAlign: 'center' }}>
          <Result status="500" />
          <Title level={2}>{t(Localization.Admin.Help.troubleshooting)}</Title>
          <Button
            target="_blank"
            rel="noopener noreferrer"
            href="https://owncast.online/docs/troubleshooting/?source=admin"
            icon={<LinkOutlined />}
            type="primary"
          >
            {t(Localization.Admin.Help.fixProblems)}
          </Button>
        </Col>
        <Col xs={24} lg={12} style={{ textAlign: 'center' }}>
          <Result status="404" />
          <Title level={2}>{t(Localization.Admin.Help.documentation)}</Title>
          <Button
            target="_blank"
            rel="noopener noreferrer"
            href="https://owncast.online/docs?source=admin"
            icon={<LinkOutlined />}
            type="primary"
          >
            {t(Localization.Admin.Help.readDocs)}
          </Button>
        </Col>
      </Row>
      <Divider />
      <Title level={2}>{t(Localization.Admin.Help.commonTasks)}</Title>
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
      <Title level={2}>{t(Localization.Admin.Help.other)}</Title>
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
