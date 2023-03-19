import axios from "axios";

const instance = axios.create({
    baseURL: "/api",
    headers: {
        "Content-Type": "application/json",
    },
    mode: 'no-cors', // disable CORS check
    withCredentials: false // do not send cookies
});

export default instance;