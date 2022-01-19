import { h, Component } from '/js/web_modules/preact.js';
import { useState } from '/js/web_modules/preact/hooks.js';

import htm from '/js/web_modules/htm.js';
import { ExternalActionButton } from './external-action-modal.js';
import {
  registerWebPushNotifications,
  isPushNotificationSupported,
} from '../notification/registerWeb.js';
import { URL_REGISTER_NOTIFICATION } from '../utils/constants.js';

const html = htm.bind(h);

export function NotifyModal({ notifications, streamName, accessToken }) {
  const [error, setError] = useState(null);
  const [loaderStyle, setLoaderStyle] = useState('none');

  const [browserRegistrationComplete, setBrowserRegistrationComplete] =
    useState(false);

  const browserNotificationsButtonState =
    Notification.permission === 'default'
      ? ''
      : 'cursor-not-allowed opacity-50';

  const { browser } = notifications;
  const { publicKey } = browser;

  let browserPushEnabled = browser.enabled;

  // Browser push notifications are only supported on Chrome currently.
  // Also make sure the browser supports them.
  if (!window.chrome || !isPushNotificationSupported()) {
    browserPushEnabled = false;
  }

  async function saveNotificationRegistration(channel, destination) {
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ channel: channel, destination: destination }),
    };

    try {
      await fetch(
        URL_REGISTER_NOTIFICATION + `?accessToken=${accessToken}`,
        options
      );
    } catch (e) {
      console.error(e);
    }
  }

  async function registerForBrowserPushButtonPressed() {
    // If it's already denied or granted, don't do anything.
    if (Notification.permission !== 'default') {
      return;
    }

    try {
      const subscription = await registerWebPushNotifications(publicKey);
      saveNotificationRegistration('BROWSER_PUSH_NOTIFICATION', subscription);
      setBrowserRegistrationComplete(true);
      setError(null);
    } catch (e) {
      setError(
        `Error registering for notifications: ${e.message}. Check your browser notification permissions and try again.`
      );
    }
  }

  function getBrowserPushButtonText() {
    if (browserRegistrationComplete) {
      return 'Done!';
    }

    let pushNotificationButtonText = 'Notify Me!';
    if (Notification.permission === 'granted') {
      pushNotificationButtonText = 'Enabled';
    } else if (Notification.permission === 'denied') {
      pushNotificationButtonText = 'Browser pushes denied';
    }
    return pushNotificationButtonText;
  }

  const pushNotificationButtonText = getBrowserPushButtonText();

  return html`
    <div class="bg-gray-100 bg-center bg-no-repeat p-4">
      <p class="text-gray-700 text-md font-semibold mb-2">
        Never miss a stream! Get notified when ${streamName} goes live.
      </p>

      <div style=${{ display: browserPushEnabled ? 'block' : 'none' }}>
        <h2 class="text-indigo-600 text-3xl font-semibold">Browser</h2>
        <p>Get notified inside your browser when ${streamName} goes live.</p>
        <p class="mt-4">
          Make sure you change the setting from ${' '}
          <em>until I close this site</em> or you won't receive future
          notifications.
        </p>
        <img class="mt-4" src="/img/browser-push-notifications-settings.png" />
        <div
          style=${error ? 'display: block' : 'display: none'}
          class="bg-red-100 border border-red-400 text-red-700 mt-5 px-4 py-3 rounded relative"
          role="alert"
        >
          ${error}
        </div>
        <button
          class="rounded-sm flex flex-row justify-center items-center overflow-hidden mt-5 px-3 py-1 text-base text-white bg-gray-800 ${browserNotificationsButtonState}"
          onClick=${registerForBrowserPushButtonPressed}
        >
          ${pushNotificationButtonText}
        </button>
      </div>

      <div
        id="follow-loading-spinner-container"
        style="display: ${loaderStyle}"
      >
        <img id="follow-loading-spinner" src="/img/loading.gif" />
      </div>
    </div>
  `;
}

export function NotifyButton({ serverName, federationInfo, onClick }) {
  const notifyAction = {
    // color: 'rgba(28, 26, 59, 1)',
    description: `Notify`,
    icon: '/img/notification-bell.svg',
    openExternally: false,
  };

  return html`
    <span id="notify-button-container">
      <${ExternalActionButton} onClick=${onClick} action=${notifyAction} />
    </span>
  `;
}
