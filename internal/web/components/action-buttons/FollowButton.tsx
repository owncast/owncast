import { Button, ButtonProps } from 'antd';

import { FC } from 'react';
import dynamic from 'next/dynamic';
import styles from './ActionButton/ActionButton.module.scss';

// Lazy loaded components

const HeartFilled = dynamic(() => import('@ant-design/icons/HeartFilled'), {
  ssr: false,
});

export type FollowButtonProps = ButtonProps & {
  onClick?: () => void;
  props?: ButtonProps;
};

export const FollowButton: FC<FollowButtonProps> = ({ onClick, props }) => (
  <Button
    {...props}
    type="primary"
    className={styles.button}
    icon={<HeartFilled />}
    onClick={onClick}
    id="follow-button"
  >
    Follow
  </Button>
);
