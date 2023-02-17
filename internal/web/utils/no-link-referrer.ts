/*
Due to Owncast's goal of being private by default, we don't want any external
links to leak the instance of Owncast as a referrer.
This observer attempts to catch any anchor tags and automatically add the
noopener and noreferrer attributes to them so the instance of Owncast isn't
passed along in the headers.

This should should be fired somewhere relatively high level in the DOM and live
for the entirety of the page.
*/

/* eslint-disable no-restricted-syntax */
export default function setupNoLinkReferrer(observationRoot: HTMLElement): void {
  // Options for the observer (which mutations to observe)
  const config = { attributes: false, childList: true, subtree: true };

  const addNoReferrer = (node: Element): void => {
    const existingAttributes = node.getAttribute('rel');
    const attributes = `${existingAttributes} noopener noreferrer`;
    node.setAttribute('rel', attributes);
  };

  // Callback function to execute when mutations are observed
  // eslint-disable-next-line func-names
  const callback = function (mutationList) {
    for (const mutation of mutationList) {
      for (const node of mutation.addedNodes) {
        // we track only elements, skip other nodes (e.g. text nodes)
        // eslint-disable-next-line no-continue
        if (!(node instanceof HTMLElement)) continue;

        if (node.tagName.toLowerCase() === 'a') {
          addNoReferrer(node);
        }
      }
    }
  };

  observationRoot.querySelectorAll('a').forEach(anchor => addNoReferrer(anchor));

  // Create an observer instance linked to the callback function
  const observer = new MutationObserver(callback);

  // Start observing the target node for configured mutations
  observer.observe(observationRoot, config);
}
