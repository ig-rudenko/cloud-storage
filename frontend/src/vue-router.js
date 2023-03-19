// импортируем библиотеку vue-router
import {createRouter} from "vue-router";

// импортируем компоненты для регистрации и входа
import Register from "./components/Register.vue";
import Login from "./components/Login.vue";
import Files from "@/components/Files.vue";

const routes = [
    {
        path: "/register", // путь для регистрации
        component: Register, // компонент для регистрации
    },
    {
        path: "/login", // путь для входа
        component: Login, // компонент для входа
    },
    {
        path: "/", // путь для входа
        component: Files, // компонент для входа
    },
]

// создаем объект роутера с массивом маршрутов
export default function (history) {
    return createRouter({
        history,
        routes
    })
}