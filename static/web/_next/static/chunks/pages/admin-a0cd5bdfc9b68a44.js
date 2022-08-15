(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[6964],{58860:function(e,n,t){(window.__NEXT_P=window.__NEXT_P||[]).push(["/admin",function(){return t(1025)}])},19105:function(e,n,t){"use strict";t.d(n,{Z:function(){return h}});var r=t(85893),s=(t(67294),t(84916)),i=t(20550),a=t(21213),o=t(53731),l=t(98147),c=s.Z.Title;function d(e,n){var t="black";return"warning"===n.level?t="orange":"error"===n.level&&(t="red"),(0,r.jsx)(i.Z,{color:t,children:e})}function u(e){return(0,r.jsx)(o.Z,{children:e})}function h(e){var n=e.logs,t=e.pageSize;if(!(null===n||void 0===n?void 0:n.length))return null;var s=[{title:"Level",dataIndex:"level",key:"level",filters:[{text:"Info",value:"info"},{text:"Warning",value:"warning"},{text:"Error",value:"Error"}],onFilter:function(e,n){return 0===n.level.indexOf(e)},render:d},{title:"Timestamp",dataIndex:"time",key:"time",render:function(e){var n=new Date(e);return(0,l.Z)(n,"pp P")},sorter:function(e,n){return new Date(e.time).getTime()-new Date(n.time).getTime()},sortDirections:["descend","ascend"],defaultSortOrder:"descend"},{title:"Message",dataIndex:"message",key:"message",render:u}];return(0,r.jsxs)("div",{className:"logs-section",children:[(0,r.jsx)(c,{children:"Logs"}),(0,r.jsx)(a.Z,{size:"middle",dataSource:n,columns:s,rowKey:function(e){return e.time},pagination:{pageSize:t||20}})]})}},35660:function(e,n,t){"use strict";t.d(n,{Z:function(){return v}});var r=t(85893),s=t(8751),i=t(11475),a=t(25968),o=t(6226),l=t(74763),c=t(84916),d=t(14670),u=t(71577),h=t(41664),m=t.n(h),f=t(67294),x=t(99519);function v(e){var n=e.showTroubleshootButton,t=(0,f.useContext)(x.aC).health;if(!t)return null;var h=t.healthy,v=t.healthPercentage,p=t.message,j=t.representation,g="#3f8600",w="info";return v<80?(g="#cf000f",w="error"):v<30&&(g="#f0ad4e",w="error"),(0,r.jsxs)("div",{children:[(0,r.jsxs)(a.Z,{gutter:8,children:[(0,r.jsx)(o.Z,{span:12,children:(0,r.jsx)(l.Z,{title:"Healthy Stream",value:h?"Yes":"No",valueStyle:{color:g},prefix:h?(0,r.jsx)(s.Z,{}):(0,r.jsx)(i.Z,{})})}),(0,r.jsx)(o.Z,{span:12,children:(0,r.jsx)(l.Z,{title:"Playback Health",value:v,valueStyle:{color:g},suffix:"%"})})]}),(0,r.jsx)(a.Z,{style:{display:j<100&&0!==j?"grid":"none"},children:(0,r.jsxs)(c.Z.Text,{type:"secondary",style:{textAlign:"center",fontSize:"0.7em",opacity:"0.3"},children:["Stream health represents ",j,"% of all known players. Other player status is unknown."]})}),(0,r.jsx)(a.Z,{gutter:16,style:{width:"100%",display:p?"grid":"none",marginTop:"10px"},children:(0,r.jsx)(o.Z,{span:24,children:(0,r.jsx)(d.Z,{message:p,type:w,showIcon:!0,action:n&&(0,r.jsx)(m(),{passHref:!0,href:"/stream-health",children:(0,r.jsx)(u.Z,{size:"small",type:"text",style:{color:"black"},children:"TROUBLESHOOT"})})})})})]})}v.defaultProps={showTroubleshootButton:!0}},1025:function(e,n,t){"use strict";t.r(n),t.d(n,{default:function(){return K}});var r=t(34051),s=t.n(r),i=t(85893),a=t(67294),o=t(41080),l=t(74763),c=t(97751),d=t(25968),u=t(6226),h=t(24019),m=t(87547),f=t(19642),x=t(45938),v=t(99519),p=t(19105),j=t(66567),g=t(63179),w=t(78346),Z=t(27482),y=t(84916),b=t(41664),k=t.n(b),N=t(4213),S=t(54907),O=t(98147),C=t(64777);function P(e,n,t,r,s,i,a){try{var o=e[i](a),l=o.value}catch(c){return void t(c)}o.done?n(l):Promise.resolve(l).then(r,s)}function T(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}var _=S.Z.Panel,E=y.Z.Title,D=y.Z.Link;function U(e){var n=e.title,t=e.url,r=e.content_html,s=e.date_published,a=new Date(s),o=(0,O.Z)(a,"MMM dd, yyyy, HH:mm");return(0,i.jsx)("article",{children:(0,i.jsx)(S.Z,{children:(0,i.jsxs)(_,{header:n,children:[(0,i.jsxs)("p",{className:"timestamp",children:[o," (",(0,i.jsx)(D,{href:"".concat("https://owncast.online").concat(t),target:"_blank",rel:"noopener noreferrer",children:"Link"}),")"]}),(0,i.jsx)("div",{dangerouslySetInnerHTML:{__html:r}})]},t)})})}function I(){var e=(0,a.useState)([]),n=e[0],t=e[1],r=(0,a.useState)(!0),l=r[0],c=r[1],d=function(){var e,n=(e=s().mark((function e(){var n;return s().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return c(!1),e.prev=1,e.next=4,(0,C.kg)("https://owncast.online/news/index.json");case 4:(null===(n=e.sent)||void 0===n?void 0:n.items.length)>0&&t(n.items),e.next=11;break;case 8:e.prev=8,e.t0=e.catch(1),console.log("==== error",e.t0);case 11:case"end":return e.stop()}}),e,null,[[1,8]])})),function(){var n=this,t=arguments;return new Promise((function(r,s){var i=e.apply(n,t);function a(e){P(i,r,s,a,o,"next",e)}function o(e){P(i,r,s,a,o,"throw",e)}a(void 0)}))});return function(){return n.apply(this,arguments)}}();(0,a.useEffect)((function(){d()}),[]);var u=l?(0,i.jsx)(o.Z,{loading:!0,active:!0}):null,h=l||0!==n.length?null:(0,i.jsx)("div",{children:"No news."});return(0,i.jsxs)("section",{className:"news-feed form-module",children:[(0,i.jsx)(E,{level:2,children:"News & Updates from Owncast"}),u,n.map((function(e){return(0,a.createElement)(U,function(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{},r=Object.keys(t);"function"===typeof Object.getOwnPropertySymbols&&(r=r.concat(Object.getOwnPropertySymbols(t).filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable})))),r.forEach((function(n){T(e,n,t[n])}))}return e}({},e,{key:e.url}))})),h]})}var z=y.Z.Paragraph,B=y.Z.Text,L=y.Z.Title,F=c.Z.Meta;function H(e){var n,r,s,o,l=e.logs,h=void 0===l?[]:l,m=e.config,f=((0,a.useContext)(v.aC)||{}).serverConfig,x=f.streamKey,y=f.rtmpServerPort,b=(null===(n=t.g.window)||void 0===n?void 0:n.location.hostname)||"";b&&y&&(o=function(e,n){return"rtmp://".concat(e.replace(/(^\w+:|^)\/\//,""),":").concat(n,"/live")}(b,y));var S=[{icon:(0,i.jsx)(j.Z,{twoToneColor:"#6f42c1"}),title:"Use your broadcasting software",content:(0,i.jsxs)("div",{children:[(0,i.jsx)("a",{href:"https://owncast.online/docs/broadcasting/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn how to point your existing software to your new server and start streaming your content."}),(0,i.jsxs)("div",{className:"stream-info-container",children:[(0,i.jsx)(B,{strong:!0,className:"stream-info-label",children:"Streaming URL:"}),o&&(0,i.jsx)(z,{className:"stream-info-box",copyable:!0,children:o}),(0,i.jsx)(B,{strong:!0,className:"stream-info-label",children:"Stream Key:"}),(0,i.jsx)(z,{className:"stream-info-box",copyable:{text:x},children:"*********************"})]})]})},{icon:(0,i.jsx)(g.Z,{twoToneColor:"#f9826c"}),title:"Embed your video onto other sites",content:(0,i.jsx)("div",{children:(0,i.jsx)("a",{href:"https://owncast.online/docs/embed?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn how you can add your Owncast stream to other sites you control."})})}];return(null===m||void 0===m?void 0:m.chatDisabled)||S.push({icon:(0,i.jsx)(w.Z,{twoToneColor:"#0366d6"}),title:"Chat is disabled",content:(0,i.jsx)("span",{children:"Chat will continue to be disabled until you begin a live stream."})}),(null===m||void 0===m||null===(r=m.yp)||void 0===r?void 0:r.enabled)||S.push({icon:(0,i.jsx)(Z.Z,{twoToneColor:"#D18BFE"}),title:"Find an audience on the Owncast Directory",content:(0,i.jsxs)("div",{children:["List yourself in the Owncast Directory and show off your stream. Enable it in"," ",(0,i.jsx)(k(),{href:"/config-public-details",children:"settings."})]})}),(null===m||void 0===m||null===(s=m.federation)||void 0===s?void 0:s.enabled)||S.push({icon:(0,i.jsx)("img",{alt:"fediverse",width:"20px",src:"fediverse-white.png"}),title:"Add your Owncast instance to the Fediverse",content:(0,i.jsxs)("div",{children:[(0,i.jsx)(k(),{href:"/config-federation",children:"Enable Owncast social"})," features to have your instance join the Fediverse, allowing people to follow, share and engage with your live stream."]})}),(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(d.Z,{children:(0,i.jsx)(u.Z,{span:12,offset:6,children:(0,i.jsxs)("div",{className:"offline-intro",children:[(0,i.jsx)("span",{className:"logo",children:(0,i.jsx)(N.Z,{variant:"simple"})}),(0,i.jsxs)("div",{children:[(0,i.jsx)(L,{level:2,children:"No stream is active"}),(0,i.jsx)("p",{children:"You should start one."})]})]})})}),(0,i.jsxs)(d.Z,{gutter:[16,16],className:"offline-content",children:[(0,i.jsx)(u.Z,{span:12,xs:24,sm:24,md:24,lg:12,className:"list-section",children:S.map((function(e){return(0,i.jsx)(c.Z,{size:"small",bordered:!1,children:(0,i.jsx)(F,{avatar:e.icon,title:e.title,description:e.content})},e.title)}))}),(0,i.jsx)(u.Z,{span:12,xs:24,sm:24,md:24,lg:12,children:(0,i.jsx)(I,{})})]}),(0,i.jsx)(p.Z,{logs:h,pageSize:5})]})}var M=t(35660),A=t(70241);function Q(e,n,t,r,s,i,a){try{var o=e[i](a),l=o.value}catch(c){return void t(c)}o.done?n(l):Promise.resolve(l).then(r,s)}function V(e){return(0,i.jsxs)("ul",{className:"statistics-list",children:[(0,i.jsxs)("li",{children:[e.videoCodec||"Unknown"," @ ",e.videoBitrate||"Unknown"," kbps"]}),(0,i.jsxs)("li",{children:[e.framerate||"Unknown"," fps"]}),(0,i.jsxs)("li",{children:[e.width," x ",e.height]})]})}function K(){var e,n,t=(0,a.useContext)(v.aC),r=t||{},j=r.broadcaster,g=r.serverConfig,w=j||{},Z=w.remoteAddr,y=w.streamDetails,b=(null===y||void 0===y?void 0:y.encoder)||"Unknown encoder",k=(0,a.useState)([]),N=k[0],S=k[1],O=function(){var e,n=(e=s().mark((function e(){var n;return s().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,(0,C.rQ)(C.WQ);case 3:n=e.sent,S(n),e.next=10;break;case 7:e.prev=7,e.t0=e.catch(0),console.log("==== error",e.t0);case 10:case"end":return e.stop()}}),e,null,[[0,7]])})),function(){var n=this,t=arguments;return new Promise((function(r,s){var i=e.apply(n,t);function a(e){Q(i,r,s,a,o,"next",e)}function o(e){Q(i,r,s,a,o,"throw",e)}a(void 0)}))});return function(){return n.apply(this,arguments)}}(),P=function(){O()};if((0,a.useEffect)((function(){P();var e;return e=setInterval(P,C.NE),function(){clearInterval(e)}}),[]),(0,A.Qr)(g)||(0,A.Qr)(t))return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(o.Z,{active:!0}),(0,i.jsx)(o.Z,{active:!0}),(0,i.jsx)(o.Z,{active:!0})]});if(!j)return(0,i.jsx)(H,{logs:N,config:g});var T=null===t||void 0===t||null===(e=t.currentBroadcast)||void 0===e||null===(n=e.outputSettings)||void 0===n?void 0:n.map((function(e){var n=e.audioPassthrough,t=e.videoPassthrough,r=e.audioBitrate,s=e.videoBitrate,a=e.framerate,o=n?"".concat(y.audioCodec||"Unknown",", ").concat(y.audioBitrate," kbps"):"".concat(r||"Unknown"," kbps"),c=t?"".concat(y.videoBitrate||"Unknown"," kbps, ").concat(y.framerate," fps ").concat(y.width," x ").concat(y.height):"".concat(s||"Unknown"," kbps, ").concat(a," fps");return(0,i.jsxs)("div",{className:"stream-details-item-container",children:[(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Outbound Video Stream",value:c}),(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Outbound Audio Stream",value:o})]})})),_=t.viewerCount,E=t.sessionPeakViewerCount,D="".concat(y.audioCodec,", ").concat(y.audioBitrate||"Unknown"," kbps"),U=new Date(j.time);return(0,i.jsxs)("div",{className:"home-container",children:[(0,i.jsxs)("div",{className:"sections-container",children:[(0,i.jsx)("div",{className:"online-status-section",children:(0,i.jsxs)(c.Z,{size:"small",type:"inner",className:"online-details-card",children:[(0,i.jsxs)(d.Z,{gutter:[16,16],align:"middle",children:[(0,i.jsx)(u.Z,{span:8,sm:24,md:8,children:(0,i.jsx)(l.Z,{title:"Stream started ".concat((0,f.Z)(U,Date.now())),value:(0,x.Z)(U),prefix:(0,i.jsx)(h.Z,{})})}),(0,i.jsx)(u.Z,{span:8,sm:24,md:8,children:(0,i.jsx)(l.Z,{title:"Viewers",value:_,prefix:(0,i.jsx)(m.Z,{})})}),(0,i.jsx)(u.Z,{span:8,sm:24,md:8,children:(0,i.jsx)(l.Z,{title:"Peak viewer count",value:E,prefix:(0,i.jsx)(m.Z,{})})})]}),(0,i.jsx)(M.Z,{})]})}),(0,i.jsxs)(d.Z,{gutter:[16,16],className:"section stream-details-section",children:[(0,i.jsxs)(u.Z,{className:"stream-details",span:12,sm:24,md:24,lg:12,children:[(0,i.jsx)(c.Z,{size:"small",title:"Outbound Stream Details",type:"inner",className:"outbound-details",children:T}),(0,i.jsxs)(c.Z,{size:"small",title:"Inbound Stream Details",type:"inner",children:[(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Input",value:"".concat(b," ").concat((0,A.t5)(Z))}),(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Inbound Video Stream",value:y,formatter:V}),(0,i.jsx)(l.Z,{className:"stream-details-item",title:"Inbound Audio Stream",value:D})]})]}),(0,i.jsx)(u.Z,{span:12,xs:24,sm:24,md:24,lg:12,children:(0,i.jsx)(I,{})})]})]}),(0,i.jsx)("br",{}),(0,i.jsx)(p.Z,{logs:N,pageSize:5})]})}}},function(e){e.O(0,[3698,1213,8147,7961,5938,7751,4763,2429,3016,9774,2888,179],(function(){return n=58860,e(e.s=n);var n}));var n=e.O();_N_E=n}]);