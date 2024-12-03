import {Avatar, List, message, Tooltip} from "antd";
import {formatDate} from "../../util/ntil";
import {LikeFilled, LikeOutlined, MessageOutlined} from "@ant-design/icons";
import React, {useEffect, useState} from "react";
import {getChildrenComment} from "../../services/commentService";


const ChildrenCommentItem = ({childrenComment,replayHandle}) => {

    const taggleReplyStatus = () =>{
        replayHandle(childrenComment.comment_id)
    }


    const getReplyContentFormat = (parent_comment_user_nickname,comment_content) => {
        return '回复 '+parent_comment_user_nickname+":"+comment_content;
    }
    return (
        <div style={{ marginBottom: "16px", borderBottom: "1px solid #f0f0f0", paddingBottom: "8px" }}>
            <div style={{ display: "flex", alignItems: "flex-start" }}>
                <Avatar style={{ marginRight: "8px" }}>
                    {childrenComment.nickname.charAt(0)} {/* 显示昵称的首字母 */}
                </Avatar>
                <div style={{ flex: 1 }}>
                    <div style={{ fontWeight: "bold", marginBottom: "4px" }}>{childrenComment.nickname}</div>
                    <div style={{ marginBottom: "8px" }}>{getReplyContentFormat(childrenComment.parent_comment_user_nickname,
                        childrenComment.comment_content)}</div>
                    <div style={{ color: "#aaa", fontSize: "12px", marginBottom: "8px" }}>
                        <Tooltip title={childrenComment.comment_updated_at}>
                            <span>{formatDate(childrenComment.comment_updated_at)}</span>
                        </Tooltip>
                    </div>
                    <div style={{ display: "flex", gap: "8px" }}>
                        <span>
                            {childrenComment.comment_like_count > 0 ? (
                                <>
                                    <LikeFilled style={{ color: "blue" }} /> {childrenComment.comment_like_count}
                                </>
                            ) : (
                                <LikeOutlined />
                            )}
                        </span>
                        <span onClick={() => taggleReplyStatus()}>
                            <MessageOutlined /> 回复
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );

}



const ChildrenCommentList = ({ rootCommentId,showChildReplies,setCurrentReplyCommentIdHandle }) => {
    const [childrenCpmmentList, setChildrenCpmmentList] = useState([]);
    useEffect(()=>{
        if(showChildReplies){
            fetchChildrenComments()
        }
    },[showChildReplies])

    const replayHandle =  (replyCommentId) => {
        setCurrentReplyCommentIdHandle(replyCommentId);
    }

    const fetchChildrenComments = async () => {
        try{
            const response = await getChildrenComment(rootCommentId);
            if (response.status === 200) {
                console.log(response.data.children_comments);
                if(response.data.children_comments){
                    setChildrenCpmmentList([response.data.children_comments]);
                }else{
                    setChildrenCpmmentList([]);
                }
            }
        }catch (error) {
            message.error(error.response?.error || "回复评论获取失败，请稍后再试");
        }
    }
    const renderChildrenComment= (data) => {
        return data.map((comment) => (
            <ChildrenCommentItem childrenComment={comment} replayHandle={replayHandle}/>
        ));
    };
    return (
        <>
            <List
                dataSource={childrenCpmmentList}
                renderItem={(item) => <>{renderChildrenComment(item)}</>}
                locale={{ emptyText: "暂无回复评论" }}
            />
        </>
    )
}

export default ChildrenCommentList;
