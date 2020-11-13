import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { STATUS, fetchData, FETCH_INTERVAL } from './apis';

const initialState = {
  broadcastActive: false,
  broadcaster: null,
  online: false,
  viewerCount: 0,
  sessionPeakViewerCount: 0,
  overallPeakViewerCount: 0,
  disableUpgradeChecks: true,
  versionNumber: '0.0.0',
};

export const ServerStatusContext = React.createContext(initialState);

const ServerStatusProvider = ({ children }) => {
  const [status, setStatus] = useState(initialState);

  const getStatus = async () => {
    try {
      const result = await fetchData(STATUS);
      setStatus({ ...result });

    } catch (error) {
      // setBroadcasterStatus({ ...broadcasterStatus, message: error.message });
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getStatus();
    getStatusIntervalId = setInterval(getStatus, FETCH_INTERVAL);
  
    // returned function will be called on component unmount 
    return () => {
      clearInterval(getStatusIntervalId);
    }
  }, [])

  return (
    <ServerStatusContext.Provider value={status}>
      {children}
    </ServerStatusContext.Provider>
  );
}

ServerStatusProvider.propTypes = {
  children: PropTypes.element.isRequired,
};

export default ServerStatusProvider;