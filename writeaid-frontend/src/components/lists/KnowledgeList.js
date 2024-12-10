import React, { useState, useEffect } from "react";
import { Collapse, Button, Popover, Menu, message } from "antd";
import { MoreOutlined } from "@ant-design/icons";
import { deleteKnowledgeBase, getKnowledgeBaseList } from "../../services/knowledgeService";
import DocumentListDrawer from "./DocumentListDrawer";
import { useNavigate } from "react-router-dom";
import { useReloadContext } from "../../contexts/ReloadContext";
import NewKnowledgeBaseModal from "../models/NewKnowledgeBaseModal";

const { Panel } = Collapse;

const KnowledgeList = () => {
    const [knowledgeList, setKnowledgeList] = useState([]);
    const { kbListReloadTrigger, setkbListReloadTrigger } = useReloadContext();
    const [isNewKBModelVisible, setIsNewKBModelVisible] = useState(false);
    const [activeKey, setActiveKey] = useState(["1"]); // 控制展开的面板
    const navigate = useNavigate();

    // 获取知识库列表
    const fetchKnowledgeBaseList = async () => {
        try {
            const response = await getKnowledgeBaseList();
            if (response.status === 200) {
                if (response.data.knowledge_bases) {
                    setKnowledgeList(response.data.knowledge_bases);
                } else {
                    setKnowledgeList([]);
                }
            } else {
                message.error("知识库拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "知识库拉取失败");
        }
    };

    useEffect(() => {
        fetchKnowledgeBaseList();
    }, [kbListReloadTrigger]);

    const handleNewKBModelOpen = () => {
        setIsNewKBModelVisible(true);
    };

    const handleNewKBClose = () => {
        setIsNewKBModelVisible(false);
    };

    const KBListButtonItem = ({ kb }) => {
        const [popoverVisible, setPopoverVisible] = useState(false);
        const [drawerVisible, setDrawerVisible] = useState(false);

        const onSelectKnowledgeBase = () => {
            navigate(`/knowledgeBase/${kb.kb_id}`);
        };

        const handleDelete = async (kbId) => {
            try {
                const response = await deleteKnowledgeBase(kbId);
                if (response.status === 200) {
                    message.success("知识库删除成功");
                    fetchKnowledgeBaseList();
                } else {
                    message.error("知识库删除失败");
                }
            } catch (error) {
                message.error(error.response?.data?.error || "系统错误，知识库删除失败");
            }
        };

        const renderMenu = (kbId) => (
            <Menu
                onClick={({ key }) => {
                    if (key === "1") {
                        handleDelete(kbId);
                    } else if (key === "2") {
                        message.info("设置权限");
                    }
                    setPopoverVisible(false);
                }}
            >
                <Menu.Item key="1">删除</Menu.Item>
                <Menu.Item key="2">设置权限</Menu.Item>
            </Menu>
        );

        return (
            <div
                key={kb.kb_id}
                style={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    marginBottom: "8px",
                }}
            >
                <Button
                    style={{ flex: 1, textAlign: "center" }}
                    onClick={() => onSelectKnowledgeBase()}
                >
                    {kb.kb_name}
                </Button>
                <Popover
                    content={renderMenu(kb.kb_id)}
                    trigger="click"
                    open={popoverVisible}
                    onOpenChange={setPopoverVisible}
                >
                    <Button
                        icon={<MoreOutlined />}
                        style={{ marginLeft: 8, border: "none", backgroundColor: "transparent" }}
                    />
                </Popover>
                <DocumentListDrawer
                    kbId={kb.kb_id}
                    kbName={kb.kb_name}
                    visible={drawerVisible}
                    onClose={() => setDrawerVisible(false)}
                />
            </div>
        );
    };

    return (
        <>
            <Collapse
                activeKey={activeKey}  // 使用 activeKey 来控制面板展开
                onChange={(key) => setActiveKey(key)} // 更新 activeKey
                style={{ border: "1px solid #f0f0f0", borderRadius: "4px" }}
            >
                <Panel
                    header={
                        <div
                            style={{
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                            }}
                        >
                            <span style={{ fontSize: "18px", color: "#1890ff" }}>知识库</span>
                            <Button type="primary" onClick={handleNewKBModelOpen}>
                                创建知识库
                            </Button>
                        </div>
                    }
                    key="1"
                >
                    {knowledgeList.map((kb) => (
                        <KBListButtonItem key={kb.kb_id} kb={kb} />
                    ))}
                </Panel>
            </Collapse>
            <NewKnowledgeBaseModal
                visible={isNewKBModelVisible}
                onCancel={handleNewKBClose}
            />
        </>
    );
};

export default KnowledgeList;
