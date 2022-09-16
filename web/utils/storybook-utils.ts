/**
 * Returns an object that can be used to disable a value in a Storybook story's "Controls" panel.
 * To use: spread the returned object into the 'argTypes' property of story's ComponentMeta.
 *
 * @param controls to be hidden. Can be a single control, or many.
 * @returns an object that should be spread into argTypes.
 *
 * @example `{ argTypes: { ...disableControls('foo', 'bar') } }`
 */
export const disableControls = (...controls: string[]) =>
	controls
		.map((control) => ({
			[control]: {
				control: false,
			},
		}))
		.reduce((prev, current) => ({ ...prev, ...current }), {});
