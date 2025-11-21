import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';
import styles from './ChatModeratorNotification.module.scss';
import Icon from '../../../assets/images/moderator.svg';

export const ChatModeratorNotification = () => (
  <div className={styles.chatModerationNotification}>
    <Icon className={styles.icon} />
    <Translation
      translationKey={Localization.Frontend.Chat.moderatorNotification}
      defaultText="You are now a moderator."
    />
  </div>
);
