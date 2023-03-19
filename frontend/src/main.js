import { createApp } from 'vue'
import App from '@/App.vue'
import {createWebHistory} from "vue-router";
import createRouter from "@/vue-router";

import store from "@/store"; // #

import setupInterceptors from './services/setupInterceptors';

setupInterceptors(store);

import '@/assets/main.css'

createApp(App)
    .use(store)  // #
    .use(createRouter(createWebHistory()))
    .mount('#app')