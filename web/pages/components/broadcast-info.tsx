import React, { useState, useEffect } from 'react';


export default function BroadcastInfo(props) {
const { remoteAddr, streamDetails, time } = props;
 
  return (
    <div style={{border: '1px solid green', width: '100%'}}>
      <h2>Broadcast Info</h2>
      <p>Remote Address: {remoteAddr}</p>  
      <p>Time: {(new Date(time)).toLocaleTimeString()}</p>  
      <p>Stream Details: {JSON.stringify(streamDetails)}</p>  
    </div>
  );
}
