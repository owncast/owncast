(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[5405],{48312:function(e,t,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/",function(){return n(44486)}])},37527:function(e,t,n){"use strict";n.d(t,{h:function(){return E},Z:function(){return M}});var o=n(85893),s=n(97183),i=n(94199),a=n(20550),l=n(94184),r=n.n(l),c=n(66516),d=n(38504),h=n(71577),u=n(26713),m=n(86548),x=n(94149),p=n(38545),f=n(87547),j=n(39398),_=n(4480),g=n(67294),w=n(49218),v=n(5152),y=n.n(v),b=n(77466),C=n(79252),k=n.n(C);let N=y()(()=>n.e(4761).then(n.bind(n,14761)).then(e=>e.Modal),{loadableGenerated:{webpack:()=>[14761]}}),S=y()(()=>Promise.all([n.e(8939),n.e(9096)]).then(n.bind(n,59096)).then(e=>e.NameChangeModal),{loadableGenerated:{webpack:()=>[59096]}}),Z=y()(()=>Promise.all([n.e(4381),n.e(5535)]).then(n.bind(n,65535)).then(e=>e.AuthModal),{loadableGenerated:{webpack:()=>[65535]}}),F=e=>{let{username:t}=e,[n,s]=(0,g.useState)(!1),[i,a]=(0,g.useState)(!1),[l,r]=(0,_.FV)(b.ZA),v=(0,_.sJ)(b.Q),y=()=>{r(!l)},C=()=>{s(!0)};(0,w.y1)("c",y,{enableOnContentEditable:!1},[l]);let F=(0,_.sJ)(b.db);if(!F)return null;let{displayName:L}=F,T=(0,o.jsxs)(c.Z,{children:[(0,o.jsx)(c.Z.Item,{icon:(0,o.jsx)(m.Z,{}),onClick:()=>C(),children:"Change name"},"0"),(0,o.jsx)(c.Z.Item,{icon:(0,o.jsx)(x.Z,{}),onClick:()=>a(!0),children:"Authenticate"},"1"),v.chatAvailable&&(0,o.jsx)(c.Z.Item,{icon:(0,o.jsx)(p.Z,{}),onClick:()=>y(),children:"Toggle chat"},"3")]});return(0,o.jsxs)("div",{id:"user-menu",className:"".concat(k().root),children:[(0,o.jsx)(d.Z,{overlay:T,trigger:["click"],children:(0,o.jsx)(h.Z,{type:"primary",icon:(0,o.jsx)(f.Z,{style:{marginRight:".5rem"}}),children:(0,o.jsxs)(u.Z,{children:[t||L,(0,o.jsx)(j.Z,{})]})})}),(0,o.jsx)(N,{title:"Change Chat Display Name",open:n,handleCancel:()=>s(!1),children:(0,o.jsx)(S,{})}),(0,o.jsx)(N,{title:"Authenticate",open:i,handleCancel:()=>a(!1),children:(0,o.jsx)(Z,{})})]})};var L=n(50738),T=n(31764),A=n.n(T);let{Header:D}=s.Z,E=e=>{let{name:t="Your stream title",chatAvailable:n,chatDisabled:s}=e;return(0,o.jsxs)(D,{className:r()(["".concat(A().header)],"global-header"),children:[(0,o.jsxs)("div",{className:"".concat(A().logo),children:[(0,o.jsx)(L.C,{variant:"contrast"}),(0,o.jsx)("span",{className:"global-header-text",children:t})]}),n&&!s&&(0,o.jsx)(F,{}),!n&&!s&&(0,o.jsx)(i.Z,{title:"Chat is available when the stream is live.",placement:"left",children:(0,o.jsx)(a.Z,{style:{cursor:"pointer"},children:"Chat offline"})})]})};var M=E},14761:function(e,t,n){"use strict";n.r(t),n.d(t,{Modal:function(){return d}});var o=n(85893),s=n(85402),i=n(26303),a=n(11382),l=n(67294),r=n(77011),c=n.n(r);let d=e=>{let{title:t,url:n,open:r,handleOk:d,handleCancel:h,afterClose:u,height:m,width:x,children:p}=e,[f,j]=(0,l.useState)(!!n),_=n&&(0,o.jsx)("iframe",{title:t,src:n,width:"100%",height:"100%",sandbox:"allow-same-origin allow-scripts allow-popups allow-forms",frameBorder:"0",allowFullScreen:!0,onLoad:()=>j(!1)});return(0,o.jsx)(s.Z,{title:t,open:r,onOk:d,onCancel:h,afterClose:u,bodyStyle:{padding:"0px",minHeight:m,height:null!=m?m:"100%"},width:x,zIndex:9999,footer:null,centered:!0,destroyOnClose:!0,children:(0,o.jsxs)(o.Fragment,{children:[f&&(0,o.jsx)(i.Z,{active:f,style:{padding:"10px"},paragraph:{rows:10}}),_&&(0,o.jsx)("div",{style:{display:f?"none":"inline"},children:_}),p&&(0,o.jsx)("div",{className:c().content,children:p}),f&&(0,o.jsx)(a.Z,{className:c().spinner,spinning:f,size:"large"})]})})};t.default=d,d.defaultProps={url:void 0,children:void 0,handleOk:void 0,handleCancel:void 0,afterClose:void 0}},51513:function(e,t,n){"use strict";n.d(t,{R:function(){return c}});var o=n(85893),s=n(27049),i=n(24019),a=n(45938),l=n(88335),r=n.n(l);let c=e=>{let t,{streamName:n,customText:l,lastLive:c,notificationsEnabled:d,fediverseAccount:h,onNotifyClick:u,onFollowClick:m}=e;return t=l||(!l&&d&&h?(0,o.jsxs)("span",{children:["This stream is offline. You can"," ",(0,o.jsx)("span",{role:"link",tabIndex:0,className:r().actionLink,onClick:u,children:"be notified"})," ","the next time ",n," goes live or"," ",(0,o.jsx)("span",{role:"link",tabIndex:0,className:r().actionLink,onClick:m,children:"follow"})," ",h," on the Fediverse."]}):!l&&d?(0,o.jsxs)("span",{children:["This stream is offline."," ",(0,o.jsx)("span",{role:"link",tabIndex:0,className:r().actionLink,onClick:u,children:"Be notified"})," ","the next time ",n," goes live."]}):!l&&h?(0,o.jsxs)("span",{children:["This stream is offline."," ",(0,o.jsx)("span",{role:"link",tabIndex:0,className:r().actionLink,onClick:m,children:"Follow"})," ",h," on the Fediverse to see the next time ",n," goes live."]}):"This stream is offline. Check back soon!"),(0,o.jsx)("div",{className:r().outerContainer,children:(0,o.jsxs)("div",{className:r().innerContainer,children:[(0,o.jsx)("div",{className:r().header,children:n}),(0,o.jsx)(s.Z,{className:r().separator}),(0,o.jsx)("div",{className:r().bodyText,children:t}),c&&(0,o.jsxs)("div",{className:r().lastLiveDate,children:[(0,o.jsx)(i.Z,{className:r().clockIcon}),"Last live ".concat((0,a.Z)(new Date(c))," ago.")]})]})})}},69357:function(e,t,n){"use strict";n.d(t,{X:function(){return d}});var o=n(85893),s=n(45938),i=n(68730),a=n(67294),l=n(31326),r=n(37970),c=n.n(r);let d=e=>{let t,{online:n,lastConnectTime:r,lastDisconnectTime:d,viewerCount:h}=e,[,u]=(0,a.useState)(new Date);(0,a.useEffect)(()=>{let e=setInterval(()=>u(new Date),1e3);return()=>{clearInterval(e)}},[]);let m="";if(n&&r){let x=function(e){let t=(0,i.Z)({start:e,end:new Date});return t.days>1?"".concat(t.days," days ").concat(t.hours," hours"):t.hours>=1?"".concat(t.hours," hours ").concat(t.minutes," minutes"):"".concat(t.minutes," minutes ").concat(t.seconds," seconds")}(new Date(r));m=n?"Live for  ".concat(x):"Offline",t=h>0&&(0,o.jsxs)("div",{className:c().right,children:[(0,o.jsx)("span",{children:(0,o.jsx)(l.Z,{})}),(0,o.jsx)("span",{children:" ".concat(h)})]})}else!n&&(m="Offline",d&&(t="Last live ".concat((0,s.Z)(new Date(d))," ago.")));return(0,o.jsxs)("div",{className:c().statusbar,children:[(0,o.jsx)("div",{children:m}),(0,o.jsx)("div",{children:t})]})};d.defaultProps={lastConnectTime:null,lastDisconnectTime:null}},44486:function(e,t,n){"use strict";n.r(t),n.d(t,{default:function(){return eG}});var o=n(85893),s=n(97183),i=n(4480),a=n(9008),l=n.n(a),r=n(67294),c=n(77466),d=n(84381),h=n(11382),u=n(94184),m=n.n(u),x=n(5152),p=n.n(x),f=n(72581),j=n(67917),_=n.n(j);let g=e=>{let{version:t}=e;return(0,o.jsxs)("div",{className:_().footer,children:[(0,o.jsxs)("span",{children:["Powered by ",(0,o.jsx)("a",{href:"https://owncast.online",children:t})]}),(0,o.jsxs)("span",{className:_().links,children:[(0,o.jsx)("a",{href:"https://owncast.online/docs",target:"_blank",rel:"noreferrer",children:"Documentation"}),(0,o.jsx)("a",{href:"https://owncast.online/help",target:"_blank",rel:"noreferrer",children:"Contribute"}),(0,o.jsx)("a",{href:"https://github.com/owncast/owncast",target:"_blank",rel:"noreferrer",children:"Source"})]})]})};var w=n(10808),v=n.n(w);let y=e=>{let{content:t}=e,n=(0,i.sJ)(c.hz),s=(0,i.sJ)(c.g1),{version:a}=s;return(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)("div",{className:v().pageContentContainer,children:(0,o.jsx)("div",{className:v().customPageContent,dangerouslySetInnerHTML:{__html:t}})}),n&&(0,o.jsx)(g,{version:a})]})};var b=n(17725),C=n.n(b),k=n(87412),N=n(12341),S=n.n(N);let Z=p()(()=>Promise.all([n.e(1272),n.e(8700),n.e(2852),n.e(4977),n.e(1496)]).then(n.bind(n,94977)).then(e=>e.ChatContainer),{loadableGenerated:{webpack:()=>[94977]}}),F=()=>{let e=(0,i.sJ)(c.db),t=(0,i.sJ)(c.pH);if(!e)return null;let{id:n,isModerator:s,displayName:a}=e;return(0,o.jsx)(k.ZP,{className:S().root,collapsedWidth:0,width:320,children:(0,o.jsx)(Z,{messages:t,usernameToHighlight:a,chatUserId:n,isModerator:s})})};var L=n(12270),T=n.n(L);let A=e=>{let{children:t}=e;return(0,o.jsx)("div",{className:"".concat(T().row),children:t})};var D=n(71577),E=n(13959),M=n.n(E);let I=e=>{let{action:t,primary:n=!0,externalActionSelected:s}=e,{title:i,description:a,icon:l,color:r}=t;return(0,o.jsxs)(D.Z,{type:n?"primary":"default",className:"".concat(M().button),onClick:()=>s(t),style:{backgroundColor:r},children:[l&&(0,o.jsx)("img",{src:l,className:"".concat(M().icon),alt:a}),i]})};var H=n(51513),z=n(34447);let B=e=>{let{onClick:t,props:n}=e;return(0,o.jsx)(D.Z,{...n,type:"primary",className:M().button,icon:(0,o.jsx)(z.Z,{}),onClick:t,children:"Follow"})};var O=n(71578);let J=e=>{let{onClick:t,text:n}=e;return(0,o.jsx)(D.Z,{type:"primary",className:"".concat(M().button),icon:(0,o.jsx)(O.Z,{}),onClick:t,children:n||"Notify"})};var P=n(53731),R=n(79216),U=n(74933),G=n.n(U);let V=e=>{let{src:t}=e;return(0,o.jsx)("div",{className:G().root,children:(0,o.jsx)("div",{className:G().container,children:(0,o.jsx)(R.Z,{src:t,alt:"Logo",className:G().image,rootClassName:G().image})})})};var Y=n(25675),q=n.n(Y),Q=n(573),W=n.n(Q);let K=e=>{let{links:t}=e;return(0,o.jsx)("div",{className:W().links,children:t.map(e=>(0,o.jsx)("a",{href:e.url,className:W().link,target:"_blank",rel:"noreferrer me",children:(0,o.jsx)(q(),{src:e.icon||"/img/platformlogos/default.svg",alt:e.platform,title:e.platform,className:W().link,width:"30",height:"30"})},e.platform))})};var X=n(47900),$=n.n(X);let ee=e=>{let{name:t,title:n,summary:s,logo:i,tags:a,links:l}=e;return(0,o.jsx)("div",{className:$().root,children:(0,o.jsxs)("div",{className:$().logoTitleSection,children:[(0,o.jsx)("div",{className:$().logo,children:(0,o.jsx)(V,{src:i})}),(0,o.jsxs)("div",{className:$().titleSection,children:[(0,o.jsx)("div",{className:m()($().title,$().row,"header-title"),children:t}),(0,o.jsx)("div",{className:m()($().subtitle,$().row,"header-subtitle"),children:(0,o.jsx)(P.Z,{children:n||s})}),(0,o.jsx)("div",{className:m()($().tagList,$().row),children:a.length>0&&a.map(e=>(0,o.jsxs)("span",{children:["#",e,"\xa0"]},e))}),(0,o.jsx)("div",{className:m()($().socialLinks,$().row),children:(0,o.jsx)(K,{links:l})})]})]})})};var et=n(69357),en=n(26303),eo=n(25968),es=n(6226),ei=n(3698),ea=n(24093),el=n(69833),er=n.n(el);let ec=e=>{let{follower:t}=e;return(0,o.jsx)("div",{className:er().follower,children:(0,o.jsx)("a",{href:t.link,target:"_blank",rel:"noreferrer",children:(0,o.jsxs)(eo.Z,{wrap:!1,children:[(0,o.jsx)(es.Z,{span:6,children:(0,o.jsx)(ea.C,{src:t.image,alt:"Avatar",className:er().avatar,children:(0,o.jsx)("img",{src:"/logo",alt:"Logo",className:er().placeholder})})}),(0,o.jsxs)(es.Z,{children:[(0,o.jsx)(eo.Z,{children:t.name}),(0,o.jsx)(eo.Z,{className:er().account,children:t.username})]})]})})})};var ed=n(21890),eh=n.n(ed);let eu=e=>{let{name:t}=e,[n,s]=(0,r.useState)([]),[i,a]=(0,r.useState)(0),[l,c]=(0,r.useState)(1),[d,h]=(0,r.useState)(!0),u=async()=>{try{let e=await fetch("".concat("/api/followers","?page=").concat(l)),t=await e.json(),{results:n,total:o}=t;s(n),a(o)}catch(i){console.error(i)}};(0,r.useEffect)(()=>{u()},[]),(0,r.useEffect)(()=>{u()},[l]),(0,r.useEffect)(()=>{h(!1)},[n]);let m=(0,o.jsxs)("div",{className:eh().noFollowers,children:[(0,o.jsx)("h2",{children:"Be the first follower!"}),(0,o.jsxs)("p",{children:["Owncast"!==t?t:"This server"," is a part of the"," ",(0,o.jsx)("a",{href:"https://owncast.online/join-fediverse",children:"Fediverse"}),", an interconnected network of independent users and servers."]}),(0,o.jsxs)("p",{children:["By following ","Owncast"!==t?t:"this server"," you'll be able to get updates from the stream, share it with others, and and show your appreciation when it goes live, all from your own Fediverse account."]}),(0,o.jsx)(B,{})]}),x=(0,o.jsx)(en.Z,{active:!0,paragraph:{rows:3}});return d?x:(null==n?void 0:n.length)?(0,o.jsxs)("div",{className:eh().followers,children:[(0,o.jsx)(eo.Z,{wrap:!0,gutter:[10,10],children:n.map(e=>(0,o.jsx)(es.Z,{children:(0,o.jsx)(ec,{follower:e},e.link)},e.link))}),(0,o.jsx)(ei.Z,{current:l,pageSize:24,total:Math.ceil(i/24)||1,onChange(e){c(e)},hideOnSinglePage:!0})]}):m};var em=n(14761),ex=n(66516),ep=n(38504),ef=n(49647),ej=n(60198),e_=n(89705),eg=n(97038),ew=n.n(eg);let ev="notify",ey="follow",eb=e=>{let{actions:t,externalActionSelected:n,notifyItemSelected:s,followItemSelected:i,showFollowItem:a,showNotifyItem:l}=e,r=e=>{if(e.key===ev){s();return}if(e.key===ey){i();return}let o=t.find(t=>t.url===e.key);n(o)},c=t.map(e=>({key:e.url,label:(0,o.jsxs)("span",{className:ew().item,children:[e.icon&&(0,o.jsx)("img",{className:ew().icon,src:e.icon,alt:e.title})," ",e.title]})}));a&&c.unshift({key:ey,label:(0,o.jsxs)("span",{className:ew().item,children:[(0,o.jsx)(ef.Z,{className:ew().icon})," Follow this stream"]})}),l&&c.unshift({key:ev,label:(0,o.jsxs)("span",{className:ew().item,children:[(0,o.jsx)(ej.Z,{className:ew().icon}),"Notify when live"]})});let d=(0,o.jsx)(ex.Z,{items:c,onClick:r});return(0,o.jsx)(ep.Z,{overlay:d,placement:"bottomRight",trigger:["click"],className:ew().menu,children:(0,o.jsx)("div",{className:ew().buttonWrap,children:(0,o.jsx)(D.Z,{type:"default",onClick:e=>e.preventDefault(),size:"large",icon:(0,o.jsx)(e_.Z,{size:6,style:{rotate:"90deg"}})})})})};var eC=n(26713),ek=n(14670),eN=n(69677),eS=n(66009),eZ=n.n(eS);let eF=e=>{let{handleClose:t,account:n,name:s}=e,[i,a]=(0,r.useState)(null),[l,c]=(0,r.useState)(!1),[d,u]=(0,r.useState)(!1),[m,x]=(0,r.useState)(null),p=e=>{a(e),function(e){let t=e.replace(/^@+/,"");return/^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(String(t).toLowerCase())}(e)?c(!0):c(!1)},f=()=>{window.open("https://owncast.online/join-fediverse","_blank")},j=async()=>{if(l){u(!0);try{let e=i.replace(/^@+/,""),n=await fetch("/api/remotefollow",{method:"POST",body:JSON.stringify({account:e})}),o=await n.json();if(o.redirectUrl&&(window.open(o.redirectUrl,"_blank"),t()),!o.success){x(o.message),u(!1);return}if(!o.redirectUrl){x("Unable to follow."),u(!1);return}}catch(s){x(s.message)}u(!1)}};return(0,o.jsxs)(eC.Z,{direction:"vertical",children:[(0,o.jsxs)("div",{className:eZ().header,children:["By following this stream you'll get notified on the Fediverse when it goes live. Now is a great time to",(0,o.jsx)("a",{href:"https://owncast.online/join-fediverse",target:"_blank",rel:"noreferrer",children:"learn about the Fediverse"}),"if it's new to you."]}),(0,o.jsxs)(h.Z,{spinning:d,children:[m&&(0,o.jsx)(ek.Z,{message:"Follow Error",description:m,type:"error",showIcon:!0}),(0,o.jsxs)("div",{className:eZ().account,children:[(0,o.jsx)("img",{src:"/logo",alt:"logo",className:eZ().logo}),(0,o.jsxs)("div",{className:eZ().username,children:[(0,o.jsx)("div",{className:eZ().name,children:s}),(0,o.jsx)("div",{children:n})]})]}),(0,o.jsxs)("div",{children:[(0,o.jsx)("div",{className:eZ().instructions,children:"Enter your username @server to follow"}),(0,o.jsx)(eN.Z,{value:i,size:"large",onChange:e=>p(e.target.value),placeholder:"Your fediverse account @account@server",defaultValue:i}),(0,o.jsx)("div",{className:eZ().footer,children:"You'll be redirected to your Fediverse server and asked to confirm the action."})]}),(0,o.jsxs)(eC.Z,{className:eZ().buttons,children:[(0,o.jsx)(D.Z,{disabled:!l,type:"primary",onClick:j,children:"Follow"}),(0,o.jsx)(D.Z,{onClick:f,type:"primary",children:"Join the Fediverse"})]})]})]})},{Content:eL}=s.Z,eT=p()(()=>n.e(6160).then(n.bind(n,66160)).then(e=>e.BrowserNotifyModal),{loadableGenerated:{webpack:()=>[66160]}}),eA=p()(()=>n.e(7815).then(n.bind(n,17815)).then(e=>e.NotifyReminderPopup),{loadableGenerated:{webpack:()=>[17815]}}),eD=p()(()=>Promise.all([n.e(2544),n.e(7902),n.e(2239),n.e(7018)]).then(n.bind(n,8888)).then(e=>e.OwncastPlayer),{loadableGenerated:{webpack:()=>[8888]}}),eE=p()(()=>Promise.all([n.e(1272),n.e(8700),n.e(2852),n.e(4977),n.e(1496)]).then(n.bind(n,94977)).then(e=>e.ChatContainer),{loadableGenerated:{webpack:()=>[94977]}}),eM=e=>{let{name:t,streamTitle:n,summary:s,tags:i,socialHandles:a,extraPageContent:l}=e,r=(0,o.jsx)(y,{content:l}),c=(0,o.jsx)(eu,{name:t});return(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)("div",{className:C().lowerHalf,children:(0,o.jsx)(ee,{name:t,title:n,summary:s,tags:i,links:a,logo:"/logo"})}),(0,o.jsx)("div",{className:C().lowerSection,children:(0,o.jsx)(d.Z,{defaultActiveKey:"0",items:[{label:"About",key:"2",children:r},{label:"Followers",key:"3",children:c}]})})]})},eI=e=>{let{name:t,streamTitle:n,summary:s,tags:i,socialHandles:a,extraPageContent:l,messages:c,currentUser:h,showChat:u,actions:x,setExternalActionToDisplay:p,setShowNotifyPopup:f,setShowFollowModal:j}=e;if(!h)return null;let _=(0,r.useRef)(),{id:g,displayName:w}=h,v=u&&(0,o.jsx)(eE,{messages:c,usernameToHighlight:w,chatUserId:g,isModerator:!1}),b=(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)(ee,{name:t,title:n,summary:s,tags:i,links:a,logo:"/logo"}),(0,o.jsx)(y,{content:l})]}),k=(0,o.jsx)(eu,{name:t}),N="".concat(function(e){let[t,n]=(0,r.useState)(0),o=()=>{if(!e.current)return;let t=e.current.getBoundingClientRect().top,{innerHeight:o}=window;n(o-t)};return(0,r.useEffect)(()=>(o(),window.addEventListener("resize",o),()=>{window.removeEventListener("resize",o)}),[]),t}(_),"px"),S=(e,t)=>(0,o.jsxs)("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"start"},children:[(0,o.jsx)(t,{...e,style:{width:"85%"}}),(0,o.jsx)(eb,{showFollowItem:!0,showNotifyItem:!0,actions:x,externalActionSelected:p,notifyItemSelected:()=>f(!0),followItemSelected:()=>j(!0)})]});return(0,o.jsx)("div",{className:m()(C().lowerSectionMobile),ref:_,style:{height:N},children:(0,o.jsx)(d.Z,{defaultActiveKey:"0",items:[u&&{label:"Chat",key:"0",children:v},{label:"About",key:"2",children:b},{label:"Followers",key:"3",children:k}],renderTabBar:S})})},eH=e=>{let{externalActionToDisplay:t,setExternalActionToDisplay:n}=e,{title:s,description:i,url:a}=t;return(0,o.jsx)(em.Modal,{title:i||s,url:a,open:!!t,height:"80vh",handleCancel:()=>n(null)})},ez=()=>{let e=(0,i.sJ)(c.Q),t=(0,i.sJ)(c.g1),n=(0,i.sJ)(c.pT),s=(0,i.sJ)(c.di),a=(0,i.sJ)(c.db),[l,d]=(0,i.FV)(c.hz),u=(0,i.sJ)(c.j$),m=(0,i.sJ)(c.YW),{viewerCount:x,lastConnectTime:p,lastDisconnectTime:j,streamTitle:_}=(0,i.sJ)(c.RI),{extraPageContent:g,version:w,name:v,summary:y,socialHandles:b,tags:k,externalActions:N,offlineMessage:S,chatDisabled:Z,federation:L,notifications:T}=t,[D,E]=(0,r.useState)(!1),[M,z]=(0,r.useState)(!1),[O,P]=(0,r.useState)(!1),{account:R}=L,{browser:U}=T,{enabled:G}=U,[V,Y]=(0,r.useState)(null),q=e=>{let{openExternally:t,url:n}=e;t?window.open(n,"_blank"):Y(e)},Q=N.map(e=>(0,o.jsx)(I,{action:e,externalActionSelected:q},e.url)),W=()=>{let e=parseInt((0,f.$o)(f.dA.userVisitCount),10);Number.isNaN(e)&&(e=0),(0,f.qQ)(f.dA.userVisitCount,e+1),e>2&&!(0,f.$o)(f.dA.hasDisplayedNotificationModal)&&E(!0)},K=()=>{z(!1),E(!1),(0,f.qQ)(f.dA.hasDisplayedNotificationModal,!0)},X=()=>{let e=window.innerWidth;void 0===l&&(e<=768?d(!0):d(!1)),!l&&e<=768&&d(!0),l&&e>768&&d(!1)};(0,r.useEffect)(()=>(W(),X(),window.addEventListener("resize",X),()=>{window.removeEventListener("resize",X)}),[]);let $=!Z&&s&&n;return(0,o.jsxs)(o.Fragment,{children:[(0,o.jsxs)("div",{className:C().main,children:[(0,o.jsx)(h.Z,{wrapperClassName:C().loadingSpinner,size:"large",spinning:e.appLoading,children:(0,o.jsxs)(eL,{className:C().root,children:[(0,o.jsxs)("div",{className:C().mainSection,children:[(0,o.jsxs)("div",{className:C().topSection,children:[m&&(0,o.jsx)(eD,{source:"/hls/stream.m3u8",online:m}),!m&&!e.appLoading&&(0,o.jsx)(H.R,{streamName:v,customText:S,notificationsEnabled:G,fediverseAccount:R,lastLive:j,onNotifyClick:()=>z(!0),onFollowClick:()=>P(!0)}),m&&(0,o.jsx)(et.X,{online:m,lastConnectTime:p,lastDisconnectTime:j,viewerCount:x})]}),(0,o.jsx)("div",{className:C().midSection,children:(0,o.jsxs)("div",{className:C().buttonsLogoTitleSection,children:[!l&&(0,o.jsxs)(A,{children:[Q,(0,o.jsx)(B,{size:"small",onClick:()=>P(!0)}),(0,o.jsx)(eA,{open:D,notificationClicked:()=>z(!0),notificationClosed:()=>K(),children:(0,o.jsx)(J,{onClick:()=>z(!0)})})]}),(0,o.jsx)(em.Modal,{title:"Notify",open:M,afterClose:()=>K(),handleCancel:()=>K(),children:(0,o.jsx)(eT,{})})]})}),l?(0,o.jsx)(eI,{name:v,streamTitle:_,summary:y,tags:k,socialHandles:b,extraPageContent:g,messages:u,currentUser:a,showChat:$,actions:N,setExternalActionToDisplay:q,setShowNotifyPopup:z,setShowFollowModal:P}):(0,o.jsx)(eM,{name:v,streamTitle:_,summary:y,tags:k,socialHandles:b,extraPageContent:g})]}),$&&!l&&(0,o.jsx)(F,{})]})}),!l&&!1]}),V&&(0,o.jsx)(eH,{externalActionToDisplay:V,setExternalActionToDisplay:Y}),(0,o.jsx)(em.Modal,{title:"Follow ".concat(v),open:O,handleCancel:()=>P(!1),width:"550px",children:(0,o.jsx)(eF,{account:R,name:v,handleClose:()=>P(!1)})})]})};var eB=n(37527),eO=n(85402);let eJ=e=>{let{title:t,message:n}=e;return(0,o.jsx)(eO.Z,{title:t,visible:!0,footer:null,closable:!1,keyboard:!1,width:900,centered:!0,className:"modal",children:(0,o.jsx)("p",{style:{fontSize:"1.3rem"},children:n})})},eP=()=>{let e=(0,i.sJ)(c.j$),t=(0,i.sJ)(c.RI),n=!1,o="",s=()=>{n=!0,o=window.document.title},a=()=>{n=!1,window.document.title=o},l=()=>{window.addEventListener("blur",s),window.addEventListener("focus",a)};return(0,r.useEffect)(()=>(o=window.document.title,l(),()=>{window.removeEventListener("focus",a),window.removeEventListener("blur",s)}),[]),(0,r.useEffect)(()=>{let{online:e}=t;n&&e&&(window.document.title="\uD83D\uDCAC :: ".concat(o))},[e]),(0,r.useEffect)(()=>{if(!n)return;let{online:e}=t;e?window.document.title=" \uD83D\uDFE2 :: ".concat(o):e||(window.document.title=" \uD83D\uDD34 :: ".concat(o))},[c.RI]),null},eR=()=>(0,o.jsx)("script",{id:"server-side-hydration",dangerouslySetInnerHTML:{__html:"\n	window.configHydration = {{.ServerConfigJSON}};\n	window.statusHydration = {{.StatusJSON}};\n	"}}),eU=()=>{let[e]=(0,i.FV)(c.hz),t=(0,i.sJ)(c.g1),{name:n,title:a,customStyles:d,version:h}=t,u=(0,i.sJ)(c.di),m=(0,i.sJ)(c.ap),x=(0,r.useRef)(null),{chatDisabled:p}=t;return(0,r.useEffect)(()=>{!function(e){let t=e=>{let t=e.getAttribute("rel");e.setAttribute("rel","".concat(t," noopener noreferrer"))},n=function(e){for(let n of e)for(let o of n.addedNodes)o instanceof HTMLElement&&"a"===o.tagName.toLowerCase()&&t(o)};e.querySelectorAll("a").forEach(e=>t(e));let o=new MutationObserver(n);o.observe(e,{attributes:!1,childList:!0,subtree:!0})}(x.current)},[]),(0,o.jsxs)(o.Fragment,{children:[(0,o.jsxs)(l(),{children:[(0,o.jsx)(eR,{}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"57x57",href:"/img/favicon/apple-icon-57x57.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"60x60",href:"/img/favicon/apple-icon-60x60.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"72x72",href:"/img/favicon/apple-icon-72x72.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"76x76",href:"/img/favicon/apple-icon-76x76.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"114x114",href:"/img/favicon/apple-icon-114x114.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"120x120",href:"/img/favicon/apple-icon-120x120.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"144x144",href:"/img/favicon/apple-icon-144x144.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"152x152",href:"/img/favicon/apple-icon-152x152.png"}),(0,o.jsx)("link",{rel:"apple-touch-icon",sizes:"180x180",href:"/img/favicon/apple-icon-180x180.png"}),(0,o.jsx)("link",{rel:"icon",type:"image/png",sizes:"192x192",href:"/img/favicon/android-icon-192x192.png"}),(0,o.jsx)("link",{rel:"icon",type:"image/png",sizes:"32x32",href:"/img/favicon/favicon-32x32.png"}),(0,o.jsx)("link",{rel:"icon",type:"image/png",sizes:"96x96",href:"/img/favicon/favicon-96x96.png"}),(0,o.jsx)("link",{rel:"icon",type:"image/png",sizes:"16x16",href:"/img/favicon/favicon-16x16.png"}),(0,o.jsx)("link",{rel:"manifest",href:"/manifest.json"}),(0,o.jsx)("link",{href:"/api/auth/provider/indieauth"}),(0,o.jsx)("meta",{name:"msapplication-TileColor",content:"#ffffff"}),(0,o.jsx)("meta",{name:"msapplication-TileImage",content:"/img/favicon/ms-icon-144x144.png"}),(0,o.jsx)("meta",{name:"theme-color",content:"#ffffff"}),(0,o.jsx)("style",{children:d}),(0,o.jsx)("base",{target:"_blank"})]}),(0,o.jsxs)(l(),{children:[n?(0,o.jsx)("title",{children:n}):(0,o.jsx)("title",{children:"{{.Name}}"}),(0,o.jsx)("meta",{name:"description",content:"{{.Summary}}"}),(0,o.jsx)("meta",{property:"og:title",content:"{{.Name}}"}),(0,o.jsx)("meta",{property:"og:site_name",content:"{{.Name}}"}),(0,o.jsx)("meta",{property:"og:url",content:"{{.RequestedURL}}"}),(0,o.jsx)("meta",{property:"og:description",content:"{{.Summary}}"}),(0,o.jsx)("meta",{property:"og:type",content:"video.other"}),(0,o.jsx)("meta",{property:"video:tag",content:"{{.TagsString}}"}),(0,o.jsx)("meta",{property:"og:image",content:"{{.Thumbnail}}"}),(0,o.jsx)("meta",{property:"og:image:url",content:"{{.Thumbnail}}"}),(0,o.jsx)("meta",{property:"og:image:alt",content:"{{.Image}}"}),(0,o.jsx)("meta",{property:"og:video",content:"{{.RequestedURL}}/embed/video"}),(0,o.jsx)("meta",{property:"og:video:secure_url",content:"{{.RequestedURL}}/embed/video"}),(0,o.jsx)("meta",{property:"og:video:height",content:"315"}),(0,o.jsx)("meta",{property:"og:video:width",content:"560"}),(0,o.jsx)("meta",{property:"og:video:type",content:"text/html"}),(0,o.jsx)("meta",{property:"og:video:actor",content:"{{.Name}}"}),(0,o.jsx)("meta",{property:"twitter:title",content:"{{.Name}}"}),(0,o.jsx)("meta",{property:"twitter:url",content:"{{.RequestedURL}}"}),(0,o.jsx)("meta",{property:"twitter:description",content:"{{.Summary}}"}),(0,o.jsx)("meta",{property:"twitter:image",content:"{{.Image}}"}),(0,o.jsx)("meta",{property:"twitter:card",content:"player"}),(0,o.jsx)("meta",{property:"twitter:player",content:"{{.RequestedURL}}/embed/video"}),(0,o.jsx)("meta",{property:"twitter:player:width",content:"560"}),(0,o.jsx)("meta",{property:"twitter:player:height",content:"315"})]}),(0,o.jsx)(c.me,{}),(0,o.jsx)(eP,{}),(0,o.jsxs)(s.Z,{ref:x,style:{minHeight:"100vh"},children:[(0,o.jsx)(eB.h,{name:a||n,chatAvailable:u,chatDisabled:p}),(0,o.jsx)(ez,{}),m&&(0,o.jsx)(eJ,{title:m.title,message:m.message}),!e&&(0,o.jsx)(g,{version:h})]})]})};function eG(){return(0,o.jsx)(eU,{})}},13959:function(e){e.exports={button:"ActionButton_button__z5Z2c",icon:"ActionButton_icon__EPp7Q"}},97038:function(e){e.exports={item:"ActionButtonMenu_item__OJQdr",buttonWrap:"ActionButtonMenu_buttonWrap__WQ9kt",icon:"ActionButtonMenu_icon__edY1D",menu:"ActionButtonMenu_menu__GChDk"}},12270:function(e){e.exports={row:"ActionButtonRow_row__SiEGe"}},47900:function(e){e.exports={root:"ContentHeader_root__HaUG0",row:"ContentHeader_row__9Q8gH",logoTitleSection:"ContentHeader_logoTitleSection__Z8pUc",logo:"ContentHeader_logo__wo_HN",titleSection:"ContentHeader_titleSection___6Y15",title:"ContentHeader_title__E_DsI",subtitle:"ContentHeader_subtitle__n1Wew",tagList:"ContentHeader_tagList__rx3jY"}},79252:function(e){e.exports={root:"UserDropdown_root__IdxfQ","ant-space":"UserDropdown_ant-space__XJTZ3","ant-space-item":"UserDropdown_ant-space-item__w4nC2"}},66009:function(e){e.exports={header:"FollowModal_header__la1ji",buttons:"FollowModal_buttons__tt4Mc",instructions:"FollowModal_instructions__HiKFF",footer:"FollowModal_footer__AjucH",account:"FollowModal_account__cmHkm",logo:"FollowModal_logo__Ew8xK",username:"FollowModal_username__A_OTw",name:"FollowModal_name__Sf_TP"}},17725:function(e){e.exports={root:"Content_root__h1mNK",mainSection:"Content_mainSection__Gk78Y",topSection:"Content_topSection__JIZi0",lowerSection:"Content_lowerSection__BZHYI",lowerSectionMobile:"Content_lowerSectionMobile__hRr0_",leftCol:"Content_leftCol__U2TDq",loadingSpinner:"Content_loadingSpinner__mDlYC",main:"Content_main__XVf63"}},10808:function(e){e.exports={pageContentContainer:"CustomPageContent_pageContentContainer__EG4tU",customPageContent:"CustomPageContent_customPageContent__Mr981",summary:"CustomPageContent_summary___Zgps"}},67917:function(e){e.exports={footer:"Footer_footer__mPuvf",links:"Footer_links__7bBxV"}},31764:function(e){e.exports={header:"Header_header__U4Ro1",logo:"Header_logo__HLZ6Z"}},74933:function(e){e.exports={root:"Logo_root__jKiJC",container:"Logo_container__A4UYT",image:"Logo_image__Ahkom"}},77011:function(e){e.exports={spinner:"Modal_spinner__GiSS0",content:"Modal_content__h9my9"}},88335:function(e){e.exports={outerContainer:"OfflineBanner_outerContainer__3AbsB",innerContainer:"OfflineBanner_innerContainer__zTm13",bodyText:"OfflineBanner_bodyText__nNNy0",separator:"OfflineBanner_separator___j_Ss",lastLiveDate:"OfflineBanner_lastLiveDate___UZdO",clockIcon:"OfflineBanner_clockIcon__s0DB_",header:"OfflineBanner_header__Vu20o",footer:"OfflineBanner_footer__o3Zl5",actionLink:"OfflineBanner_actionLink__b4Mwa"}},12341:function(e){e.exports={root:"Sidebar_root__8HE0A"}},573:function(e){e.exports={link:"SocialLinks_link___CcSm",links:"SocialLinks_links__gOAb7"}},37970:function(e){e.exports={statusbar:"Statusbar_statusbar__AtVnB"}},21890:function(e){e.exports={followers:"FollowerCollection_followers__e_EUS",noFollowers:"FollowerCollection_noFollowers__UaCVW"}},69833:function(e){e.exports={follower:"SingleFollower_follower__EyBDI",avatar:"SingleFollower_avatar__V9jHG",account:"SingleFollower_account__Z66vo",placeholder:"SingleFollower_placeholder__CgsfJ"}}},function(e){e.O(0,[8939,3903,4267,4381,5938,6395,5360,250,2858,8980,7466,9774,2888,179],function(){return e(e.s=48312)}),_N_E=e.O()}]);