import { FC, ReactNode } from 'react';
import styles from './ActionButtonRow.module.scss';

export type ActionButtonRowProps = {
  children: ReactNode;
};

export const ActionButtonRow: FC<ActionButtonRowProps> = ({ children }) => (
  <div className={styles.row}>{children}</div>
);
