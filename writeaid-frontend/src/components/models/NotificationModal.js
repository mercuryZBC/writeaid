import React from 'react';
import { Modal, List, Avatar, Badge } from 'antd';

const NotificationModal = ({ isVisible, notifications, handleClose }) => (
    <Modal
        title="通知"
        visible={isVisible}
        footer={null}
        onCancel={handleClose}
    >
        <List
            dataSource={notifications}
            renderItem={(item) => (
                <List.Item key={item.nt_time}>
                    <List.Item.Meta
                        avatar={<Avatar src={item.nt_avatar_link} />}
                        title={<a href={item.nt_link}>{item.nt_message_content}</a>}
                        description={`${item.nt_time} - ${item.nt_link_describe}`}
                    />
                    {!item.nt_is_check && <Badge color="red" text="未读" />}
                </List.Item>
            )}
        />
    </Modal>
);

export default NotificationModal;
