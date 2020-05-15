import React from 'react';
import {NavLink, useRouteMatch} from "react-router-dom";
import './navbar.css';
import Searchbar from "../searchbar/searchbar";
import { useHeartBeat } from "../../hooks/useHeartBeat";

const Navbar = () => {
    const { isConnected } = useHeartBeat();

    return (
        <>
            <nav className='nav' role='navigation'>
                <ul className='main'>
                    <li>
                        <Searchbar focus={true}/>
                    </li>
                    <li>
                        <span>{isConnected ? 'Connected': 'Disconnected'}</span>
                    </li>
                    <li>
                        <button><i className="fas fa-cog"></i></button>
                    </li>
                </ul>
                <ul className='options'>
                    <li className='nav-option'>
                        <NavLink activeClassName="active" to="/share">Share</NavLink>
                    </li>
                    <li className='nav-option'>
                        <NavLink activeClassName="active" to="/files">Files</NavLink>
                    </li>
                    <li className='nav-option'>
                        <NavLink activeClassName="active" to="/recent">Recent</NavLink>
                    </li>
                </ul>
            </nav>
        </>
    )
}

export default Navbar;
