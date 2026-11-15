import { Button } from 'antd';
import { FC } from 'react';
import dynamic from 'next/dynamic';
import { Translation } from '../ui/Translation/Translation';
import { Localization } from '../../types/localization';
import styles from './ActionButton/ActionButton.module.scss';

// Lazy loaded components

const BellFilled = dynamic(() => import('@ant-design/icons/BellFilled'), {
  ssr: false,
});

export type NotifyButtonProps = {
  text?: string;
  onClick?: () => void;
};

export const NotifyButton: FC<NotifyButtonProps> = ({ onClick, text }) => {
  return (
    <Button
      type="primary"
      className={styles.button}
      icon={<BellFilled />}
      onClick={onClick}
      id="notify-button"
    >
      {text || (
        <Translation
          translationKey={Localization.Frontend.NotifyButton.label}
          defaultText="Notify"
        />
      )}
    </Button>
  );
};
