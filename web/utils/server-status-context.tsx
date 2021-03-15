// TODO: add a notication after updating info that changes will take place either on a new stream or server restart. may be different for each field.

import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { STATUS, fetchData, FETCH_INTERVAL, SERVER_CONFIG } from './apis';
import { ConfigDetails, UpdateArgs, ExternalAction } from '../types/config-section';
import { DEFAULT_VARIANT_STATE } from './config-constants';

export const initialServerConfigState: ConfigDetails = {
  streamKey: '',
  instanceDetails: {
    extraPageContent: '',
    logo: '',
    name: '',
    nsfw: false,
    socialHandles: [],
    streamTitle: '',
    summary: '',
    welcomeMessage: '',
    tags: [],
    title: '',
  },
  ffmpegPath: '',
  rtmpServerPort: '',
  webServerPort: '',
  chatDisabled: false,
  s3: {
    accessKey: '',
    acl: '',
    bucket: '',
    enabled: false,
    endpoint: '',
    region: '',
    secret: '',
    servingEndpoint: '',
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
  externalActions: [],
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
};

export const ServerStatusContext = React.createContext({
  ...initialServerStatusState,
  serverConfig: initialServerConfigState,

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setFieldInConfigState: (args: UpdateArgs) => null,
});

const ServerStatusProvider = ({ children }) => {
  const [status, setStatus] = useState(initialServerStatusState);
  const [config, setConfig] = useState(initialServerConfigState);

  const getStatus = async () => {
    try {
      const result = await fetchData(STATUS);
      setStatus({ ...result });
    } catch (error) {
      // todo
    }
  };
  const getConfig = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      setConfig(result);
    } catch (error) {
      // todo
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

  const providerValue = {
    ...status,
    serverConfig: config,

    setFieldInConfigState,
  };
  return (
    <ServerStatusContext.Provider value={providerValue}>{children}</ServerStatusContext.Provider>
  );
};

ServerStatusProvider.propTypes = {
  children: PropTypes.element.isRequired,
};

export default ServerStatusProvider;
