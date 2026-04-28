import type { App, DirectiveBinding } from 'vue'
import { useUserStore } from '@/stores/user'

type PermissionValue = string | string[] | undefined

function applyPermission(el: HTMLElement, binding: DirectiveBinding<PermissionValue>) {
  const userStore = useUserStore()
  const allowed = userStore.hasPermission(binding.value)
  el.style.display = allowed ? (el.dataset.permissionDisplay || '') : 'none'
}

export function setupPermissionDirective(app: App) {
  app.directive('permission', {
    mounted(el: HTMLElement, binding: DirectiveBinding<PermissionValue>) {
      el.dataset.permissionDisplay = el.style.display || ''
      applyPermission(el, binding)
    },
    updated(el: HTMLElement, binding: DirectiveBinding<PermissionValue>) {
      applyPermission(el, binding)
    },
  })
}
