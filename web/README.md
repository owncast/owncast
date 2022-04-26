# Owncast Web

## Owncast Web Frontend

The Owncast web frontend is a [Next.js](https://nextjs.org/) project with [React](https://reactjs.org/) components, [TypeScript](https://www.typescriptlang.org/), [Sass](https://sass-lang.com/) styling, using [Ant Design](https://ant.design/) UI components.

### Getting Started

**First**, install the dependencies.

```npm install```

### Run the web project

Make sure you're running an instance of Owncast on localhost:8080, as your copy of the admin will look to use that as the API.

**Next**, start the web project with npm.
  
  ```npm run dev```

### Update the project

You can add or edit a pages by modifying `pages/something.js`. The page auto-updates as you edit the file.

[Routes](https://nextjs.org/docs/api-reference/next/router) will automatically be available for this new page components.

Since this project hits API endpoints you should make requests in [`componentDidMount`](https://reactjs.org/docs/react-component.html#componentdidmount), and not in [`getStaticProps`](https://nextjs.org/docs/basic-features/data-fetching), since they're not static and we don't want to fetch them at build time, but instead at runtime.

A list of API end points can be found here:
https://owncast.online/api/development/

### Admin Authentication

The pages until `/admin` require authentication to make API calls.

Auth: HTTP Basic

username: admin

pw: [your streamkey]

### Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!

## Style guide and components

We are currently experimenting with using [Storybook](https://storybook.js.org/) to build components, experiment with styling, and have a single place to find colors, fonts, and other styles.

To work with Storybook:

```npm run storybook```