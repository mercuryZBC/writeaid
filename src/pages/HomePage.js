import React, {createContext, useEffect, useState} from "react";
import {Layout, Input, Button, Modal, Form, message} from "antd";
import {Overview} from "../components/Overview"
import {SideBar} from "../components/SideBar";
import {MarkdownEditor} from "../components/MarkdownEditor"
import {Router, useNavigate} from "react-router-dom";
import * as userService from "../services/userService";
import {ReloadProvider} from "../context/ReloadContext";
import {delToken} from "../util/jwt";



const { Sider, Content } = Layout;

const HomePage = () => {
    const navigate = useNavigate();
    const [userInfo, setUserInfo] = useState({
        userid: null,
        nickname: "",
        email:""
    });

    useEffect(() => {
        // 创建异步函数
        const fetchUserInfo = async () => {
            try {
                const response = await userService.getUserInfo();
                const data = response.data;
                if (response.status === 200) {
                    setUserInfo(data);
                }
            } catch (error) {
                if(error.response && error.response.status === 401) {
                    delToken();
                    navigate('/login')
                }else{
                    message.error(error.response?.data?.error || "获取用户数据失败");
                }
            }
        };

        // 调用异步函数
        fetchUserInfo();
    }, []); // 依赖数组为空，表示只在组件挂载时运行

    return (
        <Layout style={{ height: "100vh" }}>
           <ReloadProvider>
               {/* 左侧区域 */}
               <Sider width={300} theme="light" style={{ borderRight: "1px solid #f0f0f0" }}>
                   <SideBar/>
               </Sider>

               {/* 右侧区域 */}
               <Content style={{ padding: "24px", background: "#fff" }}>
                   <Overview/>
               </Content>
           </ReloadProvider>
        </Layout>
    );
};

export default HomePage;
