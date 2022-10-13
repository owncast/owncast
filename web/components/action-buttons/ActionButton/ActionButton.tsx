import { Button } from 'antd';
import { FC, useState } from 'react';
import { Modal } from '../../ui/Modal/Modal';
import { ExternalAction } from '../../../interfaces/external-action';
import styles from './ActionButton.module.scss';

export type ActionButtonProps = {
  action: ExternalAction;
  primary?: boolean;
};

export const ActionButton: FC<ActionButtonProps> = ({
  action: { url, title, description, icon, color, openExternally },
  primary = true,
}) => {
  const [showModal, setShowModal] = useState(false);

  const onButtonClicked = () => {
    if (openExternally) {
      window.open(url, '_blank');
    } else {
      setShowModal(true);
    }
  };

  return (
    <>
      <Button
        type={primary ? 'primary' : 'default'}
        className={`${styles.button}`}
        onClick={onButtonClicked}
        style={{ backgroundColor: color }}
      >
        {icon && <img src={icon} className={`${styles.icon}`} alt={description} />}
        {title}
      </Button>
      <Modal
        title={description || title}
        url={url}
        open={showModal}
        height="80vh"
        handleCancel={() => setShowModal(false)}
      />
    </>
  );
};
