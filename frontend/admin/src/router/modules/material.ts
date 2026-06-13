const Layout = () => import("@/layout/index.vue");

export default {
  path: "/material",
  name: "Material",
  component: Layout,
  redirect: "/material/city",
  meta: {
    icon: "ep/office-building",
    title: "业务管理",
    rank: 2
  },
  children: [
    {
      path: "/material/city",
      name: "MaterialCity",
      component: () => import("@/views/material/city/index.vue"),
      meta: {
        title: "城市管理",
        roles: ["admin"]
      }
    },
    {
      path: "/material/merchant",
      name: "MaterialMerchant",
      component: () => import("@/views/material/merchant/index.vue"),
      meta: {
        title: "商家管理",
        roles: ["admin"]
      }
    },
    {
      path: "/material/order",
      name: "MaterialOrder",
      component: () => import("@/views/material/order/index.vue"),
      meta: {
        title: "订单监控",
        roles: ["admin"]
      }
    }
  ]
} satisfies RouteConfigsTable;
