import OwncastPlayer from '../../../components/video/OwncastPlayer';

export default function VideoEmbed() {
  const online = false;
  return (
    <div className="video-embed">
      <OwncastPlayer source="/hls/stream.m3u8" online={online} />
    </div>
  );
}
