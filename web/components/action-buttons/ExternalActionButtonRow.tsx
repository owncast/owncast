import { ExternalAction } from '../interfaces/external-action.interface';
import ExternalActionButton from './ExternalActionButton';
import s from './ExternalActionButtons.module.scss';

interface Props {
  actions: ExternalAction[];
}

export default function ExternalActionButtonRow(props: Props) {
  const { actions } = props;

  return (
    <div className={`${s.row}`}>
      {actions.map(action => (
        <ExternalActionButton key={action.id} action={action} />
      ))}
    </div>
  );
}
