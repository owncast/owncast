import { createContext } from 'react';
import { ServerStatus } from '../interfaces/server-status.model';

const ENDPOINT = `/api/status`;

export interface ServerStatusStaticService {
  getStatus(): Promise<ServerStatus>;
}

class ServerStatusService {
  public static async getStatus(): Promise<ServerStatus> {
    const response = await fetch(ENDPOINT);
    const status = await response.json();
    return status;
  }
}

export const ServerStatusServiceContext =
  createContext<ServerStatusStaticService>(ServerStatusService);
