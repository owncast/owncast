import { FC } from 'react';
import styles from './Logo.module.scss';

export type LogoProps = {
  src: string;
};

export const Logo: FC<LogoProps> = ({ src }) => (
  <div className={styles.root}>
    <div className={styles.container}>
      <img src={src} alt="Logo" className={styles.image} loading="lazy" />
    </div>
  </div>
);
