import { Button } from 'antd';
import { BellFilled } from '@ant-design/icons';
import s from './ActionButton/ActionButton.module.scss';

interface Props {
  onClick: () => void;
}

export default function NotifyButton({ onClick }: Props) {
  return (
    <Button type="primary" className={`${s.button}`} icon={<BellFilled />} onClick={onClick}>
      Notify
    </Button>
  );
}
