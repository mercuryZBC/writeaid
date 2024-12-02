import React, { useState, useEffect } from "react";
import { Input, Button, Tooltip, Avatar, message } from "antd";
import { LikeOutlined, LikeFilled, MessageOutlined } from "@ant-design/icons";
import {formatDate} from "../util/ntil";
import {replyDocumentComment} from "../services/commentService";



const CommentItem = ({ rootCommentId,comment}) => {
    const [replies, setReplies] = useState([]);
    const [showReplies, setShowReplies] = useState(false); // 控制二级评论是否显示
    const [newReply, setNewReply] = useState(""); // 回复内容
    const [isReplying, setIsReplying] = useState(false); // 判断是否在回复
    const handleReply = async (parentId) => {
        if (!newReply.trim()) return;
        // 这里假设提交回复的逻辑
        try {
            console.log(rootCommentId+" "+comment.comment_id+" "+comment.doc_id+" "+newReply);
            const response = await replyDocumentComment(rootCommentId,comment.comment_id,comment.doc_id,newReply);
            if (response.status === 200) {
                message.success('评论发表成功');
            }
        }catch(error) {
            message.error(error?.response?.data?.error || '评论发表失败');
        }
        setIsReplying(false)
        setNewReply("");
    };

    const toggleReplies = async () => {
        if (showReplies) {
            setShowReplies(false);
        } else {

        }
    };

    return (
        <div style={{ marginBottom: "16px", borderBottom: "1px solid #f0f0f0", paddingBottom: "8px" }}>
            <div style={{ display: "flex", alignItems: "flex-start" }}>
                <Avatar style={{ marginRight: "8px" }}>
                    {comment.nickname.charAt(0)} {/* 显示昵称的首字母 */}
                </Avatar>
                <div style={{ flex: 1 }}>
                    <div style={{ fontWeight: "bold", marginBottom: "4px" }}>{comment.nickname}</div>
                    <div style={{ marginBottom: "8px" }}>{comment.comment_content}</div>
                    <div style={{ color: "#aaa", fontSize: "12px", marginBottom: "8px" }}>
                        <Tooltip title={comment.last_updated_at}>
                            <span>{formatDate(comment.last_updated_at)}</span>
                        </Tooltip>
                    </div>
                    <div style={{ display: "flex", gap: "8px" }}>
                        <span>
                            {comment.comment_like_count > 0 ? (
                                <>
                                    <LikeFilled style={{ color: "blue" }} /> {comment.comment_like_count}
                                </>
                            ) : (
                                <LikeOutlined />
                            )}
                        </span>
                        <span onClick={() => setIsReplying(!isReplying)}>
                            <MessageOutlined /> 回复
                        </span>
                        {comment.have_children_comment && (
                            <span onClick={toggleReplies}>
                                {showReplies ? "收起回复" : "展开回复"}
                            </span>
                        )}
                    </div>
                </div>
            </div>

            {/* 显示二级评论 */}
            {showReplies && replies.length > 0 && (
                <div style={{ paddingLeft: "48px", marginTop: "8px" }}>
                    {replies.map((reply) => (
                        <CommentItem rootCommentId={rootCommentId} comment={reply} />
                    ))}
                </div>
            )}

            {/* 显示回复输入框 */}
            {isReplying && (
                <div style={{ paddingLeft: "48px", marginTop: "8px" }}>
                    <Input.TextArea
                        rows={3}
                        value={newReply}
                        onChange={(e) => setNewReply(e.target.value)}
                        placeholder="输入回复内容..."
                    />
                    <Button
                        type="primary"
                        onClick={handleReply}
                        style={{ marginTop: "10px" }}
                    >
                        提交回复
                    </Button>
                </div>
            )}
        </div>
    );
};

export default CommentItem;
