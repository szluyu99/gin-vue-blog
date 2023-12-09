import{M as oe,k,bv as jn,bw as Wn,x as Ne,W as tn,z as r,bs as Hn,bx as Kn,cb as an,aG as vn,L as T,as as _,a3 as J,O as on,Q as te,cc as Un,ah as ln,b0 as Gn,T as B,aK as ne,U as Le,X as pn,cd as rn,a4 as Ge,ak as Ce,bn as gn,aq as G,at as Xe,bf as bn,R as H,ce as qn,az as Zn,a0 as Se,ax as en,aB as qe,aA as sn,a7 as Yn,aj as Jn,be as Qn,a9 as Xn,an as Ae,cf as et,cg as nt,aw as tt,F as Ze,af as ot,ad as it,ch as lt,S as dn,av as rt,bD as at,P as st,bi as dt,bj as nn,bk as ct,bl as ut,bm as ft,bo as ht,aE as vt,bp as cn,ci as pt,br as gt,bq as bt,Y as ee,cj as mt,al as wt,a5 as yt,ck as xt,aR as Ct}from"./index-7e9df7d9.js";import{u as mn,a as St}from"./Input-5db271d7.js";import{V as Ot,F as Rt}from"./RadioGroup-757891b8.js";import{g as Ft}from"./get-slot-1efb97e5.js";function Xt(e){switch(e){case"tiny":return"mini";case"small":return"tiny";case"medium":return"small";case"large":return"medium";case"huge":return"large"}throw Error(`${e} has no smaller size.`)}function zt(e){switch(typeof e){case"string":return e||void 0;case"number":return String(e);default:return}}function Ye(e){const t=e.filter(o=>o!==void 0);if(t.length!==0)return t.length===1?t[0]:o=>{e.forEach(s=>{s&&s(o)})}}const pe="v-hidden",Mt=Kn("[v-hidden]",{display:"none!important"}),un=oe({name:"Overflow",props:{getCounter:Function,getTail:Function,updateCounter:Function,onUpdateOverflow:Function},setup(e,{slots:t}){const o=k(null),s=k(null);function f(){const{value:v}=o,{getCounter:a,getTail:R}=e;let b;if(a!==void 0?b=a():b=s.value,!v||!b)return;b.hasAttribute(pe)&&b.removeAttribute(pe);const{children:g}=v,x=v.offsetWidth,C=[],P=t.tail?R==null?void 0:R():null;let p=P?P.offsetWidth:0,S=!1;const I=v.children.length-(t.tail?1:0);for(let w=0;w<I-1;++w){if(w<0)continue;const N=g[w];if(S){N.hasAttribute(pe)||N.setAttribute(pe,"");continue}else N.hasAttribute(pe)&&N.removeAttribute(pe);const $=N.offsetWidth;if(p+=$,C[w]=$,p>x){const{updateCounter:W}=e;for(let A=w;A>=0;--A){const j=I-1-A;W!==void 0?W(j):b.textContent=`${j}`;const K=b.offsetWidth;if(p-=C[A],p+K<=x||A===0){S=!0,w=A-1,P&&(w===-1?(P.style.maxWidth=`${x-K}px`,P.style.boxSizing="border-box"):P.style.maxWidth="");break}}}}const{onUpdateOverflow:m}=e;S?m!==void 0&&m(!0):(m!==void 0&&m(!1),b.setAttribute(pe,""))}const u=jn();return Mt.mount({id:"vueuc/overflow",head:!0,anchorMetaName:Wn,ssr:u}),Ne(f),{selfRef:o,counterRef:s,sync:f}},render(){const{$slots:e}=this;return tn(this.sync),r("div",{class:"v-overflow",ref:"selfRef"},[Hn(e,"default"),e.counter?e.counter():r("span",{style:{display:"inline-block"},ref:"counterRef"}),e.tail?e.tail():null])}});function wn(e,t){t&&(Ne(()=>{const{value:o}=e;o&&an.registerHandler(o,t)}),vn(()=>{const{value:o}=e;o&&an.unregisterHandler(o)}))}const Tt=oe({name:"Checkmark",render(){return r("svg",{xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 16 16"},r("g",{fill:"none"},r("path",{d:"M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z",fill:"currentColor"})))}}),Pt=oe({name:"Empty",render(){return r("svg",{viewBox:"0 0 28 28",fill:"none",xmlns:"http://www.w3.org/2000/svg"},r("path",{d:"M26 7.5C26 11.0899 23.0899 14 19.5 14C15.9101 14 13 11.0899 13 7.5C13 3.91015 15.9101 1 19.5 1C23.0899 1 26 3.91015 26 7.5ZM16.8536 4.14645C16.6583 3.95118 16.3417 3.95118 16.1464 4.14645C15.9512 4.34171 15.9512 4.65829 16.1464 4.85355L18.7929 7.5L16.1464 10.1464C15.9512 10.3417 15.9512 10.6583 16.1464 10.8536C16.3417 11.0488 16.6583 11.0488 16.8536 10.8536L19.5 8.20711L22.1464 10.8536C22.3417 11.0488 22.6583 11.0488 22.8536 10.8536C23.0488 10.6583 23.0488 10.3417 22.8536 10.1464L20.2071 7.5L22.8536 4.85355C23.0488 4.65829 23.0488 4.34171 22.8536 4.14645C22.6583 3.95118 22.3417 3.95118 22.1464 4.14645L19.5 6.79289L16.8536 4.14645Z",fill:"currentColor"}),r("path",{d:"M25 22.75V12.5991C24.5572 13.0765 24.053 13.4961 23.5 13.8454V16H17.5L17.3982 16.0068C17.0322 16.0565 16.75 16.3703 16.75 16.75C16.75 18.2688 15.5188 19.5 14 19.5C12.4812 19.5 11.25 18.2688 11.25 16.75L11.2432 16.6482C11.1935 16.2822 10.8797 16 10.5 16H4.5V7.25C4.5 6.2835 5.2835 5.5 6.25 5.5H12.2696C12.4146 4.97463 12.6153 4.47237 12.865 4H6.25C4.45507 4 3 5.45507 3 7.25V22.75C3 24.5449 4.45507 26 6.25 26H21.75C23.5449 26 25 24.5449 25 22.75ZM4.5 22.75V17.5H9.81597L9.85751 17.7041C10.2905 19.5919 11.9808 21 14 21L14.215 20.9947C16.2095 20.8953 17.842 19.4209 18.184 17.5H23.5V22.75C23.5 23.7165 22.7165 24.5 21.75 24.5H6.25C5.2835 24.5 4.5 23.7165 4.5 22.75Z",fill:"currentColor"}))}}),kt=T("empty",`
 display: flex;
 flex-direction: column;
 align-items: center;
 font-size: var(--n-font-size);
`,[_("icon",`
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 font-size: var(--n-icon-size);
 line-height: var(--n-icon-size);
 color: var(--n-icon-color);
 transition:
 color .3s var(--n-bezier);
 `,[J("+",[_("description",`
 margin-top: 8px;
 `)])]),_("description",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 `),_("extra",`
 text-align: center;
 transition: color .3s var(--n-bezier);
 margin-top: 12px;
 color: var(--n-extra-text-color);
 `)]),Bt=Object.assign(Object.assign({},te.props),{description:String,showDescription:{type:Boolean,default:!0},showIcon:{type:Boolean,default:!0},size:{type:String,default:"medium"},renderIcon:Function}),It=oe({name:"Empty",props:Bt,setup(e){const{mergedClsPrefixRef:t,inlineThemeDisabled:o}=on(e),s=te("Empty","-empty",kt,Un,e,t),{localeRef:f}=mn("Empty"),u=ln(Gn,null),v=B(()=>{var g,x,C;return(g=e.description)!==null&&g!==void 0?g:(C=(x=u==null?void 0:u.mergedComponentPropsRef.value)===null||x===void 0?void 0:x.Empty)===null||C===void 0?void 0:C.description}),a=B(()=>{var g,x;return((x=(g=u==null?void 0:u.mergedComponentPropsRef.value)===null||g===void 0?void 0:g.Empty)===null||x===void 0?void 0:x.renderIcon)||(()=>r(Pt,null))}),R=B(()=>{const{size:g}=e,{common:{cubicBezierEaseInOut:x},self:{[ne("iconSize",g)]:C,[ne("fontSize",g)]:P,textColor:p,iconColor:S,extraTextColor:I}}=s.value;return{"--n-icon-size":C,"--n-font-size":P,"--n-bezier":x,"--n-text-color":p,"--n-icon-color":S,"--n-extra-text-color":I}}),b=o?Le("empty",B(()=>{let g="";const{size:x}=e;return g+=x[0],g}),R,e):void 0;return{mergedClsPrefix:t,mergedRenderIcon:a,localizedDescription:B(()=>v.value||f.value.description),cssVars:o?void 0:R,themeClass:b==null?void 0:b.themeClass,onRender:b==null?void 0:b.onRender}},render(){const{$slots:e,mergedClsPrefix:t,onRender:o}=this;return o==null||o(),r("div",{class:[`${t}-empty`,this.themeClass],style:this.cssVars},this.showIcon?r("div",{class:`${t}-empty__icon`},e.icon?e.icon():r(pn,{clsPrefix:t},{default:this.mergedRenderIcon})):null,this.showDescription?r("div",{class:`${t}-empty__description`},e.default?e.default():this.localizedDescription):null,e.extra?r("div",{class:`${t}-empty__extra`},e.extra()):null)}});function _t(e,t){return r(gn,{name:"fade-in-scale-up-transition"},{default:()=>e?r(pn,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>r(Tt)}):null})}const fn=oe({name:"NBaseSelectOption",props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){const{valueRef:t,pendingTmNodeRef:o,multipleRef:s,valueSetRef:f,renderLabelRef:u,renderOptionRef:v,labelFieldRef:a,valueFieldRef:R,showCheckmarkRef:b,nodePropsRef:g,handleOptionClick:x,handleOptionMouseEnter:C}=ln(rn),P=Ge(()=>{const{value:m}=o;return m?e.tmNode.key===m.key:!1});function p(m){const{tmNode:w}=e;w.disabled||x(m,w)}function S(m){const{tmNode:w}=e;w.disabled||C(m,w)}function I(m){const{tmNode:w}=e,{value:N}=P;w.disabled||N||C(m,w)}return{multiple:s,isGrouped:Ge(()=>{const{tmNode:m}=e,{parent:w}=m;return w&&w.rawNode.type==="group"}),showCheckmark:b,nodeProps:g,isPending:P,isSelected:Ge(()=>{const{value:m}=t,{value:w}=s;if(m===null)return!1;const N=e.tmNode.rawNode[R.value];if(w){const{value:$}=f;return $.has(N)}else return m===N}),labelField:a,renderLabel:u,renderOption:v,handleMouseMove:I,handleMouseEnter:S,handleClick:p}},render(){const{clsPrefix:e,tmNode:{rawNode:t},isSelected:o,isPending:s,isGrouped:f,showCheckmark:u,nodeProps:v,renderOption:a,renderLabel:R,handleClick:b,handleMouseEnter:g,handleMouseMove:x}=this,C=_t(o,e),P=R?[R(t,o),u&&C]:[Ce(t[this.labelField],t,o),u&&C],p=v==null?void 0:v(t),S=r("div",Object.assign({},p,{class:[`${e}-base-select-option`,t.class,p==null?void 0:p.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:o,[`${e}-base-select-option--grouped`]:f,[`${e}-base-select-option--pending`]:s,[`${e}-base-select-option--show-checkmark`]:u}],style:[(p==null?void 0:p.style)||"",t.style||""],onClick:Ye([b,p==null?void 0:p.onClick]),onMouseenter:Ye([g,p==null?void 0:p.onMouseenter]),onMousemove:Ye([x,p==null?void 0:p.onMousemove])}),r("div",{class:`${e}-base-select-option__content`},P));return t.render?t.render({node:S,option:t,selected:o}):a?a({node:S,option:t,selected:o}):S}}),hn=oe({name:"NBaseSelectGroupHeader",props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){const{renderLabelRef:e,renderOptionRef:t,labelFieldRef:o,nodePropsRef:s}=ln(rn);return{labelField:o,nodeProps:s,renderLabel:e,renderOption:t}},render(){const{clsPrefix:e,renderLabel:t,renderOption:o,nodeProps:s,tmNode:{rawNode:f}}=this,u=s==null?void 0:s(f),v=t?t(f,!1):Ce(f[this.labelField],f,!1),a=r("div",Object.assign({},u,{class:[`${e}-base-select-group-header`,u==null?void 0:u.class]}),v);return f.render?f.render({node:a,option:f}):o?o({node:a,option:f,selected:!1}):a}}),$t=T("base-select-menu",`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[T("scrollbar",`
 max-height: var(--n-height);
 `),T("virtual-list",`
 max-height: var(--n-height);
 `),T("base-select-option",`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[_("content",`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),T("base-select-group-header",`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),T("base-select-menu-option-wrapper",`
 position: relative;
 width: 100%;
 `),_("loading, empty",`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),_("loading",`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),_("action",`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),T("base-select-group-header",`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),T("base-select-option",`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[G("show-checkmark",`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),J("&::before",`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),J("&:active",`
 color: var(--n-option-text-color-pressed);
 `),G("grouped",`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),G("pending",[J("&::before",`
 background-color: var(--n-option-color-pending);
 `)]),G("selected",`
 color: var(--n-option-text-color-active);
 `,[J("&::before",`
 background-color: var(--n-option-color-active);
 `),G("pending",[J("&::before",`
 background-color: var(--n-option-color-active-pending);
 `)])]),G("disabled",`
 cursor: not-allowed;
 `,[Xe("selected",`
 color: var(--n-option-text-color-disabled);
 `),G("selected",`
 opacity: var(--n-option-opacity-disabled);
 `)]),_("check",`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[bn({enterScale:"0.5"})])])]),At=oe({name:"InternalSelectMenu",props:Object.assign(Object.assign({},te.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:"medium"},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:"label"},valueField:{type:String,default:"value"},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,onToggle:Function}),setup(e){const t=te("InternalSelectMenu","-internal-select-menu",$t,qn,e,H(e,"clsPrefix")),o=k(null),s=k(null),f=k(null),u=B(()=>e.treeMate.getFlattenedNodes()),v=B(()=>Zn(u.value)),a=k(null);function R(){const{treeMate:l}=e;let c=null;const{value:L}=e;L===null?c=l.getFirstAvailableNode():(e.multiple?c=l.getNode((L||[])[(L||[]).length-1]):c=l.getNode(L),(!c||c.disabled)&&(c=l.getFirstAvailableNode())),q(c||null)}function b(){const{value:l}=a;l&&!e.treeMate.getNode(l.key)&&(a.value=null)}let g;Se(()=>e.show,l=>{l?g=Se(()=>e.treeMate,()=>{e.resetMenuOnOptionsChange?(e.autoPending?R():b(),tn(V)):b()},{immediate:!0}):g==null||g()},{immediate:!0}),vn(()=>{g==null||g()});const x=B(()=>en(t.value.self[ne("optionHeight",e.size)])),C=B(()=>qe(t.value.self[ne("padding",e.size)])),P=B(()=>e.multiple&&Array.isArray(e.value)?new Set(e.value):new Set),p=B(()=>{const l=u.value;return l&&l.length===0});function S(l){const{onToggle:c}=e;c&&c(l)}function I(l){const{onScroll:c}=e;c&&c(l)}function m(l){var c;(c=f.value)===null||c===void 0||c.sync(),I(l)}function w(){var l;(l=f.value)===null||l===void 0||l.sync()}function N(){const{value:l}=a;return l||null}function $(l,c){c.disabled||q(c,!1)}function W(l,c){c.disabled||S(c)}function A(l){var c;Ae(l,"action")||(c=e.onKeyup)===null||c===void 0||c.call(e,l)}function j(l){var c;Ae(l,"action")||(c=e.onKeydown)===null||c===void 0||c.call(e,l)}function K(l){var c;(c=e.onMousedown)===null||c===void 0||c.call(e,l),!e.focusable&&l.preventDefault()}function ce(){const{value:l}=a;l&&q(l.getNext({loop:!0}),!0)}function Q(){const{value:l}=a;l&&q(l.getPrev({loop:!0}),!0)}function q(l,c=!1){a.value=l,c&&V()}function V(){var l,c;const L=a.value;if(!L)return;const le=v.value(L.key);le!==null&&(e.virtualScroll?(l=s.value)===null||l===void 0||l.scrollTo({index:le}):(c=f.value)===null||c===void 0||c.scrollTo({index:le,elSize:x.value}))}function ue(l){var c,L;!((c=o.value)===null||c===void 0)&&c.contains(l.target)&&((L=e.onFocus)===null||L===void 0||L.call(e,l))}function ge(l){var c,L;!((c=o.value)===null||c===void 0)&&c.contains(l.relatedTarget)||(L=e.onBlur)===null||L===void 0||L.call(e,l)}sn(rn,{handleOptionMouseEnter:$,handleOptionClick:W,valueSetRef:P,pendingTmNodeRef:a,nodePropsRef:H(e,"nodeProps"),showCheckmarkRef:H(e,"showCheckmark"),multipleRef:H(e,"multiple"),valueRef:H(e,"value"),renderLabelRef:H(e,"renderLabel"),renderOptionRef:H(e,"renderOption"),labelFieldRef:H(e,"labelField"),valueFieldRef:H(e,"valueField")}),sn(et,o),Ne(()=>{const{value:l}=f;l&&l.sync()});const fe=B(()=>{const{size:l}=e,{common:{cubicBezierEaseInOut:c},self:{height:L,borderRadius:le,color:Oe,groupHeaderTextColor:Re,actionDividerColor:Fe,optionTextColorPressed:be,optionTextColor:me,optionTextColorDisabled:re,optionTextColorActive:U,optionOpacityDisabled:we,optionCheckColor:se,actionTextColor:ze,optionColorPending:he,optionColorActive:ve,loadingColor:Me,loadingSize:Te,optionColorActivePending:Pe,[ne("optionFontSize",l)]:ye,[ne("optionHeight",l)]:xe,[ne("optionPadding",l)]:Z}}=t.value;return{"--n-height":L,"--n-action-divider-color":Fe,"--n-action-text-color":ze,"--n-bezier":c,"--n-border-radius":le,"--n-color":Oe,"--n-option-font-size":ye,"--n-group-header-text-color":Re,"--n-option-check-color":se,"--n-option-color-pending":he,"--n-option-color-active":ve,"--n-option-color-active-pending":Pe,"--n-option-height":xe,"--n-option-opacity-disabled":we,"--n-option-text-color":me,"--n-option-text-color-active":U,"--n-option-text-color-disabled":re,"--n-option-text-color-pressed":be,"--n-option-padding":Z,"--n-option-padding-left":qe(Z,"left"),"--n-option-padding-right":qe(Z,"right"),"--n-loading-color":Me,"--n-loading-size":Te}}),{inlineThemeDisabled:ie}=e,Y=ie?Le("internal-select-menu",B(()=>e.size[0]),fe,e):void 0,ae={selfRef:o,next:ce,prev:Q,getPendingTmNode:N};return wn(o,e.onResize),Object.assign({mergedTheme:t,virtualListRef:s,scrollbarRef:f,itemSize:x,padding:C,flattenedNodes:u,empty:p,virtualListContainer(){const{value:l}=s;return l==null?void 0:l.listElRef},virtualListContent(){const{value:l}=s;return l==null?void 0:l.itemsElRef},doScroll:I,handleFocusin:ue,handleFocusout:ge,handleKeyUp:A,handleKeyDown:j,handleMouseDown:K,handleVirtualListResize:w,handleVirtualListScroll:m,cssVars:ie?void 0:fe,themeClass:Y==null?void 0:Y.themeClass,onRender:Y==null?void 0:Y.onRender},ae)},render(){const{$slots:e,virtualScroll:t,clsPrefix:o,mergedTheme:s,themeClass:f,onRender:u}=this;return u==null||u(),r("div",{ref:"selfRef",tabindex:this.focusable?0:-1,class:[`${o}-base-select-menu`,f,this.multiple&&`${o}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},this.loading?r("div",{class:`${o}-base-select-menu__loading`},r(Jn,{clsPrefix:o,strokeWidth:20})):this.empty?r("div",{class:`${o}-base-select-menu__empty`,"data-empty":!0,"data-action":!0},Xn(e.empty,()=>[r(It,{theme:s.peers.Empty,themeOverrides:s.peerOverrides.Empty})])):r(Qn,{ref:"scrollbarRef",theme:s.peers.Scrollbar,themeOverrides:s.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},{default:()=>t?r(Ot,{ref:"virtualListRef",class:`${o}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:v})=>v.isGroup?r(hn,{key:v.key,clsPrefix:o,tmNode:v}):v.ignored?null:r(fn,{clsPrefix:o,key:v.key,tmNode:v})}):r("div",{class:`${o}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(v=>v.isGroup?r(hn,{key:v.key,clsPrefix:o,tmNode:v}):r(fn,{clsPrefix:o,key:v.key,tmNode:v})))}),Yn(e.action,v=>v&&[r("div",{class:`${o}-base-select-menu__action`,"data-action":!0,key:"action"},v),r(Rt,{onFocus:this.onTabOut,key:"focus-detector"})]))}}),Et=J([T("base-selection",`
 position: relative;
 z-index: auto;
 box-shadow: none;
 width: 100%;
 max-width: 100%;
 display: inline-block;
 vertical-align: bottom;
 border-radius: var(--n-border-radius);
 min-height: var(--n-height);
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[T("base-loading",`
 color: var(--n-loading-color);
 `),T("base-selection-tags","min-height: var(--n-height);"),_("border, state-border",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border: var(--n-border);
 border-radius: inherit;
 transition:
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),_("state-border",`
 z-index: 1;
 border-color: #0000;
 `),T("base-suffix",`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[_("arrow",`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),T("base-selection-overlay",`
 display: flex;
 align-items: center;
 white-space: nowrap;
 pointer-events: none;
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 padding: var(--n-padding-single);
 transition: color .3s var(--n-bezier);
 `,[_("wrapper",`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),T("base-selection-placeholder",`
 color: var(--n-placeholder-color);
 `,[_("inner",`
 max-width: 100%;
 overflow: hidden;
 `)]),T("base-selection-tags",`
 cursor: pointer;
 outline: none;
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 display: flex;
 padding: var(--n-padding-multiple);
 flex-wrap: wrap;
 align-items: center;
 width: 100%;
 vertical-align: bottom;
 background-color: var(--n-color);
 border-radius: inherit;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),T("base-selection-label",`
 height: var(--n-height);
 display: inline-flex;
 width: 100%;
 vertical-align: bottom;
 cursor: pointer;
 outline: none;
 z-index: auto;
 box-sizing: border-box;
 position: relative;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 border-radius: inherit;
 background-color: var(--n-color);
 align-items: center;
 `,[T("base-selection-input",`
 font-size: inherit;
 line-height: inherit;
 outline: none;
 cursor: pointer;
 box-sizing: border-box;
 border:none;
 width: 100%;
 padding: var(--n-padding-single);
 background-color: #0000;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 caret-color: var(--n-caret-color);
 `,[_("content",`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),_("render-label",`
 color: var(--n-text-color);
 `)]),Xe("disabled",[J("&:hover",[_("state-border",`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),G("focus",[_("state-border",`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),G("active",[_("state-border",`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),T("base-selection-label","background-color: var(--n-color-active);"),T("base-selection-tags","background-color: var(--n-color-active);")])]),G("disabled","cursor: not-allowed;",[_("arrow",`
 color: var(--n-arrow-color-disabled);
 `),T("base-selection-label",`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[T("base-selection-input",`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),_("render-label",`
 color: var(--n-text-color-disabled);
 `)]),T("base-selection-tags",`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),T("base-selection-placeholder",`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),T("base-selection-input-tag",`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[_("input",`
 font-size: inherit;
 font-family: inherit;
 min-width: 1px;
 padding: 0;
 background-color: #0000;
 outline: none;
 border: none;
 max-width: 100%;
 overflow: hidden;
 width: 1em;
 line-height: inherit;
 cursor: pointer;
 color: var(--n-text-color);
 caret-color: var(--n-caret-color);
 `),_("mirror",`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),["warning","error"].map(e=>G(`${e}-status`,[_("state-border",`border: var(--n-border-${e});`),Xe("disabled",[J("&:hover",[_("state-border",`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),G("active",[_("state-border",`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),T("base-selection-label",`background-color: var(--n-color-active-${e});`),T("base-selection-tags",`background-color: var(--n-color-active-${e});`)]),G("focus",[_("state-border",`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),T("base-selection-popover",`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),T("base-selection-tag-wrapper",`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[J("&:last-child","padding-right: 0;"),T("tag",`
 font-size: 14px;
 max-width: 100%;
 `,[_("content",`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),Nt=oe({name:"InternalSelection",props:Object.assign(Object.assign({},te.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:""},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:"label"},valueField:{type:String,default:"value"},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:"medium"},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(e){const t=k(null),o=k(null),s=k(null),f=k(null),u=k(null),v=k(null),a=k(null),R=k(null),b=k(null),g=k(null),x=k(!1),C=k(!1),P=k(!1),p=te("InternalSelection","-internal-selection",Et,nt,e,H(e,"clsPrefix")),S=B(()=>e.clearable&&!e.disabled&&(P.value||e.active)),I=B(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):Ce(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),m=B(()=>{const i=e.selectedOption;if(i)return i[e.labelField]}),w=B(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function N(){var i;const{value:h}=t;if(h){const{value:E}=o;E&&(E.style.width=`${h.offsetWidth}px`,e.maxTagCount!=="responsive"&&((i=b.value)===null||i===void 0||i.sync()))}}function $(){const{value:i}=g;i&&(i.style.display="none")}function W(){const{value:i}=g;i&&(i.style.display="inline-block")}Se(H(e,"active"),i=>{i||$()}),Se(H(e,"pattern"),()=>{e.multiple&&tn(N)});function A(i){const{onFocus:h}=e;h&&h(i)}function j(i){const{onBlur:h}=e;h&&h(i)}function K(i){const{onDeleteOption:h}=e;h&&h(i)}function ce(i){const{onClear:h}=e;h&&h(i)}function Q(i){const{onPatternInput:h}=e;h&&h(i)}function q(i){var h;(!i.relatedTarget||!(!((h=s.value)===null||h===void 0)&&h.contains(i.relatedTarget)))&&A(i)}function V(i){var h;!((h=s.value)===null||h===void 0)&&h.contains(i.relatedTarget)||j(i)}function ue(i){ce(i)}function ge(){P.value=!0}function fe(){P.value=!1}function ie(i){!e.active||!e.filterable||i.target!==o.value&&i.preventDefault()}function Y(i){K(i)}function ae(i){if(i.key==="Backspace"&&!l.value&&!e.pattern.length){const{selectedOptions:h}=e;h!=null&&h.length&&Y(h[h.length-1])}}const l=k(!1);let c=null;function L(i){const{value:h}=t;if(h){const E=i.target.value;h.textContent=E,N()}e.ignoreComposition&&l.value?c=i:Q(i)}function le(){l.value=!0}function Oe(){l.value=!1,e.ignoreComposition&&Q(c),c=null}function Re(i){var h;C.value=!0,(h=e.onPatternFocus)===null||h===void 0||h.call(e,i)}function Fe(i){var h;C.value=!1,(h=e.onPatternBlur)===null||h===void 0||h.call(e,i)}function be(){var i,h;if(e.filterable)C.value=!1,(i=v.value)===null||i===void 0||i.blur(),(h=o.value)===null||h===void 0||h.blur();else if(e.multiple){const{value:E}=f;E==null||E.blur()}else{const{value:E}=u;E==null||E.blur()}}function me(){var i,h,E;e.filterable?(C.value=!1,(i=v.value)===null||i===void 0||i.focus()):e.multiple?(h=f.value)===null||h===void 0||h.focus():(E=u.value)===null||E===void 0||E.focus()}function re(){const{value:i}=o;i&&(W(),i.focus())}function U(){const{value:i}=o;i&&i.blur()}function we(i){const{value:h}=a;h&&h.setTextContent(`+${i}`)}function se(){const{value:i}=R;return i}function ze(){return o.value}let he=null;function ve(){he!==null&&window.clearTimeout(he)}function Me(){e.active||(ve(),he=window.setTimeout(()=>{w.value&&(x.value=!0)},100))}function Te(){ve()}function Pe(i){i||(ve(),x.value=!1)}Se(w,i=>{i||(x.value=!1)}),Ne(()=>{tt(()=>{const i=v.value;i&&(e.disabled?i.removeAttribute("tabindex"):i.tabIndex=C.value?-1:0)})}),wn(s,e.onResize);const{inlineThemeDisabled:ye}=e,xe=B(()=>{const{size:i}=e,{common:{cubicBezierEaseInOut:h},self:{borderRadius:E,color:ke,placeholderColor:De,textColor:Ve,paddingSingle:je,paddingMultiple:We,caretColor:Be,colorDisabled:Ie,textColorDisabled:_e,placeholderColorDisabled:He,colorActive:Ke,boxShadowFocus:$e,boxShadowActive:de,boxShadowHover:n,border:d,borderFocus:y,borderHover:M,borderActive:F,arrowColor:O,arrowColorDisabled:z,loadingColor:D,colorActiveWarning:X,boxShadowFocusWarning:Ue,boxShadowActiveWarning:xn,boxShadowHoverWarning:Cn,borderWarning:Sn,borderFocusWarning:On,borderHoverWarning:Rn,borderActiveWarning:Fn,colorActiveError:zn,boxShadowFocusError:Mn,boxShadowActiveError:Tn,boxShadowHoverError:Pn,borderError:kn,borderFocusError:Bn,borderHoverError:In,borderActiveError:_n,clearColor:$n,clearColorHover:An,clearColorPressed:En,clearSize:Nn,arrowSize:Ln,[ne("height",i)]:Dn,[ne("fontSize",i)]:Vn}}=p.value;return{"--n-bezier":h,"--n-border":d,"--n-border-active":F,"--n-border-focus":y,"--n-border-hover":M,"--n-border-radius":E,"--n-box-shadow-active":de,"--n-box-shadow-focus":$e,"--n-box-shadow-hover":n,"--n-caret-color":Be,"--n-color":ke,"--n-color-active":Ke,"--n-color-disabled":Ie,"--n-font-size":Vn,"--n-height":Dn,"--n-padding-single":je,"--n-padding-multiple":We,"--n-placeholder-color":De,"--n-placeholder-color-disabled":He,"--n-text-color":Ve,"--n-text-color-disabled":_e,"--n-arrow-color":O,"--n-arrow-color-disabled":z,"--n-loading-color":D,"--n-color-active-warning":X,"--n-box-shadow-focus-warning":Ue,"--n-box-shadow-active-warning":xn,"--n-box-shadow-hover-warning":Cn,"--n-border-warning":Sn,"--n-border-focus-warning":On,"--n-border-hover-warning":Rn,"--n-border-active-warning":Fn,"--n-color-active-error":zn,"--n-box-shadow-focus-error":Mn,"--n-box-shadow-active-error":Tn,"--n-box-shadow-hover-error":Pn,"--n-border-error":kn,"--n-border-focus-error":Bn,"--n-border-hover-error":In,"--n-border-active-error":_n,"--n-clear-size":Nn,"--n-clear-color":$n,"--n-clear-color-hover":An,"--n-clear-color-pressed":En,"--n-arrow-size":Ln}}),Z=ye?Le("internal-selection",B(()=>e.size[0]),xe,e):void 0;return{mergedTheme:p,mergedClearable:S,patternInputFocused:C,filterablePlaceholder:I,label:m,selected:w,showTagsPanel:x,isComposing:l,counterRef:a,counterWrapperRef:R,patternInputMirrorRef:t,patternInputRef:o,selfRef:s,multipleElRef:f,singleElRef:u,patternInputWrapperRef:v,overflowRef:b,inputTagElRef:g,handleMouseDown:ie,handleFocusin:q,handleClear:ue,handleMouseEnter:ge,handleMouseLeave:fe,handleDeleteOption:Y,handlePatternKeyDown:ae,handlePatternInputInput:L,handlePatternInputBlur:Fe,handlePatternInputFocus:Re,handleMouseEnterCounter:Me,handleMouseLeaveCounter:Te,handleFocusout:V,handleCompositionEnd:Oe,handleCompositionStart:le,onPopoverUpdateShow:Pe,focus:me,focusInput:re,blur:be,blurInput:U,updateCounter:we,getCounter:se,getTail:ze,renderLabel:e.renderLabel,cssVars:ye?void 0:xe,themeClass:Z==null?void 0:Z.themeClass,onRender:Z==null?void 0:Z.onRender}},render(){const{status:e,multiple:t,size:o,disabled:s,filterable:f,maxTagCount:u,bordered:v,clsPrefix:a,onRender:R,renderTag:b,renderLabel:g}=this;R==null||R();const x=u==="responsive",C=typeof u=="number",P=x||C,p=r(lt,null,{default:()=>r(St,{clsPrefix:a,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var I,m;return(m=(I=this.$slots).arrow)===null||m===void 0?void 0:m.call(I)}})});let S;if(t){const{labelField:I}=this,m=V=>r("div",{class:`${a}-base-selection-tag-wrapper`,key:V.value},b?b({option:V,handleClose:()=>{this.handleDeleteOption(V)}}):r(Ze,{size:o,closable:!V.disabled,disabled:s,onClose:()=>{this.handleDeleteOption(V)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>g?g(V,!0):Ce(V[I],V,!0)})),w=()=>(C?this.selectedOptions.slice(0,u):this.selectedOptions).map(m),N=f?r("div",{class:`${a}-base-selection-input-tag`,ref:"inputTagElRef",key:"__input-tag__"},r("input",Object.assign({},this.inputProps,{ref:"patternInputRef",tabindex:-1,disabled:s,value:this.pattern,autofocus:this.autofocus,class:`${a}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),r("span",{ref:"patternInputMirrorRef",class:`${a}-base-selection-input-tag__mirror`},this.pattern)):null,$=x?()=>r("div",{class:`${a}-base-selection-tag-wrapper`,ref:"counterWrapperRef"},r(Ze,{size:o,ref:"counterRef",onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:s})):void 0;let W;if(C){const V=this.selectedOptions.length-u;V>0&&(W=r("div",{class:`${a}-base-selection-tag-wrapper`,key:"__counter__"},r(Ze,{size:o,ref:"counterRef",onMouseenter:this.handleMouseEnterCounter,disabled:s},{default:()=>`+${V}`})))}const A=x?f?r(un,{ref:"overflowRef",updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:"100%",display:"flex",overflow:"hidden"}},{default:w,counter:$,tail:()=>N}):r(un,{ref:"overflowRef",updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:"100%",display:"flex",overflow:"hidden"}},{default:w,counter:$}):C?w().concat(W):w(),j=P?()=>r("div",{class:`${a}-base-selection-popover`},x?w():this.selectedOptions.map(m)):void 0,K=P?{show:this.showTagsPanel,trigger:"hover",overlap:!0,placement:"top",width:"trigger",onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover}:null,Q=(this.selected?!1:this.active?!this.pattern&&!this.isComposing:!0)?r("div",{class:`${a}-base-selection-placeholder ${a}-base-selection-overlay`},r("div",{class:`${a}-base-selection-placeholder__inner`},this.placeholder)):null,q=f?r("div",{ref:"patternInputWrapperRef",class:`${a}-base-selection-tags`},A,x?null:N,p):r("div",{ref:"multipleElRef",class:`${a}-base-selection-tags`,tabindex:s?void 0:0},A,p);S=r(it,null,P?r(ot,Object.assign({},K,{scrollable:!0,style:"max-height: calc(var(--v-target-height) * 6.6);"}),{trigger:()=>q,default:j}):q,Q)}else if(f){const I=this.pattern||this.isComposing,m=this.active?!I:!this.selected,w=this.active?!1:this.selected;S=r("div",{ref:"patternInputWrapperRef",class:`${a}-base-selection-label`},r("input",Object.assign({},this.inputProps,{ref:"patternInputRef",class:`${a}-base-selection-input`,value:this.active?this.pattern:"",placeholder:"",readonly:s,disabled:s,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),w?r("div",{class:`${a}-base-selection-label__render-label ${a}-base-selection-overlay`,key:"input"},r("div",{class:`${a}-base-selection-overlay__wrapper`},b?b({option:this.selectedOption,handleClose:()=>{}}):g?g(this.selectedOption,!0):Ce(this.label,this.selectedOption,!0))):null,m?r("div",{class:`${a}-base-selection-placeholder ${a}-base-selection-overlay`,key:"placeholder"},r("div",{class:`${a}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,p)}else S=r("div",{ref:"singleElRef",class:`${a}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label!==void 0?r("div",{class:`${a}-base-selection-input`,title:zt(this.label),key:"input"},r("div",{class:`${a}-base-selection-input__content`},b?b({option:this.selectedOption,handleClose:()=>{}}):g?g(this.selectedOption,!0):Ce(this.label,this.selectedOption,!0))):r("div",{class:`${a}-base-selection-placeholder ${a}-base-selection-overlay`,key:"placeholder"},r("div",{class:`${a}-base-selection-placeholder__inner`},this.placeholder)),p);return r("div",{ref:"selfRef",class:[`${a}-base-selection`,this.themeClass,e&&`${a}-base-selection--${e}-status`,{[`${a}-base-selection--active`]:this.active,[`${a}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${a}-base-selection--disabled`]:this.disabled,[`${a}-base-selection--multiple`]:this.multiple,[`${a}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},S,v?r("div",{class:`${a}-base-selection__border`}):null,v?r("div",{class:`${a}-base-selection__state-border`}):null)}});function Ee(e){return e.type==="group"}function yn(e){return e.type==="ignored"}function Je(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function Lt(e,t){return{getIsGroup:Ee,getIgnored:yn,getKey(s){return Ee(s)?s.name||s.key||"key-required":s[e]},getChildren(s){return s[t]}}}function Dt(e,t,o,s){if(!t)return e;function f(u){if(!Array.isArray(u))return[];const v=[];for(const a of u)if(Ee(a)){const R=f(a[s]);R.length&&v.push(Object.assign({},a,{[s]:R}))}else{if(yn(a))continue;t(o,a)&&v.push(a)}return v}return f(e)}function Vt(e,t,o){const s=new Map;return e.forEach(f=>{Ee(f)?f[o].forEach(u=>{s.set(u[t],u)}):s.set(f[t],f)}),s}const jt=J([T("select",`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 `),T("select-menu",`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[bn({originalTransition:"background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)"})])]),Wt=Object.assign(Object.assign({},te.props),{to:nn.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:"bottom-start"},widthMode:{type:String,default:"trigger"},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:"label"},valueField:{type:String,default:"value"},childrenField:{type:String,default:"children"},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:"show"},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},onChange:[Function,Array],items:Array}),eo=oe({name:"Select",props:Wt,setup(e){const{mergedClsPrefixRef:t,mergedBorderedRef:o,namespaceRef:s,inlineThemeDisabled:f}=on(e),u=te("Select","-select",jt,pt,e,t),v=k(e.defaultValue),a=H(e,"value"),R=dn(a,v),b=k(!1),g=k(""),x=B(()=>{const{valueField:n,childrenField:d}=e,y=Lt(n,d);return rt(V.value,y)}),C=B(()=>Vt(Q.value,e.valueField,e.childrenField)),P=k(!1),p=dn(H(e,"show"),P),S=k(null),I=k(null),m=k(null),{localeRef:w}=mn("Select"),N=B(()=>{var n;return(n=e.placeholder)!==null&&n!==void 0?n:w.value.placeholder}),$=at(e,["items","options"]),W=[],A=k([]),j=k([]),K=k(new Map),ce=B(()=>{const{fallbackOption:n}=e;if(n===void 0){const{labelField:d,valueField:y}=e;return M=>({[d]:String(M),[y]:M})}return n===!1?!1:d=>Object.assign(n(d),{value:d})}),Q=B(()=>j.value.concat(A.value).concat($.value)),q=B(()=>{const{filter:n}=e;if(n)return n;const{labelField:d,valueField:y}=e;return(M,F)=>{if(!F)return!1;const O=F[d];if(typeof O=="string")return Je(M,O);const z=F[y];return typeof z=="string"?Je(M,z):typeof z=="number"?Je(M,String(z)):!1}}),V=B(()=>{if(e.remote)return $.value;{const{value:n}=Q,{value:d}=g;return!d.length||!e.filterable?n:Dt(n,q.value,d,e.childrenField)}});function ue(n){const d=e.remote,{value:y}=K,{value:M}=C,{value:F}=ce,O=[];return n.forEach(z=>{if(M.has(z))O.push(M.get(z));else if(d&&y.has(z))O.push(y.get(z));else if(F){const D=F(z);D&&O.push(D)}}),O}const ge=B(()=>{if(e.multiple){const{value:n}=R;return Array.isArray(n)?ue(n):[]}return null}),fe=B(()=>{const{value:n}=R;return!e.multiple&&!Array.isArray(n)?n===null?null:ue([n])[0]||null:null}),ie=st(e),{mergedSizeRef:Y,mergedDisabledRef:ae,mergedStatusRef:l}=ie;function c(n,d){const{onChange:y,"onUpdate:value":M,onUpdateValue:F}=e,{nTriggerFormChange:O,nTriggerFormInput:z}=ie;y&&ee(y,n,d),F&&ee(F,n,d),M&&ee(M,n,d),v.value=n,O(),z()}function L(n){const{onBlur:d}=e,{nTriggerFormBlur:y}=ie;d&&ee(d,n),y()}function le(){const{onClear:n}=e;n&&ee(n)}function Oe(n){const{onFocus:d,showOnFocus:y}=e,{nTriggerFormFocus:M}=ie;d&&ee(d,n),M(),y&&re()}function Re(n){const{onSearch:d}=e;d&&ee(d,n)}function Fe(n){const{onScroll:d}=e;d&&ee(d,n)}function be(){var n;const{remote:d,multiple:y}=e;if(d){const{value:M}=K;if(y){const{valueField:F}=e;(n=ge.value)===null||n===void 0||n.forEach(O=>{M.set(O[F],O)})}else{const F=fe.value;F&&M.set(F[e.valueField],F)}}}function me(n){const{onUpdateShow:d,"onUpdate:show":y}=e;d&&ee(d,n),y&&ee(y,n),P.value=n}function re(){ae.value||(me(!0),P.value=!0,e.filterable&&_e())}function U(){me(!1)}function we(){g.value="",j.value=W}const se=k(!1);function ze(){e.filterable&&(se.value=!0)}function he(){e.filterable&&(se.value=!1,p.value||we())}function ve(){ae.value||(p.value?e.filterable?_e():U():re())}function Me(n){var d,y;!((y=(d=m.value)===null||d===void 0?void 0:d.selfRef)===null||y===void 0)&&y.contains(n.relatedTarget)||(b.value=!1,L(n),U())}function Te(n){Oe(n),b.value=!0}function Pe(n){b.value=!0}function ye(n){var d;!((d=S.value)===null||d===void 0)&&d.$el.contains(n.relatedTarget)||(b.value=!1,L(n),U())}function xe(){var n;(n=S.value)===null||n===void 0||n.focus(),U()}function Z(n){var d;p.value&&(!((d=S.value)===null||d===void 0)&&d.$el.contains(gt(n))||U())}function i(n){if(!Array.isArray(n))return[];if(ce.value)return Array.from(n);{const{remote:d}=e,{value:y}=C;if(d){const{value:M}=K;return n.filter(F=>y.has(F)||M.has(F))}else return n.filter(M=>y.has(M))}}function h(n){E(n.rawNode)}function E(n){if(ae.value)return;const{tag:d,remote:y,clearFilterAfterSelect:M,valueField:F}=e;if(d&&!y){const{value:O}=j,z=O[0]||null;if(z){const D=A.value;D.length?D.push(z):A.value=[z],j.value=W}}if(y&&K.value.set(n[F],n),e.multiple){const O=i(R.value),z=O.findIndex(D=>D===n[F]);if(~z){if(O.splice(z,1),d&&!y){const D=ke(n[F]);~D&&(A.value.splice(D,1),M&&(g.value=""))}}else O.push(n[F]),M&&(g.value="");c(O,ue(O))}else{if(d&&!y){const O=ke(n[F]);~O?A.value=[A.value[O]]:A.value=W}Ie(),U(),c(n[F],n)}}function ke(n){return A.value.findIndex(y=>y[e.valueField]===n)}function De(n){p.value||re();const{value:d}=n.target;g.value=d;const{tag:y,remote:M}=e;if(Re(d),y&&!M){if(!d){j.value=W;return}const{onCreate:F}=e,O=F?F(d):{[e.labelField]:d,[e.valueField]:d},{valueField:z,labelField:D}=e;$.value.some(X=>X[z]===O[z]||X[D]===O[D])||A.value.some(X=>X[z]===O[z]||X[D]===O[D])?j.value=W:j.value=[O]}}function Ve(n){n.stopPropagation();const{multiple:d}=e;!d&&e.filterable&&U(),le(),d?c([],[]):c(null,null)}function je(n){!Ae(n,"action")&&!Ae(n,"empty")&&n.preventDefault()}function We(n){Fe(n)}function Be(n){var d,y,M,F,O;if(!e.keyboard){n.preventDefault();return}switch(n.key){case" ":if(e.filterable)break;n.preventDefault();case"Enter":if(!(!((d=S.value)===null||d===void 0)&&d.isComposing)){if(p.value){const z=(y=m.value)===null||y===void 0?void 0:y.getPendingTmNode();z?h(z):e.filterable||(U(),Ie())}else if(re(),e.tag&&se.value){const z=j.value[0];if(z){const D=z[e.valueField],{value:X}=R;e.multiple&&Array.isArray(X)&&X.some(Ue=>Ue===D)||E(z)}}}n.preventDefault();break;case"ArrowUp":if(n.preventDefault(),e.loading)return;p.value&&((M=m.value)===null||M===void 0||M.prev());break;case"ArrowDown":if(n.preventDefault(),e.loading)return;p.value?(F=m.value)===null||F===void 0||F.next():re();break;case"Escape":p.value&&(bt(n),U()),(O=S.value)===null||O===void 0||O.focus();break}}function Ie(){var n;(n=S.value)===null||n===void 0||n.focus()}function _e(){var n;(n=S.value)===null||n===void 0||n.focusInput()}function He(){var n;p.value&&((n=I.value)===null||n===void 0||n.syncPosition())}be(),Se(H(e,"options"),be);const Ke={focus:()=>{var n;(n=S.value)===null||n===void 0||n.focus()},focusInput:()=>{var n;(n=S.value)===null||n===void 0||n.focusInput()},blur:()=>{var n;(n=S.value)===null||n===void 0||n.blur()},blurInput:()=>{var n;(n=S.value)===null||n===void 0||n.blurInput()}},$e=B(()=>{const{self:{menuBoxShadow:n}}=u.value;return{"--n-menu-box-shadow":n}}),de=f?Le("select",void 0,$e,e):void 0;return Object.assign(Object.assign({},Ke),{mergedStatus:l,mergedClsPrefix:t,mergedBordered:o,namespace:s,treeMate:x,isMounted:dt(),triggerRef:S,menuRef:m,pattern:g,uncontrolledShow:P,mergedShow:p,adjustedTo:nn(e),uncontrolledValue:v,mergedValue:R,followerRef:I,localizedPlaceholder:N,selectedOption:fe,selectedOptions:ge,mergedSize:Y,mergedDisabled:ae,focused:b,activeWithoutMenuOpen:se,inlineThemeDisabled:f,onTriggerInputFocus:ze,onTriggerInputBlur:he,handleTriggerOrMenuResize:He,handleMenuFocus:Pe,handleMenuBlur:ye,handleMenuTabOut:xe,handleTriggerClick:ve,handleToggle:h,handleDeleteOption:E,handlePatternInput:De,handleClear:Ve,handleTriggerBlur:Me,handleTriggerFocus:Te,handleKeydown:Be,handleMenuAfterLeave:we,handleMenuClickOutside:Z,handleMenuScroll:We,handleMenuKeydown:Be,handleMenuMousedown:je,mergedTheme:u,cssVars:f?void 0:$e,themeClass:de==null?void 0:de.themeClass,onRender:de==null?void 0:de.onRender})},render(){return r("div",{class:`${this.mergedClsPrefix}-select`},r(ct,null,{default:()=>[r(ut,null,{default:()=>r(Nt,{ref:"triggerRef",inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e,t;return[(t=(e=this.$slots).arrow)===null||t===void 0?void 0:t.call(e)]}})}),r(ft,{ref:"followerRef",show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===nn.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?"target":void 0,minWidth:"target",placement:this.placement},{default:()=>r(gn,{name:"fade-in-scale-up-transition",appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e,t,o;return this.mergedShow||this.displayDirective==="show"?((e=this.onRender)===null||e===void 0||e.call(this),ht(r(At,Object.assign({},this.menuProps,{ref:"menuRef",onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,(t=this.menuProps)===null||t===void 0?void 0:t.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:"medium",renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[(o=this.menuProps)===null||o===void 0?void 0:o.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange}),{empty:()=>{var s,f;return[(f=(s=this.$slots).empty)===null||f===void 0?void 0:f.call(s)]},action:()=>{var s,f;return[(f=(s=this.$slots).action)===null||f===void 0?void 0:f.call(s)]}}),this.displayDirective==="show"?[[vt,this.mergedShow],[cn,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[cn,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}}),Ht=()=>mt,Kt={name:"Space",self:Ht},Ut=Kt;let Qe;const Gt=()=>{if(!wt)return!0;if(Qe===void 0){const e=document.createElement("div");e.style.display="flex",e.style.flexDirection="column",e.style.rowGap="1px",e.appendChild(document.createElement("div")),e.appendChild(document.createElement("div")),document.body.appendChild(e);const t=e.scrollHeight===1;return document.body.removeChild(e),Qe=t}return Qe},qt=Object.assign(Object.assign({},te.props),{align:String,justify:{type:String,default:"start"},inline:Boolean,vertical:Boolean,size:{type:[String,Number,Array],default:"medium"},wrapItem:{type:Boolean,default:!0},itemStyle:[String,Object],wrap:{type:Boolean,default:!0},internalUseGap:{type:Boolean,default:void 0}}),no=oe({name:"Space",props:qt,setup(e){const{mergedClsPrefixRef:t,mergedRtlRef:o}=on(e),s=te("Space","-space",void 0,Ut,e,t),f=yt("Space",o,t);return{useGap:Gt(),rtlEnabled:f,mergedClsPrefix:t,margin:B(()=>{const{size:u}=e;if(Array.isArray(u))return{horizontal:u[0],vertical:u[1]};if(typeof u=="number")return{horizontal:u,vertical:u};const{self:{[ne("gap",u)]:v}}=s.value,{row:a,col:R}=xt(v);return{horizontal:en(R),vertical:en(a)}})}},render(){const{vertical:e,align:t,inline:o,justify:s,itemStyle:f,margin:u,wrap:v,mergedClsPrefix:a,rtlEnabled:R,useGap:b,wrapItem:g,internalUseGap:x}=this,C=Ct(Ft(this));if(!C.length)return null;const P=`${u.horizontal}px`,p=`${u.horizontal/2}px`,S=`${u.vertical}px`,I=`${u.vertical/2}px`,m=C.length-1,w=s.startsWith("space-");return r("div",{role:"none",class:[`${a}-space`,R&&`${a}-space--rtl`],style:{display:o?"inline-flex":"flex",flexDirection:e?"column":"row",justifyContent:["start","end"].includes(s)?"flex-"+s:s,flexWrap:!v||e?"nowrap":"wrap",marginTop:b||e?"":`-${I}`,marginBottom:b||e?"":`-${I}`,alignItems:t,gap:b?`${u.vertical}px ${u.horizontal}px`:""}},!g&&(b||x)?C:C.map((N,$)=>r("div",{role:"none",style:[f,{maxWidth:"100%"},b?"":e?{marginBottom:$!==m?S:""}:R?{marginLeft:w?s==="space-between"&&$===m?"":p:$!==m?P:"",marginRight:w?s==="space-between"&&$===0?"":p:"",paddingTop:I,paddingBottom:I}:{marginRight:w?s==="space-between"&&$===m?"":p:$!==m?P:"",marginLeft:w?s==="space-between"&&$===0?"":p:"",paddingTop:I,paddingBottom:I}]},N)))}});export{eo as N,Xt as a,no as b,It as c,At as d,Lt as e,Ye as m,Ut as s};
