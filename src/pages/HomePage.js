import React, {useEffect} from "react";
import {Layout,message} from "antd";
import {OverviewContent} from "../components/OverviewContent"
import {HomeSideBar} from "../components/sidebars/HomeSideBar";
import {useNavigate} from "react-router-dom";
import * as userService from "../services/userService";
import {ReloadProvider} from "../context/ReloadContext";
import {delToken} from "../util/jwt";



const { Sider, Content } = Layout;

const HomePage = () => {
    const navigate = useNavigate();
    useEffect(() => {
        // 确定当前用户是否已经登录，否则跳转到登录界面
        const fetchUserInfo = async () => {
            try {
                const response = await userService.getUserInfo();
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
                   <HomeSideBar/>
               </Sider>

               {/* 右侧区域 */}
               <Content style={{ padding: "24px", background: "#fff" }}>
                   <OverviewContent/>
               </Content>
           </ReloadProvider>
        </Layout>
    );
};

export default HomePage;
