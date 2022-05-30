import { Button } from 'antd';
import { useState } from 'react';
import Modal from '../ui/Modal/Modal';
import { ExternalAction } from '../../interfaces/external-action';
import s from './ActionButton.module.scss';

interface Props {
  action: ExternalAction;
}

export default function ActionButton({
  action: { url, title, description, icon, openExternally },
}: Props) {
  const [showModal, setShowModal] = useState(false);

  const buttonClicked = () => {
    if (openExternally) {
      window.open(url, '_blank');
    } else {
      setShowModal(true);
    }
  };

  return (
    <>
      <Button type="primary" className={`${s.button}`} onClick={buttonClicked}>
        <img src={icon} className={`${s.icon}`} alt={description} />
        {title}
      </Button>
      <Modal
        title={description || title}
        url={url}
        visible={showModal}
        height="80vh"
        handleCancel={() => setShowModal(false)}
      />
    </>
  );
}
