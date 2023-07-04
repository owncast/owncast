import { Spin, Modal as AntModal } from 'antd';
import React, { FC, ReactNode, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { ComponentError } from '../ComponentError/ComponentError';
import styles from './Modal.module.scss';

export type ModalProps = {
  title: string;
  url?: string;
  open: boolean;
  handleOk?: () => void;
  handleCancel?: () => void;
  afterClose?: () => void;
  children?: ReactNode;
  height?: string;
  width?: string;
};

export const Modal: FC<ModalProps> = ({
  title,
  url,
  open,
  handleOk,
  handleCancel,
  afterClose,
  height,
  width,
  children,
}) => {
  const [loading, setLoading] = useState(!!url);

  let defaultHeight = '100%';
  let defaultWidth = '520px';
  if (url) {
    defaultHeight = '70vh';
    defaultWidth = '900px';
  }

  const modalContentBodyStyle = {
    padding: '0px',
    minHeight: height,
    height: height ?? defaultHeight,
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
      style={{ display: 'block' }}
      // eslint-disable-next-line react/no-unknown-property
      onLoad={() => setLoading(false)}
    />
  );

  const iframeDisplayStyle = loading ? 'none' : 'inline';

  return (
    <AntModal
      title={title}
      open={open}
      onOk={handleOk}
      onCancel={handleCancel}
      afterClose={afterClose}
      bodyStyle={modalContentBodyStyle}
      width={width ?? defaultWidth}
      zIndex={999}
      footer={null}
      centered
      destroyOnClose
      className={styles.modal}
    >
      <ErrorBoundary
        // eslint-disable-next-line react/no-unstable-nested-components
        fallbackRender={({ error, resetErrorBoundary }) => (
          <ComponentError
            componentName="Modal"
            message={error.message}
            retryFunction={resetErrorBoundary}
          />
        )}
      >
        <div id="modal-container" style={{ height: '100%' }}>
          {iframe && <div style={{ display: iframeDisplayStyle }}>{iframe}</div>}
          {children && <div className={styles.content}>{children}</div>}
          {loading && (
            <Spin className={styles.spinner} spinning={loading} size="large" tip={title} />
          )}
        </div>
      </ErrorBoundary>
    </AntModal>
  );
};

Modal.defaultProps = {
  url: undefined,
  children: undefined,
  handleOk: undefined,
  handleCancel: undefined,
  afterClose: undefined,
};
