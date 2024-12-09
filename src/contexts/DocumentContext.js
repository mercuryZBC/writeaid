import React, { createContext, useContext, useState } from "react";

const DocumentContext = createContext();

export const useDocumentContext = () => {
    return useContext(DocumentContext);
};

export const DocumentProvider = ({ children }) => {
    const [documents, setDocuments] = useState({});

    const addDocument = (id, title, content, hash) => {
        setDocuments((prev) => ({
            ...prev,
            [id]: { id, title, content, hash },
        }));
    };

    const getDocument = (id) => {
        return documents[id];
    };

    return (
        <DocumentContext.Provider value={{ documents, addDocument, getDocument }}>
            {children}
        </DocumentContext.Provider>
    );
};
