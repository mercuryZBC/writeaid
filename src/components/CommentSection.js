import React, { useState} from "react";
import {Input, Button, message} from "antd";
import {createDocumentComment} from "../services/commentService";
import CommentList from "./lists/CommentList";


const { TextArea } = Input;

const CommentSection = ({docId}) => {
    const [newComment, setNewComment] = useState(""); // 新评论内容
    const [refresh, setRefresh] = useState(false);

    const handleAddComment =async () => {
        if (!newComment.trim()) return;
        try {
            const response = await createDocumentComment(docId,newComment);
            if (response.status === 200) {
                message.success('评论发表成功');
            }
        }catch(error) {
            message.error(error?.response?.data?.error || '评论发表失败');
        }
        setRefresh(!refresh);
        setNewComment("");
    };

    return (
        <div>
            <h3>评论</h3>
            <CommentList docId={docId} refresh={refresh} />
            <div style={{ marginTop: "10px" }}>
                <TextArea
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                    rows={3}
                    placeholder={"添加评论..."}
                />
                <Button
                    type="primary"
                    onClick={handleAddComment}
                    style={{ marginTop: "10px" }}>
                    {"提交评论"}
                </Button>
            </div>
        </div>
    );
};

export default CommentSection;
