(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[177],{94194:function(e,t,r){(window.__NEXT_P=window.__NEXT_P||[]).push(["/hardware-info",function(){return r(61003)}])},89270:function(e,t,r){"use strict";r.d(t,{Z:function(){return s}});var n=r(85893),o=r(31877),c=r(92616),a=r.n(c),l=r(58091),i=r(60727);function u(e){var t={};return e.forEach((function(e){var r=new Date(e.time),n=(0,l.Z)(r,"H:mma");t[n]=e.value})),t}function s(e){var t=e.data,r=e.title,o=e.color,c=e.unit,a=e.dataCollections,l=[];return t&&t.length>0&&l.push({name:r,color:o,data:u(t)}),a.forEach((function(e){l.push({name:e.name,data:u(e.data),color:e.color})})),(0,n.jsx)("div",{className:"line-chart-container",children:(0,n.jsx)(i.wW,{xtitle:"Time",ytitle:r,suffix:c,legend:"bottom",color:o,data:l,download:r})})}a().use(o.Z),s.defaultProps={dataCollections:[],data:[],title:""}},34440:function(e,t,r){"use strict";r.d(t,{Z:function(){return p}});var n=r(85893),o=r(17256),c=r(58008),a=r(74763),l=r(97751);function i(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function u(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{},n=Object.keys(r);"function"===typeof Object.getOwnPropertySymbols&&(n=n.concat(Object.getOwnPropertySymbols(r).filter((function(e){return Object.getOwnPropertyDescriptor(r,e).enumerable})))),n.forEach((function(t){i(e,t,r[t])}))}return e}var s=o.Z.Text,f={title:"",value:0,prefix:null,color:"",progress:!1,centered:!1,formatter:null};function d(e){var t=e.title,r=e.value,o=e.prefix,a=e.color,l=r>90?"red":a,i=(0,n.jsxs)("div",{children:[o,(0,n.jsx)("div",{children:(0,n.jsx)(s,{type:"secondary",children:t})}),(0,n.jsx)("div",{children:(0,n.jsxs)(s,{type:"secondary",children:[r,"%"]})})]});return(0,n.jsx)(c.Z,{type:"dashboard",percent:r,width:120,strokeColor:{"0%":a,"90%":l},format:function(){return i}})}function v(e){var t=e.title,r=e.value,o=e.prefix,c=e.formatter;return(0,n.jsx)(a.Z,{title:t,value:r,prefix:o,formatter:c})}function p(e){var t=e.progress?d:v,r=e.centered?{display:"flex",alignItems:"center",justifyContent:"center"}:{};return(0,n.jsx)(l.Z,{type:"inner",children:(0,n.jsx)("div",{style:r,children:(0,n.jsx)(t,u({},e))})})}d.defaultProps=f,v.defaultProps=f,p.defaultProps=f},61003:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return k}});var n=r(28520),o=r.n(n),c=r(85893),a=r(1413),l=r(67294),i={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M956.9 845.1L896.4 632V168c0-17.7-14.3-32-32-32h-704c-17.7 0-32 14.3-32 32v464L67.9 845.1C60.4 866 75.8 888 98 888h828.8c22.2 0 37.6-22 30.1-42.9zM200.4 208h624v395h-624V208zm228.3 608l8.1-37h150.3l8.1 37H428.7zm224 0l-19.1-86.7c-.8-3.7-4.1-6.3-7.8-6.3H398.2c-3.8 0-7 2.6-7.8 6.3L371.3 816H151l42.3-149h638.2l42.3 149H652.7z"}}]},name:"laptop",theme:"outlined"},u=r(42135),s=function(e,t){return l.createElement(u.Z,(0,a.Z)((0,a.Z)({},e),{},{ref:t,icon:i}))};s.displayName="LaptopOutlined";var f=l.forwardRef(s),d={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M632 888H392c-4.4 0-8 3.6-8 8v32c0 17.7 14.3 32 32 32h192c17.7 0 32-14.3 32-32v-32c0-4.4-3.6-8-8-8zM512 64c-181.1 0-328 146.9-328 328 0 121.4 66 227.4 164 284.1V792c0 17.7 14.3 32 32 32h264c17.7 0 32-14.3 32-32V676.1c98-56.7 164-162.7 164-284.1 0-181.1-146.9-328-328-328zm127.9 549.8L604 634.6V752H420V634.6l-35.9-20.8C305.4 568.3 256 484.5 256 392c0-141.4 114.6-256 256-256s256 114.6 256 256c0 92.5-49.4 176.3-128.1 221.8z"}}]},name:"bulb",theme:"outlined"},v=function(e,t){return l.createElement(u.Z,(0,a.Z)((0,a.Z)({},e),{},{ref:t,icon:d}))};v.displayName="BulbOutlined";var p=l.forwardRef(v),h={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M893.3 293.3L730.7 130.7c-7.5-7.5-16.7-13-26.7-16V112H144c-17.7 0-32 14.3-32 32v736c0 17.7 14.3 32 32 32h736c17.7 0 32-14.3 32-32V338.5c0-17-6.7-33.2-18.7-45.2zM384 184h256v104H384V184zm456 656H184V184h136v136c0 17.7 14.3 32 32 32h320c17.7 0 32-14.3 32-32V205.8l136 136V840zM512 442c-79.5 0-144 64.5-144 144s64.5 144 144 144 144-64.5 144-144-64.5-144-144-144zm0 224c-44.2 0-80-35.8-80-80s35.8-80 80-80 80 35.8 80 80-35.8 80-80 80z"}}]},name:"save",theme:"outlined"},m=function(e,t){return l.createElement(u.Z,(0,a.Z)((0,a.Z)({},e),{},{ref:t,icon:h}))};m.displayName="SaveOutlined";var x=l.forwardRef(m),y=r(17256),j=r(25968),g=r(6226),b=r(58827),w=r(89270),Z=r(34440);function O(e,t,r,n,o,c,a){try{var l=e[c](a),i=l.value}catch(u){return void r(u)}l.done?t(i):Promise.resolve(i).then(n,o)}function P(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function E(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{},n=Object.keys(r);"function"===typeof Object.getOwnPropertySymbols&&(n=n.concat(Object.getOwnPropertySymbols(r).filter((function(e){return Object.getOwnPropertyDescriptor(r,e).enumerable})))),n.forEach((function(t){P(e,t,r[t])}))}return e}function k(){var e,t,r,n,a=(0,l.useState)({cpu:[],memory:[],disk:[],message:""}),i=a[0],u=a[1],s=(n=o().mark((function e(){var t;return o().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,(0,b.rQ)(b.nx);case 3:t=e.sent,u(E({},t)),e.next=10;break;case 7:e.prev=7,e.t0=e.catch(0),u(E({},i,{message:e.t0.message}));case 10:case"end":return e.stop()}}),e,null,[[0,7]])})),function(){var e=this,t=arguments;return new Promise((function(r,o){var c=n.apply(e,t);function a(e){O(c,r,o,a,l,"next",e)}function l(e){O(c,r,o,a,l,"throw",e)}a(void 0)}))});if((0,l.useEffect)((function(){var e;return s(),e=setInterval(s,b.NE),function(){clearInterval(e)}}),[]),!i.cpu)return null;var d=null===(e=i.cpu[i.cpu.length-1])||void 0===e?void 0:e.value,v=null===(t=i.memory[i.memory.length-1])||void 0===t?void 0:t.value,h=null===(r=i.disk[i.disk.length-1])||void 0===r?void 0:r.value,m=[{name:"CPU",color:"#B63FFF",data:i.cpu},{name:"Memory",color:"#2087E2",data:i.memory},{name:"Disk",color:"#FF7700",data:i.disk}];return(0,c.jsxs)(c.Fragment,{children:[(0,c.jsx)(y.Z.Title,{children:"Hardware Info"}),(0,c.jsx)("br",{}),(0,c.jsxs)("div",{children:[(0,c.jsxs)(j.Z,{gutter:[16,16],justify:"space-around",children:[(0,c.jsx)(g.Z,{children:(0,c.jsx)(Z.Z,{title:m[0].name,value:"".concat(d||0),prefix:(0,c.jsx)(f,{style:{color:m[0].color}}),color:m[0].color,progress:!0,centered:!0})}),(0,c.jsx)(g.Z,{children:(0,c.jsx)(Z.Z,{title:m[1].name,value:"".concat(v||0),prefix:(0,c.jsx)(p,{style:{color:m[1].color}}),color:m[1].color,progress:!0,centered:!0})}),(0,c.jsx)(g.Z,{children:(0,c.jsx)(Z.Z,{title:m[2].name,value:"".concat(h||0),prefix:(0,c.jsx)(x,{style:{color:m[2].color}}),color:m[2].color,progress:!0,centered:!0})})]}),(0,c.jsx)(w.Z,{title:"% used",dataCollections:m,color:"#FF7700",unit:"%"})]})]})}}},function(e){e.O(0,[570,91,903,138,42,958,774,888,179],(function(){return t=94194,e(e.s=t);var t}));var t=e.O();_N_E=t}]);