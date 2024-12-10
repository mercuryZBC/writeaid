import React, { useState } from "react";
import { Layout, Button, Typography } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import KnowledgeList from "../lists/KnowledgeList";
import SearchModal from "../models/SearchModal";

const { Sider } = Layout;
const { Title } = Typography;

export const HomeSideBar = () => {
    const navigate = useNavigate();
    const [searchModalVisible, setSearchModalVisible] = useState(false);

    const handleNavigateHome = () => {
        navigate("/home");
    };

    const handleOpenSearch = () => {
        setSearchModalVisible(true);
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
                <Button
                    type="primary"
                    icon={<SearchOutlined />}
                    onClick={handleOpenSearch}
                    style={{
                        marginLeft: "auto", // 将按钮推到右侧
                        borderRadius: "4px",
                    }}
                >
                    搜索
                </Button>
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
