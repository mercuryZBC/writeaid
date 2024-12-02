import React, {useEffect, useState} from "react";
import {List, Input, Button, Tooltip, Avatar, message} from "antd";
import { LikeOutlined, LikeFilled, MessageOutlined } from "@ant-design/icons";
import {createDocumentComment, getDocumentRootComment} from "../services/commentService";
import CommentItem from "./CommentItem";
import {formatDate} from "../util/ntil";

const { TextArea } = Input;

const CommentSection = ({docId}) => {
    const [commentList, setCommentList] = useState([]); // 评论列表
    const [newComment, setNewComment] = useState(""); // 新评论内容
    const [refresh, setRefresh] = useState(false);
    useEffect(() => {
        if(docId){
            fetchComments()
        }
    },[docId,refresh])

    const  fetchComments = async () => {
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
    const renderComments = (data) => {
        return data.map((comment) => (
            <CommentItem rootCommentId={comment.comment_id} comment={comment}>
            </CommentItem>
        ));
    };

    return (
        <div>
            <h3>评论</h3>
            <List
                dataSource={commentList}
                renderItem={(item) => <>{renderComments([item])}</>}
                locale={{ emptyText: "暂无评论" }}
            />
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
