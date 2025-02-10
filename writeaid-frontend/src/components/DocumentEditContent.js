import React, {useEffect, useState} from "react";
import "vditor/dist/index.css";
import "../styles/editorPage.css"
import {Layout, message, Spin} from "antd";
import {getDocumentById, updateDocumentById} from "../services/documentService";
import Vditor from "vditor";
import '../App.css'
import {data} from "react-router-dom";

const {Header, Content} = Layout;

const DocumentEditContent = ({docId}) => {
    const [vd, setVd] = useState(null);  // 不使用类型注释
    const [docTitle, setDocTitle] = useState("");
    const [saveStatus, setSaveStatus] = useState(""); // 保存状态
    useEffect(() => {
        console.log(docId);
        if(docId){
            setTimeout(()=>{
                setVd(vd =>{
                    fetchDocumentData(vd);
                    fetchDocumentInfo();
                    return vd;
                })
            },60);
        }
    }, [docId]);

    // 自动保存逻辑
    useEffect(() => {
        if (!vd) return;

        const intervalId = setInterval(() => {
            console.log("自动保存...");
            saveDocumentData();
        }, 5000); // 每 10 秒自动保存一次

        return () => clearInterval(intervalId); // 组件卸载时清除定时器
    }, [vd]);


    // 监听键盘事件
    // 监听键盘事件，保存文档
    useEffect(() => {
        const handleKeyDown = (event) => {
            if ((event.ctrlKey || event.metaKey) && event.key === "s") {
                event.preventDefault(); // 阻止默认保存操作
                saveDocumentData(); // 执行保存操作
            }
        };

        window.addEventListener("keydown", handleKeyDown);
        return () => {
            window.removeEventListener("keydown", handleKeyDown);
        };
    }, [vd]); // 依赖于 `vd`

    useEffect(() => {
        const vditor = new Vditor("vditor", {
            after: () => {
                vditor.setValue('');
                setVd(vditor);
            },
        });
        // 清理函数
        return () => {
            vd?.destroy();  // 销毁 Vditor 实例
            setVd(undefined);  // 重置状态
        };
    },[]);


    // 获取文档内容
    const fetchDocumentData = async (vditorInstance) => {
        try {
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                if (vditorInstance) {
                    console.log(data.doc_content);
                    vditorInstance.setValue(data.doc_content);
                }
            } else {
                message.error("文档信息拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "文档信息拉取失败");
        }
    };


    const fetchDocumentInfo = async () =>{
            try {
                const response = await getDocumentById(docId);
                if (response.status === 200) {
                    const data = response.data;
                    setDocTitle(data["doc_title"]);
                } else {
                    message.error("文档对应的知识库ID获取失败");
                }
            } catch (error) {
                message.error("文档对应的知识库ID获取失败");
            }
    };

    // 保存文档内容
    const saveDocumentData = async () => {
        setVd(vd =>{
            if(!vd){
                message.error("编辑器尚未初始化");
                return;
            }
            return vd;
        })
        try {
            const content = vd.getValue();
            const blob = new Blob([content], { type: "text/plain" });
            const formData = new FormData();
            formData.append("file", blob, `${docId}.txt`);
            setSaveStatus("自动保存中...");
            const response = await updateDocumentById(docId, formData);
            if (response.status === 200) {
                const now = new Date();
                const hours = String(now.getHours()).padStart(2, '0');   // 获取小时并补充为两位数
                const minutes = String(now.getMinutes()).padStart(2, '0'); // 获取分钟并补充为两位数
                const seconds = String(now.getSeconds()).padStart(2, '0'); // 获取秒并补充为两位数
                // message.success("最近保存时间", `${hours}:${minutes}:${seconds}`);
                setSaveStatus(`最近保存时间: ${hours}:${minutes}:${seconds}`);
            } else {
                message.error("文档保存失败");
            }
        } catch (error) {
            message.error("文档保存失败");
        }
    };



    return (
        <>
            <Header style={{
                height: "40px",  // 设置固定高度
                padding: "0 6px",
                background: "#fff",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center" // 确保内容垂直居中
            }}>
                <div className="navigatorBar">
                    <span className="docTitle">{docTitle}</span>
                    <span className="saveStatus">
                    {saveStatus === "自动保存中..." ? <Spin size="small"/> : null} {saveStatus}
                </span>
                </div>
            </Header>
            <Content>
                <div id="vditor" className="vditor"></div>
            </Content>
        </>
    );
};


export default DocumentEditContent;

