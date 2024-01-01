<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu } from 'naive-ui'
import { usePermissionStore, useTagStore } from '@/store'
import { renderIcon } from '@/utils'

const router = useRouter() // 获取路由信息
const curRoute = useRoute() // 进行路由跳转
const permissionStore = usePermissionStore()
const tagStore = useTagStore()

const activeKey = computed(() => curRoute.meta?.activeMenu || curRoute.name)

const menuOptions = computed(() => {
  return permissionStore.menus
    .map(item => getMenuItem(item))
    .sort((a, b) => a.order - b.order)
})

// 点击标签, 自动展开菜单栏, 选中对应菜单
const menuRef = ref(null)
watch(curRoute, async () => {
  await nextTick()
  menuRef.value.showOption()
})

function resolvePath(basePath, path) {
  if (isExternal(path)) {
    return path
  }
  return (
    `/${[basePath, path]
      .filter(path => !!path && path !== '/')
      .map(path => path.replace(/(^\/)|(\/$)/g, ''))
      .join('/')}`
  )
}

// 根据路由获取菜单项
function getMenuItem(route, basePath = '') {
  let menuItem = {
    label: route.meta?.title || route.name,
    key: route.name,
    path: resolvePath(basePath, route.path),
    icon: getIcon(route.meta),
    order: route.meta?.order || 0,
  }

  const visibleChildren = route.children?.filter(e => e.name && !e.isHidden) ?? []

  if (!visibleChildren.length) {
    return menuItem
  }

  // 目录不展示子菜单
  if (route.isCatalogue) {
    const singleRoute = visibleChildren[0]
    menuItem = {
      label: singleRoute.meta?.title || singleRoute.name,
      key: singleRoute.name,
      path: resolvePath(menuItem.path, singleRoute.path),
      icon: getIcon(singleRoute.meta),
      order: menuItem.order,
    }
    const visibleItems = singleRoute.children
      ?.filter(item => item.name && !item.isHidden) ?? []
    if (visibleItems.length === 1) {
      menuItem = getMenuItem(visibleItems[0], menuItem.path)
    }
    else if (visibleItems.length > 1) {
      menuItem.children = visibleItems
        .map(item => getMenuItem(item, menuItem.path))
        .sort((a, b) => a.order - b.order)
    }
  }
  else {
    menuItem.children = visibleChildren
      .map(item => getMenuItem(item, menuItem.path))
      .sort((a, b) => a.order - b.order)
  }

  return menuItem
}

function getIcon(meta) {
  if (meta?.icon) {
    return renderIcon(meta.icon, { size: 18 })
  }
  return null
}

function handleMenuSelect(_, item) {
  if (isExternal(item.path)) {
    window.open(item.path)
    return
  }

  if (item.path === curRoute.path) {
    return
  }

  // 如果 tagStore 中没有该 tag, 需要重新渲染
  if (!tagStore.tags.some(e => e.path === item.path)) {
    tagStore.updateAliveKey(item.key)
  }
  router.push(item.path)
}

/**
 * 是否是外链
 * @param {string} path
 * @returns {boolean} 是否是外链
 */
function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}
</script>

<template>
  <NMenu
    ref="menuRef"
    class="side-menu"
    :indent="18"
    :collapsed-icon-size="22"
    :collapsed-width="64"
    :options="menuOptions"
    :value="activeKey"
    @update:value="handleMenuSelect"
  />
</template>

<style lang="scss">
.side-menu:not(.n-menu--collapsed) {
  .n-menu-item-content {
    &::before {
      left: 5px;
      right: 5px;
    }

    &.n-menu-item-content--selected,
    &:hover {
      &::before {
        border-left: 4px solid var(--primary-color);
      }
    }
  }
}
</style>
