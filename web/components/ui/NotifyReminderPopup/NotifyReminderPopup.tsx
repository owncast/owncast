import React, { useState, useEffect, FC } from 'react';
import dynamic from 'next/dynamic';
import styles from './NotifyReminderPopup.module.scss';
import { Popover } from '../Popover/Popover';

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
  const [openPopup, setOpenPopup] = useState(open);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setOpenPopup(open);
  }, [open]);

  useEffect(() => {
    setMounted(true);
  }, []);

  const title = <div className={styles.title}>Stay updated!</div>;

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
        aria-label="Follow"
        className={styles.closebutton}
        onClick={popupClosed}
      >
        <CloseOutlined />
      </button>
      <div className={styles.contentbutton}>Click and never miss future streams!</div>
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
