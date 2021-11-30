(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[457],{50291:function(e,n,i){(window.__NEXT_P=window.__NEXT_P||[]).push(["/config-federation",function(){return i(6898)}])},10063:function(e,n,i){"use strict";i.d(n,{Q:function(){return c},Z:function(){return d}});var r=i(85893),a=i(67294),t=i(44068),s=i(20550),l=i(48419),o=t.Z.Title,c="#5a67d8";function d(e){var n=(0,a.useState)(""),i=n[0],t=n[1],d=e.title,u=e.description,h=e.placeholder,f=e.maxLength,p=e.values,v=e.handleDeleteIndex,m=e.handleCreateString,g=e.submitStatus;return(0,r.jsxs)("div",{className:"edit-string-array-container",children:[(0,r.jsx)(o,{level:3,className:"section-title",children:d}),(0,r.jsx)("p",{className:"description",children:u}),(0,r.jsx)("div",{className:"edit-current-strings",children:null===p||void 0===p?void 0:p.map((function(e,n){return(0,r.jsx)(s.Z,{closable:!0,onClose:function(){v(n)},color:c,children:e},"tag-".concat(e,"-").concat(n))}))}),(0,r.jsx)("div",{className:"add-new-string-section",children:(0,r.jsx)(l.ZP,{fieldName:"string-input",value:i,onChange:function(e){var n=e.value;t(n)},onPressEnter:function(){var e=i.trim();m(e),t("")},maxLength:f,placeholder:h,status:g})})]})}d.defaultProps={maxLength:50,description:null,submitStatus:null}},15976:function(e,n,i){"use strict";i.d(n,{Z:function(){return f}});var r=i(28520),a=i.n(r),t=i(85893),s=i(67294),l=i(94594),o=i(83200),c=i(78464),d=i(25964),u=i(35159);function h(e,n,i,r,a,t,s){try{var l=e[t](s),o=l.value}catch(c){return void i(c)}l.done?n(o):Promise.resolve(o).then(r,a)}function f(e){var n,i=(0,s.useState)(null),r=i[0],f=i[1],p=null,v=((0,s.useContext)(u.aC)||{}).setFieldInConfigState,m=e.apiPath,g=e.checked,b=e.reversed,x=void 0!==b&&b,j=e.configPath,w=void 0===j?"":j,k=e.disabled,y=void 0!==k&&k,P=e.fieldName,Z=e.label,N=e.tip,C=e.useSubmit,S=e.onChange,E=function(){f(null),clearTimeout(p),p=null},F=(n=a().mark((function e(n){var i;return a().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!C){e.next=6;break}return f((0,o.kg)(o.Jk)),i=x?!n:n,e.next=5,(0,d.Si)({apiPath:m,data:{value:i},onSuccess:function(){v({fieldName:P,value:i,path:w}),f((0,o.kg)(o.zv))},onError:function(e){f((0,o.kg)(o.Un,"There was an error: ".concat(e)))}});case 5:p=setTimeout(E,d.sI);case 6:S&&S(n);case 7:case"end":return e.stop()}}),e)})),function(){var e=this,i=arguments;return new Promise((function(r,a){var t=n.apply(e,i);function s(e){h(t,r,a,s,l,"next",e)}function l(e){h(t,r,a,s,l,"throw",e)}s(void 0)}))}),_=null!==r&&r.type===o.Jk;return(0,t.jsxs)("div",{className:"formfield-container toggleswitch-container",children:[Z&&(0,t.jsx)("div",{className:"label-side",children:(0,t.jsx)("span",{className:"formfield-label",children:Z})}),(0,t.jsxs)("div",{className:"input-side",children:[(0,t.jsxs)("div",{className:"input-group",children:[(0,t.jsx)(l.Z,{className:"switch field-".concat(P),loading:_,onChange:F,defaultChecked:g,checked:g,checkedChildren:"ON",unCheckedChildren:"OFF",disabled:y}),(0,t.jsx)(c.Z,{status:r})]}),(0,t.jsx)("p",{className:"field-tip",children:N})]})]})}f.defaultProps={apiPath:"",checked:!1,reversed:!1,configPath:"",disabled:!1,label:"",tip:"",useSubmit:!1,onChange:null}},6898:function(e,n,i){"use strict";i.r(n),i.d(n,{default:function(){return k}});var r=i(85893),a=i(56516),t=i(71577),s=i(44068),l=i(25968),o=i(6226),c=i(67294),d=i(45697),u=i.n(d),h=i(48419),f=i(50197),p=i(15976),v=i(10063),m=i(25964),g=i(35159),b=i(83200);function x(e,n,i){return n in e?Object.defineProperty(e,n,{value:i,enumerable:!0,configurable:!0,writable:!0}):e[n]=i,e}function j(e){for(var n=1;n<arguments.length;n++){var i=null!=arguments[n]?arguments[n]:{},r=Object.keys(i);"function"===typeof Object.getOwnPropertySymbols&&(r=r.concat(Object.getOwnPropertySymbols(i).filter((function(e){return Object.getOwnPropertyDescriptor(i,e).enumerable})))),r.forEach((function(n){x(e,n,i[n])}))}return e}function w(e){var n=e.cancelPressed,i=e.okPressed;return(0,r.jsxs)(a.Z,{width:"70%",title:"Enable Social Features",visible:!0,onCancel:n,footer:(0,r.jsxs)("div",{children:[(0,r.jsx)(t.Z,{onClick:n,children:"Do not enable"}),(0,r.jsx)(t.Z,{type:"primary",onClick:i,children:"Enable Social Features"})]}),children:[(0,r.jsx)(s.Z.Title,{level:3,children:"How do Owncast's social features work?"}),(0,r.jsxs)(s.Z.Paragraph,{children:["Owncast's social features are accomplished by having your server join the"," ",(0,r.jsx)("a",{href:"https://en.wikipedia.org/wiki/Fediverse",rel:"noopener noreferrer",target:"_blank",children:"Fediverse"}),", a decentralized, open, collection of independent servers, like yours."]}),(0,r.jsxs)(s.Z.Paragraph,{children:["Anybody with an account on any participating compatible server can follow your stream and get notified when you go live as well as share and ",(0,r.jsx)("i",{children:"Like"})," your posts."]}),(0,r.jsx)(s.Z.Title,{level:3,children:"What do you need to know?"}),(0,r.jsxs)("ul",{children:[(0,r.jsx)("li",{children:"These features are brand new. Given the variability of interfacing with the rest of the world, bugs are possible. Please report anything that you think isn't working quite right."}),(0,r.jsx)("li",{children:"You must always host your Owncast server with SSL using a https url."}),(0,r.jsx)("li",{children:"If you ever change your server name or Fediverse username you will be seen as a completely different user on the Fediverse, and the old user will disappear. It is best to never change these once you set them."}),(0,r.jsxs)("li",{children:["Turning on ",(0,r.jsx)("i",{children:"Private mode"})," will allow you to manually approve each follower and limit the visibility of your posts to followers only."]})]}),(0,r.jsx)(s.Z.Title,{level:3,children:"Learn more about The Fediverse"}),(0,r.jsx)(s.Z.Paragraph,{children:"If these concepts are new you should discover more about what The Fediverse has to offer."}),(0,r.jsxs)("div",{style:{textAlign:"center"},children:[(0,r.jsxs)(l.Z,{children:[(0,r.jsx)(o.Z,{span:6,children:(0,r.jsx)("a",{href:"https://en.wikipedia.org/wiki/Fediverse",rel:"noopener noreferrer",target:"_blank",children:"Fediverse on Wikpedia"})}),(0,r.jsx)(o.Z,{span:6,children:(0,r.jsx)("a",{href:"https://fediverse.party/",rel:"noopener noreferrer",target:"_blank",children:"Fediverse Party"})}),(0,r.jsx)(o.Z,{span:6,children:(0,r.jsx)("a",{href:"https://pleroma.social/",rel:"noopener noreferrer",target:"_blank",children:"Pleroma"})}),(0,r.jsx)(o.Z,{span:6,children:(0,r.jsx)("a",{href:"https://www.youtube.com/watch?v=S57uhCQBEk0",rel:"noopener noreferrer",target:"_blank",children:"Fediverse Explained Video"})})]}),(0,r.jsxs)(l.Z,{style:{marginTop:"20px"},children:[(0,r.jsx)(o.Z,{span:6,children:(0,r.jsx)("a",{href:"https://joinmastodon.org/",rel:"noopener noreferrer",target:"_blank",children:"Mastodon"})}),(0,r.jsx)(o.Z,{span:6,children:(0,r.jsx)("a",{href:"https://join.misskey.page/en/",rel:"noopener noreferrer",target:"_blank",children:"Misskey"})})]})]})]})}function k(){var e=function(){E(null)},n=function(){try{(0,m.Si)({apiPath:m.pE,data:{value:t.blockedDomains},onSuccess:function(){N({fieldName:"forbiddenUsernames",value:t.forbiddenUsernames}),E(b.zv),setTimeout(e,m.sI)},onError:function(n){E((0,b.kg)(b.Un,n)),setTimeout(e,m.sI)}})}catch(n){console.error(n),E(b.Un)}},i=s.Z.Title,a=(0,c.useState)(null),t=a[0],d=a[1],u=(0,c.useState)(!1),k=u[0],y=u[1],P=(0,c.useContext)(g.aC)||{},Z=P.serverConfig,N=P.setFieldInConfigState,C=(0,c.useState)(null),S=C[0],E=C[1],F=Z.federation,_=Z.yp,T=Z.instanceDetails,O=F.enabled,D=F.isPrivate,U=F.username,L=F.goLiveMessage,I=F.showEngagement,R=F.blockedDomains,M=_.instanceUrl,q=T.nsfw,V=function(e){var n=e.fieldName,i=e.value;d(j({},t,x({},n,i)))};if((0,c.useEffect)((function(){d({enabled:O,isPrivate:D,username:U,goLiveMessage:L,showEngagement:I,blockedDomains:R,nsfw:q,instanceUrl:_.instanceUrl})}),[Z,_]),!t)return null;var W=""!==M,z=M.startsWith("https://");return(0,r.jsxs)("div",{children:[(0,r.jsx)(i,{children:"Configure Social Features"}),"TODO: Explain what the Fediverse is here and talk about what happens if you were to enable this feature.",(0,r.jsxs)(l.Z,{children:[(0,r.jsxs)(o.Z,{span:15,className:"form-module",style:{marginRight:"15px"},children:[(0,r.jsx)(p.Z,j({fieldName:"enabled",onChange:function(e){e?y(!0):d(j({},t,{enabled:!1}))}},m.Kl,{checked:t.enabled,disabled:!W||!z})),(0,r.jsx)(f.ZP,j({fieldName:"instanceUrl"},m.yi,{value:t.instanceUrl,initialValue:_.instanceUrl,type:h.xA,onChange:V,onSubmit:function(){var e=""!==t.instanceUrl,n=t.instanceUrl.startsWith("https://");e&&n||((0,m.Si)({apiPath:m.Kl.apiPath,data:{value:!1}}),d(j({},t,{enabled:!1})))},disabled:!O})),(0,r.jsx)(p.Z,j({fieldName:"isPrivate"},m.LC,{checked:t.isPrivate,disabled:!O})),(0,r.jsx)(p.Z,j({fieldName:"nsfw",useSubmit:!0},m.B_,{checked:t.nsfw,disabled:!W})),(0,r.jsx)(f.ZP,j({required:!0,fieldName:"username"},m.Xc,{value:t.username,initialValue:U,onChange:V,disabled:!O})),(0,r.jsx)(f.ZP,j({fieldName:"goLiveMessage"},m.BF,{type:h.Sk,value:t.goLiveMessage,initialValue:L,onChange:V,disabled:!O})),(0,r.jsx)(p.Z,j({fieldName:"showEngagement"},m.FE,{checked:t.showEngagement,disabled:!O}))]}),(0,r.jsx)(o.Z,{span:8,className:"form-module",children:(0,r.jsx)(v.Z,{title:m.dR.label,placeholder:m.dR.placeholder,description:m.dR.tip,values:t.blockedDomains,handleDeleteIndex:function(e){t.blockedDomains.splice(e,1),n()},handleCreateString:function(e){var i;try{i=new URL(e).host}catch(r){i=e}t.blockedDomains.push(i),V({fieldName:"blockedDomains",value:t.blockedDomains}),n()},submitStatus:(0,b.kg)(S)})})]}),k&&(0,r.jsx)(w,{cancelPressed:function(){y(!1),d(j({},t,{enabled:!1}))},okPressed:function(){y(!1),d(j({},t,{enabled:!0}))}})]})}w.propTypes={cancelPressed:u().func.isRequired,okPressed:u().func.isRequired}}},function(e){e.O(0,[674,774,888,179],(function(){return n=50291,e(e.s=n);var n}));var n=e.O();_N_E=n}]);