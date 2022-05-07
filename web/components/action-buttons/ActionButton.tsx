import { Button } from 'antd';
import { ExternalAction } from '../interfaces/external-action.interface';
import s from './ActionButton.module.scss';

interface Props {
  action: ExternalAction;
}

export default function ActionButton(props: Props) {
  const { action } = props;
  const { url, title, description, icon, color, openExternally } = action;

  return (
    <Button type="primary" className={`${s.button}`} style={{ backgroundColor: color }}>
      <img src={icon} className={`${s.icon}`} alt={description} />
      {title}
    </Button>
  );
}
