import{ah as G,b0 as Q,T as C,H as X,a3 as x,L as Y,aq as B,as as Z,M as ee,O as te,k as T,x as H,a0 as z,R as P,Q as I,U as le,z as b,o as ne,c as oe,w as i,m as E,a as s,b as y,u as l,l as re,t as w,y as ae,d as se,F,b1 as ie,N as k,A as ue,B as R}from"./index-7e9df7d9.js";import{_ as de}from"./CommonPage-7491aea4.js";import{_ as ce}from"./QueryItem-0f282244.js";import{_ as he}from"./CrudModal-d60ad9bd.js";import{_ as me}from"./CrudTable-121684b7.js";import{u as fe}from"./useCRUD-f64dec52.js";import{N as pe}from"./Input-5db271d7.js";import{N as ge,a as g}from"./FormItem-e99f0659.js";import{N as be}from"./Popconfirm-6b45faf2.js";import"./AppPage-848b5b84.js";import"./Space-f54cd109.js";import"./RadioGroup-757891b8.js";import"./get-slot-1efb97e5.js";import"./Checkbox-ac1591ea.js";import"./Forward-8d8448e7.js";function _e(n,e){const a=G(Q,null);return C(()=>n.hljs||(a==null?void 0:a.mergedHljsRef.value))}const ye=n=>{const{textColor2:e,fontSize:a,fontWeightStrong:h,textColor3:m}=n;return{textColor:e,fontSize:a,fontWeightStrong:h,"mono-3":"#a0a1a7","hue-1":"#0184bb","hue-2":"#4078f2","hue-3":"#a626a4","hue-4":"#50a14f","hue-5":"#e45649","hue-5-2":"#c91243","hue-6":"#986801","hue-6-2":"#c18401",lineNumberTextColor:m}},$e={name:"Code",common:X,self:ye},ve=$e,je=x([Y("code",`
 font-size: var(--n-font-size);
 font-family: var(--n-font-family);
 `,[B("show-line-numbers",`
 display: flex;
 `),Z("line-numbers",`
 user-select: none;
 padding-right: 12px;
 text-align: right;
 transition: color .3s var(--n-bezier);
 color: var(--n-line-number-text-color);
 `),B("word-wrap",[x("pre",`
 white-space: pre-wrap;
 word-break: break-all;
 `)]),x("pre",`
 margin: 0;
 line-height: inherit;
 font-size: inherit;
 font-family: inherit;
 `),x("[class^=hljs]",`
 color: var(--n-text-color);
 transition: 
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `)]),({props:n})=>{const e=`${n.bPrefix}code`;return[`${e} .hljs-comment,
 ${e} .hljs-quote {
 color: var(--n-mono-3);
 font-style: italic;
 }`,`${e} .hljs-doctag,
 ${e} .hljs-keyword,
 ${e} .hljs-formula {
 color: var(--n-hue-3);
 }`,`${e} .hljs-section,
 ${e} .hljs-name,
 ${e} .hljs-selector-tag,
 ${e} .hljs-deletion,
 ${e} .hljs-subst {
 color: var(--n-hue-5);
 }`,`${e} .hljs-literal {
 color: var(--n-hue-1);
 }`,`${e} .hljs-string,
 ${e} .hljs-regexp,
 ${e} .hljs-addition,
 ${e} .hljs-attribute,
 ${e} .hljs-meta-string {
 color: var(--n-hue-4);
 }`,`${e} .hljs-built_in,
 ${e} .hljs-class .hljs-title {
 color: var(--n-hue-6-2);
 }`,`${e} .hljs-attr,
 ${e} .hljs-variable,
 ${e} .hljs-template-variable,
 ${e} .hljs-type,
 ${e} .hljs-selector-class,
 ${e} .hljs-selector-attr,
 ${e} .hljs-selector-pseudo,
 ${e} .hljs-number {
 color: var(--n-hue-6);
 }`,`${e} .hljs-symbol,
 ${e} .hljs-bullet,
 ${e} .hljs-link,
 ${e} .hljs-meta,
 ${e} .hljs-selector-id,
 ${e} .hljs-title {
 color: var(--n-hue-2);
 }`,`${e} .hljs-emphasis {
 font-style: italic;
 }`,`${e} .hljs-strong {
 font-weight: var(--n-font-weight-strong);
 }`,`${e} .hljs-link {
 text-decoration: underline;
 }`]}]),we=Object.assign(Object.assign({},I.props),{language:String,code:{type:String,default:""},trim:{type:Boolean,default:!0},hljs:Object,uri:Boolean,inline:Boolean,wordWrap:Boolean,showLineNumbers:Boolean,internalFontSize:Number,internalNoHighlight:Boolean}),L=ee({name:"Code",props:we,setup(n,{slots:e}){const{internalNoHighlight:a}=n,{mergedClsPrefixRef:h,inlineThemeDisabled:m}=te(),p=T(null),$=a?{value:void 0}:_e(n),d=(t,r,u)=>{const{value:c}=$;return!c||!(t&&c.getLanguage(t))?null:c.highlight(u?r.trim():r,{language:t}).value},S=C(()=>n.inline||n.wordWrap?!1:n.showLineNumbers),_=()=>{if(e.default)return;const{value:t}=p;if(!t)return;const{language:r}=n,u=n.uri?window.decodeURIComponent(n.code):n.code;if(r){const f=d(r,u,n.trim);if(f!==null){if(n.inline)t.innerHTML=f;else{const N=t.querySelector(".__code__");N&&t.removeChild(N);const j=document.createElement("pre");j.className="__code__",j.innerHTML=f,t.appendChild(j)}return}}if(n.inline){t.textContent=u;return}const c=t.querySelector(".__code__");if(c)c.textContent=u;else{const f=document.createElement("pre");f.className="__code__",f.textContent=u,t.innerHTML="",t.appendChild(f)}};H(_),z(P(n,"language"),_),z(P(n,"code"),_),a||z($,_);const q=I("Code","-code",je,ve,n,h),v=C(()=>{const{common:{cubicBezierEaseInOut:t,fontFamilyMono:r},self:{textColor:u,fontSize:c,fontWeightStrong:f,lineNumberTextColor:N,"mono-3":j,"hue-1":V,"hue-2":D,"hue-3":M,"hue-4":U,"hue-5":J,"hue-5-2":W,"hue-6":K,"hue-6-2":A}}=q.value,{internalFontSize:O}=n;return{"--n-font-size":O?`${O}px`:c,"--n-font-family":r,"--n-font-weight-strong":f,"--n-bezier":t,"--n-text-color":u,"--n-mono-3":j,"--n-hue-1":V,"--n-hue-2":D,"--n-hue-3":M,"--n-hue-4":U,"--n-hue-5":J,"--n-hue-5-2":W,"--n-hue-6":K,"--n-hue-6-2":A,"--n-line-number-text-color":N}}),o=m?le("code",C(()=>`${n.internalFontSize||"a"}`),v,n):void 0;return{mergedClsPrefix:h,codeRef:p,mergedShowLineNumbers:S,lineNumbers:C(()=>{let t=1;const r=[];let u=!1;for(const c of n.code)c===`
`?(u=!0,r.push(t++)):u=!1;return u||r.push(t++),r.join(`
`)}),cssVars:m?void 0:v,themeClass:o==null?void 0:o.themeClass,onRender:o==null?void 0:o.onRender}},render(){var n,e;const{mergedClsPrefix:a,wordWrap:h,mergedShowLineNumbers:m,onRender:p}=this;return p==null||p(),b("code",{class:[`${a}-code`,this.themeClass,h&&`${a}-code--word-wrap`,m&&`${a}-code--show-line-numbers`],style:this.cssVars,ref:"codeRef"},m?b("pre",{class:`${a}-code__line-numbers`},this.lineNumbers):null,(e=(n=this.$slots).default)===null||e===void 0?void 0:e.call(n))}}),Ce=se("span",{class:"i-material-symbols:playlist-remove"},null,-1),Ie=Object.assign({name:"操作日志"},{__name:"index",setup(n){function e(o){switch(o){case"GET":return"info";case"POST":return"success";case"PUT":return"warning";case"DELETE":return"error";default:return"info"}}const a=T(null),h=T({keyword:""}),{modalVisible:m,modalLoading:p,handleDelete:$,modalForm:d,modalFormRef:S,handleView:_}=fe({name:"日志",doDelete:E.deleteOperationLogs,refresh:()=>{var o;return(o=a.value)==null?void 0:o.handleSearch()}});H(()=>{var o;(o=a.value)==null||o.handleSearch()});const q=[{type:"selection",width:20,fixed:"left"},{title:"系统模块",key:"opt_module",width:70,align:"center",ellipsis:{tooltip:!0}},{title:"操作类型",key:"opt_type",width:70,align:"center",ellipsis:{tooltip:!0}},{title:"请求方法",key:"request_method",width:80,align:"center",ellipsis:{tooltip:!0},render(o){return b(F,{type:e(o.request_method)},{default:()=>o.request_method})}},{title:"操作人员",key:"nickname",width:80,align:"center",ellipsis:{tooltip:!0}},{title:"登录IP",key:"ip_address",width:80,align:"center",ellipsis:{tooltip:!0}},{title:"登录地址",key:"ip_source",width:80,align:"center",ellipsis:{tooltip:!0}},{title:"发布时间",key:"created_at",align:"center",width:80,render(o){return b(k,{size:"small",type:"text",ghost:!0},{default:()=>ue(o.created_at),icon:R("mdi:update",{size:18})})}},{title:"操作",key:"actions",width:120,align:"center",fixed:"right",render(o){return[b(k,{size:"small",quaternary:!0,type:"info",onClick:()=>_(o)},{default:()=>"查看",icon:R("ic:outline-remove-red-eye",{})}),b(be,{onPositiveClick:()=>$([o.id],!1)},{trigger:()=>b(k,{size:"small",quaternary:!0,type:"error",style:"margin-left: 15px;"},{default:()=>"删除",icon:R("material-symbols:delete-outline",{})}),default:()=>b("div",{},"确定删除该日志吗?")})]}}];function v(o){const{copy:t}=ie();t(JSON.stringify(JSON.parse(o),null,2)),window.$message.success("内容已复制到剪切板!")}return(o,t)=>(ne(),oe(de,{title:"操作日志"},{action:i(()=>{var r;return[s(l(k),{type:"error",disabled:!((r=a.value)!=null&&r.selections.length),onClick:t[0]||(t[0]=u=>{var c;return l($)((c=a.value)==null?void 0:c.selections)})},{icon:i(()=>[Ce]),default:i(()=>[y(" 批量删除 ")]),_:1},8,["disabled"])]}),default:i(()=>[s(me,{ref_key:"$table",ref:a,"query-items":h.value,"onUpdate:queryItems":t[3]||(t[3]=r=>h.value=r),columns:q,"get-data":l(E).getOperationLogs},{queryBar:i(()=>[s(ce,{label:"模块名","label-width":50},{default:i(()=>[s(l(pe),{value:h.value.keyword,"onUpdate:value":t[1]||(t[1]=r=>h.value.keyword=r),clearable:"",type:"text",placeholder:"请输入模块名或描述",onKeydown:t[2]||(t[2]=re(r=>{var u;return(u=a.value)==null?void 0:u.handleSearch()},["enter"]))},null,8,["value"])]),_:1})]),_:1},8,["query-items","get-data"]),s(he,{visible:l(m),"onUpdate:visible":t[6]||(t[6]=r=>ae(m)?m.value=r:null),title:"日志详情","show-footer":!1,loading:l(p),width:"full"},{default:i(()=>[s(l(ge),{ref_key:"modalFormRef",ref:S,"label-placement":"left","label-align":"left","label-width":90,model:l(d)},{default:i(()=>[s(l(g),{label:"操作模块: ",path:"opt_module"},{default:i(()=>[y(w(l(d).opt_module),1)]),_:1}),s(l(g),{label:"请求地址: ",path:"opt_url"},{default:i(()=>[y(w(l(d).opt_url),1)]),_:1}),s(l(g),{label:"请求方法: ",path:"request_method"},{default:i(()=>[s(l(F),{type:e(l(d).request_method)},{default:i(()=>[y(w(l(d).request_method),1)]),_:1},8,["type"])]),_:1}),s(l(g),{label:"操作类型: ",path:"opt_type"},{default:i(()=>[y(w(l(d).opt_type),1)]),_:1}),s(l(g),{label:"操作方法: ",path:"opt_method"},{default:i(()=>[s(l(L),{code:l(d).opt_method,"code-wrap":"",language:"json"},null,8,["code"])]),_:1}),s(l(g),{label:"操作人员: ",path:"nickname"},{default:i(()=>[y(w(l(d).nickname),1)]),_:1}),s(l(g),{label:"请求参数: ",path:"request_param"},{default:i(()=>[s(l(L),{class:"word-wrap cursor-pointer p-7",code:JSON.stringify(JSON.parse(l(d).request_param),null,2),language:"json",onClick:t[4]||(t[4]=r=>v(l(d).request_param))},null,8,["code"])]),_:1}),s(l(g),{label:"返回数据: ",path:"response_data"},{default:i(()=>[s(l(L),{class:"cursor-pointer p-7",code:JSON.stringify(JSON.parse(l(d).response_data),null,2),language:"json",onClick:t[5]||(t[5]=r=>v(l(d).response_data))},null,8,["code"])]),_:1})]),_:1},8,["model"])]),_:1},8,["visible","loading"])]),_:1}))}});export{Ie as default};
