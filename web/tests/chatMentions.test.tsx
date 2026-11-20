import { act, fireEvent, render, screen } from '@testing-library/react';
import { createStore, Provider } from 'jotai';
import { TextEncoder } from 'node:util';
import WebsocketService from '../services/websocket-service';
import { ChatEvent, MessageType } from '../interfaces/socket-events';
import { ChatUserMessage } from '../components/chat/ChatUserMessage/ChatUserMessage';
import { ChatTextField } from '../components/chat/ChatTextField/ChatTextField';
import { websocketServiceAtom } from '../components/stores/ClientConfigStore';

jest.mock('next/dynamic', () => () => () => null);
jest.mock('sanitize-html', () => jest.fn());
jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

const originalFetch = global.fetch;
const originalTextEncoder = global.TextEncoder;
const originalExecCommand = document.execCommand;

beforeEach(() => {
  global.fetch = jest.fn().mockResolvedValue({ json: async () => [] });
  global.TextEncoder = TextEncoder;
  document.execCommand = jest.fn((_command, _showUI, html) => {
    const selection = document.getSelection();
    const range = selection.getRangeAt(0);
    const fragment = range.createContextualFragment(html);
    const lastNode = fragment.lastChild;
    range.deleteContents();
    range.insertNode(fragment);
    range.setStartAfter(lastNode);
    range.collapse(true);
    return true;
  });
});

afterEach(() => {
  global.fetch = originalFetch;
  global.TextEncoder = originalTextEncoder;
  document.execCommand = originalExecCommand;
});

function renderComposer() {
  const websocket = new WebsocketService('', '', '');
  const send = jest.spyOn(websocket, 'send').mockImplementation(() => {});
  const store = createStore();
  store.set(websocketServiceAtom, websocket);
  render(
    <Provider store={store}>
      <ChatTextField
        enabled
        focusInput={false}
        mentionCandidates={[
          { id: 'alice', displayName: 'Alice' },
          { id: 'artist', displayName: 'Alice & Co' },
        ]}
      />
    </Provider>,
  );
  return { send, editor: screen.getByRole('combobox') };
}

async function typeMessage(editor: HTMLElement, text: string) {
  await act(async () => {
    editor.focus();
    editor.textContent = text;
    document.getSelection().collapse(editor.firstChild, text.length);
  });
}

test('selects a participant by keyboard and sends the exact UTF-8 range', async () => {
  const { editor, send } = renderComposer();
  await typeMessage(editor, '👋 @ali');
  expect(screen.getAllByRole('option')).toHaveLength(2);

  fireEvent.keyDown(editor, { key: 'ArrowDown' });
  expect(screen.getByRole('option', { name: '@Alice & Co' })).toHaveAttribute(
    'aria-selected',
    'true',
  );
  await act(async () => {
    fireEvent.keyDown(editor, { key: 'Enter' });
  });
  expect(send).not.toHaveBeenCalled();

  await act(async () => {
    fireEvent.keyDown(editor, { key: 'Enter' });
  });
  expect(send).toHaveBeenCalledWith({
    type: 'CHAT',
    body: '👋 @Alice & Co',
    mentions: [{ userId: 'artist', start: 5, end: 16 }],
  });
});

test('sends an unselected name as ordinary text after dismissing suggestions', async () => {
  const { editor, send } = renderComposer();
  await typeMessage(editor, '@Alice');
  fireEvent.keyDown(editor, { key: 'Escape' });
  expect(screen.queryByRole('listbox')).not.toBeInTheDocument();

  await act(async () => {
    fireEvent.keyDown(editor, { key: 'Enter' });
  });
  expect(send).toHaveBeenCalledWith({ type: 'CHAT', body: '@Alice' });
});

test('shows the stored highlight only in the sender’s view', () => {
  const message = new ChatEvent({
    id: 'message',
    timestamp: new Date(),
    type: MessageType.CHAT,
    body: '<p>Hello <span data-chat-mention="artist">@Alice &amp; Co</span></p>',
    user: { id: 'sender', displayName: 'Sender', displayColor: 1, createdAt: new Date() },
  });
  const props = {
    message,
    highlightString: '',
    showModeratorMenu: false,
    sameUserAsLast: false,
    isAuthorModerator: false,
    isAuthorAuthenticated: false,
    isAuthorBot: false,
  };
  const { rerender } = render(<ChatUserMessage {...props} sentBySelf />);
  expect(screen.getByText('@Alice & Co').tagName).toBe('MARK');

  rerender(<ChatUserMessage {...props} sentBySelf={false} />);
  expect(screen.getByText('Hello @Alice & Co').tagName).not.toBe('MARK');

  rerender(<ChatUserMessage {...props} sentBySelf={false} highlightString="Alice & Co" />);
  expect(screen.getByText('Alice & Co').tagName).toBe('MARK');
});
