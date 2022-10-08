import { FC, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { Button, ButtonProps } from 'antd';
import { HeartFilled } from '@ant-design/icons';
import { Modal } from '~/components/ui/Modal/Modal';
import { FollowModal } from '~/components/modals/FollowModal/FollowModal';
import { clientConfigStateAtom } from '~/components/stores/ClientConfigStore';
import { ClientConfig } from '~/interfaces/client-config.model';
import styles from './ActionButton/ActionButton.module.scss';

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
