import ServerStatus from '../models/ServerStatus';

const ENDPOINT = `http://localhost:8080/api/status`;

class ServerStatusService {
  public static async getStatus(): Promise<ServerStatus> {
    const response = await fetch(ENDPOINT);
    const status = await response.json();
    return status;
  }
}

export default ServerStatusService;
