import { CheckCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import { Alert, Button, Col, Row, Statistic } from 'antd';
import Link from 'next/link';
import React, { useContext } from 'react';
import { ServerStatusContext } from '../utils/server-status-context';

export default function StreamHealthOverview() {
  const serverStatusData = useContext(ServerStatusContext);
  const { health } = serverStatusData;

  if (!health) {
    return null;
  }

  const { healthy, healthPercentage, message } = health;
  console.log(healthPercentage);
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
      <Row gutter={16}>
        <Col span={12}>
          <Statistic
            title="Healthy Stream"
            value={healthy ? 'Yes' : 'No'}
            valueStyle={{ color }}
            prefix={healthy ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}
          />
        </Col>
        <Col span={12}>
          <Statistic
            title="Playback Health"
            value={healthPercentage}
            valueStyle={{ color }}
            suffix="%"
          />
        </Col>
      </Row>
      <Row gutter={16} style={{ display: message ? 'grid' : 'none', marginTop: '10px' }}>
        <Col span={24}>
          <Alert
            message={message}
            type={icon}
            showIcon
            action={
              <Link passHref href="/stream-health">
                <Button size="small" type="text" style={{ color: 'black' }}>
                  TROUBLESHOOT
                </Button>
              </Link>
            }
          />
        </Col>
      </Row>
    </div>
  );
}
