import { Modal } from 'antd';
import { FC, useState } from 'react';

export type FatalErrorStateModalProps = {
  title: string;
  message: string;
};

export const FatalErrorStateModal: FC<FatalErrorStateModalProps> = ({ title, message }) => {
  const [isOpen, setOpen] = useState(true);

  const handleCancel = () => {
    setOpen(false);
  };

  const handleOk = () => {
    setOpen(false);
  };

  return (
    <Modal
      title={title}
      footer={null}
      closable={true}
      keyboard={false}
      width={900}
      centered
      open={isOpen}
      onCancel={handleCancel}
      onOk={handleOk}
      className="modal"
    >
      <p style={{ fontSize: '1.3rem' }}>{message}</p>
    </Modal>
  );
};
