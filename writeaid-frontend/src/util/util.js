export const formatDate = (dateString) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
    }).format(date);
};

// 计算 SHA-256 哈希值的函数
export const calculateSHA256 = async (content) => {
    const encoder = new TextEncoder(); // 将文本转换为 Uint8Array
    const data = encoder.encode(content); // 将文档内容编码为字节数组

    // 使用 Web Crypto API 计算 SHA-256 哈希
    const hashBuffer = await crypto.subtle.digest("SHA-256", data);

    // 将 ArrayBuffer 转换为十六进制字符串
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

    return hashHex; // 返回计算出来的哈希值
};
