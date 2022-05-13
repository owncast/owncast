import { Spin, Skeleton, Modal as AntModal } from 'antd';
import React, { ReactNode, useState } from 'react';
import s from './Modal.module.scss';

interface Props {
  title: string;
  url?: string;
  visible: boolean;
  handleOk?: () => void;
  handleCancel?: () => void;
  afterClose?: () => void;
  children?: ReactNode;
}

export default function Modal(props: Props) {
  const { title, url, visible, handleOk, handleCancel, afterClose, children } = props;
  const [loading, setLoading] = useState(!!url);

  const modalStyle = {
    padding: '0px',
    height: '80vh',
  };

  const iframe = url && (
    <iframe
      title={title}
      src={url}
      width="100%"
      height="100%"
      sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
      frameBorder="0"
      allowFullScreen
      onLoad={() => setLoading(false)}
    />
  );

  const iframeDisplayStyle = loading ? 'none' : 'inline';

  return (
    <AntModal
      title={title}
      visible={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      afterClose={afterClose}
      bodyStyle={modalStyle}
      width="70%"
      zIndex={9999}
      footer={null}
      centered
      destroyOnClose
    >
      <>
        {loading && (
          <Skeleton active={loading} style={{ padding: '10px' }} paragraph={{ rows: 10 }} />
        )}

        {iframe && <div style={{ display: iframeDisplayStyle }}>{iframe}</div>}
        {children && <div className={s.content}>{children}</div>}
        {loading && <Spin className={s.spinner} spinning={loading} size="large" />}
      </>
    </AntModal>
  );
}

Modal.defaultProps = {
  url: undefined,
  children: undefined,
  handleOk: undefined,
  handleCancel: undefined,
  afterClose: undefined,
};
