import React, {useEffect, useState} from "react";
import "vditor/dist/index.css";
import {message} from "antd";
import {getDocumentById, updateDocumentById} from "../services/documentService";
import remarkGfm from 'remark-gfm';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { dracula } from 'react-syntax-highlighter/dist/esm/styles/prism';

const DocumentEditContent = ({docId}) => {
    const [content, setContent] = useState('');

    useEffect(() => {
        console.log(docId);
        if(docId){
            setTimeout(()=>{
                fetchDocumentData();
            },60);
        }
    }, [docId]);

    // 监听键盘事件
    useEffect(() => {
        if(content){
            const handleKeyDown = (event) => {
                // 检测 Ctrl+S 或 Command+S 快捷键
                if ((event.ctrlKey || event.metaKey) && event.key === "s") {
                    event.preventDefault(); // 阻止默认保存操作
                    saveDocumentData(); // 执行保存操作
                }
            };

            // 添加键盘事件监听
            window.addEventListener("keydown", handleKeyDown);

            // 清理事件监听
            return () => {
                window.removeEventListener("keydown", handleKeyDown);
            };
        }
    },[content]); // 当 content 变化时，重新设置监听



    // 获取文档内容
    const fetchDocumentData = async () => {
        try {
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                setContent(data.doc_content);
                console.log(data.doc_content);
            } else {
                message.error("文档信息拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "文档信息拉取失败");
        }
    };

    // 保存文档内容
    const saveDocumentData = async () => {
        try {
            const blob = new Blob([content], { type: "text/plain" }); // 创建文本文件
            const formData = new FormData();
            formData.append('file', blob, docId + '.txt'); // 使用 docId 作为文件名

            // 发送 PUT 请求
            const response = await updateDocumentById(docId, formData);

            // 根据响应状态显示提示
            if (response.status === 200) {
                message.success("文档保存成功");
            } else {
                message.error("文档保存失败");
            }
        } catch (error) {
            message.error("文档保存失败");
        }
    };



    return (
        <div style={{ display: 'flex', height: '100vh' }}>
            {/* 输入区域 */}
            <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="在这里输入 Markdown..."
                style={{
                    width: '50%',
                    padding: '10px',
                    fontSize: '16px',
                    borderRight: '1px solid #ddd',
                }}
            />

            {/* 渲染区域 */}
            <div
                style={{
                    width: '50%',
                    padding: '10px',
                    overflowY: 'scroll',
                    backgroundColor: '#f9f9f9',
                }}
            >
                <ReactMarkdown
                    children={content}
                    remarkPlugins={[remarkGfm]}
                    components={{
                        code({ node, inline, className, children, ...props }) {
                            const match = /language-(\w+)/.exec(className || '');
                            return !inline && match ? (
                                <SyntaxHighlighter
                                    style={dracula}
                                    language={match[1]}
                                    PreTag="div"
                                    {...props}
                                >
                                    {String(children).replace(/\n$/, '')}
                                </SyntaxHighlighter>
                            ) : (
                                <code className={className} {...props}>
                                    {children}
                                </code>
                            );
                        },
                    }}
                />
            </div>
        </div>
    );
};

export default DocumentEditContent;

