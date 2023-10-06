(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[6559],{70502:function(e,t){"use strict";Object.defineProperty(t,"__esModule",{value:!0}),t.getRenderPropValue=void 0,t.getRenderPropValue=function(e){return e?"function"==typeof e?e():e:null}},27847:function(e,t,n){"use strict";var r=n(75263).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=t.SizeContextProvider=void 0;var a=r(n(67294)),o=a.createContext("default");t.SizeContextProvider=function(e){var t=e.children,n=e.size;return a.createElement(o.Consumer,null,function(e){return a.createElement(o.Provider,{value:n||e},t)})},t.default=o},71511:function(e,t,n){"use strict";var r=n(75263).default,a=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var o=a(n(10434)),l=a(n(38416)),i=a(n(18698)),s=a(n(27424)),u=a(n(94184)),c=a(n(48555)),f=n(75531),d=r(n(67294)),p=n(31929),m=a(n(60872)),v=n(67046);a(n(13594));var g=a(n(27847)),__rest=function(e,t){var n={};for(var r in e)Object.prototype.hasOwnProperty.call(e,r)&&0>t.indexOf(r)&&(n[r]=e[r]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var a=0,r=Object.getOwnPropertySymbols(e);a<r.length;a++)0>t.indexOf(r[a])&&Object.prototype.propertyIsEnumerable.call(e,r[a])&&(n[r[a]]=e[r[a]]);return n},y=d.forwardRef(function(e,t){var n,r,a,y=d.useContext(g.default),h=d.useState(1),b=(0,s.default)(h,2),j=b[0],w=b[1],x=d.useState(!1),E=(0,s.default)(x,2),O=E[0],P=E[1],C=d.useState(!0),S=(0,s.default)(C,2),_=S[0],N=S[1],k=d.useRef(null),R=d.useRef(null),T=(0,f.composeRef)(t,k),M=d.useContext(p.ConfigContext).getPrefixCls,setScaleParam=function(){if(R.current&&k.current){var t=R.current.offsetWidth,n=k.current.offsetWidth;if(0!==t&&0!==n){var r=e.gap,a=void 0===r?4:r;2*a<n&&w(n-2*a<t?(n-2*a)/t:1)}}};d.useEffect(function(){P(!0)},[]),d.useEffect(function(){N(!0),w(1)},[e.src]),d.useEffect(function(){setScaleParam()},[e.gap]);var Z=e.prefixCls,D=e.shape,z=void 0===D?"circle":D,A=e.size,B=void 0===A?"default":A,L=e.src,V=e.srcSet,F=e.icon,I=e.className,U=e.alt,W=e.draggable,q=e.children,H=e.crossOrigin,G=__rest(e,["prefixCls","shape","size","src","srcSet","icon","className","alt","draggable","children","crossOrigin"]),Q="default"===B?y:B,X=Object.keys("object"===(0,i.default)(Q)&&Q||{}).some(function(e){return["xs","sm","md","lg","xl","xxl"].includes(e)}),J=(0,m.default)(X),Y=d.useMemo(function(){if("object"!==(0,i.default)(Q))return{};var e=Q[v.responsiveArray.find(function(e){return J[e]})];return e?{width:e,height:e,lineHeight:"".concat(e,"px"),fontSize:F?e/2:18}:{}},[J,Q]),K=M("avatar",Z),$=(0,u.default)((n={},(0,l.default)(n,"".concat(K,"-lg"),"large"===Q),(0,l.default)(n,"".concat(K,"-sm"),"small"===Q),n)),ee=d.isValidElement(L),et=(0,u.default)(K,$,(r={},(0,l.default)(r,"".concat(K,"-").concat(z),!!z),(0,l.default)(r,"".concat(K,"-image"),ee||L&&_),(0,l.default)(r,"".concat(K,"-icon"),!!F),r),I),en="number"==typeof Q?{width:Q,height:Q,lineHeight:"".concat(Q,"px"),fontSize:F?Q/2:18}:{};if("string"==typeof L&&_)a=d.createElement("img",{src:L,draggable:W,srcSet:V,onError:function(){var t=e.onError;!1!==(t?t():void 0)&&N(!1)},alt:U,crossOrigin:H});else if(ee)a=L;else if(F)a=F;else if(O||1!==j){var er="scale(".concat(j,") translateX(-50%)"),ea="number"==typeof Q?{lineHeight:"".concat(Q,"px")}:{};a=d.createElement(c.default,{onResize:setScaleParam},d.createElement("span",{className:"".concat(K,"-string"),ref:R,style:(0,o.default)((0,o.default)({},ea),{msTransform:er,WebkitTransform:er,transform:er})},q))}else a=d.createElement("span",{className:"".concat(K,"-string"),style:{opacity:0},ref:R},q);return delete G.onError,delete G.gap,d.createElement("span",(0,o.default)({},G,{style:(0,o.default)((0,o.default)((0,o.default)({},en),Y),G.style),className:et,ref:T}),a)});t.default=y},51289:function(e,t,n){"use strict";var r=n(75263).default,a=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var o=a(n(38416)),l=a(n(94184)),i=a(n(45598)),s=r(n(67294)),u=n(31929),c=a(n(62443)),f=n(47419),d=a(n(71511)),p=n(27847);t.default=function(e){var t=s.useContext(u.ConfigContext),n=t.getPrefixCls,r=t.direction,a=e.prefixCls,m=e.className,v=e.maxCount,g=e.maxStyle,y=e.size,h=n("avatar-group",a),b=(0,l.default)(h,(0,o.default)({},"".concat(h,"-rtl"),"rtl"===r),void 0===m?"":m),j=e.children,w=e.maxPopoverPlacement,x=e.maxPopoverTrigger,E=(0,i.default)(j).map(function(e,t){return(0,f.cloneElement)(e,{key:"avatar-key-".concat(t)})}),O=E.length;if(v&&v<O){var P=E.slice(0,v),C=E.slice(v,O);return P.push(s.createElement(c.default,{key:"avatar-popover-key",content:C,trigger:void 0===x?"hover":x,placement:void 0===w?"top":w,overlayClassName:"".concat(h,"-popover")},s.createElement(d.default,{style:g},"+".concat(O-v)))),s.createElement(p.SizeContextProvider,{size:y},s.createElement("div",{className:b,style:e.style},P))}return s.createElement(p.SizeContextProvider,{size:y},s.createElement("div",{className:b,style:e.style},E))}},84960:function(e,t,n){"use strict";var r=n(64836).default;t.ZP=void 0;var a=r(n(71511)),o=r(n(51289)),l=a.default;l.Group=o.default,t.ZP=l},62443:function(e,t,n){"use strict";var r=n(75263).default,a=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var o=a(n(10434)),l=r(n(67294)),i=n(70502),s=n(53683),u=n(31929),c=a(n(94055)),__rest=function(e,t){var n={};for(var r in e)Object.prototype.hasOwnProperty.call(e,r)&&0>t.indexOf(r)&&(n[r]=e[r]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var a=0,r=Object.getOwnPropertySymbols(e);a<r.length;a++)0>t.indexOf(r[a])&&Object.prototype.propertyIsEnumerable.call(e,r[a])&&(n[r[a]]=e[r[a]]);return n},Overlay=function(e){var t=e.title,n=e.content,r=e.prefixCls;return l.createElement(l.Fragment,null,t&&l.createElement("div",{className:"".concat(r,"-title")},(0,i.getRenderPropValue)(t)),l.createElement("div",{className:"".concat(r,"-inner-content")},(0,i.getRenderPropValue)(n)))},f=l.forwardRef(function(e,t){var n=e.prefixCls,r=e.title,a=e.content,i=e._overlay,f=e.placement,d=e.trigger,p=e.mouseEnterDelay,m=e.mouseLeaveDelay,v=e.overlayStyle,g=__rest(e,["prefixCls","title","content","_overlay","placement","trigger","mouseEnterDelay","mouseLeaveDelay","overlayStyle"]),y=l.useContext(u.ConfigContext).getPrefixCls,h=y("popover",n),b=y(),j=l.useMemo(function(){return i||(r||a?l.createElement(Overlay,{prefixCls:h,title:r,content:a}):null)},[i,r,a,h]);return l.createElement(c.default,(0,o.default)({placement:void 0===f?"top":f,trigger:void 0===d?"hover":d,mouseEnterDelay:void 0===p?.1:p,mouseLeaveDelay:void 0===m?.1:m,overlayStyle:void 0===v?{}:v},g,{prefixCls:h,ref:t,overlay:j,transitionName:(0,s.getTransitionName)(b,"zoom-big",g.transitionName)}))});t.default=f},89277:function(e,t,n){"use strict";var r=n(75263).default,a=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var o=a(n(10434));a(n(18698));var l=r(n(67294));a(n(13594));var i=a(n(28460)),__rest=function(e,t){var n={};for(var r in e)Object.prototype.hasOwnProperty.call(e,r)&&0>t.indexOf(r)&&(n[r]=e[r]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var a=0,r=Object.getOwnPropertySymbols(e);a<r.length;a++)0>t.indexOf(r[a])&&Object.prototype.propertyIsEnumerable.call(e,r[a])&&(n[r[a]]=e[r[a]]);return n},s=l.forwardRef(function(e,t){var n=e.ellipsis,r=e.rel,a=__rest(e,["ellipsis","rel"]),s=(0,o.default)((0,o.default)({},a),{rel:void 0===r&&"_blank"===a.target?"noopener noreferrer":r});return delete s.navigate,l.createElement(i.default,(0,o.default)({},s,{ref:t,ellipsis:!!n,component:"a"}))});t.default=s},21987:function(e,t,n){"use strict";var r=n(75263).default,a=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var o=a(n(10434)),l=r(n(67294)),i=a(n(28460)),s=l.forwardRef(function(e,t){return l.createElement(i.default,(0,o.default)({ref:t},e,{component:"div"}))});t.default=s},15394:function(e,t,n){"use strict";var r=n(75263).default,a=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var o=a(n(10434)),l=a(n(18698)),i=a(n(18475)),s=r(n(67294));a(n(13594));var u=a(n(28460)),__rest=function(e,t){var n={};for(var r in e)Object.prototype.hasOwnProperty.call(e,r)&&0>t.indexOf(r)&&(n[r]=e[r]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var a=0,r=Object.getOwnPropertySymbols(e);a<r.length;a++)0>t.indexOf(r[a])&&Object.prototype.propertyIsEnumerable.call(e,r[a])&&(n[r[a]]=e[r[a]]);return n},c=s.forwardRef(function(e,t){var n=e.ellipsis,r=__rest(e,["ellipsis"]),a=s.useMemo(function(){return n&&"object"===(0,l.default)(n)?(0,i.default)(n,["expandable","rows"]):n},[n]);return s.createElement(u.default,(0,o.default)({ref:t},r,{ellipsis:a,component:"span"}))});t.default=c},53740:function(e,t,n){"use strict";var r=n(64836).default;t.default=void 0;var a=r(n(89277)),o=r(n(21987)),l=r(n(15394)),i=r(n(34528)),s=r(n(89652)).default;s.Text=l.default,s.Link=a.default,s.Title=i.default,s.Paragraph=o.default,t.default=s},13882:function(e,t,n){"use strict";function requiredArgs(e,t){if(t.length<e)throw TypeError(e+" argument"+(e>1?"s":"")+" required, but only "+t.length+" present")}n.d(t,{Z:function(){return requiredArgs}})},19013:function(e,t,n){"use strict";n.d(t,{Z:function(){return toDate}});var r=n(71002),a=n(13882);function toDate(e){(0,a.Z)(1,arguments);var t=Object.prototype.toString.call(e);return e instanceof Date||"object"===(0,r.Z)(e)&&"[object Date]"===t?new Date(e.getTime()):"number"==typeof e||"[object Number]"===t?new Date(e):(("string"==typeof e||"[object String]"===t)&&"undefined"!=typeof console&&(console.warn("Starting with v2.0.0-beta.1 date-fns doesn't accept strings as date arguments. Please use `parseISO` to parse strings. See: https://github.com/date-fns/date-fns/blob/master/docs/upgradeGuide.md#string-arguments"),console.warn(Error().stack)),new Date(NaN))}},10887:function(e,t,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/admin/chat/emojis",function(){return n(52345)}])},52345:function(e,t,n){"use strict";n.r(t);var r=n(85893),a=n(84960),o=n(65400),l=n(70302),i=n(5789),s=n(55673),u=n(94055),c=n(53740),f=n(28465),d=n(67294),p=n(5152),m=n.n(p),v=n(50988),g=n(3111),y=n(60942),h=n(24770),b=n(42124),j=n(90695);let{Meta:w}=l.default,x=m()(()=>Promise.resolve().then(n.t.bind(n,40753,23)),{loadableGenerated:{webpack:()=>[40753]},ssr:!1}),{Title:E,Paragraph:O}=c.default,Emoji=()=>{let[e,t]=(0,d.useState)([]),[n,c]=(0,d.useState)(!1),[p,m]=(0,d.useState)(null),[j,P]=(0,d.useState)(null),C=null,resetStates=()=>{m(null),clearTimeout(C),C=null};async function getEmojis(){c(!0);try{let e=await (0,g.rQ)("/api/emoji");t(e)}catch(e){console.error("error fetching emojis",e)}c(!1)}async function handleDelete(e){let t="/".concat(e.split("/").slice(3).join("/"));console.log(t),c(!0),m((0,h.kg)(h.Jk,"Deleting emoji..."));try{let e=await (0,g.rQ)(g.Ff,{method:"POST",data:{name:t}});if(e instanceof Error)throw e;m((0,h.kg)(h.zv,"Emoji deleted")),C=setTimeout(resetStates,b.sI)}catch(e){m((0,h.kg)(h.Un,"".concat(e))),c(!1),C=setTimeout(resetStates,b.sI)}getEmojis()}async function handleEmojiUpload(){c(!0);try{m((0,h.kg)(h.Jk,"Converting emoji..."));let e=await new Promise((e,t)=>{if(!y.dr.includes(j.type)){let e="File type is not supported: ".concat(j.type);return t(e)}(0,y.y3)(j,t=>e({name:j.name,url:t}))});m((0,h.kg)(h.Jk,"Uploading emoji..."));let t=await (0,g.rQ)(g.Qc,{method:"POST",data:{name:e.name,data:e.url}});if(t instanceof Error)throw t;m((0,h.kg)(h.zv,"Emoji uploaded successfully!")),getEmojis()}catch(e){m((0,h.kg)(h.Un,"".concat(e)))}C=setTimeout(resetStates,b.sI),c(!1)}return(0,d.useEffect)(()=>{getEmojis()},[]),(0,r.jsxs)("div",{children:[(0,r.jsx)(E,{children:"Emojis"}),(0,r.jsx)(O,{children:"Here you can upload new custom emojis for usage in the chat. When uploading a new emoji, the filename without extension will be used as emoji name. Additionally, emoji names are case-insensitive. For best results, ensure all emoji have unique names."}),(0,r.jsx)("br",{}),(0,r.jsx)(f.Z,{name:"emoji",listType:"picture",className:"emoji-uploader",showUploadList:!1,accept:y.dr.join(","),beforeUpload:P,customRequest:handleEmojiUpload,disabled:n,children:(0,r.jsx)(o.default,{type:"primary",disabled:n,children:"Upload new emoji"})}),(0,r.jsx)(v.Z,{status:p}),(0,r.jsx)("br",{}),(0,r.jsx)(s.Z,{children:e.map(e=>(0,r.jsx)(i.Z,{style:{padding:"10px"},children:(0,r.jsx)(l.default,{style:{width:120,marginTop:16},actions:[],children:(0,r.jsx)(w,{description:[(0,r.jsxs)("div",{style:{display:"flex",justifyItems:"center",alignItems:"center",flexDirection:"column",gap:"20px"},children:[(0,r.jsx)(u.default,{title:e.name,children:(0,r.jsx)(a.ZP,{style:{height:50,width:50},src:e.url})}),(0,r.jsx)(o.default,{size:"small",type:"ghost",title:"Delete emoji",style:{position:"absolute",right:0,top:0,height:24,width:24,border:"none",color:"gray"},onClick:()=>handleDelete(e.url),icon:(0,r.jsx)(x,{})})]})]})})},e.name))}),(0,r.jsx)("br",{})]})};Emoji.getLayout=function(e){return(0,r.jsx)(j.l,{page:e})},t.default=Emoji},60942:function(e,t,n){"use strict";n.d(t,{Z7:function(){return r},dr:function(){return a},kR:function(){return readableBytes},y3:function(){return getBase64}});let r=2097152,a=["image/png","image/jpeg","image/gif"];function getBase64(e,t){let n=new FileReader;n.addEventListener("load",()=>t(n.result)),n.readAsDataURL(e)}function readableBytes(e){let t=Math.floor(Math.log(e)/Math.log(1024)),n=1*Number((e/Math.pow(1024,t)).toFixed(2));return"".concat(n," ").concat(["B","KB","MB","GB","TB","PB","EB","ZB","YB"][t])}},60057:function(e,t,n){"use strict";n.r(t),n.d(t,{default:function(){return j}});var r=n(4942),a=n(1413),o=n(97685),l=n(45987),i=n(67294),s=n(81263),u=n(94184),c=n.n(u),f={adjustX:1,adjustY:1},d=[0,0],p={topLeft:{points:["bl","tl"],overflow:f,offset:[0,-4],targetOffset:d},topCenter:{points:["bc","tc"],overflow:f,offset:[0,-4],targetOffset:d},topRight:{points:["br","tr"],overflow:f,offset:[0,-4],targetOffset:d},bottomLeft:{points:["tl","bl"],overflow:f,offset:[0,4],targetOffset:d},bottomCenter:{points:["tc","bc"],overflow:f,offset:[0,4],targetOffset:d},bottomRight:{points:["tr","br"],overflow:f,offset:[0,4],targetOffset:d}},m=n(15105),v=n(75164),g=n(88603),y=m.Z.ESC,h=m.Z.TAB,b=["arrow","prefixCls","transitionName","animation","align","placement","placements","getPopupContainer","showAction","hideAction","overlayClassName","overlayStyle","visible","trigger","autoFocus"],j=i.forwardRef(function(e,t){var n,u,f,d,m,j,w,x,E,O,P,C,S,_,N,k,R=e.arrow,T=void 0!==R&&R,M=e.prefixCls,Z=void 0===M?"rc-dropdown":M,D=e.transitionName,z=e.animation,A=e.align,B=e.placement,L=e.placements,V=e.getPopupContainer,F=e.showAction,I=e.hideAction,U=e.overlayClassName,W=e.overlayStyle,q=e.visible,H=e.trigger,G=void 0===H?["hover"]:H,Q=e.autoFocus,X=(0,l.Z)(e,b),J=i.useState(),Y=(0,o.Z)(J,2),K=Y[0],$=Y[1],ee="visible"in e?q:K,et=i.useRef(null);i.useImperativeHandle(t,function(){return et.current}),f=(u={visible:ee,setTriggerVisible:$,triggerRef:et,onVisibleChange:e.onVisibleChange,autoFocus:Q}).visible,d=u.setTriggerVisible,m=u.triggerRef,j=u.onVisibleChange,w=u.autoFocus,x=i.useRef(!1),E=function(){if(f&&m.current){var e,t,n,r;null===(e=m.current)||void 0===e||null===(t=e.triggerRef)||void 0===t||null===(n=t.current)||void 0===n||null===(r=n.focus)||void 0===r||r.call(n),d(!1),"function"==typeof j&&j(!1)}},O=function(){var e,t,n,r,a=(0,g.tS)(null===(e=m.current)||void 0===e?void 0:null===(t=e.popupRef)||void 0===t?void 0:null===(n=t.current)||void 0===n?void 0:null===(r=n.getElement)||void 0===r?void 0:r.call(n))[0];return null!=a&&!!a.focus&&(a.focus(),x.current=!0,!0)},P=function(e){switch(e.keyCode){case y:E();break;case h:var t=!1;x.current||(t=O()),t?e.preventDefault():E()}},i.useEffect(function(){return f?(window.addEventListener("keydown",P),w&&(0,v.Z)(O,3),function(){window.removeEventListener("keydown",P),x.current=!1}):function(){x.current=!1}},[f]);var getOverlayElement=function(){var t=e.overlay;return"function"==typeof t?t():t},getMenuElement=function(){var e=getOverlayElement();return i.createElement(i.Fragment,null,T&&i.createElement("div",{className:"".concat(Z,"-arrow")}),e)},en=I;return en||-1===G.indexOf("contextMenu")||(en=["click"]),i.createElement(s.Z,(0,a.Z)((0,a.Z)({builtinPlacements:void 0===L?p:L},X),{},{prefixCls:Z,ref:et,popupClassName:c()(U,(0,r.Z)({},"".concat(Z,"-show-arrow"),T)),popupStyle:W,action:G,showAction:F,hideAction:en||[],popupPlacement:void 0===B?"bottomLeft":B,popupAlign:A,popupTransitionName:D,popupAnimation:z,popupVisible:ee,stretch:(C=e.minOverlayWidthMatchTrigger,S=e.alignPoint,"minOverlayWidthMatchTrigger"in e?C:!S)?"minWidth":"",popup:"function"==typeof e.overlay?getMenuElement:getMenuElement(),onPopupVisibleChange:function(t){var n=e.onVisibleChange;$(t),"function"==typeof n&&n(t)},onPopupClick:function(t){var n=e.onOverlayClick;$(!1),n&&n(t)},getPopupContainer:V}),(N=(_=e.children).props?_.props:{},k=c()(N.className,void 0!==(n=e.openClassName)?n:"".concat(Z,"-open")),ee&&_?i.cloneElement(_,{className:k}):_))})}},function(e){e.O(0,[5596,1130,4104,9403,1024,3942,971,6697,1664,1749,1700,2122,3068,6300,8531,7423,8465,695,9774,2888,179],function(){return e(e.s=10887)}),_N_E=e.O()}]);