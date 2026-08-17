import PlaybackMetrics from '../components/video/metrics/playback';

type Listener = () => void;

const makePlayer = () => {
  const listeners = new Map<string, Listener>();
  return {
    buffered: () => [],
    currentTime: () => 0,
    ended: () => false,
    error: () => null,
    off: jest.fn(),
    on: jest.fn((event: string, listener: Listener) => listeners.set(event, listener)),
    paused: () => false,
    playbackRate: () => 1,
    seeking: () => false,
    tech: () => undefined,
    trigger: (event: string) => listeners.get(event)?.(),
  };
};

describe('PlaybackMetrics', () => {
  beforeEach(() => {
    global.fetch = jest.fn(() => Promise.resolve({})) as jest.Mock;
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  test('reports ended when a started player is torn down', async () => {
    jest.useFakeTimers();
    const player = makePlayer();
    const videojs = { Vhs: {}, xhr: jest.fn() };
    const metrics = new PlaybackMetrics(player, videojs);

    player.trigger('playing');
    await Promise.resolve();
    metrics.stop();
    await Promise.resolve();

    const reports = (global.fetch as jest.Mock).mock.calls.map(([, init]) => JSON.parse(init.body));
    expect(reports.map(report => report.sta)).toEqual(['p', 'e']);
  });

  test('reports a pause after playback starts', async () => {
    jest.useFakeTimers();
    const player = makePlayer();
    const metrics = new PlaybackMetrics(player, { Vhs: {}, xhr: jest.fn() });
    expect(metrics).toBeDefined();

    player.trigger('playing');
    player.trigger('pause');
    await Promise.resolve();

    const reports = (global.fetch as jest.Mock).mock.calls.map(([, init]) => JSON.parse(init.body));
    expect(reports.map(report => report.sta)).toEqual(['p', 'a']);
  });
});
