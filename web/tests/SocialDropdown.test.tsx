import { useState } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { SocialDropdown } from '../components/admin/SocialDropdown';

const iconList = [
  {
    key: 'bandcamp',
    platform: 'Bandcamp',
    icon: '/img/platformlogos/bandcamp.svg',
  },
  {
    key: 'mastodon',
    platform: 'Mastodon',
    icon: '/img/platformlogos/mastodon.svg',
  },
];

describe('SocialDropdown', () => {
  it('keeps keyboard selection and the popup within the social platform modal', async () => {
    const onSelected = jest.fn();

    const TestDropdown = () => {
      const [selectedOption, setSelectedOption] = useState('');

      return (
        <div className="ant-modal-container">
          <SocialDropdown
            iconList={iconList}
            selectedOption={selectedOption}
            onSelected={value => {
              setSelectedOption(value);
              onSelected(value);
            }}
          />
        </div>
      );
    };

    const { container } = render(<TestDropdown />);
    const modalContainer = container.querySelector('.ant-modal-container');
    const combobox = screen.getByRole('combobox', { name: 'Social Platform' });

    act(() => combobox.focus());
    fireEvent.keyDown(combobox, { key: 'ArrowDown', code: 'ArrowDown', keyCode: 40 });

    const listbox = await screen.findByRole('listbox');
    expect(modalContainer).toContainElement(listbox);

    fireEvent.keyDown(combobox, { key: 'Enter', code: 'Enter', keyCode: 13 });

    await waitFor(() => expect(onSelected).toHaveBeenCalledWith('bandcamp'));
    expect(combobox).toHaveFocus();
    expect(container.querySelector('.ant-select-content-value')).toHaveTextContent('Bandcamp');
  });
});
