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
      <style global jsx>{`
        .modal .ant-modal-content,
        .modal .ant-modal-header {
          background-color: var(--warning-color);
        }
      `}</style>
      <p>{message}</p>
    </Modal>
  );
}
