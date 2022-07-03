import s from './Logo.module.scss';

interface Props {
  src: string;
}

export default function Logo({ src }: Props) {
  return (
    <div className={s.root}>
      <div className={s.container}>
        <div className={s.image} style={{ backgroundImage: `url(${src})` }} />
      </div>
    </div>
  );
}
