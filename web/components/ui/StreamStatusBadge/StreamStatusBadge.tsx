import type { FC } from 'react';
import { Badge } from 'antd';
import classNames from 'classnames';
import styles from './StreamStatusBadge.module.scss';

export type StreamStatusBadgeProps = {
  isOnline: boolean;
  className?: string;
};

export const StreamStatusBadge: FC<StreamStatusBadgeProps> = ({ isOnline, className }) => (
  <Badge
    status={isOnline ? 'success' : 'default'}
    text={isOnline ? 'LIVE' : 'OFFLINE'}
    className={classNames(styles.badge, { [styles.online]: isOnline }, className)}
  />
);
