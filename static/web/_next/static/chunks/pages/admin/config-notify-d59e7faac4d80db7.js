(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[4440],{15746:function(e,n,t){"use strict";var i=t(21584);n.Z=i.Z},97183:function(e,n,t){"use strict";var i=t(2897),a=t(7293),s=i.ZP;s.Header=i.h4,s.Footer=i.$_,s.Content=i.VY,s.Sider=a.Z,n.Z=s},71230:function(e,n,t){"use strict";var i=t(92820);n.Z=i.Z},94594:function(e,n,t){"use strict";t.d(n,{Z:function(){return y}});var i=t(87462),a=t(4942),s=t(50888),l=t(94184),r=t.n(l),o=t(97685),c=t(45987),d=t(67294),u=t(21770),h=t(15105),f=d.forwardRef(function(e,n){var t,i=e.prefixCls,s=void 0===i?"rc-switch":i,l=e.className,f=e.checked,p=e.defaultChecked,m=e.disabled,x=e.loadingIcon,g=e.checkedChildren,b=e.unCheckedChildren,v=e.onClick,y=e.onChange,j=e.onKeyDown,C=(0,c.Z)(e,["prefixCls","className","checked","defaultChecked","disabled","loadingIcon","checkedChildren","unCheckedChildren","onClick","onChange","onKeyDown"]),w=(0,u.Z)(!1,{value:f,defaultValue:p}),k=(0,o.Z)(w,2),Z=k[0],N=k[1];function E(e,n){var t=Z;return m||(N(t=e),null==y||y(t,n)),t}var S=r()(s,l,(t={},(0,a.Z)(t,"".concat(s,"-checked"),Z),(0,a.Z)(t,"".concat(s,"-disabled"),m),t));return d.createElement("button",Object.assign({},C,{type:"button",role:"switch","aria-checked":Z,disabled:m,className:S,ref:n,onKeyDown:function(e){e.which===h.Z.LEFT?E(!1,e):e.which===h.Z.RIGHT&&E(!0,e),null==j||j(e)},onClick:function(e){var n=E(!Z,e);null==v||v(n,e)}}),x,d.createElement("span",{className:"".concat(s,"-inner")},Z?g:b))});f.displayName="Switch";var p=t(53124),m=t(98866),x=t(97647),g=t(68349),b=function(e,n){var t={};for(var i in e)Object.prototype.hasOwnProperty.call(e,i)&&0>n.indexOf(i)&&(t[i]=e[i]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var a=0,i=Object.getOwnPropertySymbols(e);a<i.length;a++)0>n.indexOf(i[a])&&Object.prototype.propertyIsEnumerable.call(e,i[a])&&(t[i[a]]=e[i[a]]);return t},v=d.forwardRef(function(e,n){var t,l=e.prefixCls,o=e.size,c=e.disabled,u=e.loading,h=e.className,v=b(e,["prefixCls","size","disabled","loading","className"]),y=d.useContext(p.E_),j=y.getPrefixCls,C=y.direction,w=d.useContext(x.Z),k=d.useContext(m.Z),Z=j("switch",l),N=d.createElement("div",{className:"".concat(Z,"-handle")},u&&d.createElement(s.Z,{className:"".concat(Z,"-loading-icon")})),E=r()((t={},(0,a.Z)(t,"".concat(Z,"-small"),"small"===(o||w)),(0,a.Z)(t,"".concat(Z,"-loading"),u),(0,a.Z)(t,"".concat(Z,"-rtl"),"rtl"===C),t),void 0===h?"":h);return d.createElement(g.Z,{insertExtraNode:!0},d.createElement(f,(0,i.Z)({},v,{prefixCls:Z,className:E,disabled:(null!=c?c:k)||u,ref:n,loadingIcon:N})))});v.__ANT_SWITCH=!0;var y=v},93645:function(e,n,t){"use strict";t.d(n,{u:function(){return a}});var i={ceil:Math.ceil,round:Math.round,floor:Math.floor,trunc:function(e){return e<0?Math.ceil(e):Math.floor(e)}};function a(e){return e?i[e]:i.trunc}},59910:function(e,n,t){"use strict";t.d(n,{Z:function(){return s}});var i=t(19013),a=t(13882);function s(e,n){return(0,a.Z)(2,arguments),(0,i.Z)(e).getTime()-(0,i.Z)(n).getTime()}},11699:function(e,n,t){"use strict";t.d(n,{Z:function(){return l}});var i=t(59910),a=t(13882),s=t(93645);function l(e,n,t){(0,a.Z)(2,arguments);var l=(0,i.Z)(e,n)/1e3;return(0,s.u)(null==t?void 0:t.roundingMethod)(l)}},7148:function(e,n,t){(window.__NEXT_P=window.__NEXT_P||[]).push(["/admin/config-notify",function(){return t(4391)}])},44716:function(e,n,t){"use strict";t.d(n,{Z:function(){return d}});var i=t(85893),a=t(67294),s=t(94594),l=t(91332),r=t(57520),o=t(24044),c=t(38631);let d=e=>{let{apiPath:n,checked:t,reversed:d=!1,configPath:u="",disabled:h=!1,fieldName:f,label:p,tip:m,useSubmit:x,onChange:g}=e,[b,v]=(0,a.useState)(null),y=null,j=(0,a.useContext)(c.aC),{setFieldInConfigState:C}=j||{},w=()=>{v(null),clearTimeout(y),y=null},k=async e=>{if(x){v((0,l.kg)(l.Jk));let t=d?!e:e;await (0,o.Si)({apiPath:n,data:{value:t},onSuccess:()=>{C({fieldName:f,value:t,path:u}),v((0,l.kg)(l.zv))},onError:e=>{v((0,l.kg)(l.Un,"There was an error: ".concat(e)))}}),y=setTimeout(w,o.sI)}g&&g(e)},Z=null!==b&&b.type===l.Jk;return(0,i.jsxs)("div",{className:"formfield-container toggleswitch-container",children:[p&&(0,i.jsx)("div",{className:"label-side",children:(0,i.jsx)("span",{className:"formfield-label",children:p})}),(0,i.jsxs)("div",{className:"input-side",children:[(0,i.jsxs)("div",{className:"input-group",children:[(0,i.jsx)(s.Z,{className:"switch field-".concat(f),loading:Z,onChange:k,defaultChecked:t,checked:t,checkedChildren:"ON",unCheckedChildren:"OFF",disabled:h}),(0,i.jsx)(r.E,{status:b})]}),(0,i.jsx)("p",{className:"field-tip",children:m})]})]})};d.defaultProps={apiPath:"",checked:!1,reversed:!1,configPath:"",disabled:!1,label:"",tip:"",useSubmit:!1,onChange:null}},4391:function(e,n,t){"use strict";t.r(n),t.d(n,{default:function(){return S}});var i=t(85893),a=t(85818),s=t(14670),l=t(71230),r=t(15746),o=t(71577),c=t(67294),d=t(41664),u=t.n(d),h=t(38631),f=t(67032),p=t(57520),m=t(24044),x=t(44716),g=t(91332);let{Title:b}=a.Z,v=()=>{let e=(0,c.useContext)(h.aC),{serverConfig:n,setFieldInConfigState:t}=e||{},{notifications:a}=n||{},{discord:s}=a||{},{enabled:l,webhook:r,goLiveMessage:d}=s||{},[u,v]=(0,c.useState)({}),[y,j]=(0,c.useState)(null),[C,w]=(0,c.useState)(!1);(0,c.useEffect)(()=>{v({enabled:l,webhook:r,goLiveMessage:d})},[a,s]);let k=()=>""!==r&&""!==d,Z=e=>{let{fieldName:n,value:t}=e;v({...u,[n]:t}),w(k())},N=()=>{j(null),clearTimeout(null)},E=async()=>{await (0,m.Si)({apiPath:"/notifications/discord",data:{value:u},onSuccess:()=>{t({fieldName:"discord",value:u,path:"notifications"}),j((0,g.kg)(g.zv,"Updated.")),setTimeout(N,m.sI)},onError:e=>{j((0,g.kg)(g.Un,e)),setTimeout(N,m.sI)}})},S=e=>{Z({fieldName:"enabled",value:e})};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(b,{children:"Discord"}),(0,i.jsx)("p",{className:"description reduced-margins",children:"Let your Discord channel know each time you go live."}),(0,i.jsxs)("p",{className:"description reduced-margins",children:[(0,i.jsx)("a",{href:"https://support.discord.com/hc/en-us/articles/228383668",target:"_blank",rel:"noreferrer",children:"Create a webhook"})," ","under ",(0,i.jsx)("i",{children:"Edit Channel / Integrations"})," on your Discord channel and provide it below."]}),(0,i.jsx)(x.Z,{apiPath:"",fieldName:"discordEnabled",label:"Enable Discord",checked:u.enabled,onChange:S}),(0,i.jsx)("div",{style:{display:u.enabled?"block":"none"},children:(0,i.jsx)(f.nv,{...m.oy.webhookUrl,required:!0,value:u.webhook,onChange:Z})}),(0,i.jsx)("div",{style:{display:u.enabled?"block":"none"},children:(0,i.jsx)(f.nv,{...m.oy.goLiveMessage,required:!0,value:u.goLiveMessage,onChange:Z})}),(0,i.jsx)(o.Z,{type:"primary",onClick:E,style:{display:C?"inline-block":"none",position:"relative",marginLeft:"auto",right:"0",marginTop:"20px"},children:"Save"}),(0,i.jsx)(p.E,{status:y})]})},{Title:y}=a.Z,j=()=>{let e=(0,c.useContext)(h.aC),{serverConfig:n,setFieldInConfigState:t}=e||{},{notifications:a}=n||{},{browser:s}=a||{},{enabled:l,goLiveMessage:r}=s||{},[d,u]=(0,c.useState)({}),[b,v]=(0,c.useState)(null),[j,C]=(0,c.useState)(!1);(0,c.useEffect)(()=>{u({enabled:l,goLiveMessage:r})},[a,s]);let w=()=>!0,k=e=>{let{fieldName:n,value:t}=e;console.log(n,t),u({...d,[n]:t}),C(w())},Z=e=>{k({fieldName:"enabled",value:e})},N=()=>{v(null),clearTimeout(null)},E=async()=>{await (0,m.Si)({apiPath:"/notifications/browser",data:{value:d},onSuccess:()=>{t({fieldName:"browser",value:d,path:"notifications"}),v((0,g.kg)(g.zv,"Updated.")),setTimeout(N,m.sI)},onError:e=>{v((0,g.kg)(g.Un,e)),setTimeout(N,m.sI)}})};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(y,{children:"Browser Alerts"}),(0,i.jsx)("p",{className:"description reduced-margins",children:"Viewers can opt into being notified when you go live with their browser."}),(0,i.jsx)("p",{className:"description reduced-margins",children:"Not all browsers support this."}),(0,i.jsx)(x.Z,{apiPath:"",fieldName:"enabled",label:"Enable browser notifications",onChange:Z,checked:d.enabled}),(0,i.jsx)("div",{style:{display:d.enabled?"block":"none"},children:(0,i.jsx)(f.nv,{...m.mv.goLiveMessage,required:!0,type:f.Sk,value:d.goLiveMessage,onChange:k})}),(0,i.jsx)(o.Z,{type:"primary",style:{display:j?"inline-block":"none",position:"relative",marginLeft:"auto",right:"0",marginTop:"20px"},onClick:E,children:"Save"}),(0,i.jsx)(p.E,{status:b})]})},{Title:C}=a.Z,w=()=>{let e=(0,c.useContext)(h.aC),{serverConfig:n}=e||{},{federation:t}=n||{},{enabled:a}=t||{},[s,l]=(0,c.useState)({});return(0,c.useEffect)(()=>{l({enabled:a})},[a]),(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(C,{children:"Fediverse Social"}),(0,i.jsx)("p",{className:"description",children:"Enabling the Fediverse social features will not just alert people to when you go live, but also enable other functionality."}),(0,i.jsxs)("p",{children:["Fediverse social features:"," ",(0,i.jsx)("span",{style:{color:t.enabled?"green":"red"},children:s.enabled?"Enabled":"Disabled"})]}),(0,i.jsx)(u(),{passHref:!0,href:"/admin/config-federation/",children:(0,i.jsx)(o.Z,{type:"primary",style:{position:"relative",marginLeft:"auto",right:"0",marginTop:"20px"},children:"Configure"})})]})};var k=t(78353),Z=t(89126),N=t(34261);let{Title:E}=a.Z;function S(){let[e,n]=(0,c.useState)(null),t=(0,c.useContext)(h.aC),{serverConfig:a}=t||{},{yp:d}=a,{instanceUrl:f}=d,[p,x]=(0,c.useState)(!1);(0,c.useEffect)(()=>{n({instanceUrl:f})},[d]);let g=()=>{p&&n({...e,enabled:!1})},b=t=>{let{fieldName:i,value:a}=t;x((0,Z.jv)(a)),n({...e,[i]:a})},y=""!==f,C=!y&&(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(s.Z,{message:"You must set your server URL before you can enable this feature.",type:"warning",showIcon:!0}),(0,i.jsx)("br",{}),(0,i.jsx)(k.$7,{fieldName:"instanceUrl",...m.yi,value:(null==e?void 0:e.instanceUrl)||"",initialValue:d.instanceUrl,type:k.xA,onChange:b,onSubmit:g,required:!0})]});return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(E,{children:"Notifications"}),(0,i.jsxs)("p",{className:"description",children:["Let your viewers know when you go live by supporting any of the below notification channels."," ",(0,i.jsx)("a",{href:"https://owncast.online/docs/notifications/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn more about live notifications."})]}),C,(0,i.jsxs)(l.Z,{children:[(0,i.jsx)(r.Z,{span:10,className:"form-module ".concat(y?"":"disabled"),style:{margin:"5px",display:"flex",flexDirection:"column"},children:(0,i.jsx)(j,{})}),(0,i.jsx)(r.Z,{span:10,className:"form-module ".concat(y?"":"disabled"),style:{margin:"5px",display:"flex",flexDirection:"column"},children:(0,i.jsx)(v,{})}),(0,i.jsx)(r.Z,{span:10,className:"form-module ".concat(y?"":"disabled"),style:{margin:"5px",display:"flex",flexDirection:"column"},children:(0,i.jsx)(w,{})}),(0,i.jsxs)(r.Z,{span:10,className:"form-module ".concat(y?"":"disabled"),style:{margin:"5px",display:"flex",flexDirection:"column"},children:[(0,i.jsx)(E,{children:"Custom"}),(0,i.jsx)("p",{className:"description",children:"Build your own notifications by using custom webhooks."}),(0,i.jsx)(u(),{passHref:!0,href:"/admin/webhooks",children:(0,i.jsx)(o.Z,{type:"primary",style:{position:"relative",marginLeft:"auto",right:"0",marginTop:"20px"},children:"Create"})})]})]})]})}S.getLayout=function(e){return(0,i.jsx)(N.l,{page:e})}},9008:function(e,n,t){e.exports=t(42636)},11163:function(e,n,t){e.exports=t(96885)}},function(e){e.O(0,[173,5874,2184,2364,4931,5402,5257,5818,1664,8014,9915,8035,4261,9774,2888,179],function(){return e(e.s=7148)}),_N_E=e.O()}]);