// * 前端导出, 传入文件内容和文件名称
export function downloadFile(content, fileName) {
  const aEle = document.createElement('a') // 创建下载链接
  aEle.download = fileName // 设置下载的名称
  aEle.style.display = 'none'// 隐藏的可下载链接
  // 字符内容转变成 blob 地址
  const blob = new Blob([content])
  aEle.href = URL.createObjectURL(blob)
  // 绑定点击时间
  document.body.appendChild(aEle)
  aEle.click()
  // 然后移除
  document.body.removeChild(aEle)
}

// 使用 iframe 实现下载
// export function downloadFile(url) {
//   // url = 'https://static.talkxj.com/markdown/博客技术总结.md'
//   const iframe = document.createElement('iframe')
//   iframe.style.display = 'none' // 防止影响页面
//   iframe.style.height = 0 // 防止影响页面
//   iframe.src = convertImgUrl(url)
//   document.body.appendChild(iframe)
//   setTimeout(() => {
//     iframe.remove()
//   }, 5 * 60 * 1000)
// }

// 通过 a 标签实现下载
// export function downloadFile(url) {
//   const a = document.createElement('a')
//   a.href = convertImgUrl(url) // 文件链接
//   a.target = '_blank'
//   a.download = '下载的文件名' // 文件名，跨域资源 download 无效
//   a.click()
//   a.remove()
// }

// fetch 读取文件内容再转二进制流下载 FIXME: 跨域
// export function downloadFile(url) {
//   const link = document.createElement('a')
//   fetch(convertImgUrl(url))
//     .then(res => res.blob())
//     .then((blob) => { // 将链接地址字符内容转变成 blob 地址
//       link.href = URL.createObjectURL(blob)
//       link.download = ''
//       document.body.appendChild(link)
//       link.click()
//       document.body.removeChild(link)
//     })
// }
