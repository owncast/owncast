# Tips for creating a new Admin form

### Layout
- Give your page or form a title. Feel free to use Ant Design's `<Title>` component.
- Give your form a description inside of a `<p className="description" />` tag.

- Use some Ant Design `Row` and `Col`'s to layout your forms if you want to spread them out into responsive columns.

- Use the `form-module` CSS class if you want to add a visual separation to a grouping of items.
  


### Form fields
- Feel free to use the pre-styled `<TextField>` text form field or the `<ToggleSwitch>` compnent, in a group of form fields together. These have been styled and laid out to match each other.

- `Slider`'s - If your form uses an Ant Slider component, follow this recommended markup of CSS classes to maintain a consistent look and feel to other Sliders in the app. 
  ```
  <div className="segment-slider-container">
    <Slider ...props />
    <p className="selected-value-note">{selected value}</p>
  </div>
  ```

### Submit Statuses
- It would be nice to display indicators of success/warnings to let users know if something has been successfully updated on the server. It has a lot of steps (sorry, but it could probably be optimized), but it'll provide a consistent way to display messaging. 

  - See `reset-yp.tsx` for an example of using `submitStatus` with `useState()` and the `<FormStatusIndicator>` component to achieve this.

### Styling
- This admin site chooses to have a generally Dark color palette, but with colors that are  different from Ant design's _dark_ stylesheet, so that style sheet is not included. This results in a very large `ant-overrides.scss` file to reset colors on frequently used Ant components in the system. If you find yourself a new Ant Component that has not yet been used in this app, feel free to add a reset style for that component to the overrides stylesheet.

  - Take a look at `variables.scss` CSS file if you want to give some elements custom css colors.


---
---
# Creating Admin forms the Config section 
First things first..

## General Config data flow in this React app

- When the Admin app loads, the `ServerStatusContext` (in addition to checking server `/status` on a timer) makes a call to the  `/serverconfig` API to get your config details. This data will be stored as **`serverConfig`** in app state, and _provided_ to the app via `useContext` hook.  

- The `serverConfig` in state is be the central source of data that pre-populates the forms.

- The `ServerStatusContext` also provides a method for components to update the serverConfig state, called `setFieldInConfigState()`.

- After you have updated a config value in a form field, and successfully submitted it through its endpoint, you should call `setFieldInConfigState` to update the global state with the new value.


## Suggested Config Form Flow
- *NOTE: Each top field of the serverConfig has its own API update endpoint.*


There many steps here, but they are highly suggested to ensure that Config values are updated and displayed properly throughout the entire admin form.

For each form input (or group of inputs) you make, you should:
  1. Get the field values that you want out of `serverConfig` from ServerStatusContext with `useContext`.
  2. Next we'll have to put these field values of interest into a `useState` in each grouping.  This will help you edit the form. 
  3. Because ths config data is populated asynchronously,  Use a `useEffect` to check when that data has arrived before putting it into state.
  4. You will be using the state's value to populate the `defaultValue` and the `value` props of each Ant input component (`Input`, `Toggle`, `Switch`, `Select`, `Slider` are currently used).
  5. When an `onChange` event fires for each type of input component, you will update the local state of each page with the changed value.
  6. Depending on the form, an `onChange` of the input component, or a subsequent `onClick` of a submit button will take the value from local state and POST the field's API.
  7. `onSuccess` of the post, you should update the global app state with the new value.

There are also a variety of other local states to manage the display of error/success messaging.

- It is recommended that you use `form-textfield-with-submit` and `form-toggleswitch`(with `useSubmit=true`) Components to edit Config fields.

Examples of Config form groups where individual form fields submitting to the update API include:
- `edit-instance-details.tsx`
- `edit-server-details.tsx`

Examples of Config form groups where there is 1 submit button for the entire group include: 
- `edit-storage.tsx`


---
#### Notes about `form-textfield-with-submit` and `form-togglefield` (with useSubmit=true)
- The text field is intentionally designed to make it difficult for the user to submit bad data.
- If you make a change on a field, a Submit buttton will show up that you have to click to update. That will be the only way you can update it.
- If you clear out a field that is marked as Required, then exit/blur the field, it will repopulate with its original value. 

- Both of these elements are specifically meant to be used with updating `serverConfig` fields, since each field requires its own endpoint.

- Give these fields a bunch of props, and they will display labelling, some helpful UI around tips, validation messaging, as well as submit the update for you.

- (currently undergoing re-styling and TS cleanup)

- NOTE: you don't have to use these components. Some form groups may require a customized UX flow where you're better off using the Ant components straight up.

