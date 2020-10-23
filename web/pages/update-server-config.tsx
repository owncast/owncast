import React, { useState, useEffect, useContext } from 'react';
import { SERVER_CONFIG, fetchData, FETCH_INTERVAL } from './utils/apis';
import { BroadcastStatusContext } from './utils/broadcast-status-context';

export default function ServerConfig() {
  const context = useContext(BroadcastStatusContext);
  const { broadcastActive } = context || {};
  

  const [clients, setClients] = useState({});

  const getInfo = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      console.log("viewers result", result)

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
      <h2>Server Config</h2>
      <p>Display this data all pretty, most things will be editable in the future, not now.</p>
      <div style={{border: '1px solid pink', height: '300px', width: '100%', overflow:'auto'}}>
        {JSON.stringify(clients)}
      </div>
    </div>
  );
}
