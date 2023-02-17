import React from 'react';

import EditInstanceDetails from './EditInstanceDetails';
import EditInstanceTags from './EditInstanceTags';
import EditSocialLinks from './EditSocialLinks';
import EditPageContent from './EditPageContent';

// eslint-disable-next-line react/function-component-definition
export default function PublicFacingDetails() {
  return (
    <div className="config-public-details-page">
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
    </div>
  );
}
