import { ServerStatus } from '../interfaces/server-status.model';

const ENDPOINT = `/api/status`;

class ServerStatusService {
  public static async getStatus(): Promise<ServerStatus> {
    const response = await fetch(ENDPOINT);
    const status = await response.json();
    return status;
  }
}

export default ServerStatusService;
