import { Button } from 'antd';
import { useState } from 'react';
import Modal from '../../ui/Modal/Modal';
import { ExternalAction } from '../../../interfaces/external-action';
import s from './ActionButton.module.scss';

interface Props {
  action: ExternalAction;
  primary?: boolean;
}
ActionButton.defaultProps = {
  primary: true,
};

export default function ActionButton({
  action: { url, title, description, icon, color, openExternally },
  primary = false,
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
      <Button
        type={primary ? 'primary' : 'default'}
        className={`${s.button}`}
        onClick={buttonClicked}
        style={{ backgroundColor: color }}
      >
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
