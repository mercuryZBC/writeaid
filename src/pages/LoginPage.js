import React, { useState, useEffect } from "react";
import { Form, Input, Button, message, Tabs, Image } from "antd";
import { loginRequest, registerRequest } from "../services/authService";
import { getCaptcha } from "../services/utilService"; // 假设已经有此接口
import { useNavigate } from 'react-router-dom';
import CryptoJS from "crypto-js";
import {getToken, setToken} from "../util/jwt";
import background from "../assets/images/login_background.jpg";

const LoginPage = () => {
    const [activeTab, setActiveTab] = useState("login");
    const [captchaId, setCaptchaId] = useState(""); // 存储验证码ID
    const [captchaImage, setCaptchaImage] = useState(""); // 存储验证码图片
    const [captchaType, setCaptchaType] = useState("login"); // 验证码类型：登录或注册
    const navigate = useNavigate();
    // 获取验证码
    const fetchCaptcha = async () => {
        try {
            const response = await getCaptcha(); // 传递验证码类型（login 或 register）
            const data = response.data;
            setCaptchaId(data["captchaId"]);
            setCaptchaImage(data["captcha"]);
        } catch (error) {
            message.error("验证码加载失败！");
        }
    };

    // 登录处理
    const handleLogin = async (values) => {
        const hashedPassword = CryptoJS.MD5(values.password).toString(CryptoJS.enc.Base64);
        try {
            const response = await loginRequest(values.email, hashedPassword, captchaId,values.captcha);
            if (response.status) {
                message.success("登录成功！");
                const data = response.data;
                setToken(data['access_token'],data['expires_in'])
                navigate('/home')
            }
        } catch (error) {
            message.error(error.response?.data?.error || "Failed to login");
        }
    };

    // 注册处理
    const handleRegister = async (values) => {
        if (values.password === values.confirmPassword) {
            const hashedPassword = CryptoJS.MD5(values.password).toString(CryptoJS.enc.Base64);
            try {
                const response = await registerRequest(values.email, values.nickname, hashedPassword, captchaId,values.captcha);
                if (response.status === 200) {
                    message.success("注册成功！");
                }
            } catch (error) {
                message.error(error.response.data['error']);
            }
        } else {
            message.error("两次密码输入不一致！");
        }
    };

    useEffect(() => {
        if(getToken() !== ""){
            navigate("/home");
        }
        fetchCaptcha(); // 初始化时获取验证码
    }, [captchaType]);

    // 点击刷新验证码
    const handleCaptchaClick = () => {
        fetchCaptcha();
    };

    // 更改验证码类型（登录/注册）
    const handleTabChange = (key) => {
        setActiveTab(key);
        setCaptchaType(key); // 切换验证码类型
    };

    const backgroundStyle = {
        width: "100vw",
        height: "100vh",
        backgroundImage: `url(${background})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
    };

    return (
        <div style={backgroundStyle}>
            <div style={{ width: 400, padding: 20, backgroundColor: "#fff", borderRadius: 8, boxShadow: "0 2px 8px rgba(0,0,0,0.1)" }}>
                <Tabs
                    defaultActiveKey="login"
                    activeKey={activeTab}
                    onChange={handleTabChange}
                    centered
                >
                    <items tab="登录" key="login">
                        <Form name="login" onFinish={handleLogin}>
                            <Form.Item
                                name="email"
                                rules={[{ required: true, message: "请输入邮箱！" }, { type: "email", message: "请输入有效的邮箱！" }]}
                            >
                                <Input placeholder="邮箱" />
                            </Form.Item>

                            <Form.Item
                                name="password"
                                rules={[{ required: true, message: "请输入密码！" }]}
                            >
                                <Input.Password placeholder="密码" />
                            </Form.Item>

                            <Form.Item
                                name="captcha"
                                rules={[{ required: true, message: "请输入验证码！" }]}
                            >
                                <div style={{ display: "flex", alignItems: "center" }}>
                                    <Image
                                        src={captchaImage}
                                        alt="验证码"
                                        preview={false}
                                        style={{ width: 120, marginRight: 10, cursor: "pointer" }}
                                        onClick={handleCaptchaClick} // 点击刷新验证码
                                    />
                                    <Input placeholder="验证码" />
                                </div>
                            </Form.Item>

                            <Form.Item>
                                <Button type="primary" htmlType="submit" style={{ width: "100%" }} onClick={handleCaptchaClick}>
                                    登录
                                </Button>
                            </Form.Item>
                        </Form>
                    </items>
                    <items tab="注册" key="register">
                        <Form name="register" onFinish={handleRegister}>
                            <Form.Item
                                name="email"
                                rules={[{ required: true, message: "请输入邮箱！" }, { type: "email", message: "请输入有效的邮箱！" }]}
                            >
                                <Input placeholder="邮箱" />
                            </Form.Item>

                            <Form.Item
                                name="nickname"
                                rules={[{ required: true, message: "请输入昵称！" }]}
                            >
                                <Input placeholder="昵称" />
                            </Form.Item>

                            <Form.Item
                                name="password"
                                rules={[{ required: true, message: "请输入密码！" }]}
                            >
                                <Input.Password placeholder="密码" />
                            </Form.Item>

                            <Form.Item
                                name="confirmPassword"
                                rules={[{ required: true, message: "请确认密码！" }]}
                            >
                                <Input.Password placeholder="确认密码" />
                            </Form.Item>

                            <Form.Item
                                name="captcha"
                                rules={[{ required: true, message: "请输入验证码！" }]}
                            >
                                <div style={{ display: "flex", alignItems: "center" }}>
                                    <Image
                                        src={captchaImage}
                                        alt="验证码"
                                        preview={false}
                                        style={{ width: 120, marginRight: 10, cursor: "pointer" }}
                                        onClick={handleCaptchaClick} // 点击刷新验证码
                                    />
                                    <Input placeholder="验证码" />
                                </div>
                            </Form.Item>

                            <Form.Item>
                                <Button type="primary" htmlType="submit" style={{ width: "100%" }} onClick={handleCaptchaClick}>
                                    注册
                                </Button>
                            </Form.Item>
                        </Form>
                    </items>
                </Tabs>
            </div>
        </div>
    );
};

export default LoginPage;
