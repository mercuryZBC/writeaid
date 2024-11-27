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

export const createDocument = (kb_id,doc_title) => {
    return apiWithDocument.post("/document/createDocument",
        {
            "kb_id":kb_id,
            "doc_title":doc_title,
        })
}

export const getDocumentListByKbId = (kb_id) => {
    return apiWithDocument.get(`/document/getDocumentListByKbId/${kb_id}`)
}

export const getDocumentById = (doc_id) => {
    return apiWithDocument.get(`/document/getDocument/${doc_id}`)
}

export const deleteDocument = (doc_id) => {
    return apiWithDocument.delete(`/document/deleteDocument/${doc_id}`)
}

export const updateDocumentById = (doc_id, doc_file_form_data) => {
    return apiWithDocument.put(`/document/updateDocument/${doc_id}`, doc_file_form_data, {
        headers: {
            "Content-Type": "multipart/form-data",
        },
    });
};

export const getRecentViewDocuments = (start,offset) =>{
    return apiWithDocument.get(`/document/recentViewDocument?start=${start}&offset=${offset}`);
}

