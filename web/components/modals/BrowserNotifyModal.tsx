/* eslint-disable react/no-unescaped-entities */
/* eslint-disable @typescript-eslint/no-unused-vars */
import { Row, Col, Switch, Typography } from 'antd';
import { useState } from 'react';

const { Title } = Typography;

// interface Props {}

export default function BrowserNotifyModal() {
  const [enabled, setEnabled] = useState(false);

  const onSwitchToggle = (checked: Boolean) => {
    setEnabled(true);
  };

  return (
    <div>
      <Row align="top">
        <Col span={12}>
          <Switch defaultChecked={enabled} checked={enabled} onChange={onSwitchToggle} />{' '}
          {enabled ? 'Enabled' : 'Disabled'}
        </Col>
        <Col span={12}>
          You'll need to allow your browser to receive notifications from Owncast Nightly, first.
          Fake push notification prompt example goes here.
        </Col>
      </Row>
      <Row align="top">
        <Title>Browser Notifications</Title>
        Get notified right in the browser each time this stream goes live. Blah blah blah more
        description text goes here.
      </Row>
    </div>
  );
}
