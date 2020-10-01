This is a [Next.js](https://nextjs.org/) project with [TypeScript](https://www.typescriptlang.org/), [Sass](https://sass-lang.com/) styling, using [Ant Design](https://ant.design/) UI components.

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000/admin](http://localhost:3000/admin) with your browser to see the result.

You can start editing a page by modifying `pages/something.js`. The page auto-updates as you edit the file.

Add new pages by adding files to the `pages` directory and [routes](https://nextjs.org/docs/api-reference/next/router) will be available for this new page component.

Since this project hits API endpoints you should make requests in [`componentDidMount`](https://reactjs.org/docs/react-component.html#componentdidmount), and not in [`getStaticProps`](https://nextjs.org/docs/basic-features/data-fetching), since they're not static and we don't want to fetch them at build time, but instead at runtime.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!
