import s from './ChatUserBadge.module.scss';

interface Props {
  badge: string;
  userColor: number;
}

export default function ChatUserBadge(props: Props) {
  const { badge, userColor } = props;
  const color = `var(--theme-user-colors-${userColor})`;
  const style = { color, borderColor: color };

  return (
    <span style={style} className={s.badge}>
      {badge}
    </span>
  );
}
