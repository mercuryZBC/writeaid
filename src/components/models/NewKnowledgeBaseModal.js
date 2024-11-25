import React, { useState } from 'react';
import { Modal, Input, Checkbox, Button, message } from 'antd';
import { BookOutlined } from '@ant-design/icons';
import * as documentService from "../../services/knowledgeService";
import {useReloadContext} from "../../context/ReloadContext";

const NewKnowledgeBaseModal = ({ visible, onCancel}) => {
    const [knowledgeBaseName, setKnowledgeBaseName] = useState('');
    const [knowledgeBaseDescription, setKnowledgeBaseDescription] = useState('');
    const [isPublic, setIsPublic] = useState(false);
    // 重新加载知识库列表
    const {kbListReloadTrigger, setkbListReloadTrigger} = useReloadContext();
    // 提交创建知识库
    const handleCreate = async () => {
        if (!knowledgeBaseName) {
            message.error('请输入知识库名');
            return;
        }
        try{
            const response = await documentService.createKnowledgeBase(
                knowledgeBaseName,knowledgeBaseDescription,isPublic);
            if(response.status === 200) {
                setkbListReloadTrigger(!kbListReloadTrigger);
                message.success("知识库创建成功")
            }
        }catch(error){
            console.log(error);
            message.error(error.response?.data?.error || "知识库创建失败");
        }
        // 重置状态
        setKnowledgeBaseName('');
        setKnowledgeBaseDescription('');
        setIsPublic(false);
        onCancel();
    };

    return (
        <Modal
            title="创建知识库"
            open={visible}
            onCancel={onCancel}
            footer={null}
            width={400}
        >
            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                {/* 知识库名称 */}
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <BookOutlined style={{ marginRight: '8px' }} />
                    <Input
                        placeholder="请输入知识库名"
                        value={knowledgeBaseName}
                        onChange={(e) => setKnowledgeBaseName(e.target.value)}
                    />
                </div>

                {/* 知识库简介 */}
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <BookOutlined style={{ marginRight: '8px' }} />
                    <Input.TextArea
                        placeholder="请输入知识库简介"
                        value={knowledgeBaseDescription}
                        onChange={(e) => setKnowledgeBaseDescription(e.target.value)}
                        rows={4}  // 设置显示的行数
                        style={{
                            resize: 'none',  // 禁用输入框的大小调整
                        }}
                    />
                </div>

                {/* 是否公开 */}
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <Checkbox
                        checked={isPublic}
                        onChange={(e) => setIsPublic(e.target.checked)}
                    />
                    <span style={{ marginLeft: '8px' }}>是否公开</span>
                </div>

                {/* 提交按钮 */}
                <div style={{ marginTop: '16px', textAlign: 'right' }}>
                    <Button type="primary" onClick={handleCreate}>
                        创建
                    </Button>
                </div>
            </div>
        </Modal>
    );
};

export default NewKnowledgeBaseModal;
