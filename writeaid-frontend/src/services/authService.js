import axios from "axios";
import {baseApi} from "./baseService";

const apiWithAuth = axios.create({
    ...baseApi.defaults,
    headers:{
        ...baseApi.defaults.headers,
    }
});

// 登录请求
export const loginRequest = (email, password,captchaId,captchaValue) => {
    console.log(email,password,captchaId,captchaValue);
    return apiWithAuth.post("/auth/login", {
        "email":email,
        "password":password,
        "captchaId":captchaId,
        "captchaValue": captchaValue
    });
};

// 注册请求
export const registerRequest = (email,nickname,password,captchaId,captchaValue) => {
    return apiWithAuth.post("/auth/register", {
        "email": email,
        "nickname": nickname,
        "password": password,
        "captchaId": captchaId,
        "captchaValue":captchaValue,
    });
};
