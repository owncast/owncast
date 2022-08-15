import React, { useState } from 'react';
import PropTypes from 'prop-types';

export const AlertMessageContext = React.createContext({
  message: null,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setMessage: (text?: string) => null,
});

const AlertMessageProvider = ({ children }) => {
  const [message, setMessage] = useState('');

  const providerValue = {
    message,
    setMessage,
  };
  return (
    <AlertMessageContext.Provider value={providerValue}>{children}</AlertMessageContext.Provider>
  );
};

AlertMessageProvider.propTypes = {
  children: PropTypes.element.isRequired,
};

export default AlertMessageProvider;
