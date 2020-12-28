// DEFAULT VALUES

export const DEFAULT_NAME = 'Owncast User';
export const DEFAULT_TITLE = 'Owncast Server';
export const DEFAULT_SUMMARY = '';

export const TEXT_MAXLENGTH = 255;


// Creating this so that it'll be easier to change values in one place, rather than looking for places to change it in a sea of JSX.
// key is the input's `fieldName`
//serversummary
export const TEXTFIELD_DEFAULTS = {
  name: {
    apiPath: '/name',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Your name',
    tip: 'This is your name that shows up on things and stuff.',
  },

  summary: {
    apiPath: '/summary',
    defaultValue: DEFAULT_NAME,
    maxLength: TEXT_MAXLENGTH,
    placeholder: DEFAULT_NAME,
    configPath: 'instanceDetails',
    label: 'Summary',
    tip: 'A brief blurb about what your stream is about.',
  }
}