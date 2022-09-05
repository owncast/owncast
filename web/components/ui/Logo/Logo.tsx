import { Image } from 'antd';
import s from './Logo.module.scss';

interface Props {
  src: string;
}

export const Logo = ({ src }: Props) => (
  <div className={s.root}>
    <div className={s.container}>
      <Image src={src} alt="Logo" className={s.image} rootClassName={s.image} />
    </div>
  </div>
);
export default Logo;
