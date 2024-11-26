import React, { useState } from "react";
import { Layout, Input, Typography } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import KnowledgeList from "./lists/KnowledgeList";
import SearchModal from "./models/SearchModal";

const { Sider } = Layout;
const { Title } = Typography;

export const SideBar = () => {
    const navigate = useNavigate();
    const [searchModalVisible, setSearchModalVisible] = useState(false);

    const handleNavigateHome = () => {
        navigate("/home");
    };
    return (
        <Sider
            width={300}
            theme="light"
            style={{
                borderRight: "1px solid #f0f0f0",
                backgroundColor: "#fff",
                borderRadius: "8px",
                boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
                height: "100vh",
            }}
        >
            {/* 顶部功能区 */}
            <div
                style={{
                    display: "flex",
                    alignItems: "center",
                    padding: "16px",
                    gap: "8px",
                    backgroundColor: "#fafafa",
                    borderTopLeftRadius: "8px",
                    borderTopRightRadius: "8px",
                }}
            >
                <Title
                    level={3}
                    onClick={handleNavigateHome}
                    style={{
                        cursor: "pointer",
                        color: "#1890ff",
                        margin: 0,
                        fontSize: "18px",
                    }}
                >
                    WriteAid
                </Title>
                <Input.Search
                    placeholder="搜索知识库/文档"
                    enterButton={<SearchOutlined />}
                    onSearch={() => setSearchModalVisible(true)}
                    style={{
                        flex: 1,
                        borderRadius: "4px",
                    }}
                />
            </div>

            {/* 知识库列表 */}
            <KnowledgeList />

            {/* 搜索框 */}
            <SearchModal
                visible={searchModalVisible}
                onClose={() => setSearchModalVisible(false)}
            />
        </Sider>
    );
};
