import styles from './ChatModeratorNotification.module.scss';
import Icon from '../../../assets/images/moderator.svg';

export const ChatModeratorNotification = () => (
  <div className={styles.chatModerationNotification}>
    <Icon className={styles.icon} />
    You are now a moderator.
  </div>
);
