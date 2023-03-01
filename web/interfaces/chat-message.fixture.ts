import { ChatMessage } from './chat-message.model';
import { MessageType } from './socket-events';
import { spidermanUser, grootUser } from './user.fixture';
import { User } from './user.model';

export const createMessages = (
  basicMessages: Array<{ body: string; user: User }>,
): Array<ChatMessage> => {
  const baseDate = new Date(2022, 1, 3).valueOf();
  return basicMessages.map(
    ({ body, user }, index): ChatMessage => ({
      body,
      user,
      id: index.toString(),
      type: MessageType.CHAT,
      timestamp: new Date(baseDate + 1_000 * index),
    }),
  );
};

export const exampleChatHistory = createMessages([
  {
    body: 'So, how do you like my new suit?',
    user: spidermanUser,
  },
  {
    body: 'Im am Groot.',
    user: grootUser,
  },
  {
    body: 'Really? That bad?',
    user: spidermanUser,
  },
  {
    body: 'Im am Groot!',
    user: grootUser,
  },
  {
    body: 'But what about the new web slingers?',
    user: spidermanUser,
  },
  {
    body: 'Im am Groooooooooooooooot.',
    user: grootUser,
  },
  {
    body: "Ugh, come on, they aren't THAT big!",
    user: spidermanUser,
  },
  {
    body: 'I am Groot.',
    user: grootUser,
  },
  {
    body: "Fine then. I don't like your new leaves either!",
    user: spidermanUser,
  },
  {
    body: 'I AM GROOT!!!!!',
    user: grootUser,
  },
]);
