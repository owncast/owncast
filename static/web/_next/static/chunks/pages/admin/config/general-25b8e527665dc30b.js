(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[1871],{39856:function(e,t,a){(window.__NEXT_P=window.__NEXT_P||[]).push(["/admin/config/general",function(){return a(60060)}])},4349:function(e,t,a){"use strict";a.d(t,{Q:function(){return d},Y:function(){return u}});var s=a(85893),n=a(67294),l=a(53740),i=a(59361),r=a(67032),o=a(57520);let{Title:c}=l.default,d="#5a67d8",u=e=>{let{title:t,description:a,placeholder:l,maxLength:u,values:h,handleDeleteIndex:m,handleCreateString:p,submitStatus:x,continuousStatusMessage:f}=e,[g,j]=(0,n.useState)("");return(0,s.jsxs)("div",{className:"edit-string-array-container",children:[(0,s.jsx)(c,{level:3,className:"section-title",children:t}),(0,s.jsx)("p",{className:"description",children:a}),(0,s.jsx)("div",{className:"edit-current-strings",children:null==h?void 0:h.map((e,t)=>(0,s.jsx)(i.Z,{closable:!0,onClose:()=>{m(t)},color:d,children:e},"tag-".concat(e,"-").concat(t)))}),f&&(0,s.jsx)("div",{className:"continuous-status-section",children:(0,s.jsx)(o.E,{status:f})}),(0,s.jsx)("div",{className:"add-new-string-section",children:(0,s.jsx)(r.nv,{fieldName:"string-input",value:g,onChange:e=>{let{value:t}=e;j(t)},onPressEnter:()=>{let e=g.trim();p(e),j("")},maxLength:u,placeholder:l,status:x})})]})};u.defaultProps={maxLength:50,description:null,submitStatus:null,continuousStatusMessage:null}},44716:function(e,t,a){"use strict";a.d(t,{Z:function(){return d}});var s=a(85893),n=a(67294),l=a(40987),i=a(91332),r=a(57520),o=a(24044),c=a(38631);let d=e=>{let{apiPath:t,checked:a,reversed:d=!1,configPath:u="",disabled:h=!1,fieldName:m,label:p,tip:x,useSubmit:f,onChange:g}=e,[j,v]=(0,n.useState)(null),y=null,b=(0,n.useContext)(c.a),{setFieldInConfigState:k}=b||{},N=()=>{v(null),clearTimeout(y),y=null},w=async e=>{if(f){v((0,i.kg)(i.Jk));let a=d?!e:e;await (0,o.Si)({apiPath:t,data:{value:a},onSuccess:()=>{k({fieldName:m,value:a,path:u}),v((0,i.kg)(i.zv))},onError:e=>{v((0,i.kg)(i.Un,"There was an error: ".concat(e)))}}),y=setTimeout(N,o.sI)}g&&g(e)},S=null!==j&&j.type===i.Jk;return(0,s.jsxs)("div",{className:"formfield-container toggleswitch-container",children:[p&&(0,s.jsx)("div",{className:"label-side",children:(0,s.jsx)("span",{className:"formfield-label",children:p})}),(0,s.jsxs)("div",{className:"input-side",children:[(0,s.jsxs)("div",{className:"input-group",children:[(0,s.jsx)(l.Z,{className:"switch field-".concat(m),loading:S,onChange:w,defaultChecked:a,checked:a,checkedChildren:"ON",unCheckedChildren:"OFF",disabled:h}),(0,s.jsx)(r.E,{status:j})]}),(0,s.jsx)("p",{className:"field-tip",children:x})]})]})};d.defaultProps={apiPath:"",checked:!1,reversed:!1,configPath:"",disabled:!1,label:"",tip:"",useSubmit:!1,onChange:null}},97469:function(e,t,a){"use strict";a.d(t,{Z:function(){return S}});var s=a(85893),n=a(67294),l=a(53740),i=a(2307),r=a(65400),o=a(56697),c=a(51024),d=a(5152),u=a.n(d),h=a(64749),m=a(24044);let p=e=>{let{iconList:t,selectedOption:a,onSelected:n}=e,l=""===a?null:a;return(0,s.jsxs)("div",{className:"social-dropdown-container",children:[(0,s.jsx)("p",{className:"description",children:"If you are looking for a platform name not on this list, please select Other and type in your own name. A logo will not be provided."}),(0,s.jsxs)("div",{className:"formfield-container",children:[(0,s.jsx)("div",{className:"label-side",children:(0,s.jsx)("span",{className:"formfield-label",children:"Social Platform"})}),(0,s.jsx)("div",{className:"input-side",children:(0,s.jsxs)(h.default,{style:{width:240},className:"social-dropdown",placeholder:"Social platform...",defaultValue:l,value:l,onSelect:e=>{n&&n(e)},children:[t.map(e=>{let{platform:t,icon:a,key:n}=e;return(0,s.jsxs)(h.default.Option,{className:"social-option",value:n,children:[(0,s.jsx)("span",{className:"option-icon",children:(0,s.jsx)("img",{src:a,alt:"",className:"option-icon"})}),(0,s.jsx)("span",{className:"option-label",children:t})]},"platform-".concat(n))}),(0,s.jsx)(h.default.Option,{className:"social-option",value:m.z_,children:"Other..."},"platform-".concat(m.z_))]})})]})]})};var x=a(81453),f=a(38631),g=a(53899),j=a(67032),v=a(91332),y=a(57520);let{Title:b}=l.default,k=u()(()=>Promise.resolve().then(a.t.bind(a,18244,23)),{loadableGenerated:{webpack:()=>[18244]},ssr:!1}),N=u()(()=>Promise.resolve().then(a.t.bind(a,22802,23)),{loadableGenerated:{webpack:()=>[22802]},ssr:!1}),w=u()(()=>a.e(7949).then(a.t.bind(a,77949,23)),{loadableGenerated:{webpack:()=>[77949]},ssr:!1});function S(){var e,t;let[a,l]=(0,n.useState)([]),[d,u]=(0,n.useState)([]),[h,S]=(0,n.useState)(!1),[C,T]=(0,n.useState)(!1),[E,P]=(0,n.useState)(!1),[O,Z]=(0,n.useState)(-1),[_,z]=(0,n.useState)(m.wC),[L,I]=(0,n.useState)(null),U=(0,n.useContext)(f.a),{serverConfig:B,setFieldInConfigState:F}=U||{},{instanceDetails:A}=B,{socialHandles:D}=A,M=async()=>{try{let e=await (0,x.rQ)(x.$i,{auth:!1}),t=Object.keys(e).map(t=>({key:t,...e[t]}));l(t)}catch(e){console.log(e)}},R=e=>a.find(t=>t.key===e)||!1,J=""!==_.platform&&!a.find(e=>e.key===_.platform);(0,n.useEffect)(()=>{M()},[]),(0,n.useEffect)(()=>{A.socialHandles&&u(D)},[A]);let V=()=>{I(null),clearTimeout(null)},G=()=>{S(!1),Z(-1),T(!1),P(!1),z({...m.wC})},H=()=>{G()},W=(e,t)=>{z({..._,[e]:t})},Y=async e=>{await (0,m.Si)({apiPath:m.c9,data:{value:e},onSuccess:()=>{F({fieldName:"socialHandles",value:e,path:"instanceDetails"}),P(!1),H(),I((0,v.kg)(v.zv)),setTimeout(V,m.sI)},onError:e=>{I((0,v.kg)(v.Un,"There was an error: ".concat(e))),P(!1),setTimeout(V,m.sI)}})},K=e=>{let t=[...d];t.splice(e,1),Y(t)},$=e=>{if(e<=0||e>=d.length)return;let t=[...d],a=t[e-1];t[e-1]=t[e],t[e]=a,Y(t)},X=e=>{if(e<0||e>=d.length-1)return;let t=[...d],a=t[e+1];t[e+1]=t[e],t[e]=a,Y(t)},Q=[{title:"Social Link",dataIndex:"",key:"combo",render:(e,t)=>{let{platform:a,url:n}=t,l=R(a);if(!l)return(0,s.jsx)("div",{className:"social-handle-cell",children:(0,s.jsxs)("p",{className:"option-label",children:[(0,s.jsx)("strong",{children:a}),(0,s.jsx)("span",{className:"handle-url",title:n,children:n})]})});let{icon:i,platform:r}=l;return(0,s.jsxs)("div",{className:"social-handle-cell",children:[(0,s.jsx)("span",{className:"option-icon",children:(0,s.jsx)("img",{src:i,alt:"",className:"option-icon"})}),(0,s.jsxs)("p",{className:"option-label",children:[(0,s.jsx)("strong",{children:r}),(0,s.jsx)("span",{className:"handle-url",title:n,children:n})]})]})}},{title:"",dataIndex:"",key:"edit",render:(e,t,a)=>(0,s.jsxs)("div",{className:"actions",children:[(0,s.jsx)(r.default,{size:"small",onClick:()=>{let e=d[a];Z(a),z({...e}),S(!0),R(e.platform)||T(!0)},children:"Edit"}),(0,s.jsx)(r.default,{icon:(0,s.jsx)(N,{}),size:"small",hidden:0===a,onClick:()=>$(a)}),(0,s.jsx)(r.default,{icon:(0,s.jsx)(k,{}),size:"small",hidden:a===d.length-1,onClick:()=>X(a)}),(0,s.jsx)(r.default,{className:"delete-button",icon:(0,s.jsx)(w,{}),size:"small",onClick:()=>K(a)})]})}],q={disabled:(e=_.url,"xmpp"===(t=_.platform)?!(0,g.Kf)(e,"xmpp"):"matrix"===t?!(0,g.bu)(e):!(0,g.jv)(e))},ee=(0,s.jsxs)("div",{className:"other-field-container formfield-container",children:[(0,s.jsx)("div",{className:"label-side"}),(0,s.jsx)("div",{className:"input-side",children:(0,s.jsx)(c.default,{placeholder:"Other platform name",defaultValue:_.platform,onChange:e=>{let{value:t}=e.target;W("platform",t)}})})]});return(0,s.jsxs)("div",{className:"social-links-edit-container",children:[(0,s.jsx)(b,{level:3,className:"section-title",children:"Your Social Handles"}),(0,s.jsx)("p",{className:"description",children:"Add all your social media handles and links to your other profiles here."}),(0,s.jsx)(y.E,{status:L}),(0,s.jsx)(i.Z,{className:"social-handles-table",pagination:!1,size:"small",rowKey:e=>"".concat(e.platform,"-").concat(e.url),columns:Q,dataSource:d}),(0,s.jsx)(o.default,{title:"Edit Social Handle",open:h,onOk:()=>{P(!0);let e=d.length?[...d]:[];-1===O?e.push(_):e.splice(O,1,_),Y(e)},onCancel:H,confirmLoading:E,okButtonProps:q,children:(0,s.jsxs)("div",{className:"social-handle-modal-content",children:[(0,s.jsx)(p,{iconList:a,selectedOption:J?m.z_:_.platform,onSelected:e=>{e===m.z_?(T(!0),W("platform","")):(T(!1),W("platform",e))}}),C&&ee,(0,s.jsx)("br",{}),(0,s.jsx)(j.nv,{fieldName:"social-url",label:"URL",placeholder:{mastodon:"https://mastodon.social/@username",twitter:"https://twitter.com/username"}[_.platform]||"Url to page",value:_.url,onChange:e=>{let{value:t}=e;W("url",t)},useTrim:!0,type:"url",pattern:g.ax}),(0,s.jsx)(y.E,{status:L})]})}),(0,s.jsx)("br",{}),(0,s.jsx)(r.default,{type:"primary",onClick:()=>{G(),S(!0)},children:"Add a new social link"})]})}},60060:function(e,t,a){"use strict";a.r(t),a.d(t,{default:function(){return ed}});var s=a(85893),n=a(67294),l=a(20838),i=a(65400),r=a(56697),o=a(53740),c=a(48825),d=a(95089),u=a(58909),h=a(76538),m=a(78353),p=a(38631),x=a(24044),f=a(44716),g=a(28465),j=a(5152),v=a.n(j),y=a(57520),b=a(91332),k=a(81453),N=a(80693);let w=v()(()=>Promise.resolve().then(a.t.bind(a,628,23)),{loadableGenerated:{webpack:()=>[628]},ssr:!1}),S=v()(()=>a.e(9876).then(a.t.bind(a,29876,23)),{loadableGenerated:{webpack:()=>[29876]},ssr:!1}),C=()=>{var e;let[t,a]=(0,n.useState)(null),[l,r]=(0,n.useState)(!1),[o,c]=(0,n.useState)(0),d=(0,n.useContext)(p.a),{setFieldInConfigState:u,serverConfig:h}=d||{},m=null==h?void 0:null===(e=h.instanceDetails)||void 0===e?void 0:e.logo,[f,j]=(0,n.useState)(null),v=null,{apiPath:C,tip:T}=x.TEXTFIELD_PROPS_LOGO,E=()=>{j(null),clearTimeout(v),v=null},P=async()=>{t!==m&&(j((0,b.kg)(b.Jk)),await (0,x.Si)({apiPath:C,data:{value:t},onSuccess:()=>{u({fieldName:"logo",value:t,path:""}),j((0,b.kg)(b.zv)),r(!1),c(Math.floor(100*Math.random()))},onError:e=>{j((0,b.kg)(b.Un,"There was an error: ".concat(e))),r(!1)}}),v=setTimeout(E,x.sI))},O="".concat(k.WB,"logo?random=").concat(o);return(0,s.jsxs)("div",{className:"formfield-container logo-upload-container",children:[(0,s.jsx)("div",{className:"label-side",children:(0,s.jsx)("span",{className:"formfield-label",children:"Logo"})}),(0,s.jsxs)("div",{className:"input-side",children:[(0,s.jsxs)("div",{className:"input-group",children:[(0,s.jsx)("img",{src:O,alt:"avatar",className:"logo-preview"}),(0,s.jsx)(g.Z,{name:"logo",listType:"picture",className:"avatar-uploader",showUploadList:!1,accept:N.dr.join(","),beforeUpload:e=>(r(!0),new Promise((t,s)=>{if(e.size>N.Z7){let t="File size is too big: ".concat((0,N.kR)(e.size));return j((0,b.kg)(b.Un,"There was an error: ".concat(t))),v=setTimeout(E,x.sI),r(!1),s()}if(!N.dr.includes(e.type)){let t="File type is not supported: ".concat(e.type);return j((0,b.kg)(b.Un,"There was an error: ".concat(t))),v=setTimeout(E,x.sI),r(!1),s()}(0,N.y3)(e,e=>{a(e),setTimeout(()=>t(),100)})})),customRequest:P,disabled:l,children:l?(0,s.jsx)(w,{style:{color:"white"}}):(0,s.jsx)(i.default,{icon:(0,s.jsx)(S,{})})})]}),(0,s.jsx)(y.E,{status:f}),(0,s.jsx)("p",{className:"field-tip",children:T})]})]})},{Title:T}=o.default,E=e=>{let{cancelPressed:t,okPressed:a}=e;return(0,s.jsxs)(r.default,{width:"70%",title:"Owncast Directory",visible:!0,onCancel:t,footer:(0,s.jsxs)("div",{children:[(0,s.jsx)(i.default,{onClick:t,children:"Do not share my server."}),(0,s.jsx)(i.default,{type:"primary",onClick:a,children:"Understood. Share my server publicly."})]}),children:[(0,s.jsx)(o.default.Title,{level:3,children:"What is the Owncast Directory?"}),(0,s.jsxs)(o.default.Paragraph,{children:["Owncast operates a public directory at"," ",(0,s.jsx)("a",{href:"https://directory.owncast.online",children:"directory.owncast.online"})," to share your video streams with more people, while also using these as examples for others. Live streams and servers listed on the directory may optionally be shared on other platforms and applications."]}),(0,s.jsx)(o.default.Title,{level:3,children:"Disclaimers and Responsibility"}),(0,s.jsx)(o.default.Paragraph,{children:(0,s.jsxs)("ul",{children:[(0,s.jsx)("li",{children:"By enabling this feature you are granting explicit rights to Owncast to share your stream to the public via the directory, as well as other sites, applications and any platform where the Owncast project may be promoting Owncast-powered streams including social media."}),(0,s.jsx)("li",{children:"There is no obligation to list any specific server or topic. Servers can and will be removed at any time for any reason."}),(0,s.jsx)("li",{children:"Any server that is streaming Not Safe For Work (NSFW) content and does not have the NSFW toggle enabled on their server will be removed."}),(0,s.jsx)("li",{children:"Any server streaming harmful, hurtful, misleading or hateful content in any way will not be listed."}),(0,s.jsx)("li",{children:"You may reach out to the Owncast team to report any objectionable content or content that you believe should not be be publicly listed."}),(0,s.jsx)("li",{children:"You have the right to free software and to build whatever you want with it. But there is no obligation for others to share it."})]})})]})};function P(){let[e,t]=(0,n.useState)(null),a=(0,n.useContext)(p.a),{serverConfig:l}=a||{},{instanceDetails:r,yp:o,hideViewerCount:g,disableSearchIndexing:j}=l,{instanceUrl:v}=o,[k,N]=(0,n.useState)(null),[w,S]=(0,n.useState)(!1);if((0,n.useEffect)(()=>{t({...r,...o,hideViewerCount:g,disableSearchIndexing:j})},[r,o]),!e)return null;let P=a=>{a&&S(!0),t({...e,yp:{enabled:a}})},O=a=>{let{fieldName:s,value:n}=a;t({...e,[s]:n})},Z=""!==v;return(0,s.jsxs)("div",{className:"edit-general-settings",children:[(0,s.jsx)(T,{level:3,className:"section-title",children:"Configure Instance Details"}),(0,s.jsx)("br",{}),(0,s.jsx)(m.$7,{fieldName:"name",...x.RE,value:e.name,initialValue:r.name,onChange:O}),(0,s.jsx)(m.$7,{fieldName:"instanceUrl",...x.cj,value:e.instanceUrl,initialValue:o.instanceUrl,type:m.xA,onChange:O,onSubmit:()=>{""===e.instanceUrl&&!0===o.enabled&&(0,x.Si)({apiPath:x.AP,data:{value:!1}})}}),(0,s.jsx)(m.$7,{fieldName:"summary",...x.rs,type:m.Sk,value:e.summary,initialValue:r.summary,onChange:O}),(0,s.jsxs)("div",{style:{marginBottom:"50px",marginRight:"150px"},children:[(0,s.jsxs)("div",{style:{display:"flex",width:"80vh",justifyContent:"space-between",alignItems:"end"},children:[(0,s.jsx)("p",{style:{margin:"20px",marginRight:"10px",fontWeight:"400"},children:"Offline Message:"}),(0,s.jsx)(d.ZP,{value:e.offlineMessage,...x.rd,theme:u.FZ,height:"150px",width:"450px",onChange:e=>{O({fieldName:"offlineMessage",value:e})},extensions:[(0,c.markdown)({base:c.markdownLanguage,codeLanguages:h.M})]})]}),(0,s.jsx)("div",{className:"field-tip",children:"The offline message is displayed to your page visitors when you're not streaming. Markdown is supported."}),(0,s.jsx)(i.default,{type:"primary",onClick:()=>{(0,x.Si)({apiPath:x.bJ,data:{value:e.offlineMessage}}),N((0,b.kg)(b.zv)),setTimeout(()=>{N(null)},2e3)},style:{margin:"10px",float:"right"},children:"Save Message"}),(0,s.jsx)(y.Z,{status:k})]}),(0,s.jsx)(C,{}),(0,s.jsx)(f.Z,{fieldName:"hideViewerCount",useSubmit:!0,...x._X,checked:e.hideViewerCount,onChange:function(e){O({fieldName:"hideViewerCount",value:e})}}),(0,s.jsx)(f.Z,{fieldName:"disableSearchIndexing",useSubmit:!0,...x.il,checked:e.disableSearchIndexing,onChange:function(e){O({fieldName:"disableSearchIndexing",value:e})}}),(0,s.jsx)("br",{}),(0,s.jsxs)("p",{className:"description",children:["Increase your audience by appearing in the"," ",(0,s.jsx)("a",{href:"https://directory.owncast.online",target:"_blank",rel:"noreferrer",children:(0,s.jsx)("strong",{children:"Owncast Directory"})}),". This is an external service run by the Owncast project."," ",(0,s.jsx)("a",{href:"https://owncast.online/docs/directory/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn more"}),"."]}),!o.instanceUrl&&(0,s.jsxs)("p",{className:"description",children:["You must set your ",(0,s.jsx)("strong",{children:"Server URL"})," above to enable the directory."]}),(0,s.jsxs)("div",{className:"config-yp-container",children:[(0,s.jsx)(f.Z,{fieldName:"enabled",useSubmit:!0,...x.P,checked:e.enabled,disabled:!Z,onChange:P}),(0,s.jsx)(f.Z,{fieldName:"nsfw",useSubmit:!0,...x.EY,checked:e.nsfw,disabled:!Z})]}),w&&(0,s.jsx)(E,{cancelPressed:function(){S(!1),P(!1),O({fieldName:"enabled",value:!1})},okPressed:function(){S(!1),O({fieldName:"enabled",value:!0}),t({...e,yp:{enabled:!0}})}})]})}var O=a(59361),Z=a(67032),_=a(4349);let{Title:z}=o.default;function L(){let[e,t]=(0,n.useState)(""),[a,l]=(0,n.useState)(null),i=(0,n.useContext)(p.a),{serverConfig:r,setFieldInConfigState:o}=i||{},{instanceDetails:c}=r,{tags:d=[]}=c,{apiPath:u,maxLength:h,placeholder:m,configPath:f}=x.FIELD_PROPS_TAGS,g=null;(0,n.useEffect)(()=>()=>{clearTimeout(g)},[]);let j=()=>{l(null),clearTimeout(g=null)},v=async e=>{l((0,b.kg)(b.Jk)),await (0,x.Si)({apiPath:u,data:{value:e},onSuccess:()=>{o({fieldName:"tags",value:e,path:f}),l((0,b.kg)(b.zv,"Tags updated.")),t(""),g=setTimeout(j,x.sI)},onError:e=>{l((0,b.kg)(b.Un,e)),g=setTimeout(j,x.sI)}})},y=e=>{j();let t=[...d];t.splice(e,1),v(t)};return(0,s.jsxs)("div",{className:"tag-editor-container",children:[(0,s.jsx)(z,{level:3,className:"section-title",children:"Add Tags"}),(0,s.jsx)("p",{className:"description",children:"This is a great way to categorize your Owncast server on the Directory!"}),(0,s.jsx)("div",{className:"edit-current-strings",children:d.map((e,t)=>(0,s.jsx)(O.Z,{closable:!0,onClose:()=>{y(t)},color:_.Q,children:e},"tag-".concat(e,"-").concat(t)))}),(0,s.jsx)("div",{className:"add-new-string-section",children:(0,s.jsx)(Z.nv,{fieldName:"tag-input",value:e,className:"new-tag-input",onChange:e=>{let{value:s}=e;a||l(null),t(s)},onPressEnter:()=>{j();let t=e.trim();if(""===t){l((0,b.kg)(b.dG,"Please enter a tag"));return}if(d.some(e=>e.toLowerCase()===t.toLowerCase())){l((0,b.kg)(b.dG,"This tag is already used!"));return}let a=[...d,t];v(a)},maxLength:h,placeholder:m,status:a})})]})}var I=a(97469);let{Title:U}=o.default;function B(){let[e,t]=(0,n.useState)(""),[a,l]=(0,n.useState)(null),[r,o]=(0,n.useState)(!1),m=(0,n.useContext)(p.a),{serverConfig:f,setFieldInConfigState:g}=m||{},{instanceDetails:j}=f,{extraPageContent:v}=j,k=null,N=()=>{l(null),o(!1),clearTimeout(k),k=null};async function w(){l((0,b.kg)(b.Jk)),await (0,x.Si)({apiPath:x.AA,data:{value:e},onSuccess:t=>{g({fieldName:"extraPageContent",value:e,path:"instanceDetails"}),l((0,b.kg)(b.zv,t))},onError:e=>{l((0,b.kg)(b.Un,e))}}),k=setTimeout(N,x.sI)}return(0,n.useEffect)(()=>{t(v)},[j]),(0,s.jsxs)("div",{className:"edit-page-content",children:[(0,s.jsx)(U,{level:3,className:"section-title",children:"Custom Page Content"}),(0,s.jsxs)("p",{className:"description",children:["Edit the content of your page by using simple"," ",(0,s.jsx)("a",{href:"https://www.markdownguide.org/basic-syntax/",target:"_blank",rel:"noopener noreferrer",children:"Markdown syntax"}),"."]}),(0,s.jsx)(d.ZP,{value:e,placeholder:"Enter your custom page content here...",theme:u.FZ,height:"200px",onChange:function(e){t(e),e===v||r?e===v&&r&&o(!1):o(!0)},extensions:[(0,c.markdown)({base:c.markdownLanguage,codeLanguages:h.M})]}),(0,s.jsx)("br",{}),(0,s.jsxs)("div",{className:"page-content-actions",children:[r&&(0,s.jsx)(i.default,{type:"primary",onClick:w,children:"Save"}),(0,s.jsx)(y.E,{status:a})]})]})}function F(){return(0,s.jsxs)("div",{className:"config-public-details-page",children:[(0,s.jsxs)("p",{className:"description",children:["The following are displayed on your site to describe your stream and its content."," ",(0,s.jsx)("a",{href:"https://owncast.online/docs/website/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn more."})]}),(0,s.jsxs)("div",{className:"top-container",children:[(0,s.jsx)("div",{className:"form-module instance-details-container",children:(0,s.jsx)(P,{})}),(0,s.jsxs)("div",{className:"form-module social-items-container ",children:[(0,s.jsx)("div",{className:"form-module tags-module",children:(0,s.jsx)(L,{})}),(0,s.jsx)("div",{className:"form-module social-handles-container",children:(0,s.jsx)(I.Z,{})})]})]}),(0,s.jsx)("div",{className:"form-module page-content-module",children:(0,s.jsx)(B,{})})]})}var A=a(5789),D=a(68469),M=a(55673),R=a(36155),J=a(74048),V=a(21987),G=a(34528),H=a(48120);let{Title:W}=o.default,Y=()=>{let[e,t]=(0,n.useState)("/* Enter custom CSS here */"),[a,l]=(0,n.useState)(null),[r,o]=(0,n.useState)(!1),c=(0,n.useContext)(p.a),{serverConfig:h,setFieldInConfigState:m}=c||{},{instanceDetails:f}=h,{customStyles:g}=f,j=null,v=()=>{l(null),o(!1),clearTimeout(j),j=null};async function k(){l((0,b.kg)(b.Jk)),await (0,x.Si)({apiPath:x.d$,data:{value:e},onSuccess:t=>{m({fieldName:"customStyles",value:e,path:"instanceDetails"}),l((0,b.kg)(b.zv,t))},onError:e=>{l((0,b.kg)(b.Un,e))}}),j=setTimeout(v,x.sI)}(0,n.useEffect)(()=>{t(g)},[f]);let N=n.useCallback(e=>{t(e),e===g||r?e===g&&r&&o(!1):o(!0)},[]);return(0,s.jsxs)("div",{className:"edit-custom-css",children:[(0,s.jsx)(W,{level:3,className:"section-title",children:"Customize your page styling with CSS"}),(0,s.jsxs)("p",{className:"description",children:["Customize the look and feel of your Owncast instance by overriding the CSS styles of various components on the page. Refer to the"," ",(0,s.jsx)("a",{href:"https://owncast.online/docs/website/",rel:"noopener noreferrer",target:"_blank",children:"CSS & Components guide"}),"."]}),(0,s.jsx)("p",{className:"description",children:"Please input plain CSS text, as this will be directly injected onto your page during load."}),(0,s.jsx)(d.ZP,{value:e,placeholder:"/* Enter custom CSS here */",theme:u.FZ,height:"200px",extensions:[(0,H.css)()],onChange:N}),(0,s.jsx)("br",{}),(0,s.jsxs)("div",{className:"page-content-actions",children:[r&&(0,s.jsx)(i.default,{type:"primary",onClick:k,children:"Save"}),(0,s.jsx)(y.E,{status:a})]})]})};var K=a(8118),$=a.n(K);let{Panel:X}=D.default,Q="/appearance",q=[{name:"theme-color-users-0",description:""},{name:"theme-color-users-1",description:""},{name:"theme-color-users-2",description:""},{name:"theme-color-users-3",description:""},{name:"theme-color-users-4",description:""},{name:"theme-color-users-5",description:""},{name:"theme-color-users-6",description:""},{name:"theme-color-users-7",description:""}],ee=[{name:"theme-color-background-main",description:"Background"},{name:"theme-color-action",description:"Action"},{name:"theme-color-action-hover",description:"Action Hover"},{name:"theme-color-components-primary-button-border",description:"Primary Button Border"},{name:"theme-color-components-primary-button-text",description:"Primary Button Text"},{name:"theme-color-components-chat-background",description:"Chat Background"},{name:"theme-color-components-chat-text",description:"Text: Chat"},{name:"theme-color-components-text-on-dark",description:"Text: Light"},{name:"theme-color-components-text-on-light",description:"Text: Dark"},{name:"theme-color-background-header",description:"Header/Footer"},{name:"theme-color-components-content-background",description:"Page Content"},{name:"theme-color-components-video-status-bar-background",description:"Video Status Bar Background"},{name:"theme-color-components-video-status-bar-foreground",description:"Video Status Bar Foreground"}],et=[{name:"theme-rounded-corners",description:"Corner radius"}],ea=[...ee,...q,...et].reduce((e,t)=>(e[t.name]={name:t.name,description:t.description},e),{}),es=n.memo(e=>{let{value:t,name:a,description:n,onChange:l}=e;return(0,s.jsxs)(A.Z,{span:3,children:[(0,s.jsx)("input",{type:"color",id:a,name:n,title:n,value:t,className:$().colorPicker,onChange:e=>l(a,e.target.value,n)}),(0,s.jsx)("div",{style:{padding:"2px"},children:n})]},a)}),en=e=>{let{variables:t,updateColor:a}=e,n=t.map(e=>{let{name:t,description:n,value:l}=e;return(0,s.jsx)(es,{value:l,name:t,description:n,onChange:a},t)});return(0,s.jsx)(s.Fragment,{children:n})};function el(){var e,t,a,l,r,o,c,d,u;let h=(0,n.useContext)(p.a),{serverConfig:m,setFieldInConfigState:f}=h,{instanceDetails:g}=m,{appearanceVariables:j}=g,[v,k]=(0,n.useState)(),[N,w]=(0,n.useState)(),[S,C]=(0,n.useState)(null),T=()=>{C(null),clearTimeout(null)},E=()=>{let e={};[...ee,...q,...et].forEach(t=>{let a=getComputedStyle(document.documentElement).getPropertyValue("--".concat(t.name));e[t.name]={value:a.trim(),description:t.description}}),k(e)};(0,n.useEffect)(()=>{E()},[]),(0,n.useEffect)(()=>{if(0===Object.keys(j).length)return;let e={};Object.keys(j).forEach(t=>{var a;e[t]={value:j[t],description:(null===(a=ea[t])||void 0===a?void 0:a.description)||""}}),w(e)},[j]);let P=(0,n.useCallback)((e,t,a)=>{w(s=>({...s,[e]:{value:t,description:a}}))},[]),O=async()=>{await (0,x.Si)({apiPath:Q,data:{value:{}},onSuccess:()=>{C((0,b.kg)(b.zv,"Updated.")),setTimeout(T,x.sI),w({})},onError:e=>{C((0,b.kg)(b.Un,e)),setTimeout(T,x.sI)}})},Z=async()=>{let e={};Object.keys(N).forEach(t=>{e[t]=N[t].value}),await (0,x.Si)({apiPath:Q,data:{value:e},onSuccess:()=>{C((0,b.kg)(b.zv,"Updated.")),setTimeout(T,x.sI),f({fieldName:"appearanceVariables",value:e,path:"instanceDetails"})},onError:e=>{C((0,b.kg)(b.Un,e)),setTimeout(T,x.sI)}})},_=e=>{P("theme-rounded-corners","".concat(e.toString(),"px"),"")};if(!v)return(0,s.jsx)("div",{children:"Loading..."});let z=e=>e.map(e=>{let t=(null==N?void 0:N[e.name])?N:v,{name:a,description:s}=e,{value:n}=t[a];return{name:a,description:s,value:n}});return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsxs)(J.default,{direction:"vertical",children:[(0,s.jsx)(G.default,{children:"Customize Appearance"}),(0,s.jsx)(V.default,{children:"The following colors are used across the user interface."}),(0,s.jsx)("div",{children:(0,s.jsxs)(D.default,{defaultActiveKey:["1"],children:[(0,s.jsxs)(X,{header:(0,s.jsx)("strong",{children:"Section Colors"}),children:[(0,s.jsx)("p",{children:"Certain sections of the interface can be customized by selecting new colors for them."}),(0,s.jsx)(M.Z,{gutter:[16,16],children:(0,s.jsx)(en,{variables:z(ee),updateColor:P})})]},"1"),(0,s.jsx)(X,{header:(0,s.jsx)("strong",{children:"Chat User Colors"}),children:(0,s.jsx)(M.Z,{gutter:[16,16],children:(0,s.jsx)(en,{variables:z(q),updateColor:P})})},"2"),(0,s.jsxs)(X,{header:(0,s.jsx)("strong",{children:"Other Settings"}),children:["How rounded should corners be?",(0,s.jsxs)(M.Z,{gutter:[16,16],children:[(0,s.jsx)(A.Z,{span:12,children:(0,s.jsx)(R.Z,{min:0,max:20,tooltip:{formatter:null},onChange:e=>{_(e)},value:Number(null!==(d=null!==(c=null==N?void 0:null===(t=N["theme-rounded-corners"])||void 0===t?void 0:null===(e=t.value)||void 0===e?void 0:e.replace("px",""))&&void 0!==c?c:null==v?void 0:null===(l=v["theme-rounded-corners"])||void 0===l?void 0:null===(a=l.value)||void 0===a?void 0:a.replace("px",""))&&void 0!==d?d:0)})}),(0,s.jsx)(A.Z,{span:4,children:(0,s.jsx)("div",{style:{width:"100px",height:"30px",borderRadius:"".concat(null!==(u=null==N?void 0:null===(r=N["theme-rounded-corners"])||void 0===r?void 0:r.value)&&void 0!==u?u:null==v?void 0:null===(o=v["theme-rounded-corners"])||void 0===o?void 0:o.value),backgroundColor:"var(--theme-color-palette-7)"}})})]})]},"4")]})}),(0,s.jsxs)(J.default,{direction:"horizontal",children:[(0,s.jsx)(i.default,{type:"primary",onClick:Z,children:"Save Colors"}),(0,s.jsx)(i.default,{type:"ghost",onClick:O,children:"Reset to Defaults"})]}),(0,s.jsx)(y.E,{status:S})]}),(0,s.jsx)("div",{className:"form-module page-content-module",children:(0,s.jsx)(Y,{})})]})}var ei=a(34261),er=a(122);let{Title:eo}=o.default,ec=()=>{let[e,t]=(0,n.useState)("/* Enter custom Javascript here */"),[a,l]=(0,n.useState)(null),[r,o]=(0,n.useState)(!1),c=(0,n.useContext)(p.a),{serverConfig:h,setFieldInConfigState:m}=c||{},{instanceDetails:f}=h,{customJavascript:g}=f,j=null,v=()=>{l(null),o(!1),clearTimeout(j),j=null};async function k(){l((0,b.kg)(b.Jk)),await (0,x.Si)({apiPath:x.JZ,data:{value:e},onSuccess:t=>{m({fieldName:"customJavascript",value:e,path:"instanceDetails"}),l((0,b.kg)(b.zv,t))},onError:e=>{l((0,b.kg)(b.Un,e))}}),j=setTimeout(v,x.sI)}(0,n.useEffect)(()=>{t(g)},[f]);let N=n.useCallback(e=>{t(e),e===g||r?e===g&&r&&o(!1):o(!0)},[]);return(0,s.jsxs)("div",{className:"edit-custom-css",children:[(0,s.jsx)(eo,{level:3,className:"section-title",children:"Customize your page with Javascript"}),(0,s.jsxs)("p",{className:"description",children:["Insert custom Javascript into your Owncast page to add your own functionality or to add 3rd party scripts. Read more about how to use this feature in the"," ",(0,s.jsx)("a",{href:"https://owncast.online/docs/website/",rel:"noopener noreferrer",target:"_blank",children:"Web page documentation."}),"."]}),(0,s.jsx)("p",{className:"description",children:"Please use raw Javascript, no HTML or any script tags."}),(0,s.jsx)(d.ZP,{value:e,placeholder:"/* Enter custom Javascript here */",theme:u.FZ,height:"200px",extensions:[(0,er.javascript)()],onChange:N}),(0,s.jsx)("br",{}),(0,s.jsxs)("div",{className:"page-content-actions",children:[r&&(0,s.jsx)(i.default,{type:"primary",onClick:k,children:"Save"}),(0,s.jsx)(y.E,{status:a})]})]})};function ed(){return(0,s.jsx)("div",{className:"config-public-details-page",children:(0,s.jsx)(l.default,{defaultActiveKey:"1",centered:!0,items:[{label:"General",key:"1",children:(0,s.jsx)(F,{})},{label:"Appearance",key:"2",children:(0,s.jsx)(el,{})},{label:"Custom Scripting",key:"3",children:(0,s.jsx)(ec,{})}]})})}ed.getLayout=function(e){return(0,s.jsx)(ei.l,{page:e})}},80693:function(e,t,a){"use strict";a.d(t,{Z7:function(){return s},dr:function(){return n},kR:function(){return i},y3:function(){return l}});let s=2097152,n=["image/png","image/jpeg","image/gif"];function l(e,t){let a=new FileReader;a.addEventListener("load",()=>t(a.result)),a.readAsDataURL(e)}function i(e){let t=Math.floor(Math.log(e)/Math.log(1024)),a=1*Number((e/Math.pow(1024,t)).toFixed(2));return"".concat(a," ").concat(["B","KB","MB","GB","TB","PB","EB","ZB","YB"][t])}},8118:function(e){e.exports={colorPicker:"appearance_colorPicker__8KOjx"}}},function(e){e.O(0,[5762,5596,1130,4104,9403,1024,3942,971,6697,1664,1749,1700,2122,7752,5891,2891,4749,6627,8966,3068,4938,8469,5560,7423,5056,8465,4060,4261,9774,2888,179],function(){return e(e.s=39856)}),_N_E=e.O()}]);