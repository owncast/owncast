import { Button } from 'antd';
import { FC } from 'react';
import dynamic from 'next/dynamic';
import styles from './ActionButton/ActionButton.module.scss';

// Lazy loaded components

const BellFilled = dynamic(() => import('@ant-design/icons/BellFilled'), {
  ssr: false,
});

export type NotifyButtonProps = {
  text?: string;
  onClick?: () => void;
};

export const NotifyButton: FC<NotifyButtonProps> = ({ onClick, text }) => (
  <Button
    type="primary"
    className={styles.button}
    icon={<BellFilled />}
    onClick={onClick}
    id="notify-button"
  >
    {text || 'Notify'}
  </Button>
);
