(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[9266],{79266:function(e,n,t){"use strict";t.r(n),t.d(n,{UserDropdown:function(){return A}});var r=t(85893),o=t(66516),a=t(13013),i=t(71577),l=t(4480),s=t(67294),c=t(46977),u=t(5152),d=t.n(u),f=t(37068),p=t(44974),h=t(27345),y=t.n(h),v=t(69183);let m=d()(()=>Promise.all([t.e(2074),t.e(8244)]).then(t.t.bind(t,18244,23)),{loadableGenerated:{webpack:()=>[18244]},ssr:!1}),b=d()(()=>Promise.all([t.e(2074),t.e(775)]).then(t.t.bind(t,10775,23)),{loadableGenerated:{webpack:()=>[10775]},ssr:!1}),k=d()(()=>Promise.all([t.e(2074),t.e(6275)]).then(t.t.bind(t,6275,23)),{loadableGenerated:{webpack:()=>[6275]},ssr:!1}),w=d()(()=>Promise.all([t.e(2074),t.e(5672)]).then(t.t.bind(t,65672,23)),{loadableGenerated:{webpack:()=>[65672]},ssr:!1}),C=d()(()=>Promise.all([t.e(2074),t.e(5412)]).then(t.t.bind(t,95412,23)),{loadableGenerated:{webpack:()=>[95412]},ssr:!1}),g=d()(()=>Promise.all([t.e(5874),t.e(5402),t.e(9614)]).then(t.bind(t,29614)).then(e=>e.Modal),{loadableGenerated:{webpack:()=>[29614]},ssr:!1}),j=d()(()=>Promise.all([t.e(173),t.e(5874),t.e(5257),t.e(4041),t.e(3427)]).then(t.bind(t,7783)).then(e=>e.NameChangeModal),{loadableGenerated:{webpack:()=>[7783]},ssr:!1}),x=d()(()=>Promise.all([t.e(173),t.e(5874),t.e(5818),t.e(4526),t.e(9644)]).then(t.bind(t,39644)).then(e=>e.AuthModal),{loadableGenerated:{webpack:()=>[39644]},ssr:!1}),A=e=>{let{username:n}=e,[t,u]=(0,s.useState)(!1),[d,h]=(0,s.useState)(!1),[A,E]=(0,l.FV)(p.ZA),L=(0,l.sJ)(p.Q),_=()=>{E(!A)},S=()=>{u(!0)};(0,c.y1)("c",_,{enableOnContentEditable:!1},[A]);let P=(0,l.sJ)(p.db);if(!P)return null;let{displayName:O}=P,D=(0,r.jsxs)(o.Z,{children:[(0,r.jsx)(o.Z.Item,{icon:(0,r.jsx)(b,{}),onClick:()=>S(),children:"Change name"},"0"),(0,r.jsx)(o.Z.Item,{icon:(0,r.jsx)(k,{}),onClick:()=>h(!0),children:"Authenticate"},"1"),L.chatAvailable&&(0,r.jsx)(o.Z.Item,{icon:(0,r.jsx)(w,{}),onClick:()=>_(),"aria-expanded":A,children:A?"Hide Chat":"Show Chat"},"3")]});return(0,r.jsx)(f.SV,{fallbackRender:e=>{let{error:n,resetErrorBoundary:t}=e;return(0,r.jsx)(v.A,{componentName:"UserDropdown",message:n.message,retryFunction:t})},children:(0,r.jsxs)("div",{id:"user-menu",className:"".concat(y().root),children:[(0,r.jsx)(a.Z,{overlay:D,trigger:["click"],children:(0,r.jsxs)(i.Z,{type:"primary",icon:(0,r.jsx)(C,{className:y().userIcon}),children:[(0,r.jsx)("span",{className:y().username,children:n||O}),(0,r.jsx)(m,{})]})}),(0,r.jsx)(g,{title:"Change Chat Display Name",open:t,handleCancel:()=>u(!1),children:(0,r.jsx)(j,{})}),(0,r.jsx)(g,{title:"Authenticate",open:d,handleCancel:()=>h(!1),children:(0,r.jsx)(x,{})})]})})}},27345:function(e){e.exports={root:"UserDropdown_root__IdxfQ","ant-space":"UserDropdown_ant-space__XJTZ3","ant-space-item":"UserDropdown_ant-space-item__w4nC2",userIcon:"UserDropdown_userIcon__A5XgE",username:"UserDropdown_username__nVyPA"}},46977:function(e,n,t){"use strict";t.d(n,{y1:function(){return k}});var r=t(67294);function o(){return(o=Object.assign?Object.assign.bind():function(e){for(var n=1;n<arguments.length;n++){var t=arguments[n];for(var r in t)Object.prototype.hasOwnProperty.call(t,r)&&(e[r]=t[r])}return e}).apply(this,arguments)}t(85893);var a=["shift","alt","meta","mod","ctrl"],i={esc:"escape",return:"enter",".":"period",",":"comma","-":"slash"," ":"space","`":"backquote","#":"backslash","+":"bracketright",ShiftLeft:"shift",ShiftRight:"shift",AltLeft:"alt",AltRight:"alt",MetaLeft:"meta",MetaRight:"meta",OSLeft:"meta",OSRight:"meta",ControlLeft:"ctrl",ControlRight:"ctrl"};function l(e){return(i[e]||e).trim().toLowerCase().replace(/key|digit|numpad|arrow/,"")}function s(e,n){return void 0===n&&(n=","),e.split(n)}function c(e,n){void 0===n&&(n="+");var t=e.toLocaleLowerCase().split(n).map(function(e){return l(e)});return o({},{alt:t.includes("alt"),ctrl:t.includes("ctrl")||t.includes("control"),shift:t.includes("shift"),meta:t.includes("meta"),mod:t.includes("mod")},{keys:t.filter(function(e){return!a.includes(e)})})}"undefined"!=typeof document&&(document.addEventListener("keydown",function(e){void 0!==e.key&&d([l(e.key),l(e.code)])}),document.addEventListener("keyup",function(e){void 0!==e.key&&f([l(e.key),l(e.code)])})),"undefined"!=typeof window&&window.addEventListener("blur",function(){u.clear()});var u=new Set;function d(e){var n=Array.isArray(e)?e:[e];u.has("meta")&&u.forEach(function(e){return!a.includes(e)&&u.delete(e.toLowerCase())}),n.forEach(function(e){return u.add(e.toLowerCase())})}function f(e){var n=Array.isArray(e)?e:[e];"meta"===e?u.clear():n.forEach(function(e){return u.delete(e.toLowerCase())})}function p(e,n){var t=e.target;void 0===n&&(n=!1);var r=t&&t.tagName;return n instanceof Array?!!(r&&n&&n.some(function(e){return e.toLowerCase()===r.toLowerCase()})):!!(r&&n&&!0===n)}var h=function(e,n,t){void 0===t&&(t=!1);var r,o=n.alt,a=n.meta,i=n.mod,s=n.shift,c=n.ctrl,d=n.keys,f=e.key,p=e.code,h=e.ctrlKey,y=e.metaKey,v=e.shiftKey,m=e.altKey,b=l(p),k=f.toLowerCase();if(!t){if(!m===o&&"alt"!==k||!v===s&&"shift"!==k)return!1;if(i){if(!y&&!h)return!1}else if(!y===a&&"meta"!==k&&"os"!==k||!h===c&&"ctrl"!==k&&"control"!==k)return!1}return!!(d&&1===d.length&&(d.includes(k)||d.includes(b)))||(d?(void 0===r&&(r=","),(Array.isArray(d)?d:d.split(r)).every(function(e){return u.has(e.trim().toLowerCase())})):!d)},y=(0,r.createContext)(void 0),v=(0,r.createContext)({hotkeys:[],enabledScopes:[],toggleScope:function(){},enableScope:function(){},disableScope:function(){}}),m=function(e){e.stopPropagation(),e.preventDefault(),e.stopImmediatePropagation()},b="undefined"!=typeof window?r.useLayoutEffect:r.useEffect;function k(e,n,t,o){var a,i=(0,r.useRef)(null),u=(0,r.useRef)(!1),k=t instanceof Array?o instanceof Array?void 0:o:t,w=e instanceof Array?e.join(null==k?void 0:k.splitKey):e,C=t instanceof Array?t:o instanceof Array?o:void 0,g=(0,r.useCallback)(n,null!=C?C:[]),j=(0,r.useRef)(g);C?j.current=g:j.current=n;var x=(!function e(n,t){return n&&t&&"object"==typeof n&&"object"==typeof t?Object.keys(n).length===Object.keys(t).length&&Object.keys(n).reduce(function(r,o){return r&&e(n[o],t[o])},!0):n===t}((a=(0,r.useRef)(void 0)).current,k)&&(a.current=k),a.current),A=(0,r.useContext)(v).enabledScopes,E=(0,r.useContext)(y);return b(function(){if((null==x?void 0:x.enabled)!==!1&&(e=null==x?void 0:x.scopes,0===A.length&&e?(console.warn('A hotkey has the "scopes" option set, however no active scopes were found. If you want to use the global scopes feature, you need to wrap your app in a <HotkeysProvider>'),!0):!!(!e||A.some(function(n){return e.includes(n)})||A.includes("*")))){var e,n=function(e,n){var t;if(void 0===n&&(n=!1),!p(e,["input","textarea","select"])||p(e,null==x?void 0:x.enableOnFormTags)){if(null!==i.current&&document.activeElement!==i.current&&!i.current.contains(document.activeElement)){m(e);return}(null==(t=e.target)||!t.isContentEditable||null!=x&&x.enableOnContentEditable)&&s(w,null==x?void 0:x.splitKey).forEach(function(t){var r,o,a,i=c(t,null==x?void 0:x.combinationKey);if(h(e,i,null==x?void 0:x.ignoreModifiers)||null!=(a=i.keys)&&a.includes("*")){if(n&&u.current)return;if(("function"==typeof(r=null==x?void 0:x.preventDefault)&&r(e,i)||!0===r)&&e.preventDefault(),"function"==typeof(o=null==x?void 0:x.enabled)?!o(e,i):!0!==o&&void 0!==o){m(e);return}j.current(e,i),n||(u.current=!0)}})}},t=function(e){void 0!==e.key&&(d(l(e.code)),((null==x?void 0:x.keydown)===void 0&&(null==x?void 0:x.keyup)!==!0||null!=x&&x.keydown)&&n(e))},r=function(e){void 0!==e.key&&(f(l(e.code)),u.current=!1,null!=x&&x.keyup&&n(e,!0))},o=i.current||(null==k?void 0:k.document)||document;return o.addEventListener("keyup",r),o.addEventListener("keydown",t),E&&s(w,null==x?void 0:x.splitKey).forEach(function(e){return E.addHotkey(c(e,null==x?void 0:x.combinationKey))}),function(){o.removeEventListener("keyup",r),o.removeEventListener("keydown",t),E&&s(w,null==x?void 0:x.splitKey).forEach(function(e){return E.removeHotkey(c(e,null==x?void 0:x.combinationKey))})}}},[w,x,A]),i}}}]);