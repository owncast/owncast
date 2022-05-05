import { Button } from 'antd';
import { ExternalAction } from '../interfaces/external-action.interface';

interface Props {
  action: ExternalAction;
}

export default function ExternalActionButton(props: Props) {
  const { action } = props;
  const { url, title, description, icon, color, openExternally } = action;

  return (
    <Button type="primary" style={{ backgroundColor: color }}>
      <img src={icon} width="30px" alt={description} />
      {title}
    </Button>
  );
}
