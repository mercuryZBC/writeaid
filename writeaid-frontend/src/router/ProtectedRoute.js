import React from "react";
import { Navigate } from "react-router-dom";
import {getToken} from "../util/jwt";


const ProtectedRoute = ({ children }) => {
    // 检查 Token
    const isAuthenticated = getToken() !== "";

    // 如果未登录，重定向到登录页
    return isAuthenticated ? children : <Navigate to="/login" />;
};

export default ProtectedRoute;
