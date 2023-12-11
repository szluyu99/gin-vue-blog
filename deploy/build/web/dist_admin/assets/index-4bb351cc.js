import{ab as G,aW as X,M as C,B as Y,Y as x,F as A,ak as T,am as Q,G as Z,H as ee,e as L,p as F,V as z,K as B,J as H,O as te,s as f,o as le,f as ne,w as i,l as P,j as a,h as y,i as l,k as oe,t as w,q as re,g as se,z as E,aX as ae,N as k,v as ie}from"./index-cfa27ed4.js";import{_ as ue}from"./CommonPage-4f39ed42.js";import{_ as de}from"./QueryItem-d17f875d.js";import{_ as ce}from"./CrudModal-388e823f.js";import{_ as he}from"./CrudTable-f7b664e6.js";import{u as me}from"./useCRUD-d1880a98.js";import{N as fe}from"./Input-9d39db42.js";import{N as pe,a as _}from"./FormItem-0c6ab39e.js";import{N as ge}from"./Popconfirm-c434d7b3.js";import"./AppPage-a5c29498.js";import"./Space-c006c775.js";import"./RadioGroup-3711a790.js";import"./get-slot-1efb97e5.js";import"./Checkbox-13b79107.js";import"./Forward-d629ea4a.js";function _e(n,e){const s=G(X,null);return C(()=>n.hljs||(s==null?void 0:s.mergedHljsRef.value))}const be=n=>{const{textColor2:e,fontSize:s,fontWeightStrong:h,textColor3:m}=n;return{textColor:e,fontSize:s,fontWeightStrong:h,"mono-3":"#a0a1a7","hue-1":"#0184bb","hue-2":"#4078f2","hue-3":"#a626a4","hue-4":"#50a14f","hue-5":"#e45649","hue-5-2":"#c91243","hue-6":"#986801","hue-6-2":"#c18401",lineNumberTextColor:m}},ye={name:"Code",common:Y,self:be},ve=ye,$e=x([A("code",`
 font-size: var(--n-font-size);
 font-family: var(--n-font-family);
 `,[T("show-line-numbers",`
 display: flex;
 `),Q("line-numbers",`
 user-select: none;
 padding-right: 12px;
 text-align: right;
 transition: color .3s var(--n-bezier);
 color: var(--n-line-number-text-color);
 `),T("word-wrap",[x("pre",`
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
 }`]}]),je=Object.assign(Object.assign({},H.props),{language:String,code:{type:String,default:""},trim:{type:Boolean,default:!0},hljs:Object,uri:Boolean,inline:Boolean,wordWrap:Boolean,showLineNumbers:Boolean,internalFontSize:Number,internalNoHighlight:Boolean}),R=Z({name:"Code",props:je,setup(n,{slots:e}){const{internalNoHighlight:s}=n,{mergedClsPrefixRef:h,inlineThemeDisabled:m}=ee(),g=L(null),v=s?{value:void 0}:_e(n),d=(t,r,u)=>{const{value:c}=v;return!c||!(t&&c.getLanguage(t))?null:c.highlight(u?r.trim():r,{language:t}).value},S=C(()=>n.inline||n.wordWrap?!1:n.showLineNumbers),b=()=>{if(e.default)return;const{value:t}=g;if(!t)return;const{language:r}=n,u=n.uri?window.decodeURIComponent(n.code):n.code;if(r){const p=d(r,u,n.trim);if(p!==null){if(n.inline)t.innerHTML=p;else{const N=t.querySelector(".__code__");N&&t.removeChild(N);const j=document.createElement("pre");j.className="__code__",j.innerHTML=p,t.appendChild(j)}return}}if(n.inline){t.textContent=u;return}const c=t.querySelector(".__code__");if(c)c.textContent=u;else{const p=document.createElement("pre");p.className="__code__",p.textContent=u,t.innerHTML="",t.appendChild(p)}};F(b),z(B(n,"language"),b),z(B(n,"code"),b),s||z(v,b);const q=H("Code","-code",$e,ve,n,h),$=C(()=>{const{common:{cubicBezierEaseInOut:t,fontFamilyMono:r},self:{textColor:u,fontSize:c,fontWeightStrong:p,lineNumberTextColor:N,"mono-3":j,"hue-1":I,"hue-2":V,"hue-3":D,"hue-4":J,"hue-5":M,"hue-5-2":W,"hue-6":U,"hue-6-2":K}}=q.value,{internalFontSize:O}=n;return{"--n-font-size":O?`${O}px`:c,"--n-font-family":r,"--n-font-weight-strong":p,"--n-bezier":t,"--n-text-color":u,"--n-mono-3":j,"--n-hue-1":I,"--n-hue-2":V,"--n-hue-3":D,"--n-hue-4":J,"--n-hue-5":M,"--n-hue-5-2":W,"--n-hue-6":U,"--n-hue-6-2":K,"--n-line-number-text-color":N}}),o=m?te("code",C(()=>`${n.internalFontSize||"a"}`),$,n):void 0;return{mergedClsPrefix:h,codeRef:g,mergedShowLineNumbers:S,lineNumbers:C(()=>{let t=1;const r=[];let u=!1;for(const c of n.code)c===`
`?(u=!0,r.push(t++)):u=!1;return u||r.push(t++),r.join(`
`)}),cssVars:m?void 0:$,themeClass:o==null?void 0:o.themeClass,onRender:o==null?void 0:o.onRender}},render(){var n,e;const{mergedClsPrefix:s,wordWrap:h,mergedShowLineNumbers:m,onRender:g}=this;return g==null||g(),f("code",{class:[`${s}-code`,this.themeClass,h&&`${s}-code--word-wrap`,m&&`${s}-code--show-line-numbers`],style:this.cssVars,ref:"codeRef"},m?f("pre",{class:`${s}-code__line-numbers`},this.lineNumbers):null,(e=(n=this.$slots).default)===null||e===void 0?void 0:e.call(n))}}),we=se("span",{class:"i-material-symbols:playlist-remove"},null,-1),He=Object.assign({name:"操作日志"},{__name:"index",setup(n){function e(o){switch(o){case"GET":return"info";case"POST":return"success";case"PUT":return"warning";case"DELETE":return"error";default:return"info"}}const s=L(null),h=L({keyword:""}),{modalVisible:m,modalLoading:g,handleDelete:v,modalForm:d,modalFormRef:S,handleView:b}=me({name:"日志",doDelete:P.deleteOperationLogs,refresh:()=>{var o;return(o=s.value)==null?void 0:o.handleSearch()}});F(()=>{var o;(o=s.value)==null||o.handleSearch()});const q=[{type:"selection",width:20,fixed:"left"},{title:"系统模块",key:"opt_module",width:70,align:"center",ellipsis:{tooltip:!0}},{title:"操作类型",key:"opt_type",width:70,align:"center",ellipsis:{tooltip:!0}},{title:"请求方法",key:"request_method",width:80,align:"center",ellipsis:{tooltip:!0},render(o){return f(E,{type:e(o.request_method)},{default:()=>o.request_method})}},{title:"操作人员",key:"nickname",width:80,align:"center",ellipsis:{tooltip:!0}},{title:"登录IP",key:"ip_address",width:80,align:"center",ellipsis:{tooltip:!0}},{title:"登录地址",key:"ip_source",width:80,align:"center",ellipsis:{tooltip:!0}},{title:"发布时间",key:"created_at",align:"center",width:80,render(o){return f(k,{size:"small",type:"text",ghost:!0},{default:()=>ie(o.created_at),icon:()=>f("i",{class:"i-mdi:update"})})}},{title:"操作",key:"actions",width:120,align:"center",fixed:"right",render(o){return[f(k,{size:"small",quaternary:!0,type:"info",onClick:()=>b(o)},{default:()=>"查看",icon:()=>f("i",{class:"i-ic:outline-remove-red-eye"})}),f(ge,{onPositiveClick:()=>v([o.id],!1)},{trigger:()=>f(k,{size:"small",quaternary:!0,type:"error",style:"margin-left: 15px;"},{default:()=>"删除",icon:()=>f("i",{class:"i-material-symbols:delete-outline"})}),default:()=>f("div",{},"确定删除该日志吗?")})]}}];function $(o){const{copy:t}=ae();t(JSON.stringify(JSON.parse(o),null,2)),window.$message.success("内容已复制到剪切板!")}return(o,t)=>(le(),ne(ue,{title:"操作日志"},{action:i(()=>{var r;return[a(l(k),{type:"error",disabled:!((r=s.value)!=null&&r.selections.length),onClick:t[0]||(t[0]=u=>{var c;return l(v)((c=s.value)==null?void 0:c.selections)})},{icon:i(()=>[we]),default:i(()=>[y(" 批量删除 ")]),_:1},8,["disabled"])]}),default:i(()=>[a(he,{ref_key:"$table",ref:s,"query-items":h.value,"onUpdate:queryItems":t[3]||(t[3]=r=>h.value=r),columns:q,"get-data":l(P).getOperationLogs},{queryBar:i(()=>[a(de,{label:"模块名","label-width":50},{default:i(()=>[a(l(fe),{value:h.value.keyword,"onUpdate:value":t[1]||(t[1]=r=>h.value.keyword=r),clearable:"",type:"text",placeholder:"请输入模块名或描述",onKeydown:t[2]||(t[2]=oe(r=>{var u;return(u=s.value)==null?void 0:u.handleSearch()},["enter"]))},null,8,["value"])]),_:1})]),_:1},8,["query-items","get-data"]),a(ce,{visible:l(m),"onUpdate:visible":t[6]||(t[6]=r=>re(m)?m.value=r:null),title:"日志详情","show-footer":!1,loading:l(g),width:"full"},{default:i(()=>[a(l(pe),{ref_key:"modalFormRef",ref:S,"label-placement":"left","label-align":"left","label-width":90,model:l(d)},{default:i(()=>[a(l(_),{label:"操作模块: ",path:"opt_module"},{default:i(()=>[y(w(l(d).opt_module),1)]),_:1}),a(l(_),{label:"请求地址: ",path:"opt_url"},{default:i(()=>[y(w(l(d).opt_url),1)]),_:1}),a(l(_),{label:"请求方法: ",path:"request_method"},{default:i(()=>[a(l(E),{type:e(l(d).request_method)},{default:i(()=>[y(w(l(d).request_method),1)]),_:1},8,["type"])]),_:1}),a(l(_),{label:"操作类型: ",path:"opt_type"},{default:i(()=>[y(w(l(d).opt_type),1)]),_:1}),a(l(_),{label:"操作方法: ",path:"opt_method"},{default:i(()=>[a(l(R),{code:l(d).opt_method,"code-wrap":"",language:"json"},null,8,["code"])]),_:1}),a(l(_),{label:"操作人员: ",path:"nickname"},{default:i(()=>[y(w(l(d).nickname),1)]),_:1}),a(l(_),{label:"请求参数: ",path:"request_param"},{default:i(()=>[a(l(R),{class:"word-wrap cursor-pointer p-7",code:JSON.stringify(JSON.parse(l(d).request_param),null,2),language:"json",onClick:t[4]||(t[4]=r=>$(l(d).request_param))},null,8,["code"])]),_:1}),a(l(_),{label:"返回数据: ",path:"response_data"},{default:i(()=>[a(l(R),{class:"cursor-pointer p-7",code:JSON.stringify(JSON.parse(l(d).response_data),null,2),language:"json",onClick:t[5]||(t[5]=r=>$(l(d).response_data))},null,8,["code"])]),_:1})]),_:1},8,["model"])]),_:1},8,["visible","loading"])]),_:1}))}});export{He as default};
