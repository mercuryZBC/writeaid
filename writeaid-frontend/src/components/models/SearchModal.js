import React, { useEffect, useState } from "react";
import { Modal, Input, message, List, Typography, Row, Col, Space, Button, Avatar } from "antd";
import { personalDocumentSearch, personalKnowledgeSearch } from "../../services/searchService";
import { FileOutlined, BookOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";  // 导入 useNavigate

const { Title, Text } = Typography;

const SearchModal = ({ visible, onClose }) => {
    const [searchText, setSearchText] = useState("");
    const [documentList, setDocumentList] = useState([]);
    const [knowledgeList, setKnowledgeList] = useState([]);
    const navigate = useNavigate();  // 初始化 navigate

    useEffect(() => {
        if (searchText !== "") {
            handleSearch();
        }
    }, [searchText]);

    const handleCloseSearch = () => {
        setSearchText("");
        onClose();
    };

    const handleSearch = async () => {
        try {
            const knowledgeBasesResponse = await personalDocumentSearch(searchText);
            if (knowledgeBasesResponse.status === 200) {
                setDocumentList(knowledgeBasesResponse.data.documents || []);
            }
        } catch (error) {
            message.error(error.response?.data?.error || "搜索失败，请稍后再试");
        }
        try {
            const documentsResponse = await personalKnowledgeSearch(searchText);
            if (documentsResponse.status === 200) {
                setKnowledgeList(documentsResponse.data.knowledgeBases || []);
            }
        } catch (error) {
            message.error(error.response?.data?.error || "搜索失败，请稍后再试");
        }
    };

    // 跳转到文档详情页面
    const handleDocumentClick = (docId) => {
        navigate(`/document/${docId}`);  // 跳转到文档详情页
    };

    // 跳转到知识库详情页面
    const handleKnowledgeClick = (kbId) => {
        navigate(`/knowledgeBase/${kbId}`);  // 跳转到知识库详情页
    };

    return (
        <Modal
            title="搜索"
            open={visible}
            onCancel={handleCloseSearch}
            onOk={handleSearch}
            width={800}
        >
            <Input
                placeholder="请输入搜索关键词"
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
                style={{ marginBottom: 16 }}
            />

            {/* 文档列表 */}
            <div>
                <Title level={4}>匹配的文档</Title>
                <List
                    itemLayout="horizontal"
                    dataSource={documentList}
                    renderItem={(item) => (
                        <List.Item
                            onClick={() => handleDocumentClick(item.doc_id)}
                            style={{
                                cursor: "pointer",
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                            }}
                        >
                            <div style={{ display: "flex", alignItems: "center" }}>
                                <Avatar icon={<FileOutlined style={{ fontSize: 20, color: "#1890ff" }} />} />
                                <div style={{ marginLeft: 12 }}>
                                    <Text strong>{item.doc_title}</Text>
                                    <br />
                                    <Text type="secondary">{item.kb_name || "无知识库"}</Text>
                                </div>
                            </div>
                            <Text type="secondary">
                                {new Date(item.doc_created_at).toLocaleString()}
                            </Text>
                        </List.Item>
                    )}
                />
            </div>

            {/* 知识库列表 */}
            <div style={{ marginTop: 24 }}>
                <Title level={4}>匹配的知识库</Title>
                <List
                    itemLayout="horizontal"
                    dataSource={knowledgeList}
                    renderItem={(item) => (
                        <List.Item
                            onClick={() => handleKnowledgeClick(item.kb_id)}
                            style={{
                                cursor: "pointer",
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                            }}
                        >
                            <div style={{ display: "flex", alignItems: "center" }}>
                                <Avatar icon={<BookOutlined style={{ fontSize: 20, color: "#1890ff" }} />} />
                                <div style={{ marginLeft: 12 }}>
                                    <Text strong>{item.kb_name}</Text>
                                </div>
                            </div>
                            <Text type="secondary">
                                {new Date(item.kb_created_at).toLocaleString()}
                            </Text>
                        </List.Item>
                    )}
                />
            </div>
        </Modal>
    );
};

export default SearchModal;
