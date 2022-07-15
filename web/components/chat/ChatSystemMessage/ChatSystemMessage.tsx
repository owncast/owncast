import s from './ChatSystemMessage.module.scss';

interface Props {
  message: string;
}
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function ChatSystemMessage({ message }: Props) {
  return <div className={s.chatSystemMessage}>{message}</div>;
}
