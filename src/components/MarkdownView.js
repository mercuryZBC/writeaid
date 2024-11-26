import React from "react";
import { Input, Button } from "antd";

const { TextArea } = Input;

const MarkdownView = ({ content, setContent, onClose }) => {
    const handleSave = () => {
        onClose(); // 关闭编辑状态
    };

    return (
        <div>
            <TextArea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                rows={10}
                style={{ marginBottom: "10px" }}
            />
            <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px" }}>
                <Button onClick={onClose}>取消</Button>
                <Button type="primary" onClick={handleSave}>保存</Button>
            </div>
        </div>
    );
};

export default MarkdownView;
