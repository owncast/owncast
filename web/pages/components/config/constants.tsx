// DEFAULT VALUES
import React from 'react';
import { CheckCircleFilled, ExclamationCircleFilled } from '@ant-design/icons';
import { fetchData, SERVER_CONFIG_UPDATE_URL } from '../../../utils/apis';
import { ApiPostArgs, VideoVariant, SocialHandle } from '../../../types/config-section';

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

// CONFIG API ENDPOINTS
export const API_CUSTOM_CONTENT = '/pagecontent';
export const API_FFMPEG = '/ffmpegpath';
export const API_INSTANCE_URL = '/serverurl';
export const API_LOGO = '/logo';
export const API_NSFW_SWITCH = '/nsfw';
export const API_RTMP_PORT = '/rtmpserverport';
export const API_S3_INFO = '/s3';
export const API_SERVER_SUMMARY = '/serversummary';
export const API_SERVER_TITLE = '/servertitle';
export const API_SOCIAL_HANDLES = '/socialhandles';
export const API_STREAM_KEY = '/key';
export const API_STREAM_TITLE = '/streamtitle';
export const API_TAGS = '/tags';
export const API_USERNAME = '/name';
export const API_VIDEO_SEGMENTS = '/video/streamlatencylevel';
export const API_VIDEO_VARIANTS = '/video/streamoutputvariants';
export const API_WEB_PORT = '/webserverport';
export const API_YP_SWITCH = '/directoryenabled';


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


// Some default props to help build out a TextField
export const TEXTFIELD_PROPS_USERNAME = {
  apiPath: API_USERNAME,
  configPath: 'instanceDetails',
  maxLength: TEXT_MAXLENGTH,
  placeholder: 'username',
  label: 'User name',
  tip: 'Who are you? What name do you want viewers to know you?',
};
export const TEXTFIELD_PROPS_SERVER_TITLE = {
  apiPath: API_SERVER_TITLE,
  maxLength: TEXT_MAXLENGTH,
  placeholder: 'Owncast site name', // like "gothland"
  label: 'Server Name',
  tip: 'The name of your Owncast server',
};
export const TEXTFIELD_PROPS_STREAM_TITLE = {
  apiPath: API_STREAM_TITLE,
  maxLength: TEXT_MAXLENGTH,
  placeholder: 'Doing cool things...',
  label: 'Stream Title',
  tip: 'What is your stream about today?',
};
export const TEXTFIELD_PROPS_SERVER_SUMMARY = {
  apiPath: API_SERVER_SUMMARY,
  maxLength: 500,
  placeholder: 'Summary',
  label: 'Summary',
  tip: 'A brief blurb about what your stream is about.',
};
export const TEXTFIELD_PROPS_LOGO = {
  apiPath: API_LOGO,
  maxLength: 255,
  placeholder: '/img/mylogo.png',
  label: 'Logo',
  tip: 'Path to your logo from website root. We recommend that you use a square image that is at least 256x256. (upload functionality coming soon)',
};
export const TEXTFIELD_PROPS_STREAM_KEY = {
  apiPath: API_STREAM_KEY,
  configPath: '',
  maxLength: TEXT_MAXLENGTH,
  placeholder: 'abc123',
  label: 'Stream Key',
  tip: 'Secret stream key',
  required: true,
};
export const TEXTFIELD_PROPS_FFMPEG = {
  apiPath: API_FFMPEG,
  configPath: '',
  maxLength: TEXT_MAXLENGTH,
  placeholder: '/usr/local/bin/ffmpeg',
  label: 'FFmpeg Path',
  tip: 'Absolute file path of the FFMPEG application on your server',
  required: true,
};
export const TEXTFIELD_PROPS_WEB_PORT = {
  apiPath: API_WEB_PORT,
  configPath: '',
  maxLength: 6,
  placeholder: '8080',
  label: 'Owncast Server port',
  tip: 'What port are you serving Owncast from? Default is :8080',
  required: true,
};
export const TEXTFIELD_PROPS_RTMP_PORT = {
  apiPath: API_RTMP_PORT,
  configPath: '',
  maxLength: 6,
  placeholder: '1935',
  label: 'RTMP port',
  tip: 'What port are you receiving RTMP?',
  required: true,
};
export const TEXTFIELD_PROPS_INSTANCE_URL = {
  apiPath: API_INSTANCE_URL,
  configPath: 'yp',
  maxLength: 255,
  placeholder: 'https://owncast.mysite.com',
  label: 'Instance URL',
  tip: 'Please provide the url to your Owncast site if you enable this Directory setting.',
};
// MISC FIELDS
export const FIELD_PROPS_TAGS = {
  apiPath: API_TAGS,
  configPath: 'instanceDetails',
  maxLength: 24,
  placeholder: 'Add a new tag',
  required: true,
  label: '',
  tip: '',
};

export const FIELD_PROPS_CUSTOM_CONTENT = {
  apiPath: API_CUSTOM_CONTENT,
  configPath: 'instanceDetails',
  placeholder: '',
  label: 'Extra page content',
  tip: 'Custom markup about yourself',
};
export const FIELD_PROPS_NSFW = {
  apiPath: API_NSFW_SWITCH,
  configPath: 'instanceDetails',
  label: 'NSFW?',
  tip: "Turn this ON if you plan to steam explicit or adult content. You may want to respectfully set this flag so that unexpecting eyes won't accidentally see it from the Directory.",
};

export const FIELD_PROPS_YP = {
  apiPath: API_YP_SWITCH,
  configPath: 'yp',
  label: 'Display in the Owncast Directory?',
  tip: 'Turn this ON if you want to show up in the Owncast directory at https://directory.owncast.online.',
};

export const ENCODER_PRESETS = [
  'fast',
  'faster',
  'veryfast',
  'superfast',
  'ultrafast',
];

export const DEFAULT_VARIANT_STATE:VideoVariant = {
  framerate: 24,
  videoPassthrough: false,
  videoBitrate: 800,
  audioPassthrough: true, // if false, then CAN set audiobitrate
  audioBitrate: 0,
  encoderPreset: 'veryfast',
};

export const DEFAULT_SOCIAL_HANDLE:SocialHandle = {
  url: '',
  platform: '',
};

export const OTHER_SOCIAL_HANDLE_OPTION = 'OTHER_SOCIAL_HANDLE_OPTION';
