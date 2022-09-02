import { Button } from 'antd';
import { HeartFilled } from '@ant-design/icons';
import { useState } from 'react';
import { useRecoilValue } from 'recoil';
import Modal from '../ui/Modal/Modal';
import FollowModal from '../modals/Follow/FollowModal';
import s from './ActionButton/ActionButton.module.scss';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';
import { ClientConfig } from '../../interfaces/client-config.model';

export default function FollowButton(props: any) {
  const [showModal, setShowModal] = useState(false);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name, federation } = clientConfig;
  const { account } = federation;

  const buttonClicked = () => {
    setShowModal(true);
  };

  return (
    <>
      <Button
        {...props}
        type="primary"
        className={s.button}
        icon={<HeartFilled />}
        onClick={buttonClicked}
      >
        Follow
      </Button>
      <Modal
        title={`Follow ${name}`}
        visible={showModal}
        handleCancel={() => setShowModal(false)}
        width="550px"
        height="200px"
      >
        <FollowModal account={account} name={name} handleClose={() => setShowModal(false)} />
      </Modal>
    </>
  );
}
