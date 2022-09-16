import { FC, ReactNode } from 'react';
import styles from './ButtonRow.module.scss';

export type ButtonRowProps = {
  children: ReactNode;
};

export const ButtonRow: FC<ButtonRowProps> = ({ children }) => (
  <div className={styles.row}>{children}</div>
);
