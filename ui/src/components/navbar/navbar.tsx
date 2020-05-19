import React from 'react';
import {NavLink, useRouteMatch} from "react-router-dom";
import './navbar.css';
import Searchbar from "../searchbar/searchbar";
import { useHeartBeat } from "../../hooks/useHeartBeat";
import {useNewFileCounter} from "../../hooks/useNewFileCounter";

const Navbar = () => {
    const { isConnected } = useHeartBeat();
    const { counter, setCounter } = useNewFileCounter(0);

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
                        <div className='counter'>
                            {counter > 0 ? <span>{counter}</span> : ''}
                            <NavLink onClick={() => setCounter(0)} activeClassName="active" to="/files">My Files</NavLink>
                        </div>
                    </li>
                    <li className='nav-option'>
                        <NavLink activeClassName="active" to="/recent">Downloads</NavLink>
                    </li>
                </ul>
            </nav>
        </>
    )
}

export default Navbar;
