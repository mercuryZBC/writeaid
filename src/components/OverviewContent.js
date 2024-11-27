import { Avatar, List, message, Tabs, Typography } from "antd";
import React, { useEffect, useState } from "react";
import NavigatorBar from "./NavigatorBar";
import { getRecentViewDocuments } from "../services/documentService";
import { FileOutlined } from "@ant-design/icons"; // 使用文档图标

const { TabPane } = Tabs;
const { Text } = Typography;

export const OverviewContent = () => {
    // 本地存储状态
    const [recentViewDocList, setRecentViewDocList] = useState([]);
    const [recentEditDocList, setRecentEditDocList] = useState([]);
    const [recentCommentDocList, setRecentCommentDocList] = useState([]);

    // 将时间戳转换为可读时间
    const convertTimestampToTime = (timestamp) => {
        const date = new Date(timestamp * 1000); // 转换为毫秒
        return date.toLocaleString(); // 格式化时间为本地时间字符串
    };

    // 提取数据并显示
    const parseData = (recent_docs) => {
        const recentDocs = recent_docs.map((doc) => ({
            key: doc.doc_id,
            title: doc.doc_title,
            knowledgeBase: doc.kb_name,
            time: convertTimestampToTime(doc.timestamp),
            link: `/document/${doc.doc_id}`, // 假设文档跳转链接
        }));
        return recentDocs;
    };

    const fetchRecentViewDocuments = async () => {
        try {
            const response = await getRecentViewDocuments(0, 10);
            if (response.status === 200) {
                const recent_docs = parseData(response.data.recent_docs);
                setRecentViewDocList(recent_docs);
            } else {
                message.error("系统错误，最近浏览文档获取失败");
            }
        } catch (error) {
            message.error(error?.response?.data?.error || "系统错误，最近浏览文档获取失败");
        }
    };

    useEffect(() => {
        fetchRecentViewDocuments();
    }, []);

    // 渲染列表项
    const renderListItems = (dataSource) => (
        <List
            itemLayout="horizontal"
            dataSource={dataSource}
            renderItem={(item) => (
                <List.Item
                    onClick={() => (window.location.href = item.link)}
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
                            <Text strong>{item.title}</Text>
                            <br />
                            <Text type="secondary">{item.knowledgeBase}</Text>
                        </div>
                    </div>
                    <Text type="secondary">{item.time}</Text>
                </List.Item>
            )}
        />
    );

    return (
        <div>
            <NavigatorBar />
            <Tabs defaultActiveKey="viewed" centered>
                <TabPane tab="最近浏览" key="viewed">
                    {renderListItems(recentViewDocList)}
                </TabPane>
                <TabPane tab="最近编辑" key="edited">
                    {renderListItems(recentEditDocList)}
                </TabPane>
                <TabPane tab="最近评论" key="commented">
                    {renderListItems(recentCommentDocList)}
                </TabPane>
            </Tabs>
        </div>
    );
};
