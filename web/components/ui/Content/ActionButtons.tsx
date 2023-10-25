import { Dispatch, FC, SetStateAction } from 'react';
import dynamic from 'next/dynamic';
import { Skeleton } from 'antd';
import { ExternalAction } from '../../../interfaces/external-action';
import { ActionButtonMenu } from '../../action-buttons/ActionButtonMenu/ActionButtonMenu';
import { ActionButtonRow } from '../../action-buttons/ActionButtonRow/ActionButtonRow';
import { FollowButton } from '../../action-buttons/FollowButton';
import { NotifyButton } from '../../action-buttons/NotifyButton';
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
  setExternalActionToDisplay: any;
  externalActionSelected: (action: ExternalAction) => void;
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
  setExternalActionToDisplay,
  externalActionSelected,
}) => {
  const externalActionButtons = externalActions.map(action => (
    <ActionButton
      key={action.url || action.html}
      action={action}
      externalActionSelected={externalActionSelected}
    />
  ));

  return (
    <div className={styles.actionButtonsContainer}>
      <div className={styles.desktopActionButtons}>
        <ActionButtonRow>
          {externalActionButtons}
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
        {(supportsBrowserNotifications ||
          supportsBrowserNotifications ||
          externalActionButtons.length > 0) && (
          <ActionButtonMenu
            className={styles.actionButtonMenu}
            showFollowItem={supportFediverseFeatures}
            showNotifyItem={supportsBrowserNotifications}
            actions={externalActions}
            externalActionSelected={setExternalActionToDisplay}
            notifyItemSelected={() => setShowNotifyModal(true)}
            followItemSelected={() => setShowFollowModal(true)}
          />
        )}
      </div>
    </div>
  );
};

export default ActionButtons;
