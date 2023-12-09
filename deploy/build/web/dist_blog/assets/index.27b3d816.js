import{d as K,r as p,l as f,E as R,o as q,S as Q,c as V,f as d,x as F,g as Y,U as J,V as O,n as Z,b as H,_ as ee,I as te,B as ne,C as ae,H as X,w as se,e as z,N as le,W as oe,X as ue,j as ie,y as M,Y as j,h as de,q as re,u as P,G as ce,s as pe,i as me,p as ve,k as fe}from"./index.51c47497.js";var I=K({name:"vue3-danmaku",components:{},props:{danmus:{type:Array,required:!0,default:()=>[]},channels:{type:Number,default:0},autoplay:{type:Boolean,default:!0},loop:{type:Boolean,default:!1},useSlot:{type:Boolean,default:!1},debounce:{type:Number,default:100},speeds:{type:Number,default:200},randomChannel:{type:Boolean,default:!1},fontSize:{type:Number,default:18},top:{type:Number,default:4},right:{type:Number,default:0},isSuspend:{type:Boolean,default:!1},extraStyle:{type:String,default:""}},emits:["list-end","play-end","dm-over","dm-out","update:danmus"],setup(a,{emit:m,slots:w}){let r=p(document.createElement("div")),l=p(document.createElement("div"));const u=p(0),k=p(0);let b=0;const _=p(0),N=p(0),h=p(0),c=p(!1),n=p(!1),s=p({}),v=function(o,i,e="modelValue",t){return f({get:()=>o[e],set:S=>{i(`update:${e}`,t?t(S):S)}})}(a,m,"danmus"),y=R({channels:f(()=>a.channels||_.value),autoplay:f(()=>a.autoplay),loop:f(()=>a.loop),useSlot:f(()=>a.useSlot),debounce:f(()=>a.debounce),randomChannel:f(()=>a.randomChannel)}),g=R({height:f(()=>N.value),fontSize:f(()=>a.fontSize),speeds:f(()=>a.speeds),top:f(()=>a.top),right:f(()=>a.right)});function T(){U(),a.isSuspend&&function(){let o=[];l.value.addEventListener("mouseover",i=>{let e=i.target;e.className.includes("dm")||(e=e.closest(".dm")||e),e.className.includes("dm")&&(o.includes(e)||(m("dm-over",{el:e}),e.classList.add("pause"),o.push(e)))}),l.value.addEventListener("mouseout",i=>{let e=i.target;e.className.includes("dm")||(e=e.closest(".dm")||e),e.className.includes("dm")&&(m("dm-out",{el:e}),e.classList.remove("pause"),o.forEach(t=>{t.classList.remove("pause")}),o=[])})}(),y.autoplay&&A()}function U(){if(u.value=r.value.offsetWidth,k.value=r.value.offsetHeight,u.value===0||k.value===0)throw new Error("\u83B7\u53D6\u4E0D\u5230\u5BB9\u5668\u5BBD\u9AD8")}function A(){n.value=!1,b||(b=setInterval(()=>function(){if(!n.value&&v.value.length)if(h.value>v.value.length-1){const o=l.value.children.length;y.loop&&(o<h.value&&(m("list-end"),h.value=0),$())}else $()}(),y.debounce))}function $(o){const i=y.loop?h.value%v.value.length:h.value,e=o||v.value[i];let t=document.createElement("div");y.useSlot?t=function(S,x){return J({render:()=>O("div",{},[w.dm&&w.dm({danmu:S,index:x})])}).mount(document.createElement("div"))}(e,i).$el:(t.innerHTML=e,t.setAttribute("style",a.extraStyle),t.style.fontSize=`${g.fontSize}px`,t.style.lineHeight=`${g.fontSize}px`),t.classList.add("dm"),l.value.appendChild(t),t.style.opacity="0",Z(()=>{g.height||(N.value=t.offsetHeight),y.channels||(_.value=Math.floor(k.value/(g.height+g.top)));let S=function(x){let L=[...Array(y.channels).keys()];y.randomChannel&&(L=L.sort(()=>.5-Math.random()));for(let C of L){const B=s.value[C];if(!B||!B.length)return s.value[C]=[x],x.addEventListener("animationend",()=>s.value[C].splice(0,1)),C%y.channels;for(let E=0;E<B.length;E++){const W=G(B[E])-10;if(W<=.88*(x.offsetWidth-B[E].offsetWidth)||W<=0)break;if(E===B.length-1)return s.value[C].push(x),x.addEventListener("animationend",()=>s.value[C].splice(0,1)),C%y.channels}}return-1}(t);if(S>=0){const x=t.offsetWidth,L=g.height;t.classList.add("move"),t.dataset.index=`${i}`,t.style.opacity="1",t.style.top=S*(L+g.top)+"px",t.style.width=x+g.right+"px",t.style.setProperty("--dm-scroll-width",`-${u.value+x}px`),t.style.left=`${u.value}px`,t.style.animationDuration=u.value/g.speeds+"s",t.addEventListener("animationend",()=>{Number(t.dataset.index)!==v.value.length-1||y.loop||m("play-end",t.dataset.index),l.value&&l.value.removeChild(t)}),h.value++}else l.value.removeChild(t)})}function G(o){const i=o.offsetWidth||parseInt(o.style.width),e=o.getBoundingClientRect().right||l.value.getBoundingClientRect().right+i;return l.value.getBoundingClientRect().right-e}function D(){clearInterval(b),b=0,h.value=0}return q(()=>{T()}),Q(()=>{D()}),{container:r,dmContainer:l,hidden:c,paused:n,danmuList:v,getPlayState:function(){return!n.value},resize:function(){U();const o=l.value.getElementsByClassName("dm");for(let i=0;i<o.length;i++){const e=o[i];e.style.setProperty("--dm-scroll-width",`-${u.value+e.offsetWidth}px`),e.style.left=`${u.value}px`,e.style.animationDuration=u.value/g.speeds+"s"}},play:A,pause:function(){n.value=!0},stop:function(){s.value={},l.value.innerHTML="",n.value=!0,c.value=!1,D()},show:function(){c.value=!1},hide:function(){c.value=!0},reset:function(){N.value=0,T()},add:function(o){if(h.value===v.value.length)return v.value.push(o),v.value.length-1;{const i=h.value%v.value.length;return v.value.splice(i,0,o),i+1}},push:function(o){return v.value.push(o),v.value.length-1},insert:$}}});const he={ref:"container",class:"vue-danmaku"};(function(a,m){m===void 0&&(m={});var w=m.insertAt;if(a&&typeof document<"u"){var r=document.head||document.getElementsByTagName("head")[0],l=document.createElement("style");l.type="text/css",w==="top"&&r.firstChild?r.insertBefore(l,r.firstChild):r.appendChild(l),l.styleSheet?l.styleSheet.cssText=a:l.appendChild(document.createTextNode(a))}})(`.vue-danmaku {
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
  z-index: 100;
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
}`),I.render=function(a,m,w,r,l,u){return H(),V("div",he,[d("div",{ref:"dmContainer",class:F(["danmus",{show:!a.hidden},{paused:a.paused}])},null,2),Y(a.$slots,"default")],512)},I.__file="src/lib/Danmaku.vue";const ye=a=>(ve("data-v-9a91cdb8"),a=a(),fe(),a),ge={class:"absolute inset-x-1 top-3/10 z-5 mx-auto w-[350px] animate-zoom-in border-1 rounded-3xl px-1 py-5 text-center text-light shadow-2xl lg:w-[420px]"},xe=ye(()=>d("h1",{class:"text-2xl font-bold"}," \u7559\u8A00\u677F ",-1)),be={class:"mt-6 h-9 flex justify-center lg:mt-6"},we={class:"ml-5 text-left text-white space-y-3"},ke={class:"mt-6 flex items-center"},_e={class:"space-x-3"},Se={class:"flex items-center"},Ce={class:"absolute inset-0 top-[60px]"},Ne={class:"flex items-center rounded-3xl bg-#00000060 px-2 py-1 text-white lg:px-4 lg:py-2"},Le=["src"],Be={class:"ml-2 text-sm"},Ee={__name:"index",setup(a){const m=te(),{pageList:w}=ne(ae()),r=p(""),l=p(!1),u=p(null),k=p(!1),b=p(!1),_=p([{avatar:"https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",content:"\u5927\u5BB6\u597D\uFF0C\u6211\u662F\u4F5C\u8005\uFF0C\u6B22\u8FCE\u7ED9\u6211\u70B9\u4E00\u9897 Star!",nickname:"\u9635\u3001\u96E8"}]);q(async()=>{const c=await X.getMessages();_.value=[..._.value,...c.data]});async function N(){var n;if(r.value=r.value.trim(),!r.value){(n=window==null?void 0:window.$message)==null||n.info("\u6D88\u606F\u4E0D\u80FD\u4E3A\u7A7A!");return}const c={avatar:m.avatar,nickname:m.nickname,content:r.value};await X.saveMessage(c),u.value.push(c),r.value=""}se(k,c=>c?u.value.hide():u.value.show());const h=f(()=>{const c=w.value.find(n=>n.label==="message");return c?`background: url('${c==null?void 0:c.cover}') center center / cover no-repeat;`:'background: url("https://static.talkxj.com/config/83be0017d7f1a29441e33083e7706936.jpg") center center / cover no-repeat;'});return(c,n)=>(H(),V("div",{style:ce(h.value),class:"banner-fade-down absolute inset-x-0 h-screen overflow-hidden"},[d("div",ge,[xe,d("div",be,[z(d("input",{"onUpdate:modelValue":n[0]||(n[0]=s=>r.value=s),class:"w-3/4 border-1 rounded-2xl bg-transparent px-4 text-sm text-#eee outline-none",placeholder:"\u8BF4\u70B9\u4EC0\u4E48\u5427\uFF1F",onClick:n[1]||(n[1]=oe(s=>l.value=!0,["stop"])),onKeyup:ue(N,["enter"])},null,544),[[le,r.value]]),l.value?(H(),V("button",{key:0,class:"ml-3 animate-back-in-right border-1 rounded-2xl px-4",onClick:N}," \u53D1\u9001 ")):ie("",!0)]),d("ul",we,[d("li",ke,[M(" \u5FAA\u73AF\u64AD\u653E\uFF1A "),z(d("input",{"onUpdate:modelValue":n[2]||(n[2]=s=>b.value=s),type:"checkbox"},null,512),[[j,b.value]])]),d("li",_e,[M(" \u64CD\u4F5C\u5F39\u5E55\uFF1A "),d("button",{class:"border-1 rounded-lg p-1 text-sm",onClick:n[3]||(n[3]=(...s)=>u.value.play&&u.value.play(...s))}," \u64AD\u653E "),d("button",{class:"border-1 rounded-lg p-1 text-sm",onClick:n[4]||(n[4]=(...s)=>u.value.pause&&u.value.pause(...s))}," \u6682\u505C "),d("button",{class:"border-1 rounded-lg p-1 text-sm",onClick:n[5]||(n[5]=(...s)=>u.value.stop&&u.value.stop(...s))}," \u505C\u6B62 ")]),d("li",Se,[M(" \u9690\u85CF\u5F39\u5E55\uFF1A "),z(d("input",{"onUpdate:modelValue":n[6]||(n[6]=s=>k.value=s),type:"checkbox"},null,512),[[j,k.value]])])])]),d("div",Ce,[de(P(I),{ref_key:"dmRef",ref:u,danmus:_.value,"onUpdate:danmus":n[7]||(n[7]=s=>_.value=s),class:"h-full w-full","use-slot":"",loop:b.value,speeds:200,channels:0,top:5,"is-suspend":!0},{dm:re(({danmu:s})=>[d("div",Ne,[d("img",{class:"h-[28px] rounded-full",src:P(pe)(s.avatar),alt:"avatar"},null,8,Le),d("span",Be,me(`${s.nickname} : ${s.content}`),1)])]),_:1},8,["danmus","loop"])])],4))}},ze=ee(Ee,[["__scopeId","data-v-9a91cdb8"]]);export{ze as default};
