import React, { useState, useEffect } from 'react';
import { BROADCASTER, fetchData, FETCH_INTERVAL } from './utils/apis';
import Main from './home';

export default function Admin() {
  const [broadcasterStatus, setBroadcasterStatus] = useState({});
  const [count, setCount] = useState(0);

  const getBroadcastStatus = async () => {
    try {
      const result = await fetchData(BROADCASTER);
      const broadcastActive = !!result.broadcaster;

      console.log("====",{count, result})

      setBroadcasterStatus({ ...result, broadcastActive });
      setCount(count => count + 1);

    } catch (error) {
      setBroadcasterStatus({ ...broadcasterStatus, message: error.message });
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getBroadcastStatus();
    getStatusIntervalId = setInterval(getBroadcastStatus, FETCH_INTERVAL);
  
    // returned function will be called on component unmount 
    return () => {
      clearInterval(getStatusIntervalId);
    }
  }, [])

  
  return (
    <Main {...broadcasterStatus} />
  );
}
