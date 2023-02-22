import { VideoSettingsStaticService, VideoQuality } from './video-settings-service';

export const videoSettingsServiceMockOf = (
  videoQualities: Array<VideoQuality>,
): VideoSettingsStaticService =>
  class VideoSettingsServiceMock {
    public static async getVideoQualities(): Promise<Array<VideoQuality>> {
      return videoQualities;
    }
  };

export default videoSettingsServiceMockOf;
