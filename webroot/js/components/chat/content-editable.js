/*
Since we can't really import react-contenteditable here, I'm borrowing code for this component from here:
github.com/lovasoa/react-contenteditable/

and here:
https://stackoverflow.com/questions/22677931/react-js-onchange-event-for-contenteditable/27255103#27255103

*/
import { h, Component, createRef } from '/js/web_modules/preact.js';

export function replaceCaret(el) {
  // Place the caret at the end of the element
  const target = document.createTextNode('');
  el.appendChild(target);
  // do not move caret if element was not focused
  const isTargetFocused = document.activeElement === el;
  if (target !== null && target.nodeValue !== null && isTargetFocused) {
    var sel = window.getSelection();
    if (sel !== null) {
      var range = document.createRange();
      range.setStart(target, target.nodeValue.length);
      range.collapse(true);
      sel.removeAllRanges();
      sel.addRange(range);
    }
    if (el) el.focus();
  }
}

function normalizeHtml(str) {
  return str && str.replace(/&nbsp;|\u202F|\u00A0/g, ' ');
}



export default class ContentEditable extends Component {
  constructor(props) {
    super(props);

    this.el = createRef();

    this.lastHtml = '';

    this.emitChange = this.emitChange.bind(this);
    this.getDOMElement = this.getDOMElement.bind(this);
  }

  shouldComponentUpdate(nextProps) {
    const { props } = this;
    const el = this.getDOMElement();

    // We need not rerender if the change of props simply reflects the user's edits.
    // Rerendering in this case would make the cursor/caret jump

    // Rerender if there is no element yet... (somehow?)
    if (!el) return true;

    // ...or if html really changed... (programmatically, not by user edit)
    if (
      normalizeHtml(nextProps.html) !== normalizeHtml(el.innerHTML)
    ) {
      return true;
    }

    // Handle additional properties
    return props.disabled !== nextProps.disabled ||
      props.tagName !== nextProps.tagName ||
      props.className !== nextProps.className ||
      props.innerRef !== nextProps.innerRef;
  }

  componentDidUpdate() {
    const el = this.getDOMElement();
    if (!el) return;

    // Perhaps React (whose VDOM gets outdated because we often prevent
    // rerendering) did not update the DOM. So we update it manually now.
    if (this.props.html !== el.innerHTML) {
      el.innerHTML = this.props.html;
    }
    this.lastHtml = this.props.html;
    replaceCaret(el);
  }

  getDOMElement() {
    return (this.props.innerRef && typeof this.props.innerRef !== 'function' ? this.props.innerRef : this.el).current;
  }


  emitChange(originalEvt) {
    const el = this.getDOMElement();
    if (!el) return;

    const html = el.innerHTML;
    if (this.props.onChange && html !== this.lastHtml) {
      // Clone event with Object.assign to avoid
      // "Cannot assign to read only property 'target' of object"
      const evt = Object.assign({}, originalEvt, {
        target: {
          value: html
        }
      });
      this.props.onChange(evt);
    }
    this.lastHtml = html;
  }

  render(props) {
    const { html, innerRef } = props;
    return h(
      'div',
      {
        ...props,
        ref: typeof innerRef === 'function' ? (current) => {
          innerRef(current)
          this.el.current = current
        } : innerRef || this.el,
        onInput: this.emitChange,
        onFocus: this.props.onFocus || this.emitChange,
        onBlur: this.props.onBlur || this.emitChange,
        onKeyup: this.props.onKeyUp || this.emitChange,
        onKeydown: this.props.onKeyDown || this.emitChange,
        contentEditable: !this.props.disabled,
        dangerouslySetInnerHTML: { __html: html },
      },
      this.props.children,
    );
  }
}
