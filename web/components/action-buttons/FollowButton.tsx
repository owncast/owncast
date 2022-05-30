import { Button } from 'antd';
import { HeartFilled } from '@ant-design/icons';
import { useState } from 'react';
import Modal from '../ui/Modal/Modal';
import s from './ActionButton.module.scss';

export default function FollowButton() {
  const [showModal, setShowModal] = useState(false);

  const buttonClicked = () => {
    setShowModal(true);
  };

  return (
    <>
      <Button
        type="primary"
        className={`${s.button}`}
        icon={<HeartFilled />}
        onClick={buttonClicked}
      >
        Follow
      </Button>
      <Modal
        title="Follow <servername>"
        visible={showModal}
        handleCancel={() => setShowModal(false)}
      />
    </>
  );
}
