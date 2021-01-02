// DEFAULT VALUES
import React from 'react';
import { CheckCircleFilled, ExclamationCircleFilled } from '@ant-design/icons';

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


// Creating this so that it'll be easier to change values in one place, rather than looking for places to change it in a sea of JSX.

// key is the input's `fieldName`

export const TEXTFIELD_DEFAULTS = {
  name: {
    apiPath: '/name',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Server name',
    tip: 'This is your name that shows up on things and stuff.',
  },

  summary: {
    apiPath: '/serversummary',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Summary',
    tip: 'A brief blurb about what your stream is about.',
  },
  title: {
    apiPath: '/servertitle',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Server Title',
    tip: 'A brief blurb about what your stream is about.',
  },
  streamTitle: {
    apiPath: '/streamtitle',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Stream Title',
    tip: 'The name of your stream today.',
  },

  logo: {
    apiPath: '/logo',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Stream Title',
    tip: 'A brief blurb about what your stream is about.',
  },

  extraPageContent: {
    apiPath: '/pagecontent',
    placeholder: '',
    configPath: 'instanceDetails',
    label: 'Stream Key',
    tip: 'Custom markup about yourself',
  },


  streamKey: {
    apiPath: '/key',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: '',
    label: 'Stream Key',
    tip: 'Secret stream key',
  },
  ffmpegPath: {
    // apiPath: '/key',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: '',
    label: 'FFmpeg Path',
    tip: 'Absolute file path of the FFMPEG application on your server',
  },
  webServerPort: {
    apiPath: '/port',
    defaultValue: '',
    maxLength: 6,
    placeholder: DEFAULT_NAME,
    configPath: '',
    label: 'Server port',
    tip: 'What port are you serving Owncast from? Default is :8080',
  },

  // 
  tags: {
    apiPath: '/tags',
    defaultValue: '',
    maxLength: 24,
    placeholder: 'Add a new tag',
    configPath: 'instanceDetails',
    label: '',
    tip: '',
  }
}

