(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[194],{87547:function(e,t,r){"use strict";r.d(t,{Z:function(){return s}});var n=r(1413),i=r(67294),o={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M858.5 763.6a374 374 0 00-80.6-119.5 375.63 375.63 0 00-119.5-80.6c-.4-.2-.8-.3-1.2-.5C719.5 518 760 444.7 760 362c0-137-111-248-248-248S264 225 264 362c0 82.7 40.5 156 102.8 201.1-.4.2-.8.3-1.2.5-44.8 18.9-85 46-119.5 80.6a375.63 375.63 0 00-80.6 119.5A371.7 371.7 0 00136 901.8a8 8 0 008 8.2h60c4.4 0 7.9-3.5 8-7.8 2-77.2 33-149.5 87.8-204.3 56.7-56.7 132-87.9 212.2-87.9s155.5 31.2 212.2 87.9C779 752.7 810 825 812 902.2c.1 4.4 3.6 7.8 8 7.8h60a8 8 0 008-8.2c-1-47.8-10.9-94.3-29.5-138.2zM512 534c-45.9 0-89.1-17.9-121.6-50.4S340 407.9 340 362c0-45.9 17.9-89.1 50.4-121.6S466.1 190 512 190s89.1 17.9 121.6 50.4S684 316.1 684 362c0 45.9-17.9 89.1-50.4 121.6S557.9 534 512 534z"}}]},name:"user",theme:"outlined"},a=r(42135),c=function(e,t){return i.createElement(a.Z,(0,n.Z)((0,n.Z)({},e),{},{ref:t,icon:o}))};c.displayName="UserOutlined";var s=i.forwardRef(c)},31709:function(e,t,r){(window.__NEXT_P=window.__NEXT_P||[]).push(["/viewer-info",function(){return r(48935)}])},89270:function(e,t,r){"use strict";r.d(t,{Z:function(){return l}});var n=r(85893),i=r(31877),o=r(92616),a=r.n(o),c=r(58091),s=r(60727);function u(e){var t={};return e.forEach((function(e){var r=new Date(e.time),n=(0,c.Z)(r,"H:mma");t[n]=e.value})),t}function l(e){var t=e.data,r=e.title,i=e.color,o=e.unit,a=e.dataCollections,c=[];return t&&t.length>0&&c.push({name:r,color:i,data:u(t)}),a.forEach((function(e){c.push({name:e.name,data:u(e.data),color:e.color})})),(0,n.jsx)("div",{className:"line-chart-container",children:(0,n.jsx)(s.wW,{xtitle:"Time",ytitle:r,suffix:o,legend:"bottom",color:i,data:c,download:r})})}a().use(i.Z),l.defaultProps={dataCollections:[],data:[],title:""}},34440:function(e,t,r){"use strict";r.d(t,{Z:function(){return x}});var n=r(85893),i=r(17256),o=r(58008),a=r(74763),c=r(97751);function s(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function u(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{},n=Object.keys(r);"function"===typeof Object.getOwnPropertySymbols&&(n=n.concat(Object.getOwnPropertySymbols(r).filter((function(e){return Object.getOwnPropertyDescriptor(r,e).enumerable})))),n.forEach((function(t){s(e,t,r[t])}))}return e}var l=i.Z.Text,f={title:"",value:0,prefix:null,color:"",progress:!1,centered:!1,formatter:null};function d(e){var t=e.title,r=e.value,i=e.prefix,a=e.color,c=r>90?"red":a,s=(0,n.jsxs)("div",{children:[i,(0,n.jsx)("div",{children:(0,n.jsx)(l,{type:"secondary",children:t})}),(0,n.jsx)("div",{children:(0,n.jsxs)(l,{type:"secondary",children:[r,"%"]})})]});return(0,n.jsx)(o.Z,{type:"dashboard",percent:r,width:120,strokeColor:{"0%":a,"90%":c},format:function(){return s}})}function v(e){var t=e.title,r=e.value,i=e.prefix,o=e.formatter;return(0,n.jsx)(a.Z,{title:t,value:r,prefix:i,formatter:o})}function x(e){var t=e.progress?d:v,r=e.centered?{display:"flex",alignItems:"center",justifyContent:"center"}:{};return(0,n.jsx)(c.Z,{type:"inner",children:(0,n.jsx)("div",{style:r,children:(0,n.jsx)(t,u({},e))})})}d.defaultProps=f,v.defaultProps=f,x.defaultProps=f},48935:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return h}});var n=r(28520),i=r.n(n),o=r(85893),a=r(67294),c=r(17256),s=r(25968),u=r(6226),l=r(87547),f=r(89270),d=r(34440),v=r(35159),x=r(58827);function p(e,t,r,n,i,o,a){try{var c=e[o](a),s=c.value}catch(u){return void r(u)}c.done?t(s):Promise.resolve(s).then(n,i)}function h(){var e,t=(0,a.useContext)(v.aC)||{},r=t.online,n=t.viewerCount,h=t.overallPeakViewerCount,j=t.sessionPeakViewerCount,m=(0,a.useState)([]),w=m[0],Z=m[1],y=(e=i().mark((function e(){var t;return i().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,(0,x.rQ)(x.iV);case 3:t=e.sent,Z(t),e.next=10;break;case 7:e.prev=7,e.t0=e.catch(0),console.log("==== error",e.t0);case 10:case"end":return e.stop()}}),e,null,[[0,7]])})),function(){var t=this,r=arguments;return new Promise((function(n,i){var o=e.apply(t,r);function a(e){p(o,n,i,a,c,"next",e)}function c(e){p(o,n,i,a,c,"throw",e)}a(void 0)}))});return(0,a.useEffect)((function(){var e=null;return y(),r?(e=setInterval(y,6e4),function(){clearInterval(e)}):function(){return[]}}),[r]),w.length?(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)(c.Z.Title,{children:"Viewer Info"}),(0,o.jsx)("br",{}),(0,o.jsxs)(s.Z,{gutter:[16,16],justify:"space-around",children:[r&&(0,o.jsx)(u.Z,{span:8,md:8,children:(0,o.jsx)(d.Z,{title:"Current viewers",value:n.toString(),prefix:(0,o.jsx)(l.Z,{})})}),(0,o.jsx)(u.Z,{md:r?8:12,children:(0,o.jsx)(d.Z,{title:r?"Max viewers this session":"Max viewers last session",value:j.toString(),prefix:(0,o.jsx)(l.Z,{})})}),(0,o.jsx)(u.Z,{md:r?8:12,children:(0,o.jsx)(d.Z,{title:"All-time max viewers",value:h.toString(),prefix:(0,o.jsx)(l.Z,{})})})]}),(0,o.jsx)(f.Z,{title:"Viewers",data:w,color:"#2087E2",unit:""})]}):"no info"}}},function(e){e.O(0,[570,91,903,138,42,958,774,888,179],(function(){return t=31709,e(e.s=t);var t}));var t=e.O();_N_E=t}]);