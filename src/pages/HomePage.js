import React, {createContext, useEffect, useState} from "react";
import {Layout, Input, Button, Modal, Form, message} from "antd";
import {Overview} from "../components/Overview"
import {SideBar} from "../components/SideBar";
import {MarkdownEditor} from "../components/MarkdownEditor"
import {Router, useNavigate} from "react-router-dom";
import * as userService from "../services/userService";
import {ReloadProvider} from "../context/ReloadContext";



const { Sider, Content } = Layout;

const HomePage = () => {
    // 创建一个 Context 用于跨组件共享状态
    const knowledgeListReloadContext = createContext();
    const [selectedNode, setSelectedNode] = useState(null);
    const [searchModalVisible, setSearchModalVisible] = useState(false);
    const [createModalVisible, setCreateModalVisible] = useState(false);
    const [treeData, setTreeData] = useState([]);
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
                if(error.response.status && error.response.status === 401) {
                    navigate('/login')
                }else{
                    message.error(error.response?.data?.error || "获取用户数据失败");
                }
            }
        };

        // 调用异步函数
        fetchUserInfo();
    }, []); // 依赖数组为空，表示只在组件挂载时运行

    // 点击树节点时触发
    const onSelect = (keys, event) => {
        setSelectedNode(event.node.title);
    };

    // 打开搜索框
    const openSearchModal = () => {
        setSearchModalVisible(true);
    };

    // 打开创建知识库/文档弹窗
    const openCreateModal = () => {
        setCreateModalVisible(true);
    };

    return (
        <Layout style={{ height: "100vh" }}>
           <ReloadProvider>
               {/* 左侧区域 */}
               <Sider width={300} theme="light" style={{ borderRight: "1px solid #f0f0f0" }}>
                   <SideBar
                       treeData={treeData}
                       setTreeData={setTreeData}
                       onSelect={onSelect}
                       openSearchModal={openSearchModal}/>
               </Sider>

               {/* 右侧区域 */}
               <Content style={{ padding: "24px", background: "#fff" }}>
                   <Overview/>
               </Content>

               {/* 搜索弹窗 */}
               <Modal
                   title="搜索"
                   open={searchModalVisible}
                   onCancel={() => setSearchModalVisible(false)}
                   footer={null}
               >
                   <Input.Search placeholder="输入关键字进行搜索" enterButton />
               </Modal>
           </ReloadProvider>
        </Layout>
    );
};

export default HomePage;
