import React, {useState} from 'react';
import './searchbar.css';
import {useHistory} from 'react-router-dom'

export type SearchbarProps = {
    focus?: boolean,
}

const Searchbar = ({focus}: SearchbarProps) => {
    const [focused, setFocused] = useState(focus || false);
    let history = useHistory();

    const onSearch = (e: any) => {
        e.preventDefault();
        history.push("/downloads")
    }

    return (
        <div className={`search-bar ${focused ? 'focus': ''}`}>
            <form onSubmit={onSearch}>
                <span><i className="fas fa-search"></i></span>
                <input onBlur={() => setFocused(false)} onFocus={() => setFocused(true)} autoFocus={focused} placeholder="Search" type="text"/>
            </form>
        </div>
    )
};

export default Searchbar;