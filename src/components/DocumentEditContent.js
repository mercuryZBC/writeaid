import React, {useEffect, useRef, useState} from "react";
import Vditor from "vditor";
import "vditor/dist/index.css";
import {message, Spin} from "antd";
import {getDocumentById, updateDocumentById} from "../services/documentService";

const DocumentEditContent = ({docId}) => {
    // 创建一个可变对象：useRef创建的对象会在组件的整个生命周期保持不变，他返回一个.current属性的对象，可以随时读写该属性的值
    // 不会触发重新渲染
    // 常用于存储不需要触发重新渲染的对象，比如DOM，第三方库实例等
    const vditorRef = useRef(null); // 用于存储 Vditor 实例
    const [content, setContent] = useState(null);
    const [loading,setLoading] = useState(false);
    const [dimensions, setDimensions] = useState({
        width: window.innerWidth,
        height: window.innerHeight,
    });
    // 获取文档内容
    const fetchDocumentData = async () => {
        try {
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                setContent(data.doc_content);
                console.log(data.doc_content);
                setLoading(false);
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
            // 获取最新的文档内容
            const contentValue = vditorRef.current.getValue(); // 获取编辑器的内容
            const blob = new Blob([contentValue], { type: "text/plain" }); // 创建文本文件
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
    useEffect(() => {
        console.log(docId);
        if(docId){
            setTimeout(()=>{
                fetchDocumentData();
            },60);
        }
    }, [docId]);
    useEffect(() => {
        // 等待 DOM 渲染完成后初始化 Vditor
        setTimeout(() => {
            try {
                const container = document.getElementById("vditor");
                if (!container) {
                    console.error("Vditor container not found");
                    return;
                }

                // 初始化 Vditor
                vditorRef.current = new Vditor("vditor", {
                    height: '100vh',
                    lang: "zh_CN",
                    toolbar: [],  // 禁用工具栏（设置为空数组）
                    placeholder: "请输入内容...",
                    mode: "wysiwyg",
                    // 禁用默认边框
                    style: {
                        border: "none", // 去除边框
                    },
                });
                vditorRef.current.setValue("");
            } catch (error) {

            }
        }, 50);

        return () => {
            if (vditorRef.current) {
                vditorRef.current.destroy();
                vditorRef.current = null;
            }
        }

    },[]);

    useEffect(() => {
        if(content){

            setTimeout(() =>{
                vditorRef.current.setValue(content);
            },200)

        }
    },[content]);

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

    return (
        <Spin spinning={loading} tip="文档加载中..." size="large">
            <div
                id="vditor"
                style={{
                    width: `${dimensions.width * 1.2}px`,   // 宽度为当前窗口宽度的 120%
                    height: `${dimensions.height * 1.2}px`, // 高度为当前窗口高度的 120%
                    margin: "20px",
                    overflow: "hidden",  // 防止溢出
                }}
            />
        </Spin>
    );
};

export default DocumentEditContent;
