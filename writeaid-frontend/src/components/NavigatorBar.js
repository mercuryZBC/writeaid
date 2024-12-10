import React, { useState, useEffect } from 'react';
import { Layout, Button, Avatar, Dropdown, Menu, Badge, message, Spin } from 'antd';
import { FileAddOutlined, UploadOutlined, BellOutlined, PoweroffOutlined, UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import * as jwt from "../util/jwt";
import * as userService from "../services/userService";
import NewModel from "./models/NewModel";
import {delToken} from "../util/jwt";

const { Header } = Layout;

const NavigatorBar = () => {
    const navigate = useNavigate();
    const [isNewModalVisible, setIsNewModalVisible] = useState(false);
    const [username, setUsername] = useState(null); // 存储用户名
    const [loading, setLoading] = useState(true); // 加载状态

    // 获取用户信息
    useEffect(() => {
        const fetchUserInfo = async () => {
            try {
                const response = await userService.getUserInfo();
                if (response.status === 200) {
                    setUsername(response.data['nickname'] || "未命名用户"); // 假设返回的用户名字段为 `username`
                }
            } catch (error) {
                if(error.response && error.response.status === 401) {
                    delToken();
                    navigate('/login')
                }else{
                    message.error(error.response?.data?.error || "获取用户数据失败");
                }
            } finally {
                setLoading(false);
            }
        };

        fetchUserInfo();
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

        console.log('Logging out...');
        navigate('/login');
    };

    const handleNotificationClick = () => {
        console.log('Viewing notifications...');
    };

    const handleAvatarClick = () => {
        console.log('Viewing profile...');
    };

    const userMenu = (
        <Menu>
            <Menu.Item onClick={handleAvatarClick} icon={<UserOutlined />}>
                个人资料
            </Menu.Item>
            <Menu.Item onClick={handleLogout} icon={<PoweroffOutlined />}>
                退出登录
            </Menu.Item>
        </Menu>
    );

    return (
        <Header style={{ padding: '0 16px', background: '#fff', display: 'flex', justifyContent: 'space-between' }}>
            {/* 左侧按钮区域 */}
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

            {/* 右侧按钮区域 */}
            <div style={{ display: 'flex', alignItems: 'center' }}>
                <Badge count={5} style={{ marginRight: '16px' }}>
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
                <Dropdown overlay={userMenu} trigger={['click']}>
                    <Avatar
                        style={{ backgroundColor: '#87d068', cursor: 'pointer' }}
                        icon={<UserOutlined />}
                        size="large"
                        onClick={handleAvatarClick}
                    />
                </Dropdown>
            </div>

            {/* 创建弹窗 */}
            <NewModel visable={isNewModalVisible} onClose={handleCreateCancel} />
        </Header>
    );
};

export default NavigatorBar;
