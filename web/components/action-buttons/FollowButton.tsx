import { Button, ButtonProps } from 'antd';
import { HeartFilled } from '@ant-design/icons';
import { FC, useState } from 'react';
import { useRecoilValue } from 'recoil';
import Modal from '../ui/Modal/Modal';
import FollowModal from '../modals/FollowModal/FollowModal';
import styles from './ActionButton/ActionButton.module.scss';
import { clientConfigStateAtom } from '../stores/ClientConfigStore';
import { ClientConfig } from '../../interfaces/client-config.model';

export type FollowButtonProps = ButtonProps;

export const FollowButton: FC<FollowButtonProps> = props => {
  const [showModal, setShowModal] = useState(false);
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name, federation } = clientConfig;
  const { account } = federation;

  return (
    <>
      <Button
        {...props}
        type="primary"
        className={styles.button}
        icon={<HeartFilled />}
        onClick={() => setShowModal(true)}
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
};
