import { Button, ButtonProps } from 'antd';
import { HeartFilled } from '@ant-design/icons';

import { FC } from 'react';
import styles from './ActionButton/ActionButton.module.scss';

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
