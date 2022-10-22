import { Button } from 'antd';
import { FC } from 'react';
import { ExternalAction } from '../../../interfaces/external-action';
import styles from './ActionButton.module.scss';

export type ActionButtonProps = {
  action: ExternalAction;
  primary?: boolean;
  externalActionSelected: (action: ExternalAction) => void;
};

export const ActionButton: FC<ActionButtonProps> = ({
  action,
  primary = true,
  externalActionSelected,
}) => {
  const { title, description, icon, color } = action;

  return (
    <Button
      type={primary ? 'primary' : 'default'}
      className={`${styles.button}`}
      onClick={() => externalActionSelected(action)}
      style={{ backgroundColor: color }}
    >
      {icon && <img src={icon} className={`${styles.icon}`} alt={description} />}
      {title}
    </Button>
  );
};
