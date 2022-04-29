import { useState } from 'react';

interface Props {}

export default function ChatTextField(props: Props) {
  const [value, setValue] = useState('');
  const [showEmojis, setShowEmojis] = useState(false);

  return (
    <div>
      <input
        type="text"
        value={value}
        onChange={e => setValue(e.target.value)}
        placeholder="Type a message here then hit ENTER"
      />
      <button type="button" onClick={() => setShowEmojis(!showEmojis)}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="icon"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      </button>
    </div>
  );
}
