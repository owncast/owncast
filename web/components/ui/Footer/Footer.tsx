import { Layout } from 'antd';

const { Footer } = Layout;

interface Props {
  version: string;
}

export default function FooterComponent(props: Props) {
  const { version } = props;

  return (
    <Footer style={{ textAlign: 'center', height: '64px' }}>
      <a href="https://owncast.online">{version}</a>
    </Footer>
  );
}
