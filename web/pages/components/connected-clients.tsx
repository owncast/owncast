import React, { useState, useEffect } from 'react';
import { CONNECTED_CLIENTS, fetchData, FETCH_INTERVAL } from '../utils/apis';

export default function HardwareInfo() {
  const [clients, setClients] = useState({});

/*
geo data looks like this
  "geo": {
    "countryCode": "US",
    "regionName": "California",
    "timeZone": "America/Los_Angeles"
  }
*/

  const getInfo = async () => {
    try {
      const result = await fetchData(CONNECTED_CLIENTS);
      console.log("================ result", result)

      setClients({ ...result });

    } catch (error) {
      setClients({ ...clients, message: error.message });
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getInfo();
    getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
  
    // returned function will be called on component unmount 
    return () => {
      clearInterval(getStatusIntervalId);
    }
  }, []);

  return (
    <div>
      <h2>Connected Clients</h2>
      <p>a table of info..</p>
      <p>who's watching, how long they've been there, have they chatted? where they from?</p>
      <div style={{border: '1px solid purple', height: '300px', width: '100%', overflow:'auto'}}>
        {JSON.stringify(clients)}
      </div>
    </div>
  );
}
