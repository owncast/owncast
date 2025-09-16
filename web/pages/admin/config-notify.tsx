import { Alert, Button, Col, Row, Typography } from 'antd';
import React, { ReactElement, useContext, useEffect, useState } from 'react';
import Link from 'next/link';
import { useTranslation } from 'next-export-i18n';

import { DiscordNotify as Discord } from '../../components/admin/notification/discord';
import { BrowserNotify as Browser } from '../../components/admin/notification/browser';
import { FediverseNotify as Federation } from '../../components/admin/notification/federation';
import {
  TextFieldWithSubmit,
  TEXTFIELD_TYPE_URL,
} from '../../components/admin/TextFieldWithSubmit';
import { TEXTFIELD_PROPS_FEDERATION_INSTANCE_URL } from '../../utils/config-constants';
import { ServerStatusContext } from '../../utils/server-status-context';
import { UpdateArgs } from '../../types/config-section';
import { isValidUrl } from '../../utils/validators';

import { AdminLayout } from '../../components/layouts/AdminLayout';
import { Localization } from '../../types/localization';

const { Title } = Typography;

export default function ConfigNotify() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { yp } = serverConfig;
  const { instanceUrl } = yp;
  const [urlValid, setUrlValid] = useState(false);
  const { t } = useTranslation();

  useEffect(() => {
    setFormDataValues({
      instanceUrl,
    });
  }, [yp]);

  const handleSubmitInstanceUrl = () => {
    if (!urlValid) {
      return;
    }

    setFormDataValues({
      ...formDataValues,
      enabled: false,
    });
  };

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setUrlValid(isValidUrl(value));

    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };

  const enabled = instanceUrl !== '';
  const configurationWarning = !enabled && (
    <>
      <Alert message={t(Localization.Admin.Alerts.serverUrlRequired)} type="warning" showIcon />
      <br />
      <TextFieldWithSubmit
        fieldName="instanceUrl"
        {...TEXTFIELD_PROPS_FEDERATION_INSTANCE_URL}
        value={formDataValues?.instanceUrl || ''}
        initialValue={yp.instanceUrl}
        type={TEXTFIELD_TYPE_URL}
        onChange={handleFieldChange}
        onSubmit={handleSubmitInstanceUrl}
        required
      />
    </>
  );

  return (
    <>
      <Title>{t(Localization.Admin.Config.notifications)}</Title>
      <p className="description">
        {t(Localization.Admin.Config.notificationsDescription)}{' '}
        <a
          href="https://owncast.online/docs/notifications/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          {t(Localization.Admin.Config.notificationsLearnMore)}
        </a>
      </p>

      {configurationWarning}

      <Row>
        <Col
          span={10}
          className={`form-module ${enabled ? '' : 'disabled'}`}
          style={{ margin: '5px', display: 'flex', flexDirection: 'column' }}
        >
          <Browser />
        </Col>

        <Col
          span={10}
          className={`form-module ${enabled ? '' : 'disabled'}`}
          style={{ margin: '5px', display: 'flex', flexDirection: 'column' }}
        >
          <Discord />
        </Col>

        <Col
          span={10}
          className={`form-module ${enabled ? '' : 'disabled'}`}
          style={{ margin: '5px', display: 'flex', flexDirection: 'column' }}
        >
          <Federation />
        </Col>

        <Col
          span={10}
          className={`form-module ${enabled ? '' : 'disabled'}`}
          style={{ margin: '5px', display: 'flex', flexDirection: 'column' }}
        >
          <Title>Custom</Title>
          <p className="description">Build your own notifications by using custom webhooks.</p>

          <Link passHref href="/admin/webhooks">
            <Button
              type="primary"
              style={{
                position: 'relative',
                marginLeft: 'auto',
                right: '0',
                marginTop: '20px',
              }}
            >
              Create
            </Button>
          </Link>
        </Col>
      </Row>
    </>
  );
}

ConfigNotify.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
