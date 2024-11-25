import {Avatar, List, Tabs, Typography,} from "antd";
import React from "react";
import NavigatorBar from "./NavigatorBar";
const { TabPane } = Tabs;
const {Text } = Typography;



export const Overview = () => {
    // 假设的最近文档数据
    const recentData = {
        edited: [
            { avatar: "E", title: "编辑的文档 1", description: "这是最近编辑的文档 1" },
            { avatar: "E", title: "编辑的文档 2", description: "这是最近编辑的文档 2" },
        ],
        viewed: [
            { avatar: "V", title: "浏览的文档 1", description: "这是最近浏览的文档 1" },
            { avatar: "V", title: "浏览的文档 2", description: "这是最近浏览的文档 2" },
        ],
        commented: [
            { avatar: "C", title: "评论的文档 1", description: "这是最近评论的文档 1" },
            { avatar: "C", title: "评论的文档 2", description: "这是最近评论的文档 2" },
        ],
    };

    return (
        <div>
            <NavigatorBar/>
            <Tabs defaultActiveKey="edited" centered>
                {/* 最近编辑 */}
                <TabPane tab="最近编辑" key="edited">
                    <List
                        itemLayout="horizontal"
                        dataSource={recentData.edited}
                        renderItem={(item) => (
                            <List.Item>
                                <List.Item.Meta
                                    avatar={<Avatar>{item.avatar}</Avatar>}
                                    title={<Text>{item.title}</Text>}
                                    description={item.description}
                                />
                            </List.Item>
                        )}
                    />
                </TabPane>

                {/* 最近浏览 */}
                <TabPane tab="最近浏览" key="viewed">
                    <List
                        itemLayout="horizontal"
                        dataSource={recentData.viewed}
                        renderItem={(item) => (
                            <List.Item>
                                <List.Item.Meta
                                    avatar={<Avatar>{item.avatar}</Avatar>}
                                    title={<Text>{item.title}</Text>}
                                    description={item.description}
                                />
                            </List.Item>
                        )}
                    />
                </TabPane>

                {/* 最近评论 */}
                <TabPane tab="最近评论" key="commented">
                    <List
                        itemLayout="horizontal"
                        dataSource={recentData.commented}
                        renderItem={(item) => (
                            <List.Item>
                                <List.Item.Meta
                                    avatar={<Avatar>{item.avatar}</Avatar>}
                                    title={<Text>{item.title}</Text>}
                                    description={item.description}
                                />
                            </List.Item>
                        )}
                    />
                </TabPane>
            </Tabs>
        </div>
    )
}
