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

## Using Ant's `<Form>` with `form-textfield`.
You may see that a couple of pages (currently Public Details and Server Details page), is mainly a grouping of similar Text fields. 

`const [form] = Form.useForm();`
`form.setFieldsValue(initialValues);`


A special `TextField` component was created to be used with form.


## Potential Optimizations

Looking back at the pages with `<Form>` + `form-textfield`, t

This pattern might be overly engineered.

There are also a few patterns across all the form groups that repeat quite a bit. Perhaps these patterns could be consolidated into a custom hook that could handle all the steps.

TODO: explain how to use <Form> and how the custom `form-xxxx` components work together.


## Misc notes
- `instanceDetails` needs to be filled out before `yp.enabled` can be turned on.



## Config data structure (with default values)
```
{
  streamKey: '',
  instanceDetails: {
    tags: [],
    nsfw: false,
  },
  yp: {
    enabled: false,
    instance: '',
  },
  videoSettings: {
    videoQualityVariants: [
      {
        audioPassthrough: false,
        videoPassthrough: false,
        videoBitrate: 0,
        audioBitrate: 0,
        framerate: 0,
      },
    ],
  }
};
```

TODO:
- fill out readme for how to use form fields and about data flow w/ ant and react

- cleanup 
  - more consitent constants
  - cleanup types
  - cleanup style sheets..? make style module for each config page? (but what about ant deisgn overrides?)
- redesign
  - label+form layout - put them into a table, table of rows
  - change Encpder preset into slider
  - page headers - diff color? 
  - fix social handles icon in table
  - consolidate things into 1 page?
  - things could use smaller font?
- maybe convert common form pattern to custom hook?


Possibly over engineered

https://uxcandy.co/demo/label_pro/preview/demo_2/pages/forms/form-elements.html

https://www.bootstrapdash.com/demo/corona/jquery/template/modern-vertical/pages/forms/basic_elements.html