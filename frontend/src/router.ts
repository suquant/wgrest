import Vue from 'vue'
import Router from 'vue-router'
import Layout from '@/layout/index.vue'

Vue.use(Router)

/*
  redirect:                      if set to 'noredirect', no redirect action will be trigger when clicking the breadcrumb
  meta: {
    title: 'title'               the name showed in subMenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon showed in the sidebar
    breadcrumb: false            if false, the item will be hidden in breadcrumb (default is true)
    hidden: true                 if true, this route will not show in the sidebar (default is false)
  }
*/

export default new Router({
  mode: 'history',
  scrollBehavior: (to, from, savedPosition) => {
    if (savedPosition) {
      return savedPosition
    } else {
      return { x: 0, y: 0 }
    }
  },
  base: process.env.BASE_URL,
  routes: [
    // {
    //   path: '/login',
    //   component: () => import(/* webpackChunkName: "login" */ '@/views/login/index.vue'),
    //   meta: { hidden: true }
    // },
    {
      path: '/404',
      component: () => import(/* webpackChunkName: "404" */ '@/views/404.vue'),
      meta: { hidden: true }
    },
    {
      path: '/',
      component: Layout,
      redirect: '/devices',
      children: [
        {
          path: 'devices',
          name: 'Device',
          component: () => import('@/views/dashboard/index.vue'),
          meta: {
            title: 'Devices',
            icon: 'dashboard'
          }
        },
        {
          path: 'devices/:id',
          name: 'Peer',
          component: () => import('@/views/dashboard/_id.vue'),
          meta: { hidden: true }
        }
      ]
    },
    {
      path: '*',
      redirect: '/404',
      meta: { hidden: true }
    }
  ]
})
