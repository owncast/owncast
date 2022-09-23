import { FC, ReactNode } from 'react';
import styles from './ButtonRow.module.scss';

export type ButtonRowProps = {
  children: ReactNode;
};

/**
 * A horizontal row to wrap multiple {@link ../Button}.
 */
export const ButtonRow: FC<ButtonRowProps> = ({ children }) => (
  <div className={styles.row}>{children}</div>
);
