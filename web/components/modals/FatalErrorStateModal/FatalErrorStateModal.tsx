import { Modal } from 'antd';
import { FC } from 'react';

export type FatalErrorStateModalProps = {
  title: string;
  message: string;
};

export const FatalErrorStateModal: FC<FatalErrorStateModalProps> = ({ title, message }) => (
  <Modal
    title={title}
    visible
    footer={null}
    closable={false}
    keyboard={false}
    width={900}
    centered
    className="modal"
  >
    <p style={{ fontSize: '1.3rem' }}>{message}</p>
  </Modal>
);
