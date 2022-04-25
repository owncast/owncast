import { h } from '/js/web_modules/preact.js';
import { useState, useEffect } from '/js/web_modules/preact/hooks.js';

import htm from '/js/web_modules/htm.js';
import { ExternalActionButton } from './external-action-modal.js';
import {
  registerWebPushNotifications,
  isPushNotificationSupported,
} from '../notification/registerWeb.js';
import {
  URL_REGISTER_NOTIFICATION,
  URL_REGISTER_EMAIL_NOTIFICATION,
  HAS_DISPLAYED_NOTIFICATION_MODAL_KEY,
  USER_VISIT_COUNT_KEY,
  USER_DISMISSED_ANNOYING_NOTIFICATION_POPUP_KEY,
} from '../utils/constants.js';
import { setLocalStorage, getLocalStorage } from '../utils/helpers.js';

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
  const [browserPushPermissionsPending, setBrowserPushPermissionsPending] =
    useState(false);

  const { browser } = notifications;
  const { publicKey } = browser;

  const browserPushEnabled = browser.enabled && isPushNotificationSupported();
  let emailEnabled = false;

  // Store that the user has opened the notifications modal at least once
  // so we don't ever need to remind them to do it again.
  useEffect(() => {
    setLocalStorage(HAS_DISPLAYED_NOTIFICATION_MODAL_KEY, true);
  }, []);

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

  async function startBrowserPushRegistration() {
    // If it's already denied or granted, don't do anything.
    if (Notification.permission !== 'default') {
      return;
    }

    setBrowserPushPermissionsPending(true);
    try {
      const subscription = await registerWebPushNotifications(publicKey);
      saveNotificationRegistration('BROWSER_PUSH_NOTIFICATION', subscription);
      setError(null);
    } catch (e) {
      setError(
        `Error registering for live notifications: ${e.message}. Make sure you're not inside a private browser environment or have previously disabled notifications for this stream.`
      );
    }
    setBrowserPushPermissionsPending(false);
  }

  async function handlePushToggleChange() {
    // Nothing can be done if they already denied access.
    if (Notification.permission === 'denied') {
      return;
    }

    if (!pushEnabled) {
      startBrowserPushRegistration();
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
    let pushNotificationButtonText = html`<span id="push-notification-arrow"
        >←</span
      >
      CLICK TO ENABLE`;
    if (browserPushPermissionsPending) {
      pushNotificationButtonText = '↑ ACCEPT THE BROWSER PERMISSIONS';
    } else if (Notification.permission === 'granted') {
      pushNotificationButtonText = 'ENABLED';
    } else if (Notification.permission === 'denied') {
      pushNotificationButtonText = 'DENIED. PLEASE FIX BROWSER PERMISSIONS.';
    }
    return pushNotificationButtonText;
  }

  const pushEnabled = Notification.permission === 'granted';

  return html`
    <div class="bg-gray-100 bg-center bg-no-repeat p-6">
      <div
        style=${{ display: emailEnabled ? 'grid' : 'none' }}
        class="grid grid-cols-2 gap-10 px-5 py-8"
      >
        <div>
          <h2 class="text-slate-600 text-2xl mb-2 font-semibold">
            Email Notifications
          </h2>

          <h2>
            Get notified directly to your email when this stream goes live.
          </h2>
        </div>
        <div>
          <div class="font-semibold">Enter your email address:</div>
          <input
            class="border bg-white rounded-l w-8/12 mt-2 mb-1 mr-1 py-2 px-3 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
            value=${emailAddress}
            onInput=${onEmailInput}
            placeholder="streamlover42@gmail.com"
          />
          <button
            class="rounded-sm inline px-3 py-2 text-base text-white bg-indigo-700 ${emailNotificationButtonState}"
            onClick=${registerForEmailButtonPressed}
          >
            Save
          </button>
          <div class="text-sm mt-3 text-gray-700">
            Stop receiving emails any time by clicking the unsubscribe link in
            the email. <a href="">Learn more.</a>
          </div>
        </div>
      </div>
      <hr
        style=${{ display: pushEnabled && emailEnabled ? 'block' : 'none' }}
      />

      <div
        class="grid grid-cols-2 gap-10 px-5 py-8"
        style=${{ display: browserPushEnabled ? 'grid' : 'none' }}
      >
        <div>
          <div>
            <div
              class="text-sm border-2 p-4 border-red-300"
              style=${{
                display:
                  Notification.permission === 'denied' ? 'block' : 'none',
              }}
            >
              Browser notification permissions were denied. Please visit your
              browser settings to re-enable in order to get notifications.
            </div>
            <div
              class="form-check form-switch"
              style=${{
                display:
                  Notification.permission === 'denied' ? 'none' : 'block',
              }}
            >
              <div
                class="relative inline-block w-10 align-middle select-none transition duration-200 ease-in"
              >
                <input
                  checked=${pushEnabled || browserPushPermissionsPending}
                  disabled=${pushEnabled}
                  type="checkbox"
                  name="toggle"
                  id="toggle"
                  onchange=${handlePushToggleChange}
                  class="toggle-checkbox absolute block w-8 h-8 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                  for="toggle"
                  style=${{ width: '50px' }}
                  class="toggle-label block overflow-hidden h-8 rounded-full bg-gray-300 cursor-pointer"
                ></label>
              </div>
              <div class="ml-8 text-xs inline-block text-gray-700">
                ${getBrowserPushButtonText()}
              </div>
            </div>
          </div>
          <h2 class="text-slate-600 text-2xl mt-4 mb-2 font-semibold">
            Browser Notifications
          </h2>
          <h2>
            Get notified right in the browser each time this stream goes live.
          </h2>
        </div>
        <div>
          <div
            class="text-sm mt-3"
            style=${{ display: !pushEnabled ? 'none' : 'block' }}
          >
            To disable push notifications from ${window.location.hostname}
            ${' '} access your browser permissions for this site and turn off
            notifications.
            <div style=${{ 'margin-top': '5px' }}>
              <a href="https://owncast.online/docs/notifications"
                >Learn more.</a
              >
            </div>
          </div>
          <div
            id="browser-push-preview-box"
            class="w-full  bg-white p-4 m-2 mt-4"
            style=${{ display: pushEnabled ? 'none' : 'block' }}
          >
            <div class="text-lg text-gray-700 ml-2 my-2">
              ${window.location.toString()} wants to
            </div>
            <div class="text-sm text-gray-700 my-2">
              <svg
                class="mr-3"
                style=${{ display: 'inline-block' }}
                width="16"
                height="16"
                viewBox="0 0 16 16"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M14 12.3333V13H2V12.3333L3.33333 11V7C3.33333 4.93333 4.68667 3.11333 6.66667 2.52667C6.66667 2.46 6.66667 2.4 6.66667 2.33333C6.66667 1.97971 6.80714 1.64057 7.05719 1.39052C7.30724 1.14048 7.64638 1 8 1C8.35362 1 8.69276 1.14048 8.94281 1.39052C9.19286 1.64057 9.33333 1.97971 9.33333 2.33333C9.33333 2.4 9.33333 2.46 9.33333 2.52667C11.3133 3.11333 12.6667 4.93333 12.6667 7V11L14 12.3333ZM9.33333 13.6667C9.33333 14.0203 9.19286 14.3594 8.94281 14.6095C8.69276 14.8595 8.35362 15 8 15C7.64638 15 7.30724 14.8595 7.05719 14.6095C6.80714 14.3594 6.66667 14.0203 6.66667 13.6667"
                  fill="#676670"
                />
              </svg>

              Show notifications
            </div>
            <div class="flex flex-row justify-end">
              <button
                class="bg-blue-500 py-1 px-4 mr-4 rounded-sm text-white"
                onClick=${startBrowserPushRegistration}
              >
                Allow
              </button>
              <button
                class="bg-slate-200 py-1 px-4 rounded-sm text-gray-500 cursor-not-allowed"
                style=${{
                  'outline-width': 1,
                  'outline-color': '#e2e8f0',
                  'outline-style': 'solid',
                }}
              >
                Block
              </button>
            </div>
          </div>
          <p
            class="text-gray-700 text-sm mt-6"
            style=${{ display: pushEnabled ? 'none' : 'block' }}
          >
            You'll need to allow your browser to receive notifications from
            ${' '} ${streamName}, first.
          </p>
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

export function NotifyButton({ serverName, onClick }) {
  const hasDisplayedNotificationModal = getLocalStorage(
    HAS_DISPLAYED_NOTIFICATION_MODAL_KEY
  );

  const hasPreviouslyDismissedAnnoyingPopup = getLocalStorage(
    USER_DISMISSED_ANNOYING_NOTIFICATION_POPUP_KEY
  );

  let visits = parseInt(getLocalStorage(USER_VISIT_COUNT_KEY));
  if (isNaN(visits)) {
    visits = 0;
  }

  // Only show the annoying popup if the user has never opened the notification
  // modal previously _and_ they've visited more than 3 times.
  const [showPopup, setShowPopup] = useState(
    !hasPreviouslyDismissedAnnoyingPopup &&
      !hasDisplayedNotificationModal &&
      visits > 3
  );

  const notifyAction = {
    color: 'rgba(219, 223, 231, 1)',
    description: `Never miss a stream! Get notified when ${serverName} goes live.`,
    icon: '/img/notification-bell.svg',
    openExternally: false,
  };

  const buttonClicked = (e) => {
    onClick(e);
    setShowPopup(false);
  };

  const notifyPopupDismissedClicked = (e) => {
    e.stopPropagation();
    setShowPopup(false);
    setLocalStorage(USER_DISMISSED_ANNOYING_NOTIFICATION_POPUP_KEY, true);
  };

  return html`
    <span id="notify-button-container" class="relative">
      <div
        id="notify-button-popup"
        class="text-gray-200 p-4 rounded-md cursor-pointer"
        style=${{ display: showPopup ? 'block' : 'none' }}
        onClick=${buttonClicked}
      >
        <div class="flex justify-between items-center mb-2">
          <div class="font-bold">Stay updated!</div>
          <button
            class="popout-close-button rounded-md p-1 color-gray-500"
            onClick=${notifyPopupDismissedClicked}
          >
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </button>
        </div>
        <div>Click and never miss future streams!</div>
      </div>
      <${ExternalActionButton}
        onClick=${buttonClicked}
        action=${notifyAction}
      />
    </span>
  `;
}
