import React from 'react';
import { Typography } from 'antd';

import EditInstanceDetails from '../../components/config/EditInstanceDetails';
import EditInstanceTags from '../../components/config/EditInstanceTags';
import EditSocialLinks from '../../components/config/EditSocialLinks';
import EditPageContent from '../../components/config/EditPageContent';
import EditCustomStyles from '../../components/config/EditCustomStyles';

const { Title } = Typography;

export default function PublicFacingDetails() {
  return (
    <div className="config-public-details-page">
      <Title>General Settings</Title>
      <p className="description">
        The following are displayed on your site to describe your stream and its content.{' '}
        <a
          href="https://owncast.online/docs/website/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn more.
        </a>
      </p>

      <div className="top-container">
        <div className="form-module instance-details-container">
          <EditInstanceDetails />
        </div>

        <div className="form-module social-items-container ">
          <div className="form-module tags-module">
            <EditInstanceTags />
          </div>

          <div className="form-module social-handles-container">
            <EditSocialLinks />
          </div>
        </div>
      </div>
      <div className="form-module page-content-module">
        <EditPageContent />
      </div>
      <div className="form-module page-content-module">
        <EditCustomStyles />
      </div>
    </div>
  );
}
