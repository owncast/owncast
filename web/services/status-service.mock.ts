import { ServerStatus } from '../interfaces/server-status.model';
import { ServerStatusStaticService } from './status-service';

export const serverStatusServiceMockOf = (serverStatus: ServerStatus): ServerStatusStaticService =>
  class ServerStatusServiceMock {
    public static async getStatus(): Promise<ServerStatus> {
      return serverStatus;
    }
  };

export default serverStatusServiceMockOf;
