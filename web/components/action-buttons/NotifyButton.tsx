import { Button } from 'antd';
import { BellFilled } from '@ant-design/icons';
import { FC } from 'react';
import styles from './ActionButton/ActionButton.module.scss';

export type NotifyButtonProps = {
  text?: string;
  onClick?: () => void;
};

export const NotifyButton: FC<NotifyButtonProps> = ({ onClick, text }) => (
  <Button type="primary" className={`${styles.button}`} icon={<BellFilled />} onClick={onClick}>
    {text || 'Notify'}
  </Button>
);
