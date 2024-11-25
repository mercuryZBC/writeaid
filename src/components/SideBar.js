import React, { useEffect, useState } from "react";
import { Button, Input, Layout, Popover, Menu, Typography, message, List, Drawer } from "antd";
import { SearchOutlined, MoreOutlined } from "@ant-design/icons";
import { useNavigate } from 'react-router-dom';
import {deleteKnowledgeBase, getKnowledgeBaseList} from "../services/knowledgeService";
import {useReloadContext} from "../context/ReloadContext";

const { Sider } = Layout;
const { Title } = Typography;

export const SideBar = ({ treeData, setTreeData, onSelect, openSearchModal }) => {
    const navigate = useNavigate();
    const [visible, setVisible] = useState(false);  // 控制 Drawer 显示与隐藏
    const [selectedKB, setSelectedKB] = useState(null);  // 保存当前点击的知识库
    const [documents, setDocuments] = useState([]);  // 保存文档列表
    const {kbListReloadTrigger, setkbListReloadTrigger} = useReloadContext();
    // 获取知识库列表
    const getKnowledgeBaseListReq = async () => {
        try {
            const response = await getKnowledgeBaseList();
            if (response.status === 200) {
                const kbList = response.data['knowledge_bases'];
                const tmpTreeData = [];
                kbList.forEach((kb) => {
                    tmpTreeData.push({
                        title: <TreeNodeButton title={kb.kb_name} kbId={kb.kb_id} />,
                        key: kb.kb_id.toString(),
                    });
                });
                setTreeData(tmpTreeData);
            } else {
                message.error('知识库拉取失败');
            }
        } catch (error) {
            message.error(error.response?.data?.error || "知识库拉取失败");
        }
    };
    useEffect(() => {
        getKnowledgeBaseListReq();
    }, [kbListReloadTrigger]);

    // TreeNodeButton 组件：用于渲染列表项中的按钮和菜单
    const TreeNodeButton = ({ title, kbId }) => {
        const [popoverVisible, setPopoverVisible] = useState(false);  // 控制菜单的可见性
        const MyMenu = ({onClose}) => {
            const handleMenuClick = async (e) => {
                // 根据 e.key 来判断点击的是哪个菜单项
                if (e.key === "1") {
                    try {
                        console.log(title,kbId);
                        const response = await deleteKnowledgeBase(kbId);

                        if (response.status === 200) {
                            message.success('知识库删除成功');
                            await getKnowledgeBaseListReq();
                        } else {
                            message.error('知识库删除失败');
                        }
                    } catch (error) {
                        message.error(error.response?.data?.error || '系统错误，知识库删除失败');
                    }
                    setPopoverVisible(false);
                } else if (e.key === "2") {
                    message.success("设置权限操作");
                    setPopoverVisible(false);
                }
            }
            return (
                <Menu onClick={handleMenuClick}>
                    <Menu.Item key="1">删除</Menu.Item>
                    <Menu.Item key="2">设置权限</Menu.Item>
                </Menu>
            )
        };

        const openDocumentsDrawer = async () => {
            // 在点击按钮时打开文档列表的 Drawer，并拉取该知识库下的文档
            setSelectedKB(kbId);
            // 模拟从服务器获取文档列表
            const fetchedDocuments = [
                { id: 1, title: '文档 1' },
                { id: 2, title: '文档 2' },
                { id: 3, title: '文档 3' },
            ];
            setDocuments(fetchedDocuments);
            setVisible(true);
        };

        return (
            <div style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                width: "100%",
            }}>
                <Button
                    style={{
                        width: "100%",
                        textAlign: "left",
                        border: "none",
                        padding: "8px 16px",
                        fontSize: "14px",
                        fontWeight: "500",
                        color: "#1890ff",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                        overflow: "hidden"
                    }}
                    onClick={openDocumentsDrawer}  // 点击按钮时打开 Drawer
                >
                    {title}
                </Button>
                <Popover content={<MyMenu onClose={() => setVisible(false)}/>} trigger="click" open={popoverVisible} onOpenChange={setPopoverVisible}>
                    <Button
                        icon={<MoreOutlined />}
                        style={{
                            marginLeft: 8,
                            border: "none",
                            backgroundColor: "transparent",
                            color: "#1890ff",
                            padding: "4px",
                            fontSize: "16px",
                            cursor: "pointer",
                        }}
                    />
                </Popover>
            </div>
        );
    };

    // 点击 yuquepp 跳转到主页
    const handleWriteAidClick = () => {
        navigate('/home');
    };

    // 选择文档并执行操作
    const handleDocumentSelect = (documentId) => {
        console.log('Selected document:', documentId);
        // 这里可以进行相应的文档打开操作
        message.success(`文档 ${documentId} 已选中`);
        setVisible(false);  // 关闭 Drawer
    };

    return (
        <Sider width={300} theme="light" style={{
            borderRight: "1px solid #f0f0f0",
            backgroundColor: "#fff",
            borderRadius: "8px",
            boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
            height: "100vh",
        }}>
            {/* 顶部功能按钮区域 */}
            <div style={{
                display: "flex",
                alignItems: "center",
                padding: "16px",
                paddingBottom: "50px",
                gap: "8px",
                backgroundColor: "#fafafa",
                borderTopLeftRadius: "8px",
                borderTopRightRadius: "8px"
            }}>
                <Title
                    level={3}
                    onClick={handleWriteAidClick}
                    style={{
                        cursor: 'pointer',
                        color: "#1890ff",
                        margin: 0,
                        fontSize: "18px"
                    }}
                >
                    WriteAid
                </Title>
                <Input.Search
                    placeholder="搜索知识库/文档"
                    enterButton={<SearchOutlined />}
                    onSearch={openSearchModal}
                    style={{
                        flex: 1,
                        borderRadius: "4px",
                    }}
                />
            </div>

            {/* 知识库列表 */}
            <div style={{
                padding: "16px",
                paddingTop: "8px",
                backgroundColor: "#fafafa",
                borderRadius: "8px",
                height: "100%",
            }}>
                <List
                    dataSource={treeData}
                    renderItem={(item) => (
                        <List.Item key={item.key} style={{ padding: 0 }}>
                            {item.title}
                        </List.Item>
                    )}
                    style={{
                        fontFamily: "'Roboto', sans-serif",
                        fontSize: "14px",
                        color: "#333",
                        backgroundColor: "#fafafa",
                    }}
                />
            </div>

            {/* Drawer 弹窗：选择文档 */}
            <Drawer
                title="选择文档"
                width={400}
                onClose={() => setVisible(false)}
                visible={visible}
                bodyStyle={{ paddingBottom: 80 }}
            >
                <List
                    dataSource={documents}
                    renderItem={(doc) => (
                        <List.Item
                            key={doc.id}
                            onClick={() => handleDocumentSelect(doc.id)}
                            style={{ cursor: "pointer" }}
                        >
                            <Button style={{ width: "100%", textAlign: "left" }}>
                                {doc.title}
                            </Button>
                        </List.Item>
                    )}
                />
            </Drawer>
        </Sider>
    );
};
