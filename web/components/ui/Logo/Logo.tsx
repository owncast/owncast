import { Image } from 'antd';
import { FC } from 'react';
import styles from './Logo.module.scss';

export type LogoProps = {
  src: string;
};

export const Logo: FC<LogoProps> = ({ src }) => (
  <div className={styles.root}>
    <div className={styles.container}>
      <Image src={src} alt="Logo" className={styles.image} rootClassName={styles.image} />
    </div>
  </div>
);
export default Logo;
