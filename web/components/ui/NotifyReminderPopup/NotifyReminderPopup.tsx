import { Popover } from 'antd';
import { CloseOutlined } from '@ant-design/icons';
import React, { useState, useEffect, FC } from 'react';
import styles from './NotifyReminderPopup.module.scss';

export type NotifyReminderPopupProps = {
  visible: boolean;
  children: React.ReactNode;
  notificationClicked: () => void;
  notificationClosed: () => void;
};

export const NotifyReminderPopup: FC<NotifyReminderPopupProps> = ({
  children,
  visible,
  notificationClicked,
  notificationClosed,
}) => {
  const [visiblePopup, setVisiblePopup] = useState(visible);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setVisiblePopup(visible);
  }, [visible]);

  useEffect(() => {
    setMounted(true);
  }, []);

  const title = <div className={styles.title}>Stay updated!</div>;
  const popupStyle = {
    borderRadius: '5px',
    cursor: 'pointer',
    paddingTop: '10px',
    paddingRight: '10px',
    fontSize: '16px',
  };

  const popupClicked = e => {
    e.stopPropagation();
    notificationClicked();
  };

  const popupClosed = e => {
    e.stopPropagation();
    setVisiblePopup(false);
    notificationClosed();
  };

  const content = (
    <div onClick={popupClicked} onKeyDown={popupClicked} role="menuitem" tabIndex={0}>
      <button type="button" className={styles.closebutton} onClick={popupClosed}>
        <CloseOutlined />
      </button>
      <div className={styles.contentbutton}>
        Click and never miss
        <br />
        future streams!
      </div>
    </div>
  );

  return (
    mounted && (
      <Popover
        placement="topLeft"
        defaultVisible={visiblePopup}
        visible={visiblePopup}
        destroyTooltipOnHide
        title={title}
        content={content}
        overlayInnerStyle={popupStyle}
      >
        {children}
      </Popover>
    )
  );
};
export default NotifyReminderPopup;
