import axios from "axios";
import {baseApi} from "./baseService";
import * as jwt from '../util/jwt'

// 创建 axios 实例，配置基础 URL
const apiWithKnowledge = axios.create({
    ...baseApi.defaults,
    headers:{
        ...baseApi.defaults.headers,
    }
});

// 设置请求拦截器
apiWithKnowledge.interceptors.request.use(
    (config) => {
        const token = jwt.getToken();
        if(token){
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },(error) => Promise.reject(error));

// 设置响应拦截器
apiWithKnowledge.interceptors.response.use(
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

export const createKnowledgeBase = (kbName,kbDescription,isPublic) => {
    console.log(kbName,kbDescription,isPublic)
    return apiWithKnowledge.post("/knowledge/createKnowledgeBase",
        {
            "kb_name":kbName,
            "kb_description":kbDescription,
            "kb_is_public":isPublic,
        })
}

// 获取当前用户创建的知识库列表
export const getKnowledgeBaseList = () => {
    return apiWithKnowledge.get("/knowledge/getKnowledgeBaseList");
};



export const deleteKnowledgeBase = (kbId) => {
    return apiWithKnowledge.post("/knowledge/deleteKnowledgeBase",
        {
            "kb_id":kbId,
        })
}

export const updateKnowledgeBase = (kbId, kbName, kbDescription, isPublic) => {
    return apiWithKnowledge.post("/knowledge/updateKnowledgeBase",
        {
            "kb_id": kbId,               // 知识库的ID
            "kb_name": kbName,           // 新的知识库名称
            "kb_description": kbDescription, // 新的知识库简介
            "kb_is_public": isPublic,    // 更新后的是否公开状态
        });
}

