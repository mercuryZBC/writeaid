import {createContext, useState} from "react";

const UserContext = createContext();

export const UserProvider = ({children}) => {
    const [user, setUser] = useState(null); //email,nickname,avatar

    return (
        <UserContext.Provider value={[user, setUser]}>
            {children}
        </UserContext.Provider>
    );
};


export const setLoginInState = () =>{
    localStorage.setItem('isLoggedIn','true');
}
export const setLogoutState = () =>{
    localStorage.removeItem('isLoggedIn');
}


