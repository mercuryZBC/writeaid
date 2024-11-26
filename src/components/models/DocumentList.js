import {Button, Drawer, List, message} from "antd";
import React, {useEffect, useState} from "react";
import {getDocumentListByKbId} from "../../services/knowledgeService";


const DocumentListDrawer = ({ kb_id,visible, setVisible }) => {
    const [documentList, setDocumentList] = useState([]);
    const handleDocumentSelect = (doc_id) =>{
        // 请求文档数据
    }
    useEffect(() => {
        const getDocumentList = async () => {
            try {
                const response = await getDocumentListByKbId(kb_id);
                if(response.status === 200) { // 返回类型包括id和title
                    setDocumentList(response.data['doc_list']);
                }
            }catch (error) {
                message.error('文档列表拉取失败');
            }
        }
        getDocumentList();
    },[visible]);
    return (
        <>
            <Drawer
                title="选择文档"
                width={400}
                onClose={() => setVisible(false)}
                open={visible}
                bodyStyle={{ paddingBottom: 80 }}
            >
                <List
                    dataSource={documentList}
                    renderItem={(doc) => (
                        <List.Item
                            key={doc.doc_id}
                            onClick={() => handleDocumentSelect(doc.doc_id)}
                            style={{ cursor: "pointer" }}
                        >
                            <Button style={{ width: "100%", textAlign: "left" }}>
                                {doc.title}
                            </Button>
                        </List.Item>
                    )}
                />
            </Drawer>
        </>
    );
}

export default DocumentListDrawer;
