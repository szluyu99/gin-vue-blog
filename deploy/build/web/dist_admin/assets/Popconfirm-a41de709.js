import{aa as w,G as y,H as N,ab as O,M as f,O as j,K as g,a3 as h,s as a,N as P,a1 as B,R,aT as I,bM as $,F as C,am as m,Y as x,J as S,bQ as F,e as K,au as M,bR as U,bw as E,a9 as V,bS as q,S as b}from"./index-7d3bcf57.js";import{u as T}from"./Input-41b621cb.js";const _=w("n-popconfirm"),z={positiveText:String,negativeText:String,showIcon:{type:Boolean,default:!0},onPositiveClick:{type:Function,required:!0},onNegativeClick:{type:Function,required:!0}},k=$(z),H=y({name:"NPopconfirmPanel",props:z,setup(e){const{localeRef:n}=T("Popconfirm"),{inlineThemeDisabled:s}=N(),{mergedClsPrefixRef:t,mergedThemeRef:v,props:r}=O(_),u=f(()=>{const{common:{cubicBezierEaseInOut:o},self:{fontSize:l,iconSize:c,iconColor:d}}=v.value;return{"--n-bezier":o,"--n-font-size":l,"--n-icon-size":c,"--n-icon-color":d}}),i=s?j("popconfirm-panel",void 0,u,r):void 0;return Object.assign(Object.assign({},T("Popconfirm")),{mergedClsPrefix:t,cssVars:s?void 0:u,localizedPositiveText:f(()=>e.positiveText||n.value.positiveText),localizedNegativeText:f(()=>e.negativeText||n.value.negativeText),positiveButtonProps:g(r,"positiveButtonProps"),negativeButtonProps:g(r,"negativeButtonProps"),handlePositiveClick(o){e.onPositiveClick(o)},handleNegativeClick(o){e.onNegativeClick(o)},themeClass:i==null?void 0:i.themeClass,onRender:i==null?void 0:i.onRender})},render(){var e;const{mergedClsPrefix:n,showIcon:s,$slots:t}=this,v=h(t.action,()=>this.negativeText===null&&this.positiveText===null?[]:[this.negativeText!==null&&a(P,Object.assign({size:"small",onClick:this.handleNegativeClick},this.negativeButtonProps),{default:()=>this.localizedNegativeText}),this.positiveText!==null&&a(P,Object.assign({size:"small",type:"primary",onClick:this.handlePositiveClick},this.positiveButtonProps),{default:()=>this.localizedPositiveText})]);return(e=this.onRender)===null||e===void 0||e.call(this),a("div",{class:[`${n}-popconfirm__panel`,this.themeClass],style:this.cssVars},B(t.default,r=>s||r?a("div",{class:`${n}-popconfirm__body`},s?a("div",{class:`${n}-popconfirm__icon`},h(t.icon,()=>[a(R,{clsPrefix:n},{default:()=>a(I,null)})])):null,r):null),v?a("div",{class:[`${n}-popconfirm__action`]},v):null)}}),L=C("popconfirm",[m("body",`
 font-size: var(--n-font-size);
 display: flex;
 align-items: center;
 flex-wrap: nowrap;
 position: relative;
 `,[m("icon",`
 display: flex;
 font-size: var(--n-icon-size);
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 margin: 0 8px 0 0;
 `)]),m("action",`
 display: flex;
 justify-content: flex-end;
 `,[x("&:not(:first-child)","margin-top: 8px"),C("button",[x("&:not(:last-child)","margin-right: 8px;")])])]),W=Object.assign(Object.assign(Object.assign({},S.props),q),{positiveText:String,negativeText:String,showIcon:{type:Boolean,default:!0},trigger:{type:String,default:"click"},positiveButtonProps:Object,negativeButtonProps:Object,onPositiveClick:Function,onNegativeClick:Function}),J=y({name:"Popconfirm",props:W,__popover__:!0,setup(e){const{mergedClsPrefixRef:n}=N(),s=S("Popconfirm","-popconfirm",L,F,e,n),t=K(null);function v(i){var o;if(!(!((o=t.value)===null||o===void 0)&&o.getMergedShow()))return;const{onPositiveClick:l,"onUpdate:show":c}=e;Promise.resolve(l?l(i):!0).then(d=>{var p;d!==!1&&((p=t.value)===null||p===void 0||p.setShow(!1),c&&b(c,!1))})}function r(i){var o;if(!(!((o=t.value)===null||o===void 0)&&o.getMergedShow()))return;const{onNegativeClick:l,"onUpdate:show":c}=e;Promise.resolve(l?l(i):!0).then(d=>{var p;d!==!1&&((p=t.value)===null||p===void 0||p.setShow(!1),c&&b(c,!1))})}return M(_,{mergedThemeRef:s,mergedClsPrefixRef:n,props:e}),{setShow(i){var o;(o=t.value)===null||o===void 0||o.setShow(i)},syncPosition(){var i;(i=t.value)===null||i===void 0||i.syncPosition()},mergedTheme:s,popoverInstRef:t,handlePositiveClick:v,handleNegativeClick:r}},render(){const{$slots:e,$props:n,mergedTheme:s}=this;return a(V,E(n,k,{theme:s.peers.Popover,themeOverrides:s.peerOverrides.Popover,internalExtraClass:["popconfirm"],ref:"popoverInstRef"}),{trigger:e.activator||e.trigger,default:()=>{const t=U(n,k);return a(H,Object.assign(Object.assign({},t),{onPositiveClick:this.handlePositiveClick,onNegativeClick:this.handleNegativeClick}),e)}})}});export{J as N};
