"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[3482],{33482:function(O,T,E){function N(O){for(var T={},E=O.split(" "),N=0;N<E.length;++N)T[E[N]]=!0;return T}E.r(T),E.d(T,{pig:function(){return C}});var I="ABS ACOS ARITY ASIN ATAN AVG BAGSIZE BINSTORAGE BLOOM BUILDBLOOM CBRT CEIL CONCAT COR COS COSH COUNT COUNT_STAR COV CONSTANTSIZE CUBEDIMENSIONS DIFF DISTINCT DOUBLEABS DOUBLEAVG DOUBLEBASE DOUBLEMAX DOUBLEMIN DOUBLEROUND DOUBLESUM EXP FLOOR FLOATABS FLOATAVG FLOATMAX FLOATMIN FLOATROUND FLOATSUM GENERICINVOKER INDEXOF INTABS INTAVG INTMAX INTMIN INTSUM INVOKEFORDOUBLE INVOKEFORFLOAT INVOKEFORINT INVOKEFORLONG INVOKEFORSTRING INVOKER ISEMPTY JSONLOADER JSONMETADATA JSONSTORAGE LAST_INDEX_OF LCFIRST LOG LOG10 LOWER LONGABS LONGAVG LONGMAX LONGMIN LONGSUM MAX MIN MAPSIZE MONITOREDUDF NONDETERMINISTIC OUTPUTSCHEMA  PIGSTORAGE PIGSTREAMING RANDOM REGEX_EXTRACT REGEX_EXTRACT_ALL REPLACE ROUND SIN SINH SIZE SQRT STRSPLIT SUBSTRING SUM STRINGCONCAT STRINGMAX STRINGMIN STRINGSIZE TAN TANH TOBAG TOKENIZE TOMAP TOP TOTUPLE TRIM TEXTLOADER TUPLESIZE UCFIRST UPPER UTF8STORAGECONVERTER ",e="VOID IMPORT RETURNS DEFINE LOAD FILTER FOREACH ORDER CUBE DISTINCT COGROUP JOIN CROSS UNION SPLIT INTO IF OTHERWISE ALL AS BY USING INNER OUTER ONSCHEMA PARALLEL PARTITION GROUP AND OR NOT GENERATE FLATTEN ASC DESC IS STREAM THROUGH STORE MAPREDUCE SHIP CACHE INPUT OUTPUT STDERROR STDIN STDOUT LIMIT SAMPLE LEFT RIGHT FULL EQ GT LT GTE LTE NEQ MATCHES TRUE FALSE DUMP",A="BOOLEAN INT LONG FLOAT DOUBLE CHARARRAY BYTEARRAY BAG TUPLE MAP ",R=N(I),S=N(e),t=N(A),L=/[*+\-%<>=&?:\/!|]/;function r(O,T,E){return T.tokenize=E,E(O,T)}function n(O,T){for(var E,N=!1;E=O.next();){if("/"==E&&N){T.tokenize=U;break}N="*"==E}return"comment"}function U(O,T){var E=O.next();return'"'==E||"'"==E?r(O,T,function(O,T){for(var N,I=!1,e=!1;null!=(N=O.next());){if(N==E&&!I){e=!0;break}I=!I&&"\\"==N}return(e||!I)&&(T.tokenize=U),"error"}):/[\[\]{}\(\),;\.]/.test(E)?null:/\d/.test(E)?(O.eatWhile(/[\w\.]/),"number"):"/"==E?O.eat("*")?r(O,T,n):(O.eatWhile(L),"operator"):"-"==E?O.eat("-")?(O.skipToEnd(),"comment"):(O.eatWhile(L),"operator"):L.test(E)?(O.eatWhile(L),"operator"):(O.eatWhile(/[\w\$_]/),S&&S.propertyIsEnumerable(O.current().toUpperCase())&&!O.eat(")")&&!O.eat("."))?"keyword":R&&R.propertyIsEnumerable(O.current().toUpperCase())?"builtin":t&&t.propertyIsEnumerable(O.current().toUpperCase())?"type":"variable"}let C={name:"pig",startState:function(){return{tokenize:U,startOfLine:!0}},token:function(O,T){return O.eatSpace()?null:T.tokenize(O,T)},languageData:{autocomplete:(I+A+e).split(" ")}}}}]);