import { Layout } from 'antd';
import { useRecoilValue } from 'recoil';
import {
  ClientConfigStore,
  isChatAvailableSelector,
  clientConfigStateAtom,
  fatalErrorStateAtom,
} from '../stores/ClientConfigStore';
import { Content, Header } from '../ui';
import { ClientConfig } from '../../interfaces/client-config.model';
import { DisplayableError } from '../../types/displayable-error';
import FatalErrorStateModal from '../modals/FatalErrorModal';

function Main() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { name, title } = clientConfig;
  const isChatAvailable = useRecoilValue<boolean>(isChatAvailableSelector);
  const fatalError = useRecoilValue<DisplayableError>(fatalErrorStateAtom);

  return (
    <>
      <ClientConfigStore />
      <Layout>
        <Header name={title || name} chatAvailable={isChatAvailable} />
        <Content />
        {fatalError && (
          <FatalErrorStateModal title={fatalError.title} message={fatalError.message} />
        )}
      </Layout>
    </>
  );
}

export default Main;
