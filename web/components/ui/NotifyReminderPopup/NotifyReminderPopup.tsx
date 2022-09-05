import { Popover } from 'antd';
import { CloseOutlined } from '@ant-design/icons';
import React, { useState, useEffect } from 'react';
import s from './NotifyReminderPopup.module.scss';

interface Props {
  visible: boolean;
  children: React.ReactNode;
  notificationClicked: () => void;
  notificationClosed: () => void;
}

export const NotifyReminderPopup = (props: Props) => {
  const { children, visible, notificationClicked, notificationClosed } = props;
  const [visiblePopup, setVisiblePopup] = useState(visible);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setVisiblePopup(visible);
  }, [visible]);

  useEffect(() => {
    setMounted(true);
  }, []);

  const title = <div className={s.title}>Stay updated!</div>;
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
      <button type="button" className={s.closebutton} onClick={popupClosed}>
        <CloseOutlined />
      </button>
      <div className={s.contentbutton}>
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
