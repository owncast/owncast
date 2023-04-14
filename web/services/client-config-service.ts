import { createContext } from 'react';
import { ClientConfig } from '../interfaces/client-config.model';

const ENDPOINT = `/api/config`;

export interface ClientConfigStaticService {
  getConfig(): Promise<ClientConfig>;
}

class ClientConfigService {
  public static async getConfig(): Promise<ClientConfig> {
    const response = await fetch(ENDPOINT);
    const status = await response.json();
    return status;
  }
}

export const ClientConfigServiceContext =
  createContext<ClientConfigStaticService>(ClientConfigService);
