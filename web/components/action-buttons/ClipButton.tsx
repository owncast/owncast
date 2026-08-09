import { Button } from 'antd';
import { FC } from 'react';
import dynamic from 'next/dynamic';
import { useTranslation } from 'next-export-i18n';
import styles from './ActionButton/ActionButton.module.scss';
import { Localization } from '../../types/localization';

// Lazy loaded components

const ScissorOutlined = dynamic(() => import('@ant-design/icons/ScissorOutlined'), {
  ssr: false,
});

export type ClipButtonProps = {
  onClick?: () => void;
  size?: 'small' | 'middle' | 'large';
  disabled?: boolean;
  active?: boolean;
  remainingSeconds?: number;
};

// ClipButton starts a clip and then exposes the end action with its remaining
// maximum time.
export const ClipButton: FC<ClipButtonProps> = ({
  onClick,
  size,
  disabled,
  active = false,
  remainingSeconds = 0,
}) => {
  const { t } = useTranslation();
  const remaining = `${Math.floor(remainingSeconds / 60)}:${String(remainingSeconds % 60).padStart(
    2,
    '0',
  )}`;

  return (
    <Button
      type="primary"
      size={size}
      disabled={disabled}
      className={styles.button}
      icon={<ScissorOutlined />}
      onClick={onClick}
      id="clip-button"
    >
      {active
        ? `${t(Localization.Frontend.Clips.end)} (${remaining})`
        : t(Localization.Frontend.Clips.start)}
    </Button>
  );
};
