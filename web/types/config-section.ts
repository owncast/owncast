// TS types for elements on the Config pages

export interface TextFieldProps {
  handleResetValue?: (fieldName: string) => void;
  fieldName: string;
  initialValues?: any;
  type?: string;
  configPath?: string;
  required?: boolean;
  disabled?: boolean;
  onSubmit?: () => void;
}

export interface ToggleSwitchProps {
  fieldName: string;
  initialValues?: any;
  configPath?: string;
  disabled?: boolean;
}

export interface UpdateArgs {
  fieldName: string;
  value: string;
  path?: string;
}

export interface ApiPostArgs {
  apiPath: string,
  data: object,
  onSuccess?: (arg: any) => void,
  onError?: (arg: any) => void,
}

export interface ConfigDirectoryFields {
  enabled: boolean;
  instanceUrl: string,
}

export interface ConfigInstanceDetailsFields {
  extraPageContent: string;
  logo: string;
  name: string;
  nsfw: boolean;
  streamTitle: string;
  summary: string;
  tags: string[];
  title: string;
}

export interface VideoVariant {
  audioBitrate: number;
  audioPassthrough: number;
  encoderPreset: 'ultrafast' | 'superfast' | 'veryfast' | 'faster' | 'fast';
  framerate: number;
  videoBitrate: number;
  videoPassthrough: boolean;
}
export interface VideoSettingsFields {
  numberOfPlaylistItems: number;
  segmentLengthSeconds: number;
  videoQualityVariants: VideoVariant[],
}


export interface ConfigDetails {
  ffmpegPath: string;
  instanceDetails: ConfigInstanceDetailsFields;
  rtmpServerPort: string;
  s3: any; // tbd
  streamKey: string;
  webServerPort: string;
  yp: ConfigDirectoryFields;
  videoSettings: VideoSettingsFields;
}
