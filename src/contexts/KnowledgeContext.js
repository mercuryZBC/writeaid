import React, {createContext, useContext, useState} from "react";

const KnowledgeContext = createContext();

export const useKnowledgeBaseContext = () =>{
    return useContext(KnowledgeContext);
}

export const KnowledgeProvider = ({ children }) => {
    const [knowledgeBases, setKnowledgeBases] = useState({});

    const addKnowledge = (id,name,describe) => {
        setKnowledgeBases((prev) => ({
            ...prev,
            [id]: { id, name,describe },
        }));
    };

    const getKnowledge = (id) =>{
        return knowledgeBases[id];
    };

    return (
        <KnowledgeContext.Provider value={{knowledgeBases,addKnowledge,getKnowledge}}>
            {children}
        </KnowledgeContext.Provider>
    )
}
