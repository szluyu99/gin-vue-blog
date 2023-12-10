import{s,G as _,H as j,I as E,e as N,M as I,L as G,au as se,K as P,aa as be,S as i,Y as h,F as r,ak as S,am as C,ao as ue,aB as he,aC as fe,ab as ke,Z as ve,J as H,$ as me,aE as K,O as xe,bG as ge,a1 as pe,ac as Ce,a2 as ye,c4 as we}from"./index-7d3bcf57.js";const Re=s("svg",{viewBox:"0 0 64 64",class:"check-icon"},s("path",{d:"M50.42,16.76L22.34,39.45l-8.1-11.46c-1.12-1.58-3.3-1.96-4.88-0.84c-1.58,1.12-1.95,3.3-0.84,4.88l10.26,14.51  c0.56,0.79,1.42,1.31,2.38,1.45c0.16,0.02,0.32,0.03,0.48,0.03c0.8,0,1.57-0.27,2.2-0.78l30.99-25.03c1.5-1.21,1.74-3.42,0.52-4.92  C54.13,15.78,51.93,15.55,50.42,16.76z"})),ze=s("svg",{viewBox:"0 0 100 100",class:"line-icon"},s("path",{d:"M80.2,55.5H21.4c-2.8,0-5.1-2.5-5.1-5.5l0,0c0-3,2.3-5.5,5.1-5.5h58.7c2.8,0,5.1,2.5,5.1,5.5l0,0C85.2,53.1,82.9,55.5,80.2,55.5z"})),L=be("n-checkbox-group"),Se={min:Number,max:Number,size:String,value:Array,defaultValue:{type:Array,default:null},disabled:{type:Boolean,default:void 0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onChange:[Function,Array]},Me=_({name:"CheckboxGroup",props:Se,setup(o){const{mergedClsPrefixRef:y}=j(o),g=E(o),{mergedSizeRef:w,mergedDisabledRef:T}=g,b=N(o.defaultValue),R=I(()=>o.value),u=G(R,b),c=I(()=>{var f;return((f=u.value)===null||f===void 0?void 0:f.length)||0}),a=I(()=>Array.isArray(u.value)?new Set(u.value):new Set);function A(f,n){const{nTriggerFormInput:p,nTriggerFormChange:k}=g,{onChange:l,"onUpdate:value":v,onUpdateValue:m}=o;if(Array.isArray(u.value)){const t=Array.from(u.value),B=t.findIndex(F=>F===n);f?~B||(t.push(n),m&&i(m,t,{actionType:"check",value:n}),v&&i(v,t,{actionType:"check",value:n}),p(),k(),b.value=t,l&&i(l,t)):~B&&(t.splice(B,1),m&&i(m,t,{actionType:"uncheck",value:n}),v&&i(v,t,{actionType:"uncheck",value:n}),l&&i(l,t),b.value=t,p(),k())}else f?(m&&i(m,[n],{actionType:"check",value:n}),v&&i(v,[n],{actionType:"check",value:n}),l&&i(l,[n]),b.value=[n],p(),k()):(m&&i(m,[],{actionType:"uncheck",value:n}),v&&i(v,[],{actionType:"uncheck",value:n}),l&&i(l,[]),b.value=[],p(),k())}return se(L,{checkedCountRef:c,maxRef:P(o,"max"),minRef:P(o,"min"),valueSetRef:a,disabledRef:T,mergedSizeRef:w,toggleCheckbox:A}),{mergedClsPrefix:y}},render(){return s("div",{class:`${this.mergedClsPrefix}-checkbox-group`,role:"group"},this.$slots)}}),Te=h([r("checkbox",`
 font-size: var(--n-font-size);
 outline: none;
 cursor: pointer;
 display: inline-flex;
 flex-wrap: nowrap;
 align-items: flex-start;
 word-break: break-word;
 line-height: var(--n-size);
 --n-merged-color-table: var(--n-color-table);
 `,[S("show-label","line-height: var(--n-label-line-height);"),h("&:hover",[r("checkbox-box",[C("border","border: var(--n-border-checked);")])]),h("&:focus:not(:active)",[r("checkbox-box",[C("border",`
 border: var(--n-border-focus);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),S("inside-table",[r("checkbox-box",`
 background-color: var(--n-merged-color-table);
 `)]),S("checked",[r("checkbox-box",`
 background-color: var(--n-color-checked);
 `,[r("checkbox-icon",[h(".check-icon",`
 opacity: 1;
 transform: scale(1);
 `)])])]),S("indeterminate",[r("checkbox-box",[r("checkbox-icon",[h(".check-icon",`
 opacity: 0;
 transform: scale(.5);
 `),h(".line-icon",`
 opacity: 1;
 transform: scale(1);
 `)])])]),S("checked, indeterminate",[h("&:focus:not(:active)",[r("checkbox-box",[C("border",`
 border: var(--n-border-checked);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),r("checkbox-box",`
 background-color: var(--n-color-checked);
 border-left: 0;
 border-top: 0;
 `,[C("border",{border:"var(--n-border-checked)"})])]),S("disabled",{cursor:"not-allowed"},[S("checked",[r("checkbox-box",`
 background-color: var(--n-color-disabled-checked);
 `,[C("border",{border:"var(--n-border-disabled-checked)"}),r("checkbox-icon",[h(".check-icon, .line-icon",{fill:"var(--n-check-mark-color-disabled-checked)"})])])]),r("checkbox-box",`
 background-color: var(--n-color-disabled);
 `,[C("border",`
 border: var(--n-border-disabled);
 `),r("checkbox-icon",[h(".check-icon, .line-icon",`
 fill: var(--n-check-mark-color-disabled);
 `)])]),C("label",`
 color: var(--n-text-color-disabled);
 `)]),r("checkbox-box-wrapper",`
 position: relative;
 width: var(--n-size);
 flex-shrink: 0;
 flex-grow: 0;
 user-select: none;
 -webkit-user-select: none;
 `),r("checkbox-box",`
 position: absolute;
 left: 0;
 top: 50%;
 transform: translateY(-50%);
 height: var(--n-size);
 width: var(--n-size);
 display: inline-block;
 box-sizing: border-box;
 border-radius: var(--n-border-radius);
 background-color: var(--n-color);
 transition: background-color 0.3s var(--n-bezier);
 `,[C("border",`
 transition:
 border-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border: var(--n-border);
 `),r("checkbox-icon",`
 display: flex;
 align-items: center;
 justify-content: center;
 position: absolute;
 left: 1px;
 right: 1px;
 top: 1px;
 bottom: 1px;
 `,[h(".check-icon, .line-icon",`
 width: 100%;
 fill: var(--n-check-mark-color);
 opacity: 0;
 transform: scale(0.5);
 transform-origin: center;
 transition:
 fill 0.3s var(--n-bezier),
 transform 0.3s var(--n-bezier),
 opacity 0.3s var(--n-bezier),
 border-color 0.3s var(--n-bezier);
 `),ue({left:"1px",top:"1px"})])]),C("label",`
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 padding: var(--n-label-padding);
 font-weight: var(--n-label-font-weight);
 `,[h("&:empty",{display:"none"})])]),he(r("checkbox",`
 --n-merged-color-table: var(--n-color-table-modal);
 `)),fe(r("checkbox",`
 --n-merged-color-table: var(--n-color-table-popover);
 `))]),De=Object.assign(Object.assign({},H.props),{size:String,checked:{type:[Boolean,String,Number],default:void 0},defaultChecked:{type:[Boolean,String,Number],default:!1},value:[String,Number],disabled:{type:Boolean,default:void 0},indeterminate:Boolean,label:String,focusable:{type:Boolean,default:!0},checkedValue:{type:[Boolean,String,Number],default:!0},uncheckedValue:{type:[Boolean,String,Number],default:!1},"onUpdate:checked":[Function,Array],onUpdateChecked:[Function,Array],privateInsideTable:Boolean,onChange:[Function,Array]}),Ae=_({name:"Checkbox",props:De,setup(o){const y=N(null),{mergedClsPrefixRef:g,inlineThemeDisabled:w,mergedRtlRef:T}=j(o),b=E(o,{mergedSize(e){const{size:x}=o;if(x!==void 0)return x;if(c){const{value:d}=c.mergedSizeRef;if(d!==void 0)return d}if(e){const{mergedSize:d}=e;if(d!==void 0)return d.value}return"medium"},mergedDisabled(e){const{disabled:x}=o;if(x!==void 0)return x;if(c){if(c.disabledRef.value)return!0;const{maxRef:{value:d},checkedCountRef:z}=c;if(d!==void 0&&z.value>=d&&!n.value)return!0;const{minRef:{value:$}}=c;if($!==void 0&&z.value<=$&&n.value)return!0}return e?e.disabled.value:!1}}),{mergedDisabledRef:R,mergedSizeRef:u}=b,c=ke(L,null),a=N(o.defaultChecked),A=P(o,"checked"),f=G(A,a),n=ve(()=>{if(c){const e=c.valueSetRef.value;return e&&o.value!==void 0?e.has(o.value):!1}else return f.value===o.checkedValue}),p=H("Checkbox","-checkbox",Te,we,o,g);function k(e){if(c&&o.value!==void 0)c.toggleCheckbox(!n.value,o.value);else{const{onChange:x,"onUpdate:checked":d,onUpdateChecked:z}=o,{nTriggerFormInput:$,nTriggerFormChange:U}=b,M=n.value?o.uncheckedValue:o.checkedValue;d&&i(d,M,e),z&&i(z,M,e),x&&i(x,M,e),$(),U(),a.value=M}}function l(e){R.value||k(e)}function v(e){if(!R.value)switch(e.key){case" ":case"Enter":k(e)}}function m(e){switch(e.key){case" ":e.preventDefault()}}const t={focus:()=>{var e;(e=y.value)===null||e===void 0||e.focus()},blur:()=>{var e;(e=y.value)===null||e===void 0||e.blur()}},B=me("Checkbox",T,g),F=I(()=>{const{value:e}=u,{common:{cubicBezierEaseInOut:x},self:{borderRadius:d,color:z,colorChecked:$,colorDisabled:U,colorTableHeader:M,colorTableHeaderModal:O,colorTableHeaderPopover:V,checkMarkColor:W,checkMarkColorDisabled:Y,border:J,borderFocus:Z,borderDisabled:q,borderChecked:Q,boxShadowFocus:X,textColor:ee,textColorDisabled:oe,checkMarkColorDisabledChecked:ne,colorDisabledChecked:re,borderDisabledChecked:ae,labelPadding:ce,labelLineHeight:le,labelFontWeight:ie,[K("fontSize",e)]:te,[K("size",e)]:de}}=p.value;return{"--n-label-line-height":le,"--n-label-font-weight":ie,"--n-size":de,"--n-bezier":x,"--n-border-radius":d,"--n-border":J,"--n-border-checked":Q,"--n-border-focus":Z,"--n-border-disabled":q,"--n-border-disabled-checked":ae,"--n-box-shadow-focus":X,"--n-color":z,"--n-color-checked":$,"--n-color-table":M,"--n-color-table-modal":O,"--n-color-table-popover":V,"--n-color-disabled":U,"--n-color-disabled-checked":re,"--n-text-color":ee,"--n-text-color-disabled":oe,"--n-check-mark-color":W,"--n-check-mark-color-disabled":Y,"--n-check-mark-color-disabled-checked":ne,"--n-font-size":te,"--n-label-padding":ce}}),D=w?xe("checkbox",I(()=>u.value[0]),F,o):void 0;return Object.assign(b,t,{rtlEnabled:B,selfRef:y,mergedClsPrefix:g,mergedDisabled:R,renderedChecked:n,mergedTheme:p,labelId:ge(),handleClick:l,handleKeyUp:v,handleKeyDown:m,cssVars:w?void 0:F,themeClass:D==null?void 0:D.themeClass,onRender:D==null?void 0:D.onRender})},render(){var o;const{$slots:y,renderedChecked:g,mergedDisabled:w,indeterminate:T,privateInsideTable:b,cssVars:R,labelId:u,label:c,mergedClsPrefix:a,focusable:A,handleKeyUp:f,handleKeyDown:n,handleClick:p}=this;(o=this.onRender)===null||o===void 0||o.call(this);const k=pe(y.default,l=>c||l?s("span",{class:`${a}-checkbox__label`,id:u},c||l):null);return s("div",{ref:"selfRef",class:[`${a}-checkbox`,this.themeClass,this.rtlEnabled&&`${a}-checkbox--rtl`,g&&`${a}-checkbox--checked`,w&&`${a}-checkbox--disabled`,T&&`${a}-checkbox--indeterminate`,b&&`${a}-checkbox--inside-table`,k&&`${a}-checkbox--show-label`],tabindex:w||!A?void 0:0,role:"checkbox","aria-checked":T?"mixed":g,"aria-labelledby":u,style:R,onKeyup:f,onKeydown:n,onClick:p,onMousedown:()=>{ye("selectstart",window,l=>{l.preventDefault()},{once:!0})}},s("div",{class:`${a}-checkbox-box-wrapper`}," ",s("div",{class:`${a}-checkbox-box`},s(Ce,null,{default:()=>this.indeterminate?s("div",{key:"indeterminate",class:`${a}-checkbox-icon`},ze):s("div",{key:"check",class:`${a}-checkbox-icon`},Re)}),s("div",{class:`${a}-checkbox-box__border`}))),k)}});export{Ae as N,Me as a};
