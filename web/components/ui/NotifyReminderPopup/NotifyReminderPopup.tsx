import React, { useState, useEffect, FC } from 'react';
import dynamic from 'next/dynamic';
import styles from './NotifyReminderPopup.module.scss';
import { Popover } from '../Popover/Popover';
import { Translation } from '../Translation/Translation';
import { Localization } from '../../../types/localization';
import { useTranslation } from 'next-export-i18n';

// Lazy loaded components

const CloseOutlined = dynamic(() => import('@ant-design/icons/CloseOutlined'), {
  ssr: false,
});

export type NotifyReminderPopupProps = {
  open: boolean;
  children: React.ReactNode;
  notificationClicked: () => void;
  notificationClosed: () => void;
};

export const NotifyReminderPopup: FC<NotifyReminderPopupProps> = ({
  children,
  open,
  notificationClicked,
  notificationClosed,
}) => {
  const { t } = useTranslation();
  const [openPopup, setOpenPopup] = useState(open);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setOpenPopup(open);
  }, [open]);

  useEffect(() => {
    setMounted(true);
  }, []);

  const title = (
    <div className={styles.title}>
      <Translation
        translationKey={Localization.Frontend.NotifyReminderPopup.title}
        defaultText="Stay updated!"
      />
    </div>
  );

  const popupClicked = e => {
    e.stopPropagation();
    notificationClicked();
  };

  const popupClosed = e => {
    e.stopPropagation();
    setOpenPopup(false);
    notificationClosed();
  };

  const content = (
    <div onClick={popupClicked} onKeyDown={popupClicked} role="menuitem" tabIndex={0}>
      <button
        type="button"
        aria-label={t(Localization.Frontend.FollowAndNotify.closeReminder)}
        className={styles.closebutton}
        onClick={popupClosed}
      >
        <CloseOutlined />
      </button>
      <div className={styles.contentbutton}>
        <Translation
          translationKey={Localization.Frontend.NotifyReminderPopup.description}
          defaultText="Click and never miss future streams!"
        />
      </div>
    </div>
  );

  return (
    mounted && (
      <Popover open={openPopup} title={title} content={content}>
        {children}
      </Popover>
    )
  );
};
