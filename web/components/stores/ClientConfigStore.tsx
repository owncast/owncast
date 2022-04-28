import { useEffect } from 'react';
import { ReactElement } from 'react-markdown/lib/react-markdown';
import { atom, useRecoilState } from 'recoil';
import { makeEmptyClientConfig, ClientConfig } from '../../interfaces/client-config.model';
import ClientConfigService from '../../services/client-config-service';

// The config that comes from the API.
export const clientConfigState = atom({
  key: 'clientConfigState',
  default: makeEmptyClientConfig(),
});

export const chatCurrentlyVisible = atom({
  key: 'chatvisible',
  default: false,
});

export const chatDislayName = atom({
  key: 'chatDisplayName',
  default: '',
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
