import { Button as AntButton, ButtonProps as AntButtonProps } from 'antd';
import { FC } from 'react';

import styles from './Button.module.scss';

export type ButtonProps = Omit<AntButtonProps, 'children' | 'icon' | 'title' | 'target'> & {
	text: string;
	iconSrc?: string;
	iconAltText?: string;
};

/**
 * A customised and styled Ant Design button.
 *
 * An icon can be set. Both `iconSrc` and `iconAltText` must be set for the icon to be shown.
 */
export const Button: FC<ButtonProps> = ({ text, iconSrc, iconAltText, type = "primary", ...props }) => (
	<AntButton
		type={type}
		className={styles.button}
		{...props}
	>
		{iconSrc && iconAltText !== undefined && (
			<img src={iconSrc} className={`${styles.icon}`} alt={iconAltText}/>
		)}
		{text}
	</AntButton>
);
