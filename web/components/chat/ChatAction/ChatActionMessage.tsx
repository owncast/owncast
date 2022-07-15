import s from './ChatActionMessage.module.scss';

/* eslint-disable react/no-danger */
interface Props {
  body: string;
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function ChatActionMessage(props: Props) {
  const { body } = props;

  return <div dangerouslySetInnerHTML={{ __html: body }} className={s.chatAction} />;
}
