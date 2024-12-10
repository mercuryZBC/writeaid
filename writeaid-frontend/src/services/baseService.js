import axios from "axios";
import * as config from "../config/config";
import * as jwt from "../util/jwt";


const baseApi = axios.create({
    baseURL: config.baseURL +"/api",
    timeout: config.requestTimeout,
})
// 设置请求拦截器
baseApi.interceptors.request.use(
    (config) => {
    const token = jwt.getToken();
    if(token){
        config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
},(error) => Promise.reject(error));

// 设置响应拦截器
baseApi.interceptors.response.use(
    (response) => response,
    (error) => {
        if(error.response === 401 && error.status === 401) {
            // 如果状态码为401，说明token过期或无效
            localStorage.removeItem('jwt_token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export {baseApi}
