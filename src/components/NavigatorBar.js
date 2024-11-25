import React, { useState } from 'react';
import { Layout, Button, Avatar, Dropdown, Menu, Badge, message, Modal } from 'antd';
import { FileAddOutlined, UploadOutlined, BellOutlined, PoweroffOutlined, UserOutlined, BookOutlined, FileTextOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import * as jwt from "../util/jwt";
import * as userService from "../services/userService";
import NewModel from "./models/NewModel";
const { Header } = Layout;

const NavigatorBar = () => {
    const navigate = useNavigate();  // 使用 useNavigate
    const [isNewModalVisible, setIsNewModalVisible] = useState(false);  // 控制弹窗显示;
    // 控制新建弹窗的显示
    const showCreateModal = () => {
        setIsNewModalVisible(true);
    };
    // 控制新建弹窗的关闭
    const handleCreateCancel = () => {
        setIsNewModalVisible(false);
    };
    // 处理登出功能
    const handleLogout = async () => {
        const response = await userService.getUserInfo();
        try{
            if (response.status === 200) {
                // 执行登出操作，比如清理本地存储、重定向到登录页面等
                await userService.logoutRequest()
                jwt.delToken()
            }
        }catch (error) {
            if(error.response.status === 401){
                message.error(error.response.data?.error || "Failed to logout");
                navigate('/login');
            }
        }

        console.log('Logging out...');
        navigate('/login');  // 使用 navigate 进行页面跳转
    };

    // 处理通知点击
    const handleNotificationClick = () => {
        console.log('Viewing notifications...');
    };

    // 处理头像点击
    const handleAvatarClick = () => {
        console.log('Viewing profile...');
    };

    // 用户头像的下拉菜单
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
                {/* 通知按钮 */}
                <Badge count={5} onClick={handleNotificationClick}>
                    <Button icon={<BellOutlined />} style={{ marginRight: '16px' }} />
                </Badge>

                {/* 用户头像与下拉菜单 */}
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
            <NewModel visable={isNewModalVisible} onClose={handleCreateCancel}/>
        </Header>
    );
};

export default NavigatorBar;
