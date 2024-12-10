

// 判断 用户token是否存在
export const getToken = () => {
    const token = localStorage.getItem('jwt_token');
    const expiresIn = localStorage.getItem('expires_in');

    if (!token || !expiresIn) {
        return ""; // 如果没有token或expires_in，直接返回false
    }

    const currentTime = new Date().getTime(); // 获取当前时间
    const expirationTime = parseInt(expiresIn, 10); // 将expires_in转换为整数（秒）
    // 判断token是否过期
    if (currentTime > expirationTime) {
        // 如果过期，删除token和expires_in
        localStorage.removeItem('jwt_token');
        localStorage.removeItem('expires_in');
        return "";
    }
    return token; // 如果没有过期，返回true
};

export const setToken = (access_token,expires_in) => {
    // 将access_token和expires_in存储到localStorage中
    localStorage.setItem('jwt_token', access_token);
    localStorage.setItem('expires_in', expires_in);
}

export const delToken = () => {
    // 将access_token和expires_in存储到localStorage中
    localStorage.removeItem('jwt_token');
    localStorage.removeItem('expires_in');
}


