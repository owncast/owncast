import React from 'react';
import { EditInstanceDetails } from '../../../../../components/config/EditInstanceDetails2';

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
