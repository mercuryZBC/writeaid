import React, { useState } from "react";
import { Row, Col } from "antd";
import { Editor } from "react-draft-wysiwyg";
import { EditorState, convertToRaw } from "draft-js";
import draftToMarkdown from "draftjs-to-markdown";
import "react-draft-wysiwyg/dist/react-draft-wysiwyg.css";

export const MarkdownEditor = () => {
    const [editorState, setEditorState] = useState(EditorState.createEmpty());

    const handleEditorChange = (state) => {
        setEditorState(state);
        const markdown = draftToMarkdown(convertToRaw(state.getCurrentContent()));
        console.log("Markdown Content: ", markdown);
    };

    return (
        <div style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
            <Row style={{ flex: 1, height: "100%" }}>
                {/* WYSIWYG 编辑器 */}
                <Col span={24} style={{ height: "100%" }}>
                    <Editor
                        editorState={editorState}
                        onEditorStateChange={handleEditorChange}
                        wrapperStyle={{
                            height: "100%",
                            border: "1px solid #ddd",
                            overflow: "hidden",
                        }}
                        editorStyle={{
                            height: "calc(100% - 70px)",
                            padding: "10px",
                            lineHeight: "1.6",
                        }}
                    />
                </Col>
            </Row>
        </div>
    );
};
