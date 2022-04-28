import { useRecoilValue } from 'recoil';
import { clientConfigState } from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { Layout } from 'antd';

const { Footer } = Layout;

export default function FooterComponent() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigState);
  const { version } = clientConfig;

  return <Footer style={{ textAlign: 'center' }}>Footer: Owncast {version}</Footer>;
}
