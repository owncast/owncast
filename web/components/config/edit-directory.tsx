// Note: references to "yp" in the app are likely related to Owncast Directory
import React, { useState, useContext, useEffect } from 'react';
import { Typography } from 'antd';

import ToggleSwitch from './form-toggleswitch-with-submit';

import { ServerStatusContext } from '../../utils/server-status-context';
import { FIELD_PROPS_NSFW, FIELD_PROPS_YP } from '../../utils/config-constants';

const { Title } = Typography;

export default function EditYPDetails() {
  const [formDataValues, setFormDataValues] = useState(null);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { yp, instanceDetails } = serverConfig;
  const { nsfw } = instanceDetails;
  const { enabled, instanceUrl } = yp;

  useEffect(() => {
    setFormDataValues({
      ...yp,
      enabled,
      nsfw,
    });
  }, [yp, instanceDetails]);

  const hasInstanceUrl = instanceUrl !== '';
  if (!formDataValues) {
    return null;
  }
  return (
    <div className="config-directory-details-form">
      <Title level={3} className="section-title">
        Owncast Directory Settings
      </Title>

      <p className="description">
        Would you like to appear in the{' '}
        <a href="https://directory.owncast.online" target="_blank" rel="noreferrer">
          <strong>Owncast Directory</strong>
        </a>
        ?
      </p>

      <p style={{ backgroundColor: 'black', fontSize: '.75rem', padding: '5px' }}>
        <em>
          NOTE: You will need to have a URL specified in the <code>Instance URL</code> field to be
          able to use this.
        </em>
      </p>

      <div className="config-yp-container">
        <ToggleSwitch
          fieldName="enabled"
          {...FIELD_PROPS_YP}
          checked={formDataValues.enabled}
          disabled={!hasInstanceUrl}
        />
        <ToggleSwitch
          fieldName="nsfw"
          {...FIELD_PROPS_NSFW}
          checked={formDataValues.nsfw}
          disabled={!hasInstanceUrl}
        />
      </div>
    </div>
  );
}
