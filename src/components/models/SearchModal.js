import React, { useState } from "react";
import { Modal, Input, message } from "antd";

const SearchModal = ({ visible, onClose }) => {
    const [searchText, setSearchText] = useState("");

    const handleSearch = () => {
        message.info(`搜索: ${searchText}`);
        onClose();
    };

    return (
        <Modal
            title="搜索"
            visible={visible}
            onCancel={onClose}
            onOk={handleSearch}
        >
            <Input
                placeholder="请输入搜索关键词"
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
            />
        </Modal>
    );
};

export default SearchModal;
