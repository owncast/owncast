(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[145],{19481:function(e,t,n){"use strict";n.r(t),n.d(t,{OUTCOME_TIMEOUT:function(){return C},default:function(){return E}});var r=n(4942),s=n(93433),i=n(15861),a=n(87757),c=n.n(a),o=n(67294),u=n(4525),l=n(71577),d=n(99645),f=n(38819),p=n(68855),h=n(94184),m=n.n(h),v=n(12924),g=n(94853),b=n(9431),y=n(31097),w=n(95357),j=n(88633),x=n(85893);function O(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function N(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?O(Object(n),!0).forEach((function(t){(0,r.Z)(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):O(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function P(e){var t=e.isVisible,n=e.message,r=e.setMessage;if(!n||(0,b.Qr)(n))return null;var s=null,a=(0,o.useState)(0),u=a[0],d=a[1],h=(n||{}).id;(0,o.useEffect)((function(){return function(){clearTimeout(s)}}));var m=function(){var e=(0,i.Z)(c().mark((function e(){var i;return c().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return clearTimeout(s),d(0),e.next=4,(0,g.rQ)(g.hn,{auth:!0,method:"POST",data:{visible:!t,idArray:[h]}});case 4:(i=e.sent).success&&"changed"===i.message?(r(N(N({},n),{},{visible:!t})),d(1)):(r(N(N({},n),{},{visible:t})),d(-1)),s=setTimeout((function(){d(0)}),C);case 7:case"end":return e.stop()}}),e)})));return function(){return e.apply(this,arguments)}}(),v=(0,x.jsx)(f.Z,{style:{color:"transparent"}});u&&(v=u>0?(0,x.jsx)(f.Z,{style:{color:"var(--ant-success)"}}):(0,x.jsx)(p.Z,{style:{color:"var(--ant-warning)"}}));var O="Click to ".concat(t?"hide":"show"," this message");return(0,x.jsxs)("div",{className:"toggle-switch ".concat(t?"":"hidden"),children:[(0,x.jsx)("span",{className:"outcome-icon",children:v}),(0,x.jsx)(y.Z,{title:O,placement:"topRight",children:(0,x.jsx)(l.Z,{shape:"circle",size:"small",type:"text",icon:t?(0,x.jsx)(w.Z,{}):(0,x.jsx)(j.Z,{}),onClick:m})})]})}var Z=n(99382);function k(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function S(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?k(Object(n),!0).forEach((function(t){(0,r.Z)(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):k(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}var T=u.Z.Title;var C=3e3;function E(){var e=(0,o.useState)([]),t=e[0],n=e[1],r=(0,o.useState)([]),a=r[0],u=r[1],h=(0,o.useState)(!1),y=h[0],w=h[1],j=(0,o.useState)(null),O=j[0],N=j[1],k=(0,o.useState)(""),E=k[0],_=k[1],D=null,I=null,M=function(){var e=(0,i.Z)(c().mark((function e(){var t;return c().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,(0,g.rQ)(g.WE,{auth:!0});case 3:t=e.sent,(0,b.Qr)(t)?n([]):n(t),e.next=10;break;case 7:e.prev=7,e.t0=e.catch(0),console.log("==== error",e.t0);case 10:case"end":return e.stop()}}),e,null,[[0,7]])})));return function(){return e.apply(this,arguments)}}();(0,o.useEffect)((function(){return M(),I=setInterval((function(){M()}),g.NE),function(){clearTimeout(D),clearTimeout(I)}}),[]);var z=function(e){return e.reduce((function(e,t){var n=t.user.id;return e.some((function(e){return e.text===n}))||e.push({text:n,value:n}),e}),[]).sort((function(e,t){var n=e.text.toUpperCase(),r=t.text.toUpperCase();return n<r?-1:n>r?1:0}))}(t),A={selectedRowKeys:a,onChange:function(e){u(e)}},Q=function(e){var r=t.findIndex((function(t){return t.id===e.id}));t.splice(r,1,e),n((0,s.Z)(t))},U=function(){D=setTimeout((function(){N(null),_("")}),C)},H=function(){var e=(0,i.Z)(c().mark((function e(r){var i,o;return c().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return w(!0),e.next=3,(0,g.rQ)(g.hn,{auth:!0,method:"POST",data:{visible:r,idArray:a}});case 3:(i=e.sent).success&&"changed"===i.message?(N((0,x.jsx)(f.Z,{})),U(),o=(0,s.Z)(t),a.map((function(e){var n=o.findIndex((function(t){return t.id===e})),s=S(S({},t[n]),{},{visible:r});return o.splice(n,1,s),null})),n(o),u([])):(N((0,x.jsx)(p.Z,{})),U()),w(!1);case 6:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}(),V=[{title:"Time",dataIndex:"timestamp",key:"timestamp",className:"timestamp-col",defaultSortOrder:"descend",render:function(e){var t=new Date(e);return(0,v.Z)(t,"PP pp")},sorter:function(e,t){return new Date(e.timestamp).getTime()-new Date(t.timestamp).getTime()},width:90},{title:"User",dataIndex:"user",key:"user",className:"name-col",filters:z,onFilter:function(e,t){return t.user.id===e},sorter:function(e,t){return e.user.displayName.localeCompare(t.user.displayName)},sortDirections:["ascend","descend"],ellipsis:!0,render:function(e){var t=e.displayName;return(0,x.jsx)(Z.Z,{user:e,children:t})},width:110},{title:"Message",dataIndex:"body",key:"body",className:"message-col",width:320,render:function(e){return(0,x.jsx)("div",{className:"message-contents",dangerouslySetInnerHTML:{__html:e}})}},{title:"",dataIndex:"hiddenAt",key:"hiddenAt",className:"toggle-col",filters:[{text:"Visible messages",value:!0},{text:"Hidden messages",value:!1}],onFilter:function(e,t){return t.visible===e},render:function(e,t){return(0,x.jsx)(P,{isVisible:!e,message:t,setMessage:Q})},width:30}],F=m()({"bulk-editor":!0,active:a.length});return(0,x.jsxs)("div",{className:"chat-messages",children:[(0,x.jsx)(T,{children:"Chat Messages"}),(0,x.jsx)("p",{children:"Manage the messages from viewers that show up on your stream."}),(0,x.jsxs)("div",{className:F,children:[(0,x.jsx)("span",{className:"label",children:"Check multiple messages to change their visibility to: "}),(0,x.jsx)(l.Z,{type:"primary",size:"small",shape:"round",className:"button",loading:"show"===E&&y,icon:"show"===E&&O,disabled:!a.length||E&&"show"!==E,onClick:function(){_("show"),H(!0)},children:"Show"}),(0,x.jsx)(l.Z,{type:"primary",size:"small",shape:"round",className:"button",loading:"hide"===E&&y,icon:"hide"===E&&O,disabled:!a.length||E&&"hide"!==E,onClick:function(){_("hide"),H(!1)},children:"Hide"})]}),(0,x.jsx)(d.Z,{size:"small",className:"table-container",pagination:{defaultPageSize:100,showSizeChanger:!0},scroll:{y:540},rowClassName:function(e){return e.hiddenAt?"hidden":""},dataSource:t,columns:V,rowKey:function(e){return e.id},rowSelection:A})]})}},9007:function(e,t,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/chat/messages",function(){return n(19481)}])}},function(e){e.O(0,[662,645,924,614,354,821,774,888,179],(function(){return t=9007,e(e.s=t);var t}));var t=e.O();_N_E=t}]);