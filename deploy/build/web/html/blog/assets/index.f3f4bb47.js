import{a1 as xe,d4 as we,a5 as he,aY as h,dl as fe,a4 as ne,aZ as z,a8 as ve,j as de,R as ke,a3 as be,da as Ce,r as m,E as ge,b6 as _e,k as y,ac as Se,b2 as ae,P as k,b4 as D,dm as Be,cW as $e,b9 as oe,cX as W,O as ie,cY as E,G as me,o as ye,L as Fe,n as Re,cC as ze,b as re,c as ue,g as C,q as Ee,d as Le,ba as Ve,z as De,A as We,I as pe,w as Ae,af as Ne,cP as Me,cR as He,cQ as Te,f as Ie,s as le,e as O,l as je,u as se,H as Pe,p as Ue,i as Ke,m as Xe,h as Oe,cU as Ye,_ as qe}from"./index.d0faf39f.js";const Ge={buttonHeightSmall:"14px",buttonHeightMedium:"18px",buttonHeightLarge:"22px",buttonWidthSmall:"14px",buttonWidthMedium:"18px",buttonWidthLarge:"22px",buttonWidthPressedSmall:"20px",buttonWidthPressedMedium:"24px",buttonWidthPressedLarge:"28px",railHeightSmall:"18px",railHeightMedium:"22px",railHeightLarge:"26px",railWidthSmall:"32px",railWidthMedium:"40px",railWidthLarge:"48px"},Qe=e=>{const{primaryColor:f,opacityDisabled:p,borderRadius:u,textColor3:t}=e,b="rgba(0, 0, 0, .14)";return Object.assign(Object.assign({},Ge),{iconColor:t,textColor:"white",loadingColor:f,opacityDisabled:p,railColor:b,railColorActive:f,buttonBoxShadow:"0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)",buttonColor:"#FFF",railBorderRadiusSmall:u,railBorderRadiusMedium:u,railBorderRadiusLarge:u,buttonBorderRadiusSmall:u,buttonBorderRadiusMedium:u,buttonBorderRadiusLarge:u,boxShadowFocus:`0 0 0 2px ${we(f,{alpha:.2})}`})},Ze={name:"Switch",common:xe,self:Qe},Je=Ze,et=he("switch",`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[h("children-placeholder",`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),h("rail-placeholder",`
 display: flex;
 flex-wrap: none;
 `),h("button-placeholder",`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),he("base-loading",`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[fe({left:"50%",top:"50%",originalTransform:"translateX(-50%) translateY(-50%)"})]),h("checked, unchecked",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),h("checked",`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),h("unchecked",`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),ne("&:focus",[h("rail",`
 box-shadow: var(--n-box-shadow-focus);
 `)]),z("round",[h("rail","border-radius: calc(var(--n-rail-height) / 2);",[h("button","border-radius: calc(var(--n-button-height) / 2);")])]),ve("disabled",[ve("icon",[z("rubber-band",[z("pressed",[h("rail",[h("button","max-width: var(--n-button-width-pressed);")])]),h("rail",[ne("&:active",[h("button","max-width: var(--n-button-width-pressed);")])]),z("active",[z("pressed",[h("rail",[h("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])]),h("rail",[ne("&:active",[h("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])])])])])]),z("active",[h("rail",[h("button","left: calc(100% - var(--n-button-width) - var(--n-offset))")])]),h("rail",`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[h("button-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[fe()]),h("button",`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),z("active",[h("rail","background-color: var(--n-rail-color-active);")]),z("loading",[h("rail",`
 cursor: wait;
 `)]),z("disabled",[h("rail",`
 cursor: not-allowed;
 opacity: .5;
 `)])]),tt=Object.assign(Object.assign({},be.props),{size:{type:String,default:"medium"},value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},onChange:[Function,Array]});let j;const nt=de({name:"Switch",props:tt,setup(e){j===void 0&&(typeof CSS<"u"?typeof CSS.supports<"u"?j=CSS.supports("width","max(1px)"):j=!1:j=!0);const{mergedClsPrefixRef:f,inlineThemeDisabled:p}=ke(e),u=be("Switch","-switch",et,Je,e,f),t=Ce(e),{mergedSizeRef:b,mergedDisabledRef:v}=t,_=m(e.defaultValue),B=ge(e,"value"),x=_e(B,_),g=y(()=>x.value===e.checkedValue),S=m(!1),n=m(!1),a=y(()=>{const{railStyle:d}=e;if(!!d)return d({focused:n.value,checked:g.value})});function c(d){const{"onUpdate:value":A,onChange:N,onUpdateValue:V}=e,{nTriggerFormInput:H,nTriggerFormChange:T}=t;A&&oe(A,d),V&&oe(V,d),N&&oe(N,d),_.value=d,H(),T()}function w(){const{nTriggerFormFocus:d}=t;d()}function s(){const{nTriggerFormBlur:d}=t;d()}function P(){e.loading||v.value||(x.value!==e.checkedValue?c(e.checkedValue):c(e.uncheckedValue))}function U(){n.value=!0,w()}function K(){n.value=!1,s(),S.value=!1}function Y(d){e.loading||v.value||d.key===" "&&(x.value!==e.checkedValue?c(e.checkedValue):c(e.uncheckedValue),S.value=!1)}function M(d){e.loading||v.value||d.key===" "&&(d.preventDefault(),S.value=!0)}const X=y(()=>{const{value:d}=b,{self:{opacityDisabled:A,railColor:N,railColorActive:V,buttonBoxShadow:H,buttonColor:T,boxShadowFocus:q,loadingColor:G,textColor:Q,iconColor:Z,[W("buttonHeight",d)]:$,[W("buttonWidth",d)]:J,[W("buttonWidthPressed",d)]:ee,[W("railHeight",d)]:o,[W("railWidth",d)]:r,[W("railBorderRadius",d)]:l,[W("buttonBorderRadius",d)]:i},common:{cubicBezierEaseInOut:F}}=u.value;let R,I,te;return j?(R=`calc((${o} - ${$}) / 2)`,I=`max(${o}, ${$})`,te=`max(${r}, calc(${r} + ${$} - ${o}))`):(R=ie((E(o)-E($))/2),I=ie(Math.max(E(o),E($))),te=E(o)>E($)?r:ie(E(r)+E($)-E(o))),{"--n-bezier":F,"--n-button-border-radius":i,"--n-button-box-shadow":H,"--n-button-color":T,"--n-button-width":J,"--n-button-width-pressed":ee,"--n-button-height":$,"--n-height":I,"--n-offset":R,"--n-opacity-disabled":A,"--n-rail-border-radius":l,"--n-rail-color":N,"--n-rail-color-active":V,"--n-rail-height":o,"--n-rail-width":r,"--n-width":te,"--n-box-shadow-focus":q,"--n-loading-color":G,"--n-text-color":Q,"--n-icon-color":Z}}),L=p?Se("switch",y(()=>b.value[0]),X,e):void 0;return{handleClick:P,handleBlur:K,handleFocus:U,handleKeyup:Y,handleKeydown:M,mergedRailStyle:a,pressed:S,mergedClsPrefix:f,mergedValue:x,checked:g,mergedDisabled:v,cssVars:p?void 0:X,themeClass:L==null?void 0:L.themeClass,onRender:L==null?void 0:L.onRender}},render(){const{mergedClsPrefix:e,mergedDisabled:f,checked:p,mergedRailStyle:u,onRender:t,$slots:b}=this;t==null||t();const{checked:v,unchecked:_,icon:B,"checked-icon":x,"unchecked-icon":g}=b,S=!(ae(B)&&ae(x)&&ae(g));return k("div",{role:"switch","aria-checked":p,class:[`${e}-switch`,this.themeClass,S&&`${e}-switch--icon`,p&&`${e}-switch--active`,f&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},k("div",{class:`${e}-switch__rail`,"aria-hidden":"true",style:u},D(v,n=>D(_,a=>n||a?k("div",{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},k("div",{class:`${e}-switch__rail-placeholder`},k("div",{class:`${e}-switch__button-placeholder`}),n),k("div",{class:`${e}-switch__rail-placeholder`},k("div",{class:`${e}-switch__button-placeholder`}),a)):null)),k("div",{class:`${e}-switch__button`},D(B,n=>D(x,a=>D(g,c=>k(Be,null,{default:()=>this.loading?k($e,{key:"loading",clsPrefix:e,strokeWidth:20}):this.checked&&(a||n)?k("div",{class:`${e}-switch__button-icon`,key:a?"checked-icon":"icon"},a||n):!this.checked&&(c||n)?k("div",{class:`${e}-switch__button-icon`,key:c?"unchecked-icon":"icon"},c||n):null})))),D(v,n=>n&&k("div",{key:"checked",class:`${e}-switch__checked`},n)),D(_,n=>n&&k("div",{key:"unchecked",class:`${e}-switch__unchecked`},n)))))}});function at(e,f,p="modelValue",u){return y({get:()=>e[p],set:t=>{f(`update:${p}`,u?u(t):t)}})}var ce=de({components:{},props:{danmus:{type:Array,required:!0,default:()=>[]},channels:{type:Number,default:0},autoplay:{type:Boolean,default:!0},loop:{type:Boolean,default:!1},useSlot:{type:Boolean,default:!1},debounce:{type:Number,default:100},speeds:{type:Number,default:200},randomChannel:{type:Boolean,default:!1},fontSize:{type:Number,default:18},top:{type:Number,default:4},right:{type:Number,default:0},isSuspend:{type:Boolean,default:!1},extraStyle:{type:String,default:""}},emits:["list-end","play-end","update:danmus"],setup(e,{emit:f,slots:p}){let u=m(document.createElement("div")),t=m(document.createElement("div"));const b=m(0),v=m(0);let _=0;const B=m(0),x=m(0),g=m(0),S=m(!1),n=m(!1),a=m({}),c=at(e,f,"danmus"),w=me({channels:y(()=>e.channels||B.value),autoplay:y(()=>e.autoplay),loop:y(()=>e.loop),useSlot:y(()=>e.useSlot),debounce:y(()=>e.debounce),randomChannel:y(()=>e.randomChannel)}),s=me({height:y(()=>x.value),fontSize:y(()=>e.fontSize),speeds:y(()=>e.speeds),top:y(()=>e.top),right:y(()=>e.right)});ye(()=>{P()}),Fe(()=>{V()});function P(){U(),e.isSuspend&&N(),w.autoplay&&K()}function U(){b.value=u.value.offsetWidth,v.value=u.value.offsetHeight}function K(){n.value=!1,_||(_=setInterval(()=>Y(),w.debounce))}function Y(){if(!n.value&&c.value.length)if(g.value>c.value.length-1){const o=t.value.children.length;w.loop&&o<g.value&&(f("list-end"),g.value=0,M())}else M()}function M(o){const r=w.loop?g.value%c.value.length:g.value,l=o||c.value[r];let i=document.createElement("div");w.useSlot?i=X(l,r).$el:(i.innerHTML=l,i.setAttribute("style",e.extraStyle),i.style.fontSize=`${s.fontSize}px`,i.style.lineHeight=`${s.fontSize}px`),i.classList.add("dm"),t.value.appendChild(i),i.style.opacity="0",Re(()=>{s.height||(x.value=i.offsetHeight),w.channels||(B.value=Math.floor(v.value/(s.height+s.top)));let F=L(i);if(F>=0){const R=i.offsetWidth,I=s.height;i.classList.add("move"),i.dataset.index=`${r}`,i.style.opacity="1",i.style.top=F*(I+s.top)+"px",i.style.width=R+s.right+"px",i.style.setProperty("--dm-scroll-width",`-${b.value+R}px`),i.style.left=`${b.value}px`,i.style.animationDuration=`${b.value/s.speeds}s`,i.addEventListener("animationend",()=>{Number(i.dataset.index)===c.value.length-1&&!w.loop&&f("play-end",i.dataset.index),t.value&&t.value.removeChild(i)}),g.value++}else t.value.removeChild(i)})}function X(o,r){return ze({render(){return k("div",{},[p.dm&&p.dm({danmu:o,index:r})])}}).mount(document.createElement("div"))}function L(o){let r=[...Array(w.channels).keys()];w.randomChannel&&(r=r.sort(()=>.5-Math.random()));for(let l of r){const i=a.value[l];if(i&&i.length)for(let F=0;F<i.length;F++){const R=d(i[F])-10;if(R<=(o.offsetWidth-i[F].offsetWidth)*.88||R<=0)break;if(F===i.length-1)return a.value[l].push(o),o.addEventListener("animationend",()=>a.value[l].splice(0,1)),l%w.channels}else return a.value[l]=[o],o.addEventListener("animationend",()=>a.value[l].splice(0,1)),l%w.channels}return-1}function d(o){const r=o.offsetWidth||parseInt(o.style.width),l=o.getBoundingClientRect().right||t.value.getBoundingClientRect().right+r;return t.value.getBoundingClientRect().right-l}function A(){clearInterval(_),_=0}function N(){let o=[];t.value.addEventListener("mousemove",r=>{let l=r.target;l.className.includes("dm")||(l=l.closest(".dm")||l),l.className.includes("dm")&&(l.classList.add("pause"),o.push(l))}),t.value.addEventListener("mouseout",r=>{let l=r.target;l.className.includes("dm")||(l=l.closest(".dm")||l),l.className.includes("dm")&&(l.classList.remove("pause"),o.forEach(i=>{i.classList.remove("pause")}),o=[])})}function V(){A(),g.value=0}function H(){x.value=0,P()}function T(){a.value={},t.value.innerHTML="",n.value=!0,S.value=!1,V()}function q(){n.value=!0}function G(o){if(g.value===c.value.length)return c.value.push(o),c.value.length-1;{const r=g.value%c.value.length;return c.value.splice(r,0,o),r+1}}function Q(o){return c.value.push(o),c.value.length-1}function Z(){return!n.value}function $(){S.value=!1}function J(){S.value=!0}function ee(){U();const o=t.value.getElementsByClassName("dm");for(let r=0;r<o.length;r++){const l=o[r];l.style.setProperty("--dm-scroll-width",`-${b.value+l.offsetWidth}px`),l.style.left=`${b.value}px`,l.style.animationDuration=`${b.value/s.speeds}s`}}return{container:u,dmContainer:t,hidden:S,paused:n,danmuList:c,getPlayState:Z,resize:ee,play:K,pause:q,stop:T,show:$,hide:J,reset:H,add:G,push:Q,insert:M}}});const ot={ref:"container",class:"vue-danmaku"};function it(e,f,p,u,t,b){return re(),ue("div",ot,[C("div",{ref:"dmContainer",class:Ee(["danmus",{show:!e.hidden},{paused:e.paused}])},null,2),Le(e.$slots,"default")],512)}function lt(e,f){f===void 0&&(f={});var p=f.insertAt;if(!(!e||typeof document>"u")){var u=document.head||document.getElementsByTagName("head")[0],t=document.createElement("style");t.type="text/css",p==="top"&&u.firstChild?u.insertBefore(t,u.firstChild):u.appendChild(t),t.styleSheet?t.styleSheet.cssText=e:t.appendChild(document.createTextNode(e))}}var st=`.vue-danmaku {
  position: relative;
  overflow: hidden;
}
.vue-danmaku .danmus {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  -webkit-transition: all 0.3s;
  transition: all 0.3s;
}
.vue-danmaku .danmus.show {
  opacity: 1;
}
.vue-danmaku .danmus.paused .dm.move {
  animation-play-state: paused;
}
.vue-danmaku .danmus .dm {
  position: absolute;
  font-size: 20px;
  color: #ddd;
  white-space: pre;
  transform: translateX(0);
  transform-style: preserve-3d;
}
.vue-danmaku .danmus .dm.move {
  will-change: transform;
  animation-name: moveLeft;
  animation-timing-function: linear;
  animation-play-state: running;
}
.vue-danmaku .danmus .dm.pause {
  animation-play-state: paused;
  z-index: 10;
}
@keyframes moveLeft {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(var(--dm-scroll-width));
  }
}
@-webkit-keyframes moveLeft {
  from {
    -webkit-transform: translateX(0);
  }
  to {
    -webkit-transform: translateX(var(--dm-scroll-width));
  }
}`;lt(st);ce.render=it;ce.__file="src/lib/Danmaku.vue";const rt=e=>(Ue("data-v-0a16ac23"),e=e(),Ke(),e),ut={"animate-zoom-in":"",top:"3/10","border-1":"","rounded-10":"","shadow-2xl":"",absolute:"","inset-x-1":"","text-center":"","text-light":"","px-5":"","py-15":"","mx-auto":"","z-5":"","w-350":"",lg:"w-420 py-20"},dt=rt(()=>C("h1",{"font-bold":"","text-25":"","lg:text-30":""}," \u7559\u8A00\u677F ",-1)),ct={flex:"","justify-center":"","mt-20":"","h-36":"",lg:"mt-25 h-40"},ht=["onKeyup"],ft={"text-left":"","text-white":"","ml-20":""},vt={"mt-25":"",flex:"","items-center":""},mt={"my-20":""},pt={"mb-15":""},bt={absolute:"","inset-0":"","top-60":""},gt={bg:"#00000060","rounded-20":"","text-white":"",flex:"","items-center":"","text-15":"","px-6":"","py-4":"",lg:"text-16 px-8 py-5"},yt=de({__name:"index",setup(e){const f=Ve(),p=De(We()),u=ge(p,"pageList");let t=m("");const b=m(!1),v=m(null),_=m(!1),B=m(!1);let x=m([{avatar:"https://static.talkxj.com/avatar/15bb3e8d4144163e006f4b2776aa7bac.jpg",content:"\u5927\u5BB6\u597D\uFF0C\u6211\u662F\u4F5C\u8005\uFF0C\u6B22\u8FCE\u7ED9\u6211\u70B9\u4E00\u9897 Star!",nickname:"\u9635\u3001\u96E8"}]);ye(async()=>{const n=await pe.getMessages();x.value=[...x.value,...n.data]});async function g(){var a;if(t.value=t.value.trim(),!t.value){(a=window==null?void 0:window.$message)==null||a.info("\u6D88\u606F\u4E0D\u80FD\u4E3A\u7A7A!");return}const n={avatar:f.avatar,nickname:f.nickname,content:t.value};await pe.saveMessage(n),v.value.push(n),t.value=""}Ae(_,n=>n?v.value.hide():v.value.show());const S=y(()=>{const n=u.value.find(a=>a.label==="message");return n?`background: url('${n==null?void 0:n.cover}') center center / cover no-repeat;`:'background: url("https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg") center center / cover no-repeat;'});return(n,a)=>{const c=nt,w=Ye;return re(),ue("div",{style:Pe(se(S)),"overflow-hidden":"",absolute:"","inset-x-0":"","h-screen":"",class:"banner-fade-down"},[C("div",ut,[dt,C("div",ct,[Ne(C("input",{"onUpdate:modelValue":a[0]||(a[0]=s=>t.value=s),"border-1":"","rounded-20":"","bg-transparent":"","px-20":"",w:"3/4",text:"#eee",placeholder:"\u8BF4\u70B9\u4EC0\u4E48\u5427\uFF1F",onClick:a[1]||(a[1]=He(s=>b.value=!0,["stop"])),onKeyup:Te(g,["enter"])},null,40,ht),[[Me,t.value]]),b.value?(re(),ue("button",{key:0,"ml-12":"","rounded-20":"","px-18":"","border-1":"","animate-back-in-right":"",onClick:g}," \u53D1\u9001 ")):Ie("",!0)]),C("ul",ft,[C("li",vt,[le(" \u5FAA\u73AF\u64AD\u653E\uFF1A "),O(c,{value:B.value,"onUpdate:value":a[2]||(a[2]=s=>B.value=s),"ml-5":""},null,8,["value"])]),C("li",mt,[le(" \u64CD\u4F5C\u5F39\u5E55\uFF1A "),C("button",{"ml-5":"","border-1":"","rounded-8":"","p-3":"",onClick:a[3]||(a[3]=(...s)=>v.value.play&&v.value.play(...s))}," \u64AD\u653E "),C("button",{"mx-15":"","border-1":"","rounded-8":"","p-3":"",onClick:a[4]||(a[4]=(...s)=>v.value.pause&&v.value.pause(...s))}," \u6682\u505C "),C("button",{"border-1":"","rounded-8":"","p-3":"",onClick:a[5]||(a[5]=(...s)=>v.value.stop&&v.value.stop(...s))}," \u505C\u6B62 ")]),C("li",pt,[le(" \u9690\u85CF\u5F39\u5E55\uFF1A "),O(c,{value:_.value,"onUpdate:value":a[6]||(a[6]=s=>_.value=s),"ml-5":""},null,8,["value"])])])]),C("div",bt,[O(se(ce),{ref_key:"dmRef",ref:v,danmus:x.value,"onUpdate:danmus":a[7]||(a[7]=s=>x.value=s),"wh-full":"","use-slot":"",loop:B.value,speeds:200,channels:0,top:5,"is-suspend":!0},{dm:je(({danmu:s})=>[C("div",gt,[O(w,{round:"",size:"small",src:se(Xe)(s.avatar),"mr-10":""},null,8,["src"]),C("span",null,Oe(`${s.nickname} : ${s.content}`),1)])]),_:1},8,["danmus","loop"])])],4)}}});const wt=qe(yt,[["__scopeId","data-v-0a16ac23"]]);export{wt as default};
