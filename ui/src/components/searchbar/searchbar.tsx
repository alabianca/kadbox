import React, {useState} from 'react';
import './searchbar.css';

export type SearchbarProps = {
    focus?: boolean,
}

const Searchbar = ({focus}: SearchbarProps) => {
    const [focused, setFocused] = useState(focus || false);

    return (
        <div className={`search-bar ${focused ? 'focus': ''}`}>
            <span><i className="fas fa-search"></i></span>
            <input onBlur={() => setFocused(false)} onFocus={() => setFocused(true)} autoFocus={focused} placeholder="Search SNFS" type="text"/>
        </div>
    )
};

export default Searchbar;