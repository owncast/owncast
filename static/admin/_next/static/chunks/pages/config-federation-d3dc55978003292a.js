(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[457],{94594:function(e,n,t){"use strict";t.d(n,{Z:function(){return C}});var a=t(87462),o=t(4942),r=t(67294),i=t(97685),s=t(91),l=t(94184),c=t.n(l),d=t(21770),u=t(15105),h=r.forwardRef((function(e,n){var t,a=e.prefixCls,l=void 0===a?"rc-switch":a,h=e.className,f=e.checked,p=e.defaultChecked,m=e.disabled,v=e.loadingIcon,g=e.checkedChildren,y=e.unCheckedChildren,w=e.onClick,b=e.onChange,C=e.onKeyDown,x=(0,s.Z)(e,["prefixCls","className","checked","defaultChecked","disabled","loadingIcon","checkedChildren","unCheckedChildren","onClick","onChange","onKeyDown"]),k=(0,d.Z)(!1,{value:f,defaultValue:p}),P=(0,i.Z)(k,2),E=P[0],N=P[1];function j(e,n){var t=E;return m||(N(t=e),null===b||void 0===b||b(t,n)),t}var S=c()(l,h,(t={},(0,o.Z)(t,"".concat(l,"-checked"),E),(0,o.Z)(t,"".concat(l,"-disabled"),m),t));return r.createElement("button",Object.assign({},x,{type:"button",role:"switch","aria-checked":E,disabled:m,className:S,ref:n,onKeyDown:function(e){e.which===u.Z.LEFT?j(!1,e):e.which===u.Z.RIGHT&&j(!0,e),null===C||void 0===C||C(e)},onClick:function(e){var n=j(!E,e);null===w||void 0===w||w(n,e)}}),v,r.createElement("span",{className:"".concat(l,"-inner")},E?g:y))}));h.displayName="Switch";var f=h,p=t(50888),m=t(97202),v=t(59844),g=t(97647),y=t(21687),w=function(e,n){var t={};for(var a in e)Object.prototype.hasOwnProperty.call(e,a)&&n.indexOf(a)<0&&(t[a]=e[a]);if(null!=e&&"function"===typeof Object.getOwnPropertySymbols){var o=0;for(a=Object.getOwnPropertySymbols(e);o<a.length;o++)n.indexOf(a[o])<0&&Object.prototype.propertyIsEnumerable.call(e,a[o])&&(t[a[o]]=e[a[o]])}return t},b=r.forwardRef((function(e,n){var t,i=e.prefixCls,s=e.size,l=e.loading,d=e.className,u=void 0===d?"":d,h=e.disabled,b=w(e,["prefixCls","size","loading","className","disabled"]);(0,y.Z)("checked"in b||!("value"in b),"Switch","`value` is not a valid prop, do you mean `checked`?");var C=r.useContext(v.E_),x=C.getPrefixCls,k=C.direction,P=r.useContext(g.Z),E=x("switch",i),N=r.createElement("div",{className:"".concat(E,"-handle")},l&&r.createElement(p.Z,{className:"".concat(E,"-loading-icon")})),j=c()((t={},(0,o.Z)(t,"".concat(E,"-small"),"small"===(s||P)),(0,o.Z)(t,"".concat(E,"-loading"),l),(0,o.Z)(t,"".concat(E,"-rtl"),"rtl"===k),t),u);return r.createElement(m.Z,{insertExtraNode:!0},r.createElement(f,(0,a.Z)({},b,{prefixCls:E,className:j,disabled:h||l,ref:n,loadingIcon:N})))}));b.__ANT_SWITCH=!0,b.displayName="Switch";var C=b},50291:function(e,n,t){(window.__NEXT_P=window.__NEXT_P||[]).push(["/config-federation",function(){return t(6898)}])},15976:function(e,n,t){"use strict";t.d(n,{Z:function(){return f}});var a=t(28520),o=t.n(a),r=t(85893),i=t(67294),s=t(94594),l=t(83200),c=t(78464),d=t(25964),u=t(35159);function h(e,n,t,a,o,r,i){try{var s=e[r](i),l=s.value}catch(c){return void t(c)}s.done?n(l):Promise.resolve(l).then(a,o)}function f(e){var n,t=(0,i.useState)(null),a=t[0],f=t[1],p=null,m=((0,i.useContext)(u.aC)||{}).setFieldInConfigState,v=e.apiPath,g=e.checked,y=e.reversed,w=void 0!==y&&y,b=e.configPath,C=void 0===b?"":b,x=e.disabled,k=void 0!==x&&x,P=e.fieldName,E=e.label,N=e.tip,j=e.useSubmit,S=e.onChange,Z=function(){f(null),clearTimeout(p),p=null},O=(n=o().mark((function e(n){var t;return o().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!j){e.next=6;break}return f((0,l.kg)(l.Jk)),t=w?!n:n,e.next=5,(0,d.Si)({apiPath:v,data:{value:t},onSuccess:function(){m({fieldName:P,value:t,path:C}),f((0,l.kg)(l.zv))},onError:function(e){f((0,l.kg)(l.Un,"There was an error: ".concat(e)))}});case 5:p=setTimeout(Z,d.sI);case 6:S&&S(n);case 7:case"end":return e.stop()}}),e)})),function(){var e=this,t=arguments;return new Promise((function(a,o){var r=n.apply(e,t);function i(e){h(r,a,o,i,s,"next",e)}function s(e){h(r,a,o,i,s,"throw",e)}i(void 0)}))}),T=null!==a&&a.type===l.Jk;return(0,r.jsxs)("div",{className:"formfield-container toggleswitch-container",children:[E&&(0,r.jsx)("div",{className:"label-side",children:(0,r.jsx)("span",{className:"formfield-label",children:E})}),(0,r.jsxs)("div",{className:"input-side",children:[(0,r.jsxs)("div",{className:"input-group",children:[(0,r.jsx)(s.Z,{className:"switch field-".concat(P),loading:T,onChange:O,defaultChecked:g,checked:g,checkedChildren:"ON",unCheckedChildren:"OFF",disabled:k}),(0,r.jsx)(c.Z,{status:a})]}),(0,r.jsx)("p",{className:"field-tip",children:N})]})]})}f.defaultProps={apiPath:"",checked:!1,reversed:!1,configPath:"",disabled:!1,label:"",tip:"",useSubmit:!1,onChange:null}},6898:function(e,n,t){"use strict";t.r(n),t.d(n,{default:function(){return y}});var a=t(85893),o=t(56516),r=t(71577),i=t(17256),s=t(67294),l=t(45697),c=t.n(l),d=t(48419),u=t(50197),h=t(15976),f=t(25964),p=t(35159);function m(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function v(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{},a=Object.keys(t);"function"===typeof Object.getOwnPropertySymbols&&(a=a.concat(Object.getOwnPropertySymbols(t).filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable})))),a.forEach((function(n){m(e,n,t[n])}))}return e}function g(e){var n=e.cancelPressed,t=e.okPressed;return(0,a.jsxs)(o.Z,{width:800,title:"Enable Federation",visible:!0,onCancel:n,footer:(0,a.jsx)(r.Z,{onClick:t,children:"I Understand."}),children:[(0,a.jsx)(i.Z.Title,{level:3,children:"What is Federation?"}),(0,a.jsxs)(i.Z.Paragraph,{children:["Enabling Federation allows your Owncast server to be a part of the greater"," ",(0,a.jsx)("a",{href:"https://en.wikipedia.org/wiki/Fediverse",rel:"noopener noreferrer",target:"_blank",children:"Fediverse"}),", allowing people to be notified when you go live and for you to send posts from Owncast to your followers. It can also surface actions taken on the Fediverse such as following, sharing or liking to show who is engaging with your stream."]}),(0,a.jsx)(i.Z.Title,{level:3,children:"Understand the following"}),(0,a.jsx)(i.Z.Paragraph,{children:"Due to other servers interacting with yours there are some things to keep in mind."}),(0,a.jsxs)("ul",{children:[(0,a.jsx)("li",{children:"You must always host your Owncast server with SSL using a https url."}),(0,a.jsx)("li",{children:"If you ever change your server name or Fediverse username you will be seen as a completely different user on the Fediverse, and the old user will disappear. It is best to never change these once you set them."}),(0,a.jsx)("li",{children:"Turning on Private Federation mode will require you to manually approve each follower."})]})]})}function y(){var e=i.Z.Title,n=(0,s.useState)(null),t=n[0],o=n[1],r=(0,s.useState)(!1),l=r[0],c=r[1],y=((0,s.useContext)(p.aC)||{}).serverConfig,w=y.federation,b=y.yp,C=w.enabled,x=w.isPrivate,k=w.username,P=w.goLiveMessage,E=w.showEngagement,N=b.instanceUrl,j=function(e){var n=e.fieldName,a=e.value;o(v({},t,m({},n,a)))};if((0,s.useEffect)((function(){o({enabled:C,isPrivate:x,username:k,goLiveMessage:P,showEngagement:E,instanceUrl:b.instanceUrl})}),[y,b]),!t)return null;var S=""!==N,Z=N.startsWith("https://");return(0,a.jsxs)("div",{className:"config-server-details-form",children:[(0,a.jsx)(e,{children:"Federation Settings"}),"Explain what the Fediverse is here and talk about what happens if you were to enable this feature.",(0,a.jsxs)("div",{className:"form-module config-server-details-container",children:[(0,a.jsx)(h.Z,v({fieldName:"enabled",onChange:function(e){e?c(!0):o({enabled:!1})}},f.Kl,{checked:t.enabled,disabled:!S||!Z})),(0,a.jsx)(u.ZP,v({fieldName:"instanceUrl"},f.yi,{value:t.instanceUrl,initialValue:b.instanceUrl,type:d.xA,onChange:j,onSubmit:function(){var e=""!==t.instanceUrl,n=t.instanceUrl.startsWith("https://");e&&n||((0,f.Si)({apiPath:f.Kl.apiPath,data:{value:!1}}),o({enabled:!1}))}})),(0,a.jsx)(h.Z,v({fieldName:"isPrivate"},f.LC,{checked:t.isPrivate})),(0,a.jsx)(u.ZP,v({required:!0,fieldName:"username"},f.Xc,{value:t.username,initialValue:k,onChange:j})),(0,a.jsx)(u.ZP,v({fieldName:"goLiveMessage"},f.BF,{type:d.Sk,value:t.goLiveMessage,initialValue:P,onChange:j})),(0,a.jsx)(h.Z,v({fieldName:"showEngagement"},f.FE,{checked:t.showEngagement})),l&&(0,a.jsx)(g,{cancelPressed:function(){c(!1),o({enabled:!1})},okPressed:function(){c(!1),o({enabled:!0})}}),(0,a.jsx)("br",{}),(0,a.jsx)("br",{})]})]})}g.propTypes={cancelPressed:c().func.isRequired,okPressed:c().func.isRequired}},74204:function(e,n,t){"use strict";var a;function o(e){if("undefined"===typeof document)return 0;if(e||void 0===a){var n=document.createElement("div");n.style.width="100%",n.style.height="200px";var t=document.createElement("div"),o=t.style;o.position="absolute",o.top="0",o.left="0",o.pointerEvents="none",o.visibility="hidden",o.width="200px",o.height="150px",o.overflow="hidden",t.appendChild(n),document.body.appendChild(t);var r=n.offsetWidth;t.style.overflow="scroll";var i=n.offsetWidth;r===i&&(i=t.clientWidth),document.body.removeChild(t),a=r-i}return a}function r(e){var n=e.match(/^(.*)px$/),t=Number(null===n||void 0===n?void 0:n[1]);return Number.isNaN(t)?o():t}function i(e){if("undefined"===typeof document||!e||!(e instanceof Element))return{width:0,height:0};var n=getComputedStyle(e,"::-webkit-scrollbar"),t=n.width,a=n.height;return{width:r(t),height:r(a)}}t.d(n,{Z:function(){return o},o:function(){return i}})},64217:function(e,n,t){"use strict";t.d(n,{Z:function(){return l}});var a=t(1413),o="".concat("accept acceptCharset accessKey action allowFullScreen allowTransparency\n    alt async autoComplete autoFocus autoPlay capture cellPadding cellSpacing challenge\n    charSet checked classID className colSpan cols content contentEditable contextMenu\n    controls coords crossOrigin data dateTime default defer dir disabled download draggable\n    encType form formAction formEncType formMethod formNoValidate formTarget frameBorder\n    headers height hidden high href hrefLang htmlFor httpEquiv icon id inputMode integrity\n    is keyParams keyType kind label lang list loop low manifest marginHeight marginWidth max maxLength media\n    mediaGroup method min minLength multiple muted name noValidate nonce open\n    optimum pattern placeholder poster preload radioGroup readOnly rel required\n    reversed role rowSpan rows sandbox scope scoped scrolling seamless selected\n    shape size sizes span spellCheck src srcDoc srcLang srcSet start step style\n    summary tabIndex target title type useMap value width wmode wrap"," ").concat("onCopy onCut onPaste onCompositionEnd onCompositionStart onCompositionUpdate onKeyDown\n    onKeyPress onKeyUp onFocus onBlur onChange onInput onSubmit onClick onContextMenu onDoubleClick\n    onDrag onDragEnd onDragEnter onDragExit onDragLeave onDragOver onDragStart onDrop onMouseDown\n    onMouseEnter onMouseLeave onMouseMove onMouseOut onMouseOver onMouseUp onSelect onTouchCancel\n    onTouchEnd onTouchMove onTouchStart onScroll onWheel onAbort onCanPlay onCanPlayThrough\n    onDurationChange onEmptied onEncrypted onEnded onError onLoadedData onLoadedMetadata\n    onLoadStart onPause onPlay onPlaying onProgress onRateChange onSeeked onSeeking onStalled onSuspend onTimeUpdate onVolumeChange onWaiting onLoad onError").split(/[\s\n]+/),r="aria-",i="data-";function s(e,n){return 0===e.indexOf(n)}function l(e){var n,t=arguments.length>1&&void 0!==arguments[1]&&arguments[1];n=!1===t?{aria:!0,data:!0,attr:!0}:!0===t?{aria:!0}:(0,a.Z)({},t);var l={};return Object.keys(e).forEach((function(t){(n.aria&&("role"===t||s(t,r))||n.data&&s(t,i)||n.attr&&o.includes(t))&&(l[t]=e[t])})),l}}},function(e){e.O(0,[516,774,888,179],(function(){return n=50291,e(e.s=n);var n}));var n=e.O();_N_E=n}]);