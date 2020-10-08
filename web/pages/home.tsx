import React, { useState, useEffect } from 'react';

import BroadcastInfo from './components/broadcast-info';
import HardwareInfo from './components/hardware-info';
import ViewerInfo from './components/viewer-info';

export default function HomeView(props) {
  const { broadcastActive, broadcaster, message } = props;

   const broadcastDetails = broadcastActive ? (
    <>
      <BroadcastInfo {...broadcaster} />
      <HardwareInfo />
      <ViewerInfo />
    </>
  ) : null;

  const disconnectButton = broadcastActive ? <button type="button">Boot (Disconnect)</button> : null;
  
  return (
    <div style={{padding: '2em'}}>
      <p>
        <b>Status: {broadcastActive ? 'on' : 'off'}</b>
      </p>

      <h2>Utilities</h2>
      (these dont do anything yet)
      {disconnectButton}
      <button type="button">Change Stream Key</button>
      <button type="button">Server Config</button>

      <br />
      <br />

      {broadcastDetails}      
    </div>
  );
}
