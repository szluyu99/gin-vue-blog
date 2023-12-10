import{_ as R}from"./CrudModal-d5d9f5ee.js";import{_ as q}from"./UploadOne-7994ca3e.js";import{_ as B}from"./CommonPage-60900c13.js";import{e as p,p as O,l as n,o as c,f as P,w as i,j as a,h as V,i as e,g as o,W as h,a8 as $,a7 as j,aZ as L,q as T,N as z,a_ as E,y as A,t as M,s as b}from"./index-7d3bcf57.js";import{u as W}from"./useCRUD-87b8d7a3.js";import{N as Z}from"./Image-da15fe39.js";import{N as G,a as d}from"./FormItem-f3617c31.js";import{N as f}from"./Input-41b621cb.js";import"./Upload-c098fdfc.js";import"./Add-58674542.js";import"./AppPage-e3cb8e40.js";import"./utils-540fee40.js";const H=o("i",{class:"i-material-symbols:add"},null,-1),J={class:"flex flex-wrap justify-between"},K={class:"absolute right-10 top-5 text-white"},Q=o("span",{class:"i-ion:ellipsis-horizontal h-20 w-20 text-white hover:text-blue"},null,-1),X={class:"text-15"},Y=o("div",{class:"h-0 w-300"},null,-1),ee=o("div",{class:"h-0 w-300"},null,-1),te=o("div",{class:"h-0 w-300"},null,-1),le={class:"w-full flex items-center justify-between"},_e={__name:"index",setup(ae){const{modalVisible:u,modalTitle:w,modalLoading:x,handleAdd:y,handleDelete:N,handleEdit:k,handleSave:U,modalForm:s,modalFormRef:F}=W({name:"页面",initForm:{},doCreate:n.saveOrUpdatePage,doDelete:n.deletePage,doUpdate:n.saveOrUpdatePage,refresh:g}),v=p([]),m=p(!1),_=p(null);O(async()=>{g()});async function g(){const r=await n.getPages();v.value=r.data}function C(r){m.value=!0,_.value.previewImg=r,setTimeout(()=>m.value=!1,1e3)}function S(r,t){r==="edit"?k(t):r==="delete"&&N([t.id])}const D=[{label:"编辑",key:"edit",icon:()=>b("i",{class:"i-mingcute:edit-2-line"})},{label:"删除",key:"delete",icon:()=>b("i",{class:"i-mingcute:delete-back-line"})}];return(r,t)=>(c(),P(B,{title:"页面管理"},{action:i(()=>[a(e(z),{type:"primary",onClick:e(y)},{icon:i(()=>[H]),default:i(()=>[V(" 新建页面 ")]),_:1},8,["onClick"])]),default:i(()=>[o("div",J,[(c(!0),h(j,null,$(v.value,l=>(c(),h("div",{key:l.id,class:"relative my-10 w-300 cursor-pointer text-center"},[o("div",K,[a(e(E),{options:D,onSelect:I=>S(I,l)},{default:i(()=>[Q]),_:2},1032,["onSelect"])]),a(e(Z),{src:e(A)(l.cover),height:"170",width:"300","img-props":{style:{"border-radius":"5px"}}},null,8,["src"]),o("p",X,M(l.name),1)]))),128)),Y,ee,te]),a(R,{visible:e(u),"onUpdate:visible":t[6]||(t[6]=l=>T(u)?u.value=l:null),width:"550px",title:e(w),loading:e(x),onSave:e(U)},{default:i(()=>[a(e(G),{ref_key:"modalFormRef",ref:F,"label-placement":"left","label-align":"left","label-width":80,model:e(s)},{default:i(()=>[a(e(d),{label:"页面名称",path:"name",rule:{required:!0,message:"请输入页面名称",trigger:["input","blur"]}},{default:i(()=>[a(e(f),{value:e(s).name,"onUpdate:value":t[0]||(t[0]=l=>e(s).name=l),placeholder:"页面名称"},null,8,["value"])]),_:1}),a(e(d),{label:"页面标签",path:"label",rule:{required:!0,message:"请输入页面标签",trigger:["input","blur"]}},{default:i(()=>[a(e(f),{value:e(s).label,"onUpdate:value":t[1]||(t[1]=l=>e(s).label=l),placeholder:"页面标签"},null,8,["value"])]),_:1}),a(e(d),{label:"页面封面",path:"cover",rule:{required:!0,message:"请上传封面图片",trigger:["input","blur"]}},{default:i(()=>[o("div",le,[a(q,{ref_key:"uploadOneRef",ref:_,preview:e(s).cover,"onUpdate:preview":t[2]||(t[2]=l=>e(s).cover=l),width:300,onFinish:t[3]||(t[3]=l=>e(s).cover=l)},null,8,["preview"]),o("span",{class:L(["i-uiw:reload h-20 w-20 cursor-pointer",m.value?"animate-spin":""]),onClick:t[4]||(t[4]=l=>C(e(s).cover))},null,2)])]),_:1}),a(e(d),{label:"封面链接",path:"cover",rule:{required:!0,message:"请输入封面链接",trigger:["input","blur"]}},{default:i(()=>[a(e(f),{value:e(s).cover,"onUpdate:value":t[5]||(t[5]=l=>e(s).cover=l),type:"textarea",placeholder:"图片上传成功自动生成，或者直接复制外链"},null,8,["value"])]),_:1})]),_:1},8,["model"])]),_:1},8,["visible","title","loading","onSave"])]),_:1}))}};export{_e as default};
