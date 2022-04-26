import { useEffect } from 'react';
import { ReactElement } from 'react-markdown/lib/react-markdown';
import { atom, useRecoilState } from 'recoil';
import { makeEmptyClientConfig, ClientConfig } from '../../models/ClientConfig';
import ClientConfigService from '../../services/ClientConfigService';

export const clientConfigState = atom({
  key: 'clientConfigState',
  default: makeEmptyClientConfig(),
});

export function ClientConfigStore(): ReactElement {
  const [, setClientConfig] = useRecoilState<ClientConfig>(clientConfigState);

  const updateClientConfig = async () => {
    try {
      const config = await ClientConfigService.getConfig();
      console.log(`ClientConfig: ${JSON.stringify(config)}`);
      setClientConfig(config);
    } catch (error) {
      console.error(`ClientConfigService -> getConfig() ERROR: \n${error}`);
    }
  };

  useEffect(() => {
    updateClientConfig();
  }, []);

  return null;
}
