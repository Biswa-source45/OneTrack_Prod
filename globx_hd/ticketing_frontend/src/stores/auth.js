import { defineStore } from 'pinia';
import { decodeJWT } from '../utils/jwt';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    userType: null,
    firstLogin: false,
    user: null,
  }),
  actions: {
    setAuth(token, userType, firstLogin, user) {
      this.token = token;
  // Always decode from JWT for reliability
  const claims = decodeJWT(token);
  this.userType = claims?.role ? claims.role.toLowerCase().trim() : userType; // 'admin' or 'contact'
      this.firstLogin = firstLogin;
      this.user = user;
      localStorage.setItem('token', token);
      localStorage.setItem('userType', this.userType);
      localStorage.setItem('firstLogin', firstLogin);
      localStorage.setItem('user', JSON.stringify(user));
    },
    clearAuth() {
      this.token = null;
      this.userType = null;
      this.firstLogin = false;
      this.user = null;
      localStorage.clear();
    },
    loadAuth() {
      this.token = localStorage.getItem('token');
      this.firstLogin = localStorage.getItem('firstLogin') === 'true';
      const userStr = localStorage.getItem('user');
      try {
        this.user = userStr && userStr !== "undefined" ? JSON.parse(userStr) : null;
      } catch (e) {
        this.user = null;
      }
      // Decode JWT to get userType/role
      if (this.token) {
        const claims = decodeJWT(this.token);
        this.userType = claims?.role ? claims.role.toLowerCase().trim() : localStorage.getItem('userType');
      } else {
        this.userType = null;
      }
    }
  }
});
