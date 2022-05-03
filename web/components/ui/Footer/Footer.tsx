import { Layout } from 'antd';

const { Footer } = Layout;

export default function FooterComponent(props) {
  const { version } = props;

  return <Footer style={{ textAlign: 'center', height: '64px' }}>Footer: Owncast {version}</Footer>;
}
