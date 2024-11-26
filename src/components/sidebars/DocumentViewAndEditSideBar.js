import React, { useState, useEffect } from "react";
import { Button, List, Layout, message, Modal, Input, Row, Col, Dropdown, Menu } from "antd";
import { ArrowLeftOutlined, EllipsisOutlined } from "@ant-design/icons";
import { createDocument, deleteDocument, getDocumentListByKbId } from "../../services/documentService";
import { useNavigate } from "react-router-dom";
import { getKnowledgeBaseDetail } from "../../services/knowledgeService";

const { Sider } = Layout;

const DocumentViewAndEditSideBar = ({ kbId, docId }) => {
    const [documentList, setDocumentList] = useState([]); // 存储文档列表
    const [loading, setLoading] = useState(false); // 加载状态
    const [selectedDocId, setSelectedDocId] = useState(null); // 记录选中的文档 ID
    const [kbName, setKbName] = useState(null);
    const [isModalVisible, setIsModalVisible] = useState(false); // 控制模态框显示
    const [newDocTitle, setNewDocTitle] = useState(""); // 新建文档标题
    const navigate = useNavigate();

    // 获取文档列表和知识库信息
    const fetchDocumentList = async () => {
        try {
            // 获取知识库信息
            try {
                const response = await getKnowledgeBaseDetail(kbId);
                if (response.status === 200) {
                    const data = response.data;
                    setKbName(data["kb_name"]);
                } else {
                    message.error("知识库信息拉取失败");
                }
            } catch (error) {
                message.error(error.response?.data?.error || "知识库信息拉取失败");
            }

            // 获取文档列表
            try {
                const response = await getDocumentListByKbId(kbId);
                if (response.status === 200) {
                    const docList = response.data['doc_list'];
                    if (docList != null) {
                        setDocumentList(docList);
                    }else{
                        setDocumentList([]);
                    }
                } else {
                    message.error("获取文档列表失败");
                }
            } catch (error) {
                message.error("获取文档列表失败，请稍后重试");
            } finally {
                setLoading(false);
            }
        } catch (error) {
            message.error(error.response?.data?.error || "系统错误，文档拉取失败");
        }
    };

    const onBackToHome = () => {
        navigate("/home");
    };

    // 组件加载时获取文档列表
    useEffect(() => {
        if (kbId && docId) {
            fetchDocumentList();
        }
    }, [kbId, docId]);

    // 点击文档项的处理函数
    const onDocumentClick = (docId) => {
        setSelectedDocId(docId); // 更新选中状态
        navigate(`/document/${docId}`); // 跳转到文档页面
    };

    // 新建文档按钮点击事件
    const showModal = () => {
        setIsModalVisible(true); // 显示模态框
    };

    // 关闭模态框
    const handleCancel = () => {
        setIsModalVisible(false);
        setNewDocTitle(""); // 清空文档标题
    };

    // 提交新建文档
    const handleCreateDocument = async () => {
        if (!newDocTitle) {
            message.error("文档标题不能为空");
            return;
        }

        // 这里可以执行创建文档的 API 请求
        try {
            // 假设有一个API接口可以创建文档
            const response = await createDocument(kbId, newDocTitle);
            if (response.status === 200) {
                const newDocId = response.data['doc_id'];
                message.success("文档创建成功");
                fetchDocumentList(); // 刷新文档列表
                handleCancel(); // 关闭模态框
                navigate(`/document/${newDocId}`);
            } else {
                message.error("文档创建失败");
            }
        } catch (error) {
            message.error("创建文档时发生错误，请稍后重试");
        }
    };

    // 删除文档
    const handleDeleteDocument = async (docId) => {
        try {
            const response = await deleteDocument(docId);
            if (response.status === 200) {
                message.success("文档删除成功");
                fetchDocumentList(); // 刷新文档列表
            } else {
                message.error("文档删除失败");
            }
        } catch (error) {
            message.error("删除文档时发生错误，请稍后重试");
        }
    };

    // 渲染文档操作菜单
    const menu = (docId) => (
        <Menu>
            <Menu.Item onClick={() => handleDeleteDocument(docId)}>删除文档</Menu.Item>
        </Menu>
    );

    return (
        <Layout>
            <Sider width={300} theme="light" style={{ padding: "16px" }}>
                <div style={{ display: "flex", alignItems: "center", marginBottom: "20px" }}>
                    <Button
                        icon={<ArrowLeftOutlined />}
                        onClick={onBackToHome}
                        style={{ marginRight: "10px" }}
                    />
                    <span>{kbName}</span>
                </div>

                {/* 新建文档按钮和文档列表 */}
                <Row justify="space-between" align="middle" style={{ marginBottom: "20px" }}>
                    <Col>
                        <h3>文档列表</h3>
                    </Col>
                    <Col>
                        <Button
                            type="primary"
                            onClick={showModal}
                        >
                            新建文档
                        </Button>
                    </Col>
                </Row>

                <div>
                    <List
                        loading={loading}
                        dataSource={documentList}
                        renderItem={(doc) => (
                            <List.Item
                                key={doc.doc_id}
                                onClick={() => onDocumentClick(doc.doc_id)} // 点击事件
                                style={{
                                    cursor: "pointer",
                                    backgroundColor: doc.doc_id === selectedDocId ? "#e6f7ff" : "white", // 高亮样式
                                    padding: "10px 16px",
                                    borderRadius: "4px",
                                    display: "flex",
                                    justifyContent: "space-between",
                                    alignItems: "center"
                                }}
                            >
                                <span>{doc.doc_title}</span>
                                <Dropdown overlay={menu(doc.doc_id)} trigger={['click']}>
                                    <Button icon={<EllipsisOutlined />} />
                                </Dropdown>
                            </List.Item>
                        )}
                    />
                </div>
            </Sider>

            {/* 新建文档模态框 */}
            <Modal
                title="新建文档"
                visible={isModalVisible}
                onCancel={handleCancel}
                onOk={handleCreateDocument}
            >
                <Input
                    placeholder="输入文档标题"
                    value={newDocTitle}
                    onChange={(e) => setNewDocTitle(e.target.value)}
                />
            </Modal>
        </Layout>
    );
};

export default DocumentViewAndEditSideBar;
