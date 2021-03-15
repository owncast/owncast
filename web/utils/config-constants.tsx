// DEFAULT VALUES
import { fetchData, SERVER_CONFIG_UPDATE_URL } from './apis';
import { ApiPostArgs, VideoVariant, SocialHandle } from '../types/config-section';

export const TEXT_MAXLENGTH = 255;

export const RESET_TIMEOUT = 3000;

// CONFIG API ENDPOINTS
export const API_CUSTOM_CONTENT = '/pagecontent';
export const API_FFMPEG = '/ffmpegpath';
export const API_INSTANCE_URL = '/serverurl';
export const API_LOGO = '/logo';
export const API_NSFW_SWITCH = '/nsfw';
export const API_RTMP_PORT = '/rtmpserverport';
export const API_S3_INFO = '/s3';
export const API_SERVER_SUMMARY = '/serversummary';
export const API_SERVER_WELCOME_MESSAGE = '/welcomemessage';
export const API_SERVER_NAME = '/name';
export const API_SOCIAL_HANDLES = '/socialhandles';
export const API_STREAM_KEY = '/key';
export const API_STREAM_TITLE = '/streamtitle';
export const API_TAGS = '/tags';
export const API_USERNAME = '/name';
export const API_VIDEO_SEGMENTS = '/video/streamlatencylevel';
export const API_VIDEO_VARIANTS = '/video/streamoutputvariants';
export const API_WEB_PORT = '/webserverport';
export const API_YP_SWITCH = '/directoryenabled';
export const API_CHAT_DISABLE = '/chat/disable';
export const API_EXTERNAL_ACTIONS = '/externalactions'

export async function postConfigUpdateToAPI(args: ApiPostArgs) {
  const { apiPath, data, onSuccess, onError } = args;
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
export const TEXTFIELD_PROPS_SERVER_NAME = {
  apiPath: API_SERVER_NAME,
  maxLength: TEXT_MAXLENGTH,
  placeholder: 'Owncast site name', // like "gothland"
  label: 'Name',
  tip: 'The name of your Owncast server',
};
export const TEXTFIELD_PROPS_STREAM_TITLE = {
  apiPath: API_STREAM_TITLE,
  maxLength: 100,
  placeholder: 'Doing cool things...',
  label: 'Stream Title',
  tip: 'What is your stream about today?',
};
export const TEXTFIELD_PROPS_SERVER_SUMMARY = {
  apiPath: API_SERVER_SUMMARY,
  maxLength: 500,
  placeholder: '',
  label: 'About',
  tip: 'A brief blurb about you, your server, or what your stream is about.',
};
export const TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE = {
  apiPath: API_SERVER_WELCOME_MESSAGE,
  maxLength: 500,
  placeholder: '',
  label: 'Welcome Message',
  tip:
    'A system chat message sent to viewers when they first connect to chat. Leave blank to disable.',
};
export const TEXTFIELD_PROPS_LOGO = {
  apiPath: API_LOGO,
  maxLength: 255,
  placeholder: '/img/mylogo.png',
  label: 'Logo',
  tip:
    'Name of your logo in the data directory. We recommend that you use a square image that is at least 256x256.',
};
export const TEXTFIELD_PROPS_STREAM_KEY = {
  apiPath: API_STREAM_KEY,
  configPath: '',
  maxLength: TEXT_MAXLENGTH,
  placeholder: 'abc123',
  label: 'Stream Key',
  tip: 'Save this key somewhere safe, you will need it to stream or login to the admin dashboard!',
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
  label: 'Owncast port',
  tip: 'What port is your Owncast web server listening? Default is 8080',
  required: true,
};
export const TEXTFIELD_PROPS_RTMP_PORT = {
  apiPath: API_RTMP_PORT,
  configPath: '',
  maxLength: 6,
  placeholder: '1935',
  label: 'RTMP port',
  tip: 'What port should accept inbound broadcasts? Default is 1935',
  required: true,
};
export const TEXTFIELD_PROPS_INSTANCE_URL = {
  apiPath: API_INSTANCE_URL,
  configPath: 'yp',
  maxLength: 255,
  placeholder: 'https://owncast.mysite.com',
  label: 'Server URL',
  tip: 'The full url to your Owncast server.',
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
  tip:
    "Turn this ON if you plan to steam explicit or adult content. Please respectfully set this flag so unexpected eyes won't accidentally see it in the Directory.",
};

export const FIELD_PROPS_YP = {
  apiPath: API_YP_SWITCH,
  configPath: 'yp',
  label: 'Enable directory',
  tip: 'Turn this ON if you want to show up in the directory.',
};

export const DEFAULT_VARIANT_STATE: VideoVariant = {
  framerate: 24,
  videoPassthrough: false,
  videoBitrate: 800,
  audioPassthrough: true, // if false, then CAN set audiobitrate
  audioBitrate: 0,
  cpuUsageLevel: 3,
  scaledHeight: null,
  scaledWidth: null,
  name: '',
};

export const FIELD_PROPS_DISABLE_CHAT = {
  apiPath: API_CHAT_DISABLE,
  configPath: 'chatDisabled',
  label: 'Disable chat',
  tip: 'Hide the web chat interface.',
  useSubmit: true,
};

export const VIDEO_VARIANT_SETTING_DEFAULTS = {
  // this one is currently unused
  audioBitrate: {
    min: 600,
    max: 1200,
    defaultValue: 800,
    unit: 'kbps',
    incrementBy: 100,
    tip: 'nothing to see here',
  },
  videoPassthrough: {
    tip: 'If enabled, all other settings will be disabled. Otherwise configure as desired.',
  },
  audioPassthrough: {
    tip: 'If No is selected, then you should set your desired Audio Bitrate.',
  },
  scaledWidth: {
    fieldName: 'scaledWidth',
    label: 'Resized Width',
    maxLength: 4,
    placeholder: '1080',
    tip: "Optionally resize this content's width.",
  },
  scaledHeight: {
    fieldName: 'scaledHeight',
    label: 'Resized Height',
    maxLength: 4,
    placeholder: '720',
    tip: "Optionally resize this content's height.",
  },
};

// VIDEO VARIANT FORM - framerate
export const FRAMERATE_DEFAULTS = {
  min: 24,
  max: 120,
  defaultValue: 24,
  unit: 'fps',
  incrementBy: null,
  tip:
    'Reducing your framerate will decrease the amount of video that needs to be encoded and sent to your viewers, saving CPU and bandwidth at the expense of smoothness.  A lower value is generally is fine for most content.',
};
export const FRAMERATE_SLIDER_MARKS = {
  [FRAMERATE_DEFAULTS.min]: `${FRAMERATE_DEFAULTS.min} ${FRAMERATE_DEFAULTS.unit}`,
  30: '',
  60: '',
  90: '',
  [FRAMERATE_DEFAULTS.max]: `${FRAMERATE_DEFAULTS.max} ${FRAMERATE_DEFAULTS.unit}`,
};
export const FRAMERATE_TOOLTIPS = {
  [FRAMERATE_DEFAULTS.min]: `${FRAMERATE_DEFAULTS.min}fps - Good for film, presentations, music, low power/bandwidth servers.`,
  30: '30fps - Good for slow/casual games, chat, general purpose.',
  60: '60fps - Good for fast/action games, sports, HD video.',
  90: '90fps - Good for newer fast games and hardware.',
  [FRAMERATE_DEFAULTS.max]: `${FRAMERATE_DEFAULTS.max}fps - Experimental, use at your own risk!`,
};
// VIDEO VARIANT FORM - bitrate
export const VIDEO_BITRATE_DEFAULTS = {
  min: 600,
  max: 6000,
  defaultValue: 1200,
  unit: 'kbps',
  incrementBy: 100,
  tip: 'The overall quality of your stream is generally impacted most by bitrate.',
};

export const VIDEO_NAME_DEFAULTS = {
  fieldName: 'name',
  label: 'Name',
  maxLength: 12,
  placeholder: 'HD or Low',
  tip: 'Human-readable name for for displaying in the quality selector.',
};

export const VIDEO_BITRATE_SLIDER_MARKS = {
  [VIDEO_BITRATE_DEFAULTS.min]: `${VIDEO_BITRATE_DEFAULTS.min} ${VIDEO_BITRATE_DEFAULTS.unit}`,
  3000: 3000,
  4500: 4500,
  [VIDEO_BITRATE_DEFAULTS.max]: `${VIDEO_BITRATE_DEFAULTS.max} ${VIDEO_BITRATE_DEFAULTS.unit}`,
};
// VIDEO VARIANT FORM - encoder preset
// CPU
export const ENCODER_PRESET_SLIDER_MARKS = {
  1: 'lowest',
  2: '',
  3: '',
  4: '',
  5: 'highest',
};
export const ENCODER_PRESET_TOOLTIPS = {
  1: 'Lowest CPU usage - lowest quality video',
  2: 'Low CPU usage - low quality video',
  3: 'Medium CPU usage - average quality video',
  4: 'High CPU usage - high quality video',
  5: 'Highest CPU usage - higher quality video',
};

export const DEFAULT_SOCIAL_HANDLE: SocialHandle = {
  url: '',
  platform: '',
};

export const OTHER_SOCIAL_HANDLE_OPTION = 'OTHER_SOCIAL_HANDLE_OPTION';

export const TEXTFIELD_PROPS_S3_COMMON = {
  maxLength: 255,
};

export const S3_TEXT_FIELDS_INFO = {
  accessKey: {
    fieldName: 'accessKey',
    label: 'Access Key',
    maxLength: 255,
    placeholder: 'access key 123',
    tip: '',
  },
  acl: {
    fieldName: 'acl',
    label: 'ACL',
    maxLength: 255,
    placeholder: '',
    tip: 'Optional specific access control value to add to your content.  Generally not required.',
  },
  bucket: {
    fieldName: 'bucket',
    label: 'Bucket',
    maxLength: 255,
    placeholder: 'bucket 123',
    tip: 'Create a new bucket for each Owncast instance you may be running.',
  },
  endpoint: {
    fieldName: 'endpoint',
    label: 'Endpoint',
    maxLength: 255,
    placeholder: 'https://your.s3.provider.endpoint.com',
    tip: 'The full URL endpoint your storage provider gave you.',
  },
  region: {
    fieldName: 'region',
    label: 'Region',
    maxLength: 255,
    placeholder: 'region 123',
    tip: '',
  },
  secret: {
    fieldName: 'secret',
    label: 'Secret key',
    maxLength: 255,
    placeholder: 'your secret key',
    tip: '',
  },
  servingEndpoint: {
    fieldName: 'servingEndpoint',
    label: 'Serving Endpoint',
    maxLength: 255,
    placeholder: 'http://cdn.ss3.provider.endpoint.com',
    tip:
      'Optional URL that content should be accessed from instead of the default.  Used with CDNs and specific storage providers. Generally not required.',
  },
};
