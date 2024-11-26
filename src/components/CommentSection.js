import React, { useState } from "react";
import { List, Input, Button, Tooltip } from "antd";

const { TextArea } = Input;

const CommentSection = () => {
    const [comments, setComments] = useState([]); // 评论列表
    const [newComment, setNewComment] = useState(""); // 新的评论内容

    const handleAddComment = () => {
        if (!newComment.trim()) return;
        const newEntry = {
            author: "匿名用户",
            content: newComment,
            datetime: new Date().toLocaleString(),
        };
        setComments([...comments, newEntry]);
        setNewComment("");
    };

    return (
        <div>
            <h3>评论</h3>
            <List
                // dataSource={comments}
                // renderItem={(item) => (
                //     <Comment
                //         author={item.author}
                //         content={item.content}
                //         datetime={<Tooltip title={item.datetime}>{item.datetime}</Tooltip>}
                //     />
                // )}
            />
            <div style={{ marginTop: "10px" }}>
                <TextArea
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                    rows={3}
                    placeholder="添加评论..."
                />
                <Button type="primary" onClick={handleAddComment} style={{ marginTop: "10px" }}>
                    提交评论
                </Button>
            </div>
        </div>
    );
};

export default CommentSection;
