/* eslint-disable @typescript-eslint/no-this-alias */
/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-redeclare */
/* eslint-disable @typescript-eslint/no-use-before-define */
/* eslint-disable block-scoped-var */
/* eslint-disable max-statements-per-line */
/* eslint-disable no-cond-assign */
/* eslint-disable no-sequences */
/* eslint-disable no-var */
/* eslint-disable vars-on-top */
/* eslint-disable no-unused-expressions */
/* eslint-disable no-undef */

!(function (t, e) {
  typeof exports == 'object' && typeof module == 'object'
    ? (module.exports = e())
    : typeof define == 'function' && define.amd
      ? define('TagCloud', [], e)
      : typeof exports == 'object'
        ? (exports.TagCloud = e())
        : (t.TagCloud = e())
})(typeof self != 'undefined' ? self : this, () => {
  return (function (t) {
    function e(o) {
      if (n[o])
        return n[o].exports
      const r = (n[o] = { i: o, l: !1, exports: {} })
      return t[o].call(r.exports, r, r.exports, e), (r.l = !0), r.exports
    }
    var n = {}
    return (
      (e.m = t),
      (e.c = n),
      (e.d = function (t, n, o) {
        e.o(t, n)
          || Object.defineProperty(t, n, {
            configurable: !1,
            enumerable: !0,
            get: o,
          })
      }),
      (e.n = function (t) {
        const n
          = t && t.__esModule
            ? function () {
              return t.default
            }
            : function () {
              return t
            }
        return e.d(n, 'a', n), n
      }),
      (e.o = function (t, e) {
        return Object.prototype.hasOwnProperty.call(t, e)
      }),
      (e.p = '/dist/'),
      e((e.s = 1))
    )
  })([
    function (t, e, n) {
      'use strict'
      e.a = {
        name: 'tagCloud',
        props: {
          data: { type: Array, default: [] },
          config: { type: Object, default: null },
        },
        data() {
          return {
            option: {
              radius: 120,
              maxFont: 24,
              color: null,
              rotateAngleXbase: 500,
              rotateAngleYbase: 500,
              hover: !1,
            },
            tagList: [],
          }
        },
        created() {
          this.config != null && (this.option = Object.assign({}, this.option, this.config))
        },
        mounted() {
          this._initTags()
        },
        beforeDestroy() {
          this.timer && (clearInterval(this.timer), (this.timer = null))
        },
        watch: {
          data() {
            const t = this
            this.$nextTick(() => {
              t._initTags()
            })
          },
        },
        methods: {
          _initTags() {
            if (
              ((this.rotateAngleX = Math.PI / this.option.rotateAngleXbase),
              (this.rotateAngleY = Math.PI / this.option.rotateAngleYbase),
              this.option.hover)
            ) {
              const t = this
              this.$refs.wrapper.onmousemove = function (e) {
                (t.rotateAngleY = (e.pageX - this.offsetLeft - this.offsetWidth / 2) / 1e4),
                (t.rotateAngleX = -(e.pageY - this.offsetTop - this.offsetHeight / 2) / 1e4)
              }
            }
            else { this.$refs.wrapper.onmousemove = null }
            for (let e = 0, n = this.data.length; e < n; e++) {
              const o = Math.acos((2 * (e + 1) - 1) / n - 1)
              const r = o * Math.sqrt(n * Math.PI)
              const i = this.option.radius * Math.sin(o) * Math.cos(r)
              const a = this.option.radius * Math.sin(o) * Math.sin(r)
              const s = this.option.radius * Math.cos(o)
              this.option.color
                ? (this.$refs.tag[e].style.color = this.option.color)
                : (this.$refs.tag[e].style.color
                    = `rgb(${
                     Math.round(255 * Math.random())
                     },${
                     Math.round(255 * Math.random())
                     },${
                     Math.round(255 * Math.random())
                     })`)
              const u = { x: i, y: a, z: s, ele: this.$refs.tag[e] }
              this.tagList.push(u)
            }
            const f = this
            this.timer = setInterval(() => {
              for (let t = 0; t < f.tagList.length; t++) {
                f.rotateX(f.tagList[t]),
                f.rotateY(f.tagList[t]),
                f.setPosition(f.tagList[t], f.option.radius, f.option.maxFont)
              }
            }, 20)
          },
          setPosition(t, e, n) {
            this.$refs.wrapper
              && ((t.ele.style.transform
                = `translate(${
                 t.x + this.$refs.wrapper.offsetWidth / 2 - t.ele.offsetWidth / 2
                 }px,${
                 t.y + this.$refs.wrapper.offsetHeight / 2 - t.ele.offsetHeight / 2
                 }px)`),
              (t.ele.style.opacity = t.z / e / 2 + 0.7),
              (t.ele.style.fontSize = `${(t.z / e / 2 + 0.5) * n}px`))
          },
          rotateX(t) {
            const e = Math.cos(this.rotateAngleX)
            const n = Math.sin(this.rotateAngleX)
            const o = t.y * e - t.z * n
            const r = t.y * n + t.z * e
            ;(t.y = o), (t.z = r)
          },
          rotateY(t) {
            const e = Math.cos(this.rotateAngleY)
            const n = Math.sin(this.rotateAngleY)
            const o = t.z * n + t.x * e
            const r = t.z * e - t.x * n
            ;(t.x = o), (t.z = r)
          },
          dbclickTag() {
            if (this.timer) { clearInterval(this.timer), (this.timer = null) }
            else {
              const t = this
              this.timer = setInterval(() => {
                for (let e = 0; e < t.tagList.length; e++) {
                  t.rotateX(t.tagList[e]),
                  t.rotateY(t.tagList[e]),
                  t.setPosition(t.tagList[e], t.option.radius, t.option.maxFont)
                }
              }, 20)
            }
          },
          clickTag(t) {
            this.$emit('clickTag', t)
          },
        },
      }
    },
    function (t, e, n) {
      'use strict'
      Object.defineProperty(e, '__esModule', { value: !0 })
      const o = n(2)
      const r = {
        install(t) {
          typeof window != 'undefined' && window.Vue && (t = window.Vue),
          t.component(o.a.name, o.a)
        },
      }
      e.default = r
    },
    function (t, e, n) {
      'use strict'
      function o(t) {
        n(3)
      }
      const r = n(0)
      const i = n(9)
      const a = n(8)
      const s = o
      const u = a(r.a, i.a, !1, s, 'data-v-7f9ad8d8', null)
      e.a = u.exports
    },
    function (t, e, n) {
      let o = n(4)
      typeof o == 'string' && (o = [[t.i, o, '']]), o.locals && (t.exports = o.locals)
      n(6)('3fb9a8be', o, !0, {})
    },
    function (t, e, n) {
      (e = t.exports = n(5)(!1)),
      e.push([
        t.i,
        '.tag-cloud[data-v-7f9ad8d8]{width:300px;height:300px;position:relative;color:#333;margin:0 auto;text-align:center}.tag-cloud p[data-v-7f9ad8d8]{position:absolute;top:0;left:0;color:#333;font-family:Arial;text-decoration:none;margin:0 10px 15px 0;line-height:18px;text-align:center;font-size:16px;padding:4px 9px;display:inline-block;border-radius:3px}',
        '',
      ])
    },
    function (t, e) {
      function n(t, e) {
        const n = t[1] || ''
        const r = t[3]
        if (!r)
          return n
        if (e && typeof btoa == 'function') {
          const i = o(r)
          return [n]
            .concat(
              r.sources.map((t) => {
                return `/*# sourceURL=${r.sourceRoot}${t} */`
              }),
            )
            .concat([i])
            .join('\n')
        }
        return [n].join('\n')
      }
      function o(t) {
        return (
          `/*# sourceMappingURL=data:application/json;charset=utf-8;base64,${
          btoa(unescape(encodeURIComponent(JSON.stringify(t))))
          } */`
        )
      }
      t.exports = function (t) {
        const e = []
        return (
          (e.toString = function () {
            return this.map((e) => {
              const o = n(e, t)
              return e[2] ? `@media ${e[2]}{${o}}` : o
            }).join('')
          }),
          (e.i = function (t, n) {
            typeof t == 'string' && (t = [[null, t, '']])
            for (var o = {}, r = 0; r < this.length; r++) {
              const i = this[r][0]
              typeof i == 'number' && (o[i] = !0)
            }
            for (r = 0; r < t.length; r++) {
              const a = t[r]
              ;(typeof a[0] == 'number' && o[a[0]])
                || (n && !a[2] ? (a[2] = n) : n && (a[2] = `(${a[2]}) and (${n})`),
                e.push(a))
            }
          }),
          e
        )
      }
    },
    function (t, e, n) {
      function o(t) {
        for (let e = 0; e < t.length; e++) {
          const n = t[e]
          const o = c[n.id]
          if (o) {
            o.refs++
            for (var r = 0; r < o.parts.length; r++) o.parts[r](n.parts[r])
            for (; r < n.parts.length; r++) o.parts.push(i(n.parts[r]))
            o.parts.length > n.parts.length && (o.parts.length = n.parts.length)
          }
          else {
            for (var a = [], r = 0; r < n.parts.length; r++) a.push(i(n.parts[r]))
            c[n.id] = { id: n.id, refs: 1, parts: a }
          }
        }
      }
      function r() {
        const t = document.createElement('style')
        return (t.type = 'text/css'), l.appendChild(t), t
      }
      function i(t) {
        let e
        let n
        let o = document.querySelector(`style[${m}~="${t.id}"]`)
        if (o) {
          if (h)
            return g
          o.parentNode.removeChild(o)
        }
        if (y) {
          const i = p++
          ;(o = d || (d = r())), (e = a.bind(null, o, i, !1)), (n = a.bind(null, o, i, !0))
        }
        else {
          (o = r()),
          (e = s.bind(null, o)),
          (n = function () {
            o.parentNode.removeChild(o)
          })
        }
        return (
          e(t),
          function (o) {
            if (o) {
              if (o.css === t.css && o.media === t.media && o.sourceMap === t.sourceMap)
                return
              e((t = o))
            }
            else { n() }
          }
        )
      }
      function a(t, e, n, o) {
        const r = n ? '' : o.css
        if (t.styleSheet) { t.styleSheet.cssText = x(e, r) }
        else {
          const i = document.createTextNode(r)
          const a = t.childNodes
          a[e] && t.removeChild(a[e]), a.length ? t.insertBefore(i, a[e]) : t.appendChild(i)
        }
      }
      function s(t, e) {
        let n = e.css
        const o = e.media
        const r = e.sourceMap
        if (
          (o && t.setAttribute('media', o),
          v.ssrId && t.setAttribute(m, e.id),
          r
            && ((n += `\n/*# sourceURL=${r.sources[0]} */`),
            (n
              += `\n/*# sourceMappingURL=data:application/json;base64,${
               btoa(unescape(encodeURIComponent(JSON.stringify(r))))
               } */`)),
          t.styleSheet)
        ) { t.styleSheet.cssText = n }
        else {
          for (; t.firstChild;) t.removeChild(t.firstChild)
          t.appendChild(document.createTextNode(n))
        }
      }
      const u = typeof document != 'undefined'
      if (typeof DEBUG != 'undefined' && DEBUG && !u) {
        throw new Error(
          'vue-style-loader cannot be used in a non-browser environment. Use { target: \'node\' } in your Webpack config to indicate a server-rendering environment.',
        )
      }
      const f = n(7)
      var c = {}
      var l = u && (document.head || document.getElementsByTagName('head')[0])
      var d = null
      var p = 0
      var h = !1
      var g = function () {}
      var v = null
      var m = 'data-vue-ssr-id'
      var y
          = typeof navigator != 'undefined' && /msie [6-9]\b/.test(navigator.userAgent.toLowerCase())
      t.exports = function (t, e, n, r) {
        (h = n), (v = r || {})
        let i = f(t, e)
        return (
          o(i),
          function (e) {
            for (var n = [], r = 0; r < i.length; r++) {
              const a = i[r]
              var s = c[a.id]
              s.refs--, n.push(s)
            }
            e ? ((i = f(t, e)), o(i)) : (i = [])
            for (var r = 0; r < n.length; r++) {
              var s = n[r]
              if (s.refs === 0) {
                for (let u = 0; u < s.parts.length; u++) s.parts[u]()
                delete c[s.id]
              }
            }
          }
        )
      }
      var x = (function () {
        const t = []
        return function (e, n) {
          return (t[e] = n), t.filter(Boolean).join('\n')
        }
      })()
    },
    function (t, e) {
      t.exports = function (t, e) {
        for (var n = [], o = {}, r = 0; r < e.length; r++) {
          const i = e[r]
          const a = i[0]
          const s = i[1]
          const u = i[2]
          const f = i[3]
          const c = { id: `${t}:${r}`, css: s, media: u, sourceMap: f }
          o[a] ? o[a].parts.push(c) : n.push((o[a] = { id: a, parts: [c] }))
        }
        return n
      }
    },
    function (t, e) {
      t.exports = function (t, e, n, o, r, i) {
        let a
        let s = (t = t || {})
        const u = typeof t.default
        ;(u !== 'object' && u !== 'function') || ((a = t), (s = t.default))
        const f = typeof s == 'function' ? s.options : s
        e && ((f.render = e.render), (f.staticRenderFns = e.staticRenderFns), (f._compiled = !0)),
        n && (f.functional = !0),
        r && (f._scopeId = r)
        let c
        if (
          (i
            ? ((c = function (t) {
                (t
                  = t
                  || (this.$vnode && this.$vnode.ssrContext)
                  || (this.parent && this.parent.$vnode && this.parent.$vnode.ssrContext)),
                t || typeof __VUE_SSR_CONTEXT__ == 'undefined' || (t = __VUE_SSR_CONTEXT__),
                o && o.call(this, t),
                t && t._registeredComponents && t._registeredComponents.add(i)
              }),
              (f._ssrRegister = c))
            : o && (c = o),
          c)
        ) {
          const l = f.functional
          const d = l ? f.render : f.beforeCreate
          l
            ? ((f._injectStyles = c),
              (f.render = function (t, e) {
                return c.call(e), d(t, e)
              }))
            : (f.beforeCreate = d ? [].concat(d, c) : [c])
        }
        return { esModule: a, exports: s, options: f }
      }
    },
    function (t, e, n) {
      'use strict'
      const o = function () {
        const t = this
        const e = t.$createElement
        const n = t._self._c || e
        return n(
          'div',
          { ref: 'wrapper', staticClass: 'tag-cloud' },
          t._l(t.data, (e, o) => {
            return n(
              'p',
              {
                key: o,
                ref: 'tag',
                refInFor: !0,
                on: {
                  click(n) {
                    t.clickTag(e)
                  },
                  dblclick(n) {
                    t.dbclickTag(e)
                  },
                },
              },
              [t._v(t._s(e.name))],
            )
          }),
        )
      }
      const r = []
      const i = { render: o, staticRenderFns: r }
      e.a = i
    },
  ])
})
// # sourceMappingURL=tag-cloud.js.map
