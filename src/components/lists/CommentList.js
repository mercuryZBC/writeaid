import React, { useState, useEffect } from "react";
import {Input, Button, Tooltip, Avatar, message, List} from "antd";
import { LikeOutlined, LikeFilled, MessageOutlined } from "@ant-design/icons";
import {getDocumentRootComment, replyDocumentComment} from "../../services/commentService";
import {formatDate} from "../../util/ntil";
import ChildrenCommentList from "./ChildrenCommentList";



const CommentList = ({ docId,refresh }) => {
    const [commentList, setCommentList] = useState([]); // 评论列表

    useEffect(() => {
        if(docId){
            fetchRootComments()
        }
    },[docId,refresh])

    const  fetchRootComments = async () => {
        try{
            const response = await getDocumentRootComment(docId,0,10);
            if (response.status === 200) {
                console.log(response.data);
                if(response.data.comment_list) {
                    setCommentList(response.data.comment_list);
                }else{
                    setCommentList([]);
                }
            }
        }catch(error){
            message.error(error?.response?.error || '评论信息加载失败');
        }
    }
    const renderComments = (data) => {
        return data.map((comment) => (
            <RootCommentItem rootComment={comment}>
            </RootCommentItem>
        ));
    };
    return (
        <List
            dataSource={commentList}
            renderItem={(item) => <>{renderComments([item])}</>}
            locale={{ emptyText: "暂无评论" }}
        />
    )
}


const RootCommentItem = ({rootComment}) => {
    const [showReplies, setShowReplies] = useState(false); // 控制二级评论是否显示
    const [newReply, setNewReply] = useState(""); // 回复内容
    const [isReplying, setIsReplying] = useState(false); // 判断是否在回复
    const [currentReplyCommentId, setCurrentReplyCommentId] = useState(null);
    const handleReply = async () => {
        // 评论为空直接返回
        if (!newReply.trim()) return;
        try {
            const response = await replyDocumentComment(rootComment.comment_id,currentReplyCommentId,rootComment.doc_id,newReply);
            if (response.status === 200) {
                message.success('评论回复成功');
            }
        }catch(error) {
            message.error(error?.response?.data?.error || '评论回复失败');
        }
        setIsReplying(false)
        setNewReply("");
    };

    const setCurrentReplyCommentIdHandle = (replyCommentId) =>{
        if(replyCommentId === currentReplyCommentId){
            setIsReplying(false);
            setCurrentReplyCommentId("");
        }else{
            setIsReplying(true);
            setCurrentReplyCommentId(replyCommentId);
        }
    }

    const toggleReplies = async () => {
        if (showReplies) {
            setShowReplies(false);
        } else {  //获取次级评论
            setShowReplies(true);
        }
    };

    return (
        <div style={{ marginBottom: "16px", borderBottom: "1px solid #f0f0f0", paddingBottom: "8px" }}>
            <div style={{ display: "flex", alignItems: "flex-start" }}>
                <Avatar style={{ marginRight: "8px" }}>
                    {rootComment.nickname.charAt(0)} {/* 显示昵称的首字母 */}
                </Avatar>
                <div style={{ flex: 1 }}>
                    <div style={{ fontWeight: "bold", marginBottom: "4px" }}>{rootComment.nickname}</div>
                    <div style={{ marginBottom: "8px" }}>{rootComment.comment_content}</div>
                    <div style={{ color: "#aaa", fontSize: "12px", marginBottom: "8px" }}>
                        <Tooltip title={rootComment.last_updated_at}>
                            <span>{formatDate(rootComment.last_updated_at)}</span>
                        </Tooltip>
                    </div>
                    <div style={{ display: "flex", gap: "8px" }}>
                        <span>
                            {rootComment.comment_like_count > 0 ? (
                                <>
                                    <LikeFilled style={{ color: "blue" }} /> {rootComment.comment_like_count}
                                </>
                            ) : (
                                <LikeOutlined />
                            )}
                        </span>
                        <span onClick={() => setCurrentReplyCommentIdHandle(rootComment.comment_id)}>
                            <MessageOutlined /> 回复
                        </span>
                        {rootComment.have_children_comment && (
                            <span onClick={toggleReplies}>
                                {showReplies ? "收起回复" : "展开回复"}
                            </span>
                        )}
                    </div>
                </div>
            </div>

            {/* 显示二级评论 */}
            {showReplies && (
                <div style={{ paddingLeft: "48px", marginTop: "8px" }}>
                        <ChildrenCommentList rootCommentId={rootComment.comment_id}
                                             showChildReplies={showReplies}
                                             setCurrentReplyCommentIdHandle={setCurrentReplyCommentIdHandle}
                        />
                </div>
            )}

            {/* 显示回复输入框 */}
            {isReplying && (
                <div style={{ paddingLeft: "48px", marginTop: "8px" }}>
                    <Input.TextArea
                        rows={3}
                        value={newReply}
                        onChange={(e) => setNewReply(e.target.value)}
                        placeholder={`输入回复给评论ID ${currentReplyCommentId} 的内容...`}
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

export default CommentList;
