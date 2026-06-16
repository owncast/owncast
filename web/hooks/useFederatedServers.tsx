import { useState, useEffect } from 'react';
import { message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../types/localization';

// FederatedServerResponse matches the OpenAPI FederatedServer model the
// backend serializes verbatim. The web consumes these names directly so
// there is a single, shared vocabulary for federated-server data rather
// than a backend set and a separate frontend set.
export interface FederatedServerResponse {
  id: number;
  iri: string;
  name?: string;
  displayName?: string;
  logoUrl?: string;
  isOnline: boolean;
  streamTitle?: string;
  streamDescription?: string;
  tags?: string[];
  thumbnailUrl?: string;
  lastStatusUpdate?: string;
  addedAt: string;
}

export interface UseFederatedServersResult {
  servers: FederatedServerResponse[];
  loading: boolean;
  error: string | null;
  refetch: () => void;
  addServer: (url: string) => Promise<void>;
  removeServer: (id: number) => Promise<void>;
}

interface APIErrorResponse {
  message?: string;
  errorCode?: string;
}

// API endpoints - these will need to be implemented on the backend
const API_FEDERATED_SERVERS = '/api/federation/servers';
const API_ADD_FEDERATED_SERVER = '/api/admin/federation/servers';
const API_REMOVE_FEDERATED_SERVER = '/api/admin/federation/servers';
const UNSUPPORTED_FEATURED_STREAMS_ERROR_CODE = 'UNSUPPORTED_FEATURED_STREAMS';

function getFederatedServerErrorMessage(
  error: APIErrorResponse,
  t: (key: string, query?: object) => string,
) {
  if (error.errorCode === UNSUPPORTED_FEATURED_STREAMS_ERROR_CODE) {
    return t(Localization.Admin.FeaturedStreams.unsupportedFeaturedStreams);
  }

  return error.message || t(Localization.Admin.FeaturedStreams.failedToFeature);
}

export function useFederatedServers(isAdmin: boolean = false): UseFederatedServersResult {
  const { t } = useTranslation();
  const [servers, setServers] = useState<FederatedServerResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchServers = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(API_FEDERATED_SERVERS);

      if (!response.ok) {
        throw new Error(`Failed to fetch servers: ${response.statusText}`);
      }

      const data = await response.json();
      setServers(data.servers || []);
    } catch (err: any) {
      const errorMessage = err.message || t(Localization.Admin.FeaturedStreams.failedToFeature);
      setError(errorMessage);

      // Only show error message in admin interface
      if (isAdmin) {
        message.error(errorMessage);
      }
    } finally {
      setLoading(false);
    }
  };

  const addServer = async (url: string): Promise<void> => {
    const response = await fetch(API_ADD_FEDERATED_SERVER, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ url }),
    });

    if (!response.ok) {
      const apiError: APIErrorResponse = await response.json();
      throw new Error(getFederatedServerErrorMessage(apiError, t));
    }

    // Refetch the server list
    await fetchServers();
    message.success(t(Localization.Admin.FeaturedStreams.streamFeaturedSuccess));
  };

  const removeServer = async (id: number): Promise<void> => {
    const response = await fetch(`${API_REMOVE_FEDERATED_SERVER}/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    if (!response.ok) {
      const apiError: APIErrorResponse = await response.json();
      throw new Error(apiError.message || t(Localization.Admin.FeaturedStreams.failedToUnfeature));
    }

    // Refetch the server list
    await fetchServers();
  };

  useEffect(() => {
    fetchServers();
  }, []);

  return {
    servers,
    loading,
    error,
    refetch: fetchServers,
    addServer: isAdmin ? addServer : async () => {},
    removeServer: isAdmin ? removeServer : async () => {},
  };
}
