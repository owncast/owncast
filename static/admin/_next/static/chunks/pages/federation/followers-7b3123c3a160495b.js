(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[480],{1481:function(e,t,r){(window.__NEXT_P=window.__NEXT_P||[]).push(["/federation/followers",function(){return r(24890)}])},24890:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return K}});var n=r(28520),a=r.n(n),i=r(85893),c=r(67294),o=r(87462),s=r(4942),l=r(71002),u=r(97685),f=r(94184),d=r.n(f),p=r(4084),m=r(42550),h=r(59844),v=r(21687),g=r(24308),y=r(25378),w=c.createContext("default"),x=function(e){var t=e.children,r=e.size;return c.createElement(w.Consumer,null,(function(e){return c.createElement(w.Provider,{value:r||e},t)}))},k=w,E=function(e,t){var r={};for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&t.indexOf(n)<0&&(r[n]=e[n]);if(null!=e&&"function"===typeof Object.getOwnPropertySymbols){var a=0;for(n=Object.getOwnPropertySymbols(e);a<n.length;a++)t.indexOf(n[a])<0&&Object.prototype.propertyIsEnumerable.call(e,n[a])&&(r[n[a]]=e[n[a]])}return r},b=function(e,t){var r,n,a=c.useContext(k),i=c.useState(1),f=(0,u.Z)(i,2),w=f[0],x=f[1],b=c.useState(!1),Z=(0,u.Z)(b,2),j=Z[0],S=Z[1],O=c.useState(!0),N=(0,u.Z)(O,2),C=N[0],P=N[1],z=c.useRef(),_=c.useRef(),A=(0,m.sQ)(t,z),T=c.useContext(h.E_).getPrefixCls,I=function(){if(_.current&&z.current){var t=_.current.offsetWidth,r=z.current.offsetWidth;if(0!==t&&0!==r){var n=e.gap,a=void 0===n?4:n;2*a<r&&x(r-2*a<t?(r-2*a)/t:1)}}};c.useEffect((function(){S(!0)}),[]),c.useEffect((function(){P(!0),x(1)}),[e.src]),c.useEffect((function(){I()}),[e.gap]);var F=e.prefixCls,D=e.shape,R=e.size,Q=e.src,H=e.srcSet,M=e.icon,q=e.className,W=e.alt,X=e.draggable,B=e.children,G=e.crossOrigin,K=E(e,["prefixCls","shape","size","src","srcSet","icon","className","alt","draggable","children","crossOrigin"]),U="default"===R?a:R,V=(0,y.Z)(),J=c.useMemo((function(){if("object"!==(0,l.Z)(U))return{};var e=g.c4.find((function(e){return V[e]})),t=U[e];return t?{width:t,height:t,lineHeight:"".concat(t,"px"),fontSize:M?t/2:18}:{}}),[V,U]);(0,v.Z)(!("string"===typeof M&&M.length>2),"Avatar","`icon` is using ReactNode instead of string naming in v4. Please check `".concat(M,"` at https://ant.design/components/icon"));var L,Y=T("avatar",F),$=d()((r={},(0,s.Z)(r,"".concat(Y,"-lg"),"large"===U),(0,s.Z)(r,"".concat(Y,"-sm"),"small"===U),r)),ee=c.isValidElement(Q),te=d()(Y,$,(n={},(0,s.Z)(n,"".concat(Y,"-").concat(D),!!D),(0,s.Z)(n,"".concat(Y,"-image"),ee||Q&&C),(0,s.Z)(n,"".concat(Y,"-icon"),!!M),n),q),re="number"===typeof U?{width:U,height:U,lineHeight:"".concat(U,"px"),fontSize:M?U/2:18}:{};if("string"===typeof Q&&C)L=c.createElement("img",{src:Q,draggable:X,srcSet:H,onError:function(){var t=e.onError;!1!==(t?t():void 0)&&P(!1)},alt:W,crossOrigin:G});else if(ee)L=Q;else if(M)L=M;else if(j||1!==w){var ne="scale(".concat(w,") translateX(-50%)"),ae={msTransform:ne,WebkitTransform:ne,transform:ne},ie="number"===typeof U?{lineHeight:"".concat(U,"px")}:{};L=c.createElement(p.default,{onResize:I},c.createElement("span",{className:"".concat(Y,"-string"),ref:function(e){_.current=e},style:(0,o.Z)((0,o.Z)({},ie),ae)},B))}else L=c.createElement("span",{className:"".concat(Y,"-string"),style:{opacity:0},ref:function(e){_.current=e}},B);return delete K.onError,delete K.gap,c.createElement("span",(0,o.Z)({},K,{style:(0,o.Z)((0,o.Z)((0,o.Z)({},re),J),K.style),className:te,ref:A}),L)},Z=c.forwardRef(b);Z.displayName="Avatar",Z.defaultProps={shape:"circle",size:"default"};var j=Z,S=r(50344),O=r(96159),N=r(55241),C=function(e){var t=c.useContext(h.E_),r=t.getPrefixCls,n=t.direction,a=e.prefixCls,i=e.className,o=void 0===i?"":i,l=e.maxCount,u=e.maxStyle,f=e.size,p=r("avatar-group",a),m=d()(p,(0,s.Z)({},"".concat(p,"-rtl"),"rtl"===n),o),v=e.children,g=e.maxPopoverPlacement,y=void 0===g?"top":g,w=(0,S.Z)(v).map((function(e,t){return(0,O.Tm)(e,{key:"avatar-key-".concat(t)})})),k=w.length;if(l&&l<k){var E=w.slice(0,l),b=w.slice(l,k);return E.push(c.createElement(N.Z,{key:"avatar-popover-key",content:b,trigger:"hover",placement:y,overlayClassName:"".concat(p,"-popover")},c.createElement(j,{style:u},"+".concat(k-l)))),c.createElement(x,{size:f},c.createElement("div",{className:m,style:e.style},E))}return c.createElement(x,{size:f},c.createElement("div",{className:m,style:e.style},w))},P=j;P.Group=C;var z=P,_=r(17256),A=r(73555),T=r(71577),I=r(58091),F=r(1413),D={icon:{tag:"svg",attrs:{viewBox:"64 64 896 896",focusable:"false"},children:[{tag:"path",attrs:{d:"M678.3 642.4c24.2-13 51.9-20.4 81.4-20.4h.1c3 0 4.4-3.6 2.2-5.6a371.67 371.67 0 00-103.7-65.8c-.4-.2-.8-.3-1.2-.5C719.2 505 759.6 431.7 759.6 349c0-137-110.8-248-247.5-248S264.7 212 264.7 349c0 82.7 40.4 156 102.6 201.1-.4.2-.8.3-1.2.5-44.7 18.9-84.8 46-119.3 80.6a373.42 373.42 0 00-80.4 119.5A373.6 373.6 0 00137 888.8a8 8 0 008 8.2h59.9c4.3 0 7.9-3.5 8-7.8 2-77.2 32.9-149.5 87.6-204.3C357 628.2 432.2 597 512.2 597c56.7 0 111.1 15.7 158 45.1a8.1 8.1 0 008.1.3zM512.2 521c-45.8 0-88.9-17.9-121.4-50.4A171.2 171.2 0 01340.5 349c0-45.9 17.9-89.1 50.3-121.6S466.3 177 512.2 177s88.9 17.9 121.4 50.4A171.2 171.2 0 01683.9 349c0 45.9-17.9 89.1-50.3 121.6C601.1 503.1 558 521 512.2 521zM880 759h-84v-84c0-4.4-3.6-8-8-8h-56c-4.4 0-8 3.6-8 8v84h-84c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h84v84c0 4.4 3.6 8 8 8h56c4.4 0 8-3.6 8-8v-84h84c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8z"}}]},name:"user-add",theme:"outlined"},R=r(42135),Q=function(e,t){return c.createElement(R.Z,(0,F.Z)((0,F.Z)({},e),{},{ref:t,icon:D}))};Q.displayName="UserAddOutlined";var H=c.forwardRef(Q),M=r(58827),q=r(2766);function W(e,t,r,n,a,i,c){try{var o=e[i](c),s=o.value}catch(l){return void r(l)}o.done?t(s):Promise.resolve(s).then(n,a)}function X(e){return function(){var t=this,r=arguments;return new Promise((function(n,a){var i=e.apply(t,r);function c(e){W(i,n,a,c,o,"next",e)}function o(e){W(i,n,a,c,o,"throw",e)}c(void 0)}))}}function B(e){return function(e){if(Array.isArray(e)){for(var t=0,r=new Array(e.length);t<e.length;t++)r[t]=e[t];return r}}(e)||function(e){if(Symbol.iterator in Object(e)||"[object Arguments]"===Object.prototype.toString.call(e))return Array.from(e)}(e)||function(){throw new TypeError("Invalid attempt to spread non-iterable instance")}()}var G=_.Z.Title;function K(){var e=function(e,t){return(0,i.jsx)(A.Z,{dataSource:e,columns:t,size:"small",rowKey:function(e){return e.link},pagination:{pageSize:20}})},t=(0,c.useState)([]),r=t[0],n=t[1],o=(0,c.useState)([]),s=o[0],l=o[1],u=X(a().mark((function e(){var t,r;return a().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,(0,M.rQ)(M.HP,{auth:!0});case 3:return t=e.sent,(0,q.Qr)(t)?l([]):l(t),e.next=7,(0,M.rQ)(M.E8,{auth:!0});case 7:r=e.sent,(0,q.Qr)(r)?n([]):n(r),e.next=14;break;case 11:e.prev=11,e.t0=e.catch(0),console.log("==== error",e.t0);case 14:case"end":return e.stop()}}),e,null,[[0,11]])})));(0,c.useEffect)((function(){u()}),[]);var f=[{title:"",dataIndex:"image",key:"image",width:90,render:function(e){return(0,i.jsx)(z,{size:40,src:e})}},{title:"Name",dataIndex:"name",key:"name"},{title:"Account",dataIndex:"username",key:"username",render:function(e,t){return(0,i.jsx)("a",{href:t.link,target:"_blank",rel:"noreferrer",children:t.link})}}];function d(){return(d=X(a().mark((function e(t){return a().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return console.log("approve",t),e.prev=1,e.next=4,(0,M.rQ)(M.l9,{auth:!0,method:"POST",data:{federationIRI:t.link,approved:!0}});case 4:u(),e.next=10;break;case 7:e.prev=7,e.t0=e.catch(1),console.error(e.t0);case 10:case"end":return e.stop()}}),e,null,[[1,7]])})))).apply(this,arguments)}var p=B(f);p.unshift({title:"Approve",dataIndex:null,key:null,render:function(e){return(0,i.jsx)(T.Z,{type:"primary",icon:(0,i.jsx)(H,{}),size:"large",onClick:function(){!function(e){d.apply(this,arguments)}(e)}})},width:50}),p.push({title:"Requested",dataIndex:"timestamp",key:"requested",width:200,render:function(e){var t=new Date(e);return(0,i.jsx)(i.Fragment,{children:(0,I.Z)(t,"P")})},sorter:function(e,t){return new Date(e.timestamp).getTime()-new Date(t.timestamp).getTime()},sortDirections:["descend","ascend"],defaultSortOrder:"descend"});var m=r.length>0&&(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(G,{children:"Followers needing Approval"}),(0,i.jsxs)("p",{children:["The following people are requesting to follow your Owncast server on the"," ",(0,i.jsx)("a",{href:"https://en.wikipedia.org/wiki/Fediverse",target:"_blank",rel:"noopener noreferrer",children:"Fediverse"})," ","and be alerted to when you go live."]}),e(r,p)]}),h=B(f);return h.push({title:"Added",dataIndex:"timestamp",key:"timestamp",width:200,render:function(e){var t=new Date(e);return(0,i.jsx)(i.Fragment,{children:(0,I.Z)(t,"P")})},sorter:function(e,t){return new Date(e.timestamp).getTime()-new Date(t.timestamp).getTime()},sortDirections:["descend","ascend"],defaultSortOrder:"descend"}),(0,i.jsxs)("div",{className:"followers-section",children:[m,(0,i.jsx)(G,{children:"Followers"}),(0,i.jsxs)("p",{children:["The following accounts get notified on the"," ",(0,i.jsx)("a",{href:"https://en.wikipedia.org/wiki/Fediverse",target:"_blank",rel:"noopener noreferrer",children:"Fediverse"})," ","when you go live."]}),e(s,h)]})}}},function(e){e.O(0,[555,91,774,888,179],(function(){return t=1481,e(e.s=t);var t}));var t=e.O();_N_E=t}]);