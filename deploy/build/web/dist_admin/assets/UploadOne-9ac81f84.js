import{L as T,aq as h,M as k,O as N,Q as y,T as x,aK as S,U as V,bD as R,z as m,bH as U,C as I,k as O,a0 as D,E as P,o as g,a1 as v,a as b,w as p,n as E,c as K,u as d,d as C,bI as M,b as j}from"./index-7e9df7d9.js";import{a as F,N as H}from"./Upload-667cd4e3.js";const L=T("text",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
`,[h("strong",`
 font-weight: var(--n-font-weight-strong);
 `),h("italic",{fontStyle:"italic"}),h("underline",{textDecoration:"underline"}),h("code",`
 line-height: 1.4;
 display: inline-block;
 font-family: var(--n-font-famliy-mono);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 box-sizing: border-box;
 padding: .05em .35em 0 .35em;
 border-radius: var(--n-code-border-radius);
 font-size: .9em;
 color: var(--n-code-text-color);
 background-color: var(--n-code-color);
 border: var(--n-code-border);
 `)]),q=Object.assign(Object.assign({},y.props),{code:Boolean,type:{type:String,default:"default"},delete:Boolean,strong:Boolean,italic:Boolean,underline:Boolean,depth:[String,Number],tag:String,as:{type:String,validator:()=>!0,default:void 0}}),A=k({name:"Text",props:q,setup(e){const{mergedClsPrefixRef:l,inlineThemeDisabled:r}=N(e),o=y("Typography","-text",L,U,e,l),a=x(()=>{const{depth:s,type:c}=e,f=c==="default"?s===void 0?"textColor":`textColor${s}Depth`:S("textColor",c),{common:{fontWeightStrong:n,fontFamilyMono:u,cubicBezierEaseInOut:i},self:{codeTextColor:_,codeBorderRadius:w,codeColor:B,codeBorder:$,[f]:z}}=o.value;return{"--n-bezier":i,"--n-text-color":z,"--n-font-weight-strong":n,"--n-font-famliy-mono":u,"--n-code-border-radius":w,"--n-code-text-color":_,"--n-code-color":B,"--n-code-border":$}}),t=r?V("text",x(()=>`${e.type[0]}${e.depth||""}`),a,e):void 0;return{mergedClsPrefix:l,compitableTag:R(e,["as","tag"]),cssVars:r?void 0:a,themeClass:t==null?void 0:t.themeClass,onRender:t==null?void 0:t.onRender}},render(){var e,l,r;const{mergedClsPrefix:o}=this;(e=this.onRender)===null||e===void 0||e.call(this);const a=[`${o}-text`,this.themeClass,{[`${o}-text--code`]:this.code,[`${o}-text--delete`]:this.delete,[`${o}-text--strong`]:this.strong,[`${o}-text--italic`]:this.italic,[`${o}-text--underline`]:this.underline}],t=(r=(l=this.$slots).default)===null||r===void 0?void 0:r.call(l);return this.code?m("code",{class:a,style:this.cssVars},this.delete?m("del",null,t):t):this.delete?m("del",{class:a,style:this.cssVars},t):m(this.compitableTag||"span",{class:a,style:this.cssVars},t)}}),J=["src"],Q={class:"mb-12"},W=C("span",{class:"i-mdi:upload"},null,-1),Y={__name:"UploadOne",props:{preview:{type:String,default:""},width:{type:Number,default:120}},emits:["update:preview"],setup(e,{expose:l,emit:r}){const o=e,a=r,t=I(),s=O(o.preview);D(()=>o.preview,n=>s.value=n);function c({event:n}){const u=(n==null?void 0:n.target).response,i=JSON.parse(u);if(i.code!==0){$message==null||$message.error(i.message);return}s.value=i.data,a("update:preview",s.value)}const f=x(()=>P(s.value));return l({previewImg:s}),(n,u)=>(g(),v("div",null,[b(d(H),{action:"/api/upload",headers:{Authorization:`Bearer ${d(t)}`},"show-file-list":!1,onFinish:c},{default:p(()=>[s.value?(g(),v("img",{key:0,"border-color":"#d9d9d9",class:"cursor-pointer border-2 rounded-2rem border-dashed hover:border-color-lightblue",style:E({width:`${o.width}px`}),src:f.value,alt:"文章封面"},null,12,J)):(g(),K(d(F),{key:1},{default:p(()=>[C("div",Q,[b(d(M),{size:"58",depth:3},{default:p(()=>[W]),_:1})]),b(d(A),{class:"text-14"},{default:p(()=>[j(" 点击或者拖动文件到该区域来上传 ")]),_:1})]),_:1}))]),_:1},8,["headers"])]))}};export{Y as _};
