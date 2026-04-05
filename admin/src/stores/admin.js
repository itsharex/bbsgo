import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api, { authApi } from "@/api";

export const useAdminStore = defineStore("admin", () => {
  const user = ref(JSON.parse(localStorage.getItem("admin_user") || "null"));

  const isLoggedIn = computed(() => !!user.value);
  const isAdmin = computed(() => user.value?.role >= 2);

  async function login(loginData) {
    const res = await api.post("/login", loginData);
    if (res.user.role < 2) {
      throw new Error("需要管理员权限");
    }
    user.value = res.user;
    localStorage.setItem("admin_user", JSON.stringify(res.user));
  }

  async function logout() {
    try {
      await authApi.logout();
    } catch (e) {
      // 即使 API 调用失败，也清除本地状态
    }
    user.value = null;
    localStorage.removeItem("admin_user");
  }

  return {
    user,
    isLoggedIn,
    isAdmin,
    login,
    logout,
  };
});
