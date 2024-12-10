import React, { useState, useEffect } from 'react';
import { Modal, message, List, Input, Space } from 'antd';
import { getKnowledgeBaseList } from '../../services/knowledgeService';  // 假设有一个获取知识库列表的 API
import { createDocument } from '../../services/documentService';  // 假设有一个创建文档的 API

const NewDocumentModal = ({ visible, onClose }) => {
    const [kbList, setKbList] = useState([]);
    const [selectedKB, setSelectedKB] = useState(null);  // 当前选择的知识库
    const [docName, setDocName] = useState('');  // 文档名称
    const [loading, setLoading] = useState(false);  // 加载状态

    useEffect(() => {
        const fetchKnowledgeBases = async () => {
            try {
                const response = await getKnowledgeBaseList();
                if (response.status === 200) {
                    setKbList(response.data.knowledge_bases);
                } else {
                    message.error('知识库加载失败');
                }
            } catch (error) {
                message.error('获取知识库失败');
            }
        };
        fetchKnowledgeBases();
    }, [visible]);

    const handleCreateDocument = async () => {
        if (!selectedKB || !docName) {
            message.error('请选择知识库并填写文档名称');
            return;
        }

        setLoading(true);
        try {
            console.log(selectedKB.kb_id,docName);
            const response = await createDocument(selectedKB.kb_id,docName);

            if (response.status === 200) {
                message.success('文档创建成功');
                setDocName('');  // 重置文档名称
                setSelectedKB(null);  // 清除选择的知识库
                onClose();  // 关闭弹窗
            } else {
                message.error('文档创建失败');
            }
        } catch (error) {
            message.error('文档创建失败');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal
            title="创建文档"
            visible={visible}
            onCancel={onClose}
            onOk={handleCreateDocument}
            confirmLoading={loading}
            destroyOnClose
        >
            <div>
                <Space direction="vertical" style={{ width: '100%' }}>
                    <Input
                        placeholder="请输入文档名称"
                        value={docName}
                        onChange={(e) => setDocName(e.target.value)}
                    />
                    <List
                        bordered
                        dataSource={kbList}
                        renderItem={(kb) => (
                            <List.Item
                                style={{ cursor: 'pointer' }}
                                onClick={() => setSelectedKB(kb)}
                                className={selectedKB && selectedKB.kb_id === kb.kb_id ? 'selected' : ''}
                            >
                                {kb.kb_name}
                            </List.Item>
                        )}
                    />
                </Space>
                {selectedKB && (
                    <div style={{ marginTop: 16 }}>
                        <p>当前选择的知识库：{selectedKB.kb_name}</p>
                    </div>
                )}
            </div>
        </Modal>
    );
};

export default NewDocumentModal;
