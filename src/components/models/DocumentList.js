import {Button, Drawer, List} from "antd";
import React, {useState} from "react";


const DocumentList = ({ visible, setVisible, handleDocumentSelect }) => {

    return (
        <>
            <Drawer
                title="选择文档"
                width={400}
                onClose={() => setVisible(false)}
                visible={visible}
                bodyStyle={{ paddingBottom: 80 }}
            >
                <List
                    dataSource={documents}
                    renderItem={(doc) => (
                        <List.Item
                            key={doc.id}
                            onClick={() => handleDocumentSelect(doc.id)}
                            style={{ cursor: "pointer" }}
                        >
                            <Button style={{ width: "100%", textAlign: "left" }}>
                                {doc.title}
                            </Button>
                        </List.Item>
                    )}
                />
            </Drawer>
        </>
    );
}
