import React, { useEffect, useState } from "react";
import { Drawer, List, Button, message } from "antd";
import { useNavigate } from "react-router-dom";
import { getDocumentListByKbId } from "../../services/documentService";

const DocumentListDrawer = ({ kbId,kbName,visible, onClose }) => {
    const [documents, setDocuments] = useState([]);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    // 拉取文档列表
    const fetchDocuments = async () => {
        if (!kbId) return;
        setLoading(true);
        try {
            const response = await getDocumentListByKbId(kbId);
            if (response.status === 200) {
                const docList = response.data['doc_list'];
                if(docList != null) {
                    setDocuments(docList);
                }
            } else {
                message.error("文档拉取失败");
            }
        } catch (error) {
            message.error(error.response?.data?.error || "系统错误，文档拉取失败");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (visible) {
            fetchDocuments();
        }
    }, [kbId, visible]);

    // 点击文档后跳转到预览界面
    const handleDocumentClick = (docId) => {
        navigate(`/document/${docId}`); // 假设文档预览路径为 /document/:id
    };

    return (
        <Drawer
            title={`知识库 ${kbName} 下的文档`}
            placement="right"
            width={400}
            onClose={onClose}
            open={visible}
            destroyOnClose
        >
            <List
                loading={loading}
                dataSource={documents}
                renderItem={(item) => (
                    <List.Item
                        actions={[
                            <Button
                                type="link"
                                onClick={() => handleDocumentClick(item.doc_id)}
                            >
                                查看
                            </Button>,
                        ]}
                    >
                        {item.doc_title}
                    </List.Item>
                )}
            />
        </Drawer>
    );
};

export default DocumentListDrawer;
