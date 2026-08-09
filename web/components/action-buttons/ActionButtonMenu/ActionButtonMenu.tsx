import { FC } from 'react';
import { useTranslation } from 'next-export-i18n';
import { Button, Dropdown } from 'antd';
import type { MenuProps } from 'antd';
import classNames from 'classnames';
import dynamic from 'next/dynamic';
import styles from './ActionButtonMenu.module.scss';
import { ExternalAction, ExternalActionUtils } from '../../../interfaces/external-action';
import { Localization } from '../../../types/localization';

// Lazy loaded components

const EllipsisOutlined = dynamic(() => import('@ant-design/icons/EllipsisOutlined'), {
  ssr: false,
});

const HeartOutlined = dynamic(() => import('@ant-design/icons/HeartOutlined'), {
  ssr: false,
});

const BellOutlined = dynamic(() => import('@ant-design/icons/BellOutlined'), {
  ssr: false,
});

const ScissorOutlined = dynamic(() => import('@ant-design/icons/ScissorOutlined'), {
  ssr: false,
});

const NOTIFY_KEY = 'notify';
const FOLLOW_KEY = 'follow';
const CLIP_KEY = 'clip';

export type ActionButtonMenuProps = {
  actions: ExternalAction[];
  showFollowItem?: boolean;
  showNotifyItem?: boolean;
  // showClipItem offers a clip action to a viewer allowed to clip.
  showClipItem?: boolean;
  // clipItemEnabled is false when there is nothing to clip yet.
  clipItemEnabled?: boolean;
  clipActive?: boolean;
  clipRemainingSeconds?: number;
  externalActionSelected: (action: ExternalAction) => void;
  notifyItemSelected: () => void;
  followItemSelected: () => void;
  clipItemSelected?: () => void;
  className?: string;
};

export const ActionButtonMenu: FC<ActionButtonMenuProps> = ({
  actions,
  externalActionSelected,
  notifyItemSelected,
  followItemSelected,
  clipItemSelected,
  showFollowItem,
  showNotifyItem,
  showClipItem,
  clipItemEnabled = true,
  clipActive = false,
  clipRemainingSeconds = 0,
  className,
}) => {
  const { t } = useTranslation();
  const onClick = a => {
    if (a.key === NOTIFY_KEY) {
      notifyItemSelected();
      return;
    }
    if (a.key === FOLLOW_KEY) {
      followItemSelected();
      return;
    }
    if (a.key === CLIP_KEY) {
      clipItemSelected?.();
      return;
    }
    // Find the action using the utility function
    const action = ExternalActionUtils.findByKey(actions, a.key);
    if (action) {
      externalActionSelected(action);
    }
  };

  // Typed explicitly so entries can carry menu properties beyond key/label,
  // such as the disabled clip item.
  const items: NonNullable<MenuProps['items']> = actions.map(action => ({
    key: ExternalActionUtils.generateKey(action),
    label: (
      <span className={styles.item}>
        {action.icon && <img className={styles.icon} src={action.icon} alt={action.title} />}{' '}
        {action.title}
      </span>
    ),
  }));

  if (showFollowItem) {
    items.unshift({
      key: FOLLOW_KEY,
      label: (
        <span className={styles.item}>
          <HeartOutlined className={styles.icon} /> Follow this stream
        </span>
      ),
    });
  }

  if (showNotifyItem) {
    items.unshift({
      key: NOTIFY_KEY,
      label: (
        <span className={styles.item}>
          <BellOutlined className={styles.icon} />
          Notify when live
        </span>
      ),
    });
  }

  // Clipping is the most time-sensitive action available while a stream is
  // live, so it leads the menu.
  if (showClipItem) {
    const remaining = `${Math.floor(clipRemainingSeconds / 60)}:${String(
      clipRemainingSeconds % 60,
    ).padStart(2, '0')}`;
    items.unshift({
      key: CLIP_KEY,
      disabled: !clipItemEnabled && !clipActive,
      label: (
        <span className={styles.item}>
          <ScissorOutlined className={styles.icon} />
          {clipActive
            ? `${t(Localization.Frontend.Clips.end)} (${remaining})`
            : t(Localization.Frontend.Clips.start)}
        </span>
      ),
    });
  }

  const dropdownClasses = classNames([styles.menu, className]);

  return (
    <Dropdown
      menu={{ items, onClick }}
      placement="bottomRight"
      trigger={['click']}
      className={dropdownClasses}
    >
      <div className={styles.buttonWrap}>
        <Button
          type="default"
          onClick={e => e.preventDefault()}
          size="large"
          icon={<EllipsisOutlined size={6} style={{ rotate: '90deg' }} />}
          className={styles.menuButton}
        />
      </div>
    </Dropdown>
  );
};
