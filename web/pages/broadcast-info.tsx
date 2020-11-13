import React, { useContext } from 'react';
import { ServerStatusContext } from '../utils/server-status-context';


export default function BroadcastInfo() {
  const context = useContext(ServerStatusContext);
  const { broadcaster } = context || {};
  const { remoteAddr, time, streamDetails } = broadcaster || {};
 
  return (
    <div style={{border: '1px solid green', width: '100%'}}>
      <h2>Broadcast Info</h2>
      <p>Remote Address: {remoteAddr}</p>  
      <p>Time: {(new Date(time)).toLocaleTimeString()}</p>  
      <p>Stream Details: {JSON.stringify(streamDetails)}</p>  
    </div>
  );
}
