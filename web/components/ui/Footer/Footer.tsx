import { Layout } from 'antd';
import s from './Footer.module.scss';

const { Footer } = Layout;

interface Props {
  version: string;
}

export default function FooterComponent(props: Props) {
  const { version } = props;

  return (
    <Footer className={s.footer}>
      <a href="https://owncast.online">{version}</a>
    </Footer>
  );
}
