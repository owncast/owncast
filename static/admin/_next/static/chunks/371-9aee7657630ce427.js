"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[371],{82911:function(e,n,t){t.d(n,{Z:function(){return i}});var r=t(1413),a=t(67294),s={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 708c-22.1 0-40-17.9-40-40s17.9-40 40-40 40 17.9 40 40-17.9 40-40 40zm62.9-219.5a48.3 48.3 0 00-30.9 44.8V620c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8v-21.5c0-23.1 6.7-45.9 19.9-64.9 12.9-18.6 30.9-32.8 52.1-40.9 34-13.1 56-41.6 56-72.7 0-44.1-43.1-80-96-80s-96 35.9-96 80v7.6c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V420c0-39.3 17.2-76 48.4-103.3C430.4 290.4 470 276 512 276s81.6 14.5 111.6 40.7C654.8 344 672 380.7 672 420c0 57.8-38.1 109.8-97.1 132.5z"}}]},name:"question-circle",theme:"filled"},o=t(42135),c=function(e,n){return a.createElement(o.Z,(0,r.Z)((0,r.Z)({},e),{},{ref:n,icon:s}))};c.displayName="QuestionCircleFilled";var i=a.forwardRef(c)},84674:function(e,n,t){t.d(n,{Z:function(){return i}});var r=t(1413),a=t(67294),s={icon:function(e,n){return{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm288.5 682.8L277.7 224C258 240 240 258 224 277.7l522.8 522.8C682.8 852.7 601 884 512 884c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372c0 89-31.3 170.8-83.5 234.8z",fill:e}},{tag:"path",attrs:{d:"M512 140c-205.4 0-372 166.6-372 372s166.6 372 372 372c89 0 170.8-31.3 234.8-83.5L224 277.7c16-19.7 34-37.7 53.7-53.7l522.8 522.8C852.7 682.8 884 601 884 512c0-205.4-166.6-372-372-372z",fill:n}}]}},name:"stop",theme:"twotone"},o=t(42135),c=function(e,n){return a.createElement(o.Z,(0,r.Z)((0,r.Z)({},e),{},{ref:n,icon:s}))};c.displayName="StopTwoTone";var i=a.forwardRef(c)},27049:function(e,n,t){var r=t(87462),a=t(4942),s=t(67294),o=t(94184),c=t.n(o),i=t(59844),l=function(e,n){var t={};for(var r in e)Object.prototype.hasOwnProperty.call(e,r)&&n.indexOf(r)<0&&(t[r]=e[r]);if(null!=e&&"function"===typeof Object.getOwnPropertySymbols){var a=0;for(r=Object.getOwnPropertySymbols(e);a<r.length;a++)n.indexOf(r[a])<0&&Object.prototype.propertyIsEnumerable.call(e,r[a])&&(t[r[a]]=e[r[a]])}return t};n.Z=function(e){return s.createElement(i.C,null,(function(n){var t,o=n.getPrefixCls,i=n.direction,u=e.prefixCls,d=e.type,f=void 0===d?"horizontal":d,m=e.orientation,h=void 0===m?"center":m,p=e.className,v=e.children,x=e.dashed,b=e.plain,j=l(e,["prefixCls","type","orientation","className","children","dashed","plain"]),y=o("divider",u),g=h.length>0?"-".concat(h):h,Z=!!v,w=c()(y,"".concat(y,"-").concat(f),(t={},(0,a.Z)(t,"".concat(y,"-with-text"),Z),(0,a.Z)(t,"".concat(y,"-with-text").concat(g),Z),(0,a.Z)(t,"".concat(y,"-dashed"),!!x),(0,a.Z)(t,"".concat(y,"-plain"),!!b),(0,a.Z)(t,"".concat(y,"-rtl"),"rtl"===i),t),p);return s.createElement("div",(0,r.Z)({className:w},j,{role:"separator"}),v&&s.createElement("span",{className:"".concat(y,"-inner-text")},v))}))}},66192:function(e,n,t){t.d(n,{Z:function(){return h}});var r=t(28520),a=t.n(r),s=t(85893),o=t(56516),c=t(71577),i=t(21640),l=t(82911),u=t(84674),d=t(58827);function f(e,n,t,r,a,s,o){try{var c=e[s](o),i=c.value}catch(l){return void t(l)}c.done?n(i):Promise.resolve(i).then(r,a)}function m(e){return function(){var n=this,t=arguments;return new Promise((function(r,a){var s=e.apply(n,t);function o(e){f(s,r,a,o,c,"next",e)}function c(e){f(s,r,a,o,c,"throw",e)}o(void 0)}))}}function h(e){var n=e.user,t=e.isEnabled,r=e.label,f=e.onClick;function h(){return(h=m(a().mark((function e(n){var r,s,o;return a().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return r=n.id,s={userId:r,enabled:!t},e.prev=2,e.next=5,(0,d.rQ)(d.NM,{data:s,method:"POST",auth:!0});case 5:return o=e.sent,e.abrupt("return",o.success);case 9:e.prev=9,e.t0=e.catch(2),console.error(e.t0);case 12:return e.abrupt("return",!1);case 13:case"end":return e.stop()}}),e,null,[[2,9]])})))).apply(this,arguments)}var p=t?"ban":"unban",v=t?(0,s.jsx)(i.Z,{style:{color:"var(--ant-error)"}}):(0,s.jsx)(l.Z,{style:{color:"var(--ant-warning)"}}),x=(0,s.jsxs)(s.Fragment,{children:["Are you sure you want to ",p," ",(0,s.jsx)("strong",{children:n.displayName}),t?" and remove their messages?":"?"]});return(0,s.jsx)(c.Z,{onClick:function(){o.Z.confirm({title:"Confirm ".concat(p),content:x,onCancel:function(){},onOk:function(){return new Promise((function(e,t){var r=function(e){return h.apply(this,arguments)}(n);r?setTimeout((function(){e(r),null===f||void 0===f||f()}),3e3):t()}))},okType:"danger",okText:t?"Absolutely":null,icon:v})},size:"small",icon:t?(0,s.jsx)(u.Z,{twoToneColor:"#ff4d4f"}):null,className:"block-user-button",children:r||p})}h.defaultProps={label:"",onClick:null}},85584:function(e,n,t){t.d(n,{Z:function(){return T}});var r=t(85893),a=t(67294),s=t(56266),o=t(56516),c=t(17256),i=t(25968),l=t(6226),u=t(27049),d=t(85533),f=t(58091),m=t(96486),h=t(66192),p=t(28520),v=t.n(p),x=t(71577),b=t(21640),j=t(82911),y=t(84674),g=t(58827);function Z(e,n,t,r,a,s,o){try{var c=e[s](o),i=c.value}catch(l){return void t(l)}c.done?n(i):Promise.resolve(i).then(r,a)}function w(e){return function(){var n=this,t=arguments;return new Promise((function(r,a){var s=e.apply(n,t);function o(e){Z(s,r,a,o,c,"next",e)}function c(e){Z(s,r,a,o,c,"throw",e)}o(void 0)}))}}function C(e){var n,t=e.user,a=e.onClick;function s(){return(s=w(v().mark((function e(n,t){var r,a,s;return v().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return r=n.id,a={userId:r,isModerator:t},e.prev=2,e.next=5,(0,g.rQ)(g.jr,{data:a,method:"POST",auth:!0});case 5:return s=e.sent,e.abrupt("return",s.success);case 9:e.prev=9,e.t0=e.catch(2),console.error(e.t0);case 12:return e.abrupt("return",!1);case 13:case"end":return e.stop()}}),e,null,[[2,9]])})))).apply(this,arguments)}var c=null===(n=t.scopes)||void 0===n?void 0:n.includes("MODERATOR"),i=c?"remove moderator":"add moderator",l=c?(0,r.jsx)(b.Z,{style:{color:"var(--ant-error)"}}):(0,r.jsx)(j.Z,{style:{color:"var(--ant-warning)"}}),u=(0,r.jsxs)(r.Fragment,{children:["Are you sure you want to ",i," ",(0,r.jsx)("strong",{children:t.displayName}),"?"]});return(0,r.jsx)(x.Z,{onClick:function(){o.Z.confirm({title:"Confirm ".concat(i),content:u,onCancel:function(){},onOk:function(){return new Promise((function(e,n){var r=function(e,n){return s.apply(this,arguments)}(t,!c);r?setTimeout((function(){e(r),null===a||void 0===a||a()}),3e3):n()}))},okType:"danger",okText:c?"Yup!":null,icon:l})},size:"small",icon:c?(0,r.jsx)(y.Z,{twoToneColor:"#ff4d4f"}):null,className:"block-user-button",children:i})}C.defaultProps={onClick:null};var k=t(20643),N=t(2766);function A(e){return function(e){if(Array.isArray(e)){for(var n=0,t=new Array(e.length);n<e.length;n++)t[n]=e[n];return t}}(e)||function(e){if(Symbol.iterator in Object(e)||"[object Arguments]"===Object.prototype.toString.call(e))return Array.from(e)}(e)||function(){throw new TypeError("Invalid attempt to spread non-iterable instance")}()}function T(e){var n=e.user,t=e.connectionInfo,p=e.children,v=(0,a.useState)(!1),x=v[0],b=v[1],j=function(){b(!1)},y=n.displayName,g=n.createdAt,Z=n.previousNames,w=n.nameChangedAt,T=n.disabledAt,O=t||{},P=O.connectedAt,D=O.messageCount,E=O.userAgent,S=null,z=Z&&A(Z);Z&&Z.length>1&&w&&(S=new Date(w),z.reverse());var M=new Date(g),I=(0,f.Z)(M,"PP pp"),F=S?(0,d.Z)(S):null;return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(s.Z,{title:(0,r.jsxs)(r.Fragment,{children:["Created at: ",I,".",(0,r.jsx)("br",{})," Click for more info."]}),placement:"bottomLeft",children:(0,r.jsx)("button",{type:"button","aria-label":"Display more details about this user",className:"user-item-container",onClick:function(){b(!0)},children:p})}),(0,r.jsx)(o.Z,{destroyOnClose:!0,width:600,cancelText:"Close",okButtonProps:{style:{display:"none"}},title:"User details: ".concat(y),visible:x,onOk:j,onCancel:j,children:(0,r.jsxs)("div",{className:"user-details",children:[(0,r.jsx)(c.Z.Title,{level:4,children:y}),(0,r.jsxs)("p",{className:"created-at",children:["User created at ",I,"."]}),(0,r.jsxs)(i.Z,{gutter:16,children:[t&&(0,r.jsxs)(l.Z,{md:S?12:24,children:[(0,r.jsx)(c.Z.Title,{level:5,children:"This user is currently connected to Chat."}),(0,r.jsxs)("ul",{className:"connection-info",children:[(0,r.jsxs)("li",{children:[(0,r.jsx)("strong",{children:"Active for:"})," ",(0,d.Z)(new Date(P))]}),(0,r.jsxs)("li",{children:[(0,r.jsx)("strong",{children:"Messages sent:"})," ",D]}),(0,r.jsxs)("li",{children:[(0,r.jsx)("strong",{children:"User Agent:"}),(0,r.jsx)("br",{}),(0,N.AB)(E)]})]})]}),S&&(0,r.jsxs)(l.Z,{md:t?12:24,children:[(0,r.jsx)(c.Z.Title,{level:5,children:"This user is also seen as:"}),(0,r.jsx)("ul",{className:"previous-names-list",children:(0,m.uniq)(z).map((function(e,n){return(0,r.jsxs)("li",{className:0===n?"latest":"",children:[(0,r.jsx)("span",{className:"user-name-item",children:e}),0===n?" (Changed ".concat(F," ago)"):""]})}))})]})]}),(0,r.jsx)(u.Z,{}),T?(0,r.jsxs)(r.Fragment,{children:["This user was banned on ",(0,r.jsx)("code",{children:(0,k.u)(T)}),".",(0,r.jsx)("br",{}),(0,r.jsx)("br",{}),(0,r.jsx)(h.Z,{label:"Unban this user",user:n,isEnabled:!1,onClick:j})]}):(0,r.jsx)(h.Z,{label:"Ban this user",user:n,isEnabled:!0,onClick:j}),(0,r.jsx)(C,{user:n,onClick:j})]})})]})}T.defaultProps={connectionInfo:null}},20643:function(e,n,t){t.d(n,{u:function(){return i},Z:function(){return l}});var r=t(85893),a=t(73555),s=t(58091),o=t(85584),c=t(66192);function i(e){return(0,s.Z)(new Date(e),"MMM d H:mma")}function l(e){var n=e.data,t=[{title:"Last Known Display Name",dataIndex:"displayName",key:"displayName",render:function(e,n){return(0,r.jsx)(o.Z,{user:n,children:(0,r.jsx)("span",{className:"display-name",children:e})})}},{title:"Created",dataIndex:"createdAt",key:"createdAt",render:function(e){return i(e)},sorter:function(e,n){return new Date(e.createdAt).getTime()-new Date(n.createdAt).getTime()},sortDirections:["descend","ascend"]},{title:"Disabled at",dataIndex:"disabledAt",key:"disabledAt",defaultSortOrder:"descend",render:function(e){return e?i(e):null},sorter:function(e,n){return new Date(e.disabledAt).getTime()-new Date(n.disabledAt).getTime()},sortDirections:["descend","ascend"]},{title:"",key:"block",className:"actions-col",render:function(e,n){return(0,r.jsx)(c.Z,{user:n,isEnabled:!n.disabledAt})}}];return(0,r.jsx)(a.Z,{pagination:{hideOnSinglePage:!0},className:"table-container",columns:t,dataSource:n,size:"small",rowKey:"id"})}}}]);