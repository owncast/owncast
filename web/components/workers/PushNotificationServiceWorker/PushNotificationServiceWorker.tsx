/* eslint-disable react/no-danger */
import { FC, useEffect } from 'react';

export const PushNotificationServiceWorker: FC = () => {
  const add = () => {
    navigator.serviceWorker.register('/serviceWorker.js').then(
      registration => {
        console.debug('Service Worker registration successful with scope: ', registration.scope);
      },
      err => {
        console.error('Service Worker registration failed: ', err);
      },
    );
  };

  useEffect(() => {
    if ('serviceWorker' in navigator) {
      window.addEventListener('load', add);
    }

    return () => {
      window.removeEventListener('load', add);
    };
  }, []);

  return null;
};
