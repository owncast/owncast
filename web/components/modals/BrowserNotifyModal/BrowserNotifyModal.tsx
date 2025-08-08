import { Row, Spin, Typography, Button, Alert } from 'antd';
import React, { FC, useState } from 'react';
import UploadOutlined from '@ant-design/icons/lib/icons/UploadOutlined';
import PlusSquareOutlined from '@ant-design/icons/lib/icons/PlusSquareOutlined';
import { useRecoilValue } from 'recoil';
import { ErrorBoundary } from 'react-error-boundary';
import { useTranslation } from 'next-export-i18n';
import { accessTokenAtom, clientConfigStateAtom } from '../../stores/ClientConfigStore';
import {
  registerWebPushNotifications,
  saveNotificationRegistration,
} from '../../../services/notifications-service';
import styles from './BrowserNotifyModal.module.scss';
import { ComponentError } from '../../ui/ComponentError/ComponentError';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

import { isMobileSafariHomeScreenApp, isMobileSafariIos } from '../../../utils/helpers';
import { arePushNotificationSupported } from '../../../utils/browserPushNotifications';

const { Title } = Typography;

const NotificationsNotSupported = () => (
  <div>
    <Translation
      translationKey={Localization.Frontend.BrowserNotifyModal.unsupported}
      defaultText="Browser notifications are not supported in your browser."
    />
  </div>
);

const NotificationsNotSupportedLocal = () => (
  <div>
    <Translation
      translationKey={Localization.Frontend.BrowserNotifyModal.unsupportedLocal}
      defaultText="Browser notifications are not supported for local servers."
    />
  </div>
);

const MobileSafariInstructions = () => (
  <div>
    <Title level={3}>
      <Translation
        translationKey={Localization.Frontend.BrowserNotifyModal.iosTitle}
        defaultText="Get notified on iOS"
      />
    </Title>
    <Translation
      translationKey={Localization.Frontend.BrowserNotifyModal.iosDescription}
      defaultText="It takes a couple extra steps to make sure you get notified when your favorite streams go live."
    />
    <ol>
      <li>
        Tap the{' '}
        <strong>
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.iosShareButton}
            defaultText="share"
          />
        </strong>{' '}
        button <UploadOutlined /> in Safari.
      </li>
      <li>
        Scroll down and tap{' '}
        <strong>
          &ldquo;
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.iosAddToHomeScreen}
            defaultText="Add to Home Screen"
          />
          &rdquo;
        </strong>{' '}
        <PlusSquareOutlined />.
      </li>
      <li>
        Tap{' '}
        <strong>
          &ldquo;
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.iosAddButton}
            defaultText="Add"
          />
          &rdquo;
        </strong>
        .
      </li>
      <li>
        <Translation
          translationKey={Localization.Frontend.BrowserNotifyModal.iosNameAndTap}
          defaultText="Give this link a name and tap the new icon on your home screen"
        />
      </li>

      <li>
        <Translation
          translationKey={Localization.Frontend.BrowserNotifyModal.iosComeBack}
          defaultText="Come back to this screen and enable notifications."
        />
      </li>
      <li>
        Tap{' '}
        <strong>
          &ldquo;
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.iosAllowPrompt}
            defaultText="Allow"
          />
          &rdquo;
        </strong>{' '}
        when prompted.
      </li>
    </ol>
  </div>
);

export type PermissionPopupPreviewProps = {
  start: () => void;
};

const PermissionPopupPreview: FC<PermissionPopupPreviewProps> = ({ start }) => (
  <div id="browser-push-preview-box" className={styles.pushPreview}>
    <div className={styles.inner}>
      <div className={styles.title}>
        <Translation
          translationKey={Localization.Frontend.BrowserNotifyModal.permissionWantsTo}
          defaultText="{{hostname}} wants to"
          vars={{ hostname: window.location.hostname }}
        />
      </div>
      <div className={styles.permissionLine}>
        <svg
          className={styles.bell}
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
        <span className={styles.showNotificationsText}>
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.showNotifications}
            defaultText="Show notifications"
          />
        </span>
      </div>
      <div className={styles.buttonRow}>
        <Button
          type="primary"
          onClick={() => {
            start();
          }}
        >
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.allowButton}
            defaultText="Allow"
          />
        </Button>
        <button type="button" className={styles.disabled}>
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.blockButton}
            defaultText="Block"
          />
        </button>
      </div>
    </div>
  </div>
);

const NotificationsEnabled = () => (
  <div>
    <Title level={2}>
      <Translation
        translationKey={Localization.Frontend.BrowserNotifyModal.enabledTitle}
        defaultText="Notifications are enabled"
      />
    </Title>
    <Translation
      translationKey={Localization.Frontend.BrowserNotifyModal.enabledDescription}
      defaultText="To disable push notifications from {{hostname}} access your browser permissions for this site and turn off notifications. <a href='https://owncast.online/docs/notifications'>Learn more.</a>"
      vars={{ hostname: window.location.hostname.toString() }}
    />
  </div>
);

const NotificationsDenied = () => (
  <div>
    <Title level={2}>
      <Translation
        translationKey={Localization.Frontend.BrowserNotifyModal.deniedTitle}
        defaultText="Notifications are blocked on your device"
      />
    </Title>
    <Translation
      translationKey={Localization.Frontend.BrowserNotifyModal.deniedDescription}
      defaultText="To enable push notifications from {{hostname}} access your browser permissions for this site and turn on notifications. Then reload this page to apply your updated settings on this site. <a href='https://owncast.online/docs/notifications'>Learn more.</a>"
      vars={{ hostname: window.location.hostname.toString() }}
    />
  </div>
);

export const BrowserNotifyModal = () => {
  const { t } = useTranslation();
  const [error, setError] = useState<string>(null);
  const accessToken = useRecoilValue(accessTokenAtom);
  const config = useRecoilValue(clientConfigStateAtom);
  const [browserPushPermissionsPending, setBrowserPushPermissionsPending] =
    useState<boolean>(false);
  const notificationsPermitted =
    arePushNotificationSupported() && Notification.permission === 'granted';
  const notificationsDenied =
    arePushNotificationSupported() && Notification.permission === 'denied';

  const { notifications } = config;
  const { browser } = notifications;
  const { publicKey } = browser;

  const browserPushSupported =
    browser.enabled && (arePushNotificationSupported() || isMobileSafariHomeScreenApp());

  // If notification permissions are granted, show user info how to disable them
  if (notificationsPermitted) {
    return <NotificationsEnabled />;
  }

  if (notificationsDenied) {
    return <NotificationsDenied />;
  }

  if (isMobileSafariIos() && !isMobileSafariHomeScreenApp()) {
    return <MobileSafariInstructions />;
  }

  const startBrowserPushRegistration = async () => {
    // If notification permissions are already denied or granted, don't do anything.
    if (arePushNotificationSupported() && Notification.permission !== 'default') {
      return;
    }

    setBrowserPushPermissionsPending(true);
    try {
      const subscription = await registerWebPushNotifications(publicKey);
      saveNotificationRegistration('BROWSER_PUSH_NOTIFICATION', subscription, accessToken);
      setError(null);
    } catch (e) {
      setError(
        t(Localization.Frontend.BrowserNotifyModal.errorMessage, {
          message: e.message,
        }),
      );
    }
    setBrowserPushPermissionsPending(false);
  };

  if (window.location.hostname === 'localhost') {
    return <NotificationsNotSupportedLocal />;
  }

  if (!browserPushSupported) {
    return <NotificationsNotSupported />;
  }

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error: e, resetErrorBoundary }) => (
        <ComponentError
          componentName="BrowserNotifyModal"
          message={e.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <Spin spinning={browserPushPermissionsPending}>
        <Row className={styles.description}>
          <Translation
            translationKey={Localization.Frontend.BrowserNotifyModal.mainDescription}
            defaultText="Get notified right in the browser each time this stream goes live."
          />
          <span>
            <a href="https://owncast.online/docs/notifications/#browser-notifications">
              <Translation
                translationKey={Localization.Frontend.BrowserNotifyModal.learnMore}
                defaultText="Learn more"
              />
            </a>
            &nbsp; about Owncast browser notifications.
          </span>
        </Row>
        <Row>
          {error && (
            <Alert
              message={
                <Translation
                  translationKey={Localization.Frontend.BrowserNotifyModal.errorTitle}
                  defaultText="Browser Notification Error"
                />
              }
              description={error}
              type="error"
              closable
              className={styles.errorAlert}
              onClose={() => setError(null)}
            />
          )}
        </Row>
        <PermissionPopupPreview start={() => startBrowserPushRegistration()} />
      </Spin>
    </ErrorBoundary>
  );
};
