# About the Config editing section

An adventure with React, React Hooks and Ant Design forms.

## General data flow in this React app

### First things to note
- When the Admin app loads, the `ServerStatusContext` (in addition to checking server `/status` on a timer) makes a call to the  `/serverconfig` API to get your config details. This data will be stored as **`serverConfig`** in app state, and _provided_ to the app via `useContext` hook.  

- The `serverConfig` in state is be the central source of data that pre-populates the forms.

- The `ServerStatusContext` also provides a method for components to update the serverConfig state, called `setFieldInConfigState()`.

- After you have updated a config value in a form field, and successfully submitted it through its endpoint, you should call `setFieldInConfigState` to update the global state with the new value.

- Each top field of the serverConfig has its own API update endpoint.

### Form Flow
Each form input (or group of inputs) you make, you should 
  1. Get the field values that you want out of `serverConfig` from ServerStatusContext with `useContext`.
  2. Next we'll have to put these field values of interest into a `useState` in each grouping.  This will help you edit the form. 
  3. Because ths config data is populated asynchronously,  Use a `useEffect` to check when that data has arrived before putting it into state.
  4. You will be using the state's value to populate the `defaultValue` and the `value` props of each Ant input component (`Input`, `Toggle`, `Switch`, `Select`, `Slider` are currently used).
  5. When an `onChange` event fires for each type of input component, you will update the local state of each page with the changed value.
  6. Depending on the form, an `onChange` of the input component, or a subsequent `onClick` of a submit button will take the value from local state and POST the field's API.
  7. `onSuccess` of the post, you should update the global app state with the new value.

There are also a variety of other local states to manage the display of error/success messaging.

## Notes about `form-textfield` and `form-togglefield`
- The text field is intentionally designed to make it difficult for the user to submit bad data.
- If you make a change on a field, a Submit buttton will show up that you have to click to update. That will be the only way you can update it.
- If you clear out a field that is marked as Required, then exit/blur the field, it will repopulate with its original value. 

- Both of these elements are specifically meant to be used with updating `serverConfig` fields, since each field requires its own endpoint.

- Give these fields a bunch of props, and they will display labelling, some helpful UI around tips, validation messaging, as well as submit the update for you.

- (currently undergoing re-styling and TS cleanup)

- NOTE: you don't have to use these components. Some form groups may require a customized UX flow where you're better off using the Ant components straight up.



