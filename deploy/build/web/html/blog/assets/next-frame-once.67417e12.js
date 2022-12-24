let n=[];const a=new WeakMap;function s(){n.forEach(e=>e(...a.get(e))),n=[]}function r(e,...t){a.set(e,t),!n.includes(e)&&n.push(e)===1&&requestAnimationFrame(s)}export{r as b};
