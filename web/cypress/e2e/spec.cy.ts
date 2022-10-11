/* eslint-disable no-undef */
describe('Basic tests', () => {
  it('Can visit the page', () => {
    cy.visit('http://localhost:3030/');
  });

  it('Can change to Followers tab', () => {
    cy.contains('Followers').click();
  });

  it('Can open notify modal', () => {
    cy.contains('Notify when live').click();
  });

  it('Can close notify modal', () => {
    cy.get('.ant-modal-close').click();
  });

  it('Can open follow modal', () => {
    cy.contains('Follow').click();
  });

  it('Can close follow modal', () => {
    cy.get('.ant-modal-close').click();
  });

  // Verify content header values
  it('Has correct content header values', () => {
    cy.get('.header-title').should('have.text', 'Owncast');
    cy.get('.header-subtitle').should(
      'have.text',
      'Welcome to your new Owncast server! This description can be changed in the admin. Visit https://owncast.online/docs/configuration/ to learn more.',
    );
  });

  it('Has correct global header values', () => {
    cy.get('.global-header-text').should('have.text', 'Owncast');
  });

  // Offline banner
  it('Has correct offline banner values', () => {
    cy.contains('This stream is offline. Be notified the next time Owncast goes live.').should(
      'be.visible',
    );
  });

  // Verify the tags show up
  it('Has correct tags visible', () => {
    cy.contains('#owncast').should('be.visible');
    cy.contains('#streaming').should('be.visible');
  });
});
