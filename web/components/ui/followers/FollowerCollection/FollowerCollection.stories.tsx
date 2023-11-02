import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { action } from '@storybook/addon-actions';
import { FollowerCollection } from './FollowerCollection';

const mocks = {
  mocks: [
    {
      // The "matcher" determines if this
      // mock should respond to the current
      // call to fetch().
      matcher: {
        name: 'response',
        url: 'glob:/api/followers*',
      },
      // If the "matcher" matches the current
      // fetch() call, the fetch response is
      // built using this "response".
      response: {
        status: 200,
        body: {
          total: 100,
          results: [
            {
              link: 'https://sun.minuscule.space/users/mardijker',
              name: 'mardijker',
              username: 'mardijker@sun.minuscule.space',
              image:
                'https://sun.minuscule.space/media/336af7ae5a2bcb508308eddb30b661ee2b2e15004a50795ee3ba0653ab190a93.jpg',
              timestamp: '2022-04-27T12:12:50Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.online/users/Kallegro',
              name: '',
              username: 'Kallegro@mastodon.online',
              image: '',
              timestamp: '2022-04-26T15:24:09Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.online/users/kerfuffle',
              name: 'Kerfuffle',
              username: 'kerfuffle@mastodon.online',
              image:
                'https://files.mastodon.online/accounts/avatars/000/133/698/original/6aa73caa898b2d36.gif',
              timestamp: '2022-04-25T21:32:41Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.uno/users/informapirata',
              name: 'informapirata :privacypride:',
              username: 'informapirata@mastodon.uno',
              image:
                'https://cdn.masto.host/mastodonuno/accounts/avatars/000/060/227/original/da4c44c716a339b8.png',
              timestamp: '2022-04-25T11:38:23Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/Raeanus',
              name: 'Raeanus',
              username: 'Raeanus@gamethattune.club',
              image:
                'https://gamethattune.club/media/a6e6ccea-34f8-4c2e-b9dc-ad8cca7fafd3/DD14E3BF-1358-4961-A900-42F3495F6BE2.jpeg',
              timestamp: '2022-04-23T00:46:56Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.ml/users/latte',
              name: 'Даниил',
              username: 'latte@mastodon.ml',
              image:
                'https://mastodon.ml/system/accounts/avatars/107/837/409/059/601/386/original/c45ec2676489e363.png',
              timestamp: '2022-04-19T13:06:09Z',
              disabledAt: null,
            },
            {
              link: 'https://wienermobile.rentals/users/jprjr',
              name: 'Johnny',
              username: 'jprjr@wienermobile.rentals',
              image: '',
              timestamp: '2022-04-14T14:48:11Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/johnny',
              name: 'John Regan',
              username: 'johnny@gamethattune.club',
              image:
                'https://gamethattune.club/media/3c10cd89-866b-4604-ae40-39387fe17061/profile_large.jpg',
              timestamp: '2022-04-14T14:42:48Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.social/users/MightyOwlbear',
              name: 'Haunted Owlbear',
              username: 'MightyOwlbear@mastodon.social',
              image:
                'https://files.mastodon.social/accounts/avatars/107/246/961/007/605/352/original/a86fc3db97a6de04.jpg',
              timestamp: '2022-04-14T13:33:03Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/thelinkfloyd',
              name: 'thelinkfloyd',
              username: 'thelinkfloyd@gamethattune.club',
              image: '',
              timestamp: '2022-04-05T12:23:32Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/TheBaffler',
              name: 'TheBaffler',
              username: 'TheBaffler@gamethattune.club',
              image: '',
              timestamp: '2022-04-04T19:50:08Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/Gttjessie',
              name: 'Gttjessie',
              username: 'Gttjessie@gamethattune.club',
              image: '',
              timestamp: '2022-03-30T20:18:47Z',
              disabledAt: null,
            },
            {
              link: 'https://cybre.space/users/fractal',
              name: 'Le fractal',
              username: 'fractal@cybre.space',
              image:
                'https://cybre.ams3.digitaloceanspaces.com/accounts/avatars/000/405/126/original/f1f2832a7bf1a967.png',
              timestamp: '2022-03-30T19:46:17Z',
              disabledAt: null,
            },
            {
              link: 'https://fosstodon.org/users/jumboshrimp',
              name: 'alex 👑🦐',
              username: 'jumboshrimp@fosstodon.org',
              image:
                'https://cdn.fosstodon.org/accounts/avatars/000/320/316/original/de43cda8653ade7f.jpg',
              timestamp: '2022-03-30T18:09:54Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/nvrslep303',
              name: 'Tay',
              username: 'nvrslep303@gamethattune.club',
              image: 'https://gamethattune.club/media/5cf9bc27-8821-445a-86ce-8aa3704acf2d/pfp.jpg',
              timestamp: '2022-03-30T15:27:49Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/anKerrigan',
              name: 'anKerrigan',
              username: 'anKerrigan@gamethattune.club',
              image: '',
              timestamp: '2022-03-30T14:47:04Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/jgangsta187',
              name: 'jgangsta187',
              username: 'jgangsta187@gamethattune.club',
              image: '',
              timestamp: '2022-03-30T14:42:52Z',
              disabledAt: null,
            },
            {
              link: 'https://gamethattune.club/users/aekre',
              name: 'aekre',
              username: 'aekre@gamethattune.club',
              image: '',
              timestamp: '2022-03-30T14:41:32Z',
              disabledAt: null,
            },
            {
              link: 'https://fosstodon.org/users/owncast',
              name: 'Owncast',
              username: 'owncast@fosstodon.org',
              image:
                'https://cdn.fosstodon.org/accounts/avatars/107/017/218/425/829/465/original/f98ba4cd61f483ab.png',
              timestamp: '2022-03-29T21:38:02Z',
              disabledAt: null,
            },
            {
              link: 'https://cybre.space/users/wklew',
              name: 'wally',
              username: 'wklew@cybre.space',
              image:
                'https://cybre.ams3.digitaloceanspaces.com/accounts/avatars/000/308/727/original/7453e74f3e09b27b.jpg',
              timestamp: '2022-03-29T18:24:29Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.social/users/nvrslep303',
              name: 'Tay',
              username: 'nvrslep303@mastodon.social',
              image:
                'https://files.mastodon.social/accounts/avatars/108/041/196/166/285/851/original/fc444dd6096381af.jpg',
              timestamp: '2022-03-29T18:19:31Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.social/users/morky',
              name: '',
              username: 'morky@mastodon.social',
              image: '',
              timestamp: '2022-03-29T18:17:59Z',
              disabledAt: null,
            },
            {
              link: 'https://mastodon.social/users/jgangsta187',
              name: 'John H.',
              username: 'jgangsta187@mastodon.social',
              image: '',
              timestamp: '2022-03-29T18:15:48Z',
              disabledAt: null,
            },
            {
              link: 'https://fosstodon.org/users/meisam',
              name: 'Meisam 🇪🇺:archlinux:',
              username: 'meisam@fosstodon.org',
              image:
                'https://cdn.fosstodon.org/accounts/avatars/000/264/096/original/54b4e6db97206bda.jpg',
              timestamp: '2022-03-29T18:12:21Z',
              disabledAt: null,
            },
          ],
        },
      },
    },
  ],
};

const noFollowersMock = {
  mocks: [
    {
      // The "matcher" determines if this
      // mock should respond to the current
      // call to fetch().
      matcher: {
        name: 'response',
        url: 'glob:/api/followers*',
      },
      // If the "matcher" matches the current
      // fetch() call, the fetch response is
      // built using this "response".
      response: {
        status: 200,
        body: {
          total: 0,
          results: [],
        },
      },
    },
  ],
};

const meta = {
  title: 'owncast/Components/Followers/Followers collection',
  component: FollowerCollection,
  parameters: {
    chromatic: { diffThreshold: 0.86 },
  },
} satisfies Meta<typeof FollowerCollection>;

export default meta;

const Template: StoryFn<typeof FollowerCollection> = (args: object) => (
  <RecoilRoot>
    <FollowerCollection
      onFollowButtonClick={() => {
        action('Follow button clicked');
      }}
      name="Example stream name"
      {...args}
    />
  </RecoilRoot>
);

export const NoFollowers = {
  render: Template,

  parameters: {
    fetchMock: noFollowersMock,
  },
};

export const Example = {
  render: Template,

  parameters: {
    fetchMock: mocks,
  },
};
