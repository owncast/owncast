"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[7762],{17762:function(e,t,n){n.r(t),n.d(t,{solr:function(){return u}});var o=/[^\s\|\!\+\-\*\?\~\^\&\:\(\)\[\]\{\}\"\\]/,r=/[\|\!\+\-\*\?\~\^\&]/,a=/^(OR|AND|NOT|TO)$/i;function tokenBase(e,t){var n,u=e.next();return'"'==u?t.tokenize=function(e,t){for(var n,o=!1;null!=(n=e.next())&&(n!=u||o);)o=!o&&"\\"==n;return o||(t.tokenize=tokenBase),"string"}:r.test(u)?t.tokenize=function(e,t){return"|"==u?e.eat(/\|/):"&"==u&&e.eat(/\&/),t.tokenize=tokenBase,"operator"}:o.test(u)&&(t.tokenize=(n=u,function(e,t){for(var r,u=n;(n=e.peek())&&null!=n.match(o);)u+=e.next();return(t.tokenize=tokenBase,a.test(u))?"operator":parseFloat(r=u).toString()===r?"number":":"==e.peek()?"propertyName":"string"})),t.tokenize!=tokenBase?t.tokenize(e,t):null}let u={name:"solr",startState:function(){return{tokenize:tokenBase}},token:function(e,t){return e.eatSpace()?null:t.tokenize(e,t)}}}}]);