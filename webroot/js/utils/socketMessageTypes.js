/**
 * These are the types of messages that we can handle with the websocket.
 * Mostly used by `websocket.js` but if other components need to handle
 * different types then it can import this file.
 */
export default {
  CHAT: 'CHAT',
  PING: 'PING',
  NAME_CHANGE: 'NAME_CHANGE',
  PONG: 'PONG'
};
