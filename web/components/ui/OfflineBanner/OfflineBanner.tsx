import { Divider, Button } from 'antd';
import { NotificationFilled } from '@ant-design/icons';

import s from './OfflineBanner.module.scss';

interface Props {
  name: string;
  text: string;
}

export default function OfflineBanner({ name, text }: Props) {
  const handleShowNotificationModal = () => {
    console.log('show notification modal');
  };

  return (
    <div className={s.outerContainer}>
      <div className={s.innerContainer}>
        <div className={s.header}>{name} is currently offline.</div>
        <Divider />
        <div>{text}</div>

        <div className={s.footer}>
          <Button onClick={handleShowNotificationModal}>
            <NotificationFilled />
            Notify when Live
          </Button>
        </div>
      </div>
    </div>
  );
}
