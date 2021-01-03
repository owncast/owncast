// DEFAULT VALUES
import React from 'react';
import { CheckCircleFilled, ExclamationCircleFilled } from '@ant-design/icons';
import { fetchData, SERVER_CONFIG_UPDATE_URL } from '../../../utils/apis';
import { ApiPostArgs } from '../../../types/config-section';

export const DEFAULT_NAME = 'Owncast User';
export const DEFAULT_TITLE = 'Owncast Server';
export const DEFAULT_SUMMARY = '';

export const TEXT_MAXLENGTH = 255;

export const RESET_TIMEOUT = 3000;

export const SUCCESS_STATES = {
  success: {
    icon: <CheckCircleFilled style={{ color: 'green' }} />,
    message: 'Success!',
  },
  error: {
    icon: <ExclamationCircleFilled style={{ color: 'red' }} />,
    message: 'An error occurred.',
  },
};


export async function postConfigUpdateToAPI(args: ApiPostArgs) {
  const {
    apiPath,
    data,
    onSuccess,
    onError,
  } = args;
  const result = await fetchData(`${SERVER_CONFIG_UPDATE_URL}${apiPath}`, {
    data,
    method: 'POST',
    auth: true,
  });
  if (result.success && onSuccess) {
    onSuccess(result.message);
  } else if (onError) {
    onError(result.message);
  }
}

// Creating this so that it'll be easier to change values in one place, rather than looking for places to change it in a sea of JSX.

// key is the input's `fieldName`
// the structure of this mirrors config data
export const TEXTFIELD_DEFAULTS = {
  instanceDetails: {

    // user name
    name: {
      apiPath: '/name',
      defaultValue: '',
      maxLength: TEXT_MAXLENGTH,
      placeholder: 'username',
      label: 'User name',
      tip: 'Who are you? What name do you want viewers to know you?',
    },

    // like "goth land"
    title: {
      apiPath: '/servertitle',
      defaultValue: '',
      maxLength: TEXT_MAXLENGTH,
      placeholder: 'Owncast site name',
      label: 'Server Name',
      tip: 'The name of your Owncast server',
    },


    streamTitle: {
      apiPath: '/streamtitle',
      defaultValue: '',
      maxLength: TEXT_MAXLENGTH,
      placeholder: 'Doing cool things...',
      label: 'Stream Title',
      tip: 'What is your stream about today?',
    },

    summary: {
      apiPath: '/serversummary',
      defaultValue: '',
      maxLength: 500,
      placeholder: 'Summary',
      label: 'Summary',
      tip: 'A brief blurb about what your stream is about.',
    },
  
    logo: {
      apiPath: '/logo',
      defaultValue: '',
      maxLength: 255,
      placeholder: '/img/mylogo.png',
      label: 'Logo',
      tip: 'Path to your logo from website root. We recommend that you use a square image that is at least 256x256. (upload functionality coming soon)',
    },
  
    extraPageContent: {
      apiPath: '/pagecontent',
      placeholder: '',
      label: 'Extra page content',
      tip: 'Custom markup about yourself',
    },  

    nsfw: {
      apiPath: '/nsfw',
      placeholder: '',
      label: 'NSFW?',
      tip: "Turn this ON if you plan to steam explicit or adult content. You may want to respectfully set this flag so that unexpecting eyes won't accidentally see it from the Directory.",
    },  

    tags: {
      apiPath: '/tags',
      defaultValue: '',
      maxLength: 24,
      placeholder: 'Add a new tag',
      label: '',
      tip: '',
    },
  },
  

  streamKey: {
    apiPath: '/key',
    defaultValue: 'abc123',
    maxLength: TEXT_MAXLENGTH,
    placeholder: 'abc123',
    label: 'Stream Key',
    tip: 'Secret stream key',
    required: true,
  },

  ffmpegPath: {
    apiPath: '/ffmpegpath',
    defaultValue: '',
    maxLength: TEXT_MAXLENGTH,
    placeholder: '/usr/local/bin/ffmpeg',
    label: 'FFmpeg Path',
    tip: 'Absolute file path of the FFMPEG application on your server',
    required: true,
  },

  webServerPort: {
    apiPath: '/webserverport',
    defaultValue: '8080',
    maxLength: 6,
    placeholder: '8080',
    label: 'Owncast Server port',
    tip: 'What port are you serving Owncast from? Default is :8080',
    required: true,
  },
  rtmpServerPort: {
    apiPath: '/rtmpserverport',
    defaultValue: '1935',
    maxLength: 6,
    placeholder: '1935',
    label: 'RTMP port',
    tip: 'What port are you receiving RTMP?',
    required: true,
  },

  s3: {
    // tbd
  },

  // YP options
  yp: {
    instanceUrl: {
      apiPath: '/serverurl',
      defaultValue: 'https://owncast.mysite.com',
      maxLength: 255,
      placeholder: 'url',
      label: 'Instance URL',
      tip: 'Please provide the url to your Owncast site if you enable this Directory setting.',
    },
    enabled: {
      apiPath: '/directoryenabled',
      defaultValue: false,
      maxLength: 0,
      placeholder: '',
      label: 'Display in the Owncast Directory?',
      tip: 'Turn this ON if you want to show up in the Owncast directory at https://directory.owncast.online.',
    }
  }
}

