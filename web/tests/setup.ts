import '@testing-library/jest-dom';

// Mock window.matchMedia for antd components
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// Mock ResizeObserver for antd v6 components (Popconfirm/Tooltip triggers
// observe their targets on mount); jsdom does not implement it.
class ResizeObserverMock {
  observe = jest.fn();

  unobserve = jest.fn();

  disconnect = jest.fn();
}

Object.defineProperty(window, 'ResizeObserver', {
  writable: true,
  value: ResizeObserverMock,
});

// Ant Design v6 Select delays closing its popup with MessageChannel. jsdom
// does not implement it, so provide the minimum asynchronous behavior needed
// to exercise keyboard selection.
const MessageChannelMock = jest.fn().mockImplementation(() => {
  const port1 = { onmessage: null as null | (() => void) };

  return {
    port1,
    port2: {
      postMessage: () => window.setTimeout(() => port1.onmessage?.(), 0),
    },
  };
});

Object.defineProperty(globalThis, 'MessageChannel', {
  writable: true,
  value: MessageChannelMock,
});
