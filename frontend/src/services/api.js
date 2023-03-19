import axios from "axios";

const instance = axios.create({
    baseURL: "http://localhost:5173",
    headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "http://localhost:8080"
    },
    mode: 'no-cors', // disable CORS check
    withCredentials: false // do not send cookies
});

export default instance;