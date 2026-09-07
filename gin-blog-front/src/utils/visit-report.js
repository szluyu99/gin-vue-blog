import api from '@/api'

/*
访客上报

原来 `POST /report` 全站没人调用, 后台仪表盘的「访问量」和访客地域统计一直是 0。

用 sessionStorage 而不是 localStorage 做去重: 关掉标签页再进来算新的一次访问,
正好是趋势图想要的粒度(按天的访问次数)。同一次会话里刷新页面不重复计数。
*/
const FLAG_KEY = 'visit-reported'

export function reportVisit() {
  // 隐私模式下 sessionStorage 可能直接抛错, 不能让统计挡住页面渲染
  try {
    if (sessionStorage.getItem(FLAG_KEY)) {
      return
    }
    sessionStorage.setItem(FLAG_KEY, '1')
  }
  catch {
    // 读写不了就每次都报一次, 总比一次都不报好
  }

  // 上报失败静默: 统计失败不该给访客任何提示
  return api.report().catch(() => {})
}
