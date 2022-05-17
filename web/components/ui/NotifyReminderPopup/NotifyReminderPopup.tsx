import { Popover } from 'antd';
import { CloseOutlined } from '@ant-design/icons';
import React, { useState } from 'react';
import s from './NotifyReminderPopup.module.scss';

interface Props {
  visible: boolean;
  children: React.ReactNode[];
  notificationClicked: () => void;
  notificationClosed: () => void;
}

export default function NotifyReminderPopup(props: Props) {
  const { children, visible, notificationClicked, notificationClosed } = props;
  const [visiblePopup, setVisiblePopup] = useState(visible);

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
    <button type="button" onClick={popupClicked} className={s.contentbutton}>
      <button type="button" className={s.closebutton} onClick={popupClosed}>
        <CloseOutlined />
      </button>
      Click and never miss
      <br />
      future streams!
    </button>
  );
  return (
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
  );
}
