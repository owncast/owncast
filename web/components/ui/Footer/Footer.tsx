import { Layout } from 'antd';

const { Footer } = Layout;

export default function FooterComponent(props) {
  const { version } = props;

  return <Footer style={{ textAlign: 'center' }}>Footer: Owncast {version}</Footer>;
}
