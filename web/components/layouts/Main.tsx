import { Layout } from 'antd';
import { useRecoilValue } from 'recoil';
import {
  ClientConfigStore,
  isChatAvailableSelector,
  clientConfigStateAtom,
} from '../stores/ClientConfigStore';
import { Content, Header } from '../ui';
import { ClientConfig } from '../../interfaces/client-config.model';

function Main() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name, title } = clientConfig;
  const isChatAvailable = useRecoilValue<boolean>(isChatAvailableSelector);
  return (
    <>
      <ClientConfigStore />
      <Layout>
        <Header name={title || name} chatAvailable={isChatAvailable} />
        <Content />
      </Layout>
    </>
  );
}

export default Main;
