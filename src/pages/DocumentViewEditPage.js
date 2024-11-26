import {Layout, message} from "antd";
import {ReloadProvider} from "../context/ReloadContext";
import DocumentViewAndEditSideBar from "../components/sidebars/DocumentViewAndEditSideBar";
import {useLocation} from "react-router-dom";
import DocumentViewAndEditContent from "../components/DocumentViewAndEditContent";
import {useEffect, useState} from "react";
import {getDocumentById, getDocumentListByKbId} from "../services/documentService";



const { Sider, Content } = Layout;

const DocumentViewEditPage = () => {
    const location = useLocation();
    const [docId, setDocId] = useState(null);
    const [kbId, setKbId] = useState(null);
    const fetchKBId = async () => {
        try {
            // 获取当前文档的知识库 ID
            const response = await getDocumentById(docId);
            if (response.status === 200) {
                const data = response.data;
                setKbId(data['kb_id']);
            } else {
                message.error("文档对应的的知识库ID获取失败");
            }
        } catch (error) {
            message.error("文档对应的的知识库ID获取失败");
        }
    }

    // 获取文档列表
    const fetchKBFirstDoc = async () => {
        try {
            const response = await getDocumentListByKbId(kbId); // 从服务中获取文档列表
            if (response.status === 200) {
                const docList = response.data['doc_list'];
                if (docList != null && docList.length > 0) {
                    setDocId(docList[0]['doc_id']);
                }
            } else {
                message.error("获取文档列表失败");
            }
        } catch (error) {
            message.error("获取文档列表失败，请稍后重试");
        }
    };

    // 根据路由路径来设置 kbId 和 docId
    useEffect(() => {
        if (location.pathname.indexOf("/document") === 0) {
            setDocId(location.pathname.split('/').pop());  // 获取数组中的最后一项
        }
        if (location.pathname.indexOf("/knowledgeBase") === 0) {
            setKbId(location.pathname.split('/').pop());  // 获取数组中的最后一项
        }
    }, [location.pathname]);  // 路径变化时重新执行


    // 在 kbId 更新时调用 fetchDocumentList
    useEffect(() => {
        if (kbId && docId === null) {
            fetchKBFirstDoc();
        }
    }, [kbId]);

    useEffect(() => {
        if(docId && kbId === null){
            fetchKBId();
        }
    },[docId])

    return (
        <Layout style={{ height: "100vh" }}>
            <>
                <ReloadProvider>
                    {/*左侧区域*/}
                    <Sider width={300} theme="light" style={{ borderRight: "1px solid #f0f0f0" }}>
                        <DocumentViewAndEditSideBar kbId={kbId} docId={docId}/>
                    </Sider>

                    {/* 右侧区域 */}
                    <Content style={{ padding: "24px", background: "#fff" }}>
                        <DocumentViewAndEditContent docId={docId} />
                    </Content>
                </ReloadProvider>
            </>
        </Layout>
    );
};

export default DocumentViewEditPage;
