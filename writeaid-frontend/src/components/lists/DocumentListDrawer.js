import React, { useEffect, useState } from "react";
import { Drawer, List, Button, message, Popover, Menu, Modal, Input } from "antd";
import { useNavigate } from "react-router-dom";
import { getDocumentListByKbId, deleteDocument, createDocument } from "../../services/documentService"; // 假设 createDocument 是新建文档的 API

const DocumentListDrawer = ({ kbId, kbName, visible, onClose }) => {
    const [documents, setDocuments] = useState([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false); // 控制新建文档弹窗的显示
    const [newDocTitle, setNewDocTitle] = useState(""); // 新文档的标题
    const navigate = useNavigate();

    // 拉取文档列表
    const fetchDocuments = async () => {
        if (!kbId) return;
        setLoading(true);
        try {
            const response = await getDocumentListByKbId(kbId);
            if (response.status === 200) {
                const docList = response.data['doc_list'];
                if (docList != null) {
                    setDocuments(docList);
                }else{
                    setDocuments([]);
                }
            } else {
                message.error("文档拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "系统错误，文档拉取失败");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (visible) {
            fetchDocuments();
        }
    }, [kbId, visible]);

    // 点击文档后跳转到预览界面
    const handleDocumentClick = (docId) => {
        navigate(`/document/${docId}`); // 假设文档预览路径为 /document/:id
    };

    // 删除文档
    const handleDelete = async (docId) => {
        try {
            const response = await deleteDocument(docId); // 假设这是删除文档的 API
            if (response.status === 200) {
                message.success("文档删除成功");
                await fetchDocuments(); // 重新拉取文档列表
            } else {
                message.error("文档删除失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "系统错误，删除文档失败");
        }
    };

    // 渲染菜单，包含删除操作
    const renderMenu = (docId) => (
        <Menu
            onClick={({ key }) => {
                if (key === "delete") {
                    handleDelete(docId);
                }
            }}
        >
            <Menu.Item key="delete">删除</Menu.Item>
        </Menu>
    );

    // 显示新建文档的弹窗
    const showCreateDocumentModal = () => {
        setIsModalVisible(true);
    };

    // 关闭新建文档的弹窗
    const handleCancel = () => {
        setIsModalVisible(false);
    };

    // 提交新建文档
    const handleCreateDocument = async () => {
        if (!newDocTitle) {
            message.error("文档标题不能为空");
            return;
        }

        try {
            const response = await createDocument(kbId, newDocTitle); // 假设这是新建文档的 API
            if (response.status === 200) {
                message.success("文档创建成功");
                fetchDocuments(); // 重新拉取文档列表
                setIsModalVisible(false); // 关闭弹窗
            } else {
                message.error("文档创建失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "系统错误，创建文档失败");
        }
    };

    return (
        <Drawer
            title={
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <span>{`知识库 ${kbName} 下的文档`}</span>
                    <Button
                        type="primary"
                        size="small"
                        onClick={showCreateDocumentModal}
                    >
                        新建文档
                    </Button>
                </div>
            }
            placement="right"
            width={400}
            onClose={onClose}
            open={visible}
            destroyOnClose
        >
            <List
                loading={loading}
                dataSource={documents}
                renderItem={(item) => (
                    <List.Item
                        actions={[
                            <Button
                                type="link"
                                onClick={() => handleDocumentClick(item.doc_id)}
                            >
                                查看
                            </Button>,
                            <Popover
                                content={renderMenu(item.doc_id)}
                                trigger="click"
                            >
                                <Button type="link">选项</Button>
                            </Popover>,
                        ]}
                    >
                        {item.doc_title}
                    </List.Item>
                )}
            />

            {/* 新建文档弹窗 */}
            <Modal
                title="新建文档"
                visible={isModalVisible}
                onCancel={handleCancel}
                onOk={handleCreateDocument}
                okText="创建"
                cancelText="取消"
            >
                <Input
                    placeholder="请输入文档标题"
                    value={newDocTitle}
                    onChange={(e) => setNewDocTitle(e.target.value)}
                />
            </Modal>
        </Drawer>
    );
};

export default DocumentListDrawer;
