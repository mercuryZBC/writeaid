import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import DocumentViewPage from "./pages/DocumentViewPage";
import DocumentEditPage from "./pages/DocumentEditPage";
import ProtectedRoute from "./router/ProtectedRoute";

function App() {
    return (
        <div>
            <BrowserRouter>
                <Routes>
                    {/* 登录页面 */}
                    <Route path="/login" element={<LoginPage />} />

                    {/* 受保护的路由 */}
                    <Route
                        path="/home"
                        element={
                            <ProtectedRoute>
                                <HomePage />
                            </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/document/:docId"
                        element={
                            <ProtectedRoute>
                                <DocumentViewPage />
                            </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/knowledgeBase/:kbId"
                        element={
                            <ProtectedRoute>
                                <DocumentViewPage />
                            </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/editDocument/:docId"
                        element={
                            <ProtectedRoute>
                                <DocumentEditPage />
                            </ProtectedRoute>
                        }
                    />

                    {/* 默认路由，未匹配时重定向到 /home */}
                    <Route path="*" element={<Navigate to="/home" />} />
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default App;
