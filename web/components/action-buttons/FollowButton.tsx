import { Button } from 'antd';
import { HeartFilled } from '@ant-design/icons';
import { useState } from 'react';
import { useRecoilValue } from 'recoil';
import Modal from '../ui/Modal/Modal';
import FollowModal from '../modals/Follow/FollowModal';
import s from './ActionButton.module.scss';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';
import { ClientConfig } from '../../interfaces/client-config.model';

export default function FollowButton() {
  const [showModal, setShowModal] = useState(false);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name } = clientConfig;

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
      <Modal title={`Follow ${name}`} visible={showModal} handleCancel={() => setShowModal(false)}>
        <FollowModal handleClose={() => setShowModal(false)} />
      </Modal>
    </>
  );
}
