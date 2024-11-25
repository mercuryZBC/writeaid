import axios from "axios";
import {baseApi} from "./baseService";

// 创建 axios 实例，配置基础 URL
const apiWithUtil = axios.create({
    ...baseApi.defaults,
    headers:{
        ...baseApi.defaults.headers,
    }
});

// 登录页面验证码请求
export const  getCaptcha=()=>{
    return apiWithUtil.get("/util/getCaptcha")
}
