export default {
  GOAL_BUFFER_LENGTH: 30,
  MAX_GOAL_BUFFER_LENGTH: 60,
  BACK_BUFFER_LENGTH: 30,
  GOAL_BUFFER_LENGTH_RATE: 1,
  // 0.5 MB/s
  INITIAL_BANDWIDTH: 4194304,
  // A fudge factor to apply to advertised playlist bitrates to account for
  // temporary flucations in client bandwidth
  BANDWIDTH_VARIANCE: 1.2,
  // How much of the buffer must be filled before we consider upswitching
  BUFFER_LOW_WATER_LINE: 0,
  MAX_BUFFER_LOW_WATER_LINE: 30,
  BUFFER_LOW_WATER_LINE_RATE: 1
};
