import React, { useEffect, useState } from "react";
import { Button, message } from "antd";
import ReactMarkdown from "react-markdown";
import { LikeOutlined, LikeFilled, EditOutlined } from "@ant-design/icons";

// 模块化的组件
import CommentSection from "./CommentSection";
import {getDocumentById, getDocumentContentHash} from "../services/documentService";
import {useNavigate} from "react-router-dom";
import {DocumentProvider, useDocumentContext} from "../contexts/DocumentContext";
import {addDocumentToDB, clearDocumentsFromDB, getDocumentByDocIdFromDB} from "../localStorage/documentIDB";
import {calculateSHA256} from "../util/util";

const DocumentViewContent = ({ docId }) => {
    const [content, setContent] = useState(null); // Markdown 内容
    const [likes, setLikes] = useState(0); // 点赞数
    const [liked, setLiked] = useState(false); // 是否已点赞
    const navigate = useNavigate();
    // 获取文档内容
    const fetchDocument = async () =>{
        try {
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                setContent(data.doc_content);
                const document = {
                    doc_id: data.doc_id,
                    doc_title: data.doc_title,
                    doc_content: data.doc_content,
                    doc_content_hash: await calculateSHA256(data.doc_content)
                };
                await addDocumentToDB(document);
            } else {
                message.error("文档信息拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "文档信息拉取失败");
        }
    }
    const DocumentLoader = async () => {
        const documentCache = await getDocumentByDocIdFromDB(docId);
        if(documentCache == null) {
            await fetchDocument();
        }else {
            try{
                const response = await getDocumentContentHash(docId);
                if (response.status === 200) {
                    if(response.data["doc_content_hash"] === documentCache.doc_content_hash){
                        setContent(documentCache.doc_content);
                    }else{
                        await fetchDocument(documentCache.doc_id);
                    }
                } else {
                    message.error("文档信息拉取失败");
                }
            }catch (error) {
                message.error(error.response?.data?.error || "文档信息拉取失败");
            }
        }
    };

    useEffect(() => {
        if (docId) {
            DocumentLoader();
        }
    }, [docId]);

    // 点赞逻辑
    const handleLike = () => {
        setLiked(!liked);
        setLikes(liked ? likes - 1 : likes + 1);
    };

    // 切换编辑状态
    const toggleEdit = () => {
        // setIsEditing(!isEditing);
        navigate(`/editDocument/${docId}`);
    };

    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
            {/* 文档内容展示或编辑区域 */}
            <div style={{flex: 10, overflow: "auto", padding: "20px"}}>
                {/*预览区域*/}
                <div style={{ paddingBottom: "100px" }}>
                    <ReactMarkdown>{content}</ReactMarkdown>
                </div>
                <div style={{flex: 4, overflow: "auto", padding: "20px", backgroundColor: "#fafafa"}}>
                    <CommentSection docId={docId} />
                </div>
            </div>
            {/* 评论区域 */}


            {/* 悬浮按钮 */}
            <div
                style={{
                    position: "fixed",
                    bottom: "20px",
                    right: "20px",
                    display: "flex",
                    flexDirection: "column",
                    gap: "10px",
                }}
            >
                {/* 编辑按钮 */}
                <Button
                    type="primary"
                    icon={<EditOutlined />}
                    onClick={toggleEdit}
                    style={{
                        width: "60px",
                        height: "60px",
                        borderRadius: "30px",
                        display: "flex",
                        justifyContent: "center",
                        alignItems: "center",
                        fontSize: "18px",
                    }}
                />
                {/* 点赞按钮 */}
                <Button
                    type="default"
                    icon={liked ? <LikeFilled style={{ color: "#1890ff" }} /> : <LikeOutlined />}
                    onClick={handleLike}
                    style={{
                        width: "60px",
                        height: "60px",
                        borderRadius: "30px",
                        display: "flex",
                        justifyContent: "center",
                        alignItems: "center",
                        fontSize: "18px",
                    }}
                />
            </div>
        </div>
    );
};

export default DocumentViewContent;
