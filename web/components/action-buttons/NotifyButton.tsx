import { Button } from 'antd';
import { NotificationFilled } from '@ant-design/icons';
import { useState } from 'react';
import Modal from '../ui/Modal/Modal';
import s from './ActionButton.module.scss';
import BrowserNotifyModal from '../modals/BrowserNotify/BrowserNotifyModal';

export default function NotifyButton() {
  const [showModal, setShowModal] = useState(false);

  const buttonClicked = () => {
    setShowModal(true);
  };

  return (
    <>
      <Button
        type="primary"
        className={`${s.button}`}
        icon={<NotificationFilled />}
        onClick={buttonClicked}
      >
        Notify
      </Button>
      <Modal title="Notify" visible={showModal} handleCancel={() => setShowModal(false)}>
        <BrowserNotifyModal />
      </Modal>
    </>
  );
}
