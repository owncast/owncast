import { h, Component } from '/js/web_modules/preact.js';
import { useState } from '/js/web_modules/preact/hooks.js';

import htm from '/js/web_modules/htm.js';
import { ExternalActionButton } from './external-action-modal.js';
import {
  registerWebPushNotifications,
  isPushNotificationSupported,
} from '../notification/registerWeb.js';
import {
  URL_REGISTER_NOTIFICATION,
  URL_REGISTER_EMAIL_NOTIFICATION,
} from '../utils/constants.js';

const html = htm.bind(h);

export function NotifyModal({ notifications, streamName, accessToken }) {
  const [error, setError] = useState(null);
  const [loaderStyle, setLoaderStyle] = useState('none');
  const [emailNotificationsButtonEnabled, setEmailNotificationsButtonEnabled] =
    useState(false);
  const [emailAddress, setEmailAddress] = useState(null);
  const emailNotificationButtonState = emailNotificationsButtonEnabled
    ? ''
    : 'cursor-not-allowed opacity-50';

  const [browserRegistrationComplete, setBrowserRegistrationComplete] =
    useState(false);

  const browserNotificationsButtonState =
    Notification.permission === 'default'
      ? ''
      : 'cursor-not-allowed opacity-50';

  const { browser, email } = notifications;
  const { publicKey } = browser;

  let browserPushEnabled = browser.enabled;
  let emailEnabled = email.enabled;

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

  async function registerForEmailButtonPressed() {
    try {
      const options = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ emailAddress: emailAddress }),
      };

      try {
        await fetch(
          URL_REGISTER_EMAIL_NOTIFICATION + `?accessToken=${accessToken}`,
          options
        );
      } catch (e) {
        console.error(e);
      }
    } catch (e) {
      setError(`Error registering for email notifications: ${e.message}.`);
    }
  }

  function onEmailInput(e) {
    const { value } = e.target;

    // TODO: Add validation for email
    const valid = true;

    setEmailAddress(value);
    setEmailNotificationsButtonEnabled(valid);
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

  var gridColumns = 2;
  if (browserPushEnabled && !emailEnabled) {
    gridColumns = 1;
  } else if (!browserPushEnabled && emailEnabled) {
    gridColumns = 1;
  }

  const pushNotificationButtonText = getBrowserPushButtonText();

  return html`
    <div class="bg-gray-100 bg-center bg-no-repeat p-4">
      <p class="text-gray-700 text-md font-semibold mb-2">
        Never miss a stream! Get notified when ${streamName} goes live.
      </p>

      <div class="grid grid-cols-${gridColumns} gap-6">
        <div style=${{ display: browserPushEnabled ? 'block' : 'none' }}>
          <h2 class="text-indigo-600 text-3xl font-semibold">Browser</h2>
          <p>Get notified inside your browser when ${streamName} goes live.</p>
          <p class="mt-4">
            Make sure you change the setting from ${' '}
            <em>until I close this site</em> or you won't receive future
            notifications.
          </p>
          <img
            class="mt-4"
            src="/img/browser-push-notifications-settings.png"
          />
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
        <div style=${{ display: emailEnabled ? 'block' : 'none' }}>
          <h2 class="text-indigo-600 text-3xl font-semibold">Email</h2>
          <p>Get notified directly when ${''} ${streamName} goes live.</p>
          <p class="mt-4">
            Easily unsubscribe if you no longer want to receive them.
          </p>
          <input
            class="border bg-white rounded w-full py-2 px-3 mb-2 mt-10 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
            value=${emailAddress}
            onInput=${onEmailInput}
            placeholder="streamlover42@gmail.com"
          />
          <p class="text-gray-600 text-xs italic">
            Provide your email address.
          </p>

          <button
            class="rounded-sm flex flex-row justify-center items-center overflow-hidden mt-5 px-3 py-1 text-base text-white bg-gray-800 ${emailNotificationButtonState}"
            onClick=${registerForEmailButtonPressed}
          >
            Notify Me!
          </button>
        </div>
      </div>

      <div
        id="follow-loading-spinner-container"
        style="display: ${loaderStyle}"
      >
        <img id="follow-loading-spinner" src="/img/loading.gif" />
        <p class="text-gray-700 text-lg">Contacting your server.</p>
        <p class="text-gray-600 text-lg">Please wait...</p>
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
    <span id="fediverse-follow-button-container">
      <${ExternalActionButton} onClick=${onClick} action=${notifyAction} />
    </span>
  `;
}
