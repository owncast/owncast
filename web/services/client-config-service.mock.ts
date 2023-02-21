import { ClientConfig } from '../interfaces/client-config.model';
import { ClientConfigStaticService } from './client-config-service';

export const clientConfigServiceMockOf = (config: ClientConfig): ClientConfigStaticService =>
  class ClientConfigServiceMock {
    public static async getConfig(): Promise<ClientConfig> {
      return config;
    }
  };

export default clientConfigServiceMockOf;
