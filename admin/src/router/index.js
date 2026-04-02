import { createRouter, createWebHistory } from "vue-router";
import { useAdminStore } from "@/stores/admin";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
  },
  {
    path: "/",
    component: () => import("@/views/Layout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "Dashboard",
        component: () => import("@/views/Dashboard.vue"),
      },
      {
        path: "users",
        name: "Users",
        component: () => import("@/views/Users.vue"),
      },
      {
        path: "forums",
        name: "Forums",
        component: () => import("@/views/Forums.vue"),
      },
      {
        path: "topics",
        name: "Topics",
        component: () => import("@/views/Topics.vue"),
      },
      {
        path: "posts",
        name: "Posts",
        component: () => import("@/views/Posts.vue"),
      },
      {
        path: "tags",
        name: "Tags",
        component: () => import("@/views/Tags.vue"),
      },
      {
        path: "reports",
        name: "Reports",
        component: () => import("@/views/Reports.vue"),
      },
      {
        path: "announcements",
        name: "Announcements",
        component: () => import("@/views/Announcements.vue"),
      },
      {
        path: "config",
        name: "Config",
        component: () => import("@/views/Config.vue"),
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/Settings.vue"),
      },
      {
        path: "change-password",
        name: "ChangePassword",
        component: () => import("@/views/ChangePassword.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const adminStore = useAdminStore();
  if (to.meta.requiresAuth && !adminStore.isLoggedIn) {
    next("/login");
  } else if (to.path === "/login" && adminStore.isLoggedIn) {
    next("/");
  } else {
    next();
  }
});

export default router;
