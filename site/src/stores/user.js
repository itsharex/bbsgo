import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/api";

export const useUserStore = defineStore("user", () => {
  const token = ref(localStorage.getItem("token") || "");
  const user = ref(JSON.parse(localStorage.getItem("user") || "null"));

  const isLoggedIn = computed(() => !!token.value && !!user.value);

  async function login(loginData) {
    const res = await api.post("/login", loginData);
    token.value = res.token;
    user.value = res.user;
    localStorage.setItem("token", res.token);
    localStorage.setItem("user", JSON.stringify(res.user));
  }

  async function register(registerData) {
    const res = await api.post("/register", registerData);
    token.value = res.token;
    user.value = res.user;
    localStorage.setItem("token", res.token);
    localStorage.setItem("user", JSON.stringify(res.user));
  }

  function logout() {
    token.value = "";
    user.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  }

  async function refreshProfile() {
    const res = await api.get("/user/profile");
    user.value = res;
    localStorage.setItem("user", JSON.stringify(res));
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    logout,
    refreshProfile,
  };
});
