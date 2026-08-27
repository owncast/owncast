import { render, screen } from '@testing-library/react';
import { EmojiPicker } from '../components/chat/ChatTextField/EmojiPicker';

jest.mock('@emoji-mart/react', () => ({
  __esModule: true,
  default: ({ custom, categories }) => (
    <output data-categories={String(categories)}>{JSON.stringify(custom)}</output>
  ),
}));

beforeEach(() => {
  Object.defineProperty(document, 'querySelector', {
    configurable: true,
    value: jest.fn(() => ({ shadowRoot: { adoptedStyleSheets: [] } })),
  });
  global.CSSStyleSheet = class {
    replaceSync = jest.fn();
  } as unknown as typeof CSSStyleSheet;
});

test('updates custom emoji after the emoji API response arrives', () => {
  const { rerender } = render(<EmojiPicker customEmoji={[]} onEmojiSelect={jest.fn()} />);

  rerender(
    <EmojiPicker
      customEmoji={[{ name: 'owncat', url: '/img/emoji/owncat.png' }]}
      onEmojiSelect={jest.fn()}
    />,
  );

  expect(screen.getByText(/owncat/)).toHaveTextContent('"id":"owncat"');
  expect(screen.getByText(/owncat/)).toHaveAttribute('data-categories', 'undefined');
});
