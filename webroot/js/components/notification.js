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
  const [phoneNotificationsButtonEnabled, setPhoneNotificationsButtonEnabled] =
    useState(false);
  const [phoneNumber, setPhoneNumber] = useState(null);
  const phoneNotificationButtonState = phoneNotificationsButtonEnabled
    ? ''
    : 'cursor-not-allowed opacity-50';

  const [browserRegistrationComplete, setBrowserRegistrationComplete] =
    useState(false);

  const browserNotificationsButtonState =
    Notification.permission === 'default'
      ? ''
      : 'cursor-not-allowed opacity-50';

  const { browser, textMessages } = notifications;
  const { publicKey } = browser;

  let browserPushEnabled = browser.enabled;
  let textMessagesEnabled = textMessages.enabled;

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

  function onPhoneInput(e) {
    const { value } = e.target;

    // Since phone number validation is difficult let's just make sure
    // it has a reasonable number of digits and doesn't include letters.

    // Strip non-numeric characters.
    const validated = value.replace(/[^\d]/g, '');

    let valid = false;
    if (validated.length >= 11 && validated.length <= 16) {
      valid = true;
    }

    setPhoneNumber(validated);
    setPhoneNotificationsButtonEnabled(valid);
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
  if (browserPushEnabled && !textMessagesEnabled) {
    gridColumns = 1;
  } else if (!browserPushEnabled && textMessagesEnabled) {
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
        <div style=${{ display: textMessagesEnabled ? 'block' : 'none' }}>
          <h2 class="text-indigo-600 text-3xl font-semibold">Phone</h2>
          <p>
            No apps or 3rd party services required. Get notified directly when
            ${''} ${streamName} goes live.
          </p>
          <p class="mt-4">
            Respond to notifications with <strong>stop</strong> to stop
            receiving them.
          </p>
          <input
            class="border bg-white rounded w-full py-2 px-3 mb-2 mt-10 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
            type="tel"
            maxlength="16"
            value=${phoneNumber}
            onInput=${onPhoneInput}
            placeholder="Your phone number 14023133798"
          />
          <p class="text-gray-600 text-xs italic">
            Provide your phone number with the international country code.
          </p>

          <button
            class="rounded-sm flex flex-row justify-center items-center overflow-hidden mt-5 px-3 py-1 text-base text-white bg-gray-800 ${phoneNotificationButtonState}"
            onClick=${registerForBrowserPushButtonPressed}
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
