# How we develop components

This document outlines how we develop the components for the Owncast Web UI.

You should use this document as a guide when making changes to existing components, and adding new ones.
Working with the same development process help keep the project maintainable.

## What are components

A component in React is a custom HTML element. They're included in the DOM just like regular elements `<ChatBox /`>.

## Functional Components

In react, there's two ways to write a component: there's Class-based Components, and Functional Components.

Class-based is older and has fallen out of favor.
Functional Components are the new standard and you'll find them in most React projects written today.

See the [React Functional Component docs](https://reactjs.org/docs/components-and-props.html) for more info.

### How we write Functional Components

We've defined a pattern for how we write Functional Components in the Owncast Web UI.
There's a few ways to to write Functional Components that are common, so defining a standard helps keep this project readable and consistent.

The pattern we've settled on is:

**For stateless components:**

```tsx
export type MyNewButtonProps = {
  label: string;
  onClick: () => void;
};

export const MyNewButton: FC<MyNewButtonProps> = ({ label, onClick }) => (
  <button onClick={onCLick}>{label}</button>
);
```

**For stateful components:**

```tsx
export type MyNewButtonProps = {
  label: string;
  onClick: () => void;
};

export const MyNewButton: FC<MyNewButtonProps> = ({ label, onClick }) => {
  // do something, then call the onClick fn. e.g.:
  const handleClick = useCallback(() => {
    alert(label);
    onClick && onClick();
  }, [label, onClick]);

  return <button onClick={onCLick}>{label}</button>;
};
```

### Rationale

Since there's a lot of ways to create components, settling on one pattern helps maintain readability.
But why _this_ style?

See the discussion on the PR that introduced this pattern: [#2082](https://github.com/owncast/owncast/pull/2082).

## Error Boundaries

Components that have substantial state and internal functionality should be wrapped in an [Error Boundary](https://reactjs.org/docs/error-boundaries.html). This allows for catching unexpected errors and displaying a fallback UI.

Components that are stateless views are unlikely to throw exceptions and don't require an error boundary.

The `ComponentError` component is a pre-built error state that can be used to display an error message and a bug reporting button.

### Example

```tsx
import { ErrorBoundary } from 'react-error-boundary';

<ErrorBoundary
  fallbackRender={({ error, resetErrorBoundary }) => (
    <ComponentError
      componentName="DesktopContent"
      message={error.message}
      retryFunction={resetErrorBoundary}
    />
  )}
>
  <YourComponent />
</ErrorBoundary>
```

## Storybook

We use [Storybook](https://storybook.js.org/) to create a component library where we can see and interact with each component.

Make sure to include a `.stories.tsx` file with each (exported) component you create, and to update the stories file when making changes to existing components.

You can run the Storybook server locally with `npm run storybook`.
