import s from './ChatJoinMessage.module.scss';
import ChatUserBadge from '../ChatUserBadge/ChatUserBadge';

interface Props {
  isAuthorModerator: boolean;
  userColor: number;
  displayName: string;
}

export default function ChatJoinMessage(props: Props) {
  const { isAuthorModerator, userColor, displayName } = props;
  const color = `var(--theme-user-colors-${userColor})`;

  return (
    <div className={s.join}>
      <span style={{ color }}>
        {displayName}
        {isAuthorModerator && (
          <span>
            <ChatUserBadge badge="mod" userColor={userColor} />
          </span>
        )}
      </span>{' '}
      joined the chat.
    </div>
  );
}
