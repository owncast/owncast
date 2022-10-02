import { ClientConfig } from '../interfaces/client-config.model';

const ENDPOINT = `/api/config`;

class ClientConfigService {
  public static async getConfig(): Promise<ClientConfig> {
    const response = await fetch(ENDPOINT);
    const status = await response.json();
    return status;
  }
}

export default ClientConfigService;
