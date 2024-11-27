import { Layout, message } from "antd";
import { ReloadProvider } from "../context/ReloadContext";
import DocumentViewAndEditSideBar from "../components/sidebars/DocumentViewAndEditSideBar";
import { useLocation } from "react-router-dom";
import DocumentEditContent from "../components/DocumentEditContent";
import { useEffect, useState } from "react";
import { getDocumentById, getDocumentListByKbId } from "../services/documentService";


const { Sider, Content } = Layout;

const DocumentEditPage = () => {
    const location = useLocation();
    const [docId, setDocId] = useState(null);
    const [kbId, setKbId] = useState(null);

    const fetchKBId = async () => {
        try {
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                setKbId(data["kb_id"]);
            } else {
                message.error("文档对应的知识库ID获取失败");
            }
        } catch (error) {
            message.error("文档对应的知识库ID获取失败");
        }
    };

    const fetchKBFirstDoc = async () => {
        try {
            const response = await getDocumentListByKbId(kbId);
            if (response.status === 200) {
                const docList = response.data["doc_list"];
                if (docList && docList.length > 0) {
                    setDocId(docList[0]["doc_id"]);
                }
            } else {
                message.error("获取文档列表失败");
            }
        } catch (error) {
            message.error("获取文档列表失败，请稍后重试");
        }
    };

    useEffect(() => {
        if (location.pathname.indexOf("/editDocument") === 0) {
            setDocId(location.pathname.split("/").pop());
        }
        if (location.pathname.indexOf("/knowledgeBase") === 0) {
            setKbId(location.pathname.split("/").pop());
        }
    }, [location.pathname]);

    useEffect(() => {
        if (kbId && docId === null) {
            fetchKBFirstDoc();
        }
    }, [kbId]);

    useEffect(() => {
        if (docId && kbId === null) {
            fetchKBId();
        }
    }, [docId]);

    return (
        <Layout style={{ height: "100vh", overflow: "hidden" }}>
            <ReloadProvider>
                {/* 左侧导航区域 */}
                <Sider
                    width={300}
                    theme="light"
                    style={{
                        borderRight: "1px solid #f0f0f0",
                        position: "fixed",
                        left: 0,
                        top: 0,
                        bottom: 0,
                        overflow: "auto",
                        backgroundColor: "#fff",
                    }}
                >
                    <DocumentViewAndEditSideBar kbId={kbId} docId={docId} />
                </Sider>

                {/* 右侧编辑区域 */}
                <Layout
                    style={{
                        marginLeft: 300, // 留出侧边栏宽度
                        height: "100vh",
                        overflow: "hidden",
                    }}
                >
                    <Content
                        style={{
                            padding: "24px",
                            background: "#fff",
                            overflowY: "auto", // 允许滚动
                            height: "100%",
                        }}
                    >
                        <DocumentEditContent docId={docId} />
                    </Content>
                </Layout>
            </ReloadProvider>
        </Layout>
    );
};

export default DocumentEditPage;
