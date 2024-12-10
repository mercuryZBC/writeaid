import React, { useState } from 'react';
import { Button, Modal } from 'antd';
import { BookOutlined, FileTextOutlined } from '@ant-design/icons';
import NewKnowledgeBaseModal from './NewKnowledgeBaseModal';
import NewDocumentModal from './NewDocumentModal';

const NewModel = ({visable,onClose}) => {
    const [isNewKBModelVisible, setIsNewKBModelVisible] = useState(false);
    const [isNewDocModelVisible, setIsNewDocModelVisible] = useState(false);

    // 控制创建知识库弹窗的显示
    const showNewKBModal = () => {
        setIsNewKBModelVisible(true);
    };

    // 控制创建文档弹窗的显示
    const showNewDocModal = () => {
        setIsNewDocModelVisible(true);
    };

    // 关闭创建知识库弹窗
    const handleNewKBClose = () => {
        setIsNewKBModelVisible(false);
    };

    // 关闭创建文档弹窗
    const handleNewDocClose = () => {
        setIsNewDocModelVisible(false);
    };

    return (
        <Modal
            title="选择要创建的内容"
            open={visable}
            onCancel={onClose}
            footer={null} // 去掉默认的"确定"和"取消"按钮
            width={300}
        >
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                {/* 创建知识库按钮 */}
                <Button
                    icon={<BookOutlined />}
                    type="primary"
                    style={{ marginBottom: '16px', width: '100%' }}
                    onClick={showNewKBModal}
                >
                    创建知识库
                </Button>
                <NewKnowledgeBaseModal
                    visible={isNewKBModelVisible}
                    onCancel={handleNewKBClose}
                />
                {/* 创建文档按钮 */}
                <Button
                    icon={<FileTextOutlined />}
                    type="default"
                    style={{ width: '100%' }}
                    onClick={showNewDocModal}
                >
                    创建文档
                </Button>
                <NewDocumentModal visible={isNewDocModelVisible} onClose={handleNewDocClose} />
            </div>
        </Modal>
    );
};

export default NewModel;
