(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[2532],{5789:function(e,t,n){"use strict";t.Z=void 0;var a=n(38614).Col;t.Z=a},55673:function(e,t,n){"use strict";t.Z=void 0;var a=n(38614).Row;t.Z=a},83514:function(e,t,n){"use strict";var a=n(75263).default,l=n(64836).default;Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var s=l(n(10434)),r=l(n(38416)),o=l(n(94184)),i=a(n(67294)),c=n(31929),__rest=function(e,t){var n={};for(var a in e)Object.prototype.hasOwnProperty.call(e,a)&&0>t.indexOf(a)&&(n[a]=e[a]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var l=0,a=Object.getOwnPropertySymbols(e);l<a.length;l++)0>t.indexOf(a[l])&&Object.prototype.propertyIsEnumerable.call(e,a[l])&&(n[a[l]]=e[a[l]]);return n};t.default=function(e){var t,n=e.prefixCls,a=e.className,l=e.checked,d=e.onChange,u=e.onClick,h=__rest(e,["prefixCls","className","checked","onChange","onClick"]),f=(0,i.useContext(c.ConfigContext).getPrefixCls)("tag",n),p=(0,o.default)(f,(t={},(0,r.default)(t,"".concat(f,"-checkable"),!0),(0,r.default)(t,"".concat(f,"-checkable-checked"),l),t),a);return i.createElement("span",(0,s.default)({},h,{className:p,onClick:function(e){null==d||d(!l),null==u||u(e)}}))}},59361:function(e,t,n){"use strict";var a=n(75263).default,l=n(64836).default;t.Z=void 0;var s=l(n(38416)),r=l(n(10434)),o=l(n(27424)),i=l(n(40753)),c=l(n(94184)),d=l(n(18475)),u=a(n(67294)),h=n(31929),f=n(45471),p=l(n(61539));l(n(13594));var g=l(n(83514)),__rest=function(e,t){var n={};for(var a in e)Object.prototype.hasOwnProperty.call(e,a)&&0>t.indexOf(a)&&(n[a]=e[a]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var l=0,a=Object.getOwnPropertySymbols(e);l<a.length;l++)0>t.indexOf(a[l])&&Object.prototype.propertyIsEnumerable.call(e,a[l])&&(n[a[l]]=e[a[l]]);return n},m=new RegExp("^(".concat(f.PresetColorTypes.join("|"),")(-inverse)?$")),b=new RegExp("^(".concat(f.PresetStatusColorTypes.join("|"),")$")),v=u.forwardRef(function(e,t){var n,a=e.prefixCls,l=e.className,f=e.style,g=e.children,v=e.icon,y=e.color,x=e.onClose,w=e.closeIcon,j=e.closable,C=void 0!==j&&j,k=__rest(e,["prefixCls","className","style","children","icon","color","onClose","closeIcon","closable"]),S=u.useContext(h.ConfigContext),N=S.getPrefixCls,P=S.direction,E=u.useState(!0),O=(0,o.default)(E,2),_=O[0],F=O[1];u.useEffect(function(){"visible"in k&&F(k.visible)},[k.visible]);var isPresetColor=function(){return!!y&&(m.test(y)||b.test(y))},T=(0,r.default)({backgroundColor:y&&!isPresetColor()?y:void 0},f),D=isPresetColor(),Z=N("tag",a),U=(0,c.default)(Z,(n={},(0,s.default)(n,"".concat(Z,"-").concat(y),D),(0,s.default)(n,"".concat(Z,"-has-color"),y&&!D),(0,s.default)(n,"".concat(Z,"-hidden"),!_),(0,s.default)(n,"".concat(Z,"-rtl"),"rtl"===P),n),l),handleCloseClick=function(e){e.stopPropagation(),null==x||x(e),!e.defaultPrevented&&("visible"in k||F(!1))},I="onClick"in k||g&&"a"===g.type,L=(0,d.default)(k,["visible"]),R=v||null,B=R?u.createElement(u.Fragment,null,R,u.createElement("span",null,g)):g,V=u.createElement("span",(0,r.default)({},L,{ref:t,className:U,style:T}),B,C?w?u.createElement("span",{className:"".concat(Z,"-close-icon"),onClick:handleCloseClick},w):u.createElement(i.default,{className:"".concat(Z,"-close-icon"),onClick:handleCloseClick}):null);return I?u.createElement(p.default,null,V):V});v.CheckableTag=g.default,t.Z=v},13882:function(e,t,n){"use strict";function requiredArgs(e,t){if(t.length<e)throw TypeError(e+" argument"+(e>1?"s":"")+" required, but only "+t.length+" present")}n.d(t,{Z:function(){return requiredArgs}})},19013:function(e,t,n){"use strict";n.d(t,{Z:function(){return toDate}});var a=n(71002),l=n(13882);function toDate(e){(0,l.Z)(1,arguments);var t=Object.prototype.toString.call(e);return e instanceof Date||"object"===(0,a.Z)(e)&&"[object Date]"===t?new Date(e.getTime()):"number"==typeof e||"[object Number]"===t?new Date(e):(("string"==typeof e||"[object String]"===t)&&"undefined"!=typeof console&&(console.warn("Starting with v2.0.0-beta.1 date-fns doesn't accept strings as date arguments. Please use `parseISO` to parse strings. See: https://github.com/date-fns/date-fns/blob/master/docs/upgradeGuide.md#string-arguments"),console.warn(Error().stack)),new Date(NaN))}},18957:function(e,t,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/admin/config-federation",function(){return n(12281)}])},62642:function(e,t,n){"use strict";n.d(t,{Q:function(){return d},Y:function(){return EditValueArray}});var a=n(85893),l=n(67294),s=n(53740),r=n(59361),o=n(34568),i=n(50988);let{Title:c}=s.default,d="#5a67d8",EditValueArray=e=>{let{title:t,description:n,placeholder:s,maxLength:u,values:h,handleDeleteIndex:f,handleCreateString:p,submitStatus:g,continuousStatusMessage:m}=e,[b,v]=(0,l.useState)("");return(0,a.jsxs)("div",{className:"edit-string-array-container",children:[(0,a.jsx)(c,{level:3,className:"section-title",children:t}),(0,a.jsx)("p",{className:"description",children:n}),(0,a.jsx)("div",{className:"edit-current-strings",children:null==h?void 0:h.map((e,t)=>(0,a.jsx)(r.Z,{closable:!0,onClose:()=>{f(t)},color:d,children:e},"tag-".concat(e,"-").concat(t)))}),m&&(0,a.jsx)("div",{className:"continuous-status-section",children:(0,a.jsx)(i.E,{status:m})}),(0,a.jsx)("div",{className:"add-new-string-section",children:(0,a.jsx)(o.nv,{fieldName:"string-input",value:b,onChange:e=>{let{value:t}=e;v(t)},onPressEnter:()=>{let e=b.trim();p(e),v("")},maxLength:u,placeholder:s,status:g})})]})};EditValueArray.defaultProps={maxLength:50,description:null,submitStatus:null,continuousStatusMessage:null}},80910:function(e,t,n){"use strict";n.d(t,{Z:function(){return ToggleSwitch}});var a=n(85893),l=n(67294),s=n(40987),r=n(24770),o=n(50988),i=n(42124),c=n(3273);let ToggleSwitch=e=>{let{apiPath:t,checked:n,reversed:d=!1,configPath:u="",disabled:h=!1,fieldName:f,label:p,tip:g,useSubmit:m,onChange:b}=e,[v,y]=(0,l.useState)(null),x=null,w=(0,l.useContext)(c.a),{setFieldInConfigState:j}=w||{},resetStates=()=>{y(null),clearTimeout(x),x=null},handleChange=async e=>{if(m){y((0,r.kg)(r.Jk));let n=d?!e:e;await (0,i.Si)({apiPath:t,data:{value:n},onSuccess:()=>{j({fieldName:f,value:n,path:u}),y((0,r.kg)(r.zv))},onError:e=>{y((0,r.kg)(r.Un,"There was an error: ".concat(e)))}}),x=setTimeout(resetStates,i.sI)}b&&b(e)},C=null!==v&&v.type===r.Jk;return(0,a.jsxs)("div",{className:"formfield-container toggleswitch-container",children:[p&&(0,a.jsx)("div",{className:"label-side",children:(0,a.jsx)("span",{className:"formfield-label",children:p})}),(0,a.jsxs)("div",{className:"input-side",children:[(0,a.jsxs)("div",{className:"input-group",children:[(0,a.jsx)(s.Z,{className:"switch field-".concat(f),loading:C,onChange:handleChange,defaultChecked:n,checked:n,checkedChildren:"ON",unCheckedChildren:"OFF",disabled:h}),(0,a.jsx)(o.E,{status:v})]}),(0,a.jsx)("p",{className:"field-tip",children:g})]})]})};ToggleSwitch.defaultProps={apiPath:"",checked:!1,reversed:!1,configPath:"",disabled:!1,label:"",tip:"",useSubmit:!1,onChange:null}},12281:function(e,t,n){"use strict";n.r(t);var a=n(85893),l=n(53740),s=n(56697),r=n(65400),o=n(55673),i=n(5789),c=n(4863),d=n(67294),u=n(34568),h=n(28567),f=n(80910),p=n(62642),g=n(42124),m=n(3273),b=n(24770),v=n(90695);let FederationInfoModal=e=>{let{cancelPressed:t,okPressed:n}=e;return(0,a.jsxs)(s.default,{width:"70%",title:"Enable Social Features",visible:!0,onCancel:t,footer:(0,a.jsxs)("div",{children:[(0,a.jsx)(r.default,{onClick:t,children:"Do not enable"}),(0,a.jsx)(r.default,{type:"primary",onClick:n,children:"Enable Social Features"})]}),children:[(0,a.jsx)(l.default.Title,{level:3,children:"How do Owncast's social features work?"}),(0,a.jsxs)(l.default.Paragraph,{children:["Owncast's social features are accomplished by having your server join The"," ",(0,a.jsx)("a",{href:"https://en.wikipedia.org/wiki/Fediverse",rel:"noopener noreferrer",target:"_blank",children:"Fediverse"}),", a decentralized, open, collection of independent servers, like yours."]}),"Please"," ",(0,a.jsx)("a",{href:"https://owncast.online/docs/social",rel:"noopener noreferrer",target:"_blank",children:"read more"})," ","about these features, the details behind them, and how they work.",(0,a.jsx)(l.default.Paragraph,{}),(0,a.jsx)(l.default.Title,{level:3,children:"What do you need to know?"}),(0,a.jsxs)("ul",{children:[(0,a.jsx)("li",{children:"These features are brand new. Given the variability of interfacing with the rest of the world, bugs are possible. Please report anything that you think isn't working quite right."}),(0,a.jsx)("li",{children:"You must always host your Owncast server with SSL using a https url."}),(0,a.jsx)("li",{children:"You should not change your server name URL or social username once people begin following you, as your server will be seen as a completely different user on the Fediverse, and the old user will disappear."}),(0,a.jsxs)("li",{children:["Turning on ",(0,a.jsx)("i",{children:"Private mode"})," will allow you to manually approve each follower and limit the visibility of your posts to followers only."]})]}),(0,a.jsx)(l.default.Title,{level:3,children:"Learn more about The Fediverse"}),(0,a.jsxs)(l.default.Paragraph,{children:["If these concepts are new you should discover more about what this functionality has to offer. Visit"," ",(0,a.jsx)("a",{href:"https://owncast.online/docs/social",rel:"noopener noreferrer",target:"_blank",children:"our documentation"})," ","to be pointed at some resources that will help get you started on The Fediverse."]})]})},ConfigFederation=()=>{let{Title:e}=l.default,[t,n]=(0,d.useState)(null),[s,r]=(0,d.useState)(!1),v=(0,d.useContext)(m.a),{serverConfig:y,setFieldInConfigState:x}=v||{},[w,j]=(0,d.useState)(null),{federation:C,yp:k,instanceDetails:S}=y,{enabled:N,isPrivate:P,username:E,goLiveMessage:O,showEngagement:_,blockedDomains:F}=C,{instanceUrl:T}=k,{nsfw:D}=S,handleFieldChange=e=>{let{fieldName:a,value:l}=e;n({...t,[a]:l})};function resetBlockedDomainsSaveState(){j(null)}function saveBlockedDomains(){try{(0,g.Si)({apiPath:g.pE,data:{value:t.blockedDomains},onSuccess:()=>{x({fieldName:"forbiddenUsernames",value:t.forbiddenUsernames}),j(b.zv),setTimeout(resetBlockedDomainsSaveState,g.sI)},onError:e=>{j((0,b.kg)(b.Un,e)),setTimeout(resetBlockedDomainsSaveState,g.sI)}})}catch(e){console.error(e),j(b.Un)}}if((0,d.useEffect)(()=>{n({enabled:N,isPrivate:P,username:E,goLiveMessage:O,showEngagement:_,blockedDomains:F,nsfw:D,instanceUrl:k.instanceUrl})},[y,k]),!t)return null;let Z=""!==T,U=T.startsWith("https://"),I=!U&&(0,a.jsxs)(a.Fragment,{children:[(0,a.jsx)(c.default,{message:"You must set your server URL before you can enable this feature.",type:"warning",showIcon:!0}),(0,a.jsx)("br",{}),(0,a.jsx)(h.$7,{fieldName:"instanceUrl",...g.yi,value:t.instanceUrl,initialValue:k.instanceUrl,type:u.xA,onChange:handleFieldChange,onSubmit:()=>{let e=""!==t.instanceUrl,a=t.instanceUrl.startsWith("https://");e&&a||((0,g.Si)({apiPath:g.Kl.apiPath,data:{value:!1}}),n({...t,enabled:!1}))},required:!0})]}),L=(0,a.jsx)(c.default,{message:"Only Owncast instances available on the default SSL port 443 support this feature.",type:"warning",showIcon:!0}),R=T&&""!==new URL(T).port&&"443"!==new URL(T).port;return(0,a.jsxs)("div",{children:[(0,a.jsx)(e,{children:"Configure Social Features"}),(0,a.jsx)("p",{children:"Owncast provides the ability for people to follow and engage with your instance. It's a great way to promote alerting, sharing and engagement of your stream."}),(0,a.jsx)("p",{children:"Once enabled you'll alert your followers when you go live as well as gain the ability to compose custom posts to share any information you like."}),(0,a.jsx)("p",{children:(0,a.jsx)("a",{href:"https://owncast.online/docs/social",rel:"noopener noreferrer",target:"_blank",children:"Read more about the specifics of these social features."})}),(0,a.jsxs)(o.Z,{children:[(0,a.jsxs)(i.Z,{span:15,className:"form-module",style:{marginRight:"15px"},children:[I,R&&L,(0,a.jsx)(f.Z,{fieldName:"enabled",onChange:e=>{e?r(!0):n({...t,enabled:!1})},...g.Kl,checked:t.enabled,disabled:R||!Z||!U}),(0,a.jsx)(f.Z,{fieldName:"isPrivate",...g.LC,checked:t.isPrivate,disabled:!N}),(0,a.jsx)(f.Z,{fieldName:"nsfw",useSubmit:!0,...g.B_,checked:t.nsfw,disabled:R||!Z}),(0,a.jsx)(h.$7,{required:!0,fieldName:"username",type:u.Kx,...g.Xc,value:t.username,initialValue:E,onChange:handleFieldChange,disabled:!N}),(0,a.jsx)(h.$7,{fieldName:"goLiveMessage",...g.BF,type:u.Sk,value:t.goLiveMessage,initialValue:O,onChange:handleFieldChange,disabled:!N}),(0,a.jsx)(f.Z,{fieldName:"showEngagement",...g.FE,checked:t.showEngagement,disabled:!N})]}),(0,a.jsx)(i.Z,{span:8,className:"form-module",children:(0,a.jsx)(p.Y,{title:g.dR.label,placeholder:g.dR.placeholder,description:g.dR.tip,values:t.blockedDomains,handleDeleteIndex:function(e){t.blockedDomains.splice(e,1),saveBlockedDomains()},handleCreateString:function(e){let n;try{let t=new URL(e);n=t.host}catch(t){n=e}t.blockedDomains.push(n),handleFieldChange({fieldName:"blockedDomains",value:t.blockedDomains}),saveBlockedDomains()},submitStatus:(0,b.kg)(w)})})]}),s&&(0,a.jsx)(FederationInfoModal,{cancelPressed:function(){r(!1),n({...t,enabled:!1})},okPressed:function(){r(!1),n({...t,enabled:!0})}})]})};ConfigFederation.getLayout=function(e){return(0,a.jsx)(v.l,{page:e})},t.default=ConfigFederation}},function(e){e.O(0,[5596,1130,4104,9403,1024,3942,971,6697,1664,1749,1700,2122,4938,695,9774,2888,179],function(){return e(e.s=18957)}),_N_E=e.O()}]);