import axios from "axios";
import {baseApi} from "./baseService";
import * as jwt from '../util/jwt'

// 创建 axios 实例，配置基础 URL
const apiWithDocument = axios.create({
    ...baseApi.defaults,
    headers:{
        ...baseApi.defaults.headers,
    }
});

// 设置请求拦截器
apiWithDocument.interceptors.request.use(
    (config) => {
        const token = jwt.getToken();
        if(token){
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },(error) => Promise.reject(error));

// 设置响应拦截器
apiWithDocument.interceptors.response.use(
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

export const createDocument = (kb_id,doc_name) => {
    return apiWithDocument.post("/knowledge/createKnowledgeBase",
        {
            "kb_id":kb_id,
            "doc_name":doc_name,
        })
}
