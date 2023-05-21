import React, { useState, FC, ReactElement } from 'react';

export const AlertMessageContext = React.createContext({
  message: null,
  setMessage: (_text?: string) => null,
});

export type AlertMessageProviderProps = {
  children: ReactElement;
};

const AlertMessageProvider: FC<AlertMessageProviderProps> = ({ children }) => {
  const [message, setMessage] = useState('');

  const providerValue = {
    message,
    setMessage,
  };
  return (
    <AlertMessageContext.Provider value={providerValue}>{children}</AlertMessageContext.Provider>
  );
};

export default AlertMessageProvider;
