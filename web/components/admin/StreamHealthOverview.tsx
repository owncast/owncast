import { Alert, Button, Card, Col, Row, Statistic, Typography } from 'antd';
import dynamic from 'next/dynamic';
import Link from 'next/link';
import React, { FC, useContext } from 'react';
import { useTranslation } from 'next-export-i18n';
import { ServerStatusContext } from '../../utils/server-status-context';

// Lazy loaded components

const CheckCircleOutlined = dynamic(() => import('@ant-design/icons/CheckCircleOutlined'), {
  ssr: false,
});

const ExclamationCircleOutlined = dynamic(
  () => import('@ant-design/icons/ExclamationCircleOutlined'),
  {
    ssr: false,
  },
);

export type StreamHealthOverviewProps = {
  showTroubleshootButton?: Boolean;
};

export const StreamHealthOverview: FC<StreamHealthOverviewProps> = ({ showTroubleshootButton }) => {
  const { t } = useTranslation();
  const serverStatusData = useContext(ServerStatusContext);
  const { health } = serverStatusData;
  if (!health) {
    return null;
  }

  const { healthy, healthPercentage, message, representation } = health;
  let color = '#3f8600';
  let icon: 'success' | 'info' | 'warning' | 'error' = 'info';
  if (healthPercentage < 80) {
    color = '#cf000f';
    icon = 'error';
  } else if (healthPercentage < 30) {
    color = '#f0ad4e';
    icon = 'error';
  }

  return (
    <div>
      <Card type="inner">
        <Row gutter={8}>
          <Col span={12}>
            <Statistic
              title={t('Healthy Stream')}
              value={healthy ? t('Yes') : t('No')}
              valueStyle={{ color }}
              prefix={healthy ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}
            />
          </Col>
          <Col span={12}>
            <Statistic
              title={t('Playback Health')}
              value={healthPercentage}
              valueStyle={{ color }}
              suffix="%"
            />
          </Col>
        </Row>
        <Row style={{ display: representation < 100 && representation !== 0 ? 'grid' : 'none' }}>
          <Typography.Text
            type="secondary"
            style={{ textAlign: 'center', fontSize: '0.7em', opacity: '0.3' }}
          >
            {`${t('Stream health represents')} ${representation}% ${t('of all known players. Other player status is unknown.')}`}
          </Typography.Text>
        </Row>
        <Row
          gutter={16}
          style={{ width: '100%', display: message ? 'grid' : 'none', marginTop: '10px' }}
        >
          <Col span={24}>
            <Alert
              message={message}
              type={icon}
              showIcon
              action={
                showTroubleshootButton && (
                  <Link passHref href="/admin/stream-health">
                    <Button size="small" type="text" style={{ color: 'black' }}>
                      {t('TROUBLESHOOT')}
                    </Button>
                  </Link>
                )
              }
            />
          </Col>
        </Row>
      </Card>
    </div>
  );
};

StreamHealthOverview.defaultProps = {
  showTroubleshootButton: true,
};
