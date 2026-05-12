import { createApp } from 'vue'
import { createPinia } from 'pinia'


import './assets/main.css'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth';

const app = createApp(App);

app.use(createPinia()) // Register Pinia first
app.use(router)

const auth = useAuthStore();
auth.loadAuth();

app.mount('#app')
