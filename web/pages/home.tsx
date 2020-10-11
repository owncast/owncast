import React, { useState, useEffect } from 'react';

import BroadcastInfo from './components/broadcast-info';
import HardwareInfo from './components/hardware-info';
import ViewerInfo from './components/viewer-info';
import ConnectedClients from './components/connected-clients';

export default function HomeView(props) {
  const { broadcastActive, broadcaster, message } = props;

  const broadcastDetails = broadcastActive ? (
    <>
      <BroadcastInfo {...broadcaster} />
      <HardwareInfo />
      <ViewerInfo />
      <ConnectedClients />
    </>
  ) : null;

  const disconnectButton = broadcastActive ? <button type="button">Boot (Disconnect)</button> : null;

  return (
    <div style={{ padding: '2em' }}>
      <p>
        <b>Status: {broadcastActive ? 'on' : 'off'}</b>
      </p>

      <h2>Utilities</h2>
      <p>(these dont do anything yet)</p>

      {disconnectButton}
      <button type="button">Change Stream Key</button>
      <button type="button">Server Config</button>

      <br />
      <br />

      {broadcastDetails}
    </div>
  );
}
