import { Button } from 'antd';
import { FC } from 'react';
import cn from 'classnames';
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
      className={cn([`${styles.button}`, 'action-button'])}
      onClick={() => externalActionSelected(action)}
      style={{ backgroundColor: color }}
      title={description || title}
    >
      {icon && <img src={icon} className={styles.icon} alt={description} />}
      {title}
    </Button>
  );
};
