import React, { useState, useEffect } from 'react';
import { Layout, Button, Avatar, Dropdown, Badge, message, Spin } from 'antd';
import { FileAddOutlined, UploadOutlined, BellOutlined, UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import * as jwt from "../util/jwt";
import * as userService from "../services/userService";
import * as notificationService from "../services/notificationService";
import NewModel from "./models/NewModel";
import UserMenu from '../components/menus/UserMenu';
import NotificationModal from '../components/models/NotificationModal';
import { delToken } from "../util/jwt";

const { Header } = Layout;

const NavigatorBar = () => {
    const navigate = useNavigate();
    const [isNewModalVisible, setIsNewModalVisible] = useState(false);
    const [username, setUsername] = useState(null);
    const [loading, setLoading] = useState(true);
    const [notificationCount, setNotificationCount] = useState(0);
    const [isNotificationModalVisible, setIsNotificationModalVisible] = useState(false);
    const [notifications, setNotifications] = useState([]);

    const fetchUserInfo = async () => {
        try {
            const response = await userService.getUserInfo();
            if (response.status === 200) {
                setUsername(response.data['nickname'] || "未命名用户");
            }
        } catch (error) {
            if (error.response && error.response.status === 401) {
                delToken();
                navigate('/login');
            } else {
                message.error(error.response?.data?.error || "获取用户数据失败");
            }
        } finally {
            setLoading(false);
        }
    };

    const fetchNotifications = async () => {
        try {
            const response = await notificationService.getNotifications();
            if (response.status === 200) {
                const formattedNotifications = response.data["nt_content"].map((item) => ({
                    ...item,
                    nt_time: new Date(item.nt_time * 1000).toLocaleString(),
                }));
                setNotifications(formattedNotifications);
                console.log(formattedNotifications);
            }
        } catch (error) {
            if (error.response && error.response.status === 401) {
                delToken();
                navigate('/login');
            } else {
                console.log(error.response?.data?.error || "获取通知数据失败");
            }
        }
    };

    const fetchNotificationCount = async () => {
        try {
            const response = await notificationService.getNotificationCount();
            if (response.status === 200) {
                setNotificationCount(response.data["nt_count"]);
            }
        } catch (error) {
            if (error.response && error.response.status === 401) {
                delToken();
                navigate('/login');
            } else {
                console.log(error.response?.data?.error || "获取通知数量失败");
            }
        }
    };

    useEffect(() => {
        fetchUserInfo();
        fetchNotificationCount();
        fetchNotifications();
    }, []);

    const showCreateModal = () => {
        setIsNewModalVisible(true);
    };

    const handleCreateCancel = () => {
        setIsNewModalVisible(false);
    };

    const handleLogout = async () => {
        try {
            const response = await userService.getUserInfo();
            if (response.status === 200) {
                await userService.logoutRequest();
                jwt.delToken();
            }
        } catch (error) {
            if (error.response.status === 401) {
                message.error(error.response.data?.error || "Failed to logout");
                navigate('/login');
            }
        }
        navigate('/login');
    };

    const handleNotificationClick = () => {
        setIsNotificationModalVisible(true);
    };

    const handleNotificationModalClose = () => {
        setIsNotificationModalVisible(false);
    };

    const handleAvatarClick = () => {
        console.log('Viewing profile...');
    };

    return (
        <Header style={{ padding: '0 16px', background: '#fff', display: 'flex', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center' }}>
                <Button
                    icon={<FileAddOutlined />}
                    type="primary"
                    style={{ marginRight: '16px' }}
                    onClick={showCreateModal}
                >
                    新建
                </Button>
                <Button icon={<UploadOutlined />} type="default">
                    导入
                </Button>
            </div>
            <div style={{ display: 'flex', alignItems: 'center' }}>
                <Badge count={notificationCount} style={{ marginRight: '16px' }}>
                    <Button
                        type="text"
                        icon={<BellOutlined style={{ fontSize: '24px', color: '#08c'}} />}
                        style={{ marginRight: '20px' }}
                        onClick={handleNotificationClick}
                    />
                </Badge>
                {loading ? (
                    <Spin size="small" style={{ marginRight: '8px' }} />
                ) : (
                    <span style={{ marginRight: '16px', fontWeight: 'bold' }}>{username}</span>
                )}
                <Dropdown overlay={<UserMenu handleAvatarClick={handleAvatarClick} handleLogout={handleLogout} />} trigger={['click']}>
                    <Avatar
                        style={{ backgroundColor: '#87d068', cursor: 'pointer' }}
                        icon={<UserOutlined />}
                        size="large"
                        onClick={handleAvatarClick}
                    />
                </Dropdown>
            </div>
            <NewModel visable={isNewModalVisible} onClose={handleCreateCancel} />
            <NotificationModal
                isVisible={isNotificationModalVisible}
                notifications={notifications}
                handleClose={handleNotificationModalClose}
            />
        </Header>
    );
};

export default NavigatorBar;
