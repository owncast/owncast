import { ExternalAction } from '../interfaces/external-action.interface';
import ExternalActionButton from './ExternalActionButton';

interface Props {
  actions: ExternalAction[];
}

export default function ExternalActionButtonRow(props: Props) {
  const { actions } = props;

  return (
    <div>
      {actions.map(action => (
        <ExternalActionButton key={action.id} action={action} />
      ))}
    </div>
  );
}
