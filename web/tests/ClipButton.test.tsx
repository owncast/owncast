import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { ClipButton } from '../components/action-buttons/ClipButton';

jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({
    t: (key: string) => (key === 'Frontend.Clips.end' ? 'End clip' : key),
  }),
}));

describe('ClipButton', () => {
  it('shows the full configured countdown when clipping is active', () => {
    render(<ClipButton active remainingSeconds={30} />);

    expect(screen.getByRole('button')).toHaveTextContent('End clip (0:30)');
  });

  it('formats minute-long countdowns without losing seconds', () => {
    render(<ClipButton active remainingSeconds={65} />);

    expect(screen.getByRole('button')).toHaveTextContent('End clip (1:05)');
  });
});
