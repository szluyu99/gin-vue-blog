关于项目中为什么不使用 Prettier，参考 Antfu 大佬： [为什么我不使用 Prettier](https://antfu.me/posts/why-not-prettier-zh)

Eslint 方案采用 [https://github.com/antfu/eslint-config](https://github.com/antfu/eslint-config)，最大化减少配置


## 使用 iconify 图标

首先去图标库地址：[icones](https://icones.js.org/) 找合适的图标

### 1. 结合 unocss 使用

```html
<i i-carbon-sun />
<i class="i-carbon-sun" />
```

### 2. 结合插件 unplugin-icons 自定义标签使用

`<icon-[iconify图标名称]`

```html
<icon-ant-design:fullscreen-exit-outlined  />
<icon-ant-design:fullscreen-outlined />
```

这种方式还支持自定义 svg 图标，本项目自定义 svg 图标固定放在 src/assets/svg 下

`<icon-custom-[svg图标文件名]`

```
<icon-custom-logo />
```

具体配置参看 build/plugin/unplugin.js

### 3. 结合 Naive UI 的 NIcon 组件封装使用

```html
<!-- iconify图标 -->
<TheIcon icon="material-symbols:delete-outline" />
<!-- 自定义svg图标 -->
<TheIcon icon="logo" type="custom" />
```

封装组件参看 src/components/icon
