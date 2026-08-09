import { Dispatch, FC, SetStateAction } from 'react';
import dynamic from 'next/dynamic';
import { Skeleton } from 'antd';
import { ExternalAction, ExternalActionUtils } from '../../../interfaces/external-action';
import { ActionButtonMenu } from '../../action-buttons/ActionButtonMenu/ActionButtonMenu';
import { ActionButtonRow } from '../../action-buttons/ActionButtonRow/ActionButtonRow';
import { FollowButton } from '../../action-buttons/FollowButton';
import { NotifyButton } from '../../action-buttons/NotifyButton';
import { ClipButton } from '../../action-buttons/ClipButton';
import styles from './Content.module.scss';
import { ActionButton } from '../../action-buttons/ActionButton/ActionButton';

interface ActionButtonProps {
  supportFediverseFeatures: boolean;
  externalActions: ExternalAction[];
  supportsBrowserNotifications: boolean;
  showNotifyReminder: any;
  setShowFollowModal: Dispatch<SetStateAction<boolean>>;
  setShowNotifyModal: Dispatch<SetStateAction<boolean>>;
  disableNotifyReminderPopup: () => void;
  externalActionSelected: (action: ExternalAction) => void;
  // showClipButton renders the clip control for a registered, unbanned viewer.
  showClipButton: boolean;
  // canClip is true when pressing it would work: a live stream is identified
  // and video is actively playing.
  canClip: boolean;
  clipActive: boolean;
  clipRemainingSeconds: number;
  onClipAction: () => void;
}

const NotifyReminderPopup = dynamic(
  () => import('../NotifyReminderPopup/NotifyReminderPopup').then(mod => mod.NotifyReminderPopup),
  {
    ssr: false,
    loading: () => <Skeleton loading active paragraph={{ rows: 8 }} />,
  },
);

const ActionButtons: FC<ActionButtonProps> = ({
  supportFediverseFeatures,
  supportsBrowserNotifications,
  showNotifyReminder,
  setShowFollowModal,
  setShowNotifyModal,
  disableNotifyReminderPopup,
  externalActions,
  externalActionSelected,
  showClipButton,
  canClip,
  clipActive,
  clipRemainingSeconds,
  onClipAction,
}) => {
  const externalActionButtons = externalActions.map(action => (
    <ActionButton
      key={ExternalActionUtils.generateKey(action)}
      action={action}
      externalActionSelected={externalActionSelected}
    />
  ));

  return (
    <div className={styles.actionButtonsContainer}>
      <div className={styles.desktopActionButtons}>
        <ActionButtonRow>
          {externalActionButtons}
          {showClipButton && (
            <ClipButton
              size="small"
              disabled={!canClip && !clipActive}
              active={clipActive}
              remainingSeconds={clipRemainingSeconds}
              onClick={onClipAction}
            />
          )}
          {supportFediverseFeatures && (
            <FollowButton size="small" onClick={() => setShowFollowModal(true)} />
          )}
          {supportsBrowserNotifications && (
            <NotifyReminderPopup
              open={showNotifyReminder}
              notificationClicked={() => setShowNotifyModal(true)}
              notificationClosed={() => disableNotifyReminderPopup()}
            >
              <NotifyButton onClick={() => setShowNotifyModal(true)} />
            </NotifyReminderPopup>
          )}
        </ActionButtonRow>
      </div>
      <div className={styles.mobileActionButtons}>
        {(showClipButton || supportsBrowserNotifications || externalActions.length > 0) && (
          <ActionButtonMenu
            actions={externalActions}
            showFollowItem={supportFediverseFeatures}
            showNotifyItem={supportsBrowserNotifications}
            showClipItem={showClipButton}
            clipItemEnabled={canClip}
            clipActive={clipActive}
            clipRemainingSeconds={clipRemainingSeconds}
            externalActionSelected={externalActionSelected}
            notifyItemSelected={() => setShowNotifyModal(true)}
            followItemSelected={() => setShowFollowModal(true)}
            clipItemSelected={onClipAction}
          />
        )}
      </div>
    </div>
  );
};

export default ActionButtons;
