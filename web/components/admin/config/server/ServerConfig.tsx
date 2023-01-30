import React from 'react';
import EditInstanceDetails from '../../EditInstanceDetails2';

// eslint-disable-next-line react/function-component-definition
export default function ConfigServerDetails() {
  return (
    <div className="config-server-details-form">
      <p className="description">
        You should change your admin password from the default and keep it safe. For most people
        it&apos;s likely the other settings will not need to be changed.
      </p>
      <div className="form-module config-server-details-container">
        <EditInstanceDetails />
      </div>
    </div>
  );
}
