# Owncast Web Frontend

The Owncast web frontend is a [Next.js](https://nextjs.org/) project with [React](https://reactjs.org/) components, [TypeScript](https://www.typescriptlang.org/), [Sass](https://sass-lang.com/) styling, using [Ant Design](https://ant.design/) UI components.

## Getting Started

**First**, install the dependencies.

`npm install --include=dev`

## Components and Styles

You can start the [Storybook](https://storybook.js.org/) UI for exploring, testing, and developing components by running:

`npm run storybook`

This allows for components to be made available without the need of the server to be running and changes to be made in
isolation.

## Contribute

1. Find a component that hasn't yet been worked on by looking through the [UIv2 milestone](https://github.com/owncast/owncast/milestone/18) and the sidebar of components in storybook.
1. See if you can have an example of this functionality in action via the [Owncast Demo Server](https://watch.owncast.online) or [Owncast Nightly Build](https://nightly.owncast.online) so you know how it's supposed to work if it's interactive.
1. Visit the `Docs` tab to read any specific documentation that may have been written about how this component works.
1. Go to the `Canvas` tab of the component you selected and see if there's a Design attached to it.
1. If there is a design, then that's a starting point you can use to start building out the component.
1. If there isn't, then visit the [Owncast Demo Server](https://watch.owncast.online), the [Owncast Nightly Build](https://nightly.owncast.online), or the proposed [v2 design](https://www.figma.com/file/B6ICOn1J3dyYeoZM5kPM2A/Owncast---Review?node-id=0%3A1) for some ways to start.
1. If no design exists, then you can ask around the Owncast chat for help, or come up with your own ideas!
1. No designs are stuck in stone, and we're using this as an opportunity to level up the UI of Owncast, so all ideas are welcome.

See the extra how-to guide for components here: [Components How-to](./components/_COMPONENT_HOW_TO.md).

### Run the web project

Make sure you're running an instance of Owncast on localhost:8080, as your copy of the admin will look to use that as the API.

**Next**, start the web project with npm.

`npm run dev`

### Update the project

You can add or edit a pages by modifying `pages/something.js`. The page auto-updates as you edit the file.

[Routes](https://nextjs.org/docs/api-reference/next/router) will automatically be available for this new page components.

Since this project hits API endpoints you should make requests in [`componentDidMount`](https://reactjs.org/docs/react-component.html#componentdidmount), and not in [`getStaticProps`](https://nextjs.org/docs/basic-features/data-fetching), since they're not static and we don't want to fetch them at build time, but instead at runtime.

A list of API end points can be found here:
https://owncast.online/api/development/

### Admin Authentication

The pages under `/admin` require authentication to make API calls.

Auth: HTTP Basic

username: admin

pw: [your streamkey]

### Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.
