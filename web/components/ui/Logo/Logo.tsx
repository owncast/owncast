import { Image } from 'antd';
import { FC } from 'react';
import s from './Logo.module.scss';

export type LogoProps = {
  src: string;
};

export const Logo: FC<LogoProps> = ({ src }) => (
  <div className={s.root}>
    <div className={s.container}>
      <Image src={src} alt="Logo" className={s.image} rootClassName={s.image} />
    </div>
  </div>
);
export default Logo;
