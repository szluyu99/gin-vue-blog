import path from 'path'

/**
 * * 项目根路径
 * @descrition 结尾不带 '/'
 */
export function getRootPath() {
  return path.resolve(process.cwd()) // cwd - current work directory
}

/**
 * * 项目src路径
 * @param srcName src目录名称(默认: "src")
 * @descrition 结尾不带 '/'
 */
export function getSrcPath(srcName = 'src') {
  // path.resolve 可以理解为拼接
  // path.resolve('/root', 'src') -> /root/src
  return path.resolve(getRootPath(), srcName)
}

// 处理环境变量对象为可直接使用的值
// 'true' -> true
// '3000' -> 3000
export function convertEnv(envObj) {
  const result = {}
  if (!envObj)
    return result
  for (const envKey in envObj) {
    let envVal = envObj[envKey]
    if (['true', 'false'].includes(envVal))
      envVal = envVal === 'true'
    if (['VITE_PORT'].includes(envKey))
      envVal = +envVal
    result[envKey] = envVal
  }
  return result
}
