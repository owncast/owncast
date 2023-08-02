// TODO: add a notication after updating info that changes will take place either on a new stream or server restart. may be different for each field.

import React, { useState, useEffect, FC, ReactElement } from 'react';

import { STATUS, fetchData, FETCH_INTERVAL, SERVER_CONFIG } from './apis';
import { ConfigDetails, UpdateArgs } from '../types/config-section';
import { DEFAULT_VARIANT_STATE } from './config-constants';

const initialServerConfigState: ConfigDetails = {
  streamKeys: [],
  streamKeyOverridden: false,
  adminPassword: '',
  instanceDetails: {
    customStyles: '',
    customJavascript: '',
    extraPageContent: '',
    logo: '',
    name: '',
    nsfw: false,
    socialHandles: [],
    streamTitle: '',
    summary: '',
    tags: [],
    title: '',
    welcomeMessage: '',
    offlineMessage: '',
    appearanceVariables: {},
  },
  ffmpegPath: '',
  rtmpServerPort: '',
  webServerPort: '',
  socketHostOverride: null,
  videoServingEndpoint: '',
  s3: {
    accessKey: '',
    acl: '',
    bucket: '',
    enabled: false,
    endpoint: '',
    region: '',
    secret: '',
    pathPrefix: '',
    forcePathStyle: false,
  },
  yp: {
    enabled: false,
    instanceUrl: '',
  },
  videoSettings: {
    latencyLevel: 4,
    cpuUsageLevel: 3,
    videoQualityVariants: [DEFAULT_VARIANT_STATE],
  },
  federation: {
    enabled: false,
    isPrivate: false,
    username: '',
    goLiveMessage: '',
    showEngagement: true,
    blockedDomains: [],
  },
  notifications: {
    browser: { enabled: false, goLiveMessage: '' },
    discord: { enabled: false, webhook: '', goLiveMessage: '' },
  },
  externalActions: [],
  supportedCodecs: [],
  videoCodec: '',
  forbiddenUsernames: [],
  suggestedUsernames: [],
  chatDisabled: false,
  chatJoinMessagesEnabled: true,
  chatEstablishedUserMode: false,
  hideViewerCount: false,
  disableSearchIndexing: false,
};

const initialServerStatusState = {
  broadcastActive: false,
  broadcaster: null,
  currentBroadcast: null,
  online: false,
  viewerCount: 0,
  sessionMaxViewerCount: 0,
  sessionPeakViewerCount: 0,
  overallPeakViewerCount: 0,
  versionNumber: '0.0.0',
  streamTitle: '',
  chatDisabled: false,
  health: {
    healthy: true,
    healthPercentage: 100,
    message: '',
    representation: 0,
  },
  error: {
    type: null,
    msg: null,
  },
};

export const ServerStatusContext = React.createContext({
  ...initialServerStatusState,
  serverConfig: initialServerConfigState,

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setFieldInConfigState: (_args: UpdateArgs) => null,
});

export type ServerStatusProviderProps = {
  children: ReactElement;
};

const ServerStatusProvider: FC<ServerStatusProviderProps> = ({ children }) => {
  const [status, setStatus] = useState(initialServerStatusState);
  const [config, setConfig] = useState(initialServerConfigState);

  const getStatus = async () => {
    try {
      const result = await fetchData(STATUS);

      if (result instanceof Error) {
        throw result;
      }

      setStatus({ ...result, error: { type: null, msg: null } });
    } catch (error) {
      setStatus(initialStatus => ({
        ...initialStatus,
        error: {
          type: 'OWNCAST_SERVICE_UNREACHABLE',
          msg: 'Cannot connect to the Owncast service. Please check you are connected to the internet and the Owncast server is running.',
        },
      }));
      // todo
    }
  };
  const getConfig = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);

      if (result instanceof Error) {
        throw result;
      }

      setConfig(result);
    } catch (error) {
      console.error(error);
    }
  };

  const setFieldInConfigState = ({ fieldName, value, path }: UpdateArgs) => {
    const updatedConfig = path
      ? {
          ...config,
          [path]: {
            ...config[path],
            [fieldName]: value,
          },
        }
      : {
          ...config,
          [fieldName]: value,
        };
    setConfig(updatedConfig);
  };

  useEffect(() => {
    let getStatusIntervalId = null;

    getStatus();
    getStatusIntervalId = setInterval(getStatus, FETCH_INTERVAL);

    getConfig();

    // returned function will be called on component unmount
    return () => {
      clearInterval(getStatusIntervalId);
    };
  }, []);

  // eslint-disable-next-line react/jsx-no-constructed-context-values
  const providerValue = {
    ...status,
    serverConfig: config,

    setFieldInConfigState,
  };
  return (
    <ServerStatusContext.Provider value={providerValue}>{children}</ServerStatusContext.Provider>
  );
};

export default ServerStatusProvider;
