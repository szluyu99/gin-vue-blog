import{M as g,o as m,f as v,bJ as w,w as o,g as x,bm as d,j as r,i,N as u,h as f,t as y,n as b,bK as h}from"./index-cfa27ed4.js";const S={class:"flex justify-end space-x-20"},B={__name:"CrudModal",props:{visible:{type:Boolean,required:!0},width:{type:String,default:"600px"},title:{type:String,default:""},showFooter:{type:Boolean,default:!0},loading:{type:Boolean,default:!1},cancelText:{type:String,default:"取消"},okText:{type:String,default:"确定"}},emits:["update:visible","save"],setup(e,{emit:p}){const c=e,n=p,s=g({get:()=>c.visible,set:a=>n("update:visible",a)});return(a,t)=>(m(),v(i(h),{show:s.value,"onUpdate:show":t[2]||(t[2]=l=>s.value=l),style:b({width:e.width}),preset:"card",title:e.title,size:"huge",bordered:!1},w({default:o(()=>[d(a.$slots,"default")]),_:2},[e.showFooter?{name:"footer",fn:o(()=>[x("footer",S,[d(a.$slots,"footer",{},()=>[r(i(u),{onClick:t[0]||(t[0]=l=>s.value=!1)},{default:o(()=>[f(y(e.cancelText),1)]),_:1}),r(i(u),{loading:e.loading,type:"primary",onClick:t[1]||(t[1]=l=>n("save"))},{default:o(()=>[f(y(e.okText),1)]),_:1},8,["loading"])])])]),key:"0"}:void 0]),1032,["show","style","title"]))}};export{B as _};