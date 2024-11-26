import React, { useEffect, useState } from "react";
import { Button, Comment, Tooltip, message } from "antd";
import ReactMarkdown from "react-markdown";
import { LikeOutlined, LikeFilled, EditOutlined } from "@ant-design/icons";

// 模块化的组件
import CommentSection from "./CommentSection";
import MarkdownView from "./MarkdownView";
import { getDocumentById } from "../services/documentService";

const DocumentViewAndEditContent = ({ docId }) => {
    const [content, setContent] = useState(null); // Markdown 内容
    const [isEditing, setIsEditing] = useState(false); // 编辑状态
    const [likes, setLikes] = useState(0); // 点赞数
    const [liked, setLiked] = useState(false); // 是否已点赞

    const fetchDocumentData = async () => {
        try {
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                setContent(data.doc_content);
            } else {
                message.error("知识库信息拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "知识库信息拉取失败");
        }
    };

    useEffect(() => {
        if(docId){
            fetchDocumentData();
        }
    }, [docId]);

    // 点赞逻辑
    const handleLike = () => {
        setLiked(!liked);
        setLikes(liked ? likes - 1 : likes + 1);
    };

    // 切换编辑状态
    const toggleEdit = () => {
        setIsEditing(!isEditing);
    };

    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
            {/* 内容展示区域 */}
            <div style={{ flex: 3, overflow: "auto", padding: "20px" }}>
                {isEditing ? (
                    <MarkdownView content={content} setContent={setContent} onClose={() => setIsEditing(false)} />
                ) : (
                    <ReactMarkdown>{content}</ReactMarkdown>
                )}
            </div>

            {/* 操作区域 */}
            <div
                style={{
                    flex: 1,
                    display: "flex",
                    justifyContent: "center", // 使用center居中
                    alignItems: "center", // 垂直居中
                    padding: "20px",
                    backgroundColor: "#fff",
                    borderTop: "1px solid #ddd",
                }}
            >
                <Button type="primary" icon={<EditOutlined />} onClick={toggleEdit}>
                    编辑
                </Button>
            </div>

            {/* 点赞按钮居中显示 */}
            <div
                style={{
                    display: "flex",
                    justifyContent: "center", // 水平居中
                    alignItems: "center", // 垂直居中
                    padding: "20px",
                    backgroundColor: "#fff",
                    borderTop: "1px solid #ddd",
                    height: "80px", // 控制按钮区域的高度
                }}
            >
                <Button
                    type="default"
                    icon={liked ? <LikeFilled /> : <LikeOutlined />}
                    onClick={handleLike}
                    style={{ fontSize: "24px", padding: "10px 20px" }} // 放大按钮和文字
                >
                    {likes} 点赞
                </Button>
            </div>

            {/* 评论区域 */}
            <div style={{ flex: 4, overflow: "auto", padding: "20px"}}>
                <CommentSection />
            </div>
        </div>
    );
};

export default DocumentViewAndEditContent;
