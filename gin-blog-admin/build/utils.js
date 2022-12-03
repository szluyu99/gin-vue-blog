// import fs from 'fs'
import path from 'path'
// import dotenv from 'dotenv'

/**
 * * 项目根路径
 * @descrition 结尾不带 '/'
 */
export function getRootPath() {
  return path.resolve(process.cwd()) // current work directory
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

// const httpsReg = /^https:\/\//

// 处理环境变量对象为可直接使用的值
// 'true' -> true
// '3000' -> 3000
export function convertEnv(envObj) {
  const result = {}
  if (!envObj) return result
  for (const envKey in envObj) {
    let envVal = envObj[envKey]
    if (['true', 'false'].includes(envVal)) envVal = envVal === 'true'
    if (['VITE_PORT'].includes(envKey)) envVal = +envVal
    result[envKey] = envVal
  }
  return result
}

/**
 * 获取当前环境下生效的配置文件名
 */
// function getConfFiles() {
//   const script = process.env.npm_lifecycle_script
//   const reg = new RegExp('--mode ([a-z_\\d]+)')
//   const result = reg.exec(script)
//   if (result) {
//     const mode = result[1]
//     return ['.env', '.env.local', `.env.${mode}`]
//   }
//   return ['.env', '.env.local', '.env.production']
// }

// export function getEnvConfig(match = 'VITE_', confFiles = getConfFiles()) {
//   let envConfig = {}
//   confFiles.forEach((item) => {
//     try {
//       if (fs.existsSync(path.resolve(process.cwd(), item))) {
//         const env = dotenv.parse(fs.readFileSync(path.resolve(process.cwd(), item)))
//         envConfig = { ...envConfig, ...env }
//       }
//     } catch (e) {
//       console.error(`Error in parsing ${item}`, e)
//     }
//   })
//   const reg = new RegExp(`^(${match})`)
//   Object.keys(envConfig).forEach((key) => {
//     if (!reg.test(key)) {
//       Reflect.deleteProperty(envConfig, key)
//     }
//   })
//   return envConfig
// }
