import { Button } from 'antd';
import { NotificationFilled } from '@ant-design/icons';
import s from './ActionButton.module.scss';

interface Props {
  onClick: () => void;
}

export default function NotifyButton({ onClick }: Props) {
  return (
    <Button
      type="primary"
      className={`${s.button}`}
      icon={<NotificationFilled />}
      onClick={onClick}
    >
      Notify
    </Button>
  );
}
