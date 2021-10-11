(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[64],{10610:function(e,t,n){"use strict";n.d(t,{Z:function(){return m}});var a=n(15861),r=n(87757),i=n.n(r),s=n(67294),o=n(12028),l=n(74071),c=n(91486),d=n(95828),u=n(60293),h=n(85893);function m(e){var t=(0,s.useState)(null),n=t[0],r=t[1],m=null,f=((0,s.useContext)(u.aC)||{}).setFieldInConfigState,v=e.apiPath,p=e.checked,g=e.reversed,x=void 0!==g&&g,j=e.configPath,b=void 0===j?"":j,y=e.disabled,w=void 0!==y&&y,N=e.fieldName,k=e.label,P=e.tip,C=e.useSubmit,Z=e.onChange,S=function(){r(null),clearTimeout(m),m=null},O=function(){var e=(0,a.Z)(i().mark((function e(t){var n;return i().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!C){e.next=6;break}return r((0,l.kg)(l.Jk)),n=x?!t:t,e.next=5,(0,d.Si)({apiPath:v,data:{value:n},onSuccess:function(){f({fieldName:N,value:n,path:b}),r((0,l.kg)(l.zv))},onError:function(e){r((0,l.kg)(l.Un,"There was an error: ".concat(e)))}});case 5:m=setTimeout(S,d.sI);case 6:Z&&Z(t);case 7:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}(),I=null!==n&&n.type===l.Jk;return(0,h.jsxs)("div",{className:"formfield-container toggleswitch-container",children:[k&&(0,h.jsx)("div",{className:"label-side",children:(0,h.jsx)("span",{className:"formfield-label",children:k})}),(0,h.jsxs)("div",{className:"input-side",children:[(0,h.jsxs)("div",{className:"input-group",children:[(0,h.jsx)(o.Z,{className:"switch field-".concat(N),loading:I,onChange:O,defaultChecked:p,checked:p,checkedChildren:"ON",unCheckedChildren:"OFF",disabled:w}),(0,h.jsx)(c.Z,{status:n})]}),(0,h.jsx)("p",{className:"field-tip",children:P})]})]})}m.defaultProps={apiPath:"",checked:!1,reversed:!1,configPath:"",disabled:!1,label:"",tip:"",useSubmit:!1,onChange:null}},5352:function(e,t,n){"use strict";n.r(t),n.d(t,{default:function(){return q}});var a=n(27279),r=n(4525),i=n(36725),s=n(1635),o=n(67294),l=n(15861),c=n(29439),d=n(87757),u=n.n(d),h=n(7031),m=n(75443),f=n(45197),v=n(95828),p=n(74071),g=n(60293),x=n(91486),j=n(85893);function b(){var e=(0,o.useContext)(g.aC),t=e||{},n=t.serverConfig,a=t.setFieldInConfigState,i=n||{},s=i.videoCodec,d=i.supportedCodecs,b=r.Z.Title,y=h.Z.Option,w=(0,o.useState)(null),N=w[0],k=w[1],P=(0,o.useContext)(f.k).setMessage,C=(0,o.useState)(s),Z=C[0],S=C[1],O=(0,o.useState)(s),I=O[0],T=O[1],V=o.useState(!1),_=(0,c.Z)(V,2),U=_[0],L=_[1],E=null;(0,o.useEffect)((function(){S(s)}),[s]);var D=function(){k(null),E=null,clearTimeout(E)};function F(){return(F=(0,l.Z)(u().mark((function t(){return u().wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return S(I),T(""),L(!1),t.next=5,(0,v.Si)({apiPath:v.CQ,data:{value:I},onSuccess:function(){a({fieldName:"videoCodec",value:I,path:"videoSettings"}),k((0,p.kg)(p.zv,"Video codec updated.")),E=setTimeout(D,v.sI),e.online&&P("Your latency buffer setting will take effect the next time you begin a live stream.")},onError:function(e){k((0,p.kg)(p.Un,e)),E=setTimeout(D,v.sI)}});case 5:case"end":return t.stop()}}),t)})))).apply(this,arguments)}var A=d.map((function(e){var t=e;return"libx264"===t?t="Default (libx264)":"h264_nvenc"===t?t="NVIDIA GPU acceleration":"h264_vaapi"===t?t="VA-API hardware encoding":"h264_qsv"===t?t="Intel QuickSync":"h264_v4l2m2m"===t&&(t="Video4Linux hardware encoding"),(0,j.jsx)(y,{value:e,children:t},e)})),B="";return"libx264"===Z?B="libx264 is the default codec and generally the only working choice for shared VPS enviornments. This is likely what you should be using unless you know you have set up other options.":"h264_nvenc"===Z?B="You can use your NVIDIA GPU for encoding if you have a modern NVIDIA card with encoding cores.":"h264_vaapi"===Z?B="VA-API may be supported by your NVIDIA proprietary drivers, Mesa open-source drivers for AMD or Intel graphics.":"h264_qsv"===Z?B="Quick Sync Video is Intel's brand for its dedicated video encoding and decoding hardware. It may be an option if you have a modern Intel CPU with integrated graphics.":"h264_v4l2m2m"===Z&&(B="Video4Linux is an interface to multiple different hardware encoding platforms such as Intel and AMD."),(0,j.jsxs)(j.Fragment,{children:[(0,j.jsx)(b,{level:3,className:"section-title",children:"Video Codec"}),(0,j.jsxs)("div",{className:"description",children:["If you have access to specific hardware with the drivers and software installed for them, you may be able to improve your video encoding performance.",(0,j.jsx)("p",{children:(0,j.jsx)("a",{href:"https://owncast.online/docs/codecs?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Read the documentation about this setting before changing it or you may make your stream unplayable."})})]}),(0,j.jsxs)("div",{className:"segment-slider-container",children:[(0,j.jsx)(m.Z,{title:"Are you sure you want to change your video codec to ".concat(I," and understand what this means?"),visible:U,placement:"leftBottom",onConfirm:function(){return F.apply(this,arguments)},okText:"Yes",cancelText:"No",children:(0,j.jsx)(h.Z,{defaultValue:Z,value:Z,style:{width:"100%"},onChange:function(e){T(e),L(!0)},children:A})}),(0,j.jsx)(x.Z,{status:N}),(0,j.jsx)("p",{id:"selected-codec-note",className:"selected-value-note",children:B})]})]})}var y=n(64913),w=r.Z.Title,N={0:"Lowest",1:"",2:"",3:"",4:"Highest"},k={0:"Lowest latency, lowest error tolerance (Not recommended, may not work for all content/configurations.)",1:"Low latency, low error tolerance",2:"Medium latency, medium error tolerance (Default)",3:"High latency, high error tolerance",4:"Highest latency, highest error tolerance"};function P(){var e=(0,o.useState)(null),t=e[0],n=e[1],a=(0,o.useState)(null),r=a[0],i=a[1],s=(0,o.useContext)(g.aC),c=(0,o.useContext)(f.k).setMessage,d=s||{},h=d.serverConfig,m=d.setFieldInConfigState,b=(h||{}).videoSettings,P=null;if(!b)return null;(0,o.useEffect)((function(){i(b.latencyLevel)}),[b]);var C=function(){n(null),P=null,clearTimeout(P)},Z=function(){var e=(0,l.Z)(u().mark((function e(t){return u().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return n((0,p.kg)(p.Jk)),e.next=3,(0,v.Si)({apiPath:v.sv,data:{value:t},onSuccess:function(){m({fieldName:"latencyLevel",value:t,path:"videoSettings"}),n((0,p.kg)(p.zv,"Latency buffer level updated.")),P=setTimeout(C,v.sI),s.online&&c("Your latency buffer setting will take effect the next time you begin a live stream.")},onError:function(e){n((0,p.kg)(p.Un,e)),P=setTimeout(C,v.sI)}});case 3:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}();return(0,j.jsxs)("div",{className:"config-video-latency-container",children:[(0,j.jsx)(w,{level:3,className:"section-title",children:"Latency Buffer"}),(0,j.jsx)("p",{className:"description",children:"While it's natural to want to keep your latency as low as possible, you may experience reduced error tolerance and stability the lower you go. The lowest setting is not recommended."}),(0,j.jsxs)("p",{className:"description",children:["For interactive live streams you may want to experiment with a lower latency, for non-interactive broadcasts you may want to increase it."," ",(0,j.jsx)("a",{href:"https://owncast.online/docs/encoding#latency-buffer?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Read to learn more."})]}),(0,j.jsxs)("div",{className:"segment-slider-container",children:[(0,j.jsx)(y.Z,{tipFormatter:function(e){return k[e]},onChange:function(e){Z(e)},min:0,max:4,marks:N,defaultValue:r,value:r}),(0,j.jsx)("p",{className:"selected-value-note",children:k[r]}),(0,j.jsx)(x.Z,{status:t})]})]})}var C=n(4942),Z=n(93433),S=n(71577),O=n(99645),I=n(37614),T=n(73171),V=n(68855),_=n(94184),U=n.n(_),L=n(60764),E=n(10610);function D(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function F(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?D(Object(n),!0).forEach((function(t){(0,C.Z)(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):D(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}var A=a.Z.Panel;function B(e){var t=e.dataState,n=void 0===t?v.gX:t,o=e.onUpdateField,l=n.videoPassthrough,c=U()({"config-variant-form":!0,"video-passthrough-enabled":l});return(0,j.jsxs)("div",{className:c,children:[(0,j.jsxs)("p",{className:"description",children:[(0,j.jsx)("a",{href:"https://owncast.online/docs/video?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Learn more"})," ","about how each of these settings can impact the performance of your server."]}),l&&(0,j.jsxs)("p",{className:"passthrough-warning",children:["NOTE: Video Passthrough for this output stream variant is ",(0,j.jsx)("em",{children:"enabled"}),", disabling the below video encoding settings."]}),(0,j.jsxs)(i.Z,{gutter:16,children:[(0,j.jsx)(L.ZP,F(F({maxLength:"10"},v.SS),{},{value:n.name,onChange:function(e){o({fieldName:"name",value:e.value})}})),(0,j.jsx)(s.Z,{sm:24,md:12,children:(0,j.jsxs)("div",{className:"form-module cpu-usage-container",children:[(0,j.jsx)(r.Z.Title,{level:3,children:"CPU or GPU Utilization"}),(0,j.jsx)("p",{className:"description",children:"Reduce to improve server performance, or increase it to improve video quality."}),(0,j.jsxs)("div",{className:"segment-slider-container",children:[(0,j.jsx)(y.Z,{tipFormatter:function(e){return v.I$[e]},onChange:function(e){o({fieldName:"cpuUsageLevel",value:e})},min:1,max:Object.keys(v.t$).length,marks:v.t$,defaultValue:n.cpuUsageLevel,value:n.cpuUsageLevel,disabled:n.videoPassthrough}),(0,j.jsx)("p",{className:"selected-value-note",children:l?"CPU usage selection is disabled when Video Passthrough is enabled.":v.I$[n.cpuUsageLevel]||""})]}),(0,j.jsxs)("p",{className:"read-more-subtext",children:["This could mean GPU or CPU usage depending on your server environment."," ",(0,j.jsx)("a",{href:"https://owncast.online/docs/video/?source=admin#cpu-usage",target:"_blank",rel:"noopener noreferrer",children:"Read more about hardware performance."})]})]})}),(0,j.jsx)(s.Z,{sm:24,md:12,children:(0,j.jsxs)("div",{className:"form-module bitrate-container ".concat(n.videoPassthrough?"disabled":""),children:[(0,j.jsx)(r.Z.Title,{level:3,children:"Video Bitrate"}),(0,j.jsx)("p",{className:"description",children:v.yC.tip}),(0,j.jsxs)("div",{className:"segment-slider-container",children:[(0,j.jsx)(y.Z,{tipFormatter:function(e){return"".concat(e," ").concat(v.yC.unit)},disabled:n.videoPassthrough,defaultValue:n.videoBitrate,value:n.videoBitrate,onChange:function(e){o({fieldName:"videoBitrate",value:e})},step:v.yC.incrementBy,min:v.yC.min,max:v.yC.max,marks:v.HM}),(0,j.jsx)("p",{className:"selected-value-note",children:function(){if(l)return"Bitrate selection is disabled when Video Passthrough is enabled.";var e="".concat(n.videoBitrate).concat(v.yC.unit);return e=n.videoBitrate<2e3?"".concat(e," - Good for low bandwidth environments."):n.videoBitrate<3500?"".concat(e," - Good for most bandwidth environments."):"".concat(e," - Good for high bandwidth environments.")}()})]}),(0,j.jsx)("p",{className:"read-more-subtext",children:(0,j.jsx)("a",{href:"https://owncast.online/docs/video/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Read more about bitrates."})})]})})]}),(0,j.jsx)(a.Z,{className:"advanced-settings",children:(0,j.jsxs)(A,{header:"Advanced Settings",children:[(0,j.jsxs)(i.Z,{gutter:16,children:[(0,j.jsx)(s.Z,{sm:24,md:12,children:(0,j.jsxs)("div",{className:"form-module resolution-module",children:[(0,j.jsx)(r.Z.Title,{level:3,children:"Resolution"}),(0,j.jsxs)("p",{className:"description",children:["Resizing your content will take additional resources on your server. If you wish to optionally resize your content for this stream output then you should either set the width ",(0,j.jsx)("strong",{children:"or"})," the height to keep your aspect ratio."," ",(0,j.jsx)("a",{href:"https://owncast.online/docs/video/?source=admin",target:"_blank",rel:"noopener noreferrer",children:"Read more about resolutions."})]}),(0,j.jsx)("br",{}),(0,j.jsx)(L.ZP,F(F({type:"number"},v.dL.scaledWidth),{},{value:n.scaledWidth,onChange:function(e){var t=Number(e.value);isNaN(t)||o({fieldName:"scaledWidth",value:t||0})},disabled:n.videoPassthrough})),(0,j.jsx)(L.ZP,F(F({type:"number"},v.dL.scaledHeight),{},{value:n.scaledHeight,onChange:function(e){var t=Number(e.value);isNaN(t)||o({fieldName:"scaledHeight",value:t||0})},disabled:n.videoPassthrough}))]})}),(0,j.jsx)(s.Z,{sm:24,md:12,children:(0,j.jsxs)("div",{className:"form-module video-passthrough-module",children:[(0,j.jsx)(r.Z.Title,{level:3,children:"Video Passthrough"}),(0,j.jsxs)("div",{className:"description",children:[(0,j.jsxs)("p",{children:["Enabling video passthrough may allow for less hardware utilization, but may also make your stream ",(0,j.jsx)("strong",{children:"unplayable"}),"."]}),(0,j.jsx)("p",{children:"All other settings for this stream output will be disabled if passthrough is used."}),(0,j.jsx)("p",{children:(0,j.jsx)("a",{href:"https://owncast.online/docs/video/?source=admin#video-passthrough",target:"_blank",rel:"noopener noreferrer",children:"Read the documentation before enabling, as it impacts your stream."})})]}),(0,j.jsx)(m.Z,{disabled:!0===n.videoPassthrough,title:"Did you read the documentation about video passthrough and understand the risks involved with enabling it?",icon:(0,j.jsx)(V.Z,{}),onConfirm:function(){o({fieldName:"videoPassthrough",value:!0})},okText:"Yes",cancelText:"No",children:(0,j.jsx)("a",{href:"#",children:(0,j.jsx)(E.Z,{label:"Use Video Passthrough?",fieldName:"video-passthrough",tip:v.dL.videoPassthrough.tip,checked:n.videoPassthrough,onChange:function(e){l&&o({fieldName:"videoPassthrough",value:e})}})})})]})})]}),(0,j.jsxs)("div",{className:"form-module frame-rate-module",children:[(0,j.jsx)(r.Z.Title,{level:3,children:"Frame rate"}),(0,j.jsx)("p",{className:"description",children:v.nm.tip}),(0,j.jsxs)("div",{className:"segment-slider-container",children:[(0,j.jsx)(y.Z,{tipFormatter:function(e){return"".concat(e," ").concat(v.nm.unit)},defaultValue:n.framerate,value:n.framerate,onChange:function(e){o({fieldName:"framerate",value:e})},step:v.nm.incrementBy,min:v.nm.min,max:v.nm.max,marks:v.Xq,disabled:n.videoPassthrough}),(0,j.jsx)("p",{className:"selected-value-note",children:l?"Framerate selection is disabled when Video Passthrough is enabled.":v.x8[n.framerate]||""})]}),(0,j.jsx)("p",{className:"read-more-subtext",children:(0,j.jsx)("a",{href:"https://owncast.online/docs/video/?source=admin#framerate",target:"_blank",rel:"noopener noreferrer",children:"Read more about framerates."})})]})]},"1")})]})}function z(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function R(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?z(Object(n),!0).forEach((function(t){(0,C.Z)(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):z(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}var M=r.Z.Title;function G(){var e=(0,o.useState)(!1),t=e[0],n=e[1],a=(0,o.useState)(!1),r=a[0],i=a[1],s=(0,o.useState)(0),c=s[0],d=s[1],h=(0,o.useContext)(f.k).setMessage,m=(0,o.useState)(v.gX),b=m[0],y=m[1],w=(0,o.useState)(null),N=w[0],k=w[1],P=(0,o.useContext)(g.aC),V=P||{},_=V.serverConfig,U=V.setFieldInConfigState,L=(_||{}).videoSettings,E=(L||{}).videoQualityVariants,D=null;if(!L)return null;var F=function(){k(null),D=null,clearTimeout(D)},A=function(){n(!1),d(-1),y(v.gX)},z=function(){var e=(0,l.Z)(u().mark((function e(t){return u().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return k((0,p.kg)(p.Jk)),e.next=3,(0,v.Si)({apiPath:v.vv,data:{value:t},onSuccess:function(){U({fieldName:"videoQualityVariants",value:t,path:"videoSettings"}),i(!1),A(),k((0,p.kg)(p.zv,"Variants updated")),D=setTimeout(F,v.sI),P.online&&h("Updating your video configuration will take effect the next time you begin a new stream.")},onError:function(e){k((0,p.kg)(p.Un,e)),i(!1),D=setTimeout(F,v.sI)}});case 3:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}(),G=[{title:"Name",dataIndex:"name",render:function(e){return e||"No name"}},{title:"Video bitrate",dataIndex:"videoBitrate",key:"videoBitrate",render:function(e,t){return!e||t.videoPassthrough?"Same as source":"".concat(e," kbps")}},{title:"CPU Usage",dataIndex:"cpuUsageLevel",key:"cpuUsageLevel",render:function(e,t){return!e||t.videoPassthrough?"n/a":v.I$[e].split(" ")[0]}},{title:"",dataIndex:"",key:"edit",render:function(e){var t=e.key-1;return(0,j.jsxs)("span",{className:"actions",children:[(0,j.jsx)(S.Z,{size:"small",onClick:function(){d(t),y(E[t]),n(!0)},children:"Edit"}),(0,j.jsx)(S.Z,{className:"delete-button",icon:(0,j.jsx)(T.Z,{}),size:"small",disabled:1===E.length,onClick:function(){!function(e){var t=(0,Z.Z)(E);t.splice(e,1),z(t)}(t)}})]})}}],H=E.map((function(e,t){return R({key:t+1},e)}));return(0,j.jsxs)(j.Fragment,{children:[(0,j.jsx)(M,{level:3,className:"section-title",children:"Stream output"}),(0,j.jsx)(x.Z,{status:N}),(0,j.jsx)(O.Z,{className:"variants-table",pagination:!1,size:"small",columns:G,dataSource:H}),(0,j.jsxs)(I.Z,{title:"Edit Video Variant Details",visible:t,onOk:function(){i(!0);var e=(0,Z.Z)(E);-1===c?e.push(b):e.splice(c,1,b),z(e)},onCancel:A,confirmLoading:r,width:900,children:[(0,j.jsx)(B,{dataState:R({},b),onUpdateField:function(e){var t=e.fieldName,n=e.value;y(R(R({},b),{},(0,C.Z)({},t,n)))}}),(0,j.jsx)(x.Z,{status:N})]}),(0,j.jsx)("br",{}),(0,j.jsx)(S.Z,{type:"primary",onClick:function(){d(-1),y(v.gX),n(!0)},children:"Add a new variant"})]})}var H=a.Z.Panel,X=r.Z.Title;function q(){return(0,j.jsxs)("div",{className:"config-video-variants",children:[(0,j.jsx)(X,{children:"Video configuration"}),(0,j.jsxs)("p",{className:"description",children:["Before changing your video configuration"," ",(0,j.jsx)("a",{href:"https://owncast.online/docs/video?source=admin",target:"_blank",rel:"noopener noreferrer",children:"visit the video documentation"})," ","to learn how it impacts your stream performance. The general rule is to start conservatively by having one middle quality stream output variant and experiment with adding more of varied qualities."]}),(0,j.jsxs)(i.Z,{gutter:[16,16],children:[(0,j.jsx)(s.Z,{md:24,lg:12,children:(0,j.jsx)("div",{className:"form-module variants-table-module",children:(0,j.jsx)(G,{})})}),(0,j.jsxs)(s.Z,{md:24,lg:12,children:[(0,j.jsx)("div",{className:"form-module latency-module",children:(0,j.jsx)(P,{})}),(0,j.jsx)(a.Z,{className:"advanced-settings codec-module",children:(0,j.jsx)(H,{header:"Advanced Settings",children:(0,j.jsx)("div",{className:"form-module variants-table-module",children:(0,j.jsx)(b,{})})},"1")})]})]})]})}},79893:function(e,t,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/config-video",function(){return n(5352)}])}},function(e){e.O(0,[645,614,701,364,774,888,179],(function(){return t=79893,e(e.s=t);var t}));var t=e.O();_N_E=t}]);