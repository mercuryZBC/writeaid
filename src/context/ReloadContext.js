import React, { createContext, useState, useContext } from 'react';

// 创建 Context
const ReloadContext = createContext();

// 创建一个 Provider 组件，父组件通过它来提供值
export const ReloadProvider = ({ children }) => {
    const [kbListReloadTrigger, setkbListReloadTrigger] = useState(false);
    return (
        <ReloadContext.Provider value={{ kbListReloadTrigger, setkbListReloadTrigger }}>
            {children}
        </ReloadContext.Provider>
    );
};

// 创建一个自定义 Hook，便于在子组件中使用
export const useReloadContext = () => {
    return useContext(ReloadContext);
};
