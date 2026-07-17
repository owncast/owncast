// Appearance-theme presets shared by the Theme playground story and the
// Storybook toolbar theme switcher (.storybook/preview.js). Each preset is a
// set of the same appearance variables an admin can set from the Appearance
// config page. They flow through the two real theming paths at once: the
// Theme component turns them into --theme-* CSS variables, and AntdProvider
// maps them onto Ant Design design tokens.
const THEMES: Record<string, { label: string; variables: Record<string, string> }> = {
  default: {
    label: 'Owncast default',
    variables: {},
  },
  midnight: {
    label: 'Midnight',
    variables: {
      'theme-color-action': '#22d3ee',
      'theme-color-action-hover': '#67e8f9',
      'theme-color-background-main': '#1f2937',
      'theme-color-background-header': '#0b1220',
      'theme-color-components-text-on-light': '#e5e7eb',
      'theme-color-components-form-field-background': '#111827',
      'theme-color-components-modal-header-background': '#0b1220',
      'theme-color-components-modal-header-text': '#f9fafb',
      'theme-color-components-modal-content-background': '#1f2937',
      'theme-color-components-menu-background': '#111827',
      'theme-color-components-menu-item-text': '#e5e7eb',
      'theme-rounded-corners': '4px',
    },
  },
  forest: {
    label: 'Forest',
    variables: {
      'theme-color-action': '#15803d',
      'theme-color-action-hover': '#22c55e',
      'theme-color-background-main': '#f0fdf4',
      'theme-color-background-header': '#14532d',
      'theme-color-components-text-on-light': '#14532d',
      'theme-color-components-form-field-background': '#dcfce7',
      'theme-color-components-modal-header-background': '#14532d',
      'theme-color-components-modal-header-text': '#f0fdf4',
      'theme-color-components-modal-content-background': '#f0fdf4',
      'theme-rounded-corners': '12px',
    },
  },
  terminal: {
    label: 'Terminal',
    variables: {
      'theme-color-action': '#00ff41',
      'theme-color-action-hover': '#7dff9b',
      'theme-color-background-main': '#0a0a0a',
      'theme-color-background-header': '#000000',
      'theme-color-components-text-on-light': '#00ff41',
      'theme-color-components-form-field-background': '#001a06',
      'theme-color-components-modal-header-background': '#000000',
      'theme-color-components-modal-header-text': '#00ff41',
      'theme-color-components-modal-content-background': '#0a0a0a',
      'theme-rounded-corners': '0px',
    },
  },
};

export default THEMES;
