import { Divider, Button } from 'antd';
import { NotificationFilled } from '@ant-design/icons';
import { FC } from 'react';
import s from './OfflineBanner.module.scss';

export type OfflineBannerProps = {
  name: string;
  text: string;
};

export const OfflineBanner: FC<OfflineBannerProps> = ({ name, text }) => (
  <div className={s.outerContainer}>
    <div className={s.innerContainer}>
      <div className={s.header}>{name} is currently offline.</div>
      <Divider />
      <div>{text}</div>

      <div className={s.footer}>
        <Button type="primary" onClick={() => console.log('show notification modal')}>
          <NotificationFilled />
          Notify when Live
        </Button>
      </div>
    </div>
  </div>
);
export default OfflineBanner;
