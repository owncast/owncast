import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { BROADCASTER, fetchData, FETCH_INTERVAL } from './apis';

const initialState = {
  broadcastActive: false,
  message: '',
  broadcaster: null,
};

export const BroadcastStatusContext = React.createContext(initialState);

const BroadcastStatusProvider = ({ children }) => {
  const [broadcasterStatus, setBroadcasterStatus] = useState(initialState);

  const getBroadcastStatus = async () => {
    try {
      const result = await fetchData(BROADCASTER);
      const broadcastActive = !!result.broadcaster || result.success;
      setBroadcasterStatus({ ...result, broadcastActive });

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
    <BroadcastStatusContext.Provider value={broadcasterStatus}>
      {children}
    </BroadcastStatusContext.Provider>
  );
}

BroadcastStatusProvider.propTypes = {
  children: PropTypes.element.isRequired,
};

export default BroadcastStatusProvider;