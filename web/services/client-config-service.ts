import { ClientConfig } from '../interfaces/client-config.model';

const originalFetch = require('isomorphic-fetch');
const fetch = require('fetch-retry')(originalFetch);
const ENDPOINT = `http://localhost:8080/api/config`;

class ClientConfigService {
  public static async getConfig(): Promise<ClientConfig> {
    const response = await fetch(ENDPOINT, {
      retries: 20,
      retryDelay: (attempt, error, response) => 3 ** attempt * 1000,
    });
    const status = await response.json();
    return status;
  }
}

export default ClientConfigService;
