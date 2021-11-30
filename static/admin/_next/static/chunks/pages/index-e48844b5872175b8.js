(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[405],{45301:function(e,n,t){(window.__NEXT_P=window.__NEXT_P||[]).push(["/",function(){return t(94831)}])},824:function(e,n,t){"use strict";t.d(n,{Z:function(){return m}});var r=t(85893),s=(t(67294),t(44068)),i=t(20550),a=t(88829),o=t(53731),l=t(58091),c=s.Z.Title;function d(e,n){var t="black";return"warning"===n.level?t="orange":"error"===n.level&&(t="red"),(0,r.jsx)(i.Z,{color:t,children:e})}function u(e){return(0,r.jsx)(o.Z,{children:e})}function m(e){var n=e.logs,t=e.pageSize;if(!(null===n||void 0===n?void 0:n.length))return null;var s=[{title:"Level",dataIndex:"level",key:"level",filters:[{text:"Info",value:"info"},{text:"Warning",value:"warning"},{text:"Error",value:"Error"}],onFilter:function(e,n){return 0===n.level.indexOf(e)},render:d},{title:"Timestamp",dataIndex:"time",key:"time",render:function(e){var n=new Date(e);return(0,l.Z)(n,"pp P")},sorter:function(e,n){return new Date(e.time).getTime()-new Date(n.time).getTime()},sortDirections:["descend","ascend"],defaultSortOrder:"descend"},{title:"Message",dataIndex:"message",key:"message",render:u}];return(0,r.jsxs)("div",{className:"logs-section",children:[(0,r.jsx)(c,{children:"Logs"}),(0,r.jsx)(a.Z,{size:"middle",dataSource:n,columns:s,rowKey:function(e){return e.time},pagination:{pageSize:t||20}})]})}},94831:function(e,n,t){"use strict";t.r(n),t.d(n,{default:function(){return K}});var r=t(28520),s=t.n(r),i=t(85893),a=t(67294),o=t(41080),l=t(74763),c=t(97751),d=t(25968),u=t(6226),m=t(24019),h=t(87547),f=t(26555),v=t(85533),x=t(35159),j=t(824),p=t(66567),w=t(63179),g=t(78346),Z=t(27482),b=t(43439),y=t(44068),k=t(41664),N=t(92659),S=t(54907),C=t(58091),O=t(58827);function _(e,n,t,r,s,i,a){try{var o=e[i](a),l=o.value}catch(c){return void t(c)}o.done?n(l):Promise.resolve(l).then(r,s)}function P(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}var D=S.Z.Panel,E=y.Z.Title,T=y.Z.Link;function I(e){var n=e.title,t=e.url,r=e.content_html,s=e.date_published,a=new Date(s),o=(0,C.Z)(a,"MMM dd, yyyy, HH:mm");return(0,i.jsx)("article",{children:(0,i.jsx)(S.Z,{children:(0,i.jsxs)(D,{header:n,children:[(0,i.jsxs)("p",{className:"timestamp",children:[o," (",(0,i.jsx)(T,{href:"".concat("https://owncast.online").concat(t),target:"_blank",rel:"noopener noreferrer",children:"Link"}),")"]}),(0,i.jsx)("div",{dangerouslySetInnerHTML:{__html:r}})]},t)})})}function U(){var e,n=(0,a.useState)([]),t=n[0],r=n[1],l=(0,a.useState)(!0),c=l[0],d=l[1],u=(e=s().mark((function e(){var n;return s().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return d(!1),e.prev=1,e.next=4,(0,O.kg)("https://owncast.online/news/index.json");case 4:(null===(n=e.sent)||void 0===n?void 0:n.items.length)>0&&r(n.items),e.next=11;break;case 8:e.prev=8,e.t0=e.catch(1),console.log("==== error",e.t0);case 11:case"end":return e.stop()}}),e,null,[[1,8]])})),function(){var n=this,t=arguments;return new Promise((function(r,s){var i=e.apply(n,t);function a(e){_(i,r,s,a,o,"next",e)}function o(e){_(i,r,s,a,o,"throw",e)}a(void 0)}))});(0,a.useEffect)((function(){u()}),[]);var m=c?(0,i.jsx)(o.Z,{loading:!0,active:!0}):null,h=c||0!==t.length?null:(0,i.jsx)("div",{children:"No news."});return(0,i.jsxs)("section",{className:"news-feed form-module",children:[(0,i.jsx)(E,{level:2,children:"News & Updates from Owncast"}),m,t.map((function(e){return(0,a.createElement)(I,function(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{},r=Object.keys(t);"function"===typeof Object.getOwnPropertySymbols&&(r=r.concat(Object.getOwnPropertySymbols(t).filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable})))),r.forEach((function(n){P(e,n,t[n])}))}return e}({},e,{key:e.url}))})),h]})}var z=y.Z.Paragraph,L=y.Z.Text,B=y.Z.Title,M=c.Z.Meta;function F(e,n){return"rtmp://".concat(e.replace(/(^\w+:|^)\/\//,""),":").concat(n,"/live/")}function Q(e){var n,r,s=e.logs,o=void 0===s?[]:s,l=e.config,m=((0,a.useContext)(x.aC)||{}).serverConfig,h=m.streamKey,f=m.rtmpServerPort,v=(null===(n=t.g.window)||void 0===n?void 0:n.location.hostname)||"",y=[{icon:(0,i.jsx)(p.Z,{twoToneColor:"#6f42c1"}),title:"Use your broadcasting software",content:(0,i.jsxs)("div",{children:[(0,i.jsx)("a",{href:"https://owncast.online/docs/broadcasting/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn how to point your existing software to your new server and start streaming your content."}),(0,i.jsxs)("div",{className:"stream-info-container",children:[(0,i.jsx)(L,{strong:!0,className:"stream-info-label",children:"Streaming URL:"}),(0,i.jsx)(z,{className:"stream-info-box",copyable:!0,children:F(v,f)}),(0,i.jsx)(L,{strong:!0,className:"stream-info-label",children:"Stream Key:"}),(0,i.jsx)(z,{className:"stream-info-box",copyable:{text:h},children:"*********************"})]})]})},{icon:(0,i.jsx)(w.Z,{twoToneColor:"#f9826c"}),title:"Embed your video onto other sites",content:(0,i.jsx)("div",{children:(0,i.jsx)("a",{href:"https://owncast.online/docs/embed?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn how you can add your Owncast stream to other sites you control."})})}];return(null===l||void 0===l?void 0:l.chatDisabled)||y.push({icon:(0,i.jsx)(g.Z,{twoToneColor:"#0366d6"}),title:"Chat is disabled",content:(0,i.jsx)("span",{children:"Chat will continue to be disabled until you begin a live stream."})}),(null===l||void 0===l||null===(r=l.yp)||void 0===r?void 0:r.enabled)||y.push({icon:(0,i.jsx)(Z.Z,{twoToneColor:"#D18BFE"}),title:"Find an audience on the Owncast Directory",content:(0,i.jsxs)("div",{children:["List yourself in the Owncast Directory and show off your stream. Enable it in"," ",(0,i.jsx)(k.default,{href:"/config-public-details",children:"settings."})]})}),y.push({icon:(0,i.jsx)(b.Z,{twoToneColor:"#ffd33d"}),title:"Not sure what to do next?",content:(0,i.jsxs)("div",{children:["If you're having issues or would like to know how to customize and configure your Owncast server visit ",(0,i.jsx)(k.default,{href:"/help",children:"the help page."})]})}),(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(d.Z,{children:(0,i.jsx)(u.Z,{span:12,offset:6,children:(0,i.jsxs)("div",{className:"offline-intro",children:[(0,i.jsx)("span",{className:"logo",children:(0,i.jsx)(N.Z,{})}),(0,i.jsxs)("div",{children:[(0,i.jsx)(B,{level:2,children:"No stream is active"}),(0,i.jsx)("p",{children:"You should start one."})]})]})})}),(0,i.jsxs)(d.Z,{gutter:[16,16],className:"offline-content",children:[(0,i.jsx)(u.Z,{span:12,xs:24,sm:24,md:24,lg:12,className:"list-section",children:y.map((function(e){return(0,i.jsx)(c.Z,{size:"small",bordered:!1,children:(0,i.jsx)(M,{avatar:e.icon,title:e.title,description:e.content})},e.title)}))}),(0,i.jsx)(u.Z,{span:12,xs:24,sm:24,md:24,lg:12,children:(0,i.jsx)(U,{})})]}),(0,i.jsx)(j.Z,{logs:o,pageSize:5})]})}var V=t(2766);function A(e,n,t,r,s,i,a){try{var o=e[i](a),l=o.value}catch(c){return void t(c)}o.done?n(l):Promise.resolve(l).then(r,s)}function H(e){return(0,i.jsxs)("ul",{className:"statistics-list",children:[(0,i.jsxs)("li",{children:[e.videoCodec||"Unknown"," @ ",e.videoBitrate||"Unknown"," kbps"]}),(0,i.jsxs)("li",{children:[e.framerate||"Unknown"," fps"]}),(0,i.jsxs)("li",{children:[e.width," x ",e.height]})]})}function K(){var e,n,t,r=(0,a.useContext)(x.aC),p=r||{},w=p.broadcaster,g=p.serverConfig,Z=w||{},b=Z.remoteAddr,y=Z.streamDetails,k=(null===y||void 0===y?void 0:y.encoder)||"Unknown encoder",N=(0,a.useState)([]),S=N[0],C=N[1],_=(t=s().mark((function e(){var n;return s().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,(0,O.rQ)(O.WQ);case 3:n=e.sent,C(n),e.next=10;break;case 7:e.prev=7,e.t0=e.catch(0),console.log("==== error",e.t0);case 10:case"end":return e.stop()}}),e,null,[[0,7]])})),function(){var e=this,n=arguments;return new Promise((function(r,s){var i=t.apply(e,n);function a(e){A(i,r,s,a,o,"next",e)}function o(e){A(i,r,s,a,o,"throw",e)}a(void 0)}))}),P=function(){_()};if((0,a.useEffect)((function(){P();var e;return e=setInterval(P,O.NE),function(){clearInterval(e)}}),[]),(0,V.Qr)(g)||(0,V.Qr)(r))return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(o.Z,{active:!0}),(0,i.jsx)(o.Z,{active:!0}),(0,i.jsx)(o.Z,{active:!0})]});if(!w)return(0,i.jsx)(Q,{logs:S,config:g});var D=null===r||void 0===r||null===(e=r.currentBroadcast)||void 0===e||null===(n=e.outputSettings)||void 0===n?void 0:n.map((function(e){var n=e.audioPassthrough,t=e.videoPassthrough,r=e.audioBitrate,s=e.videoBitrate,a=e.framerate,o=n?"".concat(y.audioCodec||"Unknown",", ").concat(y.audioBitrate," kbps"):"".concat(r||"Unknown"," kbps"),c=t?"".concat(y.videoBitrate||"Unknown"," kbps, ").concat(y.framerate," fps ").concat(y.width," x ").concat(y.height):"".concat(s||"Unknown"," kbps, ").concat(a," fps");return(0,i.jsxs)("div",{className:"stream-details-item-container",children:[(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Outbound Video Stream",value:c}),(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Outbound Audio Stream",value:o})]})})),E=r.viewerCount,T=r.sessionPeakViewerCount,I="".concat(y.audioCodec,", ").concat(y.audioBitrate||"Unknown"," kbps"),z=new Date(w.time);return(0,i.jsxs)("div",{className:"home-container",children:[(0,i.jsxs)("div",{className:"sections-container",children:[(0,i.jsx)("div",{className:"online-status-section",children:(0,i.jsx)(c.Z,{size:"small",type:"inner",className:"online-details-card",children:(0,i.jsxs)(d.Z,{gutter:[16,16],align:"middle",children:[(0,i.jsx)(u.Z,{span:8,sm:24,md:8,children:(0,i.jsx)(l.Z,{title:"Stream started ".concat((0,f.Z)(z,Date.now())),value:(0,v.Z)(z),prefix:(0,i.jsx)(m.Z,{})})}),(0,i.jsx)(u.Z,{span:8,sm:24,md:8,children:(0,i.jsx)(l.Z,{title:"Viewers",value:E,prefix:(0,i.jsx)(h.Z,{})})}),(0,i.jsx)(u.Z,{span:8,sm:24,md:8,children:(0,i.jsx)(l.Z,{title:"Peak viewer count",value:T,prefix:(0,i.jsx)(h.Z,{})})})]})})}),(0,i.jsxs)(d.Z,{gutter:[16,16],className:"section stream-details-section",children:[(0,i.jsxs)(u.Z,{className:"stream-details",span:12,sm:24,md:24,lg:12,children:[(0,i.jsx)(c.Z,{size:"small",title:"Outbound Stream Details",type:"inner",className:"outbound-details",children:D}),(0,i.jsxs)(c.Z,{size:"small",title:"Inbound Stream Details",type:"inner",children:[(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Input",value:"".concat(k," ").concat((0,V.t5)(b))}),(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Inbound Video Stream",value:y,formatter:H}),(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Inbound Audio Stream",value:I})]})]}),(0,i.jsx)(u.Z,{span:12,xs:24,sm:24,md:24,lg:12,children:(0,i.jsx)(U,{})})]})]}),(0,i.jsx)("br",{}),(0,i.jsx)(j.Z,{logs:S,pageSize:5})]})}}},function(e){e.O(0,[829,91,903,138,533,429,59,774,888,179],(function(){return n=45301,e(e.s=n);var n}));var n=e.O();_N_E=n}]);