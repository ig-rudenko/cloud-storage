import axios from "axios";
import store from "./store";
import router from "./vue-router"

const instance = axios.create({
    baseURL: "http://127.0.0.1:8080",
    headers: {
        "Content-Type": "application/json",
    },
});

// добавляем interceptor для запросов
instance.interceptors.request.use(
    (config) => {
        // получаем access token из хранилища
        const token = store.state.auth.user.accessToken;
        console.log("token", token)
        if (token) {
            // добавляем токен в заголовок Authorization
            config.headers["Authorization"] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// добавляем interceptor для ответов
instance.interceptors.response.use(
    (response) => {
        return response;
    },
    async (error) => {
        // если статус ответа равен 401, значит токен просрочен
        if (error.response.status === 401) {
            // получаем refresh token из хранилища
            const refresh_token = store.state.auth.user.refreshToken;
            try {
                // отправляем запрос на обновление токена
                const res = await axios.post("/token", { refresh_token });
                // получаем новый access token из ответа
                const new_token = res.data.access_token;
                // сохраняем новый токен в хранилище с помощью мутации refreshToken
                store.commit("auth/refreshToken", new_token);
                // повторяем оригинальный запрос с новым токеном
                error.config.headers["Authorization"] = `Bearer ${new_token}`;
                return instance.request(error.config);
            } catch (err) {
                // если не удалось обновить токен, то выходим из системы и переходим на страницу входа
                store.dispatch("auth/logout");
                router.push("/login");
            }
        }
        return Promise.reject(error);
    }
);

export default instance;