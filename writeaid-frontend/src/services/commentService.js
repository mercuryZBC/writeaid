
// 创建 axios 实例，配置基础 URL
import axios from "axios";
import {baseApi} from "./baseService";
import * as jwt from "../util/jwt";

const apiWithComment = axios.create({
    ...baseApi.defaults,
    headers:{
        ...baseApi.defaults.headers,
    }
});

// 设置请求拦截器
apiWithComment.interceptors.request.use(
    (config) => {
        const token = jwt.getToken();
        if(token){
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },(error) => Promise.reject(error));

// 设置响应拦截器
apiWithComment.interceptors.response.use(
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

export const getDocumentRootComment = (docId,page,offset) =>{
    return apiWithComment.get(`/comment/getDocumentRootComment/${docId}?page=${page}&page_size=${offset}`)
}

export const deleteDocumentComment = (commentId) =>{
    return apiWithComment.delete(`/comment/deleteDocumentComment/${commentId}`)
}

export const updateDocumentComment = (commentId,commentContent) =>{
    return apiWithComment.post(`/comment/updateDocumentComment`,{
        'comment_id':commentId,
        'comment_content': commentContent
    })
}

export const createDocumentComment = (docId,commentContent) =>{
    return apiWithComment.post(`/comment/createDocumentComment`,{
        'doc_id':docId,
        'comment_content': commentContent
    })
}

export const replyDocumentComment = (rootId,parentId,docId,commentContent) =>{
    return apiWithComment.post(`/comment/replyDocumentComment`,{
        'doc_id':docId,
        'root_id':rootId,
        'parent_id': parentId,
        'comment_content': commentContent
    })
}

export const getChildrenComment = (rootId) =>{
    return apiWithComment.get(`/comment/getChildrenComment/${rootId}`)
}
