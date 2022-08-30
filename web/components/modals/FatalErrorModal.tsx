import { Modal } from 'antd';

interface Props {
  title: string;
  message: string;
}

export default function FatalErrorStateModal(props: Props) {
  const { title, message } = props;

  return (
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
}
