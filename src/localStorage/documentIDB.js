import { openDB } from "idb";

// 初始化数据库
const initDB = async () => {
    return openDB("DocumentDB", 1, {
        upgrade(db) {
            if (!db.objectStoreNames.contains("documents")) {
                const store = db.createObjectStore("documents", { keyPath: "doc_id" }); // 使用 doc_id 作为主键
                store.createIndex("doc_id", "doc_id", { unique: true }); // 创建索引
                store.createIndex("kb_d", "kb_d", { unique: false }); // 如果需要根据知识库ID查询，可以创建此索引
            }
        },
    });
};

// 添加文档到数据库
export const addDocumentToDB = async (document) => {
    const db = await initDB();
    await db.put("documents", document); // 保存文档信息
    console.log(getAllDocumentsFromDB());
};

// 获取所有文档
export const getAllDocumentsFromDB = async () => {
    const db = await initDB();
    return db.getAll("documents"); // 获取所有文档
};

// 通过 doc_id 获取指定文档
export const getDocumentByDocIdFromDB = async (docId) => {
    const db = await initDB();
    try {
        // 使用索引查询
        return await db.get("documents", docId); // 通过 doc_id 查找文档
    } catch (err) {
        console.error(err);
        return null;
    }
};

// 清空所有文档
export const clearDocumentsFromDB = async () => {
    const db = await initDB();
    await db.clear("documents"); // 清空所有文档
};
