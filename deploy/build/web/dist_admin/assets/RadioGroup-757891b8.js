import{M as X,bv as ve,bw as ge,x as me,D as pe,c2 as xe,T as A,k as E,a4 as te,ax as ie,ao as G,z as I,aN as ae,bx as Z,aQ as we,aP as Re,ag as Se,P as de,ah as ye,R as oe,S as ue,O as re,Y as W,L as O,aq as _,as as S,a3 as D,at as ne,Q as q,aK as Y,a5 as ce,U as fe,a7 as ze,cl as he,aA as ke,aR as Ce}from"./index-7e9df7d9.js";import{g as Te}from"./get-slot-1efb97e5.js";function le(e){return e&-e}class Ie{constructor(t,n){this.l=t,this.min=n;const a=new Array(t+1);for(let i=0;i<t+1;++i)a[i]=0;this.ft=a}add(t,n){if(n===0)return;const{l:a,ft:i}=this;for(t+=1;t<=a;)i[t]+=n,t+=le(t)}get(t){return this.sum(t+1)-this.sum(t)}sum(t){if(t===void 0&&(t=this.l),t<=0)return 0;const{ft:n,min:a,l:i}=this;if(t>i)throw new Error("[FinweckTree.sum]: `i` is larger than length.");let s=t*a;for(;t>0;)s+=n[t],t-=le(t);return s}getBound(t){let n=0,a=this.l;for(;a>n;){const i=Math.floor((n+a)/2),s=this.sum(i);if(s>t){a=i;continue}else if(s<t){if(n===i)return this.sum(n+1)<=t?n+1:i;n=i}else return i}return n}}let K;function Be(){return K===void 0&&("matchMedia"in window?K=window.matchMedia("(pointer:coarse)").matches:K=!1),K}let ee;function se(){return ee===void 0&&(ee="chrome"in window?window.devicePixelRatio:1),ee}const Fe=Z(".v-vl",{maxHeight:"inherit",height:"100%",overflow:"auto",minWidth:"1px"},[Z("&:not(.v-vl--show-scrollbar)",{scrollbarWidth:"none"},[Z("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",{width:0,height:0,display:"none"})])]),Pe=X({name:"VirtualList",inheritAttrs:!1,props:{showScrollbar:{type:Boolean,default:!0},items:{type:Array,default:()=>[]},itemSize:{type:Number,required:!0},itemResizable:Boolean,itemsStyle:[String,Object],visibleItemsTag:{type:[String,Object],default:"div"},visibleItemsProps:Object,ignoreItemResize:Boolean,onScroll:Function,onWheel:Function,onResize:Function,defaultScrollKey:[Number,String],defaultScrollIndex:Number,keyField:{type:String,default:"key"},paddingTop:{type:[Number,String],default:0},paddingBottom:{type:[Number,String],default:0}},setup(e){const t=ve();Fe.mount({id:"vueuc/virtual-list",head:!0,anchorMetaName:ge,ssr:t}),me(()=>{const{defaultScrollIndex:o,defaultScrollKey:r}=e;o!=null?k({index:o}):r!=null&&k({key:r})});let n=!1,a=!1;pe(()=>{if(n=!1,!a){a=!0;return}k({top:B.value,left:y})}),xe(()=>{n=!0,a||(a=!0)});const i=A(()=>{const o=new Map,{keyField:r}=e;return e.items.forEach((d,b)=>{o.set(d[r],b)}),o}),s=E(null),l=E(void 0),f=new Map,c=A(()=>{const{items:o,itemSize:r,keyField:d}=e,b=new Ie(o.length,r);return o.forEach((m,p)=>{const h=m[d],R=f.get(h);R!==void 0&&b.add(p,R)}),b}),v=E(0);let y=0;const B=E(0),z=te(()=>Math.max(c.value.getBound(B.value-ie(e.paddingTop))-1,0)),V=A(()=>{const{value:o}=l;if(o===void 0)return[];const{items:r,itemSize:d}=e,b=z.value,m=Math.min(b+Math.ceil(o/d+1),r.length-1),p=[];for(let h=b;h<=m;++h)p.push(r[h]);return p}),k=(o,r)=>{if(typeof o=="number"){u(o,r,"auto");return}const{left:d,top:b,index:m,key:p,position:h,behavior:R,debounce:H=!0}=o;if(d!==void 0||b!==void 0)u(d,b,R);else if(m!==void 0)$(m,R,H);else if(p!==void 0){const Q=i.value.get(p);Q!==void 0&&$(Q,R,H)}else h==="bottom"?u(0,Number.MAX_SAFE_INTEGER,R):h==="top"&&u(0,0,R)};let T,F=null;function $(o,r,d){const{value:b}=c,m=b.sum(o)+ie(e.paddingTop);if(!d)s.value.scrollTo({left:0,top:m,behavior:r});else{T=o,F!==null&&window.clearTimeout(F),F=window.setTimeout(()=>{T=void 0,F=null},16);const{scrollTop:p,offsetHeight:h}=s.value;if(m>p){const R=b.get(o);m+R<=p+h||s.value.scrollTo({left:0,top:m+R-h,behavior:r})}else s.value.scrollTo({left:0,top:m,behavior:r})}}function u(o,r,d){s.value.scrollTo({left:o,top:r,behavior:d})}function x(o,r){var d,b,m;if(n||e.ignoreItemResize||U(r.target))return;const{value:p}=c,h=i.value.get(o),R=p.get(h),H=(m=(b=(d=r.borderBoxSize)===null||d===void 0?void 0:d[0])===null||b===void 0?void 0:b.blockSize)!==null&&m!==void 0?m:r.contentRect.height;if(H===R)return;H-e.itemSize===0?f.delete(o):f.set(o,H-e.itemSize);const j=H-R;if(j===0)return;p.add(h,j);const P=s.value;if(P!=null){if(T===void 0){const J=p.sum(h);P.scrollTop>J&&P.scrollBy(0,j)}else if(h<T)P.scrollBy(0,j);else if(h===T){const J=p.sum(h);H+J>P.scrollTop+P.offsetHeight&&P.scrollBy(0,j)}N()}v.value++}const w=!Be();let g=!1;function C(o){var r;(r=e.onScroll)===null||r===void 0||r.call(e,o),(!w||!g)&&N()}function M(o){var r;if((r=e.onWheel)===null||r===void 0||r.call(e,o),w){const d=s.value;if(d!=null){if(o.deltaX===0&&(d.scrollTop===0&&o.deltaY<=0||d.scrollTop+d.offsetHeight>=d.scrollHeight&&o.deltaY>=0))return;o.preventDefault(),d.scrollTop+=o.deltaY/se(),d.scrollLeft+=o.deltaX/se(),N(),g=!0,we(()=>{g=!1})}}}function L(o){if(n||U(o.target)||o.contentRect.height===l.value)return;l.value=o.contentRect.height;const{onResize:r}=e;r!==void 0&&r(o)}function N(){const{value:o}=s;o!=null&&(B.value=o.scrollTop,y=o.scrollLeft)}function U(o){let r=o;for(;r!==null;){if(r.style.display==="none")return!0;r=r.parentElement}return!1}return{listHeight:l,listStyle:{overflow:"auto"},keyToIndex:i,itemsStyle:A(()=>{const{itemResizable:o}=e,r=G(c.value.sum());return v.value,[e.itemsStyle,{boxSizing:"content-box",height:o?"":r,minHeight:o?r:"",paddingTop:G(e.paddingTop),paddingBottom:G(e.paddingBottom)}]}),visibleItemsStyle:A(()=>(v.value,{transform:`translateY(${G(c.value.sum(z.value))})`})),viewportItems:V,listElRef:s,itemsElRef:E(null),scrollTo:k,handleListResize:L,handleListScroll:C,handleListWheel:M,handleItemResize:x}},render(){const{itemResizable:e,keyField:t,keyToIndex:n,visibleItemsTag:a}=this;return I(ae,{onResize:this.handleListResize},{default:()=>{var i,s;return I("div",Re(this.$attrs,{class:["v-vl",this.showScrollbar&&"v-vl--show-scrollbar"],onScroll:this.handleListScroll,onWheel:this.handleListWheel,ref:"listElRef"}),[this.items.length!==0?I("div",{ref:"itemsElRef",class:"v-vl-items",style:this.itemsStyle},[I(a,Object.assign({class:"v-vl-visible-items",style:this.visibleItemsStyle},this.visibleItemsProps),{default:()=>this.viewportItems.map(l=>{const f=l[t],c=n.get(f),v=this.$slots.default({item:l,index:c})[0];return e?I(ae,{key:f,onResize:y=>this.handleItemResize(f,y)},{default:()=>v}):(v.key=f,v)})})]):(s=(i=this.$slots).empty)===null||s===void 0?void 0:s.call(i)])}})}}),Oe=X({props:{onFocus:Function,onBlur:Function},setup(e){return()=>I("div",{style:"width: 0; height: 0",tabindex:0,onFocus:e.onFocus,onBlur:e.onBlur})}}),_e={name:String,value:{type:[String,Number,Boolean],default:"on"},checked:{type:Boolean,default:void 0},defaultChecked:Boolean,disabled:{type:Boolean,default:void 0},label:String,size:String,onUpdateChecked:[Function,Array],"onUpdate:checked":[Function,Array],checkedValue:{type:Boolean,default:void 0}},be=Se("n-radio-group");function Ee(e){const t=de(e,{mergedSize(u){const{size:x}=e;if(x!==void 0)return x;if(l){const{mergedSizeRef:{value:w}}=l;if(w!==void 0)return w}return u?u.mergedSize.value:"medium"},mergedDisabled(u){return!!(e.disabled||l!=null&&l.disabledRef.value||u!=null&&u.disabled.value)}}),{mergedSizeRef:n,mergedDisabledRef:a}=t,i=E(null),s=E(null),l=ye(be,null),f=E(e.defaultChecked),c=oe(e,"checked"),v=ue(c,f),y=te(()=>l?l.valueRef.value===e.value:v.value),B=te(()=>{const{name:u}=e;if(u!==void 0)return u;if(l)return l.nameRef.value}),z=E(!1);function V(){if(l){const{doUpdateValue:u}=l,{value:x}=e;W(u,x)}else{const{onUpdateChecked:u,"onUpdate:checked":x}=e,{nTriggerFormInput:w,nTriggerFormChange:g}=t;u&&W(u,!0),x&&W(x,!0),w(),g(),f.value=!0}}function k(){a.value||y.value||V()}function T(){k()}function F(){z.value=!1}function $(){z.value=!0}return{mergedClsPrefix:l?l.mergedClsPrefixRef:re(e).mergedClsPrefixRef,inputRef:i,labelRef:s,mergedName:B,mergedDisabled:a,uncontrolledChecked:f,renderSafeChecked:y,focus:z,mergedSize:n,handleRadioInputChange:T,handleRadioInputBlur:F,handleRadioInputFocus:$}}const Ve=O("radio",`
 line-height: var(--n-label-line-height);
 outline: none;
 position: relative;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 align-items: flex-start;
 flex-wrap: nowrap;
 font-size: var(--n-font-size);
 word-break: break-word;
`,[_("checked",[S("dot",`
 background-color: var(--n-color-active);
 `)]),S("dot-wrapper",`
 position: relative;
 flex-shrink: 0;
 flex-grow: 0;
 width: var(--n-radio-size);
 `),O("radio-input",`
 position: absolute;
 border: 0;
 border-radius: inherit;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 opacity: 0;
 z-index: 1;
 cursor: pointer;
 `),S("dot",`
 position: absolute;
 top: 50%;
 left: 0;
 transform: translateY(-50%);
 height: var(--n-radio-size);
 width: var(--n-radio-size);
 background: var(--n-color);
 box-shadow: var(--n-box-shadow);
 border-radius: 50%;
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `,[D("&::before",`
 content: "";
 opacity: 0;
 position: absolute;
 left: 4px;
 top: 4px;
 height: calc(100% - 8px);
 width: calc(100% - 8px);
 border-radius: 50%;
 transform: scale(.8);
 background: var(--n-dot-color-active);
 transition: 
 opacity .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),_("checked",{boxShadow:"var(--n-box-shadow-active)"},[D("&::before",`
 opacity: 1;
 transform: scale(1);
 `)])]),S("label",`
 color: var(--n-text-color);
 padding: var(--n-label-padding);
 font-weight: var(--n-label-font-weight);
 display: inline-block;
 transition: color .3s var(--n-bezier);
 `),ne("disabled",`
 cursor: pointer;
 `,[D("&:hover",[S("dot",{boxShadow:"var(--n-box-shadow-hover)"})]),_("focus",[D("&:not(:active)",[S("dot",{boxShadow:"var(--n-box-shadow-focus)"})])])]),_("disabled",`
 cursor: not-allowed;
 `,[S("dot",{boxShadow:"var(--n-box-shadow-disabled)",backgroundColor:"var(--n-color-disabled)"},[D("&::before",{backgroundColor:"var(--n-dot-color-disabled)"}),_("checked",`
 opacity: 1;
 `)]),S("label",{color:"var(--n-text-color-disabled)"}),O("radio-input",`
 cursor: not-allowed;
 `)])]),$e=Object.assign(Object.assign({},q.props),_e),Ue=X({name:"Radio",props:$e,setup(e){const t=Ee(e),n=q("Radio","-radio",Ve,he,e,t.mergedClsPrefix),a=A(()=>{const{mergedSize:{value:v}}=t,{common:{cubicBezierEaseInOut:y},self:{boxShadow:B,boxShadowActive:z,boxShadowDisabled:V,boxShadowFocus:k,boxShadowHover:T,color:F,colorDisabled:$,colorActive:u,textColor:x,textColorDisabled:w,dotColorActive:g,dotColorDisabled:C,labelPadding:M,labelLineHeight:L,labelFontWeight:N,[Y("fontSize",v)]:U,[Y("radioSize",v)]:o}}=n.value;return{"--n-bezier":y,"--n-label-line-height":L,"--n-label-font-weight":N,"--n-box-shadow":B,"--n-box-shadow-active":z,"--n-box-shadow-disabled":V,"--n-box-shadow-focus":k,"--n-box-shadow-hover":T,"--n-color":F,"--n-color-active":u,"--n-color-disabled":$,"--n-dot-color-active":g,"--n-dot-color-disabled":C,"--n-font-size":U,"--n-radio-size":o,"--n-text-color":x,"--n-text-color-disabled":w,"--n-label-padding":M}}),{inlineThemeDisabled:i,mergedClsPrefixRef:s,mergedRtlRef:l}=re(e),f=ce("Radio",l,s),c=i?fe("radio",A(()=>t.mergedSize.value[0]),a,e):void 0;return Object.assign(t,{rtlEnabled:f,cssVars:i?void 0:a,themeClass:c==null?void 0:c.themeClass,onRender:c==null?void 0:c.onRender})},render(){const{$slots:e,mergedClsPrefix:t,onRender:n,label:a}=this;return n==null||n(),I("label",{class:[`${t}-radio`,this.themeClass,{[`${t}-radio--rtl`]:this.rtlEnabled,[`${t}-radio--disabled`]:this.mergedDisabled,[`${t}-radio--checked`]:this.renderSafeChecked,[`${t}-radio--focus`]:this.focus}],style:this.cssVars},I("input",{ref:"inputRef",type:"radio",class:`${t}-radio-input`,value:this.value,name:this.mergedName,checked:this.renderSafeChecked,disabled:this.mergedDisabled,onChange:this.handleRadioInputChange,onFocus:this.handleRadioInputFocus,onBlur:this.handleRadioInputBlur}),I("div",{class:`${t}-radio__dot-wrapper`}," ",I("div",{class:[`${t}-radio__dot`,this.renderSafeChecked&&`${t}-radio__dot--checked`]})),ze(e.default,i=>!i&&!a?null:I("div",{ref:"labelRef",class:`${t}-radio__label`},i||a)))}}),He=O("radio-group",`
 display: inline-block;
 font-size: var(--n-font-size);
`,[S("splitor",`
 display: inline-block;
 vertical-align: bottom;
 width: 1px;
 transition:
 background-color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 background: var(--n-button-border-color);
 `,[_("checked",{backgroundColor:"var(--n-button-border-color-active)"}),_("disabled",{opacity:"var(--n-opacity-disabled)"})]),_("button-group",`
 white-space: nowrap;
 height: var(--n-height);
 line-height: var(--n-height);
 `,[O("radio-button",{height:"var(--n-height)",lineHeight:"var(--n-height)"}),S("splitor",{height:"var(--n-height)"})]),O("radio-button",`
 vertical-align: bottom;
 outline: none;
 position: relative;
 user-select: none;
 -webkit-user-select: none;
 display: inline-block;
 box-sizing: border-box;
 padding-left: 14px;
 padding-right: 14px;
 white-space: nowrap;
 transition:
 background-color .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 color: var(--n-button-text-color);
 border-top: 1px solid var(--n-button-border-color);
 border-bottom: 1px solid var(--n-button-border-color);
 `,[O("radio-input",`
 pointer-events: none;
 position: absolute;
 border: 0;
 border-radius: inherit;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 opacity: 0;
 z-index: 1;
 `),S("state-border",`
 z-index: 1;
 pointer-events: none;
 position: absolute;
 box-shadow: var(--n-button-box-shadow);
 transition: box-shadow .3s var(--n-bezier);
 left: -1px;
 bottom: -1px;
 right: -1px;
 top: -1px;
 `),D("&:first-child",`
 border-top-left-radius: var(--n-button-border-radius);
 border-bottom-left-radius: var(--n-button-border-radius);
 border-left: 1px solid var(--n-button-border-color);
 `,[S("state-border",`
 border-top-left-radius: var(--n-button-border-radius);
 border-bottom-left-radius: var(--n-button-border-radius);
 `)]),D("&:last-child",`
 border-top-right-radius: var(--n-button-border-radius);
 border-bottom-right-radius: var(--n-button-border-radius);
 border-right: 1px solid var(--n-button-border-color);
 `,[S("state-border",`
 border-top-right-radius: var(--n-button-border-radius);
 border-bottom-right-radius: var(--n-button-border-radius);
 `)]),ne("disabled",`
 cursor: pointer;
 `,[D("&:hover",[S("state-border",`
 transition: box-shadow .3s var(--n-bezier);
 box-shadow: var(--n-button-box-shadow-hover);
 `),ne("checked",{color:"var(--n-button-text-color-hover)"})]),_("focus",[D("&:not(:active)",[S("state-border",{boxShadow:"var(--n-button-box-shadow-focus)"})])])]),_("checked",`
 background: var(--n-button-color-active);
 color: var(--n-button-text-color-active);
 border-color: var(--n-button-border-color-active);
 `),_("disabled",`
 cursor: not-allowed;
 opacity: var(--n-opacity-disabled);
 `)])]);function Ae(e,t,n){var a;const i=[];let s=!1;for(let l=0;l<e.length;++l){const f=e[l],c=(a=f.type)===null||a===void 0?void 0:a.name;c==="RadioButton"&&(s=!0);const v=f.props;if(c!=="RadioButton"){i.push(f);continue}if(l===0)i.push(f);else{const y=i[i.length-1].props,B=t===y.value,z=y.disabled,V=t===v.value,k=v.disabled,T=(B?2:0)+(z?0:1),F=(V?2:0)+(k?0:1),$={[`${n}-radio-group__splitor--disabled`]:z,[`${n}-radio-group__splitor--checked`]:B},u={[`${n}-radio-group__splitor--disabled`]:k,[`${n}-radio-group__splitor--checked`]:V},x=T<F?u:$;i.push(I("div",{class:[`${n}-radio-group__splitor`,x]}),f)}}return{children:i,isButtonGroup:s}}const De=Object.assign(Object.assign({},q.props),{name:String,value:[String,Number,Boolean],defaultValue:{type:[String,Number,Boolean],default:null},size:String,disabled:{type:Boolean,default:void 0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array]}),Le=X({name:"RadioGroup",props:De,setup(e){const t=E(null),{mergedSizeRef:n,mergedDisabledRef:a,nTriggerFormChange:i,nTriggerFormInput:s,nTriggerFormBlur:l,nTriggerFormFocus:f}=de(e),{mergedClsPrefixRef:c,inlineThemeDisabled:v,mergedRtlRef:y}=re(e),B=q("Radio","-radio-group",He,he,e,c),z=E(e.defaultValue),V=oe(e,"value"),k=ue(V,z);function T(g){const{onUpdateValue:C,"onUpdate:value":M}=e;C&&W(C,g),M&&W(M,g),z.value=g,i(),s()}function F(g){const{value:C}=t;C&&(C.contains(g.relatedTarget)||f())}function $(g){const{value:C}=t;C&&(C.contains(g.relatedTarget)||l())}ke(be,{mergedClsPrefixRef:c,nameRef:oe(e,"name"),valueRef:k,disabledRef:a,mergedSizeRef:n,doUpdateValue:T});const u=ce("Radio",y,c),x=A(()=>{const{value:g}=n,{common:{cubicBezierEaseInOut:C},self:{buttonBorderColor:M,buttonBorderColorActive:L,buttonBorderRadius:N,buttonBoxShadow:U,buttonBoxShadowFocus:o,buttonBoxShadowHover:r,buttonColorActive:d,buttonTextColor:b,buttonTextColorActive:m,buttonTextColorHover:p,opacityDisabled:h,[Y("buttonHeight",g)]:R,[Y("fontSize",g)]:H}}=B.value;return{"--n-font-size":H,"--n-bezier":C,"--n-button-border-color":M,"--n-button-border-color-active":L,"--n-button-border-radius":N,"--n-button-box-shadow":U,"--n-button-box-shadow-focus":o,"--n-button-box-shadow-hover":r,"--n-button-color-active":d,"--n-button-text-color":b,"--n-button-text-color-hover":p,"--n-button-text-color-active":m,"--n-height":R,"--n-opacity-disabled":h}}),w=v?fe("radio-group",A(()=>n.value[0]),x,e):void 0;return{selfElRef:t,rtlEnabled:u,mergedClsPrefix:c,mergedValue:k,handleFocusout:$,handleFocusin:F,cssVars:v?void 0:x,themeClass:w==null?void 0:w.themeClass,onRender:w==null?void 0:w.onRender}},render(){var e;const{mergedValue:t,mergedClsPrefix:n,handleFocusin:a,handleFocusout:i}=this,{children:s,isButtonGroup:l}=Ae(Ce(Te(this)),t,n);return(e=this.onRender)===null||e===void 0||e.call(this),I("div",{onFocusin:a,onFocusout:i,ref:"selfElRef",class:[`${n}-radio-group`,this.rtlEnabled&&`${n}-radio-group--rtl`,this.themeClass,l&&`${n}-radio-group--button-group`],style:this.cssVars},s)}});export{Oe as F,Le as N,Pe as V,Ue as a};
