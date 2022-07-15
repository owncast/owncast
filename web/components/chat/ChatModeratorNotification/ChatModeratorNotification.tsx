import s from './ChatModeratorNotification.module.scss';
import Icon from '../../../assets/images/moderator.svg';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function ModeratorNotification() {
  return (
    <div className={s.chatModerationNotification}>
      <Icon className={s.icon} />
      You are now a moderator.
    </div>
  );
}
