import React from "react";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import {getToken} from "./util/jwt";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

function App() {
    const isAuthenticated = getToken() !== "";
    return (
        <div>
            <BrowserRouter>
                <Routes>
                    {/* 登录页面 */}
                    <Route path="/login" element={<LoginPage />} />
                    {/* 主页 */}
                    <Route
                        path="/home"
                        element={isAuthenticated ? <HomePage /> : <Navigate to="/login" />}
                    />
                    {/* 默认路由，未匹配时重定向到 /login */}
                    <Route path="*" element={<Navigate to="/home" />} />
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default App;
