import React, { useEffect, useState } from "react";
import { Modal, Input, message, List, Typography, Row, Col, Space, Button } from "antd";
import { personalDocumentSearch, personalKnowledgeSearch } from "../../services/searchService";
import { FileOutlined, BookOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";  // 导入 useNavigate

const { Title } = Typography;

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
    }

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
        }catch (error) {
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
                    bordered
                    dataSource={documentList}
                    renderItem={(item) => (
                        <List.Item>
                            <Button
                                type="link"
                                style={{ width: "100%", textAlign: "left" }}
                                onClick={() => handleDocumentClick(item.doc_id)}  // 点击跳转到文档详情
                            >
                                <Row style={{ width: "100%" }} gutter={16}>
                                    <Col>
                                        <FileOutlined style={{ fontSize: 24, color: "#1890ff" }} />
                                    </Col>
                                    <Col span={18}>
                                        <Space direction="vertical" size={4}>
                                            <Typography.Text strong>
                                                {item.doc_title}
                                            </Typography.Text>
                                            <Typography.Text type="secondary">
                                                {item.kb_name || "无知识库"}
                                            </Typography.Text>
                                        </Space>
                                    </Col>
                                    <Col flex="none">
                                        <Typography.Text type="secondary">
                                            {new Date(item.doc_created_at).toLocaleString()}
                                        </Typography.Text>
                                    </Col>
                                </Row>
                            </Button>
                        </List.Item>
                    )}
                />
            </div>

            {/* 知识库列表 */}
            <div style={{ marginTop: 24 }}>
                <Title level={4}>匹配的知识库</Title>
                <List
                    bordered
                    dataSource={knowledgeList}
                    renderItem={(item) => (
                        <List.Item>
                            <Button
                                type="link"
                                style={{ width: "100%", textAlign: "left" }}
                                onClick={() => handleKnowledgeClick(item.kb_id)}  // 点击跳转到知识库详情
                            >
                                <Row style={{ width: "100%" }} gutter={16}>
                                    <Col>
                                        <BookOutlined style={{ fontSize: 24, color: "#1890ff" }} />
                                    </Col>
                                    <Col span={18}>
                                        <Space direction="vertical" size={4}>
                                            <Typography.Text strong>
                                                {item.kb_name}
                                            </Typography.Text>
                                        </Space>
                                    </Col>
                                    <Col>
                                        <Typography.Text type="secondary">
                                            {new Date(item.kb_created_at).toLocaleString()}  {/* 假设有创建时间字段 */}
                                        </Typography.Text>
                                    </Col>
                                </Row>
                            </Button>
                        </List.Item>
                    )}
                />
            </div>
        </Modal>
    );
};

export default SearchModal;
