import React from 'react';
import { Menu } from 'antd';
import { UserOutlined, PoweroffOutlined } from '@ant-design/icons';

const UserMenu = ({ handleAvatarClick, handleLogout }) => (
    <Menu>
        <Menu.Item onClick={handleAvatarClick} icon={<UserOutlined />}>
            个人资料
        </Menu.Item>
        <Menu.Item onClick={handleLogout} icon={<PoweroffOutlined />}>
            退出登录
        </Menu.Item>
    </Menu>
);

export default UserMenu;
