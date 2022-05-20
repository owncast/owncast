// import s from './OfflineBanner.module.scss';

interface Props {
  text: string;
}

export default function OfflineBanner({ text }: Props) {
  return <div>{text}</div>;
}
