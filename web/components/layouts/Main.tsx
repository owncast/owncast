import { Layout } from 'antd';
import { useRecoilValue } from 'recoil';
import { ClientConfigStore, clientConfigStateAtom } from '../stores/ClientConfigStore';
import { Content, Header } from '../ui';
import { ClientConfig } from '../../interfaces/client-config.model';

function Main() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name, title } = clientConfig;

  return (
    <>
      <ClientConfigStore />
      <Layout>
        <Header name={title || name} />
        <Content />
      </Layout>
    </>
  );
}

export default Main;
